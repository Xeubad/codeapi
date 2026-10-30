# 前端全站 UI 重构调研：React Bits 与通用组件库选型对比

> 调研日期：2026-10-28
> 调研范围：React Bits（reactbits.dev）、Ant Design / Ant Design Vue、shadcn/ui / shadcn-vue、Element Plus、Naive UI、vue-bits
> 调研方法：官网实际抓取（首页/介绍页/组件目录）、GitHub API 实时数据、本项目前端代码实证（令牌层、common 组件层、视图规模）
> 目的：为「全站页面 UI 重构 + 固定一套自己的通用 UI 库」提供选型依据、差距评估与周期估算

---

## 0. 一句话结论

本项目前端是 **Vue 3**，React 生态的 React Bits / Ant Design / shadcn/ui 均不可直接使用；重构应采用「**成熟 Vue 组件库做地基 + 主题映射对齐现有设计令牌 + `App*` 封装层收口页面依赖 + vue-bits 限高视觉页点缀**」的四层架构。**推荐地基为 Element Plus**（agent 批量重构场景下语料覆盖度权重最高），当前 UI 风格评价约 7.5/10、是资产而非负债，EP 默认风格与之差距约 7-8/10 但主题映射可收回 85-90%，全站重构周期估算 **5-7 周（全职投入口径）**。

> **2026-10-31 修订（最终口径见 §7.4）**：本节推荐的「全站 5-7 周 sweep」当天经历了两次修订——
> 第一次（§7.2）以冲突面实测否决大爆炸迁移；第二次（§7.4）因 fork 上游策略改为「完全解耦 + 按需借鉴」，
> 冲突面代价消失，全站迁移从「不值得」变成「取决于 5-7 周的机会成本」。
> 推荐顺序：**A 机制层（不引库，无悔）→ B 挑 2-3 个自研密集页实测 EP 主题映射 → C 达标才排期全量**。
> vue-bits 暂不引入。调研的实证部分（§1-§6）经复核仍然成立。

---

## 1. React Bits（reactbits.dev）调研

### 1.1 定位

免费、开源的**动效 React 组件库**，作者 David Haz（GitHub `DavidHDev/react-bits`）。官方自述：「最大的开源动画、交互、可定制 React & Next.js 组件库，为人类和 AI Agent 而生」。

- GitHub 实时数据（2026-10-28）：约 **50.4k stars**、2.5k forks，2024-08 创建，持续活跃；
- 许可：MIT + Commons Clause —— 个人与商业项目可自由使用组件，但**不得单独出售、再分发组件本身**；
- 商业层：另有 Pro 付费层（pro.reactbits.dev）。

### 1.2 内容构成（五大类，合计百余个组件）

| 分类 | 典型组件 |
|---|---|
| Text Animations | Split Text、Blur Text、Decrypted Text、Glitch Text、打字机、Count Up |
| Animations | 电流边框、液态金属、光标粒子拖尾、点击火花、像素溶解、Logo 跑马灯 |
| Components | 3D 环形轮播、景深卡片轨道、WebGL 图片熔融切换器、全息卡片 |
| Backgrounds | 各类 WebGL/Canvas 粒子与着色器动态背景 |
| Micro Interactions | 光标、悬停、点击级别的细节反馈 |

### 1.3 特色

1. **源码交付，四种风味**：每个组件提供 JS/TS × CSS/Tailwind 四个版本，经 shadcn CLI 安装（`npx shadcn@latest add @react-bits/SplitText-TS-TW`），装进来的源码归项目所有、可任意改；
2. **深度 AI 友好**：组件级 "Copy for AI"、文档页 "Copy as Markdown"、shadcn MCP server 接入、全库 `llms.txt`；
3. **配套免费工具**：Background Studio（调参背景并导出代码/图片/视频）、Center Stage（截图转产品宣传图）、Shape Magic（3D 质感 Logo）、Texture Lab（160 种滤镜）、Sound Lab（UI 音效）。

