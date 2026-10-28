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

## 2026年10月29日

**主题**：上游同步模式切换——从全量 merge 改为解耦 + 选择性 cherry-pick

### 改动

- 新增决策台账 `docs/UPSTREAM_DECISIONS_CN.md`：定义四条决策车道（🔴 安全必采 / 🟠 协议按 Provider 取舍 / 🟡 功能按方向取舍 / ⚪ CI/文档默认拒绝）、Release 条目模板、台账号规则（UP-YYYY-NNNN）；内置 2026-10-29 基线盘点：待评估范围 `v0.2.17..upstream/main` 共 121 提交（v0.2.18 + v0.2.19），上游头 SHA `a3298b46ebe`，评估基线 `origin/main` = `55351a7e2`。
- `AGENTS.md`：上游更新流程由「Fetch and Analyze → Merge After Confirmation」改写为「Fetch and Analyze（评估批次）→ Cherry-Pick After Confirmation」；同步分支 `sync/upstream-*` 更名为工作分支 `pick/upstream-<release>-<short-sha>`；新增安全车道红线（🔴 修复不得因冲突工作量跳过）与「回归 merge 模式须独立决策」限制。
- `docs/FORK_MAINTENANCE_CN.md`：目标、分支表、流向图、拉取分析/确认采纳章节、分析报告要求、冲突处理、验证标准、失败与恢复全部同步改写。
- `.gitignore`：白名单加入 `!docs/UPSTREAM_DECISIONS_CN.md`（原被 `docs/*` 规则拦截，首次提交时发现）。

### 关键决策

- 停止 `git merge upstream/main` 全量同步。动机：全量模式无法拒绝不需要的上游改动，且分叉达 86/121 提交后每次同步冲突面过大；逐提交评估后以 `cherry-pick -x` 采纳，不想要的提交不再进入仓库。
- **`main` 从此不与 `upstream/main` 保持祖先关系**，二者之间禁止 merge；回归 merge 模式必须作为独立决策重新评估。这是单向门，已写入两份文档。
- 保留 `upstream` 远端与 `git fetch` 流程不变；评估范围按上游 release tag 划分，已评估边界以台账为准。
- 决策记录先于落地：台账号在评估阶段分配，落地后回填短 SHA，避免「做了但不知道为什么」。
- 发现安全修复优先于流程切换收尾：CI 红因 `golang.org/x/net` GO-2026-6860，上游待评估批次已含修复（`cacc80b86` + `0a2a29b32`），正好作为新模式的第一个评估对象，不单独抢修。

### 验证

- `git merge-tree --write-tree main upstream/main`：退出码 1，43 文件冲突，热点在 `backend/internal/service`、`backend/internal/repository`、`.github/workflows`；双方都改过的文件 32 个——已记入台账基线盘点。
- `git rev-list main..upstream/main --count`：121，与 `v0.2.17..upstream/main` 一致（`v0.2.17` 及更早已包含在 `main`）。
- `grep sync/upstream` 两份流程文档：无残留旧分支名。
- `git push origin main`：`55351a7e2..fdb9e1372`，成功。
- CI：`Security Scan` **失败**（GO-2026-6860，`ollama_cloud_usage_parser.go` 调 `html.Parse`；本次 push 前已连续失败 3 次，属历史遗留，与本批文档改动无关）；`CI`（backend-ci）在本报告完成时仍在排队，未运行。

### 产出

- `f32de745e` docs(fork): switch upstream sync to decoupled cherry-pick model
- `fdb9e1372` docs(fork): whitelist upstream decisions ledger in gitignore

---

## 2026年10月28日

**主题**：上游聚合规划文档提炼入仓 + 清理 demo 目录

### 改动

- 仓库外 `sub2api-demo/UPSTREAM_HUB_PLAN.md`（规划稿 v6，含 v1–v5 演进痕迹）提炼为 `docs/UPSTREAM_HUB_DESIGN_CN.md`：按当前代码实况修订，标注已实现/未实现（预算槽位接缝、探活三级、GroupsView 双来源、账号管理只读标记等），并汇总已知缺口（当日 9 项，同日排查三个实测问题后增至 12 项）。
- `sub2api-demo/PENDING_INTENT_ROUTING.md` 迁入 `docs/PENDING_INTENT_ROUTING_CN.md`，内容基本原样（仍为待定项）。
- 删除 `sub2api-demo/` 整个目录（含 989 行交互原型 `upstream-hub-proto.html`）。

### 关键决策

- 原型 HTML 不保留：真实 UI 已上线且经多轮迭代与原型分叉，原型只剩历史价值；目录不在版本控制内，删除即永久（管理员明确授权）。意图路由文档中对 `dlg-intent` 原型的引用改为"需要时重建"。
- 设计文档不照搬规划稿：规划期"零新增后端端点"的设想已被实际实现（专用 admin API + 新表）取代，以代码实况为准。

### 验证

- 未运行测试（纯文档变更）。

### 产出

- `b9aba6fa4` docs(upstream-hub): distill plan and pending intent routing into repo docs

**主题（同日第二次工作）**：修复管理员实测的三个问题——API 端点块消失、挂组后无模型、跨平台分组收不到条目；并新增聚合调度页

### 改动

