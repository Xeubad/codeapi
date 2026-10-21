<template>
  <!-- Custom Home Content: Full Page Mode -->
  <div v-if="hasHomeContent" class="min-h-screen">
    <!-- iframe mode -->
    <iframe
      v-if="isHomeContentUrl"
      :src="homeContent.trim()"
      class="h-screen w-full border-0"
      allowfullscreen
    ></iframe>
    <!-- HTML mode - SECURITY: homeContent is admin-only setting, XSS risk is acceptable -->
    <div v-else v-html="homeContent"></div>
  </div>

  <!-- Compact Home Page -->
  <div
    v-else-if="compactHomeEnabled"
    data-testid="compact-home"
    class="flex min-h-screen flex-col bg-gray-50 text-gray-900 dark:bg-dark-950 dark:text-white"
  >
    <header class="border-b border-gray-200 px-4 py-4 sm:px-6 dark:border-dark-800">
      <nav class="mx-auto flex max-w-5xl flex-wrap items-center justify-between gap-3 sm:gap-4">
        <div class="flex min-w-0 flex-1 items-center gap-3">
          <img
            :src="siteLogo || '/logo.png'"
            alt="Logo"
            class="h-9 w-9 shrink-0 rounded-lg object-contain"
          />
          <span class="min-w-0 truncate text-base font-semibold">{{ siteName }}</span>
        </div>
        <div class="flex max-w-full shrink-0 flex-wrap items-center justify-end gap-2">
          <LocaleSwitcher />
          <router-link
            v-if="showModelPlazaEntry"
            to="/model-plaza"
            class="flex h-10 shrink-0 items-center gap-1.5 rounded-lg px-2.5 text-sm font-medium text-gray-500 hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
            :title="t('nav.modelPlaza')"
          >
            <Icon name="grid" size="md" />
            <span class="hidden sm:inline">{{ t('nav.modelPlaza') }}</span>
          </router-link>
          <button
            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg text-gray-500 hover:bg-gray-100 dark:text-dark-400 dark:hover:bg-dark-800"
            :title="isDark ? t('home.switchToLight') : t('home.switchToDark')"
            @click="toggleTheme"
          >
            <Icon v-if="isDark" name="sun" size="md" />
            <Icon v-else name="moon" size="md" />
          </button>
          <router-link
            :to="isAuthenticated ? dashboardPath : '/login'"
            class="inline-flex min-h-10 shrink-0 items-center justify-center rounded-lg bg-gray-900 px-4 py-2 text-sm font-medium text-white hover:bg-gray-800 dark:bg-white dark:text-gray-900 dark:hover:bg-gray-200"
          >
            {{ isAuthenticated ? t('home.dashboard') : t('home.login') }}
          </router-link>
        </div>
      </nav>
    </header>

    <main class="flex min-w-0 flex-1 items-center justify-center px-4 py-16 sm:px-6">
      <div class="min-w-0 max-w-2xl text-center">
        <img
          :src="siteLogo || '/logo.png'"
          alt="Logo"
          class="mx-auto mb-6 h-20 w-20 rounded-2xl object-contain"
        />
        <h1 class="[overflow-wrap:anywhere] text-3xl font-bold md:text-4xl">{{ siteName }}</h1>
        <p class="mt-4 whitespace-pre-wrap [overflow-wrap:anywhere] text-base text-gray-600 dark:text-dark-300">{{ siteSubtitle }}</p>
        <router-link
          :to="isAuthenticated ? dashboardPath : '/login'"
          class="mt-8 inline-flex min-h-10 items-center justify-center rounded-lg bg-gray-900 px-5 py-2.5 text-sm font-medium text-white hover:bg-gray-800"
        >
          {{ isAuthenticated ? t('home.goToDashboard') : t('home.login') }}
        </router-link>
      </div>
    </main>

    <footer class="min-w-0 border-t border-gray-200 px-4 py-5 text-center text-sm text-gray-500 [overflow-wrap:anywhere] sm:px-6 dark:border-dark-800 dark:text-dark-400">
      &copy; {{ currentYear }} {{ siteName }}
    </footer>
  </div>

  <!-- Promotional Home Page -->
  <div
    v-else
    data-testid="promo-home"
    class="relative flex min-h-screen flex-col overflow-x-clip bg-gray-50 text-gray-900 dark:bg-dark-950 dark:text-gray-100"
  >
    <!-- 苹果式毛玻璃背景：低饱和色场 + 柔光团 + 颗粒 -->
    <div class="brand-backdrop" aria-hidden="true">
      <div class="brand-orb -left-24 -top-32 h-[30rem] w-[30rem] bg-[rgba(111,158,18,0.20)]"></div>
      <div class="brand-orb -right-32 top-16 h-[26rem] w-[26rem] bg-[rgba(56,132,160,0.16)]"></div>
      <div class="brand-orb bottom-[-14rem] left-1/3 h-[30rem] w-[30rem] bg-[rgba(255,243,219,0.55)] dark:bg-[rgba(47,74,13,0.35)]"></div>
      <div class="brand-grain"></div>
    </div>

    <!-- 顶栏：悬浮毛玻璃 -->
    <header class="sticky top-0 z-40 px-3 pt-3 sm:px-5 sm:pt-4">
      <nav class="brand-glass-strong mx-auto flex h-14 max-w-7xl items-center gap-2 px-3 sm:px-4">
        <router-link to="/" class="flex min-w-0 items-center gap-2.5">
          <img :src="siteLogo || '/logo.png'" alt="Logo" class="h-8 w-8 shrink-0 object-contain" />
          <span class="truncate text-[15px] font-semibold tracking-tight">{{ siteName }}</span>
        </router-link>

        <div class="ml-3 hidden items-center gap-0.5 md:flex">
          <a
            v-for="item in navItems"
            :key="item.href"
            :href="item.href"
            class="rounded-lg px-2.5 py-1.5 text-[13px] font-medium text-gray-600 transition-colors hover:bg-gray-900/5 hover:text-gray-900 dark:text-dark-300 dark:hover:bg-white/5 dark:hover:text-white"
          >
            {{ item.label }}
          </a>
        </div>

        <div class="ml-auto flex shrink-0 items-center gap-1">
          <router-link
            v-if="showModelPlazaEntry"
            to="/model-plaza"
            class="hidden h-9 items-center gap-1.5 rounded-lg px-2.5 text-[13px] font-medium text-gray-600 transition-colors hover:bg-gray-900/5 hover:text-gray-900 sm:flex dark:text-dark-300 dark:hover:bg-white/5 dark:hover:text-white"
          >
            <Icon name="grid" size="sm" />
            {{ t('nav.modelPlaza') }}
          </router-link>
          <LocaleSwitcher />
          <button
            class="flex h-9 w-9 items-center justify-center rounded-lg text-gray-500 transition-colors hover:bg-gray-900/5 hover:text-gray-900 dark:text-dark-400 dark:hover:bg-white/5 dark:hover:text-white"
            :title="isDark ? t('home.switchToLight') : t('home.switchToDark')"
            @click="toggleTheme"
          >
            <Icon v-if="isDark" name="sun" size="sm" />
            <Icon v-else name="moon" size="sm" />
          </button>
          <router-link
            :to="isAuthenticated ? dashboardPath : '/login'"
            class="ml-1 inline-flex h-9 items-center gap-1.5 rounded-lg bg-gray-900 px-3.5 text-[13px] font-medium text-white transition-colors hover:bg-gray-800 dark:bg-gray-100 dark:text-gray-900 dark:hover:bg-white"
          >
            {{ isAuthenticated ? t('home.dashboard') : t('home.login') }}
            <Icon name="arrowRight" size="xs" />
          </router-link>
        </div>
      </nav>
    </header>

    <main class="relative z-10 flex-1">
      <!-- Hero -->
      <section class="mx-auto max-w-7xl px-4 pb-8 pt-10 sm:px-6 lg:pt-16">
        <div class="grid min-w-0 grid-cols-1 items-center gap-12 lg:grid-cols-[minmax(0,1.02fr)_minmax(0,0.98fr)] lg:gap-14">
          <div class="min-w-0">
            <p
              class="inline-flex items-center gap-2 rounded-full border border-gray-900/10 bg-white/50 px-3 py-1.5 backdrop-blur-md dark:border-white/10 dark:bg-white/5"
            >
              <span class="relative flex h-1.5 w-1.5">
                <span
                  class="absolute inline-flex h-full w-full animate-ping rounded-full bg-primary-500 opacity-60 motion-reduce:hidden"
                ></span>
                <span class="relative inline-flex h-1.5 w-1.5 rounded-full bg-primary-500"></span>
              </span>
              <span class="font-mono text-[10px] font-medium uppercase tracking-[0.18em] text-gray-600 dark:text-dark-300">
                {{ t('home.hero.statusLabel') }}
              </span>
              <span class="h-3 w-px bg-gray-900/10 dark:bg-white/10"></span>
              <span class="font-mono text-[10px] tabular-nums text-gray-900 dark:text-gray-100">
                {{ heroProbeText }}
              </span>
            </p>

            <h1
              class="mt-6 text-[2.35rem] font-semibold leading-[1.08] tracking-[-0.035em] break-keep text-gray-900 dark:text-white sm:text-[3.05rem] lg:text-[3.35rem]"
            >
              <span class="block">{{ t('home.hero.titleLead') }}</span>
              <span class="block">{{ t('home.hero.titleAccent') }}</span>
            </h1>
            <div class="mt-5 h-1 w-12 rounded-full bg-primary-500"></div>
            <p class="mt-5 max-w-xl text-[15px] leading-relaxed text-gray-600 dark:text-dark-300">
              {{ t('home.hero.description') }}
            </p>

            <div class="mt-8 flex flex-wrap items-center gap-3">
              <router-link
                :to="isAuthenticated ? dashboardPath : '/login'"
                class="inline-flex h-11 items-center gap-2 rounded-xl bg-gray-900 px-5 text-sm font-medium text-white transition-colors hover:bg-gray-800 dark:bg-gray-100 dark:text-gray-900 dark:hover:bg-white"
              >
                {{ isAuthenticated ? t('home.hero.primaryCtaAuthed') : t('home.hero.primaryCta') }}
                <Icon name="arrowRight" size="sm" />
              </router-link>
              <a
                href="#latency"
                class="inline-flex h-11 items-center gap-2 rounded-xl border border-white/80 bg-white/55 px-4 text-sm font-medium text-gray-800 backdrop-blur-md transition-colors hover:bg-white/80 dark:border-white/10 dark:bg-dark-900/50 dark:text-gray-200 dark:hover:bg-dark-800/80"
              >
                <Icon name="arrowDown" size="sm" />
                {{ t('home.hero.secondaryCta') }}
              </a>
            </div>

            <p class="mt-4 max-w-lg text-[12.5px] leading-relaxed text-gray-500 dark:text-dark-400">
              {{ t('home.hero.note') }}
            </p>
          </div>

          <!-- 控制台预览 -->
          <div class="relative min-w-0 lg:pl-2">
            <div
              class="pointer-events-none absolute -inset-5 rounded-[2rem] bg-[radial-gradient(70%_70%_at_30%_20%,rgba(111,158,18,0.18),transparent_70%)] blur-2xl"
            ></div>

            <div class="brand-ink-panel relative overflow-hidden">
              <div class="flex items-center gap-1.5 border-b border-white/10 px-4 py-3">
                <span class="h-2 w-2 rounded-full bg-white/15"></span>
                <span class="h-2 w-2 rounded-full bg-white/15"></span>
                <span class="h-2 w-2 rounded-full bg-white/15"></span>
                <span class="ml-2 truncate font-mono text-[10px] uppercase tracking-[0.18em] text-gray-500">
                  {{ t('home.preview.console') }}
                </span>
                <span
                  class="ml-auto shrink-0 rounded-md border border-white/10 bg-white/5 px-1.5 py-0.5 font-mono text-[9px] uppercase tracking-[0.12em] text-gray-400"
                >
                  {{ t('home.preview.badge') }}
                </span>
                <span
                  class="shrink-0 rounded-md bg-primary-500/15 px-1.5 py-0.5 font-mono text-[9px] tabular-nums text-primary-300"
                >
                  {{ heroProbeText }}
                </span>
              </div>

              <div class="grid grid-cols-[46px_minmax(0,1fr)]">
                <aside class="border-r border-white/5 py-3">
                  <div class="flex flex-col items-center gap-2">
                    <span
                      v-for="(tab, index) in previewTabs"
                      :key="tab"
                      class="flex h-7 w-7 items-center justify-center rounded-lg"
                      :class="index === 0 ? 'bg-primary-500/20 text-primary-300' : 'text-gray-600'"
                    >
                      <Icon :name="tab" size="xs" />
                    </span>
                  </div>
                </aside>

                <div class="min-w-0 p-3 sm:p-4">
                  <div class="flex items-center justify-between gap-3">
                    <div>
                      <p class="text-[13px] font-semibold text-white">{{ t('home.preview.overview') }}</p>
                      <p class="mt-0.5 font-mono text-[9px] uppercase tracking-[0.16em] text-gray-500">
                        {{ previewDate }}
                      </p>
                    </div>
                    <span
                      class="shrink-0 rounded-md border border-white/10 bg-white/5 px-2 py-1 font-mono text-[9px] text-gray-400"
                    >
                      sk-•••• 7f2a
                    </span>
                  </div>

                  <div class="mt-3 grid grid-cols-2 gap-2">
                    <div
                      v-for="stat in previewStats"
                      :key="stat.label"
                      class="rounded-xl border border-white/[0.06] bg-white/[0.03] px-3 py-2"
                    >
                      <p class="text-[9px] uppercase tracking-[0.14em] text-gray-500">{{ stat.label }}</p>
                      <p class="mt-1 text-[17px] font-semibold leading-none tabular-nums" :class="stat.accent ? 'text-primary-300' : 'text-gray-100'">
                        {{ stat.value }}
                      </p>
                    </div>
                  </div>

                  <div class="mt-3 overflow-hidden rounded-xl border border-white/[0.06]">
                    <div class="flex items-center justify-between border-b border-white/[0.06] bg-white/[0.03] px-3 py-1.5">
                      <span class="text-[10px] font-medium text-gray-400">{{ t('home.preview.recentCalls') }}</span>
                      <span class="font-mono text-[9px] text-gray-600">200 OK</span>
                    </div>
                    <div
                      v-for="row in previewRows"
                      :key="row.time + row.model"
                      class="grid grid-cols-[34px_minmax(0,1fr)_auto_42px_6px] items-center gap-2 px-3 py-1.5 text-[10px] odd:bg-white/[0.015]"
                    >
                      <span class="font-mono text-gray-600">{{ row.time }}</span>
                      <span class="truncate font-mono text-gray-300">{{ row.model }}</span>
                      <span class="font-mono text-gray-500">{{ row.tokens }}</span>
                      <span class="text-right font-mono text-primary-300">{{ row.latency }}</span>
                      <span class="h-1.5 w-1.5 rounded-full bg-emerald-500"></span>
                    </div>
                  </div>

                  <div class="mt-3 space-y-1.5">
                    <div v-for="quota in previewQuotas" :key="quota.label" class="flex items-center gap-2">
                      <span class="w-12 shrink-0 text-[9px] text-gray-500">{{ quota.label }}</span>
                      <span class="h-1 flex-1 overflow-hidden rounded-full bg-white/[0.06]">
                        <span class="block h-full rounded-full bg-primary-500/70" :style="{ width: quota.value }"></span>
                      </span>
                      <span class="w-7 shrink-0 text-right font-mono text-[9px] text-gray-500">{{ quota.value }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>

          </div>
        </div>

        <!-- 事实条 -->
        <dl class="brand-glass mt-12 grid overflow-hidden sm:grid-cols-3">
          <div
            v-for="(fact, index) in facts"
            :key="fact.title"
            class="flex items-start gap-3 border-b border-gray-900/5 px-5 py-4 last:border-b-0 sm:border-b-0 sm:border-r sm:last:border-r-0 dark:border-white/5"
          >
            <span class="mt-0.5 font-mono text-[10px] text-primary-600 dark:text-primary-400">
              {{ String(index + 1).padStart(2, '0') }}
            </span>
            <div class="min-w-0">
              <dt class="text-[13px] font-semibold text-gray-900 dark:text-gray-100">{{ fact.title }}</dt>
              <dd class="mt-0.5 truncate text-[11.5px] text-gray-500 dark:text-dark-400">{{ fact.value }}</dd>
            </div>
          </div>
        </dl>
      </section>

      <!-- 真实访问延迟 -->
      <section id="latency" class="mx-auto max-w-7xl scroll-mt-24 px-4 py-8 sm:px-6">
        <div class="flex flex-wrap items-end justify-between gap-4">
          <div class="min-w-0">
            <p class="brand-eyebrow flex items-center gap-2">
              <span class="h-1.5 w-1.5 rounded-full bg-primary-500"></span>
              {{ t('home.latency.eyebrow') }}
            </p>
            <h2 class="mt-3 text-2xl font-semibold tracking-[-0.02em] sm:text-[1.9rem]">
              {{ t('home.latency.title') }}
            </h2>
            <p class="mt-2 max-w-3xl text-sm leading-relaxed text-gray-600 dark:text-dark-300">
              {{ t('home.latency.description') }}
            </p>
          </div>
        </div>

        <div class="brand-glass mt-6 overflow-hidden">
          <div class="grid grid-cols-1 lg:grid-cols-[minmax(0,1.18fr)_minmax(0,0.82fr)]">
            <div class="min-w-0 p-5 sm:p-7">
              <div class="flex flex-wrap items-start justify-between gap-6">
                <div>
                  <p class="brand-eyebrow">{{ t('home.latency.current') }}</p>
                  <div class="mt-2 flex items-baseline gap-1.5">
                    <span class="text-[3.2rem] font-semibold leading-none tracking-[-0.045em] tabular-nums text-gray-900 dark:text-white">
                      {{ currentMsText }}
                    </span>
                    <span class="text-sm text-gray-500 dark:text-dark-400">ms</span>
                  </div>
                  <span
                    class="mt-3 inline-flex items-center gap-1.5 rounded-full px-2 py-0.5 text-[11px] font-medium"
                    :class="latencyStatus.classes"
                  >
                    <span class="h-1.5 w-1.5 rounded-full" :class="latencyStatus.dot"></span>
                    {{ t(latencyStatus.key) }}
                  </span>
                </div>

                <dl class="grid grid-cols-2 gap-x-8 gap-y-3 sm:grid-cols-4 lg:grid-cols-2">
                  <div v-for="metric in latencyMetrics" :key="metric.label">
                    <dt class="text-[10px] uppercase tracking-[0.14em] text-gray-500 dark:text-dark-400">{{ metric.label }}</dt>
                    <dd class="mt-0.5 font-mono text-[15px] font-medium tabular-nums text-gray-900 dark:text-gray-100">
                      {{ metric.value }}
                      <span class="text-[10px] font-normal text-gray-400">{{ metric.unit }}</span>
                    </dd>
                  </div>
                </dl>
              </div>

              <div class="relative mt-6 min-w-0">
                <svg
                  viewBox="0 0 640 150"
                  preserveAspectRatio="none"
                  class="h-[128px] w-full text-gray-900/10 dark:text-white/10"
                  role="img"
                  :aria-label="t('home.latency.title')"
                >
                  <defs>
                    <linearGradient id="brand-spark-fill" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="0%" stop-color="#94c737" stop-opacity="0.32" />
                      <stop offset="100%" stop-color="#94c737" stop-opacity="0" />
                    </linearGradient>
                  </defs>
                  <g stroke="currentColor" stroke-dasharray="3 5">
                    <line x1="0" y1="16" x2="640" y2="16" />
                    <line x1="0" y1="62" x2="640" y2="62" />
                    <line x1="0" y1="108" x2="640" y2="108" />
                    <line x1="0" y1="148" x2="640" y2="148" stroke-dasharray="none" />
                  </g>
                  <path v-if="spark" :d="spark.area" fill="url(#brand-spark-fill)" />
                  <path
                    v-if="spark"
                    :d="spark.line"
                    fill="none"
                    stroke="#6f9e12"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    vector-effect="non-scaling-stroke"
                  />
                </svg>
                <p
                  v-if="!spark"
                  class="absolute inset-0 flex items-center justify-center font-mono text-[11px] text-gray-500 dark:text-dark-400"
                >
                  {{ t('home.latency.waiting') }}
                </p>
              </div>

              <div class="mt-2 flex items-center justify-between gap-4 font-mono text-[10px] text-gray-500 dark:text-dark-400">
                <span>{{ t('home.latency.window', { count: probeSamples.length }) }}</span>
                <span v-if="spark" class="tabular-nums">
                  min {{ formatMs(spark?.min) }} ms · max {{ formatMs(spark?.max) }} ms
                </span>
              </div>
            </div>

            <div class="min-w-0 border-t border-gray-900/10 p-5 dark:border-white/10 sm:p-7 lg:border-l lg:border-t-0">
              <p class="brand-eyebrow">{{ t('home.latency.recent') }}</p>

              <dl class="mt-4 grid grid-cols-2 gap-x-4 gap-y-3">
                <div v-for="detail in latestDetails" :key="detail.label">
                  <dt class="text-[10px] uppercase tracking-[0.14em] text-gray-500 dark:text-dark-400">{{ detail.label }}</dt>
                  <dd class="mt-0.5 truncate font-mono text-[12.5px] text-gray-900 dark:text-gray-100">{{ detail.value }}</dd>
                </div>
              </dl>

              <div class="mt-5 space-y-2">
                <div
                  v-for="sample in recentSamples"
                  :key="sample.at"
                  class="flex items-center gap-2.5 font-mono text-[10.5px]"
                >
                  <span class="w-14 shrink-0 text-gray-500 dark:text-dark-400">{{ formatClock(sample.at) }}</span>
                  <span class="h-1 flex-1 overflow-hidden rounded-full bg-gray-900/[0.06] dark:bg-white/[0.08]">
                    <span
                      class="block h-full rounded-full"
                      :class="sample.ok ? 'bg-primary-500/70' : 'bg-amber-500/70'"
                      :style="{ width: sampleBarWidth(sample) }"
                    ></span>
                  </span>
                  <span class="w-12 shrink-0 text-right tabular-nums text-gray-900 dark:text-gray-100">
                    {{ sample.ok ? formatMs(sample.ms) : '--' }}
                  </span>
                  <span class="h-1.5 w-1.5 shrink-0 rounded-full" :class="sample.ok ? 'bg-emerald-500' : 'bg-amber-500'"></span>
                </div>
                <p v-if="!probeSamples.length" class="font-mono text-[11px] text-gray-500 dark:text-dark-400">
                  {{ t('home.latency.noData') }}
                </p>
              </div>

              <div class="mt-5 flex items-center justify-between gap-3 border-t border-gray-900/10 pt-4 dark:border-white/10">
                <p class="text-[10.5px] leading-relaxed text-gray-500 dark:text-dark-400">
                  {{ probeStatus === 'error' ? t('home.latency.failed') : t('home.latency.source') }}
                </p>
                <button
                  type="button"
                  class="inline-flex shrink-0 items-center gap-1 rounded-lg border border-white/80 bg-white/55 px-2.5 py-1.5 text-[11px] font-medium text-gray-700 transition-colors hover:bg-white/85 dark:border-white/10 dark:bg-dark-900/50 dark:text-gray-200 dark:hover:bg-dark-800/80"
                  @click="runProbe"
                >
                  <Icon name="refresh" size="xs" />
                  {{ t('home.latency.retry') }}
                </button>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- 控制台能力 -->
      <section id="capabilities" class="mx-auto max-w-7xl scroll-mt-24 px-4 py-8 sm:px-6">
        <div class="flex flex-wrap items-end justify-between gap-4">
          <div class="min-w-0">
            <p class="brand-eyebrow">{{ t('home.capabilities.eyebrow') }}</p>
            <h2 class="mt-3 text-2xl font-semibold tracking-[-0.02em] sm:text-[1.9rem]">
              {{ t('home.capabilities.title') }}
            </h2>
            <p class="mt-2 max-w-2xl text-sm leading-relaxed text-gray-600 dark:text-dark-300">
              {{ t('home.capabilities.subtitle') }}
            </p>
          </div>
          <p class="brand-eyebrow shrink-0">{{ t('home.capabilities.mockNote') }}</p>
        </div>

        <div class="mt-6 grid gap-4 lg:grid-cols-12">
          <!-- 账号池 -->
          <article class="brand-glass p-5 sm:p-6 lg:col-span-7">
            <div class="flex items-start justify-between gap-4">
              <div>
                <h3 class="text-[15px] font-semibold text-gray-900 dark:text-gray-100">{{ t('home.capabilities.routing.title') }}</h3>
                <p class="mt-1.5 max-w-md text-[13px] leading-relaxed text-gray-600 dark:text-dark-300">
                  {{ t('home.capabilities.routing.desc') }}
                </p>
              </div>
              <span class="shrink-0 rounded-full bg-primary-500/10 px-2 py-0.5 font-mono text-[10px] text-primary-700 dark:bg-primary-500/15 dark:text-primary-300">
                {{ t('home.capabilities.routing.sticky') }}
              </span>
            </div>

            <div class="mt-5 overflow-hidden rounded-xl border border-gray-900/[0.07] bg-white/45 dark:border-white/[0.07] dark:bg-dark-950/35">
              <div class="flex items-center justify-between border-b border-gray-900/[0.06] px-3 py-2 dark:border-white/[0.06]">
                <span class="font-mono text-[10px] uppercase tracking-[0.14em] text-gray-500 dark:text-dark-400">{{ t('home.capabilities.routing.pool') }}</span>
                <span class="font-mono text-[10px] text-emerald-600 dark:text-emerald-400">3 / 3</span>
              </div>
              <div
                v-for="(account, index) in routingAccounts"
                :key="account.name"
                class="grid grid-cols-[8px_minmax(0,1fr)_auto_auto] items-center gap-2.5 border-b border-gray-900/[0.04] px-3 py-2.5 last:border-b-0 dark:border-white/[0.04]"
                :class="index === 0 ? 'bg-primary-500/[0.06]' : ''"
              >
                <span class="h-1.5 w-1.5 rounded-full bg-emerald-500"></span>
                <span class="truncate font-mono text-[11px] text-gray-800 dark:text-gray-200">{{ account.name }}</span>
                <span class="hidden font-mono text-[10px] text-gray-500 sm:inline dark:text-dark-400">{{ account.model }}</span>
                <span class="flex items-center gap-2">
                  <span class="h-1 w-16 overflow-hidden rounded-full bg-gray-900/[0.07] dark:bg-white/[0.08]">
                    <span class="block h-full rounded-full bg-primary-500/70" :style="{ width: account.load }"></span>
                  </span>
                  <span class="w-10 text-right font-mono text-[10px] text-gray-500 dark:text-dark-400">{{ account.latency }}</span>
                </span>
              </div>
            </div>
          </article>

          <!-- 计费 -->
          <article class="brand-glass p-5 sm:p-6 lg:col-span-5">
            <h3 class="text-[15px] font-semibold text-gray-900 dark:text-gray-100">{{ t('home.capabilities.billing.title') }}</h3>
            <p class="mt-1.5 text-[13px] leading-relaxed text-gray-600 dark:text-dark-300">{{ t('home.capabilities.billing.desc') }}</p>

            <div class="mt-5 rounded-xl border border-gray-900/[0.07] bg-white/45 p-3 dark:border-white/[0.07] dark:bg-dark-950/35">
              <div v-for="row in billingRows" :key="row.label" class="flex items-center gap-2.5 py-1.5">
                <span class="w-9 shrink-0 text-[10px] text-gray-500 dark:text-dark-400">{{ row.label }}</span>
                <span class="h-1.5 flex-1 overflow-hidden rounded-full bg-gray-900/[0.06] dark:bg-white/[0.07]">
                  <span class="block h-full rounded-full bg-gray-900/70 dark:bg-gray-200/70" :style="{ width: row.width }"></span>
                </span>
                <span class="w-12 shrink-0 text-right font-mono text-[10px] tabular-nums text-gray-700 dark:text-gray-300">{{ row.value }}</span>
              </div>
              <div class="mt-2 flex items-center justify-between border-t border-gray-900/[0.06] pt-2 dark:border-white/[0.06]">
                <span class="text-[10px] text-gray-500 dark:text-dark-400">{{ t('home.capabilities.billing.cost') }}</span>
                <span class="font-mono text-[13px] font-medium text-gray-900 dark:text-gray-100">$0.0043</span>
              </div>
            </div>
          </article>

          <!-- 配额 -->
          <article class="brand-glass p-5 sm:p-6 lg:col-span-5">
            <h3 class="text-[15px] font-semibold text-gray-900 dark:text-gray-100">{{ t('home.capabilities.quota.title') }}</h3>
            <p class="mt-1.5 text-[13px] leading-relaxed text-gray-600 dark:text-dark-300">{{ t('home.capabilities.quota.desc') }}</p>

            <div class="mt-5 space-y-3 rounded-xl border border-gray-900/[0.07] bg-white/45 p-3 dark:border-white/[0.07] dark:bg-dark-950/35">
              <div v-for="bar in quotaBars" :key="bar.label">
                <div class="flex items-center justify-between text-[10px]">
                  <span class="text-gray-500 dark:text-dark-400">{{ bar.label }}</span>
                  <span class="font-mono tabular-nums text-gray-700 dark:text-gray-300">{{ bar.value }}%</span>
                </div>
                <div class="mt-1 h-1.5 overflow-hidden rounded-full bg-gray-900/[0.06] dark:bg-white/[0.07]">
                  <span
                    class="block h-full rounded-full"
                    :class="bar.value > 80 ? 'bg-amber-500/80' : 'bg-primary-500/70'"
                    :style="{ width: bar.value + '%' }"
                  ></span>
                </div>
              </div>
            </div>
          </article>

          <!-- 审计 -->
          <article class="brand-glass p-5 sm:p-6 lg:col-span-7">
            <h3 class="text-[15px] font-semibold text-gray-900 dark:text-gray-100">{{ t('home.capabilities.audit.title') }}</h3>
            <p class="mt-1.5 max-w-xl text-[13px] leading-relaxed text-gray-600 dark:text-dark-300">{{ t('home.capabilities.audit.desc') }}</p>

            <div class="mt-5 overflow-hidden rounded-xl border border-gray-900/[0.07] bg-white/45 dark:border-white/[0.07] dark:bg-dark-950/35">
              <div class="grid grid-cols-[46px_minmax(0,1fr)_minmax(0,1fr)_46px_52px] gap-2 border-b border-gray-900/[0.06] px-3 py-2 font-mono text-[9px] uppercase tracking-[0.1em] text-gray-500 dark:border-white/[0.06] dark:text-dark-400">
                <span>{{ t('home.capabilities.audit.time') }}</span>
                <span>{{ t('home.capabilities.audit.key') }}</span>
                <span>{{ t('home.capabilities.audit.model') }}</span>
                <span class="text-right">{{ t('home.capabilities.audit.status') }}</span>
                <span class="text-right">{{ t('home.capabilities.audit.latency') }}</span>
              </div>
              <div
                v-for="log in auditRows"
                :key="log.time"
                class="grid grid-cols-[46px_minmax(0,1fr)_minmax(0,1fr)_46px_52px] items-center gap-2 border-b border-gray-900/[0.04] px-3 py-2 font-mono text-[10px] last:border-b-0 dark:border-white/[0.04]"
              >
                <span class="text-gray-500 dark:text-dark-400">{{ log.time }}</span>
                <span class="truncate text-gray-700 dark:text-gray-300">{{ log.key }}</span>
                <span class="truncate text-gray-700 dark:text-gray-300">{{ log.model }}</span>
                <span class="text-right text-emerald-600 dark:text-emerald-400">{{ log.status }}</span>
                <span class="text-right text-gray-500 dark:text-dark-400">{{ log.latency }}</span>
              </div>
            </div>
          </article>
        </div>
      </section>

      <!-- 快速接入 -->
      <section id="integration" class="mx-auto max-w-7xl scroll-mt-24 px-4 py-8 sm:px-6">
        <div class="brand-glass overflow-hidden">
          <div class="grid grid-cols-1 lg:grid-cols-[minmax(0,0.92fr)_minmax(0,1.08fr)]">
            <div class="min-w-0 p-5 sm:p-7">
              <p class="brand-eyebrow">{{ t('home.integration.eyebrow') }}</p>
              <h2 class="mt-3 text-2xl font-semibold tracking-[-0.02em] sm:text-[1.9rem]">
                {{ t('home.integration.title') }}
              </h2>
              <p class="mt-2 max-w-lg text-sm leading-relaxed text-gray-600 dark:text-dark-300">
                {{ t('home.integration.description') }}
              </p>

              <ul class="mt-6 space-y-3">
                <li v-for="point in integrationPoints" :key="point" class="flex items-start gap-2.5">
                  <Icon name="check" size="xs" class="mt-1 shrink-0 text-primary-600 dark:text-primary-400" />
                  <span class="text-[13px] leading-relaxed text-gray-700 dark:text-gray-300">{{ point }}</span>
                </li>
              </ul>

              <div class="mt-6 flex flex-wrap items-center gap-2">
                <span
                  v-for="endpoint in endpoints"
                  :key="endpoint"
                  class="rounded-lg border border-gray-900/10 bg-white/45 px-2.5 py-1 font-mono text-[10.5px] text-gray-700 dark:border-white/10 dark:bg-dark-950/35 dark:text-gray-300"
                >
                  {{ endpoint }}
                </span>
              </div>
            </div>

            <div class="min-w-0 border-t border-gray-900/10 p-5 dark:border-white/10 sm:p-7 lg:border-l lg:border-t-0">
              <div class="overflow-hidden rounded-xl bg-gray-900 ring-1 ring-black/40 dark:bg-dark-900">
                <div class="flex items-center justify-between border-b border-white/10 px-3 py-2">
                  <span class="font-mono text-[10px] uppercase tracking-[0.14em] text-gray-500">curl</span>
                  <button
                    type="button"
                    class="inline-flex items-center gap-1 rounded-md border border-white/10 bg-white/5 px-2 py-1 font-mono text-[10px] text-gray-300 transition-colors hover:bg-white/10"
                    @click="copySnippet"
                  >
                    <Icon :name="codeCopied ? 'check' : 'copy'" size="xs" />
                    {{
                      codeCopied
                        ? t('home.integration.copied')
                        : codeCopyFailed
                          ? t('home.integration.copyFailed')
                          : t('home.integration.copy')
                    }}
                  </button>
                </div>
                <pre class="overflow-x-auto px-3 py-3 font-mono text-[11px] leading-relaxed text-gray-300"><code>{{ snippet }}</code></pre>
              </div>
              <p class="mt-3 text-[11px] text-gray-500 dark:text-dark-400">{{ t('home.integration.compatible') }}</p>
            </div>
          </div>
        </div>
      </section>

      <!-- 模型 -->
      <section id="models" class="mx-auto max-w-7xl scroll-mt-24 px-4 py-8 sm:px-6">
        <div class="flex flex-wrap items-end justify-between gap-4">
          <div class="min-w-0">
            <h2 class="text-2xl font-semibold tracking-[-0.02em] sm:text-[1.9rem]">{{ t('home.models.title') }}</h2>
            <p class="mt-2 max-w-2xl text-sm leading-relaxed text-gray-600 dark:text-dark-300">{{ t('home.models.description') }}</p>
          </div>
          <router-link
            v-if="showModelPlazaEntry"
            to="/model-plaza"
            class="inline-flex h-9 shrink-0 items-center gap-1.5 rounded-lg border border-gray-900/10 bg-white/45 px-3 text-[12.5px] font-medium text-gray-700 backdrop-blur-md transition-colors hover:bg-white/75 dark:border-white/10 dark:bg-dark-900/45 dark:text-gray-200 dark:hover:bg-dark-800/70"
          >
            {{ t('nav.modelPlaza') }}
            <Icon name="arrowRight" size="xs" />
          </router-link>
        </div>

        <ul class="mt-5 flex flex-wrap gap-2.5">
          <li
            v-for="provider in providers"
            :key="provider.name"
            class="inline-flex items-center gap-2.5 rounded-xl border border-gray-900/10 bg-white/45 px-3.5 py-2.5 backdrop-blur-md dark:border-white/10 dark:bg-dark-900/40"
          >
            <span class="font-mono text-[10px] font-medium text-primary-700 dark:text-primary-400">{{ provider.mark }}</span>
            <span class="text-[13px] font-medium text-gray-800 dark:text-gray-200">{{ provider.name }}</span>
            <span
              class="rounded-md px-1.5 py-0.5 text-[10px]"
              :class="provider.available
                ? 'bg-primary-500/10 text-primary-700 dark:bg-primary-500/15 dark:text-primary-300'
                : 'bg-gray-900/[0.06] text-gray-500 dark:bg-white/[0.06] dark:text-dark-400'"
            >
              {{ provider.available ? t('home.providers.supported') : t('home.providers.soon') }}
            </span>
          </li>
        </ul>
      </section>

      <!-- CTA -->
      <section class="mx-auto max-w-7xl px-4 py-8 sm:px-6">
        <div class="brand-glass relative overflow-hidden px-6 py-8 sm:px-9 sm:py-10">
          <div
            class="pointer-events-none absolute inset-0 bg-[radial-gradient(110%_140%_at_100%_0%,rgba(111,158,18,0.16),transparent_58%)]"
          ></div>
          <div class="relative flex flex-col items-start justify-between gap-6 lg:flex-row lg:items-center">
            <div class="max-w-xl">
              <h2 class="text-2xl font-semibold tracking-[-0.02em] sm:text-[1.9rem]">{{ t('home.cta.title') }}</h2>
              <p class="mt-2 text-sm leading-relaxed text-gray-600 dark:text-dark-300">{{ t('home.cta.description') }}</p>
            </div>
            <div class="flex shrink-0 flex-wrap items-center gap-2.5">
              <router-link
                :to="isAuthenticated ? dashboardPath : '/login'"
                class="inline-flex h-11 items-center gap-2 rounded-xl bg-gray-900 px-5 text-sm font-medium text-white transition-colors hover:bg-gray-800 dark:bg-gray-100 dark:text-gray-900 dark:hover:bg-white"
              >
                {{ isAuthenticated ? t('home.hero.primaryCtaAuthed') : t('home.cta.primaryGuest') }}
                <Icon name="arrowRight" size="sm" />
              </router-link>
              <router-link
                to="/key-usage"
                class="inline-flex h-11 items-center gap-2 rounded-xl border border-white/80 bg-white/55 px-4 text-sm font-medium text-gray-800 backdrop-blur-md transition-colors hover:bg-white/85 dark:border-white/10 dark:bg-dark-900/50 dark:text-gray-200 dark:hover:bg-dark-800/80"
              >
                {{ t('home.cta.secondary') }}
              </router-link>
            </div>
          </div>
        </div>
      </section>
    </main>

    <!-- Footer -->
    <footer class="relative z-10 border-t border-gray-900/10 dark:border-white/10">
      <div class="mx-auto flex max-w-7xl flex-col gap-4 px-4 py-6 sm:flex-row sm:items-center sm:justify-between sm:px-6">
        <div class="flex min-w-0 items-center gap-2.5">
          <img :src="siteLogo || '/logo.png'" alt="Logo" class="h-6 w-6 shrink-0 object-contain" />
          <p class="min-w-0 truncate text-[12px] text-gray-500 dark:text-dark-400">
            &copy; {{ currentYear }} {{ siteName }}. {{ t('home.footer.allRightsReserved') }}
            <span v-if="siteVersion" class="ml-1 font-mono text-[10px] text-gray-400 dark:text-dark-500">v{{ siteVersion }}</span>
          </p>
        </div>

        <div class="flex flex-wrap items-center gap-4 text-[12px]">
          <router-link to="/key-usage" class="text-gray-500 transition-colors hover:text-gray-900 dark:text-dark-400 dark:hover:text-white">
            {{ t('home.footer.keyUsage') }}
          </router-link>
          <router-link
            v-if="showModelPlazaEntry"
            to="/model-plaza"
            class="text-gray-500 transition-colors hover:text-gray-900 dark:text-dark-400 dark:hover:text-white"
          >
            {{ t('nav.modelPlaza') }}
          </router-link>
          <router-link
            :to="isAuthenticated ? dashboardPath : '/login'"
            class="text-gray-500 transition-colors hover:text-gray-900 dark:text-dark-400 dark:hover:text-white"
          >
            {{ isAuthenticated ? t('home.dashboard') : t('home.footer.login') }}
          </router-link>
          <span class="inline-flex items-center gap-1.5 font-mono text-[10px] uppercase tracking-[0.12em] text-gray-500 dark:text-dark-400">
            <span class="h-1.5 w-1.5 rounded-full bg-emerald-500"></span>
            {{ t('home.footer.status') }}
          </span>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore, useAppStore } from '@/stores'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import Icon from '@/components/icons/Icon.vue'
