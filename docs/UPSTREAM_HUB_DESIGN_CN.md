# 上游聚合（Upstream Hub）设计文档

> 状态：MVP 已落地（2026-10-26，`bdb49a5af`）并经多轮交互迭代收口（至 `e03e5a205`，区间内 15 个 upstream-hub 提交）。
> 本文由仓库外 `sub2api-demo/UPSTREAM_HUB_PLAN.md`（规划稿 v6）提炼而来，按当前代码实况修订；规划期各版本（v1–v6）的演进痕迹已去除，只保留仍然有效的决策与尚未实施的设计。
> 意图路由为独立待定项，见 [PENDING_INTENT_ROUTING_CN.md](./PENDING_INTENT_ROUTING_CN.md)。

## 一、定位

上游聚合解决的问题：智能体对接的上游分散在多家平台（含国产 Kimi / 智谱 GLM / DeepSeek / MiniMax），没有统一分组、配置分散、模型名不统一。

**定位**：与「账号管理」同源的功能——账号管理维护的是厂商凭据（订阅账号、OAuth），上游聚合维护的是**上游厂商**（透传 Base URL + Key 的供货渠道）；两者原理相同（凭据存储、健康检查、调度），只是维护对象与引入的聚合技术不同。

**核心结论：网关不重构，功能主体是编排层。** sub2api 已具备全部中转原语（`type=apikey` 自定义 base_url 账号、composite 分组、模型路由、协议翻译矩阵、failover/冷却），上游聚合 = 现有组件的受控组合 + 一张厂家台账新表。

### 心智模型：一个数据、两个入口、三个视角

- **唯一数据**：Account 行（accounts 表）。「账号管理」与「上游聚合」不是两套数据，是同一张表的两个管理入口。
- **归属判定一条规则**：有 `extra.upstream_vendor_id` → 由上游聚合管理；没有 → 由账号管理管理。分组挂载（account_groups）与来源入口无关。
- **三个视角**：供货（聚合：厂家预算/协议/条目）→ 分配（分组：挂载谁、本组是否启用）→ 消费（客户：面板自建 key）。

## 二、复用的现有原语（已核实）

| 能力 | 现有原语 | 位置 |
|---|---|---|
| 上游透传账号 | `Account(type=apikey)` + `credentials.base_url`（Base URL + API Key 直连上游），网关全部转发路径都支持；`type=upstream` 仅 antigravity relay 有实现，hub 不使用 | `backend/internal/domain/constants.go` |
| 统一分组容器 | `Group.platform = "composite"` | `backend/ent/schema/group.go` |
| 模型归一表 | `CompositeModelRoute`：public_model → target_platform + upstream_model，exact/prefix、endpoint、priority | `backend/ent/schema/composite_model_route.go` |
| 国产厂商登记 | `PlatformSpec.CNProvider`（Kimi/智谱/DeepSeek/MiniMax），平台清单是单一真相 | `backend/internal/domain/platforms.go` |
| 账号调度 | `Account.concurrency / load_factor / priority / schedulable / temp_unschedulable_until/reason` | `backend/ent/schema/account.go` |
| 顺序路由钩子 | `Group.model_routing` + `model_routing_enabled`，调度期已消费（load-aware、会话粘性） | `backend/internal/service/gateway_scheduling.go` |
| 模型解析 | `CompositeRouteResolver.Resolve`：显式路由 → 账号归属 → 模型名探测 | `backend/internal/service/composite_route_resolver.go` |
| 分组限流 | `Group.rpm_limit`（0=不限，优先于用户级兜底） | `backend/ent/schema/group.go` |
| 失败退避 | 全局「过载冷却 / 429 冷却」设置（OverloadCooldown / RateLimit429Cooldown） | admin settings |
| 对外接入 | `APIKey` + 现有 `/v1/chat/completions`、`/v1/messages` 等端点 | 无新增端点 |

## 三、数据模型（已实现）

### 厂家台账：`upstream_vendors` 表

`backend/ent/schema/upstream_vendor.go`，迁移 `backend/migrations/900_xeu_upstream_vendors.sql`：

