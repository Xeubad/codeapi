# Sub2API Fork 二次开发与上游同步手册

本文定义 `Xeubad/codeapi` 的长期二次开发流程。日常操作可直接交给 Agent；Agent 的强制执行规则位于仓库根目录 `AGENTS.md`。

## 目标

这套流程解决三个问题：

- 私有功能持续开发，但不污染上游基线；
- 能定期获取 `Wei-Shaw/sub2api` 的修复和新功能，并且**逐提交决定采纳还是拒绝**，而不是全量吞下；
- 所有采纳、冲突决策、测试和发布都有可审计记录。

## 上游更新模式：解耦 + 选择性采纳（2026-10-29 起）

2026-10-29 之前，本 fork 通过 `git merge upstream/main` 定期全量同步。这种方式无法拒绝不需要的上游改动（包括与二开约定冲突的改动，如固定端口策略），且分叉越大冲突越重。现改为解耦模式：

- 上游镜像由 `upstream/main` 远端引用承担，它只是只读快照，不是本地分支；`main` 不与它保持祖先关系；
- 上游更新以 release 为批次（如 `v0.2.18..v0.2.19`），逐提交评估后用 `git cherry-pick -x` 拿进 `main`；不想要的提交根本不进仓库；
- 每个评估过的提交（含拒绝的）必须记入决策台账 [UPSTREAM_DECISIONS_CN.md](./UPSTREAM_DECISIONS_CN.md)，避免重复分析和遗忘拒绝理由；
- `main` 与 `upstream/main` 之间不再执行 merge；未来若要回归 merge 模式，必须作为独立决策重新评估，不能顺手进行。

决策车道与安全红线见 `AGENTS.md` 的「Fetch and Analyze（评估批次）」；四个车道：🔴 安全修复默认必采、🟠 协议/Provider 修复按实际使用决定、🟡 新功能按产品方向决定、⚪ CI/文档/重构默认拒绝。

## 仓库与分支

| 名称 | 用途 |
|------|------|
| `upstream` | 官方仓库 `Wei-Shaw/sub2api` |
| `origin` | 二开 fork `Xeubad/codeapi` |
| `upstream/main` | 只读的官方主线 |
| `main` | 二开唯一长期分支，同时是集成分支和生产分支；常规二开直接在此提交并推送 `origin main` |
| `origin/main` | 唯一真相。本地改动提交后立即推送，不长期停留本地分叉 |
| `feature/custom-*` | 确需隔离的临时改动才建，事后删除；常规工作不建分支 |
| `pick/upstream-<release>-<short-sha>` | 单次上游采纳批次的工作分支，cherry-pick 与冲突处理在其上完成，验证通过后合并进 `main` 并推送；分支保留作审计，不删除 |
| `anchor/<日期>-<用途>` | 高风险动作（reset、对齐 `main`、动已发布迁移）前的锚点，事后保留 |

本 fork 单人维护、没有预发环境，因此不设 `develop`。两条长期分支只会在每次同步时多制造一次「谁提升谁」的人工步骤，收益要等有第二套部署才出现。

推荐流向：

```text
upstream/main（只读引用，按 release 划批次）
      |
      |  逐提交评估 -> 记入 UPSTREAM_DECISIONS_CN.md
      v
pick/upstream-<release>-SHA
      |
      |  cherry-pick -x（只拿采纳清单里的提交）
      v  merge 进 main -> push origin main -> CI 全绿
     main  ==  origin/main（唯一真相）
      ^
      |  直接 commit -> push origin main -> CI 全绿
常规二开改动（功能 / 修复 / 文档）
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

### 拉取并分析上游（评估批次）

你可以直接说：

```text
拉取源仓库更新
```

Agent 会自动：

1. 检查工作区、分支和远端；
2. 添加缺失且地址明确的 `upstream`；
3. fetch 最新上游；
4. 从决策台账读出已评估边界 tag，确定本批次范围；
5. 逐提交分类初判（🔴 安全 / 🟠 协议 / 🟡 功能 / ⚪ 其他）；
6. 汇总依赖、迁移、CI、Docker 和配置变化，用 `git merge-tree` 干跑佐证冲突风险；
7. 给出工作分支名、固定 SHA、拟采纳/拟拒绝清单、风险和测试计划；
8. 等待你确认。

这句话不会修改业务分支，也不会自动 push。评估结论以台账为准，没有台账条目的提交一律视为未评估。

### 确认采纳

阅读分析报告后，可以回复：

```text
确认，按方案执行
```

Agent 会基于报告中固定的批次 SHA：

1. 先证明本地 `main` 与 `origin/main` 无分叉（有分叉则打锚点对齐后再继续）；
2. 从 `origin/main` 创建 `pick/upstream-<release>-<short-sha>`；
3. 逐个 `git cherry-pick -x` 采纳清单中的提交，解决确定性冲突；
4. 运行完整测试和构建；
5. 把工作分支合并进 `main` 并推送 `origin main`；
6. 等 CI 全绿；
7. 把每个提交的最终决策（含落地短 SHA）写入决策台账，汇报批次短 SHA、冲突处理、测试结果和残余风险，收口并补工作日志。

如果批次 SHA、工作区或基线在确认前发生变化，Agent 必须重新分析，旧确认自动失效。方案中未列出的提交不得顺手采纳。

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

## 上游分析报告应包含什么

每次报告至少包含：

- 本地基线、批次范围 tag、上游头 SHA 和比较时间；
- 提交总数与逐提交分类初判（决策车道）；
- 上游功能、修复、安全更新和依赖变化的摘要；
- 受影响的后端、前端、数据库、部署和 CI 文件；
- 本地二开提交及预期冲突区域；
- 数据库迁移文件名或语义冲突；
- `go.mod/go.sum` 与 `package.json/pnpm-lock.yaml` 变化；
- 建议的工作分支名、cherry-pick 顺序和验证命令；
- 阻塞项、回滚点和需要人工决策的业务差异。

分析阶段只允许 fetch 和只读检查，不允许先 cherry-pick 再汇报。

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

Agent 不应简单采用 `ours` 或 `theirs` 覆盖整份文件；cherry-pick 场景下还应警惕「上游提交依赖的上下文本 fork 不存在」的情况，此时在台账标拒绝并说明，不强行搬运。

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

普通功能分支按影响范围运行测试。上游采纳批次默认执行：

```bash
git diff --check
make test
make build
```

其中 `make test` 覆盖后端测试与 lint、前端 lint/typecheck 和关键 Vitest。CI 还会执行部署脚本检查和安全扫描。

若本机缺少依赖或外部服务，Agent 必须说明未运行的检查、原因和剩余风险，不能把“未运行”写成“通过”。

## 失败与恢复

- cherry-pick 失败：保留 `pick/upstream-*` 工作分支，不推送；无法落地的提交在台账标拒绝并说明。
- 测试失败：在工作分支修复或等待决策，不把失败结果推进共享分支。
- CI 失败：在原分支补提交后重推，等 CI 重新全绿；不先合再补。
- 工作区存在用户改动：不自动 stash，不开始采纳；先落成 `wip` 提交并推送 `main`。
- 上游在确认后变化：固定批次 SHA 不一致时重新分析。
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
