<template>
  <div class="card flex h-full flex-col p-5">
    <div class="flex items-center gap-3">
      <div class="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl" :class="toneClass">
        <Icon :name="icon" size="md" :stroke-width="2" />
      </div>
      <p class="text-sm font-medium text-gray-500 dark:text-gray-400">{{ label }}</p>
    </div>

    <div class="mt-5 flex flex-wrap items-end gap-x-3 gap-y-1">
      <slot name="value" />
      <MetricDelta
        v-if="deltaPrevious !== undefined"
        :current="deltaCurrent ?? 0"
        :previous="deltaPrevious"
        :mode="deltaMode"
        :positive-is-good="deltaPositiveIsGood"
        class="mb-1"
      />
    </div>

    <div v-if="$slots.meta" class="mt-2">
      <slot name="meta" />
    </div>

    <div v-if="$slots.footer" class="mt-auto border-t border-gray-100 pt-3 dark:border-dark-700">
      <slot name="footer" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import Icon from '@/components/icons/Icon.vue'
import MetricDelta from '@/components/dashboard/MetricDelta.vue'

type Tone = 'blue' | 'green' | 'emerald' | 'purple' | 'amber' | 'violet' | 'rose' | 'indigo'
type IconName = 'chart' | 'cube' | 'dollar' | 'userPlus' | 'key' | 'server' | 'bolt' | 'clock' | 'database'

const props = withDefaults(defineProps<{
  label: string
  icon: IconName
  tone?: Tone
  deltaCurrent?: number
  deltaPrevious?: number | null
  deltaMode?: 'percent' | 'points'
  deltaPositiveIsGood?: boolean
}>(), {
  tone: 'blue',
  deltaMode: 'percent',
  deltaPositiveIsGood: true,
})

const toneClasses: Record<Tone, string> = {
  blue: 'bg-blue-100 text-blue-600 dark:bg-blue-900/30 dark:text-blue-400',
  green: 'bg-green-100 text-green-600 dark:bg-green-900/30 dark:text-green-400',
  emerald: 'bg-emerald-100 text-emerald-600 dark:bg-emerald-900/30 dark:text-emerald-400',
  purple: 'bg-purple-100 text-purple-600 dark:bg-purple-900/30 dark:text-purple-400',
  amber: 'bg-amber-100 text-amber-600 dark:bg-amber-900/30 dark:text-amber-400',
  violet: 'bg-violet-100 text-violet-600 dark:bg-violet-900/30 dark:text-violet-400',
  rose: 'bg-rose-100 text-rose-600 dark:bg-rose-900/30 dark:text-rose-400',
  indigo: 'bg-indigo-100 text-indigo-600 dark:bg-indigo-900/30 dark:text-indigo-400',
}

const toneClass = computed(() => toneClasses[props.tone])
</script>
