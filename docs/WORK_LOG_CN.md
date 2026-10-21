# CodeAPI 二开工作日志

本文件按天记录 `Xeubad/codeapi` 仓库上的二开工作：当天做了什么、为什么这么做、怎么验证的、产出了哪些提交。流程与分支约定见 [FORK_MAINTENANCE_CN.md](./FORK_MAINTENANCE_CN.md)，Agent 的执行规则见仓库根目录 `AGENTS.md`。

## 记录规范

- **粒度**：一天一节，标题格式 `## 2026年10月20日`，新的一天插在最上面。
- **每节固定四块**：改动 / 关键决策 / 验证 / 产出。改动写做了什么，决策写为什么和放弃了什么，验证写跑了什么命令、结果如何，产出写提交短 SHA。
- **只写已发生的事**：没跑的检查写「未运行」和原因，不能写成「通过」；数字要来自刚跑过的命令。
- **提交号用短 SHA**，与 `git log` 一致，便于追溯。
- **不写流水账**：只记对本 fork 有长期价值的信息（决策、踩坑、风险），纯中间步骤不记。

新的一天用这个模板：

```markdown
## 2026年XX月XX日

**主题**：一句话

### 改动

- ...

### 关键决策

- ...

### 验证

- ...

### 产出

- `<short-sha>` <subject>
```

---

## 2026年10月21日

**主题**：首页与登录页二次重构 —— 苹果式毛玻璃、宣传型版面、真实访问延迟实测

### 改动

- 重写默认宣传首页（`HomeView.vue` 的 promo 分支）：悬浮玻璃顶栏、Hero + 控制台界面预览、事实条、真实延迟面板、能力四宫格、curl 接入、模型条、CTA、纯内链页脚。
- 重写登录外壳 `AuthLayout.vue` 与 `LoginView.vue`：与首页共用背景、玻璃面板、字阶与 ink/lime 色；登录逻辑（验证码、2FA、协议、Passkey、OAuth）保持不变。
- `style.css` 增加共用品牌类：`brand-backdrop` / `brand-orb` / `brand-grain` / `brand-glass(-strong)` / `brand-field` / `brand-ink-panel`。
- 首页新增真实延迟实测：浏览器每 5 秒向本站静态资源发 Range 请求（1 KB），展示当前/平均/P95/抖动/样本数、SVG 折线、TTFB、连接复用、传输量、协议与最近样本；页面隐藏时暂停，重新可见立即补测。
- 移除首页与登录页的所有外链（文档、GitHub、`target="_blank"`），compact 分支同样去掉文档外链；更新 `docUrlSanitization.spec.ts` 断言首页不再渲染外链。
- 更新 zh/en `landing.ts` 的 `home` 文案块，新增导航、Hero、延迟、能力、接入、模型、CTA、页脚与登录安全说明。

### 关键决策

- 延迟数据坚持真实测量而不是伪造多地区数字：用当前浏览器 → 本站静态资源的真实往返，UI 明确标注「当前设备 → 本站静态资源」。多地区合成探针需要后端存储与隐私设计，本轮不做。
- 视图色仍沿用控制台 token：ink `#111620/#0b0e14` + lime `#6f9e12`，lime 只做状态点、细小强调与进度条，不铺大面积背景或主按钮；主按钮继续走 ink/白。
- 保留上游的 `home_content`（iframe/HTML）与 compact 模式分支，避免破坏管理员已配置的首页覆盖能力；本次只重写默认宣传页。
- 背景改成全局共用的 `brand-backdrop` + 柔光团 + 颗粒，登录页直接复用，避免两页各写一套色场导致漂移。
- 延迟请求加 `Range: bytes=0-1023`，把每次探测传输从整张 logo（约 70 KB）降到约 1 KB；采样 5 秒一次，隐藏页暂停。
- 折线 y 轴从 0 起算且最小量程 50 ms，避免样本少时把几毫秒的波动放大成假趋势。
- 登录页只改 UI/布局，不动任何校验与认证流程；新增的「复制」按钮带 `execCommand` 兜底并如实显示失败状态。

### 验证

- `npx pnpm@9 run typecheck` rc=0。
- `npx pnpm@9 run check:i18n` 3 tests 通过。
- `npx pnpm@9 run test:run` 371 文件 / 2868 用例全部通过（含更新后的 `docUrlSanitization.spec.ts`）。
- `npx pnpm@9 run build` 通过，产物写入 `backend/internal/web/dist`。
- 浏览器实测（Vite dev + CDP 设备模拟）：390 / 1024 / 1280 / 1440 宽，浅色与深色均无横向溢出（`scrollWidth === innerWidth`）；首页无任何 `http(s)` 外链与 `target="_blank"`；延迟面板采样数随时间增长、折线真实绘制，传输量约 1.3 KB；主题切换按钮生效；「复制」在受限环境显示「复制失败」，不误报成功。
- 未运行后端测试：本次只改前端视图与文案，未触碰后端。

