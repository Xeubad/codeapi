<template>
  <AppLayout>
    <div class="space-y-4">
      <div class="page-header">
        <div>
          <h1 class="page-title">{{ t('admin.aggregation.title') }}</h1>
          <p class="page-subtitle">{{ t('admin.aggregation.description') }}</p>
        </div>
        <div class="flex flex-wrap items-center gap-2">
          <Select
            v-if="compositeGroups.length"
            v-model="selectedGroupId"
            class="w-56"
            :options="compositeGroupOptions"
            @update:model-value="loadGroupData"
          />
          <button class="btn btn-secondary btn-sm" :disabled="loading" @click="loadAll">
            {{ t('admin.aggregation.refresh') }}
          </button>
        </div>
      </div>

      <EmptyState
        v-if="!loading && compositeGroups.length === 0"
        :title="t('admin.aggregation.noCompositeGroups')"
        :description="t('admin.aggregation.noCompositeGroupsHint')"
      />

      <template v-else>
        <div class="grid gap-4 xl:grid-cols-2">
          <!-- 供给汇入 -->
          <div class="card p-4">
            <div class="mb-3 flex items-center justify-between">
              <div>
                <h2 class="text-base font-semibold text-gray-900 dark:text-white">
                  {{ t('admin.aggregation.supply') }}
                </h2>
                <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('admin.aggregation.supplyHint') }}</p>
              </div>
              <span class="badge badge-gray">{{ t('admin.aggregation.attachedCount', { count: attachedCount }) }}</span>
            </div>
            <div v-if="entries.length === 0" class="py-6 text-center text-sm text-gray-500">
              {{ t('admin.aggregation.noEntries') }}
            </div>
            <div v-else class="max-h-[26rem] space-y-3 overflow-y-auto pr-1">
              <div v-for="vendor in vendorGroups" :key="vendor.id" class="space-y-1">
                <div class="flex items-center justify-between text-xs text-gray-500 dark:text-gray-400">
                  <span class="font-medium text-gray-700 dark:text-gray-300">{{ vendor.name }}</span>
                  <span>
                    {{ t('admin.aggregation.vendorBudget', {
                      used: vendor.allocated_concurrency,
                      total: vendor.concurrency_budget || '∞'
                    }) }}
                  </span>
                </div>
                <div
                  v-for="entry in vendor.entries"
                  :key="entry.account_id"
                  class="flex items-center gap-2 rounded-lg border border-gray-200 px-2 py-1.5 dark:border-dark-600"
                  :class="isAttached(entry) ? 'bg-primary-50/60 dark:bg-primary-900/10' : ''"
                >
                  <input
                    type="checkbox"
                    class="cursor-pointer"
                    :checked="isAttached(entry)"
                    @change="toggleAttach(entry)"
                  />
                  <div class="min-w-0 flex-1">
                    <div class="truncate text-sm text-gray-900 dark:text-white">{{ entry.name }}</div>
                    <div class="text-xs text-gray-500 dark:text-gray-400">
                      {{ entry.platform }} · {{ t('admin.aggregation.modelsCount', { count: entry.models?.length || 0 }) }}
                    </div>
                  </div>
                  <span
                    v-if="isAttached(entry) && !isEnabledHere(entry)"
                    class="badge badge-gray"
                  >{{ t('admin.aggregation.disabledHere') }}</span>
                  <Toggle
                    v-if="isAttached(entry)"
                    :model-value="isEnabledHere(entry)"
                    @update:model-value="() => toggleEntryEnabled(entry)"
                  />
                </div>
              </div>
            </div>
          </div>

          <!-- 归一模型 -->
          <div class="card p-4">
            <div class="mb-3">
              <h2 class="text-base font-semibold text-gray-900 dark:text-white">
                {{ t('admin.aggregation.modelRoutes') }}
              </h2>
              <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('admin.aggregation.modelRoutesHint') }}</p>
            </div>
            <div class="table-container">
              <table class="table">
                <thead>
                  <tr>
                    <th>{{ t('admin.aggregation.publicModel') }}</th>
                    <th>{{ t('admin.aggregation.targetPlatform') }}</th>
                    <th>{{ t('admin.aggregation.upstreamModel') }}</th>
                    <th>{{ t('admin.aggregation.priority') }}</th>
                    <th>{{ t('admin.aggregation.enabled') }}</th>
                    <th class="text-right">{{ t('common.actions') }}</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="route in routes" :key="route.id">
                    <td class="font-mono text-xs">{{ route.public_model }}</td>
                    <td><span class="badge badge-primary">{{ route.target_platform }}</span></td>
                    <td class="font-mono text-xs">{{ route.upstream_model || route.public_model }}</td>
                    <td>{{ route.priority }}</td>
                    <td>
                      <Toggle :model-value="route.enabled" @update:model-value="() => toggleRouteEnabled(route)" />
                    </td>
                    <td class="text-right">
                      <button class="btn btn-sm btn-ghost text-red-600" @click="removeRoute(route)">
                        {{ t('admin.aggregation.delete') }}
                      </button>
                    </td>
                  </tr>
                  <tr v-if="routes.length === 0">
                    <td colspan="6" class="py-4 text-center text-xs text-gray-500">
                      {{ t('admin.aggregation.noRoutes') }}
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
            <div class="mt-3 grid gap-2 sm:grid-cols-2">
              <input v-model="routeForm.public_model" class="input" :placeholder="t('admin.aggregation.publicModelPlaceholder')" />
              <Select v-model="routeForm.target_platform" :options="platformOptions" />
              <input v-model="routeForm.upstream_model" class="input" :placeholder="t('admin.aggregation.upstreamModelPlaceholder')" />
              <input v-model.number="routeForm.priority" class="input" type="number" min="1" :placeholder="t('admin.aggregation.priority')" />
            </div>
            <div class="mt-2 flex items-center justify-between">
              <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('admin.aggregation.routeFormHint') }}</p>
              <button class="btn btn-primary btn-sm" :disabled="saving || !routeForm.public_model.trim()" @click="addRoute">
                {{ t('admin.aggregation.addRoute') }}
              </button>
            </div>
          </div>
        </div>

        <!-- 链路视图 -->
        <div class="card p-4">
          <div class="mb-3 flex flex-wrap items-center justify-between gap-2">
            <div>
              <h2 class="text-base font-semibold text-gray-900 dark:text-white">
                {{ t('admin.aggregation.chainView') }}
              </h2>
              <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('admin.aggregation.chainHint') }}</p>
            </div>
            <div class="flex items-center gap-3 text-xs text-gray-500 dark:text-gray-400">
              <span class="inline-flex items-center gap-1">
                <span class="h-2 w-4 rounded bg-primary-500"></span>{{ t('admin.aggregation.legendOn') }}
              </span>
              <span class="inline-flex items-center gap-1">
                <span class="h-2 w-4 rounded border border-dashed border-gray-400"></span>{{ t('admin.aggregation.legendOff') }}
              </span>
            </div>
          </div>
          <div class="grid gap-3 md:grid-cols-4">
            <div class="space-y-2">
              <h3 class="text-xs font-semibold uppercase text-gray-500 dark:text-gray-400">
                {{ t('admin.aggregation.vendors') }}
              </h3>
              <div
                v-for="vendor in vendors"
                :key="vendor.id"
                class="rounded-lg border border-gray-200 px-2 py-1.5 text-sm dark:border-dark-600"
              >
                <div class="truncate font-medium text-gray-900 dark:text-white">{{ vendor.name }}</div>
                <div class="text-xs text-gray-500 dark:text-gray-400">
                  {{ t('admin.aggregation.vendorBudget', {
                    used: vendor.allocated_concurrency,
                    total: vendor.concurrency_budget || '∞'
                  }) }}
                </div>
              </div>
            </div>
            <div class="space-y-2">
              <h3 class="text-xs font-semibold uppercase text-gray-500 dark:text-gray-400">
                {{ t('admin.aggregation.entries') }}
              </h3>
              <div
                v-for="entry in entries"
                :key="entry.account_id"
                class="rounded-lg border px-2 py-1.5 text-sm"
                :class="
                  isAttached(entry) && isEnabledHere(entry)
                    ? 'border-primary-300 bg-primary-50 text-primary-800 dark:border-primary-700 dark:bg-primary-900/20 dark:text-primary-300'
                    : 'border-dashed border-gray-300 text-gray-500 dark:border-dark-600 dark:text-dark-300'
                "
              >
                <div class="truncate">{{ entry.name }}</div>
                <div class="text-xs opacity-80">
                  {{ entry.platform }} ·
                  {{ isAttached(entry) && !isEnabledHere(entry) ? t('admin.aggregation.disabledHere') : t('admin.aggregation.modelsCount', { count: entry.models?.length || 0 }) }}
                </div>
              </div>
            </div>
            <div class="space-y-2">
              <h3 class="text-xs font-semibold uppercase text-gray-500 dark:text-gray-400">
                {{ t('admin.aggregation.group') }}
              </h3>
              <div v-if="selectedGroup" class="rounded-lg border border-primary-300 bg-primary-50 px-2 py-1.5 text-sm text-primary-800 dark:border-primary-700 dark:bg-primary-900/20 dark:text-primary-300">
                <div class="truncate font-medium">{{ selectedGroup.name }}</div>
                <div class="text-xs opacity-80">
                  {{ t('admin.aggregation.attachedCount', { count: attachedCount }) }}
                  · RPM {{ selectedGroup.rpm_limit || '∞' }}
                </div>
              </div>
            </div>
            <div class="space-y-2">
              <h3 class="text-xs font-semibold uppercase text-gray-500 dark:text-gray-400">
                {{ t('admin.aggregation.customers') }}
              </h3>
              <div class="rounded-lg border border-gray-200 px-2 py-1.5 text-sm dark:border-dark-600">
                <div class="text-gray-900 dark:text-white">
                  {{ t('admin.aggregation.activeKeys', { count: stats?.active_api_keys || 0 }) }}
                </div>
                <div class="text-xs text-gray-500 dark:text-gray-400">
                  {{ t('admin.aggregation.totalKeys', { count: stats?.total_api_keys || 0 }) }}
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores/app'
import { extractApiErrorMessage } from '@/utils/apiError'
import upstreamHubAPI, { type UpstreamEntry, type UpstreamVendor } from '@/api/admin/upstreamHub'
import groupsAPI from '@/api/admin/groups'
import { CONCRETE_PLATFORM_OPTIONS } from '@/constants/platforms'
import type { AdminGroup, CompositeModelRoute } from '@/types'
import AppLayout from '@/components/layout/AppLayout.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import Select from '@/components/common/Select.vue'
import Toggle from '@/components/common/Toggle.vue'

