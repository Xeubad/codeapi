export default {
  batchImageGuide: {
    title: '图片批量生成',
    description: '一次提交多条提示词，任务完成后可统一下载图片结果'
  },
  // Home Page
  home: {
    viewDocs: '查看文档',
    docs: '文档',
    switchToLight: '切换到浅色模式',
    switchToDark: '切换到深色模式',
    dashboard: '控制台',
    login: '登录',
    getStarted: '立即开始',
    goToDashboard: '进入控制台',
    backHome: '返回首页',
    nav: {
      latency: '实测延迟',
      integration: '快速接入'
    },
    facts: {
      models: { title: '统一入口', value: '多平台 · 单密钥' },
      sticky: { title: '粘性会话', value: '同一会话固定上游' },
      billing: { title: '实时计费', value: '按 token 逐次结算' }
    },
    hero: {
      statusLabel: '服务在线',
      probeValue: '往返 {ms} ms',
      probePending: '正在测量本机延迟',
      titleLead: '一个密钥，',
      titleAccent: '接管全部模型调用',
      description: '多家主流模型平台通过同一套接口接入；按 token 计费，登录后可在模型广场查看可用模型与价格。',
      loginCta: '登录',
      secondaryCta: 'Key 用量查询',
      primaryCtaAuthed: '进入控制台',
      note: '登录后创建 API Key 即可调用；模型与价格在模型广场查看。'
    },
    gateway: {
      panelTitle: 'GATEWAY / BASE URL',
      degraded: '连接异常',
      baseHint: '只替换 base_url，密钥与其余参数保持不变。',
      copyBase: '复制 base_url',
      modelInfo: '模型与价格',
      modelInfoNote: '登录后在模型广场查看',
      modelInfoAction: '登录查看',
      modelInfoActionAuthed: '进入控制台'
    },
    latency: {
      eyebrow: 'LIVE / 实测',
      title: '你的访问延迟，是真的',
      description: '从当前浏览器向本站静态资源发起真实请求，每 5 秒采样一次。这里没有模拟数据，数字就是此刻的网络往返。',
      current: '当前往返',
      avg: '平均',
      p95: 'P95',
      jitter: '抖动',
      sampleCount: '采样',
      ttfb: '首字节',
      connection: '连接',
      reused: '复用',
      fresh: '新建',
      transfer: '传输',
      protocol: '协议',
      window: '已采样 {count} 次',
      recent: '最近样本',
      waiting: '正在积累样本',
      failed: '测量失败，正在重试',
      fast: '流畅',
      normal: '正常',
      slow: '有波动',
      measuring: '测量中',
      retry: '立即重测',
      source: '数据来源：当前设备 → 本站静态资源，真实往返时延。',
      noData: '暂无数据'
    },
    integration: {
      eyebrow: 'QUICK START',
      title: '改一行 base_url，就能接上',
      description: '保留原有 SDK 调用方式，只把地址换成本站。Anthropic 与 OpenAI 两种请求格式都支持。',
      compatible: '兼容 Anthropic / OpenAI SDK 的调用格式',
      copy: '复制',
      copied: '已复制',
      copyFailed: '复制失败',
      point: {
        baseUrl: '把 SDK 的 base_url 指向本站',
        key: '用控制台创建的 API Key 替换原密钥',
        done: '其余参数不变，直接发起调用'
      }
    },
    cta: {
      title: '把模型调用收拢到一个控制台',
      description: '创建密钥、看用量、查延迟，都在同一个地方完成；模型与价格在模型广场查看。',
      secondary: '查询 Key 用量'
    },
    footer: {
      allRightsReserved: '保留所有权利。',
      status: '服务运行中',
      keyUsage: 'Key 用量查询',
      login: '登录'
    },
    auth: {
      secureNote: '登录态由安全 Cookie 与会话令牌维护；验证码、双因素与登录协议校验保持不变。'
    },
    console: {
      billedBy: '实时计费',
      oneKey: '单密钥多模型',
      stickySession: '命中同一上游账号',
      realtimeBilling: '按 token 实时结算'
    },
  },

  // Key Usage Query Page
  keyUsage: {
    title: 'API Key 用量查询',
    subtitle: '输入您的 API Key 以查看实时消费金额与使用状态',
    placeholder: 'sk-ant-mirror-xxxxxxxxxxxx',
    query: '查询',
    querying: '查询中...',
    privacyNote: '您的 Key 仅在浏览器本地处理，不会被存储',
    dateRange: '统计范围:',
    dateRangeToday: '今日',
    dateRange7d: '7 天',
    dateRange30d: '30 天',
    dateRange90d: '90 天',
    dateRangeCustom: '自定义',
    apply: '应用',
    used: '已使用',
    detailInfo: '详细信息',
    tokenStats: 'Token 统计',
    dailyDetail: '按日明细',
    modelStats: '模型用量统计',
    // Table headers
    date: '日期',
    model: '模型',
    requests: '请求数',
    inputTokens: '输入 Tokens',
    outputTokens: '输出 Tokens',
    cacheCreationTokens: '缓存创建',
    cacheReadTokens: '缓存读取',
    cacheWriteTokens: '缓存写入',
    totalTokens: '总 Tokens',
    cost: '费用',
    // Status
    quotaMode: 'Key 限额模式',
    walletBalance: '钱包余额',
    // Ring card titles
    totalQuota: '总额度',
    limit5h: '5 小时限额',
    limitDaily: '日限额',
    limit7d: '7 天限额',
    limitWeekly: '周限额',
    limitMonthly: '月限额',
    // Detail rows
    remainingQuota: '剩余额度',
    expiresAt: '过期时间',
    todayExpires: '(今日到期)',
    daysLeft: '({days} 天)',
    usedQuota: '已用额度',
    resetNow: '即将重置',
    subscriptionType: '订阅类型',
    billingType: '计费方式',
    subscriptionExpires: '订阅到期',
    // Usage stat cells
    todayRequests: '今日请求',
    todayInputTokens: '今日输入',
    todayOutputTokens: '今日输出',
    todayTokens: '今日 Tokens',
    todayCacheCreation: '今日缓存创建',
    todayCacheRead: '今日缓存读取',
    todayCost: '今日费用',
    rpmTpm: 'RPM / TPM',
    totalRequests: '累计请求',
    totalInputTokens: '累计输入',
    totalOutputTokens: '累计输出',
    totalTokensLabel: '累计 Tokens',
    totalCacheCreation: '累计缓存创建',
    totalCacheRead: '累计缓存读取',
    totalCost: '累计费用',
    avgDuration: '平均耗时',
    // Messages
    enterApiKey: '请输入 API Key',
    querySuccess: '查询成功',
    queryFailed: '查询失败',
    queryFailedRetry: '查询失败，请稍后重试',
    noDailyUsage: '暂无按日用量数据',
  },

  // Setup Wizard
  setup: {
    title: 'Sub2API 安装向导',
    description: '配置您的 Sub2API 实例',
    database: {
      title: '数据库配置',
      description: '连接到您的 PostgreSQL 数据库',
      host: '主机',
      port: '端口',
      username: '用户名',
      password: '密码',
      databaseName: '数据库名称',
      sslMode: 'SSL 模式',
      passwordPlaceholder: '密码',
      ssl: {
        disable: '禁用',
        require: '要求',
        verifyCa: '验证 CA',
        verifyFull: '完全验证'
      }
    },
    redis: {
      title: 'Redis 配置',
      description: '连接到您的 Redis 服务器',
      host: '主机',
      port: '端口',
      username: '用户名（可选）',
      password: '密码（可选）',
      database: '数据库',
      usernamePlaceholder: '默认用户留空',
      passwordPlaceholder: '密码',
      enableTls: '启用 TLS',
      enableTlsHint: '连接 Redis 时使用 TLS（公共 CA 证书）'
    },
    admin: {
      title: '管理员账户',
      description: '创建您的管理员账户',
      email: '邮箱',
      password: '密码',
      confirmPassword: '确认密码',
      passwordPlaceholder: '至少 8 个字符',
      confirmPasswordPlaceholder: '确认密码',
      passwordMismatch: '密码不匹配'
    },
    ready: {
      title: '准备安装',
      description: '检查您的配置并完成安装',
      database: '数据库',
      redis: 'Redis',
      adminEmail: '管理员邮箱'
    },
    status: {
      testing: '测试中...',
      success: '连接成功',
      testConnection: '测试连接',
      installing: '安装中...',
      completeInstallation: '完成安装',
      completed: '安装完成！',
      redirecting: '正在跳转到登录页面...',
      restarting: '服务正在重启，请稍候...',
      timeout: '服务重启时间超出预期，请手动刷新页面。'
    }
  },

  // Common
}