- **P0 端点块**：`KeysView` 的 `EndpointPopover` 去掉 v-if，`api_base_url` 为空时回落 `window.location.origin`；`UseKeyModal` 与 cc-switch 导入同口径。
- **条目类型修正**：hub 条目从 `type=upstream` 改为 `type=apikey`（自定义 `base_url`）。网关凭据解析 `OpenAIGatewayService.GetAccessToken` / `GatewayService.GetAccessToken` 及全仓 20+ 处 `switch account.Type` 都不认 upstream，会报 `unsupported account type`。新增迁移 `900_xeu_upstream_entry_type_apikey.sql`，只转换带 `extra.upstream_vendor_id` 且未软删的行。repo 条目查询改为按 `extra.upstream_vendor_id` 的 jsonb 谓词过滤（数字/字符串双形态），与 type 解耦。
- **模型探测**：新增 `POST /admin/upstream-entries/:id/probe-models`（只读）与 `PUT /admin/upstream-entries/:id/models`（写 `model_mapping`，from=to）；新增 `GET /admin/upstream-entries` 供聚合页跨厂家编排；`UpstreamHubService` 注入 `AccountTestService`（接口 `upstreamModelProber` 便于单测）。
- **防抹写**：`UpdateEntry` 改为以现有 credentials 为底合并（`MergePreservingSensitiveCreds` 对非敏感键是"完全由 incoming 决定"，原实现会在任意编辑后清空 `model_mapping`）；复制条目连带复制模型白名单。
- **前端**：条目标题栏显示「已登记 N 个」，下方折叠面板（探测 → 勾选/全选/清空 → 保存落库；未保存行走 `sync-upstream-preview`）；分组下拉对 composite 全量开放并打「综合」徽章；详情表新增模型数列。
- **P2**：分组平台匹配放宽为「平台相等 或 composite」（对齐 `canCopyAccountsFromGroupPlatform`）；新增独立页 `/admin/aggregation`（供给汇入 / 归一模型表 / 链路视图四列），路由、侧边栏「聚合调度」、zh/en i18n、组件测试。

### 关键决策

- **不逐个放宽网关 type 开关**：20+ 处 `switch account.Type` 全量放宽等于在网关核心做大规模手术；`apikey` + 自定义 `base_url` 是平台日常在跑的国产厂商手动 Key 形态，所有转发路径天然支持。归属规则仍是文档既定的 `extra.upstream_vendor_id`，type 从来不是判据。
- **探测结果勾选后再落库**：不自动全量写入，避免中转商上百个杂模型灌进客户 `/v1/models`。
- **聚合调度不复刻调度引擎**：composite 链路（显式路由 → 账号归属 → 模型名探测 + priority 候选池 + 冷却换号）已是 magpie 式语义，新页面只做编排与可视化。
- 放弃的备选：扩 GroupsView 的 composite 页签（6890 行继续变重）；把聚合调度塞回上游聚合页（与已定的"供货层/客户面"分层相反）。

### 验证

- 后端：`gofmt -l` 干净；`go build ./...` 通过；`go test ./internal/service/ ./internal/repository/ ./internal/handler/admin/` 全绿（service 142s）。
- 前端：`check:i18n` 3 用例通过、`vue-tsc -b --force` 0 错误、`vite build` 通过、新增 `AggregationView.spec.ts` 3 用例通过。
- **运行态端到端（本地 mock 上游）**：起 18081 的 OpenAI 兼容 mock（`/v1/models` 两个模型 + `/v1/chat/completions`），用开发库 admin JWT 走完：建厂家 → 建条目（deepseek, base_url=127.0.0.1:18081）→ probe 返回两模型 → 写白名单 → 建 composite 分组并挂载 → 建客户 key → `GET /v1/models` 返回 mock 模型（不再回落 Claude 默认表）→ `POST /v1/chat/completions` 返回 `MOCK_OK`。迁移后历史条目为 apikey，软删行（3/4/6）未受影响。
- 未运行：真实第三方上游连通性（需管理员提供 Key）；浏览器内实点（自动化被安全策略禁止注入登录态，页面由组件测试 + 构建覆盖）。

### 产出

- 本次提交（见 `git log`，按主题拆分：keys 兜底 / hub 模型探测与 apikey 条目 / 聚合调度页 / docs）。

## 2026年10月27日

**主题**：上游聚合交互收口（补记，10/26 MVP 后的多轮实测迭代）

### 改动

- 黑屏修复：条目 `groups` nil slice 序列化为 `null` 导致前端 TypeError——后端初始化空 slice + 前端 `groupRefs()` 兜底（`6477b8cc3`）。
- 复制条目：服务端 `POST /:id/duplicate` 拷贝凭据，避免明文回显（`6477b8cc3`）。
- 编辑复用添加弹窗，统一为单一「添加/编辑上游」入口（`23b4fdcc8`、`ee2081b3a`）。
- 超预算硬阻断：弹窗内红色横幅 + 保存禁用 + 点击拦截；修复实时求和漏算既有条目编辑值的 bug（`bd23b2b71`、`dcfff5056`）。
- 去掉详情卡挂载动画（布局跳动，`cdfcd466a`）。
- 分组配置并入添加/编辑弹窗；分组选择改长条下拉（`be1fb4334`、`a29af2721`、`0af031fcf`）。
- 弹窗布局定为「一 key 一行」表格 + 「+ 添加 Key」；修复 CI `x-invalid-end-tag`（行号手术遗留的多余 `</div>`，`29d0ba778`、`1b6946f2e`）。
- `BaseDialog` 新增 `max` 宽度档（95vw），行内容单行不折行（`5b48f4243`）。
- 下拉分组选择最终形态：整行点击选中、选中态=底色高亮，删除复选框与行内开关；每分组「暂不调用」改为点详情表分组色条切换（`873bb3b1c`、`e03e5a205`）。