const { t } = useI18n()
const appStore = useAppStore()

const loading = ref(false)
const saving = ref(false)
const groups = ref<AdminGroup[]>([])
const vendors = ref<UpstreamVendor[]>([])
const entries = ref<UpstreamEntry[]>([])
const routes = ref<CompositeModelRoute[]>([])
const stats = ref<{ total_api_keys: number; active_api_keys: number } | null>(null)
const selectedGroupId = ref<number | null>(null)

const compositeGroups = computed(() => groups.value.filter((group) => group.platform === 'composite'))
const compositeGroupOptions = computed(() =>
  compositeGroups.value.map((group) => ({ label: group.name, value: group.id }))
)
const selectedGroup = computed(
  () => compositeGroups.value.find((group) => group.id === selectedGroupId.value) || null
)
const platformOptions = computed(() =>
  CONCRETE_PLATFORM_OPTIONS.map((platform) => ({ label: platform.label, value: platform.value }))
)

const vendorGroups = computed(() =>
  vendors.value.map((vendor) => ({
    ...vendor,
    entries: entries.value.filter((entry) => entry.vendor_id === vendor.id)
  }))
)

const attachedCount = computed(() => entries.value.filter((entry) => isAttached(entry)).length)

function groupRef(entry: UpstreamEntry) {
  return entry.groups?.find((group) => group.group_id === selectedGroupId.value)
}

