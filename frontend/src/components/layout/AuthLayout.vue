<template>
  <div class="relative flex min-h-screen flex-col overflow-x-clip bg-gray-50 dark:bg-dark-950">
    <!-- 与首页同源的毛玻璃背景 -->
    <div class="brand-backdrop" aria-hidden="true">
      <div class="brand-orb -left-28 -top-32 h-[30rem] w-[30rem] bg-[rgba(111,158,18,0.20)]"></div>
      <div class="brand-orb -right-28 top-10 h-[26rem] w-[26rem] bg-[rgba(56,132,160,0.16)]"></div>
      <div class="brand-orb bottom-[-14rem] left-1/3 h-[28rem] w-[28rem] bg-[rgba(255,243,219,0.6)] dark:bg-[rgba(47,74,13,0.35)]"></div>
      <div class="brand-grain"></div>
    </div>

    <!-- 悬浮毛玻璃顶栏 -->
    <header class="relative z-20 px-3 pt-3 sm:px-5 sm:pt-4">
      <nav class="brand-glass-strong mx-auto flex h-14 max-w-6xl items-center gap-2 px-3 sm:px-4">
        <router-link to="/" class="flex min-w-0 items-center gap-2.5">
          <img :src="siteLogo || '/logo.png'" alt="Logo" class="h-8 w-8 shrink-0 object-contain" />
          <span class="truncate text-[15px] font-semibold tracking-tight">{{ siteName }}</span>
        </router-link>

        <div class="ml-auto flex shrink-0 items-center gap-1">
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
            to="/"
            class="ml-1 inline-flex h-9 items-center gap-1.5 rounded-lg px-2.5 text-[13px] font-medium text-gray-600 transition-colors hover:bg-gray-900/5 hover:text-gray-900 dark:text-dark-300 dark:hover:bg-white/5 dark:hover:text-white"
          >
            <Icon name="home" size="sm" />
            <span class="hidden sm:inline">{{ t('home.backHome') }}</span>
          </router-link>
        </div>
      </nav>
    </header>

    <!-- 主体 -->
    <main class="relative z-10 flex flex-1 items-center px-4 py-8 sm:px-6 lg:py-12">
      <div
        class="mx-auto grid w-full max-w-6xl grid-cols-1 items-center gap-10 lg:grid-cols-[minmax(0,1fr)_minmax(0,26.5rem)] lg:gap-16"
      >
        <!-- 品牌侧 -->
        <section class="hidden min-w-0 lg:block">
          <p class="flex items-center gap-2 text-[10px] font-medium uppercase tracking-[0.22em] text-gray-500 dark:text-dark-400">
            <span class="h-1.5 w-1.5 rounded-full bg-primary-500"></span>
            {{ siteSubtitle }}
          </p>

          <h1
            class="mt-5 text-[2.55rem] font-semibold leading-[1.1] tracking-[-0.035em] text-gray-900 dark:text-white"
          >
            <span class="block">{{ t('home.hero.titleLead') }}</span>
            <span class="block">{{ t('home.hero.titleAccent') }}</span>
          </h1>
          <div class="mt-5 h-1 w-12 rounded-full bg-primary-500"></div>
          <p class="mt-5 max-w-md text-[15px] leading-relaxed text-gray-600 dark:text-dark-300">
            {{ t('home.hero.description') }}
          </p>

          <dl class="mt-9 divide-y divide-gray-900/10 border-y border-gray-900/10 dark:divide-white/10 dark:border-white/10">
            <div v-for="(fact, index) in facts" :key="fact.title" class="flex items-baseline gap-4 py-3.5">
              <span class="w-6 shrink-0 font-mono text-[10px] text-primary-600 dark:text-primary-400">
                {{ String(index + 1).padStart(2, '0') }}
              </span>
              <div class="min-w-0">
                <dt class="text-[13px] font-semibold text-gray-900 dark:text-gray-100">{{ fact.title }}</dt>
                <dd class="mt-0.5 text-[12px] text-gray-500 dark:text-dark-400">{{ fact.value }}</dd>
              </div>
            </div>
          </dl>

          <p class="mt-6 flex items-start gap-2.5 text-[12px] leading-relaxed text-gray-500 dark:text-dark-400">
            <Icon name="shield" size="sm" class="mt-0.5 shrink-0 text-primary-600 dark:text-primary-400" />
            {{ t('home.auth.secureNote') }}
          </p>
        </section>

        <!-- 表单卡片 -->
        <section class="w-full min-w-0">
          <div class="brand-glass p-6 sm:p-8">
            <slot />
          </div>

          <div class="mt-5 text-center text-sm">
            <slot name="footer" />
          </div>
        </section>
      </div>
    </main>

    <!-- 版权 -->
    <footer class="relative z-10 px-6 pb-6 pt-2 text-center text-[11px] text-gray-500 dark:text-dark-500">
      &copy; {{ currentYear }} {{ siteName }}. All rights reserved.
    </footer>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores'
import { sanitizeUrl } from '@/utils/url'
import Icon from '@/components/icons/Icon.vue'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'

const { t } = useI18n()
const appStore = useAppStore()

const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'CodeAPI')
const siteLogo = computed(() =>
  sanitizeUrl(appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '', {
    allowRelative: true,
    allowDataUrl: true
  })
)
const siteSubtitle = computed(
  () => appStore.cachedPublicSettings?.site_subtitle || 'Subscription to API Conversion Platform'
)

const facts = computed(() => [
  { title: t('home.facts.models.title'), value: t('home.facts.models.value') },
  { title: t('home.facts.sticky.title'), value: t('home.facts.sticky.value') },
  { title: t('home.facts.billing.title'), value: t('home.facts.billing.value') }
])

const isDark = ref(document.documentElement.classList.contains('dark'))
const currentYear = computed(() => new Date().getFullYear())

function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

onMounted(() => {
  const savedTheme = localStorage.getItem('theme')
  if (
    savedTheme === 'dark' ||
    (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)
  ) {
    isDark.value = true
    document.documentElement.classList.add('dark')
  }

  appStore.fetchPublicSettings()
})
</script>