### 关键决策

- 管理端 UI 口味定稿：无挂载动画；硬校验屏内阻断（不能带病退出弹窗）；一记录一行宽表；长条下拉不用圆角药丸；选中态用底色而非复选框。
- 下拉面板用固定定位 + 背板，避免被表格 overflow 裁剪。

### 验证

- 每轮均由管理员在浏览器实测反馈驱动；push 后 CI（backend-ci + security-scan）全绿。
- 踩坑：CI 前端 eslint 报 `x-invalid-end-tag 339:7`，本地 eslint 因缓存/工作区与 HEAD 不一致一度误导——push 前必须核对"已提交内容 == 工作区"。
- CI `shell` job 两次 0 步骤 cancelled（macos-15 基建抖动），代码 job 全绿，按基建问题处理。

### 产出

- `6477b8cc3`…`e03e5a205`：`git log --oneline bdb49a5af..e03e5a205 --grep=upstream-hub` 计 15 个提交（含 1 个 gofmt），HEAD `e03e5a205`。

## 2026年10月26日

**主题**：上游聚合（Upstream Hub）MVP 落地：厂家台账 + 透传条目 + 分组引用开关

### 改动

- **数据模型**：新增 `upstream_vendors` 表（名称唯一、账号级并发预算、备注，软删）；`account_groups` 增加 `enabled` 列（本组是否调用该账号，默认 true）。迁移 `900_xeu_upstream_vendors.sql` / `900_xeu_account_group_enabled.sql`，Ent schema 同步并重新生成。
- **后端**：`internal/service/upstream_hub.go` 新增厂家 CRUD、条目（`Account type=upstream`）编排、预算超限警告、条目→分组引用与逐组开关；`internal/repository/upstream_vendor_repo.go` 持久化实现；`internal/handler/admin/upstream_hub_handler.go` + `/admin/upstream-vendors`、`/admin/upstream-entries`、`/admin/upstream-groups` 路由。
- **调度**：`queryAccountsByGroup` 在可调度路径追加 `account_groups.enabled=true` 过滤；分组账号统计 SQL 同步排除本组停用条目。
- **前端**：`api/admin/upstreamHub.ts`、`views/admin/UpstreamHubView.vue`（厂家列表、一页建上游多 key、条目管理、分组引用与开关）、路由 `/admin/upstream-hub`、侧边栏与 i18n（en+zh）。

### 关键决策

- 条目复用现有账号体系（`Account(type=upstream)` + `extra.upstream_vendor_id` 归属），不新建实体；厂家属全局，跨聚合共享预算。
- 「暂不调用」只落在 `account_groups.enabled`（仅本组候选池剔除），与账号全局 `schedulable` 两级独立并存。
- 条目创建复用 `AdminService.CreateAccount` 全套校验，`SkipDefaultGroupBind` 避免自动挂默认组。

### 验证

- `go build ./...`、`go test ./internal/service/ ./internal/repository/ ./internal/handler/admin/` 通过（service 139s）。
- 新增 `upstream_hub_test.go`：厂家校验/超预算警告/extra 归属解析/BaseURL 归一 均通过。
- 迁移 SQL 在 dev 库事务内执行后 ROLLBACK，语法与幂等性通过（未实际写入）。
- 前端 `pnpm run check:i18n`、`pnpm exec vue-tsc -b`、`pnpm build` 全部通过；`make build-backend`（embed tag）通过。
- 运行态：8080 由外部启动的进程提供服务（当前 PID 2313），`/health` 正常，入口资源 `assets/index-BoyD8vrm.js` 与 dist 一致；upstream-hub 页面/接口未逐页点验。
- `golangci-lint` 已后台执行，结果见提交后说明。

### 产出

- `bdb49a5af` feat(upstream-hub): vendor ledger, passthrough entries and per-group switches

**主题（同日第二次工作）**：首页 Gateway 面板视觉收口 + 公共页顶栏/页脚收敛

### 改动

- Hero 右侧去掉 `relative min-w-0 lg:pl-2` 的硬偏移，改为玻璃“设备框”：外层 24px 圆角 + 8px padding + 细边 + 投影，深色面板不再直接悬空。
- 迷你延迟折线从深色面板内部移到玻璃框底部的 LIVE 状态带（LIVE + 折线 + 当前 ms + 样本数），消除暗色区域底部的大块空白；整条状态带可点击跳到 `#latency`。
- 首页与登录页顶栏移除语言切换控件，默认中文；主题按钮移动到登录按钮右侧，顺序为「Key 用量 → 登录 → 主题」。
- 页脚删除版本号与右侧链接组，只保留 Logo + 版权行。

### 关键决策

- 不用倾斜/3D/夸张透视，改用“设备框 + 状态带”建立层次，保持 ink + lime + 玻璃的统一语言。
- 公开页继续不出现模型/厂商/价格信息；模型信息仍只经登录后的模型广场访问。

### 验证

- pnpm 12：`typecheck`、`test:run`（2845 用例）、`build` 全部通过。
- `go build -tags=embed` 构建通过；8080 当前进程 PID 2313，`/health` 返回 ok，入口资源 `assets/index-BoyD8vrm.js` 与 dist 一致（本次另起的 2959 因端口被占用未启动）。
- 1440 / 390、浅色 / 深色实测：`scrollWidth === innerWidth`，玻璃框与 LIVE 状态带正常；顶栏无语言切换；页脚仅版权行。