### 产出

- `af4087b4f` feat(landing): rebuild the promo home and the glass login shell（在 `feature/custom-brand-pages` 上，与首轮品牌页工作同一个 PR）

---

**主题（同日第二次工作）**：后端监听端口锁死 8080，冲突时识别占用进程并在人工确认后终止

### 改动

- 新增 `backend/internal/portguard/`（`portguard.go` + 测试）：`Listen()` 在应用初始化之前绑定并持有监听器。端口被占时打印占用者 PID / 用户 / 启动时间 / 完整命令行；TTY 下要求输入 `y` 才发 `SIGTERM`，随后 5 秒内轮询重绑；非交互环境只打印报告并返回错误。PID ≤ 1 与本进程自身标为 protected，连提问都不出现。
- `backend/internal/config/config.go`：新增 `RequiredServerPort = 8080`，`Config.Validate()` 拒绝任何非 8080 的 `server.port`（`config.yaml` 与 `SERVER_PORT` 都走这条校验）。
- `backend/cmd/server/main.go`：接入 portguard，并把 `app.Server.ListenAndServe()` 改为 `app.Server.Serve(ln)`，消除「绑定失败前就打印 Server started」的竞态与二次 bind。
- `backend/internal/setup/`：CLI 向导不再询问端口，`AUTO_SETUP` 遇到非 8080 的 `SERVER_PORT` 直接报错，写配置统一取 `config.RequiredServerPort`。
- `deploy/install.sh`：移除端口提问与 `validate_port()`，中英两套文案换成「端口已固定」提示。
- 文档：`AGENTS.md` 新增「监听端口硬性条件（8080）」小节；同步 `deploy/README.md`、`deploy/DOCKER.md`、`deploy/config.example.yaml`、`deploy/.env.example` 与三份根 README。顺带纠正 `DOCKER.md` 里并不生效的 `PORT` 变量（viper 只认 `SERVER_PORT`）。

### 关键决策

- 选 fail-fast 校验而不是代码里写死 8080：强制覆盖会让 `SERVER_PORT`/`config.yaml` 变成「改了也不生效」的死字段，排障时误导人；拒绝启动并报出实际值更可预测。
- 占用检测与确认放在 Go 二进制内的独立 package，而不是 deploy 下的 shell 脚本：裸跑二进制、systemd、Docker 三条启动路径行为一致，脚本方案只覆盖安装路径。
- 非交互只报告不杀：无人值守自动 SIGTERM 可能杀掉同机 nginx 或其它业务。代价是 systemd `Restart=always` 下端口被占会每 5 秒重跑一次并刷日志，本次未改 restart 策略。
- 不升级到 `SIGKILL`：`SIGTERM` 后 5 秒仍占用就报错退出，强杀留给人工决定。
- 监听器绑定后由 `Serve(ln)` 持有，避免「探测通过 → 释放 → 重绑」之间的 TOCTOU 窗口。
- 前端 Vite dev 的 3000 与 compose 的宿主机发布端口不受约束：前者是开发服务器，后者只是 `docker -p` 映射，容器内仍是 8080。
- 复用已有依赖 `gopsutil/v4`（net + process）与 `x/term`，`go.mod`/`go.sum` 零改动。
- 本轮起改走 GitHub PR 而不是本地合并后直接 push `main`：三个主题各开一个 PR（① 品牌页与页面重构 ② 端口锁死 8080 ③ 工作日志）。端口分支原先带着品牌页的合并提交，直接开 PR 会把两个主题混进同一份 diff，所以改成从 `origin/main` 重新 cherry-pick 出 `feature/custom-port-8080`——**代价是端口三个提交的 SHA 与最初记录不同**，本节日志按 cherry-pick 后的新 SHA 记录。

### 验证