import { sanitizeUrl } from '@/utils/url'
import { FeatureFlags, isFeatureFlagEnabled } from '@/utils/featureFlags'

const { t } = useI18n()

const authStore = useAuthStore()
const appStore = useAppStore()

// Site settings - directly from appStore (already initialized from injected config)
const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'CodeAPI')
const siteLogo = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '', { allowRelative: true, allowDataUrl: true }))
const siteSubtitle = computed(() => appStore.cachedPublicSettings?.site_subtitle || 'AI API Gateway Platform')
const siteVersion = computed(() => appStore.siteVersion || '')
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')
const hasHomeContent = computed(() => homeContent.value.trim().length > 0)
const compactHomeEnabled = computed(() => appStore.cachedPublicSettings?.compact_home_enabled === true)
const modelPlazaEnabled = computed(() => isFeatureFlagEnabled(FeatureFlags.modelPlaza))

// Check if homeContent is a URL (for iframe display)
const isHomeContentUrl = computed(() => {
  const content = homeContent.value.trim()
  return content.startsWith('http://') || content.startsWith('https://')
})

// Theme
const isDark = ref(document.documentElement.classList.contains('dark'))

// Auth state
const isAuthenticated = computed(() => authStore.isAuthenticated)
const modelPlazaRequiresAuth = computed(
  () => appStore.cachedPublicSettings?.model_plaza_require_auth === true,
)
const showModelPlazaEntry = computed(
  () => modelPlazaEnabled.value && (isAuthenticated.value || !modelPlazaRequiresAuth.value),
)
const isAdmin = computed(() => authStore.isAdmin)
const dashboardPath = computed(() => isAdmin.value ? '/admin/dashboard' : '/dashboard')