- `name`：唯一（软删感知的部分唯一索引 `WHERE deleted_at IS NULL`）
- `concurrency_budget`：账号级并发预算（上游厂商对本账号的总并发上限，如 50）；0 = 不启用共享限制
- `base_urls`：厂家的候选接入地址（jsonb 数组，最多 10 个）。条目保存时实测各地址延迟，只把最快的一个写进条目凭证；运行期不再切换。
- `notes`、软删（TimeMixin + SoftDeleteMixin）
- 删除保护：厂家下仍有条目时拒绝删除（`ErrUpstreamVendorNotEmpty`）

### 条目：`Account(type=apikey)` + `extra.upstream_vendor_id`

一条上游记录 = 一行 Account，不建新实体。**hub 界面不提供任何代理配置**：需要代理时在 IP 管理创建并绑定到账号（平台能力，`account.proxy_id`）。归属判定只读 `extra.upstream_vendor_id`（兼容 int64/int/float64/string 反序列化），**与账号 type 无关**。条目创建复用 `AdminService.CreateAccount` 全套校验，`SkipDefaultGroupBind` 避免自动挂默认组。

> **类型选择（2026-10-28 修正）**：条目最初建成 `type=upstream`，实测不可用——网关凭据解析（`OpenAIGatewayService.GetAccessToken`、`GatewayService.GetAccessToken`）以及全仓 20+ 处 `switch account.Type` 都不认 `upstream`，会落 `default:` 报 `unsupported account type`。改用 `type=apikey` + 自定义 `base_url`：这是平台日常在跑的国产厂商手动 Key 形态，所有转发路径天然支持，且 `GetOpenAIBaseURL`/`GetOpenAIApiKey` 本就对 upstream/apikey 同口径处理。历史条目由迁移 `900_xeu_upstream_entry_type_apikey.sql` 转换（只改带 `extra.upstream_vendor_id` 且未软删的行，不碰 antigravity relay）。

**一 key 一协议**（管理员确认的上游现实）：一条目只对应一种接口类型。协议 → 落库 platform → 网关端点：

| 协议类型（界面） | platform 落库值 | 对外网关端点 |
|---|---|---|
| OpenAI 兼容 | `openai` | `/v1/chat/completions` |
| OpenAI Responses | `openai`（+ responses 能力） | `/v1/responses` |
| Anthropic 兼容 | `anthropic` | `/v1/messages` |
| Gemini 兼容 | `gemini` | `/v1beta/models/...` |
| System One | `typesafe` | `/v1/systemone` |

### 分组引用：`account_groups.enabled`

迁移 `backend/migrations/900_xeu_account_group_enabled.sql`（additive，默认 true）：

- `enabled=false` =「暂不调用」，**仅本分组**候选池剔除；全局停用仍走 `Account.schedulable`，两级独立并存。
- 调度侧：`queryAccountsByGroup` 在可调度路径追加 `account_groups.enabled=true` 过滤（`backend/internal/repository/account_repo.go`）。
- 统计口径：分组账号统计 SQL 同步排除本组停用条目（`backend/internal/repository/group_repo.go`），避免"显示 3 可用实际 2"。

### 厂家预算语义

- 预算是**厂家级、全局跨分组**的：同厂家所有条目的并发限额之和 ≤ 预算。
- 录入时刻实时提示（弹窗底部"已分配/预算"实时求和，含正在编辑的行），超限**硬阻断保存**：红色横幅 + 保存按钮禁用 + 点击拦截，不能带病退出弹窗。
- 同厂家多 key 不冲突：各自独立 Account 行，靠条目名区分（建议带协议后缀，如"A-OpenAI"、"A-Claude"）。

## 四、后端实现（已上线）

| 层 | 文件 | 职责 |
|---|---|---|
| Service | `backend/internal/service/upstream_hub.go` | 厂家 CRUD、条目 CRUD + 复制（服务端拷贝凭据）、预算 enrich（EntryCount / AllocatedConcurrency / OverBudget）、分组绑定（平台匹配校验）、两级开关 |
| Repo | `backend/internal/repository/upstream_vendor_repo.go` | 厂家持久化；条目按 `extra.upstream_vendor_id`（jsonb 谓词，数字/字符串两种形态）过滤，与 type 无关；分组绑定 tx 内 delete + CreateBulk |
| Handler | `backend/internal/handler/admin/upstream_hub_handler.go` | admin API，NotFound/Conflict 错误映射 |
| 路由 | `backend/internal/server/routes/admin.go` | 见下表 |

