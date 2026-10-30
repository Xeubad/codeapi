<template>
  <div class="card flex h-full min-h-0 flex-col overflow-hidden">
    <div class="flex items-center justify-between gap-3 border-b border-gray-100 px-6 py-4 dark:border-dark-700">
      <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('dashboard.announcements.title') }}</h2>
      <div class="flex items-center gap-2">
        <span v-if="markable && unreadCount > 0" class="badge badge-primary">
          {{ t('dashboard.announcements.unreadCount', { count: unreadCount }) }}
        </span>
        <slot name="action" />
      </div>
    </div>
    <div class="flex min-h-0 flex-1 flex-col p-6">
      <div v-if="loading && items.length === 0" class="flex flex-1 items-center justify-center py-12">
        <LoadingSpinner size="lg" />
      </div>
      <div v-else-if="groups.length === 0" class="flex flex-1 flex-col items-center justify-center py-10 text-center">
        <div class="mb-3 flex h-12 w-12 items-center justify-center rounded-xl bg-gray-100 dark:bg-dark-700/60">
          <Icon name="bell" size="md" class="text-gray-400 dark:text-dark-500" />
        </div>
        <p class="text-sm font-medium text-gray-900 dark:text-white">{{ t('dashboard.announcements.empty') }}</p>
        <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ t('dashboard.announcements.emptyDescription') }}</p>
      </div>
      <div v-else class="min-h-0 flex-1 overflow-y-auto pr-1">
        <div class="relative">
          <div class="absolute bottom-2 left-[7px] top-2 w-px bg-gray-200 dark:bg-dark-700"></div>
          <section
            v-for="group in groups"
            :key="group.date"
            class="relative pb-5 pl-7 last:pb-0"
            :data-date="group.date"
          >
            <span class="absolute left-0 top-1 h-[15px] w-[15px] rounded-full border-2 border-primary-500 bg-white dark:bg-dark-800"></span>
            <time class="text-xs font-semibold text-gray-500 dark:text-gray-400">{{ group.label }}</time>
            <div class="mt-2 space-y-2">
              <article
                v-for="item in group.items"
                :key="item.id"
                class="rounded-lg border border-gray-100 bg-gray-50/60 p-3 dark:border-dark-700 dark:bg-dark-800/40"
                :data-announcement-id="item.id"
              >
                <div class="flex items-start gap-2">
                  <span
                    v-if="markable && !item.read_at"
                    class="mt-1.5 h-1.5 w-1.5 shrink-0 rounded-full bg-primary-500"
                  ></span>
                  <button
                    type="button"
                    class="min-w-0 flex-1 text-left text-sm font-medium text-gray-900 hover:text-primary-600 dark:text-white dark:hover:text-primary-400"
                    @click="toggle(item)"
                  >
                    {{ item.title }}
                  </button>
                  <time class="shrink-0 text-[11px] text-gray-400 dark:text-dark-500">{{ formatTime(item.created_at) }}</time>
                </div>
                <div
                  v-if="item.content"
                  class="markdown-body announcement-content mt-1.5 text-sm"
                  :class="{ 'line-clamp-2': expandedId !== item.id }"
                  v-html="renderedContents.get(item.id)"
                ></div>
                <button
                  v-if="item.content"
                  type="button"
                  class="mt-1 text-xs font-medium text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300"
                  @click="toggle(item)"
                >
                  {{ expandedId === item.id ? t('dashboard.announcements.collapse') : t('dashboard.announcements.expand') }}
                </button>
              </article>
            </div>
          </section>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import Icon from '@/components/icons/Icon.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import { formatDate, formatDateLocalInput, formatTime } from '@/utils/format'
import '@/styles/announcement-markdown.css'

export interface AnnouncementTimelineItem {
  id: number
  title: string
  content: string
  created_at: string
  read_at?: string | null
}

const props = withDefaults(defineProps<{
  items: AnnouncementTimelineItem[]
  loading?: boolean
  markable?: boolean
  unreadCount?: number
  limit?: number
}>(), {
  loading: false,
  markable: false,
  unreadCount: 0,
  limit: 0,
})

const emit = defineEmits<{ read: [id: number] }>()

marked.setOptions({
  breaks: true,
  gfm: true,
})

const { t } = useI18n()
const expandedId = ref<number | null>(null)

const groups = computed(() => {
  const sorted = [...props.items].sort(
    (a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
  )
  const visible = props.limit > 0 ? sorted.slice(0, props.limit) : sorted
  const byDate = new Map<string, AnnouncementTimelineItem[]>()
  for (const item of visible) {
    const date = formatDateLocalInput(new Date(item.created_at))
    const items = byDate.get(date)
    if (items) {
      items.push(item)
    } else {
      byDate.set(date, [item])
    }
  }
  return Array.from(byDate.entries()).map(([date, items]) => ({
    date,
    label: formatDate(items[0].created_at, {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      weekday: 'short',
    }),
    items,
  }))
})

const renderedContents = computed(() => {
  const map = new Map<number, string>()
  for (const item of props.items) {
    if (!item.content) continue
    map.set(item.id, DOMPurify.sanitize(marked.parse(item.content) as string))
  }
  return map
})

const toggle = (item: AnnouncementTimelineItem) => {
  const opening = expandedId.value !== item.id
  expandedId.value = opening ? item.id : null
  if (opening && props.markable && !item.read_at) {
    emit('read', item.id)
  }
}
</script>

<style scoped>
.announcement-content {
  font-size: 0.875rem;
  line-height: 1.55;
}

.announcement-content :deep(h1),
.announcement-content :deep(h2),
.announcement-content :deep(h3),
.announcement-content :deep(h4) {
  margin: 0.5rem 0 0.25rem;
  padding: 0;
  border: 0;
  font-size: 0.875rem;
  font-weight: 600;
}

.announcement-content :deep(p) {
  margin: 0 0 0.25rem;
}

.announcement-content :deep(ul),
.announcement-content :deep(ol) {
  margin: 0 0 0.25rem 1rem;
}

.announcement-content :deep(blockquote),
.announcement-content :deep(pre),
.announcement-content :deep(table),
.announcement-content :deep(img),
.announcement-content :deep(hr) {
  margin: 0.5rem 0;
}

.announcement-content :deep(blockquote::before) {
  content: none;
}

.announcement-content :deep(pre) {
  padding: 0.5rem;
}

.announcement-content :deep(th),
.announcement-content :deep(td) {
  padding: 0.25rem 0.5rem;
}
</style>