### 1.4 适用场景判断

适合：产品官网、Landing Page、作品集、活动页、AI 辅助快速原型。
不适合：后台管理系统、表单表格密集应用（它不是 Ant Design / shadcn 那类通用 UI 库）、对包体积极敏感的页面（WebGL 组件需按需懒加载）。

**对本项目的结论：只能当点缀层，且必须用 Vue 版（见 §2）。**

---

## 2. 关键事实：本项目是 Vue 3，React 库全部不可直接用

`frontend/package.json` 实测：

- **Vue 3.5 + TypeScript + Vite 5 + Tailwind CSS 3 + Pinia + vue-i18n**；
- 规模：**88 个视图（views）+ 217 个组件（components）**，另有 features 两个独立模块（channel-monitor-v2、prompt-audit）；
- 已有独立设施：chart.js/vue-chartjs（图表）、@tanstack/vue-virtual（虚拟滚动）、vue-draggable-plus（拖拽）——这些**不需要**组件库提供；
- 构建门禁：`pnpm build` = `check:i18n`（i18n key 完整性）+ `vue-tsc` + `vite build`（产物嵌入后端）。

因此：

| 原讨论对象 | 本项目可用性 | Vue 等价物 |
|---|---|---|
| React Bits | ❌ 不可用 | **vue-bits**（vue-bits.dev，同一作者，2026-10-28 实测 4.6k★，活跃） |
| Ant Design（React） | ❌ 不可用 | Ant Design Vue（21.7k★） |
| shadcn/ui（React） | ❌ 不可用 | shadcn-vue（10.7k★） |
| — | — | Element Plus（27.8k★）、Naive UI（18.6k★） |

---

## 3. 项目需求画像（实证）

对现有代码的扫描结果，决定了「我们需要组件库提供什么」：

### 3.1 已有一层手搓基础组件（`src/components/common/`，约 25 个）

使用广度 TOP（引用文件数）：

| 组件 | 被引用文件数 | 现状 |
|---|---|---|
| Select | 87 | 639 行，手搓下拉 |
| Input | 85 | 103 行，API 干净（prefix/suffix 插槽、error/hint、expose focus） |
| BaseDialog | 79 | 168 行，Teleport 实现 |
| Toggle | 39 | 29 行 |
| Pagination | 31 | — |
| ConfirmDialog | 29 | — |
| DataTable | 27 | **1152 行**，能力天花板明显 |
| EmptyState | 22 | — |
| LoadingSpinner | 16 | — |
| HelpTooltip | 15 | — |
| Toast | 8 | — |

**含义**：页面代码今天只依赖这层自研 API，重构时只要 `App*` 封装层保持同名同 props、内部换成组件库实现，页面级改动可以很小——这是渐进式策略成立的前提。

### 3.2 交互原语需求

47 处原生 `<table>`、63 处 Teleport；表单校验逻辑散落各页无统一机制；日期范围、Toast、Tooltip、Upload 均有真实使用。**结论：基础层必须是全功能组件库**（表格/表单/日期/上传是手搓的无底洞），shadcn-vue 式自建对单人 fork 不可接受。

---

## 4. 候选方案对比（2026-10-28 数据）