API 一览（全部 admin 鉴权）：

```
GET/POST   /admin/upstream-vendors
GET/PUT/DELETE /admin/upstream-vendors/:id
POST       /admin/upstream-vendors/:id/entries
GET/PUT/DELETE /admin/upstream-entries/:id
POST       /admin/upstream-entries/:id/duplicate
PUT        /admin/upstream-entries/:id/schedulable
PUT        /admin/upstream-entries/:id/groups
PUT        /admin/upstream-entries/:id/groups/:groupId/enabled
GET        /admin/upstream-groups
```

DI 走 google/wire（`cmd/server/wire_gen.go` 手工补丁，本机无 wire CLI 时按既有 provider 链手接）。

单测：`backend/internal/service/upstream_hub_test.go`（extra 归属解析、BaseURL 归一、预算 enrich/超限判定、厂家校验、删除保护）。

## 五、前端实现（已上线）

- 页面：`frontend/src/views/admin/UpstreamHubView.vue`，路由 `/admin/upstream-hub`，侧边栏「上游聚合」（cube 图标）。
- API 封装：`frontend/src/api/admin/upstreamHub.ts`（注册进 `api/admin/index.ts`）。
- i18n：`locales/{zh,en}/admin/upstreamHub.ts`，`nav.upstreamHub` 进两个 `common.ts`。

### 交互定稿（多轮实测迭代后的最终形态）

- **厂家列表**为主体（两层结构：厂家 → 条目；无"聚合列表"层）。管理条目为 toggle 式按钮，选中厂家高亮；详情卡无挂载动画（避免布局跳动）。
- **添加/编辑统一为一个弹窗**（`BaseDialog width="max"`，95vw，为此给 DialogWidth 加了 `max` 档）：
  - 一 key 一行表格布局：名称 / 平台（既有条目锁定不可改）/ Base URL / API Key（编辑时留空=保持原值）/ 并发 / 优先级 / 分组 / 操作（复制+删除），底部「+ 添加 Key」。
  - 复制既有条目走服务端 `duplicate`（凭证不出后端明文回显路径）。
  - 分组选择为**长条下拉**（非圆角药丸）：固定定位面板 + 背板（避免被表格 overflow 裁剪）；**整行点击选中，选中态=底色高亮**，行内无复选框无开关。
  - 超预算硬阻断见「三、厂家预算语义」。
- **详情表分组芯片（只读）**：只展示状态——启用=青柠底、停用=灰底并附「本组暂不调用」；**点击无任何动作**。挂载关系在添加/编辑弹窗维护，本单元启停统一在聚合调度页维护（2026-10-29 修正：状态展示位不应承担写入操作）。
- **全局可调度开关**：条目行内 toggle → `PUT /:id/schedulable`。

## 六、分组管理对接（供货层 ⇄ 客户面）

### 架构分工

```
上游聚合（供货层·管理方）           分组管理（客户面）                    客户
厂家台账 + 预算                     分组（platform = 模型厂类别）          面板用户
 中转商A（预算50）                   ├ deepseek 组 ← 挂载条目              按 user_allowed_groups
  ├ A-OpenAI  [openai类别]    ──►   │   ├ DeepSeek-V3（启用）             获授权 → 面板自建 key
  ├ A-DS      [deepseek类别]        │   └ 中转商A-DS（暂不调用）           → 智能体接入
  └ A-Claude  [anthropic类别]       ├ kimi 组 / grok 组 / openai 组 …
```

### 已确认的四个决策

1. **条目入组 = 手动勾选**：按分组的模型厂类别（现有 Group.platform 枚举）过滤条目清单勾选（走现有 account_groups M2M，带 priority）。**当前实现从聚合侧（条目弹窗的分组下拉）完成**；GroupsView 建组时的「从上游聚合引用」区块为待建 UI（机制已通）。
2. **客户 key = 复用现有体系**：客户即面板用户，user_allowed_groups 授权，面板自建 api_keys 绑分组——用量/限额/审计全继承，零新代码。
3. **暂不调用 = 仅本分组**：`account_groups.enabled`，见上。
4. **一条目一类别**：中转商 openai key 若还能转发 DeepSeek，就再建一条 A-DS [deepseek 类别]。类别即调度池，边界清晰。

### 「不影响原账号管理」的双护栏