- `go test ./internal/portguard/ -count=1 -v`：8 PASS / 1 SKIP。含真实子进程持有端口 → 回答 `y` → 断言发出 SIGTERM、端口释放、拿到可用监听器；回答 `no`、以及占用者就是本进程两种情况下均不误杀。
- `go test ./internal/config/ ./internal/setup/ ./internal/server/... ./internal/server/middleware/` 通过（新增 3 个端口校验用例）。
- `go build ./...` 与 `go vet`（4 个受影响 package）干净；按 backend Makefile 同参数 `CGO_ENABLED=0 -trimpath -ldflags` 构建通过。
- `go test ./... -count=1`（backend 全量）exit 0，52 个包 ok，无 FAIL/panic。
- 真二进制端到端（演示前本机 8080 正被 `./bin/server` PID 65977 占用）：非交互 → exit 1，报告准确列出该 PID / 用户 / 启动时间 / 命令行且占用者存活；`SERVER_PORT=3000` → exit 1 并打印 pin 错误；PTY 交互回答 `n` → exit 1，占用者仍存活。
- PTY 交互回答 `y`（用户授权后执行）→ 打印 `Sent SIGTERM to PID 65977.` 与 `Port 8080 is free now; continuing startup.`，随后走到 `repository/ent.go` 的数据库初始化，因配置的 Postgres 指向 127.0.0.1:1 而 `connection refused` 退出，证明端口守卫确实跑在应用初始化之前且监听器成功交接；PID 65977 已被终止。
- 未运行：`make build` 的前端部分（会把手上未提交的 `frontend/` 改动写进 `backend/internal/web/dist`）。
- CI `golangci-lint`（v2.14.0）首轮 fail 两处，都在 `portguard_test.go`：仓库 `.golangci.yml` 开了 `errcheck.check-type-assertions`，`ln.Addr().(*net.TCPAddr)` 这种单值断言会被报；gosec G702 把「测试自 spawn 自身二进制」判成命令注入。修复在 `ddd740649`：断言改 comma-ok，G702 按仓库既有约定逐点 `//nolint:gosec` 并写明信任边界，生产代码的检查一个没放宽。
- **踩坑**：`golangci-lint` 其实装着（`~/go/bin/golangci-lint`，不在 PATH 上），但它是 2.13.2，对 Go 1.27.2 的标准库直接 `export data version 5 is greater than maximum supported version 4`  typecheck 失败，本地跑不出真实结论，只能靠 CI。之前记的「本机未安装该二进制」是错的。
- `git diff --check` 干净。
- 已知不一致但无害：`backend/internal/setup/handler.go` 的 Web 向导仍允许 1-65535 的端口输入。该路由在本 fork 未注册，且它写出的配置会在下次启动被 `Config.Validate()` 拒绝，硬性条件不依赖它。

### 产出

- 分支 `feature/custom-port-8080`（从 `origin/main` 切出，只含端口这一主题），三个提交：
  - `a9efab94e` feat(port): pin the listen port to 8080 and guard it before startup
  - `7faa370c4` feat(deploy): drop the port prompt and document the pinned 8080
  - `76cd53d18` refactor(config): take the viper port default from RequiredServerPort
- PR #2 以 merge commit 合入 `main`（`a7d064804`），CI 反馈的 lint 修复在 `ddd740649`。

---

**主题（同日第三次工作）**：单人维护的落地流程定版为「分支 + PR」，并把本地 `main` 与 `origin` 对齐

### 改动

- 排查分叉：本地 `main` 落后 `origin/main` 12、领先 1（那个本地合并提交 `77d7ae586`）。三条核验证明本地没有独有内容，差异只是缺远端更新。
- 建 `anchor/20261021-align-main` 后把本地 `main` 对齐到 `ce4d32108`；删除 7 条已合并的本地分支（`feature/custom-hard-port-8080`、`feature/custom-port-8080`、`feature/custom-brand-pages`、`feature/custom-ui-tech-mono`、`feature/custom-release-branding`、`feature/custom-work-log`、`docs/custom-work-log-20261021`），保留 `sync/upstream-20261018-b0a02a330` 和锚点分支。
- 顺带把远端多出的 `ddd740649`（portguard 的 CI lint 修复：单值类型断言改 comma-ok，gosec G702 按仓库既有约定 `//nolint:gosec` 并写明信任边界）带回工作副本。
- `AGENTS.md` 新增「上传与合并流程（定版）」六步（开分支 / 提交自证 / 推送开 PR / 等 CI 再合 / 收口 / 高风险动作前），并改写 Repository Identity、General Working Rules、User Intent Mapping 的确认门、Fetch and Analyze 第 6 步、Merge After Confirmation 第 1/3/8/9/10 步、Release and Rollback、Work Log 中所有「本地合并后快进 `main`」的表述；补上此前缺失的 `docs/custom-*` 与新增的 `anchor/*` 两类分支约定。
- `docs/FORK_MAINTENANCE_CN.md` 同步：分支表与流向图改成 PR 驱动、「确认合并」步骤、「推送结果」、新增「日常开发落地」、「失败与恢复」补三条、提交与发布约定。