const currentYear = computed(() => new Date().getFullYear())

const navItems = computed(() => [
  { href: '#capabilities', label: t('home.nav.capabilities') },
  { href: '#latency', label: t('home.nav.latency') },
  { href: '#integration', label: t('home.nav.integration') },
  { href: '#models', label: t('home.nav.models') }
])

const facts = computed(() => [
  { title: t('home.facts.models.title'), value: t('home.facts.models.value') },
  { title: t('home.facts.sticky.title'), value: t('home.facts.sticky.value') },
  { title: t('home.facts.billing.title'), value: t('home.facts.billing.value') }
])

// ==================== Console preview (illustrative only) ====================

const previewTabs = ['grid', 'chart', 'key', 'chartBar', 'user'] as const

const previewDate = new Intl.DateTimeFormat(undefined, {
  month: '2-digit',
  day: '2-digit',
  hour: '2-digit',
  minute: '2-digit',
  hour12: false
}).format(new Date())

const previewStats = computed(() => [
  { label: t('home.preview.balance'), value: '$12.40', accent: false },
  { label: t('home.preview.todayRequests'), value: '1,284', accent: false },
  { label: t('home.preview.todayCost'), value: '$3.86', accent: false },
  { label: t('home.preview.rpm'), value: '42', accent: true }
])

