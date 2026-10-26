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
      <Transition
        enter-active-class="transition duration-200 ease-out"
        enter-from-class="opacity-0 -translate-y-1"
        leave-active-class="transition duration-150 ease-in"
        leave-to-class="opacity-0 -translate-y-1"
      >
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
                        <button
                          v-if="groupRefs(entry).length === 0"
                          class="badge badge-gray"
                          @click="openGroups(entry)"
                        >
                          {{ t('admin.upstreamHub.noGroups') }}
                        </button>
                        <div v-else class="flex flex-wrap gap-1">
                          <button
                            v-for="group in groupRefs(entry)"
                            :key="group.group_id"
                            class="badge"
                            :class="group.enabled ? 'badge-success' : 'badge-gray'"
                            :title="group.enabled ? t('admin.upstreamHub.groupEnabled') : t('admin.upstreamHub.groupDisabled')"
                            @click="toggleGroupEnabled(entry, group)"
                          >
                            {{ group.group_name }}
                          </button>
                          <button class="badge badge-primary" @click="openGroups(entry)">
                            {{ t('admin.upstreamHub.manageGroups') }}
                          </button>
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
      </Transition>
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
    <BaseDialog :show="showEntryDialog" :title="t('admin.upstreamHub.addOrEdit')" width="extra-wide" @close="closeEntriesDialog">
      <div class="space-y-3">
        <div class="table-container">
          <table class="table">
            <thead>
              <tr>
                <th>{{ t('admin.upstreamHub.entryName') }}</th>
                <th>{{ t('admin.upstreamHub.platform') }}</th>
                <th>Base URL</th>
                <th>{{ t('admin.upstreamHub.apiKey') }}</th>
                <th>{{ t('admin.upstreamHub.concurrency') }}</th>
                <th>{{ t('admin.upstreamHub.priority') }}</th>
                <th class="text-right">{{ t('common.actions') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(row, index) in entryRows" :key="index">
                <td><input v-model="row.name" class="input" type="text" :placeholder="t('admin.upstreamHub.entryNamePlaceholder')" /></td>
                <td>
                  <Select v-model="row.platform" :options="platformOptions" :disabled="!!row.accountId" />
                </td>
                <td><input v-model="row.base_url" class="input" type="text" placeholder="https://..." /></td>
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
                <td class="text-right whitespace-nowrap">
                  <button class="btn btn-sm btn-secondary" @click="copyRow(row)">{{ t('admin.upstreamHub.duplicate') }}</button>
                  <button class="btn btn-sm btn-ghost ml-1 text-red-600" @click="removeRow(row, index)">
                    {{ t('admin.upstreamHub.delete') }}
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
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
        <button class="btn btn-primary" :disabled="saving" @click="saveEntryRows">{{ t('admin.upstreamHub.save') }}</button>
      </template>
    </BaseDialog>

    <!-- Groups dialog -->
    <BaseDialog :show="showGroupsDialog" :title="t('admin.upstreamHub.selectGroups')" @close="showGroupsDialog = false">
      <p class="mb-3 text-xs text-gray-500 dark:text-gray-400">{{ t('admin.upstreamHub.selectGroupsHint') }}</p>
      <div v-if="groupOptions.length === 0" class="text-sm text-gray-500">
        {{ t('admin.upstreamHub.noGroupsForPlatform') }}
      </div>
      <div v-else class="space-y-2">
        <label
          v-for="group in groupOptions"
          :key="group.id"
          class="flex items-center justify-between rounded-lg border border-gray-100 px-3 py-2 dark:border-dark-700"
        >
          <span class="flex items-center gap-2 text-sm">
            <input v-model="selectedGroupIDs" type="checkbox" :value="group.id" class="h-4 w-4" />
            {{ group.name }}
            <span class="badge badge-primary">{{ group.platform }}</span>
          </span>
        </label>
      </div>
      <template #footer>
        <button class="btn btn-secondary" @click="showGroupsDialog = false">{{ t('admin.upstreamHub.cancel') }}</button>
        <button class="btn btn-primary" :disabled="saving" @click="saveGroups">{{ t('admin.upstreamHub.save') }}</button>
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
}
const showEntryDialog = ref(false)
const entryRows = ref<EntryRow[]>([])
const pendingDeletes = ref<Array<{ accountId: number; concurrency: number }>>([])

const showGroupsDialog = ref(false)
const groupOptions = ref<UpstreamGroupOption[]>([])
const selectedGroupIDs = ref<number[]>([])
const groupsEntry = ref<UpstreamEntry | null>(null)

const showDeleteConfirm = ref(false)
const deleteTarget = ref<{ type: 'vendor'; id: number; name: string } | null>(null)
const deleteMessage = computed(() =>
  deleteTarget.value ? t('admin.upstreamHub.deleteVendorConfirm', { name: deleteTarget.value.name }) : ''
)

const effectiveAllocated = computed(() => {
  const allocated = detail.value?.vendor.allocated_concurrency ?? 0
  const drafts = entryRows.value.reduce(
    (sum, row) => sum + (row.accountId ? 0 : Number(row.concurrency) || 0),
    0
  )
  const removed = pendingDeletes.value.reduce((sum, item) => sum + item.concurrency, 0)
  return Math.max(0, allocated + drafts - removed)
})
const addRowsOverBudget = computed(() => {
  const budget = detail.value?.vendor.concurrency_budget ?? 0
  return budget > 0 && effectiveAllocated.value > budget
})

function groupRefs(entry: UpstreamEntry): UpstreamEntryGroupRef[] {
  return entry.groups || []
}

function entryToRow(entry: UpstreamEntry): EntryRow {
  return {
    accountId: entry.account_id,
    name: entry.name,
    platform: entry.platform,
    base_url: entry.base_url,
    api_key: '',
    concurrency: entry.concurrency,
    priority: entry.priority
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
function addEntryRow() {
  entryRows.value.push({
    name: '',
    platform: detail.value?.entries[0]?.platform || platformOptions.value[0]?.value || 'openai',
    base_url: '',
    api_key: '',
    concurrency: 1,
    priority: 50
  })
}

function openEntriesDialog() {
  if (!detail.value) return
  pendingDeletes.value = []
  entryRows.value = detail.value.entries.map(entryToRow)
  if (entryRows.value.length === 0) addEntryRow()
  showEntryDialog.value = true
}

function closeEntriesDialog() {
  showEntryDialog.value = false
  entryRows.value = []
  pendingDeletes.value = []
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
      entryRows.value.push(entryToRow(created))
      appStore.showSuccess(t('admin.upstreamHub.duplicateSuccess'))
    } catch (err) {
      appStore.showError(errorMessage(err))
    }
    return
  }
  entryRows.value.push({ ...row, name: row.name ? `${row.name}-copy` : row.name })
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
      } else {
        if (!row.api_key.trim()) {
          throw new Error(t('admin.upstreamHub.fillRequiredFields'))
        }
        await upstreamHubAPI.createEntry(detail.value.vendor.id, {
          name: row.name,
          platform: row.platform,
          base_url: row.base_url,
          api_key: row.api_key,
          concurrency: Number(row.concurrency) || 1,
          priority: Number(row.priority) || 50
        })
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

async function toggleGroupEnabled(entry: UpstreamEntry, group: UpstreamEntryGroupRef) {
  try {
    await upstreamHubAPI.setEntryGroupEnabled(entry.account_id, group.group_id, !group.enabled)
    group.enabled = !group.enabled
  } catch (err) {
    appStore.showError(errorMessage(err))
  }
}

async function openGroups(entry: UpstreamEntry) {
  groupsEntry.value = entry
  selectedGroupIDs.value = groupRefs(entry).map((group) => group.group_id)
  try {
    const groups = await upstreamHubAPI.listGroups()
    groupOptions.value = groups.filter((group) => group.platform === entry.platform)
    showGroupsDialog.value = true
  } catch (err) {
    appStore.showError(extractApiErrorMessage(err, t('admin.upstreamHub.loadFailed')))
  }
}

async function saveGroups() {
  if (!groupsEntry.value) return
  saving.value = true
  try {
    const previous = groupRefs(groupsEntry.value)
    await upstreamHubAPI.setEntryGroups(groupsEntry.value.account_id, selectedGroupIDs.value)
    for (const group of previous) {
      if (!group.enabled && selectedGroupIDs.value.includes(group.group_id)) {
        await upstreamHubAPI.setEntryGroupEnabled(groupsEntry.value.account_id, group.group_id, false)
      }
    }
    appStore.showSuccess(t('admin.upstreamHub.entryUpdated'))
    showGroupsDialog.value = false
    await refreshDetailEntries()
    if (detail.value) {
      const refreshed = detail.value.entries.find((item) => item.account_id === groupsEntry.value?.account_id)
      groupsEntry.value = refreshed || null
    }
  } catch (err) {
    appStore.showError(errorMessage(err))
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>