### 关键决策

- **本地 `main` 降级为只跟随 `origin/main` 的镜像**，不再承载未推送的合并；成果一律走 PR。否决了我先前建议的「文档与小修直推 main」这种例外——单人场景下例外越多，本地与远端越容易再次分叉，今天的 12/1 分叉就是这么来的。
- 合并固定 merge commit：工作日志和确认门都锚在短 SHA 上，squash / rebase 会改写它们；回滚用 `revert -m 1` 且同样走 PR。
- 分支只能从 `origin/main` 切。混入其它主题的分支用 cherry-pick 重建，不用 `git rebase --onto` 或 `git branch -f`——今天实测这类重写动作会被权限层拦，且历史更难自证。
- stacked PR 只作例外，合并前必须先 `gh pr edit <n> --repo Xeubad/codeapi --base main`；PR #3 就是 base 指错才多开了一次收口 PR。
- `git reset --hard`、对齐 `main`、动已发布迁移之前先打 `anchor/*`，事后保留。对齐前必须先跑三条核验证明本地无独有内容。
- 半成品也要 `wip(...)` 提交并推送，工作区不留未提交改动过夜。
- `cla.yml` 的 CLA 检查记为已知噪音，是否移除 workflow 留给用户决定，Agent 不自行删改。

### 验证

- `git diff --diff-filter=A --name-only origin/main main` 为空；`git log --oneline origin/main..main` 只剩 `77d7ae586`；`git diff --stat a8a4bae0e 77d7ae586 -- backend deploy AGENTS.md docs` 为空 → 后端与文档完全一致，那个本地合并的独有差异只在首页旧状态。
- 已删除分支的 tip 全部记录在册（`5435d612f`、`ddd740649`、`af4087b4f`、`2ddc8a1eb`、`f85eb4a5f`、`db482b146`、`7f5eabc21`），其中六个已在 `origin/main` 的历史里。
- 对齐后 `git status --short --branch` 干净，`main == origin/main (ce4d32108)`，`anchor/20261021-align-main` 仍指向 `77d7ae586`。
- `go test ./internal/portguard/ ./internal/config/ ./internal/setup/ -count=1` 三个包全 ok。
- 未运行：完整 `go test ./...` 与 `make build`——本轮只动文档，代码与 `origin/main` 完全一致，没有可验证的新行为。

### 产出

- `docs/custom-solo-pr-flow`（从 `origin/main` 切出）：`docs(fork): freeze the branch-plus-PR landing flow` + 本条工作日志。**经 PR #5 合入，merge commit `8d32780d8`**（补记于同日第四次工作；流程文档的落地记录无法自含自己的合并号，这类补记统一并进下一条日志）。

---

**主题（同日第四次工作）**：README 中文化并改成 fork 视角，删除非中文版本

### 改动

- 三份 README 合并为一份：`git rm README.md README_JA.md` + `git mv README_CN.md README.md`，git 识别为 rename，历史可追溯。
- fork 化：品牌 `Sub2API` → `CodeAPI`（15 处；小写的 `sub2api` 服务名、`/opt/sub2api` 安装目录刻意不动，那是 fork 保留的二进制名与路径）；仓库地址 `Wei-Shaw/sub2api` → `Xeubad/codeapi`（`git clone` ×3、`install.sh` ×2、`docker-deploy.sh` ×1）；logo 从 `assets/logo.svg` 换成 `frontend/public/logo.png`。
- 头部：删掉指向上游的 Trendshift 徽章与三语切换链接，补一行 fork 出处说明（`Sub2API` 与上游仓库链接保留在出处里，出处不品牌化）。
- 结尾：Star History 图表（展示的是上游 star）换成「上游项目」致谢段。
- `docs/FORK_MAINTENANCE_CN.md`：「刻意保留指向上游的位置」第一条翻案，写明理由与遗留项。
- 补记上一条产出：PR #5 / merge `8d32780d8`。

### 关键决策

