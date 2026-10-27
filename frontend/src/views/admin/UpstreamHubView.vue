<template>
  <AppLayout>
    <div class="w-full min-w-0 space-y-6 pb-8">
      <header
        class="page-header mb-0 rounded-3xl bg-white p-5 shadow-sm ring-1 ring-gray-900/5 dark:bg-dark-800 dark:ring-dark-700 sm:p-6"
      >
        <div class="flex flex-wrap items-start justify-between gap-3">
          <div>
            <h1 class="page-title flex items-center gap-2 text-xl font-black text-gray-900 dark:text-white">
              <span
                class="inline-flex h-8 w-8 items-center justify-center rounded-xl bg-primary-50 text-primary-600 dark:bg-primary-900/30 dark:text-primary-400"
              >
                <Icon name="cube" size="sm" />
              </span>
              {{ t('admin.upstreamHub.title') }}
            </h1>
            <p class="page-description mt-1.5 text-xs text-gray-500 dark:text-gray-400">
              {{ t('admin.upstreamHub.description') }}
            </p>
          </div>
          <button class="btn btn-primary" @click="openCreateVendor">
            {{ t('admin.upstreamHub.createVendor') }}
          </button>
        </div>
      </header>

      <div class="card">
        <div class="card-header flex items-center justify-between">
          <div>
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('admin.upstreamHub.vendors') }}</h2>
          </div>
          <button class="btn btn-secondary btn-sm" :disabled="loading" @click="load">{{ t('common.refresh') }}</button>
        </div>
        <div class="p-6">
          <div v-if="loading" class="flex items-center gap-2 text-sm text-gray-500">
            <LoadingSpinner size="sm" /> {{ t('common.loading') }}
          </div>
          <EmptyState
            v-else-if="vendors.length === 0"
            :title="t('admin.upstreamHub.empty')"
            :description="t('admin.upstreamHub.emptyHint')"
            :action-text="t('admin.upstreamHub.createVendor')"
            @action="openCreateVendor"
          />
          <div v-else class="table-container">
            <table class="table">
              <thead>
                <tr>
                  <th>{{ t('admin.upstreamHub.vendorName') }}</th>
                  <th>{{ t('admin.upstreamHub.budget') }}</th>
                  <th>{{ t('admin.upstreamHub.allocated') }}</th>
                  <th>{{ t('admin.upstreamHub.entriesCount') }}</th>
                  <th>{{ t('admin.upstreamHub.referencedGroups') }}</th>
                  <th>{{ t('common.status') }}</th>
                  <th class="text-right">{{ t('common.actions') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="vendor in vendors" :key="vendor.id">
                  <td>
                    <div class="font-medium text-gray-900 dark:text-white">{{ vendor.name }}</div>
                    <div v-if="vendor.notes" class="text-xs text-gray-500 dark:text-gray-400">{{ vendor.notes }}</div>
                  </td>
                  <td>{{ vendor.concurrency_budget || '—' }}</td>
                  <td>
                    <span :class="vendor.over_budget ? 'text-red-600 dark:text-red-400' : ''">
                      {{ vendor.allocated_concurrency }}
                    </span>
                    <span v-if="vendor.concurrency_budget > 0" class="ml-1 text-xs text-gray-500">
                      / {{ vendor.concurrency_budget }}
                    </span>
                  </td>
                  <td>{{ vendor.entry_count }}</td>
                  <td>{{ vendor.referenced_group_count }}</td>
                  <td>
                    <span v-if="vendor.over_budget" class="badge badge-danger">
                      {{ t('admin.upstreamHub.overBudget') }}
                    </span>
                    <span v-else class="badge badge-success">{{ t('common.active') }}</span>
                  </td>
                  <td class="text-right whitespace-nowrap">
                    <button
                      class="btn btn-sm"
                      :class="detail?.vendor.id === vendor.id ? 'btn-primary' : 'btn-secondary'"
                      @click="toggleDetail(vendor.id)"
                    >
                      {{ t('admin.upstreamHub.manage') }}
                    </button>
                    <button class="btn btn-sm btn-secondary ml-1" @click="openEditVendor(vendor)">
                      {{ t('admin.upstreamHub.edit') }}
                    </button>
                    <button class="btn btn-sm btn-ghost ml-1 text-red-600" @click="askDeleteVendor(vendor)">
                      {{ t('admin.upstreamHub.delete') }}
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- Vendor detail: entries -->
      <div v-if="detail" class="space-y-6">
          <div class="card">
            <div class="card-header flex flex-wrap items-center justify-between gap-2">
              <div>
                <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
                  {{ detail.vendor.name }}
                  <span class="badge badge-primary ml-1">
                    {{ detail.vendor.concurrency_budget > 0 ? detail.vendor.concurrency_budget : '∞' }}
                  </span>
                </h2>
                <p class="text-xs text-gray-500 dark:text-gray-400">
                  {{ t('admin.upstreamHub.budgetUsage', {
                    used: detail.vendor.allocated_concurrency,
                    total: detail.vendor.concurrency_budget || '∞'
                  }) }}
                </p>
              </div>
              <button class="btn btn-primary btn-sm" @click="openEntriesDialog">
                {{ t('admin.upstreamHub.addOrEdit') }}
              </button>
            </div>
            <div class="p-6">
              <div
                v-for="warning in detail.warnings"
                :key="warning"
                class="mb-4 rounded-lg border border-amber-200 bg-amber-50 p-3 text-sm text-amber-700 dark:border-amber-800 dark:bg-amber-900/20 dark:text-amber-300"
              >
                {{ warning }}
              </div>
              <EmptyState
                v-if="detail.entries.length === 0"
                :title="t('admin.upstreamHub.empty')"
                :description="t('admin.upstreamHub.emptyHint')"
                :action-text="t('admin.upstreamHub.addOrEdit')"
                @action="openEntriesDialog"
              />
              <div v-else class="table-container">
                <table class="table">
                  <thead>
                    <tr>
                      <th>{{ t('admin.upstreamHub.entryName') }}</th>
                      <th>{{ t('admin.upstreamHub.platform') }}</th>
                      <th>Base URL</th>
                      <th>{{ t('admin.upstreamHub.concurrency') }}</th>
                      <th>{{ t('admin.upstreamHub.priority') }}</th>
                      <th>{{ t('admin.upstreamHub.groups') }}</th>
                      <th>{{ t('admin.upstreamHub.schedulable') }}</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="entry in detail.entries" :key="entry.account_id">
                      <td>
                        <div class="font-medium text-gray-900 dark:text-white">{{ entry.name }}</div>
                        <div class="text-xs text-gray-500">{{ entry.status }}</div>
                      </td>
                      <td><span class="badge badge-primary">{{ entry.platform }}</span></td>
                      <td class="max-w-[16rem] truncate text-xs" :title="entry.base_url">{{ entry.base_url }}</td>
                      <td>{{ entry.concurrency }}</td>
                      <td>{{ entry.priority }}</td>
                      <td>
                        <span v-if="groupRefs(entry).length === 0" class="text-xs text-gray-400 dark:text-dark-400">
                          {{ t('admin.upstreamHub.ungrouped') }}
                        </span>
                        <div v-else class="flex flex-col items-start gap-1">
                          <span
                            v-for="group in groupRefs(entry)"
                            :key="group.group_id"
                            class="block w-full rounded-md px-2 py-0.5 text-xs"
                            :class="
                              group.enabled
                                ? 'bg-primary-50 text-primary-700 dark:bg-primary-900/30 dark:text-primary-400'
                                : 'bg-gray-100 text-gray-500 dark:bg-dark-700 dark:text-dark-300'
                            "
                            :title="group.enabled ? t('admin.upstreamHub.groupEnabled') : t('admin.upstreamHub.groupDisabled')"
                          >
                            {{ group.group_name }}
                          </span>
                        </div>
                      </td>
                      <td>
                        <Toggle :model-value="entry.schedulable" @update:model-value="(v) => toggleSchedulable(entry, v)" />
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </div>
      </div>
    </div>

    <!-- Vendor dialog -->
    <BaseDialog :show="showVendorDialog" :title="vendorForm.id ? t('admin.upstreamHub.editVendor') : t('admin.upstreamHub.createVendor')" @close="showVendorDialog = false">
      <div class="space-y-4">
        <div>
          <label class="input-label">{{ t('admin.upstreamHub.vendorName') }}</label>
          <input v-model="vendorForm.name" class="input" type="text" :placeholder="t('admin.upstreamHub.vendorNamePlaceholder')" />
        </div>
        <div>
          <label class="input-label">{{ t('admin.upstreamHub.budget') }}</label>
          <input v-model.number="vendorForm.concurrency_budget" class="input" type="number" min="0" />
          <p class="input-hint">{{ t('admin.upstreamHub.budgetHint') }}</p>
        </div>
        <div>
          <label class="input-label">{{ t('admin.upstreamHub.notes') }}</label>
          <input v-model="vendorForm.notes" class="input" type="text" :placeholder="t('admin.upstreamHub.notesPlaceholder')" />
        </div>
      </div>
      <template #footer>
        <button class="btn btn-secondary" @click="showVendorDialog = false">{{ t('admin.upstreamHub.cancel') }}</button>
        <button class="btn btn-primary" :disabled="saving" @click="saveVendor">{{ t('admin.upstreamHub.save') }}</button>
      </template>
    </BaseDialog>

    <!-- Entries dialog: add / edit / copy in one layout -->
    <BaseDialog :show="showEntryDialog" :title="t('admin.upstreamHub.addOrEdit')" width="max" @close="closeEntriesDialog">
      <div v-if="anyRowDropdownOpen" class="fixed inset-0 z-30" @click="closeAllDropdowns"></div>
      <div class="space-y-3">
        <div class="table-container">
          <table class="table">
            <thead>
              <tr>
                <th class="min-w-[9rem] whitespace-nowrap">{{ t('admin.upstreamHub.entryName') }}</th>
                <th class="min-w-[8rem] whitespace-nowrap">{{ t('admin.upstreamHub.platform') }}</th>
                <th class="min-w-[12rem] whitespace-nowrap">Base URL</th>
                <th class="min-w-[10rem] whitespace-nowrap">{{ t('admin.upstreamHub.apiKey') }}</th>
                <th class="min-w-[6rem] whitespace-nowrap">{{ t('admin.upstreamHub.concurrency') }}</th>
                <th class="min-w-[8.5rem] whitespace-nowrap">{{ t('admin.upstreamHub.priority') }}</th>
                <th class="min-w-[9rem] whitespace-nowrap">{{ t('admin.upstreamHub.groups') }}</th>
                <th class="whitespace-nowrap text-right">{{ t('common.actions') }}</th>
              </tr>
            </thead>
            <tbody>
              <template v-for="(row, index) in entryRows" :key="index">
                <tr>
                  <td>
                    <input v-model="row.name" class="input" type="text" :placeholder="t('admin.upstreamHub.entryNamePlaceholder')" />
                  </td>
                  <td>
                    <Select v-model="row.platform" :options="platformOptions" :disabled="!!row.accountId" />
                  </td>
                  <td>
                    <input v-model="row.base_url" class="input" type="text" placeholder="https://..." />
                  </td>
                  <td>
                    <input
                      v-model="row.api_key"
                      class="input"
                      type="password"
                      :placeholder="row.accountId ? t('admin.upstreamHub.apiKeyKeepPlaceholder') : 'sk-...'"
                    />
                  </td>
                  <td><input v-model.number="row.concurrency" class="input" type="number" min="1" /></td>
                  <td><input v-model.number="row.priority" class="input" type="number" /></td>
                  <td class="relative">
                    <button
                      type="button"
                      class="flex w-full items-center justify-between rounded-lg border border-gray-200 bg-white px-3 py-1.5 text-left text-sm dark:border-dark-600 dark:bg-dark-800"
                      @click="toggleRowDropdown(row, $event)"
                    >
                      <span :class="row.groupIds.length ? '' : 'text-gray-400 dark:text-dark-400'">
                        {{
                          row.groupIds.length
                            ? t('admin.upstreamHub.selectedGroups', { count: row.groupIds.length })
                            : t('admin.upstreamHub.ungrouped')
                        }}
                      </span>
                      <svg class="h-4 w-4 text-gray-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M6 9l6 6 6-6" />
                      </svg>
                    </button>
                    <div
                      v-if="row.expanded"
                      class="dropdown fixed z-40 w-80"
                      :style="{ top: row.dropdownTop + 'px', right: row.dropdownRight + 'px' }"
                    >
                      <div class="max-h-72 overflow-y-auto py-1">
                        <div
                          v-for="group in groupsForRow(row)"
                          :key="group.id"
                          class="dropdown-item w-full cursor-pointer justify-between"
                          @click="toggleRowGroup(row, group.id)"
                        >
                          <div class="flex min-w-0 items-center gap-2 text-sm">
                            <input
                              type="checkbox"
                              class="pointer-events-none h-4 w-4"
                              :checked="row.groupIds.includes(group.id)"
                            />
                            <span class="truncate">{{ group.name }}</span>
                            <span class="badge badge-primary">{{ group.platform }}</span>
                          </div>
                          <div
                            v-if="row.groupIds.includes(group.id)"
                            class="flex items-center gap-1.5"
                            @click.stop
                          >
                            <Toggle
                              :model-value="row.groupEnabled[group.id] !== false"
                              @update:model-value="(v) => (row.groupEnabled[group.id] = v)"
                            />
                          </div>
                        </div>
                        <div v-if="groupsForRow(row).length === 0" class="px-4 py-2 text-xs text-gray-500">
                          {{ t('admin.upstreamHub.noGroupsForPlatform') }}
                        </div>
                      </div>
                    </div>
                  </td>
                  <td class="text-right whitespace-nowrap">
                    <button class="btn btn-sm btn-secondary" @click="copyRow(row)">
                      {{ t('admin.upstreamHub.duplicate') }}
                    </button>
                    <button class="btn btn-sm btn-ghost ml-1 text-red-600" @click="removeRow(row, index)">
                      {{ t('admin.upstreamHub.delete') }}
                    </button>
                  </td>
                </tr>
              </template>
            </tbody>
          </table>
        </div>
        <div
          v-if="addRowsOverBudget"
          class="rounded-lg border border-red-200 bg-red-50 p-3 text-sm text-red-700 dark:border-red-800 dark:bg-red-900/20 dark:text-red-300"
        >
          {{ t('admin.upstreamHub.overBudgetBlockedInline', {
            used: effectiveAllocated,
            total: detail?.vendor.concurrency_budget || '∞',
            over: effectiveAllocated - (detail?.vendor.concurrency_budget || 0)
          }) }}
        </div>
        <div class="flex flex-wrap items-center justify-between gap-2">
          <button class="btn btn-secondary btn-sm" @click="addEntryRow">{{ t('admin.upstreamHub.addRow') }}</button>
          <p class="text-xs" :class="addRowsOverBudget ? 'text-red-600 dark:text-red-400' : 'text-gray-500'">
            {{ t('admin.upstreamHub.budgetUsage', { used: effectiveAllocated, total: detail?.vendor.concurrency_budget || '∞' }) }}
            <template v-if="addRowsOverBudget"> — {{ t('admin.upstreamHub.overBudgetHint') }}</template>
          </p>
        </div>
      </div>
      <template #footer>
        <button class="btn btn-secondary" @click="closeEntriesDialog">{{ t('admin.upstreamHub.cancel') }}</button>
        <button class="btn btn-primary" :disabled="saving || addRowsOverBudget" @click="saveEntryRows">
          {{ t('admin.upstreamHub.save') }}
        </button>
      </template>
    </BaseDialog>



    <ConfirmDialog
      :show="showDeleteConfirm"
      :title="t('admin.upstreamHub.delete')"
      :message="deleteMessage"
      :confirm-text="t('admin.upstreamHub.delete')"
      :cancel-text="t('admin.upstreamHub.cancel')"
      :danger="true"
      @confirm="confirmDelete"
      @cancel="showDeleteConfirm = false"
    />
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores/app'
import { extractApiErrorMessage } from '@/utils/apiError'
import upstreamHubAPI, {
  type UpstreamEntry,
  type UpstreamEntryGroupRef,
  type UpstreamGroupOption,
  type UpstreamVendor,
  type UpstreamVendorDetail
} from '@/api/admin/upstreamHub'
import { CONCRETE_PLATFORM_OPTIONS } from '@/constants/platforms'
import AppLayout from '@/components/layout/AppLayout.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import Select from '@/components/common/Select.vue'
import Toggle from '@/components/common/Toggle.vue'
import Icon from '@/components/icons/Icon.vue'

const { t } = useI18n()
const appStore = useAppStore()

const vendors = ref<UpstreamVendor[]>([])
const detail = ref<UpstreamVendorDetail | null>(null)
const loading = ref(false)
const saving = ref(false)

const platformOptions = computed(() => CONCRETE_PLATFORM_OPTIONS.map((p) => ({ label: p.label, value: p.value })))

const showVendorDialog = ref(false)
const vendorForm = reactive({ id: 0, name: '', concurrency_budget: 0, notes: '' })

interface EntryRow {
  accountId?: number
  name: string
  platform: string
  base_url: string
  api_key: string
  concurrency: number
  priority: number
  groupIds: number[]
  groupEnabled: Record<number, boolean>
  expanded: boolean
  dropdownTop: number
  dropdownRight: number
}
const showEntryDialog = ref(false)
const entryRows = ref<EntryRow[]>([])
const pendingDeletes = ref<Array<{ accountId: number; concurrency: number }>>([])
const groupOptions = ref<UpstreamGroupOption[]>([])

const showDeleteConfirm = ref(false)
const deleteTarget = ref<{ type: 'vendor'; id: number; name: string } | null>(null)
const deleteMessage = computed(() =>
  deleteTarget.value ? t('admin.upstreamHub.deleteVendorConfirm', { name: deleteTarget.value.name }) : ''
)

const effectiveAllocated = computed(() =>
  entryRows.value.reduce((sum, row) => sum + (Number(row.concurrency) || 0), 0)
)
const addRowsOverBudget = computed(() => {
  const budget = detail.value?.vendor.concurrency_budget ?? 0
  return budget > 0 && effectiveAllocated.value > budget
})

function groupRefs(entry: UpstreamEntry): UpstreamEntryGroupRef[] {
  return entry.groups || []
}

function entryToRow(entry: UpstreamEntry): EntryRow {
  const groupEnabled: Record<number, boolean> = {}
  for (const group of groupRefs(entry)) {
    groupEnabled[group.group_id] = group.enabled
  }
  return {
    accountId: entry.account_id,
    name: entry.name,
    platform: entry.platform,
    base_url: entry.base_url,
    api_key: '',
    concurrency: entry.concurrency,
    priority: entry.priority,
    groupIds: groupRefs(entry).map((group) => group.group_id),
    groupEnabled,
    expanded: false,
    dropdownTop: 0,
    dropdownRight: 0
  }
}

function errorMessage(err: unknown): string {
  return extractApiErrorMessage(err, t('admin.upstreamHub.saveFailed'))
}

async function load() {
  loading.value = true
  try {
    vendors.value = await upstreamHubAPI.listVendors()
  } catch (err) {
    appStore.showError(extractApiErrorMessage(err, t('admin.upstreamHub.loadFailed')))
  } finally {
    loading.value = false
  }
}

async function toggleDetail(vendorId: number) {
  if (detail.value?.vendor.id === vendorId) {
    detail.value = null
    return
  }
  try {
    detail.value = await upstreamHubAPI.getVendor(vendorId)
  } catch (err) {
    appStore.showError(extractApiErrorMessage(err, t('admin.upstreamHub.loadFailed')))
  }
}


// ---------- vendor ----------
function openCreateVendor() {
  vendorForm.id = 0
  vendorForm.name = ''
  vendorForm.concurrency_budget = 0
  vendorForm.notes = ''
  showVendorDialog.value = true
}

function openEditVendor(vendor: UpstreamVendor) {
  vendorForm.id = vendor.id
  vendorForm.name = vendor.name
  vendorForm.concurrency_budget = vendor.concurrency_budget
  vendorForm.notes = vendor.notes || ''
  showVendorDialog.value = true
}

async function saveVendor() {
  saving.value = true
  try {
    const payload = {
      name: vendorForm.name,
      concurrency_budget: Number(vendorForm.concurrency_budget) || 0,
      notes: vendorForm.notes || null
    }
    if (vendorForm.id > 0) {
      await upstreamHubAPI.updateVendor(vendorForm.id, payload)
      appStore.showSuccess(t('admin.upstreamHub.vendorUpdated'))
    } else {
      await upstreamHubAPI.createVendor(payload)
      appStore.showSuccess(t('admin.upstreamHub.vendorCreated'))
    }
    showVendorDialog.value = false
    await load()
    if (detail.value) await refreshDetailEntries()
  } catch (err) {
    appStore.showError(errorMessage(err))
  } finally {
    saving.value = false
  }
}

async function refreshDetailEntries() {
  if (!detail.value) return
  try {
    detail.value = await upstreamHubAPI.getVendor(detail.value.vendor.id)
  } catch (err) {
    appStore.showError(extractApiErrorMessage(err, t('admin.upstreamHub.loadFailed')))
  }
}

function askDeleteVendor(vendor: UpstreamVendor) {
  deleteTarget.value = { type: 'vendor', id: vendor.id, name: vendor.name }
  showDeleteConfirm.value = true
}

async function confirmDelete() {
  if (!deleteTarget.value) return
  try {
    await upstreamHubAPI.deleteVendor(deleteTarget.value.id)
    appStore.showSuccess(t('admin.upstreamHub.vendorDeleted'))
    if (detail.value?.vendor.id === deleteTarget.value.id) detail.value = null
    showDeleteConfirm.value = false
    await load()
  } catch (err) {
    appStore.showError(errorMessage(err))
  }
}

// ---------- entries dialog ----------
function blankRow(): EntryRow {
  return {
    name: '',
    platform: detail.value?.entries[0]?.platform || platformOptions.value[0]?.value || 'openai',
    base_url: '',
    api_key: '',
    concurrency: 1,
    priority: 50,
    groupIds: [],
    groupEnabled: {},
    expanded: false,
    dropdownTop: 0,
    dropdownRight: 0
  }
}

function addEntryRow() {
  entryRows.value.push(blankRow())
}

function groupsForRow(row: EntryRow): UpstreamGroupOption[] {
  return groupOptions.value.filter((group) => group.platform === row.platform)
}

const anyRowDropdownOpen = computed(() => entryRows.value.some((row) => row.expanded))

function closeAllDropdowns() {
  for (const item of entryRows.value) item.expanded = false
}

function toggleRowDropdown(row: EntryRow, event: Event) {
  const next = !row.expanded
  closeAllDropdowns()
  if (next) {
    const rect = (event.currentTarget as HTMLElement).getBoundingClientRect()
    row.dropdownTop = rect.bottom + 4
    row.dropdownRight = Math.max(8, window.innerWidth - rect.right)
  }
  row.expanded = next
}

function toggleRowGroup(row: EntryRow, groupId: number) {
  const index = row.groupIds.indexOf(groupId)
  if (index >= 0) {
    row.groupIds.splice(index, 1)
    delete row.groupEnabled[groupId]
  } else {
    row.groupIds.push(groupId)
    row.groupEnabled[groupId] = true
  }
}

async function applyRowGroups(accountId: number, row: EntryRow) {
  await upstreamHubAPI.setEntryGroups(accountId, row.groupIds)
  for (const groupId of row.groupIds) {
    if (row.groupEnabled[groupId] === false) {
      await upstreamHubAPI.setEntryGroupEnabled(accountId, groupId, false)
    }
  }
}

async function openEntriesDialog() {
  if (!detail.value) return
  pendingDeletes.value = []
  entryRows.value = detail.value.entries.map(entryToRow)
  if (entryRows.value.length === 0) addEntryRow()
  try {
    groupOptions.value = await upstreamHubAPI.listGroups()
  } catch (err) {
    appStore.showError(extractApiErrorMessage(err, t('admin.upstreamHub.loadFailed')))
  }
  showEntryDialog.value = true
}

function closeEntriesDialog() {
  showEntryDialog.value = false
  entryRows.value = []
  pendingDeletes.value = []
  groupOptions.value = []
}

async function copyRow(row: EntryRow) {
  if (row.accountId) {
    // 已存在条目：服务端复制凭据（API Key 不回显）
    try {
      const created = await upstreamHubAPI.duplicateEntry(row.accountId, {
        name: `${row.name}-copy`,
        platform: row.platform,
        concurrency: Number(row.concurrency) || 1,
        priority: Number(row.priority) || 50
      })
      const copied = entryToRow(created)
      copied.dropdownTop = 0
      copied.dropdownRight = 0
      copied.groupIds = [...row.groupIds]
      copied.groupEnabled = { ...row.groupEnabled }
      entryRows.value.push(copied)
      appStore.showSuccess(t('admin.upstreamHub.duplicateSuccess'))
    } catch (err) {
      appStore.showError(errorMessage(err))
    }
    return
  }
  entryRows.value.push({
    ...row,
    accountId: undefined,
    name: row.name ? `${row.name}-copy` : row.name,
    groupIds: [...row.groupIds],
    groupEnabled: { ...row.groupEnabled },
    expanded: false,
    dropdownTop: 0,
    dropdownRight: 0
  })
}

function removeRow(row: EntryRow, index: number) {
  if (row.accountId) {
    pendingDeletes.value.push({ accountId: row.accountId, concurrency: row.concurrency })
  }
  entryRows.value.splice(index, 1)
}

async function saveEntryRows() {
  if (!detail.value) return
  const budget = detail.value.vendor.concurrency_budget
  if (budget > 0 && effectiveAllocated.value > budget) {
    appStore.showError(
      t('admin.upstreamHub.overBudgetBlocked', { used: effectiveAllocated.value, total: budget })
    )
    return
  }
  saving.value = true
  try {
    for (const row of entryRows.value) {
      if (!row.name.trim() || !row.base_url.trim()) {
        throw new Error(t('admin.upstreamHub.fillRequiredFields'))
      }
      if (row.accountId) {
        await upstreamHubAPI.updateEntry(row.accountId, {
          name: row.name,
          base_url: row.base_url,
          api_key: row.api_key,
          concurrency: Number(row.concurrency) || 1,
          priority: Number(row.priority) || 50
        })
        await applyRowGroups(row.accountId, row)
      } else {
        if (!row.api_key.trim()) {
          throw new Error(t('admin.upstreamHub.fillRequiredFields'))
        }
        const created = await upstreamHubAPI.createEntry(detail.value.vendor.id, {
          name: row.name,
          platform: row.platform,
          base_url: row.base_url,
          api_key: row.api_key,
          concurrency: Number(row.concurrency) || 1,
          priority: Number(row.priority) || 50
        })
        await applyRowGroups(created.account_id, row)
      }
    }
    for (const removed of pendingDeletes.value) {
      await upstreamHubAPI.deleteEntry(removed.accountId)
    }
    appStore.showSuccess(t('admin.upstreamHub.entryUpdated'))
    closeEntriesDialog()
    await refreshDetailEntries()
    await load()
  } catch (err) {
    appStore.showError(errorMessage(err))
  } finally {
    saving.value = false
  }
}

// ---------- entry row actions ----------
async function toggleSchedulable(entry: UpstreamEntry, value: boolean) {
  try {
    const updated = await upstreamHubAPI.setEntrySchedulable(entry.account_id, value)
    entry.schedulable = updated.schedulable
  } catch (err) {
    appStore.showError(errorMessage(err))
  }
}




onMounted(load)
</script>
