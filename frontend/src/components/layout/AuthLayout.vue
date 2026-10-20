<template>
  <div class="relative flex min-h-screen flex-col overflow-hidden bg-gray-50 dark:bg-dark-950">
    <!-- 背景：柔和色场 + 细颗粒，与首页同一套语言 -->
    <div class="pointer-events-none fixed inset-0 z-0 overflow-hidden">
      <div
        class="absolute inset-0 bg-[radial-gradient(85%_55%_at_18%_-8%,rgba(111,158,18,0.28),transparent_60%)] dark:bg-[radial-gradient(85%_55%_at_18%_-8%,rgba(111,158,18,0.28),transparent_60%)]"
      ></div>
      <div
        class="absolute inset-0 bg-[radial-gradient(75%_55%_at_95%_12%,rgba(21,94,117,0.26),transparent_62%)] dark:bg-[radial-gradient(75%_55%_at_95%_12%,rgba(56,132,160,0.22),transparent_62%)]"
      ></div>
      <div
        class="absolute inset-0 bg-[radial-gradient(70%_55%_at_45%_112%,rgba(255,246,220,0.85),transparent_62%)] dark:bg-[radial-gradient(70%_55%_at_45%_112%,rgba(40,56,34,0.6),transparent_62%)]"
      ></div>
      <!-- 卡片背后的柔光团，让毛玻璃有可透的内容 -->
      <div
        class="absolute -top-24 -right-16 h-[26rem] w-[26rem] rounded-full bg-[radial-gradient(circle,rgba(111,158,18,0.20),transparent_68%)] blur-2xl dark:bg-[radial-gradient(circle,rgba(111,158,18,0.26),transparent_68%)]"
      ></div>
      <div
        class="absolute bottom-[-10rem] left-1/4 h-[24rem] w-[24rem] rounded-full bg-[radial-gradient(circle,rgba(21,94,117,0.16),transparent_70%)] blur-2xl dark:bg-[radial-gradient(circle,rgba(56,132,160,0.22),transparent_70%)]"
      ></div>
      <div class="brand-grain absolute inset-0"></div>
    </div>

    <!-- 顶栏 -->
    <header class="relative z-10 border-b border-gray-200/60 dark:border-dark-800/60">
      <nav class="mx-auto flex h-16 max-w-6xl items-center justify-between gap-4 px-6">
        <router-link to="/" class="flex min-w-0 items-center gap-2.5">
          <img
            :src="siteLogo || '/logo.png'"
            alt="Logo"
            class="h-8 w-8 shrink-0 object-contain"
          />
          <span class="truncate text-[15px] font-semibold tracking-tight">{{ siteName }}</span>
        </router-link>

        <div class="flex shrink-0 items-center gap-1">
          <LocaleSwitcher />
          <button
            class="flex h-9 w-9 items-center justify-center rounded-lg text-gray-500 transition-colors hover:bg-gray-900/5 hover:text-gray-900 dark:text-dark-400 dark:hover:bg-white/5 dark:hover:text-white"
            :title="isDark ? t('nav.lightMode') : t('nav.darkMode')"
            @click="toggleTheme"
          >
            <Icon v-if="isDark" name="sun" size="sm" />
            <Icon v-else name="moon" size="sm" />
          </button>
          <router-link
            to="/"
            class="ml-1 inline-flex h-9 items-center gap-1.5 rounded-lg px-3 text-sm font-medium text-gray-600 transition-colors hover:bg-gray-900/5 hover:text-gray-900 dark:text-dark-300 dark:hover:bg-white/5 dark:hover:text-white"
          >
            <Icon name="home" size="sm" />
            <span class="hidden sm:inline">{{ t('home.backHome') }}</span>
          </router-link>
        </div>
      </nav>
    </header>

    <!-- 主体 -->
    <main class="relative z-10 flex flex-1 items-center justify-center px-4 py-10 sm:px-6 lg:py-16">
      <div
        class="grid w-full max-w-5xl items-center gap-12 lg:grid-cols-[minmax(0,1fr)_minmax(0,25rem)] lg:gap-20"
      >
        <!-- 品牌侧 -->
        <section class="hidden lg:block">
          <p class="flex items-center gap-2 text-sm font-medium text-primary-700 dark:text-primary-400">
            <span class="h-1.5 w-1.5 rounded-full bg-primary-500"></span>
            {{ siteSubtitle }}
          </p>
          <h1
            class="mt-5 text-[2.6rem] font-semibold leading-[1.12] tracking-[-0.02em] text-gray-900 dark:text-white"
          >
            <span class="block">{{ t('home.heroTitleLead') }}</span>
            <span class="block">{{ t('home.heroTitleMain') }}</span>
          </h1>
          <p class="mt-5 max-w-md text-base leading-relaxed text-gray-600 dark:text-dark-300">
            {{ t('home.heroDescription') }}
          </p>

          <ul class="mt-10 space-y-4">
            <li v-for="tag in tags" :key="tag.label" class="flex items-start gap-3">
              <span class="mt-1.5 h-1.5 w-1.5 shrink-0 rounded-full bg-primary-500"></span>
              <span class="text-sm text-gray-700 dark:text-dark-200">
                <span class="font-medium text-gray-900 dark:text-gray-100">{{ tag.label }}</span>
                <span class="text-gray-500 dark:text-dark-400"> · {{ tag.note }}</span>
              </span>
            </li>
          </ul>
        </section>

        <!-- 表单 -->
        <section class="w-full">
          <div class="auth-glass p-7 sm:p-8">
            <slot />
          </div>

          <div class="mt-6 text-center text-sm">
            <slot name="footer" />
          </div>
        </section>
      </div>
    </main>

    <!-- 版权 -->
    <footer class="relative z-10 px-6 pb-8 text-center text-xs text-gray-500 dark:text-dark-500">
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

const tags = computed(() => [
  { label: t('home.tags.subscriptionToApi'), note: t('home.console.oneKey') },
  { label: t('home.tags.stickySession'), note: t('home.console.stickySession') },
  { label: t('home.tags.realtimeBilling'), note: t('home.console.realtimeBilling') }
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

<style scoped>
/* 苹果式毛玻璃：半透明面板 + 高斯模糊 + 细高光描边 */
.auth-glass {
  @apply rounded-2xl border border-white/80 bg-white/55;
  @apply backdrop-blur-2xl backdrop-saturate-150;
  @apply shadow-[0_28px_80px_-36px_rgba(11,14,20,0.55)];
}

.dark .auth-glass {
  @apply border-white/10 bg-dark-900/55;
  @apply shadow-[0_28px_80px_-36px_rgba(0,0,0,0.85)];
}

.brand-grain {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='180' height='180'%3E%3Cfilter id='g'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.85' numOctaves='3' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='180' height='180' filter='url(%23g)' opacity='0.42'/%3E%3C/svg%3E");
  background-size: 180px 180px;
  opacity: 0.035;
}

.dark .brand-grain {
  opacity: 0.06;
}
</style>