- **翻案**：原方案把 README 的 curl 安装命令刻意指向上游以压缩冲突面。今天改指 fork——README 是用户照抄命令的地方，上游 `install.sh` 会把官方二进制装进 `/opt/sub2api` 覆盖二开改动，正是 AGENTS.md 明令禁止的路径。冲突收益抵不过这个风险。
- README 只留中文一份：GitHub 首页只自动渲染 `README.md`，保留 `README_CN.md` 会丢首页；日文版没有维护者，删除。
- 品牌替换用大小写精确匹配的 `Sub2API`，不碰小写 `sub2api`（服务名、安装目录、systemd 单元）与社区项目名 `Sub2ApiPay`，避免误伤。
- GPL 许可证、Copyright（Wesley Liddick）与服务条款风险提醒原样保留——那是上游的法律声明，fork 不改。
- `deploy/README.md` 与 `deploy/APPLE_CONTAINER.md` 还有 3+ 条指向 `Wei-Shaw/sub2api` 的安装命令（同一隐患），本次按「只处理三份 README」的既定范围不动，留待决定。
- 日志条目与本次文档改动同 PR 提交：本 PR 本身就是纯文档主题，为一行补记单开 PR 违背「PR 成本花在真实功能上」；此解释只针对文档类 PR，功能 PR 仍不夹日志。

### 验证

- `grep` 核验：`README.md` 内 `Wei-Shaw` 仅剩 2 处（头部出处行 + 上游致谢行），`Sub2API` 仅剩 1 处（出处行），`CodeAPI` 15 处、`Xeubad` 7 处；`assets/logo`、`trendshift`、`star-history`、`README_CN`/`README_JA` 引用清零。
- 端口说法复核：只剩「端口固定为 8080（本二开版本不可修改）」，无自定义端口残留。
- 全仓引用 `README_CN`/`README_JA` 的位置仅剩 `docs/FORK_MAINTENANCE_CN.md` 一处，已一并改掉。
- 未运行：`make test` / `make build`（纯 markdown 改动，不涉及代码与构建）。

### 产出

- `docs/custom-readme-cn`（从 `origin/main` 切出）：README 重写 + 删除两份非中文版 + 维护手册翻案说明 + 本条日志（含 PR #5 补记）。PR 编号与合并 SHA 同样并进下一条日志。

## 2026年10月20日

**主题**：二开收尾 —— UI 风格落地、品牌化、屏蔽 Web 安装向导、修复安装脚本，合入 `main`；随后屏蔽合规承诺与新手引导、换 logo、重构首页与登录页

### 改动

- UI 风格收尾：视图层残余 teal chrome 迁到 `accent` token，清掉残留的品牌色 glow。
- 品牌化：`CodeAPI` 成为后端品牌默认值（`internal/config/config.go` 的 webauthn `rp_display_name`、`totp_service.go`、`auth_service.go`、`setting_parse.go`、`setting_public.go`、`user_service.go`、`content_moderation.go`、`setting_features.go`、`auth_oauth_email_flow.go`、`balance_notify_service.go`、`auth_email_binding.go`），并覆盖 `payment_order.go` 的商品名前缀；前端 i18n 导览文案去 Sub2API；`logo.svg` 改为 ink + lime。
- 屏蔽 Web 安装向导：删除 `backend/cmd/server/main.go` 的 `runSetupServer()` HTTP 入口，未安装时服务直接退出并打印指引（改走 `AUTO_SETUP` 或 `--setup`）；前端移除 `/setup` 路由、启动时 `needs_setup` 重定向与 `getSetupStatus` 调用。
- `deploy/install.sh`：`config.yaml` 不存在时不再自动启动服务，收尾提示改为「先跑 `--setup`，再 `systemctl start`」。
- 修复 flaky 测试：`TestInflightEstimate_AccountMappingNoDBAndBoundedMemory` 的堆增长上限 8MB → 12MB。
- 二开护栏写入 `AGENTS.md` 的「二开范围与安全护栏」。
- 建立按天工作日志 `docs/WORK_LOG_CN.md`，并把补记规则接入 `AGENTS.md`（上游整目录忽略 `docs/*`，需在 `.gitignore` 白名单里放行）。
- 屏蔽管理员合规承诺弹窗：前端移除 `App.vue` 的弹窗挂载与 `router/index.ts` 的进入校验；后端摘掉 `AdminComplianceGuard` 中间件注册（`routes/admin.go`、`routes/payment.go`、`handler/page_handler.go`）。`internal/service/admin_compliance.go` 与 `registerAdminComplianceRoutes` 原样保留。
- 屏蔽新手引导：删除 `AppHeader.vue` 的「重新查看新手引导」下拉入口、`AppLayout.vue` 的 `useOnboardingTour` 实例化，以及 `AppSidebar/KeysView/GroupsView` 里的 `isCurrentStep/nextStep` 钩子。
- 换平台 logo：新增 `frontend/public/logo.png`（256×256 RGBA，70893 B），删除 `frontend/public/logo.svg`；同步 `index.html` favicon、侧边栏、顶栏、首页、登录页、模型广场、法务页、用法页的全部引用，并把 `backend/internal/web/embed_test.go` 的静态资源断言从 `/logo.svg` 换成 `/logo.png`。
- 首页重构为宣传页：默认首页（`data-testid="promo-home"`）改成「暗色控制台预览 + 编号能力 + 对比表 + 供应商标签」的编辑式版面，色调取自控制台的 ink + lime；保留上游的自定义首页（iframe/HTML）与 compact 模式。
- 登录页重构：新增 `AuthLayout.vue` 毛玻璃外壳（`backdrop-blur-2xl` + `backdrop-saturate-150` + 高光描边 + 长投影），与首页共用同一套色场、颗粒纹理与 headline。
- 修掉首页/登录页背景色场不显示的问题：背景层原先是 `-z-10`，会被根元素不透明底色盖住（截图像素实测整片 `#f6f7f9`），改成 `z-0` 并把 `main`/`footer` 提到 `relative z-10`。
- 修掉标题在窄屏把词拆断的问题：headline 拆成 `home.heroTitleLead` / `home.heroTitleMain` 两个 key（zh「一个密钥，」/「畅用多个 AI 模型」），不再依赖 `text-balance` 的自动断行。
- 修掉 `router/__tests__/title.spec.ts` 的过期断言：仍期望默认站点名 `Sub2API`，与 `a9df0609d` 之后的 `CodeAPI` 默认值不符（该失败在本次改动前的 HEAD 上已存在）。

