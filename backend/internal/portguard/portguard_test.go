package portguard

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"
)

const (
	childEnvKey  = "PORTGUARD_TEST_CHILD_ADDR"
	childTestRun = "TestHelperHoldsPort"
)

// freePort claims an ephemeral port and releases it so the test can bind it again.
func freePort(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("reserve free port: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	if err := ln.Close(); err != nil {
		t.Fatalf("release free port: %v", err)
	}
	return port
}

func addrOf(port int) string {
	return net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
}

func assertStillHeld(t *testing.T, port int, why string) {
	t.Helper()
	probe, err := net.DialTimeout("tcp", addrOf(port), 500*time.Millisecond)
	if err != nil {
		t.Fatalf("%s: port %d is no longer held: %v", why, port, err)
	}
	_ = probe.Close()
}

func TestListenReturnsHeldListenerWhenPortIsFree(t *testing.T) {
	port := freePort(t)
	var out bytes.Buffer

	ln, err := Listen(context.Background(), "127.0.0.1", port, &UI{Out: &out})
	if err != nil {
		t.Fatalf("Listen() error = %v, output = %q", err, out.String())
	}
	defer func() { _ = ln.Close() }()

	if ln.Addr().String() != addrOf(port) {
		t.Fatalf("listener addr = %q, want %q", ln.Addr().String(), addrOf(port))
	}
	if out.Len() != 0 {
		t.Fatalf("expected a silent success, got %q", out.String())
	}
}

func TestListenNonInteractiveReportsHolderAndNeverKills(t *testing.T) {
	port := freePort(t)
	ln, err := net.Listen("tcp", addrOf(port))
	if err != nil {
		t.Fatalf("occupy port: %v", err)
	}
	defer func() { _ = ln.Close() }()

	var out bytes.Buffer
	got, err := Listen(context.Background(), "127.0.0.1", port, &UI{Out: &out})
	if err == nil {
		t.Fatalf("Listen() succeeded while port %d was occupied", port)
	}
	if got != nil {
		t.Fatalf("Listen() returned a listener on failure: %v", got)
	}

	report := out.String()
	for _, want := range []string{
		fmt.Sprintf("Port %d is already in use", port),
		strconv.Itoa(os.Getpid()), // the holder is this very test process
		"this process itself",     // and it is marked unkillable
	} {
		if !strings.Contains(report, want) {
			t.Fatalf("occupancy report missing %q:\n%s", want, report)
		}
	}
	if !strings.Contains(err.Error(), "no terminal is attached") {
		t.Fatalf("error should explain the non-interactive refusal, got: %v", err)
	}
	assertStillHeld(t, port, "non-interactive start must not kill the holder")
}

func TestListenAnswerNoKeepsHolderAlive(t *testing.T) {
	port := freePort(t)
	child := startPortHolder(t, port)

	var out bytes.Buffer
	ui := &UI{In: strings.NewReader("no\n"), Out: &out, Interactive: true}
	ln, err := Listen(context.Background(), "127.0.0.1", port, ui)
	if err == nil {
		t.Fatalf("Listen() succeeded after declining the prompt")
	}
	if ln != nil {
		_ = ln.Close()
		t.Fatal("Listen() returned a listener after declining")
	}
	report := out.String()
	if !strings.Contains(report, "[y/N]") {
		t.Fatalf("expected a confirmation prompt, got:\n%s", report)
	}
	if !strings.Contains(report, fmt.Sprintf("PID %d", child.Process.Pid)) &&
		!strings.Contains(report, strconv.Itoa(child.Process.Pid)) {
		t.Fatalf("prompt should name the holder PID, got:\n%s", report)
	}
	if strings.Contains(report, "Sent SIGTERM") {
		t.Fatalf("declining must not signal anything, got:\n%s", report)
	}
	if !strings.Contains(err.Error(), "still held") {
		t.Fatalf("error should say the port is still held, got: %v", err)
	}
	assertStillHeld(t, port, "declining the prompt must not kill the holder")
}

// Answering yes must still not let the guard signal init or its own process.
func TestListenInteractiveRefusesProtectedHolder(t *testing.T) {
	port := freePort(t)
	ln, err := net.Listen("tcp", addrOf(port))
	if err != nil {
		t.Fatalf("occupy port: %v", err)
	}
	defer func() { _ = ln.Close() }()

	var out bytes.Buffer
	ui := &UI{In: strings.NewReader("y\n"), Out: &out, Interactive: true}
	if _, err := Listen(context.Background(), "127.0.0.1", port, ui); err == nil {
		t.Fatalf("Listen() succeeded while only a protected holder occupied the port")
	}
	report := out.String()
	if !strings.Contains(report, "Nothing is safe to terminate") {
		t.Fatalf("expected a protected-holder refusal, got:\n%s", report)
	}
	if strings.Contains(report, "[y/N]") {
		t.Fatalf("a protected holder should not even be offered for killing, got:\n%s", report)
	}
	assertStillHeld(t, port, "a protected holder must survive")
}