const previewRows = [
  { time: '14:02', model: 'claude-sonnet', tokens: '1.2k / 380', latency: '412 ms' },
  { time: '14:01', model: 'gpt', tokens: '860 / 210', latency: '268 ms' },
  { time: '14:01', model: 'claude-opus', tokens: '2.4k / 720', latency: '588 ms' },
  { time: '13:58', model: 'gemini', tokens: '640 / 180', latency: '231 ms' }
]

const previewQuotas = computed(() => [
  { label: t('home.preview.fiveHour'), value: '42%' },
  { label: t('home.preview.daily'), value: '18%' }
])

// ==================== Live latency probe ====================
// 真实测量：浏览器向本站静态资源发请求，记录往返时延与性能条目，不做任何模拟。

type ProbeSample = {
  at: number
  ms: number
  ok: boolean
  ttfb: number | null
  reused: boolean | null
  bytes: number | null
  protocol: string | null
}

const PROBE_PATH = '/logo.png'
const PROBE_INTERVAL_MS = 5000
const probeSamples = ref<ProbeSample[]>([])
const probeStatus = ref<'idle' | 'ok' | 'error'>('idle')
let probeTimer: number | undefined
let warmupTimer: number | undefined
let probeInFlight = false

function findTimingEntry(url: string): PerformanceResourceTiming | null {
  if (typeof performance.getEntriesByName !== 'function') return null
  const entries = performance.getEntriesByName(url) as PerformanceResourceTiming[]
  const last = entries[entries.length - 1]
  return last && last.responseStart > 0 ? last : null
}

