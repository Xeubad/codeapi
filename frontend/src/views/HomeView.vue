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
          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg text-gray-500 hover:bg-gray-100 dark:text-dark-400 dark:hover:bg-dark-800"
            :title="t('home.viewDocs')"
          >
            <Icon name="book" size="md" />
          </a>
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
    class="relative flex min-h-screen flex-col overflow-hidden bg-gray-50 text-gray-900 dark:bg-dark-950 dark:text-gray-100"
  >
    <!-- Backdrop: 柔和色场 + 颗粒纹理 -->
    <div class="pointer-events-none fixed inset-0 z-0 overflow-hidden">
      <div
        class="absolute inset-0 bg-[radial-gradient(90%_60%_at_12%_-10%,rgba(111,158,18,0.22),transparent_62%)] dark:bg-[radial-gradient(90%_60%_at_12%_-10%,rgba(111,158,18,0.22),transparent_62%)]"
      ></div>
      <div
        class="absolute inset-0 bg-[radial-gradient(70%_50%_at_100%_10%,rgba(21,94,117,0.14),transparent_60%)] dark:bg-[radial-gradient(70%_50%_at_100%_10%,rgba(56,132,160,0.16),transparent_60%)]"
      ></div>
      <div
        class="absolute inset-0 bg-[radial-gradient(60%_45%_at_50%_115%,rgba(244,232,190,0.55),transparent_60%)] dark:bg-[radial-gradient(60%_45%_at_50%_115%,rgba(38,54,32,0.5),transparent_60%)]"
      ></div>
      <div class="brand-grain absolute inset-0"></div>
    </div>

    <!-- Top Bar -->
    <header
      class="sticky top-0 z-40 border-b border-gray-200/70 bg-gray-50/80 backdrop-blur-xl dark:border-dark-800/70 dark:bg-dark-950/80"
    >
      <nav class="mx-auto flex h-16 max-w-6xl items-center justify-between gap-4 px-6">
        <router-link to="/" class="flex min-w-0 items-center gap-2.5">
          <img
            :src="siteLogo || '/logo.png'"
            alt="Logo"
            class="h-8 w-8 shrink-0 object-contain"
          />
          <span class="truncate text-[15px] font-semibold tracking-tight">{{ siteName }}</span>
        </router-link>

        <div class="flex shrink-0 items-center gap-1.5">
          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="hidden h-9 items-center gap-1.5 rounded-lg px-3 text-sm font-medium text-gray-600 transition-colors hover:bg-gray-900/5 hover:text-gray-900 sm:flex dark:text-dark-300 dark:hover:bg-white/5 dark:hover:text-white"
          >
            <Icon name="book" size="sm" />
            {{ t('home.docs') }}
          </a>
          <router-link
            v-if="showModelPlazaEntry"
            to="/model-plaza"
            class="hidden h-9 items-center gap-1.5 rounded-lg px-3 text-sm font-medium text-gray-600 transition-colors hover:bg-gray-900/5 hover:text-gray-900 sm:flex dark:text-dark-300 dark:hover:bg-white/5 dark:hover:text-white"
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
            class="ml-1 inline-flex h-9 items-center rounded-lg bg-gray-900 px-4 text-sm font-medium text-white transition-colors hover:bg-gray-800 dark:bg-gray-100 dark:text-gray-900 dark:hover:bg-white"
          >
            {{ isAuthenticated ? t('home.dashboard') : t('home.login') }}
          </router-link>
        </div>
      </nav>
    </header>

    <main class="relative z-10 flex-1">
      <!-- Hero -->
      <section class="mx-auto max-w-6xl px-6 pt-14 pb-16 lg:pt-20 lg:pb-24">
        <div
          class="grid items-center gap-14 lg:grid-cols-[minmax(0,1fr)_minmax(0,0.82fr)] lg:gap-20"
        >
          <!-- Copy -->
          <div>
            <p class="flex items-center gap-2 text-sm font-medium text-primary-700 dark:text-primary-400">
              <span class="h-1.5 w-1.5 rounded-full bg-primary-500"></span>
              {{ siteSubtitle }}
            </p>
            <h1
              class="mt-5 text-[2.1rem] font-semibold leading-[1.15] tracking-[-0.02em] text-gray-900 dark:text-white sm:text-[2.75rem]"
            >
              <span class="block">{{ t('home.heroTitleLead') }}</span>
              <span class="block">{{ t('home.heroTitleMain') }}</span>
            </h1>
            <p class="mt-5 max-w-xl text-base leading-relaxed text-gray-600 dark:text-dark-300">
              {{ t('home.heroDescription') }}
            </p>

            <div class="mt-9 flex flex-wrap items-center gap-3">
              <router-link
                :to="isAuthenticated ? dashboardPath : '/login'"
                class="inline-flex h-11 items-center gap-2 rounded-xl bg-gray-900 px-6 text-sm font-medium text-white transition-colors hover:bg-gray-800 dark:bg-gray-100 dark:text-gray-900 dark:hover:bg-white"
              >
                {{ isAuthenticated ? t('home.goToDashboard') : t('home.getStarted') }}
                <Icon name="arrowRight" size="sm" />
              </router-link>
              <a
                v-if="docUrl"
                :href="docUrl"
                target="_blank"
                rel="noopener noreferrer"
                class="inline-flex h-11 items-center gap-2 rounded-xl border border-gray-300 bg-white/60 px-5 text-sm font-medium text-gray-800 backdrop-blur transition-colors hover:border-gray-400 hover:bg-white dark:border-dark-700 dark:bg-dark-900/50 dark:text-gray-200 dark:hover:border-dark-600 dark:hover:bg-dark-800"
              >
                <Icon name="book" size="sm" />
                {{ t('home.viewDocs') }}
              </a>
            </div>

            <!-- 能力标签 -->
            <dl class="mt-12 grid max-w-xl grid-cols-1 gap-x-8 gap-y-4 sm:grid-cols-3">
              <div v-for="tag in tags" :key="tag.label" class="border-l-2 border-primary-500/40 pl-3">
                <dt class="text-sm font-medium text-gray-900 dark:text-gray-100">{{ tag.label }}</dt>
                <dd class="mt-0.5 text-xs text-gray-500 dark:text-dark-400">{{ tag.note }}</dd>
              </div>
            </dl>
          </div>

          <!-- 控制台预览 -->
          <div class="relative">
            <div
              class="overflow-hidden rounded-2xl bg-gray-900 shadow-[0_28px_70px_-28px_rgba(11,14,20,0.6)] ring-1 ring-gray-950/60"
            >
              <div class="flex items-center gap-2 border-b border-white/5 px-4 py-3">
                <span class="h-2 w-2 rounded-full bg-white/15"></span>
                <span class="h-2 w-2 rounded-full bg-white/15"></span>
                <span class="h-2 w-2 rounded-full bg-white/15"></span>
                <span class="ml-2 font-mono text-[11px] uppercase tracking-[0.2em] text-gray-500">
                  codeapi
                </span>
                <span
                  class="ml-auto rounded-md bg-primary-500/15 px-2 py-0.5 font-mono text-[11px] font-medium text-primary-300"
                >
                  200 OK
                </span>
              </div>

              <dl class="divide-y divide-white/5 font-mono text-[12.5px]">
                <div class="flex items-baseline justify-between gap-4 px-5 py-3">
                  <dt class="text-gray-500">endpoint</dt>
                  <dd class="text-gray-100">POST /v1/messages</dd>
                </div>
                <div class="flex items-baseline justify-between gap-4 px-5 py-3">
                  <dt class="text-gray-500">model</dt>
                  <dd class="text-gray-100">claude-sonnet</dd>
                </div>
                <div class="flex items-baseline justify-between gap-4 px-5 py-3">
                  <dt class="text-gray-500">upstream</dt>
                  <dd class="text-gray-100">claude-3 · sticky</dd>
                </div>
                <div class="flex items-baseline justify-between gap-4 px-5 py-3">
                  <dt class="text-gray-500">latency</dt>
                  <dd class="text-primary-300">412 ms</dd>
                </div>
                <div class="flex items-baseline justify-between gap-4 px-5 py-3">
                  <dt class="text-gray-500">tokens</dt>
                  <dd class="text-gray-100">1,204 in · 380 out</dd>
                </div>
                <div class="flex items-baseline justify-between gap-4 px-5 py-3">
                  <dt class="text-gray-500">cost</dt>
                  <dd class="text-gray-100">$0.0043</dd>
                </div>
              </dl>

              <div class="flex items-center justify-between gap-4 border-t border-white/5 px-5 py-3">
                <span class="font-mono text-[11px] uppercase tracking-[0.18em] text-gray-500">
                  {{ t('home.console.billedBy') }}
                </span>
                <span class="font-mono text-[12.5px] text-primary-300">balance $12.40</span>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- 痛点 -->
      <section class="border-t border-gray-200/70 dark:border-dark-800/70">
        <div class="mx-auto max-w-6xl px-6 py-16 lg:py-20">
          <h2 class="text-sm font-semibold uppercase tracking-[0.16em] text-gray-500 dark:text-dark-400">
            {{ t('home.painPoints.title') }}
          </h2>
          <div class="mt-8 grid gap-x-10 gap-y-8 sm:grid-cols-2 lg:grid-cols-4">
            <div v-for="point in painPoints" :key="point.title" class="border-t border-gray-200 pt-4 dark:border-dark-800">
              <h3 class="text-[15px] font-semibold text-gray-900 dark:text-gray-100">{{ point.title }}</h3>
              <p class="mt-2 text-sm leading-relaxed text-gray-600 dark:text-dark-400">{{ point.desc }}</p>
            </div>
          </div>
        </div>
      </section>

      <!-- 核心能力 -->
      <section class="border-t border-gray-200/70 dark:border-dark-800/70">
        <div class="mx-auto max-w-6xl px-6 py-16 lg:py-20">
          <div class="max-w-2xl">
            <h2 class="text-2xl font-semibold tracking-tight text-gray-900 dark:text-white sm:text-3xl">
              {{ t('home.solutions.title') }}
            </h2>
            <p class="mt-3 text-sm text-gray-600 dark:text-dark-400">
              {{ t('home.solutions.subtitle') }}
            </p>
          </div>

          <div class="mt-12 grid gap-x-10 gap-y-10 md:grid-cols-3">
            <div v-for="(item, index) in capabilities" :key="item.title" class="border-t border-gray-200 pt-5 dark:border-dark-800">
              <span class="font-mono text-xs text-primary-600 dark:text-primary-400">
                {{ String(index + 1).padStart(2, '0') }}
              </span>
              <h3 class="mt-3 text-lg font-semibold text-gray-900 dark:text-white">{{ item.title }}</h3>
              <p class="mt-2 text-sm leading-relaxed text-gray-600 dark:text-dark-400">{{ item.desc }}</p>
            </div>
          </div>
        </div>
      </section>

      <!-- 对比 -->
      <section class="border-t border-gray-200/70 dark:border-dark-800/70">
        <div class="mx-auto max-w-6xl px-6 py-16 lg:py-20">
          <h2 class="text-2xl font-semibold tracking-tight text-gray-900 dark:text-white sm:text-3xl">
            {{ t('home.comparison.title') }}
          </h2>

          <div class="mt-10 overflow-hidden rounded-2xl border border-gray-200/80 dark:border-dark-800">
            <div
              class="grid grid-cols-[minmax(0,0.8fr)_minmax(0,1fr)_minmax(0,1fr)] gap-4 bg-gray-100/70 px-5 py-3 text-xs font-semibold uppercase tracking-[0.12em] text-gray-500 dark:bg-dark-900/60 dark:text-dark-400 sm:px-6"
            >
              <span>{{ t('home.comparison.headers.feature') }}</span>
              <span>{{ t('home.comparison.headers.official') }}</span>
              <span class="text-primary-700 dark:text-primary-400">{{ t('home.comparison.headers.us') }}</span>
            </div>
            <div
              v-for="row in comparisonRows"
              :key="row.feature"
              class="grid grid-cols-[minmax(0,0.8fr)_minmax(0,1fr)_minmax(0,1fr)] gap-4 border-t border-gray-200/70 px-5 py-4 text-sm dark:border-dark-800 sm:px-6"
            >
              <span class="text-gray-500 dark:text-dark-400">{{ row.feature }}</span>
              <span class="text-gray-600 dark:text-dark-300">{{ row.official }}</span>
              <span class="font-medium text-gray-900 dark:text-gray-100">{{ row.us }}</span>
            </div>
          </div>
        </div>
      </section>

      <!-- 已支持的模型 -->
      <section class="border-t border-gray-200/70 dark:border-dark-800/70">
        <div class="mx-auto max-w-6xl px-6 py-16 lg:py-20">
          <div class="flex flex-wrap items-end justify-between gap-4">
            <div>
              <h2 class="text-2xl font-semibold tracking-tight text-gray-900 dark:text-white sm:text-3xl">
                {{ t('home.providers.title') }}
              </h2>
              <p class="mt-3 text-sm text-gray-600 dark:text-dark-400">
                {{ t('home.providers.description') }}
              </p>
            </div>
          </div>

          <ul class="mt-10 flex flex-wrap gap-3">
            <li
              v-for="provider in providers"
              :key="provider.name"
              class="inline-flex items-center gap-2 rounded-xl border border-gray-200 bg-white/70 px-4 py-2.5 text-sm font-medium text-gray-700 backdrop-blur dark:border-dark-800 dark:bg-dark-900/60 dark:text-gray-200"
            >
              <span class="font-mono text-xs text-primary-600 dark:text-primary-400">
                {{ provider.mark }}
              </span>
              {{ provider.name }}
              <span
                class="rounded-md px-1.5 py-0.5 text-[11px] font-medium"
                :class="provider.available
                  ? 'bg-primary-50 text-primary-700 dark:bg-primary-950/60 dark:text-primary-300'
                  : 'bg-gray-100 text-gray-500 dark:bg-dark-800 dark:text-dark-400'"
              >
                {{ provider.available ? t('home.providers.supported') : t('home.providers.soon') }}
              </span>
            </li>
          </ul>
        </div>
      </section>

      <!-- CTA -->
      <section class="border-t border-gray-200/70 dark:border-dark-800/70">
        <div class="mx-auto max-w-6xl px-6 py-16 lg:py-20">
          <div
            class="flex flex-col items-start justify-between gap-8 rounded-2xl bg-gray-900 px-8 py-12 sm:px-12 lg:flex-row lg:items-center dark:bg-dark-900"
          >
            <div class="max-w-xl">
              <h2 class="text-2xl font-semibold tracking-tight text-white sm:text-3xl">
                {{ t('home.cta.title') }}
              </h2>
              <p class="mt-3 text-sm leading-relaxed text-gray-300">{{ t('home.cta.description') }}</p>
            </div>
            <router-link
              :to="isAuthenticated ? dashboardPath : '/login'"
              class="inline-flex h-11 shrink-0 items-center gap-2 rounded-xl bg-primary-500 px-6 text-sm font-medium text-white transition-colors hover:bg-primary-400"
            >
              {{ isAuthenticated ? t('home.goToDashboard') : t('home.cta.button') }}
              <Icon name="arrowRight" size="sm" />
            </router-link>
          </div>
        </div>
      </section>
    </main>

    <!-- Footer -->
    <footer class="relative z-10 border-t border-gray-200/70 dark:border-dark-800/70">
      <div
        class="mx-auto flex max-w-6xl flex-col items-center justify-between gap-4 px-6 py-8 text-center sm:flex-row sm:text-left"
      >
        <p class="text-sm text-gray-500 dark:text-dark-400">
          &copy; {{ currentYear }} {{ siteName }}. {{ t('home.footer.allRightsReserved') }}
        </p>
        <div class="flex items-center gap-5">
          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="text-sm text-gray-500 transition-colors hover:text-gray-800 dark:text-dark-400 dark:hover:text-white"
          >
            {{ t('home.docs') }}
          </a>
          <a
            :href="githubUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="text-sm text-gray-500 transition-colors hover:text-gray-800 dark:text-dark-400 dark:hover:text-white"
          >
            GitHub
          </a>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
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
const docUrl = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.doc_url || appStore.docUrl || ''))
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

