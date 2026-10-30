# Sub2API Fork Agent Rules

本仓库是 `Wei-Shaw/sub2api` 的长期二次开发 fork，**已完全脱离上游节奏**：本 fork 按自己的产品意图做完整二开创作，上游只是一个可选的功能参考库——看到值得借鉴的内容就分析并合并进来，不值得的直接忽略，不做「跟不跟得上上游」的自我要求。

详细流程见 `docs/FORK_MAINTENANCE_CN.md`。本文是执行规则；发生冲突时，本文优先。

## Repository Identity

- 官方上游远端名：`upstream`
- 官方上游地址：`https://github.com/Wei-Shaw/sub2api.git`
- 当前 fork 远端名：`origin`
- 当前 fork 地址：`https://github.com/Xeubad/codeapi.git`（原 `Xeubad/sub2api`）
- 二开发布产物：镜像名跟随仓库名 → `ghcr.io/xeubad/codeapi`；二进制名、systemd 服务名、`INSTALL_DIR` 仍为 `sub2api` / `/opt/sub2api`
- 应用内检查更新的仓库：`Xeubad/codeapi`（禁止指回上游，否则会用官方二进制覆盖本地改动）
- 上游参考镜像：`upstream/main`（只读引用，随时可比对，但不是基线也不是义务）
- **上游更新模式（2026-10-31 起）：完全解耦 + 按需借鉴。** 不再逐提交评估、不再有「已评估到哪个 tag」的欠账概念、不执行 `git merge upstream/main`、`main` 不与 `upstream/main` 保持祖先关系。看上游只在「想知道官方最近做了什么」或「某个功能值得借鉴」时进行：值得借鉴的→分析后合并或自行实现；不值得的→直接不管。`git fetch --prune upstream` 仍是唯一允许的上游写入动作（只更新只读引用）
- 借鉴记录：`docs/UPSTREAM_DECISIONS_CN.md` 从 2026-10-31 起由「逐提交决策台账」降级为「借鉴来源记录」——只记我们**实际拿进来**的东西（来源提交/功能、以何种方式落地、改了哪些文件、为何这样整合）。不再要求记录拒绝项，也不再用于划定评估边界
- 二开主干：`main`，唯一长期分支，同时承担集成与生产。本 fork 单人维护、无预发环境，因此不设 `develop`
- **`origin/main` 是唯一真相**：本地 `main` 紧跟 `origin/main`，改动提交后立即推送，不长期停留本地分叉
- 常规二开工作直接在 `main` 上提交并推送（见「提交与推送流程」），不建功能分支。确需隔离的临时改动才建 `feature/custom-*` 分支，事后删除本地与远端分支
- 借鉴工作分支（可选）：`borrow/<主题>`。只有当一次借鉴要动多个文件或需要先验证才决定是否落地时才建；能干净 `cherry-pick -x` 的单提交借鉴直接在 `main` 上做。验证通过后合并进 `main` 并推送；一次性分支事后删除，审计靠提交里的 `cherry picked from` 行和借鉴记录
- 锚点分支：`anchor/<日期>-<用途>`，在做 reset、对齐、历史重写等高风险动作前指向当时的 tip，事后保留而不是删掉

不要修改 Go module path `github.com/Wei-Shaw/sub2api`，除非用户明确要求并接受全仓导入路径迁移。

## General Working Rules

- 开始任何修改前，先读取 `git status --short --branch`、当前分支和相关文件，并 `git fetch --prune origin` 确认本地没有落后于 `origin/main`。分叉会让下一个会话重复产出已经存在的改动。
- 工作区可能包含用户改动。不得丢弃、覆盖、自动暂存或自动 stash 用户改动。
- 常规二开直接在 `main` 上提交并推送。
- 不做 rebase、不强推（含 `main`、`borrow/*` 与 `pick/*` 历史分支）。
- 移动或推送 tag、发布、部署仍需用户明确授权。常规 `git push origin main` 属于交付动作，按「提交与推送流程」执行。
- 不使用 `git reset --hard`、`git checkout -- <file>`、`git clean` 等破坏性命令处理冲突。确需执行时，先建 `anchor/*` 锚点分支并确认锚点可达，再动手。
- 借鉴上游内容时，冲突必须理解双方语义。无法确定原功能意图时保留工作现场并询问用户；不得为了省事先把不认同的上游改动整文件覆盖掉。
- 上游借鉴与私有业务改动分开提交，避免一个提交混合两类目的；借鉴来的实现允许为适配本 fork 的设计而重写，不要求与上游逐行一致。
- 未完成也提交：切换任务或结束会话前，工作区不得留下未提交改动；半成品用 `wip(<scope>): …` 提交并推送 `main`。

