# Sub2API Fork 自主二开与上游借鉴手册

本文定义 `Xeubad/codeapi` 的长期二次开发流程。本 fork 已完全脱离上游节奏，按自己的产品意图做完整创作，上游只是可选的功能参考库。日常操作可直接交给 Agent；Agent 的强制执行规则位于仓库根目录 `AGENTS.md`。

## 目标

这套流程解决三个问题：

- 本 fork 按自己的产品意图做**完整的自主创作**，不被上游的功能优先级牵着走；
- 上游作为**可选的功能参考库**：看到值得借鉴的内容就分析并合并进来，不值得的直接不管，不存在「跟不跟得上上游」的欠账；
- 实际借鉴进来的东西有可审计来源记录，冲突决策、测试和发布有留痕。

## 上游更新模式：完全解耦 + 按需借鉴（2026-10-31 起）

演进过程：`merge 全量同步`（≤2026-10-28）→ `解耦 + 逐提交选择性采纳`（2026-10-29）→ **`完全解耦 + 按需借鉴`（2026-10-31）**。

第二阶段的模式虽然能拒绝不需要的改动，但它把「上游每一个提交都要被评估并记录」变成了长期义务：单人维护的 fork 会持续为「不采纳」付分析成本，而且产品方向一旦与上游分道，逐提交台账就变成没人回看的流水账。第三阶段取消这个义务：

- 上游镜像由 `upstream/main` 远端引用承担，只是只读快照，**不是基线、不是参照系、更不是待办来源**；`main` 不与它保持祖先关系；
- 看上游由两个场景触发：① 用户主动想知道官方最近做了什么；② 某个具体功能值得借鉴。**不借鉴就不产生任何工作项**，也不要求覆盖每一个提交；
- 采纳方式自由：能干净搬运的单提交用 `git cherry-pick -x` 保留出处；需要整合的按本 fork 的设计重写，不要求与上游逐行一致；
- `docs/UPSTREAM_DECISIONS_CN.md` 由「逐提交决策台账」降级为「**借鉴来源记录**」：只记实际拿进来的东西（来源、落地方式、改动文件、与上游的差异），不再记录拒绝项，也不再用于划定「已评估到哪个 tag」；
- 永不执行 `git merge upstream/main`：一次性合并会把「不借鉴」的内容也带进来，与解耦前提直接冲突。未来若回归 merge 模式，需作为独立决策重新评估。

保留不变的硬约束：固定监听 8080 与 portguard 语义、私有迁移命名、平台自带的校验/鉴权/限流/计费一致性/审计逻辑不得弱化、`anchor/*` 锚点纪律、直接提交并推送 `main` 的落地方式。上游在这些点上做了什么，一律以本 fork 为准。

## 仓库与分支

| 名称 | 用途 |
|------|------|
| `upstream` | 官方仓库 `Wei-Shaw/sub2api` |
| `origin` | 二开 fork `Xeubad/codeapi` |
| `upstream/main` | 只读的官方参考镜像，不是基线也不是待办来源 |
| `main` | 二开唯一长期分支，同时是集成分支和生产分支；常规二开直接在此提交并推送 `origin main` |
| `origin/main` | 唯一真相。本地改动提交后立即推送，不长期停留本地分叉 |
| `feature/custom-*` | 确需隔离的临时改动才建，事后删除；常规工作不建分支 |
| `borrow/<主题>` | 借鉴工作分支（可选）。只在一次借鉴跨多文件或需先验证时才建，事后删除；审计靠借鉴记录与提交里的出处行 |
| `pick/upstream-*` | 历史遗留：2026-10-29~10-31 逐提交采纳批次的工作分支，保留作审计，不再新建 |
| `anchor/<日期>-<用途>` | 高风险动作（reset、对齐 `main`、动已发布迁移）前的锚点，事后保留 |

本 fork 单人维护、没有预发环境，因此不设 `develop`。两条长期分支只会在每次同步时多制造一次「谁提升谁」的人工步骤，收益要等有第二套部署才出现。

推荐流向：

```text
upstream/main（只读参考，按需查看，不跟踪）
      |
      |  看到值得借鉴的功能 -> 分析
      v
borrow/<主题>（可选）  或  直接 main 上 cherry-pick -x
      |
      |  搬运 or 按本 fork 设计重写；不借鉴的什么都不发生
      v  merge 进 main -> push origin main -> CI 全绿
     main  ==  origin/main（唯一真相）
      ^
      |  直接 commit -> push origin main -> CI 全绿
自主二开创作（功能 / UI / 修复 / 文档）
      |
      v  tag（只打在 origin/main 上已合并且 CI 绿的提交）
  vX.Y.Z-custom.N
```