### 产出

- 未提交：`frontend/src/views/HomeView.vue`（本轮视觉收口；登录页语言切换移除与文案收敛已在既有提交中）

**主题（同日第三次工作）**：首页 Hero 面板换内容——从「第二套延迟看板」改成 base_url 接入卡；顺带修复前端构建

### 改动

- Hero 右侧深色面板删除：`实时往返 / 协议 / 计费` 指标行、`接入端点` 列表、玻璃框底部整条 LIVE 状态带（LIVE + 折线 + 当前 ms + 样本数）。面板内不再出现任何 ms 数字或折线。
- 换成接入卡：OpenAI / Anthropic 两行 base_url diff（官方默认地址 → 本站地址，行内复制按钮），下方鉴权头对照 `Authorization: Bearer <key>` / `x-api-key: <key>`，再保留「模型与价格 → 登录查看」引导；面板标题改 `GATEWAY / BASE URL`，右上角由探针推导的「在线/连接异常」胶囊换成本站 host。
- 语言包（zh/en 同步）：删除 `home.gateway.{live,online,latency,protocol,protocolValue,billing,billingValue,endpoints,samples,viewLatency}`，新增 `{panelTitle,baseHint,copyBase}`；剪贴板写入抽成 `writeToClipboard()`，curl 复制与 base_url 复制共用，base_url 复制状态同一时刻只高亮一个。
- 新增回归用例：`640×150` 折线只允许出现在 `#latency` 内，防止再把延迟可视化搬回首屏面板。
- **修复构建**：`frontend/src/views/admin/UpstreamHubView.vue:339` 多出一个 `</div>`（div 深度在该行变为 −1），`vite build` 直接 `Invalid end tag` 失败。该错误已随 `a29af2721` 进入 `main`，删除多余闭合标签后恢复。
- 文档：`GATEWAY_HOME_DESIGN_RESEARCH_CN.md` §12.2 / §12.5 / §12.13 / §12.15 按新面板口径改写。

### 关键决策

- 用户判定：页面下方已有整块「实测延迟」，Hero 面板再做一套同样的折线 + ms 属于重复信息 → 面板职责收敛为「怎么接上」，延迟只保留左侧一枚状态胶囊与下方实测区。上一轮只换了外框（玻璃设备框）没换内容，是误读需求。
- 面板内容必须可核对：base_url 取 `window.location.origin`，鉴权头以 `internal/server/middleware/api_key_auth.go` 实际接受的入站方案（Bearer / x-api-key / x-goog-api-key）为准；仍不写模型名、厂商、价格、倍率，公开页信息边界不变。
- 放弃的备选：控制台功能清单（与底部 CTA 区重复）、把端点列表搬进面板（与「快速接入」的 curl 与端点 chips 重复）。
- `x-goog-api-key`（Gemini CLI 兼容）没写进面板：首页只讲两条主路径，三条并排会让卡片变吵，且不属于中转站首页的核心信息。

### 验证

- pnpm 12：`check:i18n`（3 用例）、`typecheck`、`test:run` 372 文件 / **2846** 用例（含新增 1 例）、`build` 全部通过；`git diff --check` 干净。
- `make -C backend build`（`-tags embed`）通过；`make -C backend test` 52 个包 `ok`（含 `internal/service` 150.5s），无 FAIL/panic；`golangci-lint` 本机未安装，该步骤报 `No such file or directory`，**未运行**，交由 CI `backend-ci` 把关。
- 浏览器实测（dev server 3000，与 8080 同一份前端代码）：面板 innerText 无 `ms` 数字；两个复制按钮分别写入 `http://localhost:3000/v1` 与 `http://localhost:3000`，点击后文案变「已复制」、1.6s 复位且互斥（用 stub 替换 `navigator.clipboard.writeText` 验证，真实剪贴板在隐藏页无用户激活会失败）；无 `/model-plaza` 链接；控制台无 i18n 缺键或 Vue 告警。
- 1440 / 500、浅色 / 深色截图：面板高度与左栏基本对齐（约 460 vs 440），版心无横向溢出。
- 8080 现由并行会话启动的进程持有（PID 33911），其内嵌资源为改版前的 `assets/index-Bu6a4PD1.js`（已不在当前 dist 中），本轮未重启该进程，故 8080 暂未反映本次改动。

### 产出

- 未提交（本轮会话的权限模式拒绝 commit/push，改动全部留在 `main` 工作区）。落地时按主题拆三个提交：
  1. `fix(upstream-hub): drop stray closing div that broke the frontend build` → `frontend/src/views/admin/UpstreamHubView.vue`
  2. `feat(home): 首页 Hero 面板改为 base_url 接入卡` → `HomeView.vue` + zh/en `landing.ts` + `HomeView.compact.spec.ts`
  3. `docs(home): 设计文档 §12 面板口径与本次工作日志` → 本文档 + `GATEWAY_HOME_DESIGN_RESEARCH_CN.md`

**主题（同日第三次工作）**：全站 UI 重构选型调研（React Bits 与通用组件库对比）入档

### 改动

