# Sub2API Fork Agent Rules

本仓库是 `Wei-Shaw/sub2api` 的长期二次开发 fork。Agent 的目标是在保留私有改动的同时，持续、安全地吸收上游更新。

详细流程见 `docs/FORK_MAINTENANCE_CN.md`。本文是执行规则；发生冲突时，本文优先。

## Repository Identity

- 官方上游远端名：`upstream`
- 官方上游地址：`https://github.com/Wei-Shaw/sub2api.git`
- 当前 fork 远端名：`origin`
- 当前 fork 地址：`https://github.com/Xeubad/codeapi.git`（原 `Xeubad/sub2api`）
- 二开发布产物：镜像名跟随仓库名 → `ghcr.io/xeubad/codeapi`；二进制名、systemd 服务名、`INSTALL_DIR` 仍为 `sub2api` / `/opt/sub2api`
- 应用内检查更新的仓库：`Xeubad/codeapi`（禁止指回上游，否则会用官方二进制覆盖本地改动）
- 上游基线：`upstream/main`
- 二开主干：`main`，唯一长期分支，同时承担集成与生产。本 fork 单人维护、无预发环境，因此不设 `develop`
- **`origin/main` 是唯一真相**：本地 `main` 只跟随 `origin/main`，不承载任何未推送的合并或提交
- 功能分支：`feature/custom-*`、`fix/custom-*`、`docs/custom-*`，一个主题一条，从 `origin/main` 切出，**通过 PR 合入 `main`**，合入后删除本地与远端分支
- 上游同步分支：`sync/upstream-YYYYMMDD-<short-sha>`，一次性，同样**通过 PR 合入 `main`**；合入后保留作审计，不删除
- 锚点分支：`anchor/<日期>-<用途>`，在做 reset、对齐、历史重写等高风险动作前指向当时的 tip，事后保留而不是删掉

不要修改 Go module path `github.com/Wei-Shaw/sub2api`，除非用户明确要求并接受全仓导入路径迁移。

## General Working Rules

- 开始任何修改前，先读取 `git status --short --branch`、当前分支和相关文件，并 `git fetch --prune origin` 确认本地没有落后于 `origin/main`。分叉会让下一个会话重复产出已经存在的改动。
- 工作区可能包含用户改动。不得丢弃、覆盖、自动暂存或自动 stash 用户改动。
- 不在 `main` 上做任何提交。成果一律走「一个主题一条分支 + 一个 PR」。
- 已共享的 `main`、`sync/*` 和远端已存在的分支不做 rebase，不强推。本地 `main` 只允许快进到 `origin/main`。
- 未经用户明确授权，不执行 `git push origin main`、移动或推送 tag、发布、部署，也不在 PR 上点合并。特性/同步分支的 push 与 `gh pr create` 属于交付动作，按「上传与合并流程」执行。
- 不使用 `git reset --hard`、`git checkout -- <file>`、`git clean` 等破坏性命令处理冲突。确需执行时，先建 `anchor/*` 锚点分支并确认锚点可达，再动手。
- 合并冲突必须理解双方语义。无法确定业务意图时保留同步分支并询问用户。
- 上游可复用的修复与私有业务改动分开提交，避免一个提交混合两类目的。
- 未完成也提交：切换任务或结束会话前，工作区不得留下未提交改动；半成品用 `wip(<scope>): …` 提交并推送分支。

## 上传与合并流程（定版）

本节是本 fork 唯一允许的落地方式，2026-10-21 定死；要改流程先改本节，再改其它章节。所有 `gh` 命令必须显式带 `--repo Xeubad/codeapi`——本目录 `gh` 的默认仓库是上游 `Wei-Shaw/sub2api`，漏掉就把 PR 开到官方仓库去。

### 1. 开分支

```bash
git fetch --prune origin
git switch -c feature/custom-<topic> origin/main
```

分支只能从 `origin/main` 切，不从本地 `main`、不从其它特性分支切。已经混进别的主题合并提交的分支，不要原地 `git rebase --onto` 或 `git branch -f`（会被权限层拦，历史也难自证），改用 cherry-pick 重建干净分支：

```bash
git switch -c feature/custom-<topic> origin/main
git cherry-pick <旧分支上的提交范围>
```

cherry-pick 会改写 SHA，工作日志与 PR 描述一律按新 SHA 记。

### 2. 提交与自证

- 一个主题一条分支，跨主题就拆成多个 PR。
- 提交前跑与改动匹配的验证（见 Project-Specific Change Rules）。CI 在 PR 上还会跑 `backend-ci` 与 `security-scan`，本地只是提前发现问题。
- `golangci-lint` 是 CI 门禁：仓库 `.golangci.yml` 开了 `errcheck.check-type-assertions`，单值类型断言要写 comma-ok；gosec 会把「测试里 spawn 自身二进制」判成命令注入，按仓库既有约定逐点 `//nolint:gosec` 并写明信任边界，**不得**为此放宽生产代码的检查。