`main` 由 Agent 直接提交并推送（单人 fork，无 PR 环节）。不要对它 rebase 或 force push。

## Agent 指令如何生效

Codex 会在开始工作时读取仓库根目录的 `AGENTS.md`，因此该文件对整个仓库生效。更新指令后，新建任务或重新启动会话可确保加载最新内容。更深目录若以后需要特殊规则，可以添加自己的 `AGENTS.md` 或 `AGENTS.override.md`。

`CLA.md` 是贡献者许可协议，不是 Agent 指令文件；不要在其中维护开发流程。

参考：[OpenAI 官方 AGENTS.md 说明](https://learn.chatgpt.com/docs/agent-configuration/agents-md)。

## 一次性初始化

初始化前先确认当前工作区无未提交改动：

```bash
git status --short --branch
git remote -v
```

标准远端配置：

```bash
git remote add upstream https://github.com/Wei-Shaw/sub2api.git
git fetch --prune upstream
git fetch --prune origin
```

如果 `upstream` 已存在，先检查地址，不要直接覆盖：

```bash
git remote get-url upstream
```

主干 `main` 在 fork 建立时就已存在并跟踪 `origin/main`，不需要额外初始化集成分支。

把上游的 push 地址禁掉，让误推在 git 层面直接失败，而不是靠自觉：

```bash
git remote set-url --push upstream DISABLED
```

建议启用冲突复用与更清晰的冲突上下文：

```bash
git config rerere.enabled true
git config merge.conflictStyle zdiff3
```

这些操作都可以交给 Agent。Agent 会先检查现状，再说明准备修改的远端或分支。

## 与 Agent 的操作约定

### 看一眼上游做了什么

你可以直接说：

```text
看看官方有什么更新
```

Agent 会自动：

1. 检查工作区、分支和远端；
2. 添加缺失且地址明确的 `upstream`；
3. `git fetch --prune upstream`；
4. 取近况（上游 release notes 或 `git log`，**不要求覆盖每一个提交**）；
5. 产出**功能级清单**：官方新做了什么能力 → 对本 fork 是否有用 → 建议「借鉴 / 改造后借鉴 / 不管」；
6. 涉及鉴权、密钥、SSRF、注入、越权、计费一致性的条目标 🔴 并说明不采纳的风险，但是否采纳由你决定，不作为必做项；
7. 等待确认。

这句话不会修改业务分支，也不会自动 push。**报告里没有「欠账」概念**：没列出的提交不代表待办，不借鉴就不产生工作项。

### 确认借鉴

阅读清单后，可以回复：

```text
确认，按方案执行
```

Agent 会：

1. 先证明本地 `main` 与 `origin/main` 无分叉（有分叉则打锚点对齐后再继续），并确认工作区没有并行会话的未提交改动被卷入；
2. 能干净搬运的单提交直接在 `main` 上 `git cherry-pick -x`；需要整合或跨多文件的建 `borrow/<主题>`，按本 fork 的设计实现；
3. 解决确定性冲突（禁止整文件 `ours`/`theirs` 覆盖），跑与改动匹配的验证（触及调度/计费/鉴权/网关核心路径时跑完整 `make test`）；
4. 合并进 `main` 并推送 `origin main`，等 CI 全绿；
5. 把这次借鉴写入 [UPSTREAM_DECISIONS_CN.md](./UPSTREAM_DECISIONS_CN.md)（来源、落地方式、改动文件、与上游差异），收口并补工作日志。

如果清单依据的上游内容或工作区状态在确认前发生变化，Agent 必须重新出清单，旧确认自动失效。清单未列出的内容不得顺手搬进来。

### 日常开发落地

任何主题（功能、修复、文档）都按同一套动作落地，命令细节以 `AGENTS.md` 的「提交与推送流程（定版）」为准：

```bash
git fetch --prune origin
git switch main && git merge --ff-only origin/main   # 先对齐基线
# ...改代码并本地验证...
git add <files>
git commit -m "<type>(<scope>): <subject>"
git push origin main
# CI 绿 -> 补工作日志
```

要点：

- 一个主题一个提交，直接落在 `main` 并推送 `origin main`，不建功能分支、不开 PR。
- 所有 `gh` 命令显式带 `--repo Xeubad/codeapi`，本目录里 `gh` 的默认仓库是上游。
- 已推送的提交不 amend；CI 红了用新提交补，不重写历史。
- 半成品也要提交（`wip(<scope>): …`）并推送，工作区不留未提交改动过夜。

### 推送结果

`git push origin main` 属于交付动作，Agent 执行前会列出目标远端和提交范围。任何情况下默认不允许 force push。

### 发布二开版本

需要发布时说：

```text
从当前 main 打 tag 发布 v0.2.17-custom.1
```

`main` 同时承担集成与生产，因此不存在 develop→main 的提升确认；但发布确认仍然独立于同步确认——同步授权不会顺带打 tag。Agent 在打 tag 前会确认 `main` HEAD 已通过完整验证、与上一个 tag 做过差异复核、并检查迁移与配置兼容性。

发布后需要知道：版本号数字段必须严格递增，`-custom.N` 后缀不参与更新判定。

## 上游借鉴清单应包含什么

「看一眼上游」的报告以**功能**为单位，不以提交为单位。至少包含：

- 比较时间点、`upstream/main` 头 SHA，以及这次看的大致区间（例如「上次看的时间点之后」或「某个 release 之后」）；
- 官方新增/修复了哪些**能力**（人话描述，不是提交标题罗列）；
- 每条对本 fork 的判定：借鉴 / 改造后借鉴 / 不管，以及一句话理由；
- 拟借鉴项的落地影响：涉及的后端/前端/数据库/配置/CI 文件、依赖变化（`go.mod/go.sum`、`pnpm-lock.yaml`）、迁移文件名或语义冲突、与本 fork 既有设计的差异；
- 拟借鉴项的落地方式（`cherry-pick -x` 还是按本 fork 重写）、验证命令与回滚点；
- 安全类条目（鉴权、密钥、SSRF、注入、越权、计费一致性）单独标 🔴 并写清不采纳的风险——**提示风险不等于强制采纳**。

分析阶段只允许 fetch 和只读检查，不允许先搬运再汇报。报告不需要「把每个上游提交都归类」，也不留下「还有多少提交没看」的欠账。

## 冲突处理原则

### 可以由 Agent 直接处理

- import、格式化和生成代码的机械冲突；
- 双方语义一致的配置项合并；
- lockfile 可由现有包管理器确定性重建的冲突；
- 测试因接口新增而需要补齐 mock/stub。

### 必须暂停确认

- 上游和二开改变了同一计费、配额或账号调度规则；
- 权限、安全默认值或数据保留策略发生冲突；
- 同一 API 字段在双方具有不同含义；
- 数据迁移可能覆盖、删除或不可逆转换生产数据；
- 解决方案需要删除私有功能或改变用户可见行为。

Agent 不应简单采用 `ours` 或 `theirs` 覆盖整份文件；借鉴场景下还应警惕「上游改动依赖的上下文本 fork 不存在」的情况——此时按本 fork 的结构重新设计，或判定该功能不适用，不强行搬运。借鉴记录只写实际落地的内容，不要求为每个未采纳的上游提交留条目。

## 项目特有风险

### 数据库迁移

服务启动时会自动执行 `backend/migrations/*.sql`，并在数据库中记录文件名和 SHA-256。任何已执行迁移都不能修改、改名或删除。

私有迁移建议使用独立命名空间：

```text
900_xeu_add_custom_feature.sql
901_xeu_backfill_custom_feature.sql
```

创建前仍需扫描现有文件，确保完整文件名唯一。修复已发布迁移时必须新增前向迁移。

### Ent 生成代码

修改 `backend/ent/schema/` 后运行：

```bash
cd backend
go generate ./ent
```

schema 和生成代码必须在同一变更中提交。

### 前端依赖与嵌入

前端只使用 pnpm。修改依赖后必须更新 `frontend/pnpm-lock.yaml`。

生产构建会把前端产物写入并嵌入后端，因此不能只验证 Vite 开发服务器：

```bash
make build
```

### Go 与 CI 版本

Go 版本不只存在于 `backend/go.mod`，CI 和 Dockerfile 也可能固定版本。同步涉及 Go 升级时，应一起检查：

- `backend/go.mod`
- `.github/workflows/*.yml`
- `Dockerfile`
- `Dockerfile.goreleaser`
- `backend/Dockerfile`
- `deploy/Dockerfile`

### 插件优先

新增 Provider 或独立传输能力时，先评估 `backend/pkg/pluginapi/`。插件或独立 package 通常比修改核心账号调度和计费路径更容易长期同步。

## 验证标准

普通改动按影响范围运行测试。借鉴上游内容、或改动触及账号调度/计费/鉴权/网关核心路径时，默认执行完整验证：

```bash
git diff --check
make test
make build
```

其中 `make test` 覆盖后端测试与 lint、前端 lint/typecheck 和关键 Vitest。CI 还会执行部署脚本检查和安全扫描。

若本机缺少依赖或外部服务，Agent 必须说明未运行的检查、原因和剩余风险，不能把“未运行”写成“通过”。

## 失败与恢复

- cherry-pick 失败或搬运代价过高：保留 `borrow/*` 工作分支，不推送；按本 fork 的设计重写，或在借鉴清单里标为「不借鉴」并说明原因（不写入借鉴记录，因为它没进仓库）。
- 测试失败：在工作分支修复或等待决策，不把失败结果推进共享分支。
- CI 失败：在原分支补提交后重推，等 CI 重新全绿；不先合再补。
- 工作区存在用户改动（含并行会话的未提交改动）：不自动 stash、不开始搬运；先按「未完成也提交」落成 `wip` 提交并推送 `main`，或停下来报告。改不属于本主题的共享文件前先 `git show HEAD:<file>` 比对，避免覆盖别人的在途编辑。
- 上游在确认后变化：清单依据的内容不一致时重新出清单，旧确认失效。
- 本地 `main` 与 `origin/main` 分叉：先 `git diff --diff-filter=A --name-only origin/main main`、`git log origin/main..main`、`git diff --stat origin/main main` 三条核验；确认无独有内容后打 `anchor/*` 锚点再对齐，有独有内容则停下来报告。
- 推送后发现问题：使用新的 revert/fix 提交并推送 `main`，不重写共享历史，不 force push。

## 发布产物命名

仓库已从 `Xeubad/sub2api` 改名为 `Xeubad/codeapi`。发布链路的命名规则如下，改任何一处前先确认整体一致性。

- **镜像名唯一来源**是 `.github/release-tools/release-images.sh`，它取 `${GITHUB_REPOSITORY##*/}` 并小写，因此当前产物是 `ghcr.io/xeubad/codeapi`。以后再改仓库名不需要改代码。
- `.goreleaser.yaml` 的 `dockers` 与 `docker_manifests` 两段是**死配置**：`release_matrix.py` 在生成发布配置时把它们置空，CI 的 goreleaser 步骤还带 `--skip=docker`。其中残留的 `sub2api` 字面量不产生镜像，不要为品牌化去改这十几处。
- Release 正文里的 `docker pull` 文案来自 `release.footer`，已改为引用 `GITHUB_REPO_NAME`，与镜像名同源，不会各说各话。
- **二进制名 `sub2api` 不可改。** 它被 `Dockerfile.goreleaser` 的 `COPY sub2api /app/sub2api`、`release_matrix.py` 的归档名 `sub2api_{ver}_{os}_{arch}` 与 tar 成员名、`test_release_matrix.py` 的 fixture、`deploy/install.sh` 的 `SERVICE_NAME` 与 `/opt/sub2api` 共同依赖。只改「镜像名的最后一段」不影响这条链。
- 发布标签 `v0.2.7-custom.1` 可以通过校验（`release_matrix.py` 的 `VERSION_RE` 允许 `-` 后缀），流水线会把该版本号写入 `backend/cmd/server/VERSION` 再嵌入二进制。
- **已知后果：** `update_service.go` 的 `parseVersion` 丢弃 `-` 之后的一切且只比较三段数字，所以 `0.2.7-custom.1` 相对 `0.2.7` 判定为「无更新」，后台的更新与回滚按钮实际不会亮。当前决定是保留 `-custom.N` 命名并接受这一点：升级一律走 Git + 自行发版，不要把「一键更新」当升级通道。
- Docker Hub 当前未启用（仓库 Actions secrets 为空），`release-images.sh` 与 goreleaser 的 DockerHub 分支恒 skip；Telegram 发布通知同样因缺少 secret 而跳过。GHCR 侧已经自洽：`release.yml` 通知里的 `GHCR_IMAGE="ghcr.io/${REPO,,}"` 与 `release-images.sh` 同源，都取仓库名。
- 将来若启用 Docker Hub，`release-images.sh` 会自动推 `<hub>/<仓库名>`，但 `release.yml` 里还有两处字面量仍写着 `sub2api`，必须一起改：`:304` 的 `dockerhub-description` `repository:`，和 `:350` Telegram 文案里的 `DOCKER_IMAGE=`。再加上 `VersionBadge.vue` 的 `DOCKER_IMAGE`，共三处。

自己的部署入口：

```bash
# 脚本安装（装的是本 fork 的 release）
curl -sSL https://raw.githubusercontent.com/Xeubad/codeapi/main/deploy/install.sh | bash

# 容器
docker pull ghcr.io/xeubad/codeapi:latest
```

本 fork **屏蔽了 Web 安装向导**：没有安装锁时服务直接退出并打印指引，不会在主端口上起向导。首次配置走两条终端路径：

- 脚本安装（systemd）：`cd /opt/sub2api && sudo -u sub2api ./sub2api --setup`，完成后 `sudo systemctl start sub2api`。`install.sh` 在 `config.yaml` 不存在时会跳过自动启动，避免 `Restart=always` 造成崩溃循环。
- 容器部署：`AUTO_SETUP=true` 配合 `DATABASE_*`、`REDIS_*`、`ADMIN_EMAIL`、`ADMIN_PASSWORD` 无人值守安装。

保留 `internal/setup` 包是为了这两条路径；只移除了 HTTP 入口、前端路由和启动守卫。

## 刻意保留指向上游的位置

以下内容不随改名和品牌化修改，目的是压缩与上游的冲突面：

- `deploy/README.md`、`deploy/APPLE_CONTAINER.md` 的 curl 安装命令。**2026-10-21 变更**：三份根 README（英文/中文/日文）的 curl、`git clone` 与 star-history 已改指 fork，README 合并为唯一一份中文 `README.md`，删除英文与日文版。原方案「保留指向上游以压缩冲突面」被否决——用户照抄 README 的安装命令时，上游 `install.sh` 会把官方二进制装进 `/opt/sub2api` 覆盖二开改动，冲突收益远抵不过这个风险。`deploy/README.md`、`deploy/APPLE_CONTAINER.md` 仍指向 `Wei-Shaw/sub2api`，属于同一隐患，待定。
- `Dockerfile`、`deploy/Dockerfile`、`Dockerfile.goreleaser` 里的 `LABEL maintainer` 与 `org.opencontainers.image.source`。后者在构建时被 `release-images.sh` 传入的 `--label org.opencontainers.image.source=https://github.com/$GITHUB_REPOSITORY` 覆盖，实际发布镜像的 source 已经是你自己的仓库
- `deploy/sub2api.service` 与 `install.sh` 生成的 systemd `Documentation=` 行
- 法务与支付文档链接：`admin_compliance.go`、`stores/adminCompliance.ts`、`SettingsView.vue`
- 首页与 Header 的 GitHub 外链：`HomeView.vue`、`KeyUsageView.vue`、`AppHeader.vue`（如需指向 fork，各是一行常量替换）

需要改的功能性指向只有三处，已经改完：`backend/internal/service/update_service.go` 的 `githubRepo`、`frontend/src/components/common/VersionBadge.vue` 的 `GITHUB_REPO` 与 `DOCKER_IMAGE`、`deploy/install.sh` 与 `deploy/docker-deploy.sh` 的仓库地址。

## 工作日志

按天记录二开工作的过程文档在 [WORK_LOG_CN.md](./WORK_LOG_CN.md)。每次工作结束后补记当天一节，格式为 `## 2026年10月20日`，固定包含**改动 / 关键决策 / 验证 / 产出**四块。

日志记的是决策过程，不是 `git log` 的复制：为什么这么做、否决了哪些方案、踩了什么坑、验证到底跑了没跑。产出块写 `origin/main` 上可见的短 SHA——cherry-pick 重建分支或按 CI 反馈补提交都会改写 SHA，按最终落地的记，不记本地旧号。

日志随当次工作一起提交，或单独一个 `docs` 提交，不夹在功能改动里。

## 提交和发布约定

建议提交格式（标题英文 conventional commits，正文中文，scope 用主题名）：

```text
feat(port): pin the listen port to 8080 and guard it before startup
fix(scheduler): correct ...
chore(deps): bump ...
docs(fork): update maintenance workflow
```

历史里出现过的 `feat(custom)`、`feat(brand)`、`feat(ui)` 都合规，不要求统一改成主题名；但新提交优先具体 scope，方便 `git log --grep` 定位。

落地直接经 `main`：一个主题一个提交，本地验证后直接推送 `origin main`，CI 全绿为准；不 squash、不 rebase、不 force push。

二开版本使用可辨识标签：

```text
v0.2.7-custom.1
```

标签只打在 `origin/main` 上已合并且 CI 全绿的提交上。标签、Release、镜像发布和部署都需要独立明确授权。