- 新增 `docs/UI_REFACTOR_RESEARCH_CN.md`：React Bits（reactbits.dev）官网实抓调研 + Vue 生态等价物对比（vue-bits / Element Plus / Naive UI / Ant Design Vue / shadcn-vue）+ 本项目前端实证（88 视图、217 组件、common 层引用广度、令牌层评价）+ 四层目标架构与四阶段周期估算（全职 5-7 周口径）。

### 关键决策

- 本项目前端为 Vue 3，React 生态库（React Bits / Ant Design / shadcn/ui）全部排除直接使用，仅作选型参照；动效点缀用同作者的 vue-bits。
- 结论：推荐 Element Plus 做基础层 + `App*` 封装层收口页面依赖 + vue-bits 限高视觉页点缀；当前 UI 风格评价 7.5/10 属资产，EP 默认风格差距 7-8/10，主题映射可收回 85-90%。
- Element Plus 的最终采用待用户确认，Phase 0（依赖+主题映射+封装层骨架）未启动。

### 验证

- 调研文档产出，无代码变更；`pnpm build` 与后端测试未运行（不适用）。
- 外部数据（stars、组件清单、许可）为 2026-10-28 GitHub API 与官网实抓值，文档附录已注明来源。

### 产出

- `390b6f5b3` docs(ui-refactor): add research on React Bits vs Vue component libraries（含 `.gitignore` 白名单同步；工作区另有同日第二次工作等并行未提交改动，未代为提交，待统一落地）

## 2026年10月25日

**主题**：① 价格目录切换到 models.dev + 模型广场改版（块网格/弹窗/管理端比价）；② 包管理器升级 pnpm 9 → 12（供应链加固 + CVE 修复窗口外移）+ 补两个新披露漏洞的版本钉子

### 改动

**（上午）价格目录与模型广场**

- **价格目录数据源**：`pricing.remote_url` 默认从 model-price-repo 切换为 `https://models.dev/api.json`；`pricing_service.go` 的 `parsePricingData` 入口探测 models.dev 格式（厂商→模型两层嵌套），新文件 `pricing_modelsdev.go` 拍平为 LiteLLM 扁平结构后复用原解析/计费管线，探测失败回退 LiteLLM 路径。单位换算 USD/百万token → USD/token；一方厂商（anthropic/openai/google/xai/deepseek/moonshotai…）价格优先，无一方收录时按字典序稳定取首个；google/moonshotai/zhipuai 换算回 gemini/moonshot/zhipu 保证管理端按平台同步不受影响。图片输出模型 output 映射 `output_cost_per_image_token`（models.dev 该口径即图片 token 价），长上下文 tier 折算为 `*_above_<N>k_tokens` 字段复用现成折算器。
- **配置**：models.dev 不发布 .sha256 → `hash_url` 默认置空，刷新走 `update_interval_hours` 文件年龄检查；allowlist `pricing_hosts` 增加 `models.dev`；`config.example.yaml` 同步。
- **用户广场改版**：删除 8 列大表格组件，改为 `PlazaModelGrid`（模型块：名称/本站简价/倍率/能力徽章）+ `PlazaModelDialog`（弹窗：本站价格 × 倍率 + 官方参考价两区，含阶梯、1h 缓存、分时、思考倍率）；计价逻辑抽到 `plazaPricing.ts` 共用。
- **管理端比价**：转换时保留模型→厂商报价索引（原始口径），新增 `GET /admin/channels/pricing/provider-quotes`；弹窗内管理员专属「其他厂商」区（胜出厂商「计费采用」徽章排首、最低输入价高亮），非管理员不渲染不发请求。
- 测试：转换器冒烟（真实 api.json 全量）、格式误判防回归、报价索引存活回归、组件测试 21 项。

**（下午）pnpm 12 迁移**

- **pnpm 9 → 12 迁移**：`frontend/pnpm-workspace.yaml` 从占位残留改为正式配置（`allowBuilds` 白名单 esbuild/vue-demi，overrides 从 package.json `pnpm` 字段迁入——pnpm 10 起该字段不再读取）；删除 `frontend/package.json` 失效的 `pnpm` 字段与 pnpm 9 时代残留的 `frontend/.npmrc`（`ignore-scripts=false` 本就是默认值，`legacy-peer-deps` 是 npm 专用键，构建脚本策略由 allowBuilds 接管）；lockfile 用 pnpm 12.4.1 重生成。
- **CI 与镜像**：`backend-ci.yml`、`release.yml`、`security-scan.yml` 的 `pnpm/action-setup` version 9 → 12；`deploy/Dockerfile` corepack 钉 `pnpm@12.15.0`（当前 12 线最新稳定版）。
- **顺手补两个新披露漏洞的钉子**：audit 发现 dompurify ≥3.4.17（critical，SAFE_FOR_TEMPLATES bypass）与 follow-redirects ≥1.16.1（high，header 凭证泄露）有补丁版，override 钉子上调 `dompurify@<3.4.17`、新增 `follow-redirects@<1.16.1`，与既有按补丁升钉惯例一致。

### 关键决策

**（上午）**