### 3. 推送并开 PR

```bash
git push -u origin feature/custom-<topic>
gh pr create --repo Xeubad/codeapi --base main \
  --title "<type>(<scope>): <subject>" \
  --body "<Summary + Test plan>"
```

标题用英文 conventional commits，正文中文，与现有历史保持一致。

### 4. 等 CI 绿，再合

- CI 全绿才允许合并；红了修完重推，不要先合再补。
- **默认由用户在 PR 页面点合并。** Agent 只有在用户明确说「合并它」之类的话时，才执行 `gh pr merge --repo Xeubad/codeapi <n>`。
- 合并方式固定为 merge commit：不 squash、不 rebase 合并。工作日志里的短 SHA 要能在合并后的 `main` 上对上。
- 尽量不要造 stacked PR。非要造，合并前必须先把自己的 base 改回 `main`：

  ```bash
  gh pr edit <n> --repo Xeubad/codeapi --base main
  ```

  base 还指向另一条特性分支时，页面上绿色的 Merge 按钮会把它合进那条分支，而不是 `main`。
- `.github/workflows/cla.yml` 是上游带来的 CLA 机器人，会在每个 PR 上要求签署贡献者许可协议。它属于预期噪音，不是本次改动引入的问题；是否移除该 workflow 由用户决定，Agent 不自行删改。

### 5. 收口

```bash
git fetch --prune origin
git switch main
git merge --ff-only origin/main
git branch -d feature/custom-<topic>
git push origin --delete feature/custom-<topic>
```

`sync/*` 的本地与远端分支都保留不删。删除远端分支是共享状态操作，执行前照例报出分支名和它上面的最后一个提交，确认一次再删。收口后按当天条目补 `docs/WORK_LOG_CN.md`；日志改动走它自己的 `docs/custom-*` 分支与 PR，不夹在功能 PR 里。

### 6. 高风险动作之前

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

以下说法视为“获取并分析上游”，不是直接合并授权：

- “拉取源仓库更新”
- “同步上游”
- “检查上游更新”
- “看看官方有什么更新”

收到上述请求后，直接执行本文的 Fetch and Analyze 阶段，不需要先询问是否允许 fetch。分析结束后给出合并计划，并等待一次明确确认。

以下回复视为对刚刚报告的、SHA 未变化的合并计划授权：

- “确认”
- “继续”
- “合并”
- “按方案执行”

该确认授权创建同步分支、合并、解决确定性冲突、运行验证，并把同步分支推送到 `origin`、开出指向 `main` 的 PR。它不授权：在 PR 上点合并、推送 `main`、打 tag、发布或部署。

如果用户一开始明确说“拉取、分析并确认后合并”，仍必须在分析报告后停在确认门，不得把最初请求当作提前确认。

## Fetch and Analyze

当用户要求拉取或同步上游时，按以下顺序执行：

1. 读取工作区状态、当前分支、所有远端及分支跟踪关系。
2. 若 `upstream` 不存在，确认官方地址后执行：
   `git remote add upstream https://github.com/Wei-Shaw/sub2api.git`。
3. 若 `upstream` 已存在但 URL 不同，不自动覆盖；报告差异并等待用户决定。
4. 验证 `origin`。若与当前 fork 地址不同，不自动覆盖。
5. 执行 `git fetch --prune upstream`。需要了解 fork 状态时再执行 `git fetch --prune origin`。
6. 二开基线固定为 `origin/main`。先 `git fetch --prune origin`；若本地 `main` 与 `origin/main` 分叉，先按「上传与合并流程」第 6 步打锚点并把本地 `main` 对齐远端，再规划。
7. 记录待合并的完整 `upstream/main` SHA。后续确认只对该 SHA 有效。
8. 分析并汇报：
   - 二开基线与上游各自 ahead/behind 数量；
   - 上游新增提交，按功能、修复、安全、依赖、CI/部署、数据库迁移分类；
   - 变更文件和受影响模块；
   - 本地二开提交与高冲突区域；
   - 潜在文本冲突，可使用 `git merge-tree` 预检；
   - 新增或修改的迁移、Ent schema/生成代码、`go.mod/go.sum`、`package.json/pnpm-lock.yaml`、Docker/CI/配置样例；
   - 建议验证命令、预计风险和回滚点。
9. 输出明确的拟执行计划：基线分支、同步分支名、上游 SHA、预计冲突、测试范围。
10. 停止并请求用户确认。此阶段不得 merge、rebase、切换/创建同步分支或改代码。