### 关键决策

- **安全护栏成文**：原平台自带的校验、安全检测、鉴权、限流、计费一致性校验与审计逻辑不得删除、绕过或弱化；需要改动时只能在原实现基础上增强并补测试，先报告「为什么必须改、替代校验是什么」再动手。
- flaky 断言选择**放宽阈值**而不是删测试或改实现：该断言在干净主干上同样会因 MemStats 的 GC 噪声在 8MB 附近随机越界（实测 8.63 / 8.68 / 8.74MB），12MB 留出余量又不至于放过真正的按模型名缓存。
- 屏蔽 setup 采用**移除入口、保留 `internal/setup` 包**：容器走 `AUTO_SETUP=true` + 环境变量，systemd 走 `sudo -u sub2api ./sub2api --setup`，两条终端路径都保留，只是不再在主端口上开 Web 向导。
- 继续保留指向上游的位置（README、`org.opencontainers.image.source` 之外的镜像 LABEL、法务与支付文档链接、首页与 Header 的 GitHub 外链），目的是压低与上游的冲突面。
- 合规承诺**只摘闸门、不删留痕**：上游 `admin_compliance` 的服务、路由与落库（版本/时间/账号/IP/UA）全部保留，只是不再拦截管理员请求。这套逻辑属于上游自带的审计能力，按二开护栏不能删除或弱化；触发 423 的中间件注册被摘掉，管理员登录后不再被弹窗打断。
- 新手引导**掐入口、不删文件**：driver.js、`useOnboardingTour.ts`、`stores/onboarding.ts`、`styles/onboarding.css` 全部留在仓库里，只是取消自动播放与重播入口，后续同步上游时冲突面最小。
- 首页默认走宣传页，但**保留上游自定义首页能力**：`home_content`（iframe/HTML）与 compact 模式的判断逻辑沿用原样，管理员在设置里写入 `home_content` 仍会覆盖宣传页。
- 背景层用 `z-0` + 内容 `z-10`，而不是负 z-index：固定层挂 `-z-10` 时会被根元素的不透明底色盖住，实测截图像素整片是 `#f6f7f9`，色场根本没渲染出来。
- 视觉语言沿用控制台既有 token（ink `#0b0e14`、lime `#6f9e12`、hairline 分隔线、等宽控制台预览），不另起一套配色；浅色下的背景色场保持低饱和，避免变成彩虹渐变。

### 验证