async function runProbe(): Promise<void> {
  if (probeInFlight) return
  if (typeof document !== 'undefined' && document.hidden) return
  probeInFlight = true

  const startedAt = Date.now()
  const started = performance.now()
  const url = `${PROBE_PATH}?probe=${startedAt}`

  try {
    const response = await fetch(url, {
      cache: 'no-store',
      credentials: 'omit',
      headers: { Range: 'bytes=0-1023' }
    })
    if (!response.ok) throw new Error(`HTTP ${response.status}`)
    await response.arrayBuffer()

    const ms = performance.now() - started
    const entry = findTimingEntry(response.url)
    const transferSize = entry && entry.transferSize > 0 ? entry.transferSize : null

    pushSample({
      at: startedAt,
      ms,
      ok: true,
      ttfb: entry && entry.responseStart > entry.startTime
        ? entry.responseStart - entry.startTime
        : null,
      reused: entry ? entry.connectStart === 0 : null,
      bytes: transferSize,
      protocol: entry?.nextHopProtocol || null
    })
    probeStatus.value = 'ok'
  } catch {
    pushSample({ at: startedAt, ms: 0, ok: false, ttfb: null, reused: null, bytes: null, protocol: null })
    probeStatus.value = 'error'
  } finally {
    probeInFlight = false
  }
}

