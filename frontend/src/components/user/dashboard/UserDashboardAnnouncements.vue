<template>
  <AnnouncementTimeline
    :items="announcements"
    :loading="loading"
    :limit="5"
    markable
    :unread-count="unreadCount"
    @read="handleRead"
  >
    <template #action>
      <button
        type="button"
        class="text-xs font-medium text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300"
        @click="announcementStore.openListModal()"
      >
        {{ t('dashboard.announcements.viewAll') }} →
      </button>
    </template>
  </AnnouncementTimeline>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { storeToRefs } from 'pinia'
import { useAnnouncementStore } from '@/stores/announcements'
import AnnouncementTimeline from '@/components/common/AnnouncementTimeline.vue'

const { t } = useI18n()
const announcementStore = useAnnouncementStore()
const { announcements, loading } = storeToRefs(announcementStore)

const unreadCount = computed(() => announcementStore.unreadCount)

const handleRead = (id: number) => {
  void announcementStore.markAsRead(id)
}

onMounted(() => {
  void announcementStore.fetchAnnouncements()
})
</script>