- 本地容器（`postgres:18-alpine` + `redis:8-alpine`）跑通完整栈；`AUTO_SETUP` 安装后 `/health` 返回 200。
- 无安装锁时服务不监听任何端口，直接退出并打印指引。
- `go test -tags=unit ./internal/service/ -count=1` 全包通过（201.1s）。
- `main` 分支 CI 5/5 成功：`shell` / `release-helpers` / `frontend` / `test` / `golangci-lint`（run `39426046924`，headSha `52e7eaf9e`）。
- 前端全量检查：`npx pnpm@9 run typecheck` rc=0；`npx pnpm@9 run check:i18n` 3 tests 通过；`npx pnpm@9 run test:run` 371 文件 / 2868 用例全绿（其中修掉 1 个 HEAD 上已存在的失败）；`npx pnpm@9 run build` 通过，产物落到 `backend/internal/web/dist`。
- 后端：`go build ./...` rc=0；`go build -tags=embed -o bin/server ./cmd/server` rc=0；`go test -tags=unit -count=1 ./internal/web/... ./internal/server/...` 全通过（web / server / middleware / routes）。
- 浏览器实测（`127.0.0.1:8080`，管理员 `admin@codeapi.local`）：登录后直接进 `/admin/dashboard`，页面无合规弹窗、DOM 里搜不到「合规承诺」字样；顶栏用户下拉只剩「个人资料 / API 密钥 / GitHub / 退出登录」，没有「重新查看新手引导」，也没有 driver.js 遮罩；侧边栏与顶栏 logo 都是 `/logo.png`。
- 视觉走查：headless Chrome 截 `/home`、`/login` 的明暗两版；移动端改用 CDP `Emulation.setDeviceMetricsOverride` 以 390×844 真机宽度复核（Chrome 桌面版最小窗口宽 500px，直接改 `--window-size` 截图只会被裁切而不是重排）。

### 产出

- `1a21e0e8c` refactor(ui): move remaining teal chrome onto the accent token
- `a9df0609d` feat(brand): rebrand user-facing surfaces to CodeAPI
- `958aef0a1` feat(brand): make CodeAPI the backend brand default
- `a81ec2f96` fix(i18n): drop Sub2API from onboarding tour copy
- `3cecf719a` test(server): expect the new brand default in the settings contract
- `7915179fe` feat(setup): drop the HTTP wizard from the fork
- `78a4d2971` fix(install.sh): stop starting the service before it is configured
- `eb9417f4c` test(billing): widen the inflight-cache ceiling to absorb GC noise
- `2ddc8a1eb` docs(fork): pin the fork's scope and the no-weakening rule
- `52e7eaf9e` merge: 二开 UI 风格重做与首轮功能定制 (feature/custom-ui-tech-mono)
- `9917d4239` feat(admin): unhook the upstream compliance acknowledgement gate
- `2ec9f0e4d` feat(ui): disable the onboarding tour and its replay entry
- `59e3d5376` feat(brand): swap the platform logo and rebuild the home and login pages

---

## 2026年10月19日

**主题**：首次上游同步（653 提交）、分支模型收敛为单主干、UI 风格重做起步

### 改动

- 拉取并分析 `upstream/main`（653 个提交），在 `sync/upstream-20261018-b0a02a330` 上执行 `merge --no-ff`，随后把 `main` 快进到同步结果。
- 合并后重放二开的品牌化与文档提交。
- 分支模型收敛：`main` 成为唯一长期分支，取消 `develop`（本 fork 单人维护、无预发环境，两条长期分支只会在每次同步时多制造一次人工提升）。
- UI 风格重做（token 层）：主色由 teal 改为 lime（`primary.500 #6f9e12`，`700 #4d7c0f`），中性色统一为冷灰 ink 阶（`gray.50 #f6f7f9` → `gray.950 #0b0e14`），圆角改小（`lg 7px` / `xl 8px` / `2xl 10px`），按钮拍平为实心 ink（`.btn-primary` = `bg-gray-900` / 暗色 `bg-gray-100`），去掉彩色 glow。

### 关键决策

- `.gitignore` 的唯一冲突按**并集**解决，不选 ours/theirs：上游删掉了 `!docs/ANTIGRAVITY_ATTRIBUTION_429.md` 白名单（因为它同时删了那个文件），本 fork 在同一处新增了 `!docs/FORK_MAINTENANCE_CN.md`。两边意图都保留。
- 同步采用 `merge --no-ff` 而非 rebase，保留上游原始历史便于日后追溯。
- 风格方向定为「明亮 SaaS + 克制科技感」：换掉整体色调（不只换主色），但保留 DeepSeek teal 与全部支付品牌色、语义色，避免撞色。
- 支付类按钮保留品牌色，只去掉品牌色 glow。

### 验证

