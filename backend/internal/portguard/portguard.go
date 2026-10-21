// Package portguard binds the pinned listen port and, when the port is taken,
// reports which process holds it and asks before terminating it.
package portguard

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	gonet "github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
	"golang.org/x/term"
)

const (
	// After SIGTERM the port is polled this long before giving up. Escalating to
	// SIGKILL is left to the operator: a hung process is a decision for a human.
	termGrace       = 5 * time.Second
	freeRetryPeriod = 200 * time.Millisecond
)

// UI carries the launching terminal. Interactive is false when stdin is not a
// terminal (systemd, docker, piped scripts), so no unattended run can kill a process.
type UI struct {
	In          io.Reader
	Out         io.Writer
	Interactive bool
}

func NewUI(stdin *os.File, out io.Writer) *UI {
	return &UI{
		In:          stdin,
		Out:         out,
		Interactive: stdin != nil && out != nil && term.IsTerminal(int(stdin.Fd())),
	}
}

// Holder is one process listening on the target port.
type Holder struct {
	PID     int32
	Name    string
	Cmdline string
	User    string
	Started time.Time
	// Protected is a non-empty reason this process must never be terminated.
	Protected string
}

// Listen binds tcp host:port and returns the held listener, so the caller can
// Serve it without a second bind race. On a busy port it prints an occupancy
// report; in a terminal it offers to terminate the holders first.
func Listen(ctx context.Context, host string, port int, ui *UI) (net.Listener, error) {
	addr := net.JoinHostPort(host, strconv.Itoa(port))

	ln, err := net.Listen("tcp", addr)
	if err == nil {
		return ln, nil
	}
	if !isAddressInUse(err) {
		return nil, fmt.Errorf("listen on %s: %w", addr, err)
	}

	holders, lookupErr := findHolders(ctx, port, host)
	writeOccupancyReport(ui.Out, addr, port, holders, lookupErr)

	if !ui.Interactive {
		return nil, fmt.Errorf("port %d is busy and no terminal is attached to confirm: free it manually or start this service from a terminal", port)
	}
	if !confirmTerminate(ui, holders, addr, port) {
		return nil, fmt.Errorf("port %d is still held by another process", port)
	}

	for _, killErr := range terminateHolders(ui.Out, holders) {
		fmt.Fprintf(ui.Out, "  %v\n", killErr)
	}

	ln, err = waitForFreePort(ctx, addr, termGrace)
	if err != nil {
		return nil, fmt.Errorf("port %d is still busy after SIGTERM: %w", port, err)
	}
	fmt.Fprintf(ui.Out, "Port %d is free now; continuing startup.\n", port)
	return ln, nil
}

func findHolders(ctx context.Context, port int, host string) ([]Holder, error) {
	conns, err := gonet.ConnectionsWithContext(ctx, "tcp")
	if err != nil {
		return nil, err
	}

	seen := make(map[int32]bool, len(conns))
	holders := make([]Holder, 0, len(conns))
	for _, conn := range conns {
		if int(conn.Laddr.Port) != port || !strings.EqualFold(conn.Status, "LISTEN") {
			continue
		}
		if !overlaps(conn.Laddr.IP, host) || conn.Pid <= 0 || seen[conn.Pid] {
			continue
		}
		seen[conn.Pid] = true
		holders = append(holders, describe(ctx, conn.Pid))
	}
	return holders, nil
}

// overlaps reports whether a listener bound to listenerIP can block a bind on host.
func overlaps(listenerIP, host string) bool {
	if isUnspecified(host) || isUnspecified(listenerIP) {
		return true
	}
	return listenerIP == host
}

func isUnspecified(host string) bool {
	return host == "" || host == "0.0.0.0" || host == "::"
}

func describe(ctx context.Context, pid int32) Holder {
	holder := Holder{PID: pid}
	if proc, err := process.NewProcessWithContext(ctx, pid); err == nil {
		if name, err := proc.NameWithContext(ctx); err == nil {
			holder.Name = name
		}
		if cmdline, err := proc.CmdlineWithContext(ctx); err == nil {
			holder.Cmdline = cmdline
		}
		if user, err := proc.UsernameWithContext(ctx); err == nil {
			holder.User = user
		}
		if created, err := proc.CreateTimeWithContext(ctx); err == nil {
			holder.Started = time.UnixMilli(created)
		}
	}
	switch {
	case pid <= 1:
		holder.Protected = "init process"
	case int(pid) == os.Getpid():
		holder.Protected = "this process itself"
	}
	return holder
}

