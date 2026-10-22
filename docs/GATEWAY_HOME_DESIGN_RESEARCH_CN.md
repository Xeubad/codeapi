# API 中转站首页设计调研与 CodeAPI 设计语言

> 调研日期：2026-10-22
> 调研范围：AI API 聚合平台 / 中转站 / 网关类产品的首页、模型页与自定义主页
> 目的：为 CodeAPI 默认宣传首页建立一套自己的设计语言，明确“必须展示什么”，而不是照搬某个站点

---

## 0. 一句话结论

中转站首页的“像”，主要不来自配色，而来自四类证据：

1. **模型**：有哪些模型、属于哪个厂商、支持什么能力；
2. **价格**：输入/输出单价、分组倍率、计费方式；
3. **路由与状态**：请求怎么走、当前延迟多少、是否在线；
4. **产品导航**：模型广场、价格、文档、控制台、充值/登录。

视觉只负责把这些证据组织得可信、克制、统一。CodeAPI 已有真实端点和真实延迟，缺的是模型目录、价格和“网关感”的首屏结构。

---

## 1. 调研对象与方法

统一使用桌面 1440×1100 截图，记录首屏、滚动结构、导航、数据展示和视觉 token。部分站点为动态渲染或登录后页面，截图只代表调研时的状态。