## 提交与推送流程（定版）

本节是本 fork 唯一允许的落地方式，2026-10-23 由「PR 流程」简化为「直接 push main」。单人 fork、无预发环境，PR 的 review/merge 环节与 CLA 机器人都是纯 overhead，故移除；CI（`backend-ci`、`security-scan`）由 push 触发，仍然生效。所有 `gh` 命令（如关闭 PR、查 CI）仍必须显式带 `--repo Xeubad/codeapi`——本目录 `gh` 的默认仓库是上游 `Wei-Shaw/sub2api`。

### 1. 开始前

```bash
git fetch --prune origin
git switch main
git merge --ff-only origin/main
```

本地 `main` 必须先与 `origin/main` 一致再动手；在过时基线上改会产生分叉。

### 2. 提交与自证

- 一个主题一个提交，跨主题拆成多个提交。
- 提交前跑与改动匹配的验证（见 Project-Specific Change Rules）。CI 在 push 后跑 `backend-ci` 与 `security-scan`，本地只是提前发现问题。
- `golangci-lint` 是 CI 门禁：仓库 `.golangci.yml` 开了 `errcheck.check-type-assertions`，单值类型断言要写 comma-ok；gosec 会把「测试里 spawn 自身二进制」判成命令注入，按仓库既有约定逐点 `//nolint:gosec` 并写明信任边界，**不得**为此放宽生产代码的检查。

### 3. 推送

```bash
git push origin main
```

提交信息标题用英文 conventional commits（`<type>(<scope>): <subject>`），正文中文，与现有历史保持一致。push 后若 CI 红了，修完用新提交重推，不要 amend 已推送的提交。

### 4. 收口

推送后补当天 `docs/WORK_LOG_CN.md` 条目（见 Work Log）。工作区不得留未提交改动。

### 5. 高风险动作之前

执行 `git reset --hard`、把本地 `main` 对齐远端、改动已发布迁移、动锁文件之前，先打锚点：

```bash
git branch anchor/$(date +%Y%m%d)-<用途> <当前 tip>
git log -1 --oneline anchor/<...>   # 确认指向预期提交后再动手
```

对齐 `main` 之前必须先证明本地那份没有独有内容，三条都要为空或可解释：

```bash
git diff --diff-filter=A --name-only origin/main main   # 本地独有的文件
git log --oneline origin/main..main                     # 本地独有的提交
git diff --stat origin/main main                        # 内容差异，逐文件看懂是「缺远端的更新」还是「本地独有修改」
```

锚点分支事后保留。默认只留本地；若这次动作可能不可逆且本地是唯一副本（例如要丢弃只存在于本地的提交），先 `git push origin anchor/<...>`。

## User Intent Mapping

以下说法视为「看一眼上游最近做了什么，并给出值得借鉴的功能清单」，不是合并授权：

- “拉取源仓库更新”
- “同步上游”
- “检查上游更新”
- “看看官方有什么更新”

收到上述请求后，直接执行 `git fetch --prune upstream` 与只读分析，不需要先询问是否允许 fetch。分析产出是**功能级清单**（官方新做了什么能力、哪个对本 fork 有用），不是逐提交决策表；清单给完后停下等一次明确确认。

以下回复视为对刚刚那份清单中所列功能的借鉴授权：

- “确认”
- “继续”
- “采纳”
- “按方案执行”

该授权允许：落地清单中列出的功能（`cherry-pick -x` 或按本 fork 的设计重写）、解决确定性冲突、跑验证、把结果提交进 `main` 并推送 `origin main`。它不授权：打 tag、发布或部署。清单未列出的内容不得顺手搬进来；用户对单个功能的取舍有最终决定权。

如果用户一开始明确说“拉取、分析并确认后采纳”，仍必须在分析报告后停在确认门，不得把最初请求当作提前确认。**不得把「上游又出了新版本」当成本 fork 的待办**：不借鉴就不产生任何工作项。

## 上游借鉴流程（按需，不跟踪）

没有固定批次、没有「必须评估完」的范围。按以下两种触发之一执行：

**A. 用户要求看一眼上游** → 只出报告，不动代码：

