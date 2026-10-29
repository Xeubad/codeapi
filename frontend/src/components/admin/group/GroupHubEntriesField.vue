<template>
  <div
    v-if="vendorsWithEntries.length > 0"
    class="mt-4 border-t border-dashed border-gray-200 pt-4 dark:border-dark-700"
  >
    <div class="mb-2 flex flex-wrap items-center justify-between gap-2">
      <span class="text-sm font-medium text-gray-700 dark:text-gray-300">
        {{ t('admin.groups.hubRef.title') }}
      </span>
      <span class="text-xs text-gray-500 dark:text-gray-400">{{ t('admin.groups.hubRef.hint') }}</span>
    </div>
    <div class="max-h-48 space-y-2 overflow-y-auto rounded-lg border border-gray-200 p-2 dark:border-dark-600">
      <div v-for="vendor in vendorsWithEntries" :key="vendor.id">
        <div class="mb-1 text-xs font-medium text-gray-500 dark:text-gray-400">{{ vendor.name }}</div>
        <label
          v-for="entry in vendor.entries"
          :key="entry.account_id"
          class="flex items-center gap-2 rounded-md px-1 py-1 text-sm hover:bg-gray-50 dark:hover:bg-dark-700"
        >
          <input
            type="checkbox"
            class="cursor-pointer"
            :checked="modelValue.includes(entry.account_id)"
            @change="toggle(entry.account_id)"
          />
          <span class="truncate text-gray-900 dark:text-white">{{ entry.name }}</span>
          <span class="badge badge-primary">{{ platformLabel(entry.platform) }}</span>
          <span
            class="text-xs"
            :class="(entry.models?.length || 0) > 0 ? 'text-gray-500 dark:text-gray-400' : 'text-amber-600 dark:text-amber-400'"
          >
            {{ t('admin.groups.hubRef.models', { count: entry.models?.length || 0 }) }}
          </span>
          <span v-if="!entry.schedulable" class="badge badge-gray">
            {{ t('admin.groups.hubRef.unschedulable') }}
          </span>
        </label>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { platformDisplayName } from '@/constants/platformCatalog'
import type { UpstreamEntry, UpstreamVendor } from '@/api/admin/upstreamHub'

const props = defineProps<{
  platform: string
  vendors: UpstreamVendor[]
  entries: UpstreamEntry[]
  modelValue: number[]
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: number[]): void
}>()

const { t } = useI18n()

const vendorsWithEntries = computed(() =>
  props.vendors
    .map((vendor) => ({
      ...vendor,
      entries: props.entries.filter(
        (entry) =>
          entry.vendor_id === vendor.id &&
          (props.platform === 'composite' || entry.platform === props.platform)
      )
    }))
    .filter((vendor) => vendor.entries.length > 0)
)

function platformLabel(platform: string): string {
  return platformDisplayName(platform) || platform
}

function toggle(accountID: number) {
  const next = props.modelValue.includes(accountID)
    ? props.modelValue.filter((id) => id !== accountID)
    : [...props.modelValue, accountID]
  emit('update:modelValue', next)
}
</script>