| 维度 | Element Plus | Naive UI | Ant Design Vue | shadcn-vue 自建 | 维持手搓 |
|---|---|---|---|---|---|
| 组件完整度（表格/表单/日期/上传） | **9.5**（el-table-v2、async-validator） | 9（n-data-table 需 render 函数） | 8.5 | 5 | 4 |
| AI 语料/文档覆盖（agent 重构效率） | **10**（Vue 库中最高） | 7 | 8 | 7 | 2 |
| 主题贴合现有 Tailwind 令牌 | 7.5（CSS 变量覆盖，路径成熟） | **9**（JS 对象直连） | 6（less 变量，改造最重） | 9 | **10** |
| TS 类型体验（vue-tsc 门禁） | 8 | **9.5** | 7 | 8 | 5 |
| 与现有 SFC template 范式契合 | **9**（el-table-column 模板式） | 6.5（列定义要 h() 渲染函数） | 8 | 8 | 10 |
| 暗色模式 | 8.5（官方 dark 变量，配 html.dark） | **9**（darkTheme 一行） | 7 | 8 | 5 |
| 长期维护/社区 | **9.5**（27.8k★，最活跃） | 8（18.6k★） | 7.5（跟随 React 版滞后） | 8.5 | 2 |
| **加权总分（密集管理平台场景）** | **≈9.0** | ≈7.9 | ≈7.4 | ≈6.9 | ≈4.8 |

### 4.1 推荐：Element Plus

1. **本轮重构主要劳动力是 agent**：EP 语料覆盖最高，批量迁移 88 个视图时首次生成正确率显著更高、返工更少；
2. **现有 `common/` 层是天然迁移缓冲**（见 §3.1）；
3. **EP 的弱项本项目都用不上**：图表/虚拟滚动/拖拽已有独立设施，不需要 EP 出场；
4. Naive UI 更优的场景：人工手写为主 + 极度在意运行时主题切换干净度。本任务形态（批量重构）下语料权重最高，故选 EP。
5. 明确排除：Ant Design Vue（视觉改造最重、跟随滞后）；shadcn-vue 自建（表格/日期硬骨头工作量数倍）。

### 4.2 目标架构：四层

```
设计令牌层    tailwind.config.js（lime 主色 + 冷灰墨色系）——保留并扩展
基础组件层    Element Plus + 主题映射（element-theme.css，亮/暗双套）
封装层        src/components/ui/App*.vue ←─ 页面代码只认这层，不直接 import EP
点缀层        vue-bits 按需源码引入 ←─ 仅 Home/登录/公开页/空状态/404
```

工程接入要点：

```bash
pnpm add element-plus @element-plus/icons-vue
pnpm add -D unplugin-vue-components unplugin-auto-import
```

- vite 配 `ElementPlusResolver` 按需引入，避免全量打包；
- EP locale 跟随 vue-i18n 语言切换（自带校验文案国际化）；
- **已知坑**：`style.css` 的全局 `* { border-gray-200 }` 会污染 EP 组件边框，迁移时收窄作用域；
- 不碰的部分：chart.js、@tanstack/vue-virtual、vue-draggable-plus、验证码组件（Turnstile/阿里云/腾讯）、支付 SDK 组件（Stripe/Airwallex）。

---

## 5. 当前 UI 设计风格评价：约 7.5/10，是资产而非负债

### 5.1 做得好的（重构中必须保留的设计资产）

- **令牌层有真实设计工程**：lime 主色 `#6f9e12` 附对比度实测注释（500 对白字 3.19，600/700 落 AA 4.5+），并成文使用纪律（「只用于链接、选中条、图表首系列、焦点环；不作为大面积背景或主按钮底色」）；亮暗色温统一为同一冷灰墨黑族，解决过主题切换断层；
- **完整语义类体系**（`style.css` 870 行）：`.btn` 8 变体、`.card`/`.glass`/`.stat-card`（icon 四色变体+趋势）、`.badge` 6 变体、`.table`、skeleton、dropdown；
- **氛围层质感**：glass 阴影、lime glow、mesh-gradient 背景、8 组自定义动效曲线、tabular-nums 数字排版、卡片 10px/输入 12px 圆角，暗色模式一等公民。

整体属于**现代 SaaS 管理台风（Linear/Vercel 一派）**，不是「没有设计」。

### 5.2 真实短板（重构的正当动机）

- 一致性靠纪律不靠机制：88 视图 × 217 组件，drift 必然存在；
- 手搓组件能力天花板：DataTable 已 1152 行、Select 639 行，键盘导航/ARIA/虚拟化边角是无底洞；
- 表单校验散落各页，无统一机制。