- **双头管理防冲突**（规划）：聚合条目在「账号管理」列表打"上游聚合"标记并只读——**未实现**，见已知缺口。
- **调度零侵入**（已实现）：条目入组后走现有分组调度（同池同优先级同退避），聚合不另起调度路径，暂不调用只是候选池过滤。

### 建组双来源（规划，GroupsView 侧未建）

1. 来源一「从账号管理选择」：现有挂载 UI/逻辑原样保留，只列非聚合账号。
2. 来源二「从上游聚合引用」（新增）：按类别过滤厂家条目勾选；引用后分组页可逐 key「暂不调用」。
3. 组内列表统一展示：来源徽章区分；仅聚合来源行有「本组调用」开关。
4. 聚合侧回看：条目列表「被引用分组」列已实现（详情表分组色条）。

## 七、网关适配边界

**结论：网关不需要重构。** magpie 三核心行为（候选循环换号、退避休眠、会话亲和）在现有网关均已实现且更完善（流式保护、池模式同账号重试、prompt cache 透传），重构=推翻分组调度模型=触碰二开红线，负收益。

**三个接缝（全部为编排层工作，不动网关核心）：**

1. **平台疆界**：现有调度按 Group.platform 划池；聚合条目自带协议属性，通过挂分组解决（composite 天生跨平台）。
2. **复合解析优先级**：归一模型表作为 CompositeModelRoute"显式路由"喂给现有解析器（显式路由 > 账号归属 > 模型名探测），解析器零改动。
3. **供应商预算槽位**（唯一真正的网关接触面，**未实现**）：双槽位获取（先占厂家预算槽、再占条目并发槽）挂在选号成功之后、转发之前；实现在独立薄拦截层，不进 selectAccountForModel* 内部；切面位置必须保证不绕过会话亲和与计费链路。当前预算只在**录入时刻**做硬校验，运行期不占用槽位。厂家级联动冷却（同厂家任一条目触发上游 429 → 全厂家条目一起退避）同属此接缝，未实现。

**明确放弃的复刻点**：magpie 的单机多 Agent 配置改写（桌面工具形态）。取其路由语义（路由组=模型池、策略排序、失败降级、会话亲和），弃其桌面形态。

### 跨协议翻译矩阵（已核实）

| 入站 \ 上游 | OpenAI 兼容 | Anthropic 兼容 | Gemini | Responses | System One |
|---|---|---|---|---|---|
| /v1/messages | ForwardAsChatCompletions | 原生 | — | ForwardAsResponses | systemone 端点 |
| /v1/chat/completions | 原生 | ForwardAsAnthropic | — | 原生端点 | — |
| /v1beta | — | GeminiMessagesCompat | 原生 | — | — |

**连通性状态**：2026-10-28 用本地 mock 上游（OpenAI 兼容 `/v1/models` + `/v1/chat/completions`）完成端到端实测：hub 条目 → 探测模型 → 写入 `model_mapping` → 挂入 composite 分组 → 客户 key 调 `/v1/models` 返回上游真实模型、`/v1/chat/completions` 正常转发。**真实第三方上游的协议覆盖仍需管理员用真实 Key 复测**（mock 只证明路由与凭据链路，不证明各家上游的兼容差异）。

### cc-switch 设计参考

- 聚合模式：多供应商模型汇入统一列表，按所选模型路由；不做故障转移。
- 路由模式：故障转移队列 P1..Pn + 熔断器 + 按应用重试/超时。
- 映射：cc-switch"聚合"≈ 归一模型表；"故障转移"≈ priority_fallback。形态不同（单机桌面 vs 服务端网关），路由语义同源。

## 八、探活三级设计（规划，未实现）

复用 channel_monitor，控开销：

| 档位 | 机制 | 计费 | 用途 |
|---|---|---|---|
| 被动（默认） | 真实流量的失败/429 冷却状态（temp_unschedulable） | 零 | 日常存活判定 |
| quota | 读上游用量/额度接口（check_mode=quota 同款） | 零 | 定期巡检 |
| probe | 算术挑战探活（channel_monitor_challenge 同款，token 极小） | 每次≈最小计费单位 | 手动测试/低频巡检 |

不使用账号测试的全量请求（完整系统提示词 + max_tokens=1024，贵）做定期探活。探活频率与档位可配，默认被动 + 每日 quota。