1. `git fetch --prune upstream`（只读引用；若 `upstream` 不存在，确认官方地址后再 `git remote add`；URL 与官方不同则报告，不自动改）。
2. 用 `git log --oneline <上次看过的点>..upstream/main` 或直接读上游 release notes 取近况。**不要求覆盖每一个提交**，也不因为「还有提交没看」而生成待办。
3. 产出**功能级清单**：官方新增/修复了什么能力 → 对本 fork 是否有用 → 建议「借鉴 / 改造后借鉴 / 不管」。涉及安全（鉴权、密钥、SSRF、注入、越权、计费一致性）的条目单独标 🔴 并说明不采纳的风险，但**是否采纳仍由用户决定**，不得当成必做项。
4. 停在确认门。此阶段不得 merge、cherry-pick、建分支或改代码。

**B. 用户点名要某个功能** → 直接进入下面的落地步骤。

借鉴时优先按本 fork 的设计重新实现，能干净搬运才用 `cherry-pick -x`。上游与本 fork 规则冲突的部分（固定 8080 端口、portguard 语义、私有迁移命名、二开护栏、UI 设计语言）一律以本 fork 为准，不为对齐上游而回退本地规则。

## 借鉴落地步骤（用户确认后）

1. 重新检查工作区、本地 `main` 与 `origin/main`。若工作区变脏（含并行会话的改动）或所依据的上游内容已变化，停止并重新出清单；旧确认失效。
2. 单提交能干净搬运的，直接在 `main` 上 `git cherry-pick -x <sha>`（保留出处行）；需要改写或多文件整合的，建 `borrow/<主题>` 分支，按本 fork 的设计实现。冲突必须理解双方语义后解决，禁止整文件 `ours/theirs` 覆盖。放弃清单中已列出的功能必须停下说明原因，经用户同意后在借鉴记录里留痕。
3. 借鉴来的功能若触及本 fork 不存在的上下文（上游改了我们没有的代码），按本 fork 的结构重新设计或直接标为不适用，不强行搬运。
4. 特别检查：迁移文件名与 SHA、Ent schema/生成代码、锁文件、配置默认值与前端嵌入产物；依赖变化用 `go mod tidy` / `pnpm install` 重建后核对 diff。
5. 运行与改动匹配的验证：`git diff --check`、受影响 package 的测试、`make build`。借鉴触及调度、计费、鉴权或网关核心路径时跑完整 `make test`。
6. 验证失败保留现场，不推送。验证通过则合并进 `main` 并推送 `origin main`，等 CI 全绿。
7. 当天把这次借鉴写进 `docs/UPSTREAM_DECISIONS_CN.md`（来源功能/提交、落地方式、改了哪些文件、与上游的差异）并更新 `docs/WORK_LOG_CN.md`。

推送前先展示将推送的远端、分支和提交范围，再执行普通 push；禁止 force push。

**⚠️ 本 fork 与 `upstream/main` 之间永不执行 `git merge`。** `main` 与上游无祖先关系，一次性合并会把「不借鉴」的东西也带进来，与解耦前提冲突。若未来用户决定回归 merge 模式，必须作为独立决策重新评估，不得在借鉴批次中顺手进行。

## Release and Rollback

- `main` 同时是集成分支和生产分支，不存在 `develop -> main` 提升步骤。
- 发布即从 `origin/main` 上已合并且 CI 全绿的 HEAD 打 tag，版本号形如 `v0.2.17-custom.1`；不给只存在于本地工作副本的提交打 tag。
- 数字段必须严格递增：`update_service.go` 的 `parseVersion` 丢弃 `-custom.N` 后缀，更新判定只比较三段数字，同版本号的二开 tag 永远显示"无更新"。
- 创建或推送标签、触发 Release workflow、部署前必须得到明确授权。
- 回滚用新的前向提交完成：借鉴进来的提交用 `git revert <sha>` 回滚并推送 `main`；按本 fork 设计重写的功能则用前向提交撤销，不假设上游语义可逆。已应用的数据库迁移不可回退，必须新增前向修正迁移。

## Work Log

每次工作结束后，把当天的产出补记到 `docs/WORK_LOG_CN.md`。

- **一天一节**，标题格式 `## 2026年10月20日`，新的一天插在最上面；同一天的多次工作合并进同一节。
- 每节固定四块：**改动 / 关键决策 / 验证 / 产出**。改动写做了什么，决策写为什么以及放弃了什么，验证写跑了什么命令、结果如何，产出写提交短 SHA。
- 只记已发生的事。没跑的检查写「未运行」和原因，不能写成「通过」；数字要来自刚跑过的命令。
- 不写流水账。只记对本 fork 有长期价值的信息：决策、踩坑、风险、被否决的方案。
- 产出块写 `origin/main` 上可见的短 SHA。cherry-pick 重建分支、或按 CI 反馈补提交都会改写 SHA，一律按最终落地的提交记，不记本地旧 SHA。
- 日志改动随当次工作一起提交，或单独一个 `docs` 提交，不夹在功能改动里。