function pushSample(sample: ProbeSample): void {
  probeSamples.value = [...probeSamples.value, sample].slice(-40)
}

function startProbe(): void {
  void runProbe()
  // 首个样本后补一次，让折线不用等满一个采样周期才出现
  warmupTimer = window.setTimeout(() => void runProbe(), 900)
  probeTimer = window.setInterval(() => {
    void runProbe()
  }, PROBE_INTERVAL_MS)
}

function handleVisibilityChange(): void {
  if (!document.hidden) void runProbe()
}

const okSamples = computed(() => probeSamples.value.filter((sample) => sample.ok))
const latestSample = computed(() => okSamples.value[okSamples.value.length - 1] || null)
const windowSamples = computed(() => okSamples.value.slice(-12))

const currentMsText = computed(() => (latestSample.value ? formatMs(latestSample.value.ms) : '--'))

const heroProbeText = computed(() => {
  if (latestSample.value) return t('home.hero.probeValue', { ms: formatMs(latestSample.value.ms) })
  return t('home.hero.probePending')
})

const latencyMetrics = computed(() => {
  const values = windowSamples.value.map((sample) => sample.ms)
  const avg = values.length ? values.reduce((sum, value) => sum + value, 0) / values.length : null
  const sorted = [...values].sort((a, b) => a - b)
  const p95 = sorted.length ? sorted[Math.min(sorted.length - 1, Math.ceil(sorted.length * 0.95) - 1)] : null
  let jitter: number | null = null
  if (values.length > 1) {
    let total = 0
    for (let i = 1; i < values.length; i += 1) total += Math.abs(values[i] - values[i - 1])
    jitter = total / (values.length - 1)
  }

  return [
    { label: t('home.latency.avg'), value: formatMs(avg), unit: 'ms' },
    { label: t('home.latency.p95'), value: formatMs(p95), unit: 'ms' },
    { label: t('home.latency.jitter'), value: formatMs(jitter), unit: 'ms' },
    { label: t('home.latency.sampleCount'), value: String(probeSamples.value.length), unit: '' }
  ]
})