## 九、当前状态与已知缺口（截至 2026-10-28）

**已上线可用**（2026-10-28 更新）：厂家台账（预算/超限硬阻断）、条目增删改复制、并发+优先级、分组绑定（含 composite 综合分组）与每分组「暂不调用」、全局可调度开关、调度与统计尊重 enabled、`/admin/upstream-hub` 与 `/admin/aggregation` 两个界面（zh/en）、**上游模型探测 + 勾选落库**（折叠面板）、归一模型表（聚合调度页）。

**已知缺口**：

1. **真实第三方上游连通性未复测**——mock 上游已打通全链路（§七），等待真实 Base URL + Key 验证各厂商协议差异。
2. GroupsView 建组「从上游聚合引用」区块未建（机制已从聚合侧打通，聚合调度页也可挂载）。
3. 条目平台创建后不可改（弹窗中既有条目平台锁定）。
4. 复制条目到其它厂家/分组时 enabled 重置为 true。
5. `/v1/models` 未过滤"本分组暂不调用"的条目。
6. 账号管理页对聚合条目无只读标记（双头管理护栏未落；条目现在与普通 apikey 账号同型，更需要这个标记）。
7. 运行期厂家预算槽位与厂家级联动冷却未实现（§七接缝 3）。
8. 探活三级未实现（§八）。
9. 意图路由整体待定（见 [PENDING_INTENT_ROUTING_CN.md](./PENDING_INTENT_ROUTING_CN.md)）。
10. ~~条目无 `model_mapping` → 模型列表链路全断~~ → **已解决**（§十一 P1）：探测 + 勾选写 `model_mapping`，`/v1/models` 与 composite 归属解析随之生效。
11. ~~hub 的分组平台匹配严于平台自身规则~~ → **已解决**（§十一 P2）：改为「平台相等 或 composite」，前端下拉同步放开。
12. ~~API 密钥页「API 端点」块随空设置消失~~ → **已解决**（§十一 P0）：前端兜底 `window.location.origin`。

**后续实施顺序**：真实上游复测 → GroupsView 双来源 UI → 探活三级 + 预算槽位接缝 → 意图路由（待定）。

## 十、风险

1. `type=apikey` 的 hub 条目跨协议覆盖面需实测（openai / anthropic / gemini 三协议各连一家真实上游）；mock 上游已证明路由与凭据链路。
2. `Group.model_routing` 与 `CompositeModelRoute.priority` 的一致性由界面层写入顺序保证（先 routes 后 model_routing），需在测试中覆盖原子性。
3. 意图路由相关的延迟/计费/时序风险整体移至待定项。

## 十一、实施计划（2026-10-28 定稿）

三个管理员实测问题定位后确定的修复方案。决策已锁定：端点 UI 加前端兜底；探测结果**勾选确认后**再写入（不做一键全量，避免中转商上百个杂模型灌进客户列表）；聚合调度做成**独立新页** `/admin/aggregation`（不扩 GroupsView 的 composite 页签，也不塞回上游聚合页）。

### P0 · API 端点 UI 兜底（独立，先走）

| 改动 | 位置 |
|---|---|
| v-if 与传参改为 `api_base_url \|\| window.location.origin`，端点块不再因空设置消失 | `frontend/src/views/user/KeysView.vue:26-29` |
| 同口径修 `UseKeyModal` 的 `:base-url`（当前传空串） | `frontend/src/views/user/KeysView.vue:1076` |
| 正式域名仍需在 系统设置 → 站点设置 填写（兜底只保证不空白） | 管理端设置 |

### P1 · 探测上游模型（核心，同时是 P2 的前提）

复用平台已有的探测能力，不新造：`AccountTestService.FetchUpstreamSupportedModels` / `SyncUpstreamModelCatalog`（`internal/service/upstream_models.go`），已覆盖 openai、多协议 API Key 供应商（国产厂商 + 聚合平台）、gemini、anthropic、grok、antigravity（`buildUpstreamModelsRequest:904`）。管理端现成端点：`POST /admin/accounts/models/sync-upstream-preview`（不落库）、`POST /admin/accounts/:id/models/sync-upstream`。

**后端**