---

## 6. Element Plus 默认风格 vs 现有风格：差距大，但可控

| 维度 | EP 默认 | 当前系统 | 开箱差距 |
|---|---|---|---|
| 主色 | `#409EFF` 标准蓝 | lime 品牌色+使用纪律 | 🔴 品牌感全丢 |
| 圆角 | 4px | 卡片 10px / 输入 12px | 🔴 明显更陈旧 |
| 质感 | 平面、边框优先 | 分层阴影+glass+glow | 🔴 差一代 |
| 气质 | 经典国内后管（阿里云控制台风） | 现代 SaaS 产品 | 🔴 **约 7-8/10 观感差距** |

- **直接换等于视觉降级**；但 EP 视觉层几乎全由 CSS 变量驱动（`--el-color-primary`、radius、border、字重），一次系统映射可拉回 **85-90% 现有视觉语言**，目标是「看不出是 EP」；
- **不可皮肤化的剩余差距（约 2-3/10）**：交互几何与组件解剖——日期面板布局、下拉结构、表格 hover 行为、Message/Notification 形状、表单校验排布。换任何主题都消不掉，需接受；
- **收益重定位**：本次重构的主要收益不在「变好看」（现在不丑），而在**可维护性**——删掉千行手搓组件、获得无障碍与键盘导航、收敛校验机制。

---

## 7. 落地路径与周期估算

### 7.1 阶段划分

| 阶段 | 内容 | 工期（agent 辅助） |
|---|---|---|
| **Phase 0** 地基 | 依赖+vite 按需配置+EP 主题映射（对齐现有令牌，亮/暗双套）+ `ui/` 封装层骨架（约 15 个 App* 组件，API 对齐现有 common 层；Select/DataTable 最重） | 1-1.5 周 |
| **Phase 1** 试点 | 用户侧 1 页 + 管理侧 1 页，锁定 API，过 `pnpm build` | 2-3 天 |
| **Phase 2** 批量迁移 | 88 视图：~40 简单（0.5-1h/页）+ ~30 中等（1-2h）+ ~18 复杂管理页如 Groups/Accounts/Channels（2-4h） | 3-5 周 |
| **Phase 3** 收尾 | 暗色全量走查、EP locale 接 vue-i18n、a11y 抽查、vue-bits 上高视觉页、删除旧 common 基础件 | 1 周 |

**总计：全职约 5-7 周；兼职（每天 2-3h）约 8-12 周。** 复杂度分布为按目录结构推断，Phase 1 试点后应修正单页耗时。

> 本节是 10-28 的**原始估算**，保留作参考。它已被 §7.2 修订版取代：Phase 2 的全站 sweep 不做，
> 冲突面依据见 §9。A 阶段≈本节 Phase 0 的机制层部分（去掉 EP 与 15 个 `App*` 的库实现），B 阶段≈Phase 1 试点但限定在 fork 独有文件。

### 7.2 推进决策（2026-10-31 修订；同日因政策变更再修订，见 §7.4）

**修订原因**：10-28 版把「上游同步冲突」当成一次性窗口风险，实测后发现它是**永久性冲突面**（见 §9）；
同时 §4.1 选 EP 的首要理由（agent 语料覆盖度）与 §4.2 的封装层设计互相抵消。
本文档的实证部分（组件规模、令牌层、差距评估）经复核成立，被推翻的只是**执行形态与工期**。

#### 原决策（已作废）

- ~~Phase 2 批量迁移 88 视图，3-5 周~~ → 冲突面永久扩大，不做。
- ~~迁移窗口内冻结功能二开~~ → 单人 fork、`main` 即生产分支，冻结 5-7 周成本高于收益。
- ~~vue-bits 引入限高视觉页~~ → 首页/登录页已按 ink + lime + 玻璃定稿，动效点缀层与既定口味（无 AI 味、管理端不做挂载动画）相抵，当前无缺口。
- ~~Phase 3 删除旧 common 基础件~~ → 上游 90 天内对 `components/common/` 有 60 次提交，删除等于自造冲突。