## Project-Specific Change Rules

### 二开定位与安全护栏

2026-10-31 起，本 fork 的定位从「在小范围内做二次开发」改为「**完整的自主产品创作**」。功能、UI、信息架构、数据模型都可以按自己的产品意图重做，不受「只改一小部分」限制；上游做过什么不构成约束，也不构成必须跟进的义务。

**唯一不变的底线：平台自带的校验、安全检测、鉴权、限流、计费一致性校验与审计逻辑不得删除、绕过或弱化。** 这条与「是否跟上游」无关，是产品自身的安全姿态。需要改动时，只能在原实现基础上增强或精进，并同步补测试。如果某处校验确实挡住需求，先报告「为什么必须改、替代校验是什么、由哪些测试覆盖」，取得明确同意后再动；不得以“二开方便”为由先删后说。

新增功能优先落在新文件、独立 package 或 `backend/pkg/pluginapi/`，避免侵入账号调度、计费与网关核心路径——这条是**降低自研改动互相踩踏的工程约定**，不是范围限制；确实需要改核心路径时，按上一条流程说明风险后再动。

### 监听端口硬性条件（8080）

本 fork 的后端服务**只能监听 8080**，这是硬性条件，不是默认值：

- 唯一真相是 `backend/internal/config` 的 `config.RequiredServerPort`；`Config.Validate()` 对 `server.port != 8080` 直接 fail-fast 拒绝启动。
- 不得重新引入可配置端口的路径：`config.yaml server.port`、`SERVER_PORT` 环境变量、`--setup` 向导、`deploy/install.sh` 一律只能得到 8080。`AUTO_SETUP` 遇到非 8080 的 `SERVER_PORT` 直接报错。
- `deploy/docker-compose*.yml` 与 `.env` 里的 `SERVER_PORT` 只表示宿主机发布端口，容器内仍是 8080；前端 Vite dev server 的 3000 与代理目标 `VITE_DEV_PROXY_TARGET` 不受此条约束。
- 端口冲突由 `backend/internal/portguard` 处理：在初始化应用之前绑定监听并持有它。TTY 下先打印占用者（PID、用户、启动时间、命令行），得到 `y` 确认后才发 `SIGTERM`；非交互环境（systemd、Docker、管道）只打印报告并以非 0 退出，**绝不自动杀进程**。
- 护栏：PID ≤ 1 与本进程自身永不被列为可杀；`SIGTERM` 后最多轮询 5 秒，不升级到 `SIGKILL`，仍失败就报错退出。
- 上游同步时若上游改动端口默认值或引入端口配置项，保留本 fork 的固定端口与 portguard 语义，冲突一律以本 fork 为准。

### Backend

- Go 版本以 `backend/go.mod` 为准；当前 CI 还在 workflow 和 Dockerfile 中校验/固定版本。
- 修改 `backend/ent/schema/` 后运行 `cd backend && go generate ./ent`，提交生成代码。
- 修改 interface 时检查所有 mock、stub 和测试实现。
- 后端测试至少覆盖受影响 package；借鉴触及调度、计费、鉴权或网关核心路径时运行完整 `make test`。

### Frontend

- 只使用 pnpm。修改 `frontend/package.json` 时必须同步 `frontend/pnpm-lock.yaml`。
- 用户可见文本必须检查 i18n key 完整性。
- 前端最终构建到 `backend/internal/web/dist/` 并嵌入后端，不能只验证开发服务器。

### Database Migrations

- 已应用到任一环境的迁移不可修改、重命名或删除。系统按文件名和 SHA-256 校验。
- 私有迁移使用独立且不与上游冲突的命名，例如 `900_xeu_<description>.sql`，创建前仍要扫描现有文件名。
- 修正已发布迁移必须新增前向迁移。
- `*_notx.sql` 只用于允许的并发索引语句，并保持幂等。

### Extension Boundary

- 新 Provider 或独立传输能力优先评估 `backend/pkg/pluginapi/` 和 `docs/PLUGIN_DEVELOPMENT.md`。
- 能通过配置、插件或独立 package 实现的功能，避免侵入账号调度、计费和网关核心路径。

## Documentation

- 二开流程或分支约定变化时，同步更新本文件与 `docs/FORK_MAINTENANCE_CN.md`。
- 不要把 `CLA.md` 当作 Agent 指令文件；它是贡献者许可协议。
- 文档中的仓库名、命令、分支名必须与实际远端和 workflow 一致。