1. ~~放宽三处 type 守卫接受 `AccountTypeUpstream`~~ → **实施时改道**：实测发现 `type=upstream` 在网关 20+ 处 `switch account.Type` 都不被识别（`GetAccessToken` 直接报 unsupported），逐个放宽风险过大。改为条目一律 `type=apikey` + 自定义 `base_url`，并迁移历史条目；探测链对 apikey 全平台已内置支持，`upstream_models.go` 零改动。归属仍只认 `extra.upstream_vendor_id`。
2. hub 新增两个端点：
   - `POST /admin/upstream-entries/:id/probe-models` → 返回 `{models, warnings}`，**不落库**
   - `PUT /admin/upstream-entries/:id/models` body `{models:[...]}` → 写 `credentials.model_mapping`（from=to），经 `AdminService.UpdateAccount`；已确认 `admin_account.go:653` 是 `MergePreservingSensitiveCreds` **合并**语义，不会抹掉 `api_key`/`base_url`（仓库已有「Credentials 只含 model_mapping」先例：spark shadow，`admin_account.go:1445`）
   - 写完调 `InvalidateAvailableModelsCache(groupID, platform)`（`gateway_service.go:1601`），否则短 TTL 缓存仍返回旧列表
3. 未保存的行复用现成 `sync-upstream-preview`，`type` 传 `upstream`（守卫放宽后是合法值，无需伪装 apikey）。
4. 单测：三平台 × upstream 类型的守卫放宽；probe 不落库；models 写入是合并非替换。

**前端**（`UpstreamHubView.vue`）

- 每行下方**折叠态**面板，默认收起；行尾显示「已登记 N 个模型」。
- 探测结果为**按钮块**（flex-wrap）：点击即启用/停用（选中=青柠实底），不使用复选框。
- 展开后：探测按钮 → 模型清单带勾选（默认全不选）→ 全选/反选 → 确认写回该行；保存条目时随 `PUT /:id/models` 落库，新行先建条目再写模型。
- 结果带**来源徽章**（上游实时 / 公共目录候选）与可读告警；上游无模型接口时支持**手动添加模型 ID**（cc-switch 同款）。
- 探测失败按 `UpstreamModelSyncError` 展示可读原因（含底层传输错误），**不阻塞保存**。
- UI 口味沿用：无展开动画、整行可点勾选、不用圆角药丸。

**收益**：`/v1/models` 返回上游真实模型；composite 归属判据有数据可依 → P2 跨平台路由才真正可用。

### P2 · 综合分组 = 新建 `/admin/aggregation`

1. **放宽平台匹配**：`upstream_hub.go:397` 由严格相等改为 `相等 || group.Platform == composite`（对齐平台自身规则 `admin_group.go:341`，非弱化校验）；前端 `UpstreamHubView.vue:569` 下拉同步放开，composite 分组打「综合」徽章。
2. **新页面三块**：
   - 供给汇入：跨平台勾选厂家条目进综合分组，显示平台/并发/优先级/本组启用
   - 归一模型表：public_model → 目标平台 + 上游模型 + 优先级链，直接读写现有 `CompositeModelRoute`；模型候选来自 P1 探测结果
   - 链路视图：厂家 → 条目 → 分组 → 客户四列拓扑，实线=启用、虚线灰=暂不调用/未授权
3. **明确不做**：不新建调度引擎（composite 链路已实现 magpie 式 显式路由 → 账号归属 → 模型名探测 + priority 候选池 + 冷却换号，见 `composite_route_resolver.go:97`）、不改 gateway 核心、不动 GroupsView 的 composite 页签。

### 落地结果（2026-10-28 当日完成）

- **P0**：`KeysView` 端点块改用 `api_base_url || window.location.origin`，`UseKeyModal`、cc-switch 导入同口径。
- **P1 后端**：条目改 `type=apikey`（迁移 `900_xeu_upstream_entry_type_apikey.sql`）；repo 按 `extra.upstream_vendor_id` 过滤；`GET /admin/upstream-entries`、`POST /:id/probe-models`、`PUT /:id/models`；`UpdateEntry` 以现有 credentials 为底合并，避免编辑抹掉 `model_mapping`；复制条目连带复制白名单。单测覆盖探测归一化、白名单写入、超量拒绝、composite 平台放行。
- **P1 前端**：条目标题栏「已登记 N 个」+ 折叠面板（探测 → 勾选 → 保存落库）；未保存的行走 `sync-upstream-preview`；分组下拉对 composite 全量开放并打「综合」徽章；详情表新增模型数（0 时琥珀色）。
- **P2**：新页面 `/admin/aggregation`（供给汇入 / 归一模型表 / 链路视图）+ 路由、侧边栏「聚合调度」、zh/en i18n、组件测试。
- **实测**：本地 mock 上游跑通「建厂家 → 建条目 → 探测 → 落库 → 挂 composite 组 → 客户 key 调 `/v1/models` 与 `/v1/chat/completions`」；历史条目迁移后类型为 apikey，软删行未受影响。
- **未在浏览器内实点**：浏览器自动化被安全策略禁止注入登录态，两个页面由组件测试与构建覆盖，需管理员登录后走一遍。