#### 修订后：三阶段，不引库的先行，引库的限定地盘

| 阶段 | 内容 | 工期 | 上游冲突面 |
|---|---|---|---|
| **A 机制层**（不引库，无悔棋） | ① 统一表单校验机制（规则生成器 + `useForm`），新页强制、旧页「碰到才换」；② 收窄 `style.css:7` 的全局 `* { @apply border-gray-200 }` 到自有 class；③ 给 `DataTable`/`Select` 补键盘导航与 ARIA，虚拟滚动复用已有 `@tanstack/vue-virtual`；④ 建 `src/components/ui/App*.vue`，**先只转发到现有 common 组件** | 2-4 天 | ≈0 |
| **B 自有地盘**（EP 判据阶段） | 引入 Element Plus + 亮/暗主题映射，但**只在 fork 独有的文件里使用**：`views/admin/UpstreamHubView.vue`、`views/admin/AggregationView.vue`、`components/admin/group/GroupHubEntriesField.vue`、`components/modelPlaza/PlazaModelGrid.vue`、`components/modelPlaza/PlazaModelDialog.vue`（上游没有这些文件，cherry-pick 碰不到），上游同名文件一律不动 | ~1 周 | ≈0 |
| **C 机会式迁移**（长期，不设窗口） | 仅当某视图本来就要因需求/缺陷改动、且与上游无未合差异时，顺手迁到 `App*`；一页一提交，并在 `UPSTREAM_DECISIONS_CN.md` 记「此页已私有化」，后续 cherry-pick 按 ours 判定 | 持续 | 每迁一页 +1 |

**A 是无条件要做的**：§5.2 列的三笔债（千行 DataTable、639 行 Select、校验散落）换不换库都得还；
`App*` 转发层立起来后，将来无论换 EP、Naive UI 还是继续手搓，改动面都收缩到约 15 个封装件。

**B 是全站迁移的前置判据**：主题映射能否做到「看不出是 EP」，只有在没有上游冲突压力的自有页面上才测得准。
还原度不达标就停在 A，全站方案自动作废，损失 1 周而非 7 周。

**C 的停止条件**：当上游 90 天视图 churn 覆盖不到已迁页面时再扩大；否则长期停在 B。

#### 明确不做

- 全站 88 视图 sweep；为「AI 语料覆盖度」选型；现在引入 vue-bits / WebGL 背景；删除 `common/` 基础件。

#### 何时重开「全站迁移」这个决定（本小节已被 §7.4 取代，保留作推理痕迹）

1. ~~上游连续两个 release 的视图 churn <10 个文件~~ → 政策变更后不再是条件：我们已不跟踪上游 churn；
2. 或 `DataTable`/`Select` 的自维护成本实测超过 1 周/季（**仍成立**：这是自研路线的真实痛点阈值）；
3. 或 B 阶段主题映射还原度实测 ≥90% 且页面视觉确有收益（**仍成立，且现在是唯一的技术闸门**）。

### 7.3 风险与治理

| 风险 | 应对 |
|---|---|
| **上游同步冲突（原最大变量）**：Phase 2 期间 88 个视图都是 diff | **已因政策变更失效**：2026-10-31 起本 fork 完全脱离上游节奏、永不 `git merge upstream/main`、不再逐提交 cherry-pick，重写视图不再产生长期冲突面（见 §7.4）。剩下的整合成本参考 §9 |
| `pnpm build` 门禁（i18n+vue-tsc+嵌入构建） | 每页迁移完必须过门禁，无「批量改完再验」捷径；EP locale 文案、被删组件引用会被门禁拦截 |
| 全局样式污染（`* { border-gray-200 }` 等） | Phase 0 收窄作用域 |
| 支付/验证码组件回归风险 | 不在重构范围，明确排除 |
| EP 气质残留 | 接受 2-3/10 不可消除差距，主题映射做到 85-90% 还原 |

