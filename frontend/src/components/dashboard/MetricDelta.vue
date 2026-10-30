<template>
  <span v-if="delta === null" class="inline-flex items-center text-xs text-gray-400 dark:text-dark-500">
    {{ t('dashboard.metrics.noComparison') }}
  </span>
  <span v-else class="inline-flex items-center gap-1 text-xs" :class="toneClass">
    <Icon v-if="delta !== 0" :name="delta > 0 ? 'arrowUp' : 'arrowDown'" size="xs" />
    <span class="font-medium tabular-nums">{{ display }}</span>
    <span class="text-gray-400 dark:text-dark-500">{{ t('dashboard.metrics.vsYesterday') }}</span>
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'

const props = withDefaults(defineProps<{
  current: number
  previous: number | null
  mode?: 'percent' | 'points'
  positiveIsGood?: boolean
}>(), {
  mode: 'percent',
  positiveIsGood: true,
})

const { t } = useI18n()

const delta = computed<number | null>(() => {
  if (props.previous === null) return null
  if (props.mode === 'points') return props.current - props.previous
  if (props.previous === 0) return props.current === 0 ? 0 : null
  return ((props.current - props.previous) / props.previous) * 100
})

const toneClass = computed(() => {
  if (delta.value === null || delta.value === 0) return 'text-gray-500 dark:text-gray-400'
  const isGood = delta.value > 0 ? props.positiveIsGood : !props.positiveIsGood
  return isGood ? 'text-emerald-600 dark:text-emerald-400' : 'text-red-500 dark:text-red-400'
})

const display = computed(() => {
  if (delta.value === null) return ''
  const sign = delta.value > 0 ? '+' : ''
  const suffix = props.mode === 'points' ? t('dashboard.metrics.pointsShort') : '%'
  return `${sign}${delta.value.toFixed(1)}${suffix}`
})
</script>