func writeOccupancyReport(out io.Writer, addr string, port int, holders []Holder, lookupErr error) {
	if out == nil {
		return
	}
	fmt.Fprintf(out, "Port %d is already in use (bind %s: address already in use).\n", port, addr)
	if lookupErr != nil {
		fmt.Fprintf(out, "  Process lookup failed: %v\n", lookupErr)
	}
	if len(holders) == 0 {
		fmt.Fprintf(out, "  No owning process could be identified; run `lsof -nP -iTCP:%d -sTCP:LISTEN` or `ss -lptn 'sport = :%d'` as root.\n", port, port)
		return
	}
	fmt.Fprintf(out, "  %-8s %-10s %-20s %s\n", "PID", "USER", "STARTED", "COMMAND")
	for _, holder := range holders {
		started := "-"
		if !holder.Started.IsZero() {
			started = holder.Started.Format("2006-01-02 15:04:05")
		}
		command := holder.Cmdline
		if command == "" {
			command = holder.Name
		}
		if command == "" {
			command = "(unknown)"
		}
		user := holder.User
		if user == "" {
			user = "-"
		}
		fmt.Fprintf(out, "  %-8d %-10s %-20s %s\n", holder.PID, user, started, command)
		if holder.Protected != "" {
			fmt.Fprintf(out, "           protected: %s, will not be terminated\n", holder.Protected)
		}
	}
}

func confirmTerminate(ui *UI, holders []Holder, addr string, port int) bool {
	targets := make([]string, 0, len(holders))
	for _, holder := range holders {
		if holder.Protected != "" {
			continue
		}
		label := strconv.FormatInt(int64(holder.PID), 10)
		if holder.Name != "" {
			label += " (" + holder.Name + ")"
		}
		targets = append(targets, label)
	}
	if len(targets) == 0 {
		fmt.Fprintln(ui.Out, "  Nothing is safe to terminate for this port.")
		return false
	}
	fmt.Fprintf(ui.Out, "Terminate %s with SIGTERM so this service can bind %s (port %d)? [y/N]: ",
		strings.Join(targets, ", "), addr, port)

	answer, err := readAnswer(ui.In)
	if err != nil {
		return false
	}
	switch strings.ToLower(answer) {
	case "y", "yes":
		return true
	default:
		return false
	}
}

func readAnswer(in io.Reader) (string, error) {
	if in == nil {
		return "", errors.New("no input stream")
	}
	line, err := bufio.NewReader(in).ReadString('\n')
	if err != nil && line == "" {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func terminateHolders(out io.Writer, holders []Holder) []error {
	var errs []error
	for _, holder := range holders {
		if holder.Protected != "" {
			errs = append(errs, fmt.Errorf("refusing to terminate PID %d: %s", holder.PID, holder.Protected))
			continue
		}
		proc, err := os.FindProcess(int(holder.PID))
		if err != nil {
			errs = append(errs, fmt.Errorf("PID %d not found: %w", holder.PID, err))
			continue
		}
		if err := proc.Signal(syscall.SIGTERM); err != nil {
			errs = append(errs, fmt.Errorf("SIGTERM PID %d failed: %w (another user's process needs root)", holder.PID, err))
			continue
		}
		fmt.Fprintf(out, "Sent SIGTERM to PID %d.\n", holder.PID)
	}
	return errs
}

func waitForFreePort(ctx context.Context, addr string, timeout time.Duration) (net.Listener, error) {
	deadline := time.Now().Add(timeout)
	var lastErr error
	for {
		ln, err := net.Listen("tcp", addr)
		if err == nil {
			return ln, nil
		}
		if !isAddressInUse(err) {
			return nil, err
		}
		lastErr = err
		if time.Now().After(deadline) {
			return nil, lastErr
		}
		timer := time.NewTimer(freeRetryPeriod)
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil, ctx.Err()
		case <-timer.C:
		}
	}
}

func isAddressInUse(err error) bool {
	if err == nil {
		return false
	}
	var errno syscall.Errno
	if errors.As(err, &errno) {
		return errno == syscall.EADDRINUSE
	}
	// Windows surfaces the same condition as a WSAEADDRINUSE message, not a syscall.Errno.
	return strings.Contains(strings.ToLower(err.Error()), "address already in use")
}