- **保留计费链路不动，只换数据层**：转换成 LiteLLM 格式再进原管线，override/fallback 叠加、家族匹配、DeepSeek 策略全部原样生效；回滚只需把 `remote_url` 配回旧值。
- **全量 diff 后确认三处字段缺口并接受**：models.dev 无 1h 缓存写入价（Claude 按 5m 价计费，少收一半）、无 priority 档价（priority 流量按标准价）、部分图片模型无价（gpt-image-1 cost 为空）。缓解：内置 212 条 fallback 价卡继续按层合并，覆盖上述模型的关键字段；后续可经 override 文件补。
- **非整千 tier 跳过**：真实数据存在 size=200001 的 tier，LiteLLM 字段名只支持千粒度，打日志跳过（6 处）。
- **providerQuotes 曾被 fallback 合并清空**：fallback/override 以 LiteLLM 格式重入 `parsePricingData`，defer 无条件覆盖索引为 nil；改为仅当本次输入确为 models.dev 格式时更新，并加回归测试钉住。
- **8080 重启**：经用户确认后重启本地服务（旧进程 10-22 启动，无 .installed 锁，需 `DATA_DIR=./data` 环境启动），旧 LiteLLM 目录快照移走后启动时自动从 models.dev 拉取成功。

**（下午）**

- **升 12 的理由按优先级**：① pnpm 9 已在 CVE 修复窗口外（CVE-2026-55698 CVSS 8.8 只修 10.34.2+/11.5.3+；CVE-2026-24056 只修 10.28.2+）——对本 fork 危险在于威胁模型包含上游 merge；② pnpm 10-12 默认供应链防护（脚本拦截、minimumReleaseAge 1d、blockExoticSubdeps、完整性不匹配硬失败、lockfile 策略校验）；③ 与仓库安全姿态一致。
- **工作区旧状态不保留**：下午会话开始时 lockfile 有 445 行 diff（pnpm 12 隐式 install 摘掉了全部 overrides 钉子）+ 占位 workspace.yaml，属于安全倒退状态，重做而非提交。
- **本地 npmmirror 无 audit 端点**：仅影响本地直跑 `pnpm audit`，验证时用 `--registry=https://registry.npmjs.org` 绕过；CI 用官方 registry 不受影响。
- **放弃留在 pnpm 9**：9 无补丁计划，钉版本只能防依赖层防不了包管理器本体。

### 验证

**（上午）**

- 后端：`go test ./internal/service/`（138s 全绿）、`./internal/handler/...`、`./internal/server/routes/`；`go vet` 通过；`make build`（前端因网络换 pnpm 直跑 `vite build` 成功）。
- 转换正确性：真实 api.json 全量冒烟（3444→3542 模型、单位/去重/阶梯/图片映射断言）；与上游 model-price-repo 全量 diff 记录 52 处分歧（图片 token 价口径差 + deepseek 两边价不同）与 108 个缺失模型（vertex 变体为主），据此确认 fallback 层必要性。
- 真实环境端到端：本地服务从 models.dev 拉取成功（日志 `Downloaded 3542 models`，对比旧目录 262）；`provider-quotes?model=gpt-5.5` 返回 20 家厂商。
- 浏览器实测（暗色主题）：块网格渲染、弹窗开合（Escape 关闭）、长上下文阶梯弹窗（≤100K/>100K 两档）、搜索过滤、管理员比价区（openai 带计费采用徽章）；非管理员无比价区（组件测试覆盖）。

**（下午）**

- 安全钉子核对：lockfile 中 js-cookie@3.0.7、dompurify@3.4.17、postcss@8.5.29、source-map-js@1.2.2、form-data@4.0.6、follow-redirects@1.16.1 全部命中。
- `pnpm peers check`：仅 @emoji-mart/react 声明滞后一条（react 16-18 vs 实装 19），升级前已存在，非迁移引入。
- audit 例外工具兼容性：pnpm 12 `audit --json` 输出喂 `tools/check_pnpm_audit_exceptions.py`（python3.12）通过，最终仅剩 xlsx 两条已例外项。
- 前端 `make test-frontend`：lint + typecheck + 30 文件 574 用例全绿；`make build` 前端构建嵌入 `internal/web/dist` 成功（12.7s）。
- 后端 `go test ./...` 全绿（pnpm 迁移与 Go 代码无交集）；golangci-lint 本地不可用（本地二进制不支持 Go 1.27 export data），由 CI 门禁跑。
- 未运行：Docker 构建端到端验证（corepack prepare pnpm@12.15.0 行，由下次 release workflow 或手动 docker build 兜底）。
- 意外发现与处置：① 上午 `f95a86a22` 的 CI 实为红（golangci-lint gofmt），`e88cc317` 修复；② security-scan workflow 此前从未在 push 上触发过，本次首跑即暴露 go1.27.2 标准库漏洞，`21c0a2cf` 升 1.27.3 后绿。
- Docker 构建端到端（原「未运行」项补跑）：corepack 拉 pnpm 12.15.0 成功，但实测暴露两处 Dockerfile 问题——COPY 漏 `pnpm-workspace.yaml`（frozen 安装报 ERR_PNPM_LOCKFILE_CONFIG_MISMATCH）与 1536MB 堆 OOM；`1ee7db428` 修复后三阶段构建全通过，容器启动冒烟正常（v0.2.17、前端嵌入、go1.27.3）。
- **本地数据目录回归上游形态**：`backend/data/` 整体上移一级到 `backend/`（config.yaml/model_pricing.json+sha256/logs/pages/plugins），裸跑 `go run`/`./bin/server` 不再需要 `DATA_DIR`，严格贴合上游「config 与启动目录同层」与 Docker `/app/data` 两种形态；`ad842f9b6` 补 gitignore（backend/logs|pages|plugins|model_pricing*）与 dockerignore（backend/config.yaml、.installed、运行时数据）确保 dev 配置进不了 git 与镜像上下文（探针镜像实测 `COPY backend/ ./` 上下文无 config）。dev config.yaml 补 `log.output.file_path: ./logs/sub2api.log` 消除上游对非容器环境硬编码 `/app/data` 日志路径的写失败降级警告（属 dev 配置，不动代码）。
- dev 配置进不了发布产物的三条结构性确认：本地 docker build 靠 .dockerignore（实测）；CI release 镜像走 goreleaser 净上下文（归档裸二进制 + LICENSE/README/deploy 白名单，结构性排除）；二进制 `go:embed` 只嵌 VERSION 与前端 dist，config.yaml 运行时读取。
- 8080 服务重启至新布局：旧实例（10-22 启动）kill 后以 `./bin/server` 重启，/health 200、日志落 `backend/logs/sub2api.log`、价格缓存读 `backend/model_pricing.json`；确认后删除旧 `backend/data/` 目录。

