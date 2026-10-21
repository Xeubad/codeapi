export default {
  batchImageGuide: {
    title: 'Batch Image Generation',
    description: 'Submit multiple prompts in one job and download the generated images when complete'
  },
  // Home Page
  home: {
    viewDocs: 'View Documentation',
    docs: 'Docs',
    switchToLight: 'Switch to Light Mode',
    switchToDark: 'Switch to Dark Mode',
    dashboard: 'Dashboard',
    login: 'Login',
    getStarted: 'Get Started',
    goToDashboard: 'Go to Dashboard',
    backHome: 'Back to home',
    nav: {
      capabilities: 'Capabilities',
      latency: 'Live latency',
      integration: 'Quick start',
      models: 'Models'
    },
    hero: {
      statusLabel: 'Service online',
      probeValue: 'round trip {ms} ms',
      probePending: 'Measuring this device',
      titleLead: 'One key.',
      titleAccent: 'Runs every model call you make',
      description: 'Claude, GPT and Gemini behind one endpoint. Account pooling, sticky sessions and per-token billing all live in the same console.',
      primaryCta: 'Get started',
      primaryCtaAuthed: 'Open console',
      secondaryCta: 'See live latency',
      note: 'Create an API key in the console after signing in. Every call and every cost stays visible.'
    },
    facts: {
      models: { title: 'One endpoint', value: 'Claude · GPT · Gemini' },
      sticky: { title: 'Sticky sessions', value: 'Same session, same upstream' },
      billing: { title: 'Live billing', value: 'Settled per token' }
    },
    preview: {
      console: 'console / overview',
      badge: 'UI preview',
      overview: 'Today',
      balance: 'Balance',
      todayRequests: 'Requests today',
      todayCost: 'Cost today',
      rpm: 'RPM',
      recentCalls: 'Recent calls',
      fiveHour: '5 hours',
      daily: 'Daily',
      replay: 'Sticky session',
      replayNote: 'Same session hits one upstream'
    },
    latency: {
      eyebrow: 'LIVE / MEASURED',
      title: 'Your latency here is real',
      description: 'This browser requests a static asset from this site every 5 seconds. No simulated numbers: what you see is the actual round trip.',
      current: 'Current round trip',
      avg: 'Average',
      p95: 'P95',
      jitter: 'Jitter',
      sampleCount: 'Samples',
      ttfb: 'Time to first byte',
      connection: 'Connection',
      reused: 'Reused',
      fresh: 'New',
      transfer: 'Transfer',
      protocol: 'Protocol',
      window: 'Samples so far: {count}',
      recent: 'Recent samples',
      waiting: 'Collecting samples',
      failed: 'Probe failed, retrying',
      fast: 'Smooth',
      normal: 'Normal',
      slow: 'Unstable',
      measuring: 'Measuring',
      retry: 'Probe now',
      source: 'Source: this device → this site\'s static asset, real round-trip time.',
      noData: 'No data yet'
    },
    capabilities: {
      eyebrow: 'CONTROL ROOM',
      title: 'What already lives in the console',
      subtitle: 'Not a concept sketch. Every block below maps to a page you can open after signing in.',
      mockNote: 'Illustration · sample data',
      routing: {
        title: 'Account pool routing',
        desc: 'Several upstream accounts take requests in turn, failover is automatic, and nothing needs manual rewiring.',
        sticky: 'Sticky',
        pool: 'Upstream pool'
      },
      billing: {
        title: 'Per-token billing',
        desc: 'Each call settles immediately, with input, output and cache accounted separately.',
        input: 'Input',
        output: 'Output',
        cache: 'Cache',
        cost: 'This call'
      },
      quota: {
        title: 'Quotas and limits',
        desc: 'Five-hour, daily and weekly limits are independent, and the trend is visible before you hit them.',
        fiveHour: '5 hours',
        daily: 'Daily',
        weekly: 'Weekly'
      },
      audit: {
        title: 'Call detail and audit',
        desc: 'Every request keeps its time, key, model, status and latency for reconciliation and debugging.',
        time: 'Time',
        key: 'Key',
        model: 'Model',
        status: 'Status',
        latency: 'Latency'
      }
    },
    integration: {
      eyebrow: 'QUICK START',
      title: 'Change one base_url and you are in',
      description: 'Keep your existing SDK calls and point them at this site. Anthropic and OpenAI request formats are both accepted.',
      compatible: 'Works with Anthropic / OpenAI SDK request formats',
      copy: 'Copy',
      copied: 'Copied',
      copyFailed: 'Copy failed',
      point: {
        baseUrl: 'Point the SDK base_url at this site',
        key: 'Swap in a key created in the console',
        done: 'Leave the rest untouched and send the call'
      }
    },
    models: {
      title: 'Models behind one key',
      description: 'Supported models keep growing, and the console shows the live availability of each one.'
    },
    cta: {
      title: 'Bring every model call into one console',
      description: 'Create keys, top up, read usage and check latency in the same place.',
      primaryGuest: 'Sign up and start',
      secondary: 'Look up key usage'
    },
    footer: {
      allRightsReserved: 'All rights reserved.',
      status: 'Service online',
      keyUsage: 'Key usage lookup',
      login: 'Sign in'
    },
    auth: {
      secureNote: 'Sessions are held by secure cookies and tokens; captcha, two-factor and login agreement checks stay in place.'
    },
    console: {
      billedBy: 'Billed live',
      oneKey: 'One key, many models',
      stickySession: 'Stays on one upstream',
      realtimeBilling: 'Billed per token'
    },
    providers: {
      title: 'Supported AI models',
      description: 'One API, many choices',
      supported: 'Supported',
      soon: 'Coming soon',
      claude: 'Claude',
      gemini: 'Gemini',
      antigravity: 'Antigravity',
      more: 'More'
    }
  },

  // Key Usage Query Page
  keyUsage: {
    title: 'API Key Usage',
    subtitle: 'Enter your API Key to view real-time spending and usage status',
    placeholder: 'sk-ant-mirror-xxxxxxxxxxxx',
    query: 'Query',
    querying: 'Querying...',
    privacyNote: 'Your Key is processed locally in the browser and will not be stored',
    dateRange: 'Date Range:',
    dateRangeToday: 'Today',
    dateRange7d: '7 Days',
    dateRange30d: '30 Days',
    dateRange90d: '90 Days',
    dateRangeCustom: 'Custom',
    apply: 'Apply',
    used: 'Used',
    detailInfo: 'Detail Information',
    tokenStats: 'Token Statistics',
    dailyDetail: 'Daily Detail',
    modelStats: 'Model Usage Statistics',
    // Table headers
    date: 'Date',
    model: 'Model',
    requests: 'Requests',
    inputTokens: 'Input Tokens',
    outputTokens: 'Output Tokens',
    cacheCreationTokens: 'Cache Creation',
    cacheReadTokens: 'Cache Read',
    cacheWriteTokens: 'Cache Write',
    totalTokens: 'Total Tokens',
    cost: 'Cost',
    // Status
    quotaMode: 'Key Quota Mode',
    walletBalance: 'Wallet Balance',
    // Ring card titles
    totalQuota: 'Total Quota',
    limit5h: '5-Hour Limit',
    limitDaily: 'Daily Limit',
    limit7d: '7-Day Limit',
    limitWeekly: 'Weekly Limit',
    limitMonthly: 'Monthly Limit',
    // Detail rows
    remainingQuota: 'Remaining Quota',
    expiresAt: 'Expires At',
    todayExpires: '(expires today)',
    daysLeft: '({days} days)',
    usedQuota: 'Used Quota',
    resetNow: 'Resetting soon',
    subscriptionType: 'Subscription Type',
    billingType: 'Billing Type',
    subscriptionExpires: 'Subscription Expires',
    // Usage stat cells
    todayRequests: 'Today Requests',
    todayInputTokens: 'Today Input',
    todayOutputTokens: 'Today Output',
    todayTokens: 'Today Tokens',
    todayCacheCreation: 'Today Cache Creation',
    todayCacheRead: 'Today Cache Read',
    todayCost: 'Today Cost',
    rpmTpm: 'RPM / TPM',
    totalRequests: 'Total Requests',
    totalInputTokens: 'Total Input',
    totalOutputTokens: 'Total Output',
    totalTokensLabel: 'Total Tokens',
    totalCacheCreation: 'Total Cache Creation',
    totalCacheRead: 'Total Cache Read',
    totalCost: 'Total Cost',
    avgDuration: 'Avg Duration',
    // Messages
    enterApiKey: 'Please enter an API Key',
    querySuccess: 'Query successful',
    queryFailed: 'Query failed',
    queryFailedRetry: 'Query failed, please try again later',
    noDailyUsage: 'No daily usage data',
  },

  // Setup Wizard
  setup: {
    title: 'Sub2API Setup',
    description: 'Configure your Sub2API instance',
    database: {
      title: 'Database Configuration',
      description: 'Connect to your PostgreSQL database',
      host: 'Host',
      port: 'Port',
      username: 'Username',
      password: 'Password',
      databaseName: 'Database Name',
      sslMode: 'SSL Mode',
      passwordPlaceholder: 'Password',
      ssl: {
        disable: 'Disable',
        require: 'Require',
        verifyCa: 'Verify CA',
        verifyFull: 'Verify Full'
      }
    },
    redis: {
      title: 'Redis Configuration',
      description: 'Connect to your Redis server',
      host: 'Host',
      port: 'Port',
      username: 'Username (optional)',
      password: 'Password (optional)',
      database: 'Database',
      usernamePlaceholder: 'Leave empty for default user',
      passwordPlaceholder: 'Password',
      enableTls: 'Enable TLS',
      enableTlsHint: 'Use TLS when connecting to Redis (public CA certs)'
    },
    admin: {
      title: 'Admin Account',
      description: 'Create your administrator account',
      email: 'Email',
      password: 'Password',
      confirmPassword: 'Confirm Password',
      passwordPlaceholder: 'Min 8 characters',
      confirmPasswordPlaceholder: 'Confirm password',
      passwordMismatch: 'Passwords do not match'
    },
    ready: {
      title: 'Ready to Install',
      description: 'Review your configuration and complete setup',
      database: 'Database',
      redis: 'Redis',
      adminEmail: 'Admin Email'
    },
    status: {
      testing: 'Testing...',
      success: 'Connection Successful',
      testConnection: 'Test Connection',
      installing: 'Installing...',
      completeInstallation: 'Complete Installation',
      completed: 'Installation completed!',
      redirecting: 'Redirecting to login page...',
      restarting: 'Service is restarting, please wait...',
      timeout: 'Service restart is taking longer than expected. Please refresh the page manually.'
    }
  },

  // Common
}