function isAttached(entry: UpstreamEntry): boolean {
  return !!groupRef(entry)
}

function isEnabledHere(entry: UpstreamEntry): boolean {
  return groupRef(entry)?.enabled ?? false
}

function errorMessage(err: unknown): string {
  return extractApiErrorMessage(err, t('admin.aggregation.loadFailed'))
}

async function loadGroups() {
  groups.value = await groupsAPI.getAllIncludingInactive()
  if (!selectedGroupId.value || !compositeGroups.value.some((group) => group.id === selectedGroupId.value)) {
    selectedGroupId.value = compositeGroups.value[0]?.id ?? null
  }
}

async function loadGroupData() {
  if (!selectedGroupId.value) {
    routes.value = []
    stats.value = null
    return
  }
  try {
    const [routeList, groupStats] = await Promise.all([
      groupsAPI.listCompositeRoutes(selectedGroupId.value),
      groupsAPI.getStats(selectedGroupId.value)
    ])
    routes.value = routeList
    stats.value = { total_api_keys: groupStats.total_api_keys, active_api_keys: groupStats.active_api_keys }
  } catch (err) {
    appStore.showError(errorMessage(err))
  }
}

async function loadAll() {
  loading.value = true
  try {
    const [vendorList, entryList] = await Promise.all([upstreamHubAPI.listVendors(), upstreamHubAPI.listAllEntries()])
    vendors.value = vendorList
    entries.value = entryList
    await loadGroups()
    await loadGroupData()
  } catch (err) {
    appStore.showError(errorMessage(err))
  } finally {
    loading.value = false
  }
}