func TestListenTerminatesHolderAfterConfirmation(t *testing.T) {
	port := freePort(t)
	child := startPortHolder(t, port)

	var out bytes.Buffer
	ui := &UI{In: strings.NewReader("y\n"), Out: &out, Interactive: true}
	ln, err := Listen(context.Background(), "127.0.0.1", port, ui)
	if err != nil {
		t.Fatalf("Listen() after confirmation error = %v, output =\n%s", err, out.String())
	}
	defer func() { _ = ln.Close() }()

	report := out.String()
	if !strings.Contains(report, fmt.Sprintf("Sent SIGTERM to PID %d", child.Process.Pid)) {
		t.Fatalf("expected SIGTERM for the holder PID, got:\n%s", report)
	}
	if !strings.Contains(report, "is free now") {
		t.Fatalf("expected a freed-port line, got:\n%s", report)
	}

	waitErr := child.Wait()
	if waitErr == nil {
		t.Fatalf("child exited cleanly although it was terminated")
	}
	var exitErr *exec.ExitError
	if !errors.As(waitErr, &exitErr) {
		t.Fatalf("child.Wait() = %v, want a signaled exit", waitErr)
	}
}

// TestHelperHoldsPort is only meaningful as a child process (see startPortHolder).
func TestHelperHoldsPort(t *testing.T) {
	addr := os.Getenv(childEnvKey)
	if addr == "" {
		t.Skip("not a port-holding child")
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "child listen: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = ln.Close() }()
	time.Sleep(2 * time.Minute)
}

// startPortHolder runs a child that listens on 127.0.0.1:port and dies on SIGTERM.
func startPortHolder(t *testing.T, port int) *exec.Cmd {
	t.Helper()
	child := exec.Command(os.Args[0], "-test.run="+childTestRun, "-test.timeout=60s")
	child.Env = append(os.Environ(), childEnvKey+"="+addrOf(port))
	child.Stdout = os.Stderr
	child.Stderr = os.Stderr
	if err := child.Start(); err != nil {
		t.Fatalf("start child holder: %v", err)
	}
	// The test body reaps the child; this only covers an early t.Fatalf.
	t.Cleanup(func() { _ = child.Process.Kill() })

	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		probe, err := net.DialTimeout("tcp", addrOf(port), 200*time.Millisecond)
		if err == nil {
			_ = probe.Close()
			return child
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("child never listened on %s", addrOf(port))
	return nil
}

func TestTerminateHoldersRefusesProtectedProcesses(t *testing.T) {
	var out bytes.Buffer
	errs := terminateHolders(&out, []Holder{
		{PID: 1, Protected: "init process"},
		{PID: int32(os.Getpid()), Name: "portguard.test", Protected: "this process itself"},
	})
	if len(errs) != 2 {
		t.Fatalf("expected two refusals, got %v", errs)
	}
	for _, err := range errs {
		if !strings.Contains(err.Error(), "refusing to terminate") {
			t.Fatalf("expected refusal wording, got: %v", err)
		}
	}
	if out.Len() != 0 {
		t.Fatalf("nothing should be signalled, output = %q", out.String())
	}
}

func TestOverlapsBlocksBind(t *testing.T) {
	cases := []struct {
		listenerIP, host string
		want             bool
	}{
		{"0.0.0.0", "0.0.0.0", true},
		{"0.0.0.0", "127.0.0.1", true},
		{"127.0.0.1", "0.0.0.0", true},
		{"127.0.0.1", "127.0.0.1", true},
		{"192.0.2.1", "127.0.0.1", false},
		{"::", "", true},
	}
	for _, tc := range cases {
		if got := overlaps(tc.listenerIP, tc.host); got != tc.want {
			t.Errorf("overlaps(%q, %q) = %v, want %v", tc.listenerIP, tc.host, got, tc.want)
		}
	}
}

func TestIsAddressInUse(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer func() { _ = ln.Close() }()

	_, err = net.Listen("tcp", ln.Addr().String())
	if err == nil {
		t.Fatal("expected the duplicate bind to fail")
	}
	if !isAddressInUse(err) {
		t.Fatalf("isAddressInUse(%v) = false, want true", err)
	}
	if isAddressInUse(nil) {
		t.Fatal("isAddressInUse(nil) = true")
	}
	if isAddressInUse(net.ErrClosed) {
		t.Fatal("isAddressInUse(net.ErrClosed) = true")
	}
}