### 产出

- `641b4aed` feat(pricing): switch catalog source to models.dev api.json
- `102a477c` feat(plaza): redesign user model plaza with block grid and price dialog（CI 绿）
- `f95a86a22` feat(plaza): add admin-only provider comparison to model dialog
- `e15d33ae` build(frontend): upgrade pnpm 9 to 12 with supply-chain hardening
- `e88cc317` style(service): fix gofmt alignment in pricing_modelsdev.go（上午 models.dev 提交漏跑 gofmt，CI golangci-lint 门禁红，补齐两处 struct tag 对齐后 CI 绿）
- `21c0a2cf` chore(deps): bump go 1.27.2 to 1.27.3 for crypto/x509 security fix（security-scan 首次在本仓库触发 govulncheck 门禁即报标准库漏洞，passkey 登录路径可达；go.mod + 4 处镜像/CI 版本钉同步升级；最终 CI 与 Security Scan 双绿）
- `1ee7db42` fix(docker): copy pnpm-workspace.yaml and raise node heap for frontend build
- 下午收尾后 `origin/main` tip：`1ee7db42`，CI（backend-ci）与 security-scan 全绿，Docker 三阶段构建实测通过

---

## 2026年10月23日

**主题**：初始化管理员账号默认值改为 `admin@admin.local` / `admin123`

### 改动

- `backend/internal/setup/setup.go`：新增 `defaultAdminEmail`/`defaultAdminPassword` 常量；`prepareAdminCredentials` 在邮箱/密码为空时回退到固定默认值，替代上游的随机账号；删除 `generateAdminEmail`；bootstrap 打印文案从「随机生成/一次性」改为「默认值 + 请改密」。
- `backend/internal/setup/cli.go`：CLI 向导默认邮箱提示改用常量。
- `backend/internal/setup/setup_test.go`：两处「随机邮箱正则」断言改为「固定默认值」断言。
- 部署配置同步：4 个 `docker-compose*.yml` 的 `ADMIN_EMAIL`/`ADMIN_PASSWORD` 默认值、`.env.example`、`docker-deploy.sh`（去掉随机邮箱生成）、`deploy/README.md`、`APPLE_CONTAINER.md`。
- 本地开发库（Docker postgres）：把既有 `admin@codeapi.local` 更新为 `admin@admin.local`，`password_hash` 重置为 `admin123` 的 bcrypt 哈希，并同步 `auth_identities` 的 email 身份 `provider_subject`。

### 关键决策

- 弱化上游「随机账号防爆破」设计（issue #7850），换取开发/内网便利。这是 fork 二开的显式选择，经用户确认；默认值属已知弱凭据，公网部署需装后改密或经 `ADMIN_EMAIL`/`ADMIN_PASSWORD` 覆盖，已写进代码注释与 PR 描述。
- `config.example.yaml` 里的 `default.admin_email`/`admin_password` 仍是上游死配置（无消费方），本次未动，避免与上游 sync 冲突。

### 验证

- `cd backend && go build ./...` 通过。
- `cd backend && go test ./internal/setup/...` 通过。
- `bash -n deploy/docker-deploy.sh` 通过；`git diff --check` 通过。
- 本地库哈希用 `bcrypt.CompareHashAndPassword` 实测 `admin123` 匹配（true）。

### 产出

- `1ec75605f` feat(setup): default admin account to admin@admin.local/admin123（PR #7，合并提交 `cb62585c0`）

**主题（同日第二次工作）**：流程从「PR」简化为「直接 push main」

### 改动

- `AGENTS.md`：把「上传与合并流程（定版）」改写为「提交与推送流程（定版）」，移除开分支、PR、CLA、merge commit 等环节；同步更新 Repository Identity、General Working Rules、User Intent Mapping、Merge After Confirmation、Release、Work Log 各节的 PR/分支引用。
- `docs/FORK_MAINTENANCE_CN.md`：分支表、流向图、确认合并、日常开发落地、推送结果、失败恢复、工作日志、提交约定各节同步去 PR 化。

### 关键决策

- 单人 fork、无预发环境，PR 的 review/merge 与 CLA 机器人是纯 overhead；CI（backend-ci、security-scan）由 push 触发仍生效，故移除 PR 环节，改为直接 commit + push `origin main`。
- 上游同步仍保留 `sync/upstream-*` 分支做冲突处理与审计，只是合并后直接 push，不再开 PR。

### 验证

- `git diff --check` 通过（纯文档改动，无代码变更）。

### 产出