async function refreshEntries() {
  entries.value = await upstreamHubAPI.listAllEntries()
  vendors.value = await upstreamHubAPI.listVendors()
}

async function toggleAttach(entry: UpstreamEntry) {
  if (!selectedGroupId.value) return
  const otherGroups = (entry.groups || [])
    .map((group) => group.group_id)
    .filter((groupId) => groupId !== selectedGroupId.value)
  const nextGroupIds = isAttached(entry) ? otherGroups : [...otherGroups, selectedGroupId.value]
  try {
    await upstreamHubAPI.setEntryGroups(entry.account_id, nextGroupIds)
    await refreshEntries()
  } catch (err) {
    appStore.showError(errorMessage(err))
  }
}

async function toggleEntryEnabled(entry: UpstreamEntry) {
  const entryGroupRef = groupRef(entry)
  if (!entryGroupRef) return
  try {
    await upstreamHubAPI.setEntryGroupEnabled(entry.account_id, entryGroupRef.group_id, !entryGroupRef.enabled)
    await refreshEntries()
  } catch (err) {
    appStore.showError(errorMessage(err))
  }
}

const routeForm = reactive({
  public_model: '',
  target_platform: platformOptions.value[0]?.value || 'openai',
  upstream_model: '',
  priority: 100
})

async function addRoute() {
  if (!selectedGroupId.value) return
  if (!routeForm.public_model.trim()) {
    appStore.showError(t('admin.aggregation.routeRequired'))
    return
  }
  saving.value = true
  try {
    await groupsAPI.createCompositeRoute(selectedGroupId.value, {
      public_model: routeForm.public_model.trim(),
      match_type: 'exact',
      target_platform: routeForm.target_platform,
      upstream_model: routeForm.upstream_model.trim() || undefined,
      endpoint: 'any',
      priority: Number(routeForm.priority) || 100,
      enabled: true
    })
    routeForm.public_model = ''
    routeForm.upstream_model = ''
    appStore.showSuccess(t('admin.aggregation.routeCreated'))
    await loadGroupData()
  } catch (err) {
    appStore.showError(errorMessage(err))
  } finally {
    saving.value = false
  }
}

async function toggleRouteEnabled(route: CompositeModelRoute) {
  if (!selectedGroupId.value) return
  try {
    await groupsAPI.updateCompositeRoute(selectedGroupId.value, route.id, {
      public_model: route.public_model,
      match_type: route.match_type,
      target_platform: route.target_platform,
      upstream_model: route.upstream_model,
      endpoint: route.endpoint,
      priority: route.priority,
      enabled: !route.enabled,
      notes: route.notes
    })
    await loadGroupData()
  } catch (err) {
    appStore.showError(errorMessage(err))
  }
}

async function removeRoute(route: CompositeModelRoute) {
  if (!selectedGroupId.value) return
  try {
    await groupsAPI.deleteCompositeRoute(selectedGroupId.value, route.id)
    appStore.showSuccess(t('admin.aggregation.routeDeleted'))
    await loadGroupData()
  } catch (err) {
    appStore.showError(errorMessage(err))
  }
}

onMounted(loadAll)
</script>