| 类型 | 站点 | 用途 |
|---|---|---|
| 国际模型聚合/路由 | [OpenRouter](https://openrouter.ai/) | 数据极简、规模感、模型/供应商墙 |
| 国际 API 网关 | [AiHubMix](https://aihubmix.com/) | 网关叙事、SDK 兼容、模型价格卡 |
| 国内模型市场 | [302.AI](https://302.ai/) | 搜索、分类、模型卡、促销位 |
| NewAPI 自定义主页 | [清羽飞扬模板](https://gist.liushen.fun/LiuShen/new-api-custom-homepage) | 首屏运维卡、模型路由条 |
| 国内中转站（反例） | [DMXAPI](https://dmxapi.cn/) | 识别“AI 味”来源 |
| 国内模型云 | [硅基流动](https://www.siliconflow.cn/) | 大厂 SaaS 风格对照 |
| 模型列表页 | [OpenRouter Models](https://openrouter.ai/models) / [AiHubMix Models](https://aihubmix.com/models) | 模型/价格/指标如何组织 |
| NewAPI 默认前端 | 管理后台截图观察 | 中转站用户的导航预期 |

---

## 2. 案例拆解

### 2.1 OpenRouter：用数据建立规模感

**首屏结构**

- 顶栏：搜索框 + Models / Rankings / Orl / Benchmarks / Chat / Pricing / Docs + 荧光绿 Sign Up。
- Hero：纯黑背景，两行超大标题 “The Unified Interface For Every Model”，副标题 “Better prices, better uptime, no subscriptions.”。
- 主 CTA “Get API Key” 为荧光黄绿实心按钮，次 CTA “Discover Models” 为描边按钮。
- 数据墙：600T+ Monthly Tokens、10M+ Global Users、80+ Providers、500+ Models。
- 四张能力卡：多模态统一入口、跨供应商可用性、价格与性能曲线、数据策略；卡片内直接放模型图标墙、路由图、延迟曲线等“产品碎片”。
- 下方接 Featured Models 列表。

**模型页**

- 左侧多组筛选：输入模态、折扣、上下文长度、价格、厂商、参数、是否零数据保留等。
- 右侧模型行：模型名、描述、厂商、发布日期、输入/输出价、token 用量。
- 没有大面积插画，全部是可扫读的数据。

**视觉**

- 黑/白为主，单一荧光黄绿强调色，无渐变。
- 大字号标题 + 小字号等宽标签；卡片 1px 灰边，圆角很小。
- 信息密度高，但层级清楚：规模数字 → 能力 → 模型。

**为什么像中转站**

首屏不解释“我们是谁”，而是直接给规模、供应商数量、模型数量和价格入口。

**可借鉴**

- 真实规模数字（模型数、分组数、供应商数）。
- 模型/供应商墙。
- 单一强调色 + 黑白底色。
- 价格、上下文、吞吐等硬指标。

**不借鉴**

- 全黑冷感与 CodeAPI 现有浅色玻璃体系冲突。
- 没有真实数据时不要做“数据墙”。

---

### 2.2 AiHubMix：最完整的“网关叙事”

**首屏结构**

- 白底，点阵世界地图作背景。
- 标题 “The Unified Gateway for AI Models”，副标题 “Access every major model through the SDKs you already use.”。
- 三个信任点：710+ models、Official cloud routes、Pay as you go。
- 主 CTA “Get API Key”，次 CTA “Explore models”。
- 三张能力卡：
  1. All major models, one key：模型 logo 汇聚到一个 key；
  2. Drop-in for your SDK：代码 diff，只改 `base_url`；
  3. Auto-failover across clouds：多家云之间的路由图。
- 下方直接接产品 dashboard 预览：余额、今日花费、Claude Code 连接。

**模型页**

- 顶部 Filter：输出模态、分类、厂商 chips（OpenAI、Anthropic、Google、DeepSeek、Qwen 等）、排序。
- 模型卡：厂商 logo、New/折扣 badge、Copy ID、输入/输出价、上下文长度、TTFB、吞吐、简介。
- 价格信息密度很高，但卡片留白足够，不拥挤。

**视觉**

- 白/浅蓝灰画布，蓝色主 CTA；点阵地图是最主要的“技术感”来源。
- 卡片 1px 边、约 16px 圆角；没有大面积渐变。
- 品牌色只出现在 CTA、链接和少量 badge。

**为什么像中转站**

首屏三件事就是中转站的核心：模型接入、SDK 兼容、跨云故障转移；下面直接给余额与账单。

**可借鉴**

- 网关叙事：一个 key、只改 base_url、自动切换。
- 厂商字标 + 模型卡。
- 代码 diff 与路由图。
- 余额/用量预览。
- 模型卡的 TTFB、吞吐、价格。

**不借鉴**

- 过于通用的 SaaS 白；CodeAPI 可以保留自己的玻璃与 ink 面板。

---

### 2.3 302.AI：模型市场的组织方式

**首屏结构**

- 顶部公告条 + 工具导航（管理后台、快捷使用、技术支持、更新日志、公告）。
- 大搜索框；“应用超市 / API 超市 / 客户端 / 模型排行榜 / GitHub / 文章资讯”二级导航。
- 左侧品类列表（语言大模型、图片生成、视频生成、音视频处理、RAG 等）。
- 中间 Hero 轮播；右侧 4×2 快捷入口（数据汇总、应用管理、API Keys、Agent 管理等）。
- 下方“特价专区”，模型卡带厂商 logo、折扣角标、模型名与简介。

**视觉**

- 紫/蓝渐变，入口多、卡片多，信息密度最高。
- 更像“模型市场/应用市场”，不是品牌官网。

**为什么像中转站**

搜索、分类、模型卡、促销位组合起来，用户一眼知道“这里能买什么模型”。

**可借鉴**

- 模型市场结构、分类/筛选、促销位、快捷入口。
- 模型卡上的厂商 logo 与折扣 badge。

**不借鉴**

- 紫蓝渐变与信息过载。
- “应用超市”式的大量入口，会把首页变成导航站。

---

### 2.4 NewAPI 自定义主页（清羽飞扬）：首屏运维卡

**首屏结构**

- 浅暖灰/玻璃画布，右侧一块深色运维卡：
  - ALL SYSTEMS REFINED
  - LATENCY 82ms / UPTIME 99.9% / ROUTES 50+
  - MODEL ROUTING：Claude、Gemini、DeepSeek、Qwen 的占比条
  - `POST /v1/chat/completions` 代码块
- 下方 4 个数字条：50+ 上游、100+ 模型、10+ 协同能力、24/7 稳定访问。
- 再往下是“高级感来自秩序，而不是堆满元素”“一个入口，连接不同模型生态”等章节。

**视觉 token（来自模板源码）**

- 浅色：背景 `#fbfaf7`，主色 `#1d7f8c`，金色 `#c79f62`，卡片 `rgba(255,255,255,.62)`，大圆角 24–34px，柔光阴影。
- 深色：背景 `#090d12`，主色 `#7fe3eb`，金色 `#e3bf7e`。
- 玻璃 + 暖白 + teal/gold，整体偏“高级模板”，不是传统中转站的紫蓝霓虹。

**为什么像中转站**

首屏就是网关运行状态：延迟、路由、模型占比、接口代码，而不是品牌口号。

**可借鉴**

- 首屏运维卡：延迟 / 模型 / 路由。
- 模型路由占比条。
- 数字条。
- 暖白 + 深色双主题。

**不借鉴**

- `UPTIME 99.9%` 如果没有真实监控数据，不要写。
- 34px 大圆角与大量玻璃叠加，容易变成“模板感”；CodeAPI 应控制在 12–20px，玻璃只做容器。

---

### 2.5 DMXAPI：典型的“AI 味”反例

- 紫棕渐变背景、emoji、卡通动物、发光大字、二维码悬浮。
- 信息以口号为主，缺少模型、价格、状态等可验证内容。
- 视觉上同时使用高饱和渐变、emoji、插画、发光，导致“AI 生成感”很强。

**结论**：不采用。渐变、emoji、卡通、无数据支撑的口号，是“AI 味”的主要来源。

---

### 2.6 硅基流动：模型云厂商风格

- 紫光渐变轮播、通用 SaaS 大标题、产品矩阵。
- 更偏“模型云/基础设施厂商”，不是中转站。
- 可借鉴产品矩阵的分层方式；视觉不借鉴。

---

### 2.7 NewAPI 默认前端：导航预期

从管理界面的顶部导航可以观察中转站用户熟悉的入口：

> 主页 / 控制台 / 模型广场 / 排行榜 / 文档 / 关于

这说明用户对“中转站首页”的预期不是纯营销页，而是能进入**模型广场、价格、文档、控制台**的产品入口。

---

## 3. 横向对比

| 站点 | 首屏主角 | 模型展示 | 价格展示 | 状态/延迟 | 导航 | 数据密度 | 材质/视觉 |
|---|---|---|---|---|---|---|---|
| OpenRouter | 规模数字 + 能力卡 | 模型/供应商墙 | 模型页输入/输出价 | 有 uptime/性能叙事 | 产品导航 | 高 | 纯黑 + 荧光绿，无玻璃 |
| AiHubMix | 网关叙事 + 三张卡 | 厂商 logo + 模型卡 | 输入/输出、上下文、TTFB、吞吐 | 余额/账单预览 | 产品导航 | 高 | 白底 + 点阵地图 + 蓝 |
| 302.AI | 搜索 + 品类 + 轮播 | 模型卡/特价区 | 折扣 + 模型卡 | 少量 | 市场导航 | 很高 | 紫蓝渐变、入口多 |
| NewAPI 自定义 | 运维卡 + 数字条 | 模型路由占比 | 无明确价格 | **延迟/路由/uptime** | 少量 | 中高 | 暖白玻璃 + teal/gold |
| DMXAPI | 口号 + 插画 | 无 | 无 | 无 | 简单 | 低 | 紫棕渐变 + emoji |
| 硅基流动 | 品牌轮播 | 产品矩阵 | 价格页 | 无 | 大厂导航 | 中 | 紫光渐变 |

---

## 4. 中转站首页的六个必备模块

### 4.1 首屏网关状态卡

至少包含：

- 在线状态；
- 当前延迟（真实测量，不是装饰数字）；
- 模型数 / 分组数（有真实数据才显示）；
- 3–5 条模型路由（模型名、分组、倍率或价格、状态）；
- 一个真实端点 + 一段 curl / 代码。

作用：让用户在第一屏确认“这是一个能用的网关”，而不是一张品牌海报。

### 4.2 模型目录

至少包含：

- 模型名；
- 厂商（Claude / GPT / Gemini / …）；
- 所属分组；
- 输入/输出价格（每百万 token）或计费方式；
- 分组倍率；
- 状态（可用/维护/即将推出）。

形式可以是卡片、紧凑表格或“卡片 + 表格”混合。首屏只放 4–6 个代表模型，完整目录放模型广场。

### 4.3 接入方式

- base_url；
- API Key；
- 3 步接入（改地址、换 key、直接调用）；
- curl / SDK 片段；
- 复制按钮。

### 4.4 价格与倍率说明

- 按 token / 按次 / 按图片；
- 分组倍率；
- 官方参考价 vs 实收价（如果有）；
- “以控制台结算为准”的说明。

### 4.5 稳定性证据

- 真实客户端往返延迟；
- 最近样本；
- TTFB、连接复用、传输量、协议；
- 如无真实 uptime 监控，**不要写 99.9%**。

### 4.6 产品导航与转化入口

- 模型广场 / 价格 / 接入 / 实测延迟；
- 控制台 / 登录 / 注册；
- Key 用量查询；
- 可选：公告、充值、兑换码、排行榜。

---

## 5. 当前 CodeAPI 首页盘点（2026-10-22）

### 已有

- 真实接入端点（`/v1/messages`、`/v1/chat/completions`）。
- 真实客户端延迟实测（当前/平均/P95/抖动、折线、TTFB、连接复用、传输量、协议）。
- 模型厂商 chips（Claude / GPT / Gemini / Antigravity）。
- curl 接入片段与复制。
- 玻璃、ink + lime、与登录页统一。
- 无外链。

### 缺失

- **模型级目录**：只有厂商，没有具体模型。
- **价格**：没有输入/输出单价、分组倍率、计费方式。
- **网关首屏**：右侧面板在 `model_plaza_enabled=false` 时只剩端点和“价目未开放”，视觉上偏空。
- **产品导航**：目前是“实测延迟 / 快速接入 / 模型”营销锚点，缺少模型广场、价格、控制台。
- **数据密度**：留白大、玻璃面积大，更像品牌官网，不像 API 基础设施。

### 根因

不是“玻璃不对”，而是**信息结构不对**：中转站用户第一眼要找的四件事（模型、价格、路由/状态、接入），当前只完整回答了“接入”和“延迟”。

---

## 6. CodeAPI 设计语言：Gateway Ledger（网关台账）

### 6.1 关键词

**技术台账 / 运维面板 / 克制的玻璃 / 数据优先 / 一个信号色**

不是“AI 科技感”，也不是“苹果展厅”。是把网关的运行数据像台账一样摊开：模型、价格、路由、延迟，一行一行可核对。

### 6.2 设计原则

1. **先给证据，再讲愿景**：首屏先回答模型、价格、稳定性、接入方式。
2. **一屏四件事**：模型、价格、路由/状态、接入。
3. **数据必须真实**：没有数据就明确说“未公开/以控制台为准”，不补假数、不写假 uptime。
4. **玻璃只做容器**：顶栏 + 一张信息卡；主体用实心浅色面板和细线表格。
5. **一个强调色**：lime 只做状态点、关键数字、价格强调；语义色（emerald/amber）只做状态。
6. **数字等宽、标签小写大写**：所有数值用 mono、tabular-nums、右对齐；标签用 mono 大写小字。
7. **密度优先**：缩小圆角与留白，用 1px 细线、分组标题、紧凑行高组织信息。

### 6.3 Token

| 用途 | 浅色 | 深色 |
|---|---|---|
| 画布 | `#f6f7f4` / `#f1f3f6` | `#0b0e14` |
| 主文本 | `#111620` | `#eef2f6` |
| 次文本 | `#5b6472` | `#a7b1bd` |
| 弱文本 | `#8b93a1` | `#737f8f` |
| Ink 面板 | `#111620` / `#0b0e14` | `#111620` |
| Lime 主色 | `#6f9e12` / `#558310` | `#b3d766` |
| 细线 | `rgba(17,22,32,.10)` | `rgba(255,255,255,.10)` |
| 玻璃 | `rgba(255,255,255,.55)` + blur 18px | `rgba(20,27,38,.55)` + blur 18px |
| 圆角 | 12 / 16 / 20px | 同左 |

字体：

- 标题：系统 sans（现有字阶）。
- 标签/数字/代码：mono。
- 桌面 H1：44–54px；移动 32–36px。
- 正文 14–15px；数据标签 9–11px。

布局：

- 版心 `max-w-7xl`（1280px），左右 `px-4 sm:px-6`。
- 区块间距 40–56px；卡片内边距 20–28px。
- 桌面 12 列；Hero 左 45% / 右 55%。

### 6.4 页面结构

```
┌ 顶栏（玻璃）────────────────────────────────────┐
│ Logo   模型与价格  接入  实测延迟   [在线/延迟] [登录] │
├─────────────────────────────────────────────────┤
│ Hero 左：标题 + 说明 + CTA                        │
│ Hero 右：GATEWAY / LIVE 网关面板                  │
│   在线 · 延迟 · 模型数 · 分组数                    │
│   模型路由列表（模型/分组/价格/状态）              │
│   端点 + 迷你折线                                 │
├─────────────────────────────────────────────────┤
│ 事实条：统一入口 / 会话保持 / 实时计费              │
├─────────────────────────────────────────────────┤
│ 模型与价格（卡片或紧凑表格）                       │
├─────────────────────────────────────────────────┤
│ 接入（3 步 + curl + 复制）                        │
├─────────────────────────────────────────────────┤
│ 实测延迟（真实折线 + 明细 + 最近样本）              │
├─────────────────────────────────────────────────┤
│ CTA + 页脚                                       │
└─────────────────────────────────────────────────┘
```

### 6.5 模块规格

#### A. 顶栏

- 左：Logo + 站点名。
- 中：模型与价格、接入、实测延迟（桌面）。
- 右：在线状态 + 实时往返（桌面）、Key 用量查询（≥sm）、主题切换、登录/控制台。
- 移动端：只保留品牌、主题、登录；状态可收进 Hero。
- 模型广场入口仅在 `model_plaza_enabled` 时出现。

#### B. Hero 网关面板

- 标题栏：`GATEWAY / LIVE` + 在线 pill。
- 指标：实时往返（真实）、模型数、分组数。
- 路由列表：4 行；每行厂商字标、模型名、分组、输入价或倍率、状态点。
- 底部：`POST /v1/messages` + 迷你折线 + 样本数。
- 空状态：模型目录未公开时，显示 4 个厂商行 + `目录未公开`，不显示假价格。

#### C. 模型与价格

- 有模型广场数据：
  - 6 张模型卡或 6 行表格；
  - 列：模型、厂商、分组、输入/输出价（/1M）、倍率；
  - 右上“完整模型广场”入口；
  - 底部“价格来自模型广场，实际以控制台结算为准”。
- 无数据：
  - 4 个厂商卡（Claude / GPT / Gemini / Antigravity）+ “已支持”；
  - 一句“登录控制台查看完整模型与价格”；
  - 不显示假价格。

价格计算：

```
生效倍率 = user_rate_multiplier ?? rate_multiplier
实收单价 = 基础单价 × 生效倍率
展示 = formatScaled(实收单价, 1_000_000, 2) + " / 1M"
```

#### D. 接入

- 左：三步说明 + 两个端点 chips。
- 右：深色终端卡 + curl + 复制。
- 保留现有 `home.integration.*` 文案。

#### E. 实测延迟

- 当前往返大字 + 状态 pill。
- 平均、P95、抖动、样本数。
- 折线从 0 起算、最小量程 50ms。
- 右侧：首字节、连接、传输、协议、最近 5 条样本。
- 底部：数据来源说明 + 立即重测。
- 不写 uptime，不写模拟数据。

#### F. CTA / 页脚

- CTA：开始使用 + Key 用量查询。
- 页脚：Logo、版权、版本、Key 用量查询、登录、服务状态。

### 6.6 数据来源与边界

| 数据 | 来源 | 边界 |
|---|---|---|
| 模型/分组/倍率/价格 | `GET /api/v1/model-plaza` | 开关关闭时返回 404，前端显示“目录未公开” |
| 端点 | 当前 `window.location.origin` | 仅内链，无外部跳转 |
| 实时延迟 | 浏览器请求本站静态资源（Range 1KB） | 真实测量，不伪造 |
| 首字节/连接/传输/协议 | Resource Timing | 浏览器不支持时显示 `--` |
| 站点名/Logo | 公开设置 | 沿用现有注入逻辑 |

**当前已知前提**：`model_plaza_enabled=false`。要让首页的价格区显示真实数据，需要在后台开启模型广场；否则只能展示厂商与模型名，不能展示价格。

### 6.7 响应式

| 断点 | 行为 |
|---|---|
| < 640 | 单列；Hero 面板在文案下方；模型卡单列；顶栏隐藏中间导航与 Key 用量 |
| 640–1023 | 模型卡两列；事实条三列；顶栏保留主题/登录 |
| ≥ 1024 | Hero 两列；模型卡三列；顶栏完整导航 |
| ≥ 1280 | 版心 1280，最大留白固定 |

### 6.8 文案规则

- 不写“最强/最便宜/无限并发”等无证据的宣传。
- 不写没有监控数据支撑的 uptime。
- 价格必须带单位（/1M tokens）与“以控制台为准”。
- 模型目录未公开时直接说明，不用“敬请期待”掩盖。
- 所有“查看文档/模型广场”链接必须站内；不引入外部跳转。

### 6.9 验收清单

- [ ] 1440 / 1024 / 390 三档无横向溢出。
- [ ] 浅色 / 深色均通过对比度检查。
- [ ] 首屏可见：在线状态、实时延迟、模型数/分组数（有数据时）、至少 4 条路由、端点。
- [ ] 模型区有真实模型与价格，或明确的空状态，不出现假数据。
- [ ] 延迟面板为真实采样，折线随时间增长。
- [ ] 顶栏导航均指向站内。
- [ ] `typecheck`、`check:i18n`、`test:run`、`build` 全部通过。
- [ ] 登录页与首页使用同一套背景、玻璃、字阶与 ink/lime token。

---

## 7. 落地顺序

1. **确认数据前提**：是否开启模型广场；若不开启，接受“模型名可展示、价格显示未公开”。
2. **Hero 网关面板 + 顶栏产品导航**：先把“网关感”立起来。
3. **模型与价格区**：有数据走卡片，无数据走厂商卡 + 控制台 CTA。
4. **接入 / 延迟 / CTA 的密度调整**：缩小留白与圆角，统一细线与 mono 标签。
5. **登录页对齐**：只核对 token 与字阶，不改认证流程。
6. **视觉验收**：三档宽度、双主题、真实数据、无外链。

---

## 8. 待定问题

1. 是否开启 `model_plaza_enabled`，让首页直接展示真实模型与价格？
2. 厂商标识使用文字字标（CL/GPT/GM/AG）还是官方 logo？建议先用文字字标，避免版权与风格不统一。
3. 是否需要公告/促销条？如果需要，必须有真实运营内容，不做装饰性公告。
4. 是否需要“模型排行榜/用量榜”？需要真实聚合数据；没有就只做“模型目录”。
5. 玻璃强度：建议只保留顶栏 + 一块信息卡，其余实心；是否接受降低玻璃比例？

---

## 9. 参考链接

- [OpenRouter 首页](https://openrouter.ai/)
- [OpenRouter Models](https://openrouter.ai/models)
- [AiHubMix 首页](https://aihubmix.com/)
- [AiHubMix Models](https://aihubmix.com/models)
- [302.AI](https://302.ai/)
- [DMXAPI](https://dmxapi.cn/)
- [硅基流动](https://www.siliconflow.cn/)
- [NewAPI 自定义主页模板（清羽飞扬）](https://gist.liushen.fun/LiuShen/new-api-custom-homepage)
- [AxonHub / NewAPI 文章](https://blog.liushen.fun/posts/b404081c/)
- [New API 仓库](https://github.com/QuantumNous/new-api)

---

# 第二部分：本项目路由盘点与 Home 设计方案

> 分析日期：2026-10-22
> 依据：`frontend/src/router/index.ts`（988 行）、`/api/v1/settings/public` 当前返回值、`HomeView.vue` / `ModelPlazaView.vue` 现状
> 目标：首页只链接真实存在且当前可用的页面；所有重点信息都有真实数据来源；不因为“参考站有”就加假入口

## 10. 本项目路由盘点

### 10.1 公开路由

| 路径 | name | 鉴权 | 额外守卫 / 说明 | 首页是否可链 |
|---|---|---|---|---|
| `/home` | Home | 否 | 默认宣传页；`backend_mode_enabled=true` 且未登录时会被后端模式守卫重定向到 `/login` | 首页自身 |
| `/` | — | — | 重定向 `/home` | 是（Logo） |
| `/login` | Login | 否 | 已登录访问会跳到用户/管理员控制台 | 是（登录） |
| `/register` | Register | 否 | 页面内部受 `registration_enabled` 控制；后端模式下仅允许有 pending 会话时访问 | 仅注册开启时 |
| `/email-verify` | EmailVerify | 否 | 邮箱验证流程 | 否 |
| `/auth/callback` | OAuthCallback | 否 | 通用 OAuth 回调 | 否 |
| `/auth/linuxdo/callback` | LinuxDoOAuthCallback | 否 | LinuxDo 回调 | 否 |
| `/auth/wechat/callback` | WeChatOAuthCallback | 否 | 微信回调 | 否 |
| `/auth/wechat/payment/callback` | WeChatPaymentOAuthCallback | 否 | 微信支付 OAuth 回调 | 否 |
| `/auth/dingtalk/callback` | DingTalkOAuthCallback | 否 | 钉钉回调 | 否 |
| `/auth/dingtalk/email-completion` | dingtalk-email-completion | 否 | 钉钉邮箱补全 | 否 |
| `/auth/oidc/callback` | OIDCOAuthCallback | 否 | OIDC 回调 | 否 |
| `/forgot-password` | ForgotPassword | 否 | 忘记密码 | 否（登录页内链） |
| `/reset-password` | ResetPassword | 否 | 重置密码 | 否 |
| `/key-usage` | KeyUsage | 否 | 公开 Key 用量查询；后端模式白名单内 | **是（次 CTA）** |
| `/legal/:documentId` | LegalDocument | 否 | 法务文档；后端模式白名单内 | 页脚可选（需配置文档） |
| `/model-plaza` | ModelPlaza | 否 | **受 `model_plaza_enabled` + `model_plaza_require_auth` 双重控制**；关闭时重定向 `/home`；强制登录时跳 `/login` | **仅开启时** |
| `/payment/result` | PaymentResult | 否 | 支付结果回调页；后端模式白名单内 | 否 |
| `/payment/stripe` | StripePayment | 否 | Stripe 支付页 | 否 |
| `/payment/airwallex` | AirwallexPayment | 否 | Airwallex 支付页 | 否 |
| `/payment/stripe-popup` | StripePopup | 否 | Stripe 弹窗页 | 否 |

### 10.2 用户路由（需要登录）

| 路径 | name | 说明 | 首页是否可链 |
|---|---|---|---|
| `/dashboard` | Dashboard | 用户控制台；管理员在首页 CTA 中应走 `/admin/dashboard` | **是（已登录主 CTA）** |
| `/keys` | Keys | API Key 管理与创建 | 经控制台进入 |
| `/batch-image`（别名 `/docs/batch-image`） | BatchImageGuide | 批量图片指南，登录后可用 | 否 |
| `/usage` | Usage | 用户调用记录 | 经控制台进入 |
| `/redeem` | Redeem | 兑换码 | 否 |
| `/affiliate` | Affiliate | 邀请返利；受 `affiliate_enabled` 控制 | 否 |
| `/available-channels` | UserAvailableChannels | 用户可用渠道；受 `available_channels_enabled` 控制 | 否 |
| `/profile` | Profile | 个人资料 | 否 |
| `/subscriptions` | Subscriptions | 我的订阅；受 `requiresSubscription` 控制 | 否 |
| `/purchase` | PurchaseSubscription | 购买订阅；受 `requiresPayment` + `payment_enabled` 控制 | 否（当前支付关闭） |
| `/orders` | OrderList | 我的订单；受 `requiresPayment` 控制 | 否（当前支付关闭） |
| `/payment/qrcode` | PaymentQRCode | 扫码支付；受 `requiresPayment` 控制 | 否 |
| `/custom/:id` | CustomPage | 登录后自定义页面 | 否 |
| `/monitor` | ChannelStatus | 用户渠道监控；需要登录，非公开状态页 | 否（目前不是公开 uptime 页） |

### 10.3 管理端路由（需要登录 + 管理员）

| 路径 | name | 说明 |
|---|---|---|
| `/admin` | — | 重定向 `/admin/dashboard` |
| `/admin/dashboard` | AdminDashboard | 管理员总览 |
| `/admin/ops` | AdminOps | 运维 |
| `/admin/audit-logs` | AdminAuditLogs | 审计日志 |
| `/admin/users` | AdminUsers | 用户管理 |
| `/admin/groups` | AdminGroups | **分组/模型/倍率管理（首页模型价格的数据源之一）** |
| `/admin/channels` | — | 重定向 `/admin/channels/pricing` |
| `/admin/channels/pricing` | AdminChannels | **渠道与模型定价（首页模型价格的数据源之一）** |
| `/admin/channels/monitor` | AdminChannelMonitor | 渠道监控 |
| `/admin/subscriptions` | AdminSubscriptions | 订阅管理 |
| `/admin/accounts` | AdminAccounts | 账号池 |
| `/admin/plugins` | AdminPlugins | 插件 |
| `/admin/announcements` | AdminAnnouncements | 公告管理（**没有公开公告路由**） |
| `/admin/proxies` | AdminProxies | 代理 |
| `/admin/redeem` | AdminRedeem | 兑换码管理 |
| `/admin/promo-codes` | AdminPromoCodes | 优惠码 |
| `/admin/settings` | AdminSettings | **模型广场开关、强制登录、首页内容等设置** |
| `/admin/risk-control` | AdminRiskControl | 风控 |
| `/admin/prompt-audit` | AdminPromptAudit | 提示词审计 |
| `/admin/usage` | AdminUsage | 全站用量 |
| `/admin/affiliates` | — | 重定向 `/admin/affiliates/invites` |
| `/admin/affiliates/invites` | AdminAffiliateInvites | 邀请 |
| `/admin/affiliates/rebates` | AdminAffiliateRebates | 返利 |
| `/admin/affiliates/transfers` | AdminAffiliateTransfers | 转账 |
| `/admin/orders/dashboard` | AdminPaymentDashboard | 支付看板 |
| `/admin/orders` | AdminOrders | 订单 |
| `/admin/orders/plans` | AdminPaymentPlans | 订阅套餐 |
| `/:pathMatch(.*)*` | NotFound | 404 |

### 10.4 会直接影响首页设计的关键发现

1. **首页是公开的，但后端模式下未登录用户会被强制到 `/login`**：`backend_mode_enabled=false` 时才展示宣传页。首页设计默认按“正常模式”处理。
2. **`registration_enabled=false`（当前）**：现在首页状态条写“开放注册”、CTA 写“注册并开始”都与实际设置不符。首页必须按开关切换文案。
3. **`model_plaza_enabled=false`（当前）**：`/model-plaza` 会被守卫重定向 `/home`，首页不能无条件链过去，也不能展示真实模型/价格。
4. **`model_plaza_require_auth=false`（当前）**：一旦开启模型广场，匿名用户可访问；若改为 `true`，首页的模型 CTA 必须指向 `/login?redirect=/model-plaza`。
5. **`payment_enabled=false`（当前）**：首页不要出现充值/购买/套餐入口。
6. **`doc_url=''`，且没有公开文档路由**：参考站的 “Docs” 导航不能照搬；`/docs/batch-image` 是登录后页面。首页的“接入”区就是当前唯一的公开接入说明。
7. **没有公开排行榜、公告、状态页**：`/monitor` 需要登录，公告只有管理端。首页不要加这三类导航或模块。
8. **`/key-usage` 是公开的、后端模式白名单内**：这是首页最可靠的次 CTA。
9. **模型价格数据来自管理端的 `/admin/groups` + `/admin/channels/pricing`，公开出口是 `/api/v1/model-plaza`**：首页应复用同一个公开接口，而不是另起一套假数据。
10. **首页 CTA 必须按登录态分流**：未登录 → `/login` 或 `/register`；已登录 → `/dashboard`；管理员 → `/admin/dashboard`。

## 11. 与调研对象的导航/信息结构比对

| 参考站导航/模块 | OpenRouter | AiHubMix | 302.AI | NewAPI | 本项目对应路由 | 当前可用性 | 首页处理 |
|---|---|---|---|---|---|---|---|
| Models / 模型 | Models | Models | API 超市 | 模型广场 | `/model-plaza` | 受开关控制（当前关闭） | 开启后进入导航；关闭时用首页 `#models` 锚点 + 厂商兜底 |
| Pricing / 价格 | Pricing | 模型卡价格 | 特价专区 | 价格页 | `/model-plaza` 内的价目 | 同上 | 与模型合并为“模型与价格”区 |
| Docs / 文档 | Docs | Docs | 文章资讯 | 文档 | 无公开路由 | ✗ | 用首页“接入”区替代；不做假文档入口 |
| Rankings / 排行榜 | Rankings | — | 模型排行榜 | 排行榜 | 无公开路由 | ✗ | 不展示 |
| Chat / Playground | Chat | Playground | — | — | 无 | ✗ | 不展示 |
| Apps / 应用 | — | Apps | 应用超市 | — | 无 | ✗ | 不展示 |
| Console / 控制台 | — | Dashboard | 管理后台 | 控制台 | `/dashboard`、`/admin/dashboard` | 需登录 | 已登录主 CTA |
| Login / Sign up | Sign Up | Sign in | 登录/注册 | 登录 | `/login`、`/register` | 注册受开关控制 | 按登录态 + 注册开关分流 |
| Key 用量 | — | — | — | — | `/key-usage` | 公开 | **我们的差异项，放次 CTA** |
| Status / Uptime | 有可用性叙事 | — | — | 有 uptime 数字 | `/monitor`（需登录） | 非公开 | 用真实客户端延迟替代，不写 uptime |
| Announcements / 公告 | — | — | 公告 | — | 仅 `/admin/announcements` | 无公开出口 | 不展示 |
| Recharge / 余额 | — | 余额预览 | 钱包 | 充值 | `/purchase`、`/orders` | 支付关闭 + 需登录 | 不展示 |

**结论**：参考站看起来“模块很多”，是因为它们有对应的公开路由和数据；本项目当前只有 **首页、登录/注册、模型广场（可选）、Key 用量、控制台** 这几个公开/半公开入口。首页设计应该围绕这五个入口组织，而不是补齐参考站的导航数量。

## 12. Home 页面设计方案 v2（Gateway Ledger · 收敛版）

> v2 依据（2026-10-22 确认）：
> 1. `model_plaza_enabled` 不开启；
> 2. `registration_enabled` 不开启；
> 3. **模型信息只能登录后在模型广场访问**——公开首页不出现任何模型名、厂商列表、价格、倍率。

### 12.1 公开信息边界（最高优先级）

| 信息 | 公开首页 | 登录后 |
|---|---|---|
| 具体模型名 | ✗ | 模型广场 |
| 厂商列表 / 模型 logo | ✗ | 模型广场 |
| 价格 / 倍率 / 计费档位 | ✗ | 模型广场 |
| 协议兼容（OpenAI / Anthropic 格式） | ✓ | ✓ |
| 接入端点（/v1/messages 等） | ✓ | ✓ |
| 实时延迟（真实探针） | ✓ | ✓ |
| Key 用量查询 | ✓ | ✓ |
| 控制台 / 充值 / 订阅 | 登录入口 | 控制台 |

落地含义：

- 首页删除“模型覆盖”区，不出现厂商卡与模型 chips；
- Hero 文案不点名 Claude / GPT / Gemini 等具体模型/厂商，改用“多家主流模型平台”；
- Gateway 面板展示协议、端点、真实延迟，不展示模型路由；
- 想看模型与价格的唯一路径：登录 → 模型广场。

> 运营提示：要让“登录后在模型广场查看”成立，后台需要开启模型广场并勾选“需要登录”
> （`model_plaza_enabled=true` + `model_plaza_require_auth=true`）。首页不依赖该开关，只负责把访客送到登录页。

### 12.2 页面结构

```
┌ 顶栏（玻璃）────────────────────────────────────────┐
│ Logo   接入   实测延迟            Key 用量  主题  登录 │
├─────────────────────────────────────────────────────┤
│ Hero 左                        Hero 右 GATEWAY/LIVE │
│ 服务在线 · 往返 X ms            延迟 / 协议 / 计费     │
│ H1 + 说明 + 登录 CTA            协议与端点            │
│ Key 用量次 CTA                  迷你折线              │
│                                 登录查看模型与价格 →   │
├─────────────────────────────────────────────────────┤
│ 事实条：统一入口 / 会话保持 / 实时计费                  │
├─────────────────────────────────────────────────────┤
│ 接入：3 步 + 端点 + curl + 复制                       │
├─────────────────────────────────────────────────────┤
│ 实测延迟：真实折线 + TTFB/连接/传输/协议 + 最近样本     │
├─────────────────────────────────────────────────────┤
│ CTA（登录/进入控制台 + Key 用量） + 页脚               │
└─────────────────────────────────────────────────────┘
```

公开首页一共 6 个区块：顶栏、Hero、事实条、接入、实测延迟、CTA/页脚。**没有模型区。**

### 12.3 顶栏

- 左：Logo + 站点名。
- 中：接入（`#integration`）、实测延迟（`#latency`）。
- 右：Key 用量（`/key-usage`）、主题切换、登录/控制台。
- 移动端：品牌 + 主题 + 登录；导航链接进“菜单”。
- 不出现：模型、价格、文档、排行榜、公告、充值。

### 12.4 Hero 左侧

- 状态条：`服务在线`（探针失败时 `连接异常`）+ `往返 {ms} ms`（真实）。
- H1：`一个密钥，` / `接管全部模型调用`。
- 说明（不点名厂商）：`多家主流模型平台通过同一套接口接入；按 token 计费，登录后可在模型广场查看可用模型与价格。`
- 主 CTA：
  - 未登录：`登录` → `/login`；
  - 已登录：`进入控制台` → `/dashboard`（管理员 `/admin/dashboard`）。
- 次 CTA：`Key 用量查询` → `/key-usage`。
- 辅助说明：`登录后创建 API Key 即可调用，模型与价格在模型广场查看。`

### 12.5 Hero 右侧：Gateway 面板

- 标题栏：`GATEWAY / LIVE` + 在线状态（真实探针）。
- 指标（3 格）：
  1. 实时往返：`{ms}`（真实，lime 强调）；
  2. 接入协议：`OpenAI / Anthropic`；
  3. 计费方式：`按 token`。
- 协议与端点块：
  - `POST /v1/messages`
  - `POST /v1/chat/completions`
  - `OpenAI 兼容 · Anthropic 原生`
- 迷你折线 + 样本数（真实）。
- 底部 CTA：`登录查看模型与价格` → `/login`（已登录 → 控制台）。
- **不出现**：模型路由行、厂商字标、模型数、价格。

### 12.6 事实条

沿用三个事实，去掉厂商名：

- `统一入口` — `多平台 · 单密钥`；
- `会话保持` — `同一会话固定上游`；
- `实时计费` — `按 token 逐次结算`。

### 12.7 接入区

- 3 步：改 base_url、换 API Key、直接调用；
- 端点 chips：`POST /v1/messages`、`POST /v1/chat/completions`；
- 右侧深色终端卡 + curl + 复制；
- 说明：`兼容 Anthropic / OpenAI SDK 的调用格式`。
- 不引导到任何模型清单；模型与价格统一指向“登录后模型广场”。

### 12.8 实测延迟区

与 v1 相同：当前往返、平均、P95、抖动、样本数、真实折线（从 0 起算、最小量程 50ms）、
TTFB / 连接 / 传输 / 协议、最近 5 条样本、数据来源说明、立即重测。

### 12.9 CTA 与页脚

- CTA：主按钮同 Hero（登录/进入控制台）；次按钮 `查询 Key 用量`。
- 页脚：版权、版本、`Key 用量查询`、`登录`（或 `控制台`）、`服务运行中`。
- 无注册、无充值、无外链。

### 12.10 状态矩阵（收敛后）

| # | 登录态 | 条件 | 行为 |
|---|---|---|---|
| 1 | 游客 | 正常模式 | 全部模块可见；主 CTA 登录；无任何模型信息 |
| 2 | 已登录（用户） | 正常模式 | 主 CTA 进入控制台；其余相同 |
| 3 | 已登录（管理员） | 正常模式 | 控制台指向 `/admin/dashboard` |
| 4 | 任意 | `backend_mode_enabled` | 未登录访客被守卫送 `/login`，首页不承担宣传职责 |
| 5 | 任意 | `home_content` / `compact_home_enabled` | 优先渲染自定义/紧凑首页，v2 不渲染 |

注册开关与模型广场开关均关闭，不再出现在矩阵中。

### 12.11 视觉规范

沿用 Gateway Ledger：

- 画布：暖中性浅色 / 深色 ink；背景 `brand-backdrop` + 柔光团 + 颗粒。
- 玻璃：顶栏 + Hero 网关卡（ink 深色）+ 事实条；接入/延迟区实心浅色 + 1px 细线。
- 强调色：lime 只用于状态点、实时数字、登录 CTA 上的关键数字；emerald/amber 只做状态。
- 字体：标题 sans，数据/标签 mono，数值 `tabular-nums`。
- 圆角：12 / 16 / 20px；区块间距 40–56px；卡片内边距 20–28px。
- 与登录页共用：`brand-backdrop`、`brand-glass`、`brand-field`、按钮形态、字阶。
- 禁止：紫蓝渐变、emoji、卡通、发光大字、假数据、整页玻璃、模型/厂商/价格信息。

### 12.12 文案规则

- 不出现任何模型名、厂商名（Claude / GPT / Gemini / Antigravity 等一律不写）。
- 不出现价格、倍率、“按模型计费”等表述；计费只说“按 token”。
- 不写“开放注册”“注册并开始”。
- 不写没有监控数据支撑的 uptime。
- “查看模型与价格”必须与登录动作绑定（未登录 → `/login`；已登录 → 控制台）。
- 现有 `home.hero.description`、`home.facts.models.value` 中的厂商枚举需要同步改写。

### 12.13 验收清单

- [ ] 公开首页全文检索不到 `Claude`、`GPT`、`Gemini`、`Antigravity`、`$`、`/1M`、`倍率`。
- [ ] 首屏可见：在线状态、真实往返、协议、两个端点、登录入口。
- [ ] 主 CTA：游客 → `/login`；已登录 → 对应控制台；不出现注册。
- [ ] 不请求 `/api/v1/model-plaza`，不渲染模型广场链接。
- [ ] 事实条与 Hero 说明不含厂商枚举。
- [ ] 1440 / 1024 / 390 无横向溢出；浅色/深色统一。
- [ ] `typecheck`、`check:i18n`、`test:run`、`build` 通过。
- [ ] 登录页与首页共用同一套 token；登录页自身的品牌文案如含厂商枚举，一并按本规则收敛。

### 12.14 实施顺序

1. **删模型区**：移除首页“模型覆盖/模型与价格”区与 Gateway 面板中的模型路由行。
2. **改文案**：Hero 说明、事实条去掉厂商枚举；状态条去掉“开放注册”；CTA 全部改为登录语义。
3. **Gateway 面板重排**：指标改为 延迟/协议/计费，新增协议与端点块 + 登录查看 CTA。
4. **顶栏**：去掉模型入口，导航收敛为 接入/实测延迟/Key 用量/登录。
5. **密度与视觉收口**：玻璃收缩、细线、mono、圆角、移动端菜单。
6. **回归**：§12.13 验收清单 + 三档宽度 + 双主题。
7. **登录页一致性**：检查 AuthLayout 左栏文案是否含厂商枚举，同步收敛。
8. **运营侧（非代码）**：如需“登录后在模型广场查看”，在后台开启模型广场并勾选强制登录。

### 12.15 实现状态（2026-10-22）

- 已按 v2 落地 `frontend/src/views/HomeView.vue`、zh/en `landing.ts` 与 `HomeView.compact.spec.ts`。
- 8080 已重新嵌入并重启（PID 71457，`assets/index-DikISsyc.js`），`/health` 正常。
- 实测公开首页不含 `Claude/GPT/Gemini/Antigravity`、不含 `/model-plaza` 链接、不请求模型广场接口；延迟采样正常。
- `model_plaza_enabled` 与 `registration_enabled` 仍保持关闭；后台开启模型广场（并勾选需要登录）后，登录用户可在模型广场查看模型与价格，首页代码无需再改。