### 7.4 政策变更后的再评估（2026-10-31 二次修订，以本节为最终口径）

同日 fork 的上游策略从「解耦 + 逐提交选择性采纳」改为「**完全解耦 + 按需借鉴**」
（`AGENTS.md`、`docs/FORK_MAINTENANCE_CN.md`、`docs/UPSTREAM_DECISIONS_CN.md` 已同步）。
这正好抽掉了 §7.2 的否决项：

- 不再 merge 上游、不再逐提交采纳 ⇒ **重写 89 个视图不再产生长期冲突面**。上游改过的视图只有在我们要它某个功能时才需要人工整合，那是主动选择，不是被动义务。
- 因此 §9.3 的推论 1（冲突面永久扩到 89）与推论 2（窗口撞上 4-7 个 release）**不再构成否决**；§9.1/§9.2 的数字仍然有效，但用途从「否决依据」降级为「整合成本参考」。

拿掉冲突代价后，全站迁移的账只剩三项，且都未被证伪：

| 项 | 结论 |
|---|---|
| 工期 5-7 周（全职）/ 8-12 周（兼职） | **仍是主要代价**。单人 fork、`main` 即生产分支，这段时间的功能开发机会成本不因解耦而消失 |
| EP 气质残留 2-3/10 | 仍需接受；能否压到可接受只有实测知道 |
| 可维护性收益（删千行手搓件、a11y、键盘导航、统一校验） | 真实，但其中「统一校验」与「a11y/键盘导航」在 A 阶段**不引库也能拿到大半** |

**最终推荐顺序**（替代 §7.2 的 A/B/C）：

1. **A 机制层照做**（2-4 天，无悔，且是 B 的前置）：统一校验机制、收窄全局 `* { border }`、补 `DataTable`/`Select` 的键盘导航与 ARIA、建 `ui/App*.vue` 转发层。
2. **B 试点范围扩大**：不再受限于 fork 独有的 5 个文件，可以直接挑 2-3 个我们**本来就打算自己重做**的密集管理页（`GroupsView` / `AccountsView` / `ChannelsView` 任一），把「主题映射还原度 ≥90%」实测出来。
3. **C 由「机会式」升级为「可计划」**：B 达标且你愿意投入工期 → 全站 sweep 按模块排期做，每模块一提交；B 不达标 → 永久停在「EP 只用于新页面 + 自研基础件继续维护」的混合态，A 阶段的收益已经拿到手。

一句话：**政策变更把全站迁移从「不值得」改成了「取决于你怎么看 5-7 周的机会成本」**——工程阻碍消失，剩下的是产品优先级问题，而这个由你定，不再由上游节奏定。

---

## 8. 附录：调研数据来源

- reactbits.dev 首页/`/get-started/introduction`/`/get-started/index` 实抓（2026-10-28）；
- GitHub API（2026-10-28）：DavidHDev/react-bits（50,461★）、DavidHDev/vue-bits（4,621★）、unovue/shadcn-vue（10,741★）、element-plus/element-plus（27,826★）、tusen-ai/naive-ui（18,589★）、vueComponent/ant-design-vue（21,695★）；
- 本项目实测：`frontend/package.json`、`frontend/tailwind.config.js`、`frontend/src/style.css`（870 行）、`frontend/src/components/common/`（含 DataTable.vue 1152 行、Select.vue 639 行）、views/components 计数、基础件引用广度扫描。

---

## 9. 冲突面实测（2026-10-31）

§7.2 第一版修订建立在这一节的数字上。**注意前提已在同日变更**：本 fork 之后不再逐提交采纳上游，
因此下面的冲突面数字不再是「否决依据」，而是「若将来要整合某个上游功能，人工成本有多大」的参考。
最终口径见 §7.4。测量口径：`git fetch --prune upstream` 后，
以上游 `main` 与本地 `main` 的 merge-base 为参照点，按路径统计提交与文件集合。

