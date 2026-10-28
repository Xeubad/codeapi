# 上游聚合（Upstream Hub）设计文档

> 状态：MVP 已落地（2026-10-26，`bdb49a5af`）并经多轮交互迭代收口（至 `e03e5a205`，区间内 15 个 upstream-hub 提交）。
> 本文由仓库外 `sub2api-demo/UPSTREAM_HUB_PLAN.md`（规划稿 v6）提炼而来，按当前代码实况修订；规划期各版本（v1–v6）的演进痕迹已去除，只保留仍然有效的决策与尚未实施的设计。
> 意图路由为独立待定项，见 [PENDING_INTENT_ROUTING_CN.md](./PENDING_INTENT_ROUTING_CN.md)。

## 一、定位

上游聚合解决的问题：智能体对接的上游分散在多家平台（含国产 Kimi / 智谱 GLM / DeepSeek / MiniMax），没有统一分组、配置分散、模型名不统一。

**定位**：与「账号管理」同源的功能——账号管理维护的是厂商凭据（订阅账号、OAuth），上游聚合维护的是**上游厂商**（透传 Base URL + Key 的供货渠道）；两者原理相同（凭据存储、健康检查、调度），只是维护对象与引入的聚合技术不同。

**核心结论：网关不重构，功能主体是编排层。** sub2api 已具备全部中转原语（`type=upstream` 账号、composite 分组、模型路由、协议翻译矩阵、failover/冷却），上游聚合 = 现有组件的受控组合 + 一张厂家台账新表。

### 心智模型：一个数据、两个入口、三个视角

- **唯一数据**：Account 行（accounts 表）。「账号管理」与「上游聚合」不是两套数据，是同一张表的两个管理入口。
- **归属判定一条规则**：有 `extra.upstream_vendor_id` → 由上游聚合管理；没有 → 由账号管理管理。分组挂载（account_groups）与来源入口无关。
- **三个视角**：供货（聚合：厂家预算/协议/条目）→ 分配（分组：挂载谁、本组是否启用）→ 消费（客户：面板自建 key）。

## 二、复用的现有原语（已核实）

| 能力 | 现有原语 | 位置 |
|---|---|---|
| 上游透传账号 | `AccountTypeUpstream = "upstream"`（Base URL + API Key 直连上游），完整校验链保留 | `backend/internal/domain/constants.go` |
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
- `notes`、软删（TimeMixin + SoftDeleteMixin）
- 删除保护：厂家下仍有条目时拒绝删除（`ErrUpstreamVendorNotEmpty`）

### 条目：`Account(type=upstream)` + `extra.upstream_vendor_id`

一条上游记录 = 一行 Account，不建新实体。归属判定读 `extra.upstream_vendor_id`（兼容 int64/int/float64/string 反序列化）。条目创建复用 `AdminService.CreateAccount` 全套校验，`SkipDefaultGroupBind` 避免自动挂默认组。

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
| Repo | `backend/internal/repository/upstream_vendor_repo.go` | 厂家持久化；条目查询按 `type=upstream` + extra 归属过滤；分组绑定 tx 内 delete + CreateBulk |
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
- **详情表分组色条**：可点击按钮，切换该条目在该分组的「暂不调用」（`account_groups.enabled`），启用=青柠底色、停用=灰。
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

**唯一未验证缺口**：翻译路径对 `type=upstream` 账号是否生效（与 apikey 账号大概率同路，但无测试覆盖）——连通性实测的核心问题，也是"能否复刻 magpie/cc-switch"的最终判据。**待真实上游凭据实测。**

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

**已上线可用**：厂家台账（预算/超限硬阻断）、条目增删改复制、并发+优先级、分组绑定与每分组「暂不调用」、全局可调度开关、调度与统计尊重 enabled、`/admin/upstream-hub` 完整界面（zh/en）。

**已知缺口**：

1. **真实上游连通性测试未跑**——等待真实 Base URL + Key，验证翻译矩阵对 `type=upstream` 的覆盖（§七）。
2. GroupsView 建组「从上游聚合引用」区块未建（机制已从聚合侧打通）。
3. 条目平台创建后不可改（弹窗中既有条目平台锁定）。
4. 复制条目到其它厂家/分组时 enabled 重置为 true。
5. `/v1/models` 未过滤"本分组暂不调用"的条目。
6. 账号管理页对聚合条目无只读标记（双头管理护栏未落）。
7. 运行期厂家预算槽位与厂家级联动冷却未实现（§七接缝 3）。
8. 探活三级未实现（§八）。
9. 意图路由整体待定（见 [PENDING_INTENT_ROUTING_CN.md](./PENDING_INTENT_ROUTING_CN.md)）。

**后续实施顺序**（原规划保留）：连通性实测 → GroupsView 双来源 UI → 探活三级 + 预算槽位接缝 → 意图路由（待定）。

## 十、风险

1. `upstream` 类型账号跨协议覆盖面需实测（openai / anthropic / gemini 三协议各连一家真实上游）。
2. `Group.model_routing` 与 `CompositeModelRoute.priority` 的一致性由界面层写入顺序保证（先 routes 后 model_routing），需在测试中覆盖原子性。
3. 意图路由相关的延迟/计费/时序风险整体移至待定项。