// GitHub URL
const githubUrl = 'https://github.com/Xeubad/codeapi'

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

// Current year for footer
const currentYear = computed(() => new Date().getFullYear())

const tags = computed(() => [
  { label: t('home.tags.subscriptionToApi'), note: t('home.console.oneKey') },
  { label: t('home.tags.stickySession'), note: t('home.console.stickySession') },
  { label: t('home.tags.realtimeBilling'), note: t('home.console.realtimeBilling') }
])

const painPoints = computed(() => [
  { title: t('home.painPoints.items.expensive.title'), desc: t('home.painPoints.items.expensive.desc') },
  { title: t('home.painPoints.items.complex.title'), desc: t('home.painPoints.items.complex.desc') },
  { title: t('home.painPoints.items.unstable.title'), desc: t('home.painPoints.items.unstable.desc') },
  { title: t('home.painPoints.items.noControl.title'), desc: t('home.painPoints.items.noControl.desc') }
])

const capabilities = computed(() => [
  { title: t('home.features.unifiedGateway'), desc: t('home.features.unifiedGatewayDesc') },
  { title: t('home.features.multiAccount'), desc: t('home.features.multiAccountDesc') },
  { title: t('home.features.balanceQuota'), desc: t('home.features.balanceQuotaDesc') }
])

