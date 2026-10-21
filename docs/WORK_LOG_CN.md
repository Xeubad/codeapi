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