- 同步后 `main` 相对 `upstream/main` 为 ahead 19 / behind 0，基线完全包含上游。
- 上游带来的变更覆盖：Go 1.27.0 → 1.27.2、`golang.org/x/net` v0.58.0 → v0.60.0、axios ^1.18 → ^1.20、vue ^3.4 → ^3.5、`source-map-js` 提升；`xlsx` 补安全公告豁免说明（CVE-2023-30533 / CVE-2024-22363，版本仍是 ^0.18.5，仓库只做写不读）；数据库迁移 240–245；4 个 Ent schema 变更（`composite_model_route.go`、`group.go`、`payment_order.go`、`user_platform_quota.go`）。

### 产出

- `5ca23d628` chore(sync): merge upstream at b0a02a330 and reapply fork branding
- `f85eb4a5f` docs(fork): collapse to a single long-lived branch
- `4eeb14ec5` docs(fork): name both DockerHub literals in release.yml
- `d14abfe51` feat(ui): retoken palette to lime accent on a unified neutral scale
- `c839fb9e7` refactor(ui): flatten buttons to solid ink and drop coloured glows
- `0ec5197f7` fix(ui): drop brand-coloured glows and teal chrome left in views

---

## 2026年10月18日

**主题**：fork 落地 —— 仓库改名 `codeapi`、发布链路品牌化、二开流程成文

### 改动

- GitHub 仓库 `Xeubad/sub2api` 改名为 `Xeubad/codeapi`，本地 `origin` 同步改名。
- 远端与 Git 配置：`origin` = `Xeubad/codeapi`，`upstream` = `Wei-Shaw/sub2api`（只读，push 地址置为 `DISABLED`，让误推在 git 层直接失败）；启用 `rerere.enabled` 与 `merge.conflictStyle zdiff3`。
- 发布链路：镜像名改为跟随仓库名（`.github/release-tools/release-images.sh` 取 `${GITHUB_REPOSITORY##*/}` 并小写），产物即 `ghcr.io/xeubad/codeapi`；`.goreleaser.yaml` 的 `name_template` 改为 `CodeAPI {{.Version}}`，正文 `docker pull` 文案改引用 `GITHUB_REPO_NAME`，与镜像名同源。
- 应用内更新：后端 `backend/internal/service/update_service.go` 的 `githubRepo` 改为 `Xeubad/codeapi`；前端 `frontend/src/components/common/VersionBadge.vue` 的 `GITHUB_REPO` 与 `DOCKER_IMAGE` 同步。
- 部署脚本：`deploy/install.sh` 与 `deploy/docker-deploy.sh` 的仓库地址改指自己。
- 文档：新增 `AGENTS.md`（二开执行规则）与 `docs/FORK_MAINTENANCE_CN.md`（流程手册）。

### 关键决策

- 更新器指向自己仓库并**保留 `-custom.N` tag 方案**。已知后果：`update_service.go` 的 `parseVersion` 丢弃 `-` 之后的一切且只比较三段数字，同号二开 tag 永远判定为「无更新」，后台的更新/回滚按钮实际不会亮。因此**一键更新不作为升级通道**，升级一律走 Git + 自行发版。
- 不改 README，只写自己的部署文档 —— README 属于上游内容，改动会持续制造合并冲突。
- 二进制名、systemd 服务名、`INSTALL_DIR` 仍为 `sub2api` / `/opt/sub2api`：它们被 `Dockerfile.goreleaser` 的 `COPY sub2api /app/sub2api`、`release_matrix.py` 的归档名与 tar 成员名、`test_release_matrix.py` 的 fixture、`deploy/install.sh` 共同依赖。
- `.goreleaser.yaml` 的 `dockers` / `docker_manifests` 两段是**死配置**（`release_matrix.py` 生成时置空，CI 的 goreleaser 步骤带 `--skip=docker`），其中残留的 `sub2api` 字面量不产生镜像，不为品牌化去改那十几处。
- 应用内检查更新的仓库必须指自己 —— 指回上游会让「一键更新」用官方二进制覆盖本地改动。

### 验证

- 改名后 `gh repo view` 确认仓库为 `Xeubad/codeapi`。
- 确认镜像名与 Release 通知里的 `GHCR_IMAGE` 同源：两者都取仓库名，改仓库名不需要再改代码。
- Docker Hub 未启用（仓库 Actions secrets 为空），`release-images.sh` 与 goreleaser 的 DockerHub 分支恒 skip；GHCR 侧自洽。

### 产出

- `8a164492d` feat(release): derive image name from repository and brand releases as CodeAPI
- `33ff7e81e` fix(update): point release consumers at the fork repository
- `0aab31cf4` docs(fork): record agent rules, sync workflow and release naming