### 第二轮修正（2026-10-29）

- **分组状态只读**：详情表分组芯片不再可点，只展示「启用 / 本组暂不调用」；本单元启停统一在聚合调度页。
- **探测加固**（对齐 magpie + cc-switch）：
  - 传输错误 / 429 / 5xx 重试 3 次（400ms、1200ms 退避）。实测中转站会间歇 TLS reset（同一 URL 连续请求 000/200 交替），单次探测必然误报。
  - 上游 404/405 时按变体端点重试（`{base}/models`、`{base}/v1/models`，保留路径前缀）。
  - 全部失败回退 **models.dev 目录候选**（先按 base_url 匹配供应商，再按平台映射），结果带 `source` 与含底层原因的 `warning`；目录请求自身也重试。
  - 前端显示来源徽章，并支持**手动添加模型 ID**（上游无 `/models` 时的最后兜底）。
- **聚合调度改为调度单元列表**：单元 = composite 分组。新建单元直接在页面内创建 composite 分组；详情三块为供给汇入（挂载 + 本单元启停）、归一模型表、链路视图；分组管理对 composite 分组新增「聚合调度」入口（`/admin/aggregation?unit=<id>`）。供方编排与需方配置（倍率/RPM/授权/客户 key）彻底分离。

### 第三轮修正（2026-10-29）：多 BaseURL 保存时定稿、模型按钮块

- **厂家多 BaseURL**：`upstream_vendors.base_urls`（jsonb，≤10）。「管理地址」入口在**厂家卡片头部**（紧挨「添加/编辑」），弹窗内多行增删 + 逐地址测速（可选测速出口）；厂家弹窗与条目行都不再配置地址，条目名称下方只读显示当前固化值。
- **保存即定稿**：`POST /admin/upstream-base-urls/latency` 并发测各候选地址——**不带鉴权地 GET {base}/v1/models，计时到响应头**（DNS+TCP+TLS+首包），**拿到任何 HTTP 响应都算可达**，只统计延迟；仅传输错误算不可达。测速后按延迟升序重排并默认选中第一个可达地址；管理员可点击改选首选。**地址池顺序即优先级**：`CreateEntry`/`UpdateEntry` 直接把第一个地址写入 `credentials.base_url`，服务端保存时不再测速。**运行期不切换地址**，避免请求途中换线路导致行为漂移。
- **探测结果按钮块**：模型以按钮块展示，点击切换启用（选中=青柠实底），无复选框；手动添加保留。
- 账号页的模型同步保持原有 404/405 语义，端点变体只在 hub 探测路径启用（`fetchUpstreamModelList(..., allowEndpointVariants)`）。
- **代理不在 hub 配置**：条目级 `proxy_id` 与「测速出口」已移除——代理定义与运行期绑定额外一律走 IP 管理（绑定在账号上）。地址弹窗可**选用** IP 管理已有代理做本次测速（仅该次请求生效、不落库），默认直连。
- **地址弹窗**：顶部「添加地址 / 测速」两个按钮；每行 = 地址输入 + 删除，测速结果直接回填在该地址下方（完整换行展示，不截断）。保存时仍优先选可达地址中最快的一个；全部不可达则按第一个地址保存（本机直连不通不等于上游不可用）。

### 验证口径

`gofmt -l` 干净 → `go build ./...` + `go test`（受影响 package）→ `pnpm check:i18n` + `vue-tsc` + `vite build` → `make build-backend`（embed）→ 重启 8080 浏览器实测（探测真实上游、跨平台综合分组调用）→ 按主题拆提交 push `main` → 补 `docs/WORK_LOG_CN.md`。