工作区不干净时仍可 fetch 和做只读分析，但报告中必须列出受影响文件。不得进入合并阶段，也不得自动 stash：先按「未完成也提交」把改动落成 `wip` 提交并推送分支，再进入合并阶段。

## Merge After Confirmation

用户确认刚刚的计划后：

1. 重新检查工作区、本地 `main`、`origin/main` 和 `upstream/main` SHA。
2. 若工作区变脏或上游 SHA 已变化，停止并重新分析；旧确认失效。
3. 从 `origin/main` 创建计划中的 `sync/upstream-*` 分支。
4. 执行 `git merge --no-ff upstream/main`。不得用 rebase 改写二开历史。
5. 解决确定性冲突，完成合并提交，并记录每个冲突采用上游、本地或组合方案的原因。
6. 特别检查迁移文件名冲突、已发布迁移被修改、Ent 生成代码、锁文件、配置默认值与前端嵌入产物。
7. 运行与变更匹配的验证。上游同步默认运行完整验证：
   - `git diff --check`
   - `make test`
   - `make build`
8. 若验证失败，保留同步分支和现场，报告失败原因；不得推送，也不得开 PR。
9. 若验证通过，按「上传与合并流程」推送同步分支，开出 base 为 `main` 的 PR，等 CI 全绿。合并在 GitHub 上完成，本地不预先合并 `main`。
10. 汇报 PR 编号、同步分支短 SHA、冲突决策、测试结果、未验证风险和建议的下一步；合并后收口，并把冲突决策写进当天工作日志。

不要自动删除同步分支。推送前先展示将推送的远端、分支和提交范围，再执行普通 push；禁止 force push，禁止由 Agent 直接推送 `main`。

若 `origin/main` 恰好是 `upstream/main` 的祖先（领先 0、落后 N），不要先合功能分支再同步：先开一条纯快进的同步 PR（零冲突）合入，再基于更新后的 `main` 开第二条特性 PR。这样大批量上游提交不进入三方合并，冲突只在功能分支那一侧出现一次。

## Release and Rollback

- `main` 同时是集成分支和生产分支，不存在 `develop -> main` 提升步骤。
- 发布即从 `origin/main` 上已合并且 CI 全绿的 HEAD 打 tag，版本号形如 `v0.2.17-custom.1`；不给只存在于本地工作副本的提交打 tag。
- 数字段必须严格递增：`update_service.go` 的 `parseVersion` 丢弃 `-custom.N` 后缀，更新判定只比较三段数字，同版本号的二开 tag 永远显示"无更新"。
- 创建或推送标签、触发 Release workflow、部署前必须得到明确授权。
- 回滚用新的前向提交完成，并同样走分支 + PR：对同步合并执行 `git revert -m 1 <merge-commit>`；已应用的数据库迁移不可回退，必须新增前向修正迁移。

## Work Log

每次工作结束后，把当天的产出补记到 `docs/WORK_LOG_CN.md`。

- **一天一节**，标题格式 `## 2026年10月20日`，新的一天插在最上面；同一天的多次工作合并进同一节。
- 每节固定四块：**改动 / 关键决策 / 验证 / 产出**。改动写做了什么，决策写为什么以及放弃了什么，验证写跑了什么命令、结果如何，产出写提交短 SHA。
- 只记已发生的事。没跑的检查写「未运行」和原因，不能写成「通过」；数字要来自刚跑过的命令。
- 不写流水账。只记对本 fork 有长期价值的信息：决策、踩坑、风险、被否决的方案。
- 产出块写 PR 编号和合并后在 `origin/main` 上可见的短 SHA。cherry-pick 重建分支、或按 CI 反馈补提交都会改写 SHA，一律按最终落地的提交记，不记本地旧 SHA。
- 日志改动走它自己的 `docs/custom-*` 分支与 PR，不夹在功能 PR 里。

## Project-Specific Change Rules

### 二开范围与安全护栏

用户已明确的二开范围仅限：

- UI 与页面样式；
- 账号转发逻辑；
- 平台功能中的小部分逻辑；
- 新增特定功能。

**原平台自带的校验、安全检测、鉴权、限流、计费一致性校验与审计逻辑不得删除、绕过或弱化。** 需要改动时，只能在原实现基础上增强或精进，并同步补测试。如果某处校验确实挡住需求，先报告「为什么必须改、替代校验是什么、由哪些测试覆盖」，取得明确同意后再动；不得以“二开方便”为由先删后说。

新增功能优先落在新文件、独立 package 或 `backend/pkg/pluginapi/`，避免侵入账号调度、计费与网关核心路径。

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
- 后端测试至少覆盖受影响 package；合并上游时运行完整测试。

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
