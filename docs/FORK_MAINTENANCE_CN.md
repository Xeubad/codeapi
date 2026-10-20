# Sub2API Fork 二次开发与上游同步手册

本文定义 `Xeubad/codeapi` 的长期二次开发流程。日常操作可直接交给 Agent；Agent 的强制执行规则位于仓库根目录 `AGENTS.md`。

## 目标

这套流程解决三个问题：

- 私有功能持续开发，但不污染上游基线；
- 能定期获取 `Wei-Shaw/sub2api` 的修复和新功能；
- 所有合并、冲突决策、测试和发布都有可审计记录。

## 仓库与分支

| 名称 | 用途 |
|------|------|
| `upstream` | 官方仓库 `Wei-Shaw/sub2api` |
| `origin` | 二开 fork `Xeubad/codeapi` |
| `upstream/main` | 只读的官方主线 |
| `main` | 二开唯一长期分支，同时是集成分支和生产分支 |
| `feature/custom-*` | 私有功能开发，合入 `main` 后删除 |
| `fix/custom-*` | 私有缺陷修复，合入 `main` 后删除 |
| `sync/upstream-*` | 单次上游同步与冲突处理，快进 `main` 后保留作审计 |

本 fork 单人维护、没有预发环境，因此不设 `develop`。两条长期分支只会在每次同步时多制造一次「谁提升谁」的人工步骤，收益要等有第二套部署才出现。

推荐流向：

```text
upstream/main
      |
      v  merge --no-ff
sync/upstream-YYYYMMDD-SHA
      |
      v  merge --ff-only
     main  <--- feature/custom-* / fix/custom-*
      |
      v  tag
  vX.Y.Z-custom.N
```

`main` 是共享分支。不要对它 rebase 或 force push。

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

### 拉取并分析上游

你可以直接说：

```text
拉取源仓库更新
```

Agent 会自动：

1. 检查工作区、分支和远端；
2. 添加缺失且地址明确的 `upstream`；
3. fetch 最新上游；
4. 计算 ahead/behind；
5. 汇总新增提交和受影响模块；
6. 预检冲突、迁移、依赖、CI、Docker 和配置变化；
7. 给出同步分支、固定上游 SHA、风险和测试计划；
8. 等待你确认。

这句话不会直接修改业务分支，也不会自动 push。

### 确认合并

阅读分析报告后，可以回复：

```text
确认，按方案执行
```

Agent 会基于报告中固定的上游 SHA：

1. 从 `main` 创建 `sync/upstream-*`；
2. 使用 merge 吸收 `upstream/main`；
3. 解决确定性冲突；
4. 运行完整测试和构建；
5. 测试通过后，将本地 `main` 快进到同步结果；
6. 汇报冲突处理、测试结果和残余风险。

如果上游 SHA、工作区或基线在确认前发生变化，Agent 必须重新分析，旧确认自动失效。

### 推送结果

本地验证通过后，另行说：

```text
把这次同步结果推送到 origin
```

Agent 会先列出目标远端、分支和提交范围，再执行普通 push。默认不允许 force push。

### 发布二开版本

需要发布时说：

```text
从当前 main 打 tag 发布 v0.2.17-custom.1
```

`main` 同时承担集成与生产，因此不存在 develop→main 的提升确认；但发布确认仍然独立于同步确认——同步授权不会顺带打 tag。Agent 在打 tag 前会确认 `main` HEAD 已通过完整验证、与上一个 tag 做过差异复核、并检查迁移与配置兼容性。

发布后需要知道：版本号数字段必须严格递增，`-custom.N` 后缀不参与更新判定。

## 上游分析报告应包含什么

每次报告至少包含：

- 本地基线、上游 SHA 和比较时间；
- ahead/behind 数量；
- 上游新增功能、修复、安全更新和依赖变化；
- 受影响的后端、前端、数据库、部署和 CI 文件；
- 本地二开提交及预期冲突区域；
- 数据库迁移文件名或语义冲突；
- `go.mod/go.sum` 与 `package.json/pnpm-lock.yaml` 变化；
- 建议的同步分支名和验证命令；
- 阻塞项、回滚点和需要人工决策的业务差异。

分析阶段只允许 fetch 和只读检查，不允许先合并再汇报。

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

Agent 不应简单采用 `ours` 或 `theirs` 覆盖整份文件。

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

普通功能分支按影响范围运行测试。上游同步默认执行：

```bash
git diff --check
make test
make build
```

其中 `make test` 覆盖后端测试与 lint、前端 lint/typecheck 和关键 Vitest。CI 还会执行部署脚本检查和安全扫描。

若本机缺少依赖或外部服务，Agent 必须说明未运行的检查、原因和剩余风险，不能把“未运行”写成“通过”。

## 失败与恢复

- 合并失败：保留 `sync/upstream-*` 分支，不更新 `main`。
- 测试失败：在同步分支修复或等待决策，不把失败结果推进共享分支。
- 工作区存在用户改动：不自动 stash，不开始合并。
- 上游在确认后变化：固定 SHA 不一致时重新分析。
- 已更新本地 `main` 但尚未 push：保留同步分支作为审计和恢复点。
- 已推送后发现问题：使用新的 revert/fix 提交，不重写共享历史。

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

- `README.md` / `README_CN.md` / `README_JA.md` 的 curl 安装命令与 star-history 图片；`deploy/README.md`、`deploy/APPLE_CONTAINER.md`
- `Dockerfile`、`deploy/Dockerfile`、`Dockerfile.goreleaser` 里的 `LABEL maintainer` 与 `org.opencontainers.image.source`。后者在构建时被 `release-images.sh` 传入的 `--label org.opencontainers.image.source=https://github.com/$GITHUB_REPOSITORY` 覆盖，实际发布镜像的 source 已经是你自己的仓库
- `deploy/sub2api.service` 与 `install.sh` 生成的 systemd `Documentation=` 行
- 法务与支付文档链接：`admin_compliance.go`、`stores/adminCompliance.ts`、`SettingsView.vue`
- 首页与 Header 的 GitHub 外链：`HomeView.vue`、`KeyUsageView.vue`、`AppHeader.vue`（如需指向 fork，各是一行常量替换）

需要改的功能性指向只有三处，已经改完：`backend/internal/service/update_service.go` 的 `githubRepo`、`frontend/src/components/common/VersionBadge.vue` 的 `GITHUB_REPO` 与 `DOCKER_IMAGE`、`deploy/install.sh` 与 `deploy/docker-deploy.sh` 的仓库地址。

## 提交和发布约定

建议提交格式：

```text
feat(custom): add ...
fix(custom): correct ...
chore(sync): merge upstream at <short-sha>
docs(fork): update maintenance workflow
```

二开版本使用可辨识标签：

```text
v0.2.7-custom.1
```

标签、Release、镜像发布和部署都需要独立明确授权。