const latencyStatus = computed(() => {
  if (probeStatus.value === 'error') {
    return {
      key: 'home.latency.failed',
      classes: 'bg-amber-500/10 text-amber-700 dark:bg-amber-500/15 dark:text-amber-300',
      dot: 'bg-amber-500'
    }
  }
  if (!latestSample.value) {
    return {
      key: 'home.latency.measuring',
      classes: 'bg-gray-900/[0.06] text-gray-600 dark:bg-white/[0.06] dark:text-dark-300',
      dot: 'bg-gray-400'
    }
  }
  if (latestSample.value.ms <= 120) {
    return {
      key: 'home.latency.fast',
      classes: 'bg-emerald-500/10 text-emerald-700 dark:bg-emerald-500/15 dark:text-emerald-300',
      dot: 'bg-emerald-500'
    }
  }
  if (latestSample.value.ms <= 300) {
    return {
      key: 'home.latency.normal',
      classes: 'bg-primary-500/10 text-primary-700 dark:bg-primary-500/15 dark:text-primary-300',
      dot: 'bg-primary-500'
    }
  }
  return {
    key: 'home.latency.slow',
    classes: 'bg-amber-500/10 text-amber-700 dark:bg-amber-500/15 dark:text-amber-300',
    dot: 'bg-amber-500'
  }
})

