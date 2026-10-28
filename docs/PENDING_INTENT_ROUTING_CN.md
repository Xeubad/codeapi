# 待定：意图路由（Intent Routing）详细设计 — 暂缓项

> 状态：**待定文档，暂不实现**。等上游聚合（Upstream Hub，见 [UPSTREAM_HUB_DESIGN_CN.md](./UPSTREAM_HUB_DESIGN_CN.md)）真实上游连通性跑通后再启动本项设计。
> 决策来源：2026-10-24，管理员确认意图路由从 MVP 拆出，放置在网关**后置**，由管理员人工配置；不绑定 Jev，分类器可插拔。
> 本文由仓库外 `sub2api-demo/PENDING_INTENT_ROUTING.md` 提炼迁入。

## 一、已确定的约束

| 约束 | 内容 |
|---|---|
| 位置 | 网关**后置**：不参与主路由决策链，作为转发链路的后处理/旁路阶段存在。主路径（归一模型 → priority_fallback 降级）永远独立可用 |
| 配置方式 | 管理员人工配置，不做自动学习/自动建议 |
| 分类器 | 可插拔：任意 OpenAI 兼容端点或 TypeSafe（Jev）只是其一；配置形态对齐 `ContentModerationConfig`（engine + base_url/model/api_key/timeout，默认 3000ms、上限 30000ms，见 `backend/internal/service/content_moderation.go`） |
| 失败语义 | 分类器超时/出错/未知标签 → 一切照旧，绝不阻塞、绝不拒绝请求 |
| 限额 | 分类调用本身计入所在聚合的限额，防刷量 |

## 二、后置架构的三种可能形态（待论证）

v0 方案里意图分类在**路由前**（pre-routing）执行：分类 → 改写模型/调候选顺序 → 再进入归一路由。管理员确认的后置（post-gateway）意味着分类挪到主路由决策**之后**，可能的形态：

- **形态 A（改写后重路由）**：主路由先按归一表选好目标，后置分类器旁路判定意图，命中规则且目标不同时改写模型并重新走一次归一解析。代价是多一跳重路由；收益是主路径零改动、分类完全旁路。
- **形态 B（响应侧干预）**：分类只用于记录、审计、偏好统计，不影响当次请求路由，为未来的自动建议积累数据。最保守。
- **形态 C（异步预判）**：会话首请求分类后缓存意图标签（带 TTL），后续请求直接查缓存决定路由偏置，不在请求路径上等分类器。延迟最优，但多轮对话意图漂移需要处理。

选型要回答的问题：后置分类的触发点在 handler 链的哪一层？改写模型后计费归属怎么算？与现有 `CompositeRouteResolver` 的解析时序如何交错？

## 三、启动前需要回答的问题清单

1. 后置分类的具体挂载点（gateway handler / 调度器后 / 转发前最后一刻），及其对现有计费、审计链路的影响面。
2. 意图标签集合从哪里来：管理员手工定义枚举（v0 方案）已定，但标签之间的优先级、冲突规则需细化。
3. 分类文本抽取策略：首条用户消息？拼接最近 N 轮？截断长度上限？多轮会话意图漂移如何处理。
4. 分类器配置的存储位置：setting 机制（`GET/PUT /admin/settings/...`）还是独立表（迁移 `900_xeu_*`）——取决于是否需要按聚合分组分别配置。
5. 分类测试面板（配置自证）的交互与安全边界（只读预览、不产生真实转发）。
6. 观测指标：分类命中率、改写率、降级触发率、分类延迟分布——需要在 MVP 的用量统计里预留埋点点位。

## 四、v0 方案中可直接复用的设计资产

```go
// internal/service/relayroute/（届时新包）
type IntentClassifierConfig struct {
    Engine    string  // "openai_compatible" | "typesafe"
    BaseURL   string
    Model     string
    APIKey    string  // 加密存储
    TimeoutMS int     // 默认 3000，上限 30000（对齐内容审核）
}

type IntentRule struct {
    Label   string  // 意图标签（枚举，自动注入分类提示词）
    Action  string  // "rewrite_model" | "prefer_platform"（后置架构下动作语义需重审）
    Target  string
    Enabled bool
}
```

- 提示词自动生成（枚举制：只允许返回已定义标签或 unknown）。
- 两个引擎适配器先例：openai_compatible 走 Chat Completions；typesafe 走问答式分类（`content_moderation_typesafe.go`）。
- 分类测试面板 UI：原型 `sub2api-demo/upstream-hub-proto.html` 的 `dlg-intent` 弹窗与测试区已画出（该原型文件已随 demo 目录删除，需要时从 git 历史或本描述重建）。