**样本点**：`upstream/main = d851a3847`、`origin/main = a17f27d6b`、merge-base `b0a02a330`（= v0.2.17，2026-10-16）。

### 9.1 上游节奏

| 指标 | 90 天 | 30 天 | 自基线（15 天） |
|---|---|---|---|
| 上游提交总数 | 2348 | 422 | — |
| 触及 `frontend/src/views` 的提交 | 287 | 48 | — |
| 被改过的**视图文件**（去重） | **58 / 89** | 20 | 7 |
| 触及 `frontend/src/components/common/` 的提交 | 60 | — | 0 |

- 发布密度：v0.2.14（10-07）→ v0.2.19（10-28），**21 天 6 个 release**；基线之后已出 v0.2.18、v0.2.19。
- 自基线以来上游改动 422 个文件，其中 `frontend/` 63 个（≈15%）。
- `common/` 层 90 天 60 次提交的分布：`PlatformTypeBadge` 8、`Select` 5、`PlatformIcon` 5、`BaseDialog` 4、`GroupSelector` 3、`GroupBadge` 3。

### 9.2 fork 现状（同一基线）

| 指标 | 值 |
|---|---|
| 视图总数 / 组件总数（不含 `__tests__`） | 89 / 228（§2 记的 88/217 是 10-28 当时的值） |
| 与基线相比有差异的视图文件 | **18 / 89**（含已采纳的上游改动，非全是私有改动） |
| fork 独有、上游根本没有的前端文件 | **16** |
| └ 其中可承载 UI 的 | **5**：`views/admin/UpstreamHubView.vue`、`views/admin/AggregationView.vue`、`components/admin/group/GroupHubEntriesField.vue`、`components/modelPlaza/PlazaModelGrid.vue`、`components/modelPlaza/PlazaModelDialog.vue` |
| └ 其余 11 个 | 语言包 4（`{en,zh}/admin/{aggregation,upstreamHub}.ts`）、测试 5（`*.spec.ts`）、非可视模块 2（`plazaPricing.ts`、`api/admin/upstreamHub.ts`） |
| 当前 dist 体积（未压缩） | 6.6 MB |

### 9.3 三条推论

1. ~~**全站 sweep 的冲突面是永久性的。**~~ **（2026-10-31 政策变更后失效）** 原推理：迁移 89 个视图后，上游在任意 90 天窗口内会改动其中 58 个，而逐提交 `cherry-pick -x` 意味着此后几乎每个前端上游提交都要重解模板冲突。现在不再逐提交采纳，该成本不再存在；数字保留作整合成本参考。
2. ~~**窗口本身的代价。**~~ **（失效）** 原推理基于「迁移期还要边迁边采纳上游批次」。解耦后不存在这个耦合，5-7 周剩下的纯粹是功能开发的机会成本。
3. **B 阶段冲突面为 0。** fork 独有的 16 个前端文件（其中可承载 UI 的只有 5 个）在上游没有对应行，cherry-pick 碰不到，
   所以 EP 与主题映射可以在这块地上无风险实测——这正是把它设成判据阶段的原因。代价是场地很小：5 个文件不足以覆盖表格/表单/日期/上传全谱，
   B 阶段的结论只能证明「EP 主题映射能不能贴近现有令牌」，不能证明「全站迁得动」。

### 9.4 未实测项（不得当结论引用）

- **EP 引入的包体增量**：未安装 `element-plus`，只有「`ElementPlusResolver` 按需引入会显著小于全量」这一定性判断。
- **主题映射还原度**：§6 的「85-90%」是推断值，需在 B 阶段用截图实测（浅色/深色 × 1440/390）。
- **A 阶段的 a11y 收益**：`DataTable`/`Select` 补键盘导航后的实际可访问性未做 axe 类扫描。