const latestDetails = computed(() => {
  const sample = latestSample.value
  return [
    {
      label: t('home.latency.ttfb'),
      value: sample?.ttfb != null ? `${formatMs(sample.ttfb)} ms` : '--'
    },
    {
      label: t('home.latency.connection'),
      value: sample?.reused == null
        ? '--'
        : sample.reused
          ? t('home.latency.reused')
          : t('home.latency.fresh')
    },
    {
      label: t('home.latency.transfer'),
      value: formatBytes(sample?.bytes ?? null)
    },
    {
      label: t('home.latency.protocol'),
      value: sample?.protocol || '--'
    }
  ]
})

const recentSamples = computed(() => probeSamples.value.slice(-5).reverse())

const spark = computed(() => {
  const points = probeSamples.value.slice(-32).filter((sample) => sample.ok)
  if (points.length < 2) return null

  const width = 640
  const height = 150
  const top = 14
  const bottom = 138
  const values = points.map((point) => point.ms)
  const min = Math.min(...values)
  const max = Math.max(...values)
  // 从 0 起算，避免样本少时把几毫秒的波动放大成陡峭趋势
  const domainMax = Math.max(50, Math.ceil(max / 25) * 25)
  const step = width / (points.length - 1)

  const coords = values.map((value, index) => {
    const x = index * step
    const y = bottom - (value / domainMax) * (bottom - top)
    return `${x.toFixed(1)},${y.toFixed(1)}`
  })

  const line = coords.map((coord, index) => `${index === 0 ? 'M' : 'L'}${coord}`).join(' ')
  const area = `${line} L${width},${height} L0,${height} Z`

  return { line, area, min, max }
})

const sampleMax = computed(() => Math.max(1, ...okSamples.value.map((sample) => sample.ms)))

function sampleBarWidth(sample: ProbeSample): string {
  if (!sample.ok) return '6%'
  return `${Math.max(6, Math.round((sample.ms / sampleMax.value) * 100))}%`
}

function formatMs(value: number | null | undefined): string {
  if (value == null) return '--'
  return value < 10 ? value.toFixed(1) : String(Math.round(value))
}

function formatBytes(value: number | null): string {
  if (value == null || value <= 0) return '--'
  if (value < 1024) return `${value} B`
  return `${(value / 1024).toFixed(1)} KB`
}

function formatClock(at: number): string {
  return new Intl.DateTimeFormat(undefined, {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false
  }).format(new Date(at))
}

// ==================== Integration snippet ====================

const siteOrigin = typeof window === 'undefined' ? '' : window.location.origin

const snippet = computed(() => [
  `curl "${siteOrigin}/v1/messages" \\`,
  '  -H "x-api-key: $CODEAPI_KEY" \\',
  '  -H "anthropic-version: 2023-06-01" \\',
  `  -d '{`,
  '    "model": "claude-sonnet",',
  '    "max_tokens": 512,',
  '    "messages": [{"role": "user", "content": "hello"}]',
  `  }'`
].join('\n'))

const endpoints = computed(() => [
  `POST ${siteOrigin}/v1/messages`,
  `POST ${siteOrigin}/v1/chat/completions`
])

const integrationPoints = computed(() => [
  t('home.integration.point.baseUrl'),
  t('home.integration.point.key'),
  t('home.integration.point.done')
])

const codeCopied = ref(false)
const codeCopyFailed = ref(false)
let copyResetTimer: number | undefined

async function copySnippet(): Promise<void> {
  let copied = false

  try {
    await navigator.clipboard.writeText(snippet.value)
    copied = true
  } catch {
    copied = false
  }

  if (!copied) {
    const textarea = document.createElement('textarea')
    try {
      textarea.value = snippet.value
      textarea.setAttribute('readonly', '')
      textarea.style.position = 'fixed'
      textarea.style.opacity = '0'
      document.body.appendChild(textarea)
      textarea.select()
      copied = document.execCommand('copy')
    } catch {
      copied = false
    } finally {
      textarea.remove()
    }
  }

  codeCopied.value = copied
  codeCopyFailed.value = !copied
  window.clearTimeout(copyResetTimer)
  copyResetTimer = window.setTimeout(() => {
    codeCopied.value = false
    codeCopyFailed.value = false
  }, 1600)
}

// ==================== Capability illustrations ====================

const routingAccounts = [
  { name: 'claude-pro-01', model: 'sonnet', load: '38%', latency: '42 ms' },
  { name: 'claude-pro-02', model: 'opus', load: '61%', latency: '51 ms' },
  { name: 'openai-pool-01', model: 'gpt', load: '24%', latency: '63 ms' }
]

const billingRows = computed(() => [
  { label: t('home.capabilities.billing.input'), value: '1,204', width: '82%' },
  { label: t('home.capabilities.billing.output'), value: '380', width: '30%' },
  { label: t('home.capabilities.billing.cache'), value: '896', width: '58%' }
])

const quotaBars = computed(() => [
  { label: t('home.capabilities.quota.fiveHour'), value: 42 },
  { label: t('home.capabilities.quota.daily'), value: 18 },
  { label: t('home.capabilities.quota.weekly'), value: 9 }
])

const auditRows = [
  { time: '14:02', key: 'key-01', model: 'claude-sonnet', status: '200', latency: '412 ms' },
  { time: '14:01', key: 'key-01', model: 'gpt', status: '200', latency: '268 ms' },
  { time: '14:01', key: 'key-02', model: 'claude-opus', status: '200', latency: '588 ms' },
  { time: '13:58', key: 'key-02', model: 'gemini', status: '200', latency: '231 ms' }
]

const providers = computed(() => [
  { name: t('home.providers.claude'), mark: 'CL', available: true },
  { name: 'GPT', mark: 'GPT', available: true },
  { name: t('home.providers.gemini'), mark: 'GM', available: true },
  { name: t('home.providers.antigravity'), mark: 'AG', available: true },
  { name: t('home.providers.more'), mark: '+', available: false }
])

// ==================== Theme ====================

function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

function initTheme() {
  const savedTheme = localStorage.getItem('theme')
  if (
    savedTheme === 'dark' ||
    (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)
  ) {
    isDark.value = true
    document.documentElement.classList.add('dark')
  }
}

onMounted(() => {
  initTheme()

  authStore.checkAuth()

  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }

  if (!hasHomeContent.value && !compactHomeEnabled.value) {
    startProbe()
    document.addEventListener('visibilitychange', handleVisibilityChange)
  }
})

onBeforeUnmount(() => {
  if (probeTimer !== undefined) window.clearInterval(probeTimer)
  if (warmupTimer !== undefined) window.clearTimeout(warmupTimer)
  document.removeEventListener('visibilitychange', handleVisibilityChange)
  window.clearTimeout(copyResetTimer)
})
</script>