const comparisonRows = computed(() => [
  { feature: t('home.comparison.items.pricing.feature'), official: t('home.comparison.items.pricing.official'), us: t('home.comparison.items.pricing.us') },
  { feature: t('home.comparison.items.models.feature'), official: t('home.comparison.items.models.official'), us: t('home.comparison.items.models.us') },
  { feature: t('home.comparison.items.management.feature'), official: t('home.comparison.items.management.official'), us: t('home.comparison.items.management.us') },
  { feature: t('home.comparison.items.stability.feature'), official: t('home.comparison.items.stability.official'), us: t('home.comparison.items.stability.us') },
  { feature: t('home.comparison.items.control.feature'), official: t('home.comparison.items.control.official'), us: t('home.comparison.items.control.us') }
])

const providers = computed(() => [
  { name: t('home.providers.claude'), mark: 'CL', available: true },
  { name: 'GPT', mark: 'GPT', available: true },
  { name: t('home.providers.gemini'), mark: 'GM', available: true },
  { name: t('home.providers.antigravity'), mark: 'AG', available: true },
  { name: t('home.providers.more'), mark: '+', available: false }
])

// Toggle theme
function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

// Initialize theme
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

  // Check auth state
  authStore.checkAuth()

  // Ensure public settings are loaded (will use cache if already loaded from injected config)
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
})
</script>

<style scoped>
/* 细颗粒纹理：让大面积平涂不呆板，同时避免常见的彩色光斑装饰 */
.brand-grain {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='180' height='180'%3E%3Cfilter id='g'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.85' numOctaves='3' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='180' height='180' filter='url(%23g)' opacity='0.42'/%3E%3C/svg%3E");
  background-size: 180px 180px;
  opacity: 0.035;
}

.dark .brand-grain {
  opacity: 0.06;
}
</style>