- `4c89c1de0` docs(worklog): log 2026-10-23 default admin account change
- `d6ee4bf21` chore(workflow): drop PR flow, commit directly to main

---

## 2026年10月22日

**主题**：修复首页版心丢失与真实延迟面板被移除

### 改动

- `HomeView.vue`：给顶栏、Hero、接入、模型、CTA、页脚恢复 `mx-auto max-w-7xl`，修复后续一版改动把页面改成通栏导致的「无主题」观感。
- 重新接入真实延迟实测：恢复探针脚本（5 秒一次、Range 1 KB、隐藏页暂停）、Hero 实时往返状态与完整延迟面板；同时保留后续一版的真实端点 / 真实价目面板，并把实时往返条并入 Hero 面板。
- 延迟面板恢复当前 / 平均 / P95 / 抖动 / 样本数、SVG 折线、TTFB、连接复用、传输量、协议和最近样本。
- zh/en `landing.ts` 恢复 `home.nav.latency`、`home.hero.probeValue/probePending` 与完整 `home.latency` 文案；保留 `home.panel` / `home.plaza` 键。
- 移动端顶栏隐藏「Key 用量查询」按钮，避免品牌名被截断（该入口仍保留在 Hero 与页脚）。

### 关键决策

- 不回滚后续一版「真实端点 + 真实价目」的方向：只恢复版心与真实延迟，不重新引入带示例数字的控制台预览。
- 延迟仍按当前设备 → 本站静态资源的真实往返计算，价目取不到时保持「价目未开放」，不补假数据。
- 折线继续从 0 起算、最小量程 50 ms，避免几毫秒波动被画成陡峭趋势。

### 验证

- `npx pnpm@9 run typecheck` rc=0；`npx pnpm@9 run check:i18n` 3 tests 通过。
- `npx pnpm@9 run test:run` 371 文件 / 2868 用例全部通过；`npx pnpm@9 run build` 通过。
- 重新 `go build -tags=embed` 并重启 8080（PID 66960），`/health` 返回 ok，入口资源与 `backend/internal/web/dist` 一致（`assets/index-svF_cZ-k.js`）。
- CDP 实测：1440 宽 `scrollWidth=1440`，390 宽 `scrollWidth=390`；`#latency` 存在，延迟面板采样 3 次、折线已绘制；顶栏在 390 宽显示完整 `CodeAPI`。
- 未运行后端测试：本次只改前端视图与文案。

### 产出

- 未提交（工作区改动待确认）：`HomeView.vue`、`landing.ts`×2（README 的既有改动未触碰）


**主题（同日第二次工作）**：首页收敛为 Gateway Ledger v2，模型信息改为登录后模型广场可见

### 改动

- 重写默认宣传首页 `HomeView.vue`：删除模型覆盖/价目/路由行，公开页只保留协议、端点、真实延迟与登录入口。
- Hero 右侧 Gateway 面板改为「实时延迟 / 协议 / 计费」三指标 + 端点列表 + 「模型与价格 → 登录查看」+ 真实迷你折线。
- 所有 CTA 改为登录/控制台语义；移除「注册」「开放注册」「注册并开始」与充值文案。
- 公开文案去掉厂商枚举：`home.hero.description`、`home.facts.models.value`、CTA 说明统一改为“多平台 / 模型广场查看”。
- 接入示例中的 `model` 从真实模型名改为 `your-model`，公开页不泄露模型信息。
- 默认首页不再渲染 `/model-plaza` 链接；更新 `HomeView.compact.spec.ts` 对应断言。
- 登录页复用同一组 i18n，左栏说明与事实条同步收敛，无厂商名。
- 页脚删除版本号与右侧链接组，只保留版权行。
- 首页移除语言切换控件，默认中文；主题按钮移动到登录按钮右侧。
- 登录页顶栏同步移除语言切换控件，保留主题切换与「返回首页」。

### 关键决策

- 模型信息（模型名/厂商/价格/倍率）只允许登录后在模型广场访问；公开首页零模型信息，不请求 `/api/v1/model-plaza`。
- 注册保持关闭，首页只出现“登录”；游客主 CTA 固定 `/login`，已登录走对应控制台。
- 真实延迟保留为公开页的差异化展示；没有 uptime 监控就不写 uptime。
- 玻璃只保留顶栏、Gateway 面板与事实条；其余用实心面板和细线，避免“整页玻璃”。

### 验证

- `npx pnpm@9 run typecheck` rc=0；`npx pnpm@9 run check:i18n` 3 tests 通过。
- `npx pnpm@9 run test:run` 371 文件 / 2868 用例全部通过；`npx pnpm@9 run build` 通过。
- 重新 `go build -tags=embed` 并重启 8080（PID 71457），`/health` 返回 ok，入口资源 `assets/index-DikISsyc.js` 与 dist 一致。
- CDP 实测 8080：DOM 中 `Claude/GPT/Gemini/Antigravity` 命中为 false，`/model-plaza` 链接不存在；页面链接仅 `#integration`、`#latency`、`/key-usage`、`/login`；顶栏无语言切换，右侧顺序为「Key 用量 → 登录 → 主题」；页脚仅版权行、无版本号与链接；延迟采样 3 次；1440/390 无横向溢出。
- 未运行后端测试：本次只改前端视图与文案。

### 产出

- 未提交（工作区改动待确认）：`HomeView.vue`、`landing.ts`×2、`HomeView.compact.spec.ts`

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
