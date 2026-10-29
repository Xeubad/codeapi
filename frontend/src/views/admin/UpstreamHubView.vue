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
              <div class="flex items-center gap-2">
                <button class="btn btn-secondary btn-sm" @click="openBaseURLDialog">
                  {{ t('admin.upstreamHub.manageBaseURLs') }}
                  <span
                    class="badge ml-1"
                    :class="vendorBaseURLs.length ? 'badge-primary' : 'badge-gray text-amber-600 dark:text-amber-400'"
                  >{{ vendorBaseURLs.length }}</span>
                </button>
                <button class="btn btn-primary btn-sm" @click="openEntriesDialog">
                  {{ t('admin.upstreamHub.addOrEdit') }}
                </button>
              </div>
            </div>
            <div class="p-6">
              <div
                v-if="vendorBaseURLs.length === 0"
                class="mb-4 rounded-lg border border-amber-200 bg-amber-50 p-3 text-sm text-amber-700 dark:border-amber-800 dark:bg-amber-900/20 dark:text-amber-300"
              >
                {{ t('admin.upstreamHub.baseURLsMissingBanner') }}
              </div>
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
                      <th>{{ t('admin.upstreamHub.models') }}</th>
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
                        <span
                          class="text-sm"
                          :class="
                            (entry.models?.length || 0) > 0
                              ? 'text-gray-700 dark:text-gray-300'
                              : 'text-amber-600 dark:text-amber-400'
                          "
                        >
                          {{ t('admin.upstreamHub.modelsRegistered', { count: entry.models?.length || 0 }) }}
                        </span>
                      </td>
                      <td>
                        <span v-if="groupRefs(entry).length === 0" class="text-xs text-gray-400 dark:text-dark-400">
                          {{ t('admin.upstreamHub.ungrouped') }}
                        </span>
                        <div v-else class="flex flex-col items-stretch gap-1">
                          <!-- 只读状态：分组挂载在添加/编辑弹窗维护，「本组暂不调用」在聚合调度页维护 -->
                          <span
                            v-for="group in groupRefs(entry)"
                            :key="group.group_id"
                            class="block rounded-md px-2 py-0.5 text-xs"
                            :class="
                              group.enabled
                                ? 'bg-primary-50 text-primary-700 dark:bg-primary-900/30 dark:text-primary-400'
                                : 'bg-gray-100 text-gray-500 dark:bg-dark-700 dark:text-dark-300'
                            "
                            :title="group.enabled ? t('admin.upstreamHub.groupEnabled') : t('admin.upstreamHub.groupDisabled')"
                          >
                            {{ group.group_name }}
                            <span v-if="!group.enabled" class="ml-1 opacity-80">· {{ t('admin.upstreamHub.groupDisabled') }}</span>
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
                <th class="min-w-[10rem] whitespace-nowrap">{{ t('admin.upstreamHub.apiKey') }}</th>
                <th class="min-w-[6rem] whitespace-nowrap">{{ t('admin.upstreamHub.concurrency') }}</th>
                <th class="min-w-[8.5rem] whitespace-nowrap">{{ t('admin.upstreamHub.priority') }}</th>
                <th class="min-w-[9.5rem] whitespace-nowrap">{{ t('admin.upstreamHub.models') }}</th>
                <th class="min-w-[9rem] whitespace-nowrap">{{ t('admin.upstreamHub.groups') }}</th>
                <th class="whitespace-nowrap text-right">{{ t('common.actions') }}</th>
              </tr>
            </thead>
            <tbody>
              <template v-for="(row, index) in entryRows" :key="index">
                <tr>
                  <td>
                    <input v-model="row.name" class="input" type="text" :placeholder="t('admin.upstreamHub.entryNamePlaceholder')" />
                    <div v-if="row.base_url" class="mt-0.5 truncate text-xs text-gray-500 dark:text-gray-400" :title="row.base_url">
                      {{ row.base_url }}
                    </div>
                  </td>
                  <td>
                    <Select v-model="row.platform" :options="platformOptions" :disabled="!!row.accountId" />
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
                  <td>
                    <span
                      class="text-sm"
                      :class="row.models.length ? 'text-gray-700 dark:text-gray-300' : 'text-amber-600 dark:text-amber-400'"
                    >
                      {{ t('admin.upstreamHub.modelsRegistered', { count: row.models.length }) }}
                    </span>
                  </td>
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
                          class="dropdown-item w-full cursor-pointer"
                          :class="
                            row.groupIds.includes(group.id)
                              ? 'bg-primary-50 text-primary-700 dark:bg-primary-900/20 dark:text-primary-300'
                              : ''
                          "
                          @click="toggleRowGroup(row, group.id)"
                        >
                          <span class="truncate">{{ group.name }}</span>
                          <span v-if="group.platform === 'composite'" class="badge badge-primary">
                            {{ groupPlatformLabel(group) }}
                          </span>
                          <span v-else class="badge badge-primary">{{ group.platform }}</span>
                        </div>
                        <div v-if="groupsForRow(row).length === 0" class="px-4 py-2 text-xs text-gray-500">
                          {{ t('admin.upstreamHub.noGroupsForPlatform') }}
                        </div>
                      </div>
                    </div>
                  </td>
                  <td class="text-right whitespace-nowrap">
                    <button
                      class="btn btn-sm btn-secondary"
                      :disabled="row.probeLoading"
                      @click="toggleModelDropdown(row, $event)"
                    >
                      {{ row.probeLoading ? t('admin.upstreamHub.probing') : t('admin.upstreamHub.probeModels') }}
                    </button>
                    <div
                      v-if="row.modelDropdownOpen"
                      class="dropdown fixed z-40 w-[30rem] text-left"
                      :style="{ top: row.modelDropdownTop + 'px', right: row.modelDropdownRight + 'px' }"
                    >
                      <div class="space-y-2 p-2">
                        <div class="flex flex-wrap items-center gap-2">
                          <span
                            v-if="row.probeSource"
                            class="badge"
                            :class="row.probeSource === 'models.dev' ? 'badge-gray' : 'badge-primary'"
                          >
                            {{
                              row.probeSource === 'models.dev'
                                ? t('admin.upstreamHub.probeSourceCatalog')
                                : t('admin.upstreamHub.probeSourceUpstream')
                            }}
                          </span>
                          <span class="text-xs text-gray-500 dark:text-dark-300">
                            {{ t('admin.upstreamHub.modelsSelected', { count: row.models.length }) }}
                          </span>
                          <div v-if="row.probeModels.length" class="ml-auto flex items-center gap-3">
                            <button class="text-xs text-primary-600 hover:underline" @click="selectAllProbeModels(row)">
                              {{ t('admin.upstreamHub.selectAll') }}
                            </button>
                            <button class="text-xs text-gray-500 hover:underline" @click="clearRowModels(row)">
                              {{ t('admin.upstreamHub.clearAll') }}
                            </button>
                          </div>
                        </div>
                        <p
                          v-if="row.probeError"
                          class="whitespace-normal break-all rounded-md bg-red-50 px-2 py-1 text-xs text-red-700 dark:bg-red-900/20 dark:text-red-300"
                        >
                          {{ row.probeError }}
                        </p>
                        <p
                          v-if="row.probeWarning"
                          class="whitespace-normal break-all rounded-md bg-amber-50 px-2 py-1 text-xs text-amber-700 dark:bg-amber-900/20 dark:text-amber-300"
                        >
                          {{ row.probeWarning }}
                        </p>
                        <div v-if="row.probeModels.length" class="flex max-h-52 flex-wrap gap-1.5 overflow-y-auto">
                          <button
                            v-for="model in row.probeModels"
                            :key="model"
                            type="button"
                            class="rounded-md px-2 py-1 font-mono text-xs transition-colors"
                            :class="
                              row.models.includes(model)
                                ? 'bg-primary-500 text-white'
                                : 'bg-gray-100 text-gray-700 hover:bg-gray-200 dark:bg-dark-700 dark:text-dark-200 dark:hover:bg-dark-600'
                            "
                            @click="toggleRowModel(row, model)"
                          >
                            {{ model }}
                          </button>
                        </div>
                        <p v-else-if="!row.probeError" class="text-xs text-gray-400 dark:text-dark-400">
                          {{ t('admin.upstreamHub.probeHint') }}
                        </p>
                        <div class="flex items-center gap-2 border-t border-gray-200 pt-2 dark:border-dark-600">
                          <input
                            v-model="row.manualModel"
                            class="input flex-1 font-mono text-xs"
                            type="text"
                            :placeholder="t('admin.upstreamHub.manualModelPlaceholder')"
                            @keyup.enter="addManualModel(row)"
                          />
                          <button
                            class="btn btn-secondary btn-sm shrink-0"
                            :disabled="!row.manualModel.trim()"
                            @click="addManualModel(row)"
                          >
                            {{ t('admin.upstreamHub.addModel') }}
                          </button>
                        </div>
                      </div>
                    </div>
                    <button class="btn btn-sm btn-secondary ml-1" @click="copyRow(row)">
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



    <!-- Base URL pool dialog：只维护地址与连通性，不涉及模型 -->
    <BaseDialog
      :show="showBaseURLDialog"
      :title="t('admin.upstreamHub.baseURLDialogTitle')"
      @close="showBaseURLDialog = false"
    >
      <div class="space-y-3">
        <div class="flex flex-wrap items-center gap-2">
          <button class="btn btn-primary btn-sm" @click="addBaseURLRow">
            {{ t('admin.upstreamHub.addBaseURL') }}
          </button>
          <span class="text-xs text-gray-500 dark:text-gray-400">{{ t('admin.upstreamHub.measureProxy') }}</span>
          <Select v-model="baseURLForm.proxyId" :options="proxyOptions" class="w-44" />
          <button
            class="btn btn-primary btn-sm"
            :disabled="baseURLTesting || !canMeasureBaseURLs"
            @click="testBaseURLs"
          >
            {{ baseURLTesting ? t('admin.upstreamHub.measuring') : t('admin.upstreamHub.measureBaseURLs') }}
          </button>
          <span v-if="baseURLTestError" class="whitespace-normal break-all text-xs text-red-600 dark:text-red-400">
            {{ baseURLTestError }}
          </span>
        </div>
        <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('admin.upstreamHub.baseURLDialogHint') }}</p>
        <p class="text-xs text-gray-400 dark:text-dark-400">{{ t('admin.upstreamHub.measureProxyHint') }}</p>

        <div
          v-if="baseURLForm.urls.length === 0"
          class="rounded-lg border border-dashed border-gray-300 p-6 text-center dark:border-dark-600"
        >
          <p class="text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.upstreamHub.baseURLsEmptyTitle') }}</p>
          <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ t('admin.upstreamHub.baseURLsEmptyDesc') }}</p>
        </div>
        <div v-else class="space-y-2">
          <div
            v-for="(_, index) in baseURLForm.urls"
            :key="index"
            class="flex cursor-pointer items-center gap-2 rounded-lg border px-2 py-1.5"
            :class="
              isInvalidBaseURLRow(baseURLForm.urls[index])
                ? 'border-red-300 dark:border-red-700'
                : isSelectedBaseURL(baseURLForm.urls[index])
                  ? 'border-primary-400 bg-primary-50/50 dark:border-primary-700 dark:bg-primary-900/10'
                  : 'border-gray-200 dark:border-dark-600'
            "
            @click="selectBaseURL(baseURLForm.urls[index])"
          >
            <span
              class="flex h-4 w-4 shrink-0 items-center justify-center rounded-full border"
              :class="
                isSelectedBaseURL(baseURLForm.urls[index])
                  ? 'border-primary-500'
                  : 'border-gray-300 dark:border-dark-500'
              "
            >
              <span
                v-if="isSelectedBaseURL(baseURLForm.urls[index])"
                class="h-2 w-2 rounded-full bg-primary-500"
              />
            </span>
            <input
              v-model="baseURLForm.urls[index]"
              class="input flex-1 font-mono text-xs"
              type="text"
              placeholder="https://api.example.com/v1"
              @click.stop
            />
            <span
              v-if="isInvalidBaseURLRow(baseURLForm.urls[index])"
              class="w-28 shrink-0 text-right text-xs text-red-600 dark:text-red-400"
            >
              {{ t('admin.upstreamHub.baseURLInvalidShort') }}
            </span>
            <span
              v-else-if="baseURLResultFor(baseURLForm.urls[index])"
              class="w-28 shrink-0 cursor-default text-right text-xs"
              :class="
                baseURLResultFor(baseURLForm.urls[index])?.ok
                  ? 'text-primary-600 dark:text-primary-400'
                  : 'text-red-600 dark:text-red-400'
              "
              :title="baseURLResultFor(baseURLForm.urls[index])?.error || baseURLResultFor(baseURLForm.urls[index])?.url"
            >
              {{
                baseURLResultFor(baseURLForm.urls[index])?.ok
                  ? `${baseURLResultFor(baseURLForm.urls[index])?.latency_ms}ms`
                  : (baseURLResultFor(baseURLForm.urls[index])?.status_code
                      ? `HTTP ${baseURLResultFor(baseURLForm.urls[index])?.status_code}`
                      : t('admin.upstreamHub.probeFailed'))
              }}
            </span>
            <span v-else class="w-28 shrink-0" />
            <button
              class="btn btn-sm btn-ghost shrink-0 text-red-600"
              @click.stop="removeBaseURLRow(index)"
            >
              {{ t('admin.upstreamHub.delete') }}
            </button>
          </div>
        </div>
      </div>
      <template #footer>
        <button class="btn btn-secondary" @click="showBaseURLDialog = false">{{ t('admin.upstreamHub.cancel') }}</button>
        <button class="btn btn-primary" :disabled="baseURLSaving || !canSaveBaseURLs" @click="saveBaseURLs">
          {{ t('admin.upstreamHub.saveBaseURLs') }}
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
import { syncUpstreamModelsPreview } from '@/api/admin/accounts'
import proxiesAPI from '@/api/admin/proxies'
import type { Proxy } from '@/types'
import upstreamHubAPI, {
  type UpstreamBaseURLLatency,
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
  // 已登记的模型白名单（model_mapping 的键），保存时若有改动才回写
  models: string[]
  modelsDirty: boolean
  modelDropdownOpen: boolean
  modelDropdownTop: number
  modelDropdownRight: number
  probeLoading: boolean
  probeError: string
  probeModels: string[]
  probeSource: string
  probeWarning: string
  manualModel: string
}
const showEntryDialog = ref(false)
const entryRows = ref<EntryRow[]>([])
const pendingDeletes = ref<Array<{ accountId: number; concurrency: number }>>([])
const groupOptions = ref<UpstreamGroupOption[]>([])
const proxies = ref<Proxy[]>([])
const proxyOptions = computed(() => [
  { label: t('admin.upstreamHub.proxyDirect'), value: 0 },
  ...proxies.value.map((proxy) => ({ label: proxy.name, value: proxy.id }))
])

const showDeleteConfirm = ref(false)
const deleteTarget = ref<{ type: 'vendor'; id: number; name: string } | null>(null)
const deleteMessage = computed(() =>
  deleteTarget.value ? t('admin.upstreamHub.deleteVendorConfirm', { name: deleteTarget.value.name }) : ''
)

const effectiveAllocated = computed(() =>
  entryRows.value.reduce((sum, row) => sum + (Number(row.concurrency) || 0), 0)
)
const vendorBaseURLs = computed(() => detail.value?.vendor.base_urls || [])
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
    dropdownRight: 0,
    models: [...(entry.models || [])],
    modelsDirty: false,
    modelDropdownOpen: false,
    modelDropdownTop: 0,
    modelDropdownRight: 0,
    probeLoading: false,
    probeError: '',
    probeModels: [],
    probeSource: '',
    probeWarning: '',
    manualModel: ''
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
    dropdownRight: 0,
    models: [],
    modelsDirty: false,
    modelDropdownOpen: false,
    modelDropdownTop: 0,
    modelDropdownRight: 0,
    probeLoading: false,
    probeError: '',
    probeModels: [],
    probeSource: '',
    probeWarning: '',
    manualModel: ''
  }
}

function addEntryRow() {
  entryRows.value.push(blankRow())
}

function groupsForRow(row: EntryRow): UpstreamGroupOption[] {
  // 综合分组天生跨平台，任何条目都可以挂；具体平台分组仍要求平台一致。
  return groupOptions.value.filter(
    (group) => group.platform === row.platform || group.platform === 'composite'
  )
}

function groupPlatformLabel(group: UpstreamGroupOption): string {
  return group.platform === 'composite' ? t('admin.upstreamHub.compositeGroup') : group.platform
}

const anyRowDropdownOpen = computed(() =>
  entryRows.value.some((row) => row.expanded || row.modelDropdownOpen)
)

function closeAllDropdowns() {
  for (const item of entryRows.value) {
    item.expanded = false
    item.modelDropdownOpen = false
  }
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

// 探测按钮：打开下拉并立即探测（每次点击都重新拉一次上游列表）。
async function toggleModelDropdown(row: EntryRow, event: Event) {
  const willOpen = !row.modelDropdownOpen
  closeAllDropdowns()
  if (!willOpen) return
  const rect = (event.currentTarget as HTMLElement).getBoundingClientRect()
  row.modelDropdownTop = rect.bottom + 4
  row.modelDropdownRight = Math.max(8, window.innerWidth - rect.right)
  row.modelDropdownOpen = true
  await probeRowModels(row)
}

function toggleRowModel(row: EntryRow, model: string) {
  const index = row.models.indexOf(model)
  if (index >= 0) {
    row.models.splice(index, 1)
  } else {
    row.models.push(model)
  }
  row.modelsDirty = true
}

function selectAllProbeModels(row: EntryRow) {
  row.models = [...new Set([...row.models, ...row.probeModels])].sort()
  row.modelsDirty = true
}

function clearRowModels(row: EntryRow) {
  row.models = []
  row.modelsDirty = true
}

async function probeRowModels(row: EntryRow) {
  row.probeError = ''
  row.probeSource = ''
  row.probeWarning = ''
  row.probeLoading = true
  try {
    let models: string[] = []
    if (row.accountId) {
      const probe = await upstreamHubAPI.probeEntryModels(row.accountId)
      models = probe.models || []
      row.probeSource = probe.source || 'upstream'
      row.probeWarning = probe.warning || ''
    } else {
      if (!row.base_url.trim() || !row.api_key.trim()) {
        row.probeError = t('admin.upstreamHub.probeNeedsCredentials')
        return
      }
      const result = await syncUpstreamModelsPreview({
        platform: row.platform,
        type: 'apikey',
        base_url: row.base_url.trim(),
        api_key: row.api_key.trim()
      })
      models = (result.models || []).map((model) => model.trim()).filter(Boolean)
      row.probeSource = 'upstream'
    }
    row.probeModels = [...new Set(models)].sort()
    if (row.probeModels.length === 0) {
      row.probeError = t('admin.upstreamHub.probeEmpty')
    }
  } catch (err) {
    row.probeError = extractApiErrorMessage(err, t('admin.upstreamHub.probeFailed'))
  } finally {
    row.probeLoading = false
  }
}

// 手动补录：上游没有模型接口或列表缺项时（cc-switch 同款），直接写模型 ID。
function addManualModel(row: EntryRow) {
  const value = row.manualModel.trim()
  if (!value) return
  if (!row.models.includes(value)) {
    row.models.push(value)
    row.models.sort()
    row.modelsDirty = true
  }
  if (!row.probeModels.includes(value)) {
    row.probeModels.push(value)
    row.probeModels.sort()
  }
  row.manualModel = ''
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
    models: [...row.models],
    modelsDirty: row.models.length > 0,
    modelDropdownOpen: false,
    modelDropdownTop: 0,
    modelDropdownRight: 0,
    probeLoading: false,
    probeError: '',
    probeModels: [],
    probeSource: '',
    probeWarning: '',
    manualModel: '',
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

// ---------- 接入地址池（厂家级，添加/编辑页内维护） ----------
const showBaseURLDialog = ref(false)
const baseURLForm = reactive({ urls: [] as string[], proxyId: 0, selectedUrl: '' })
const baseURLTesting = ref(false)
const baseURLSaving = ref(false)
const baseURLTestResults = ref<UpstreamBaseURLLatency[]>([])
const baseURLTestError = ref('')

const cleanBaseURLs = computed(() => baseURLForm.urls.map((url) => url.trim()).filter(Boolean))

async function openBaseURLDialog() {
  baseURLForm.urls = vendorBaseURLs.value.length ? [...vendorBaseURLs.value] : []
  baseURLForm.proxyId = 0
  baseURLForm.selectedUrl = vendorBaseURLs.value[0] || ''
  baseURLTestResults.value = []
  baseURLTestError.value = ''
  showBaseURLDialog.value = true
  if (proxies.value.length === 0) {
    try {
      proxies.value = await proxiesAPI.getAll()
    } catch (err) {
      appStore.showError(extractApiErrorMessage(err, t('admin.upstreamHub.loadFailed')))
    }
  }
}

function isValidBaseURL(raw: string): boolean {
  const value = raw.trim()
  if (!value) return false
  try {
    const parsed = new URL(value)
    return (parsed.protocol === 'http:' || parsed.protocol === 'https:') && parsed.host !== ''
  } catch {
    return false
  }
}

// 非空但格式非法的行才报错；空行忽略。
function isInvalidBaseURLRow(raw: string): boolean {
  return raw.trim() !== '' && !isValidBaseURL(raw)
}

const dirtyBaseURLs = computed(() => baseURLForm.urls.map((url) => url.trim()).filter(Boolean))
const hasInvalidBaseURLs = computed(() => dirtyBaseURLs.value.some((url) => !isValidBaseURL(url)))
const canMeasureBaseURLs = computed(() => dirtyBaseURLs.value.length > 0 && !hasInvalidBaseURLs.value)
const canSaveBaseURLs = computed(() => canMeasureBaseURLs.value)

function addBaseURLRow() {
  baseURLForm.urls.push('')
}

function removeBaseURLRow(index: number) {
  baseURLForm.urls.splice(index, 1)
}

function normalizeURLForCompare(raw: string): string {
  return raw.trim().replace(/\/+$/, '').toLowerCase()
}

// 测速结果按地址回填到对应行；后端返回的是归一化后的 URL。
function baseURLResultFor(raw: string): UpstreamBaseURLLatency | undefined {
  const target = normalizeURLForCompare(raw)
  if (!target) return undefined
  return baseURLTestResults.value.find((item) => normalizeURLForCompare(item.url) === target)
}

// 测速后按延迟升序排列：可达的在前（快→慢），不可达的排最后；并默认选中第一个可达地址。
function sortBaseURLsByLatency() {
  const rank = (url: string) => {
    const result = baseURLResultFor(url)
    if (!result) return Number.MAX_SAFE_INTEGER
    return result.ok ? result.latency_ms : Number.MAX_SAFE_INTEGER - 1
  }
  baseURLForm.urls = [...baseURLForm.urls].sort((a, b) => rank(a) - rank(b))
  const fastest = baseURLForm.urls.find((url) => baseURLResultFor(url)?.ok)
  baseURLForm.selectedUrl = (fastest || baseURLForm.urls[0] || '').trim()
}

function isSelectedBaseURL(raw: string): boolean {
  const selected = normalizeURLForCompare(baseURLForm.selectedUrl)
  return selected !== '' && normalizeURLForCompare(raw) === selected
}

function selectBaseURL(raw: string) {
  if (!isValidBaseURL(raw)) return
  baseURLForm.selectedUrl = raw.trim()
}

// 保存时把「首选地址」放到第一位，顺序即优先级。
function orderedBaseURLs(): string[] {
  const urls = cleanBaseURLs.value
  const selectedKey = normalizeURLForCompare(baseURLForm.selectedUrl)
  if (!selectedKey) return urls
  const primary = urls.filter((url) => normalizeURLForCompare(url) === selectedKey)
  const rest = urls.filter((url) => normalizeURLForCompare(url) !== selectedKey)
  return [...primary, ...rest]
}

async function testBaseURLs() {
  const urls = cleanBaseURLs.value
  if (urls.length === 0) return
  baseURLTesting.value = true
  baseURLTestError.value = ''
  try {
    baseURLTestResults.value = await upstreamHubAPI.measureBaseURLs({
      base_urls: urls,
      proxy_id: baseURLForm.proxyId || undefined
    })
    sortBaseURLsByLatency()
  } catch (err) {
    baseURLTestError.value = extractApiErrorMessage(err, t('admin.upstreamHub.measureFailed'))
  } finally {
    baseURLTesting.value = false
  }
}

async function saveBaseURLs() {
  if (!detail.value) return
  if (!canSaveBaseURLs.value) return
  const urls = orderedBaseURLs()
  baseURLSaving.value = true
  try {
    await upstreamHubAPI.updateVendor(detail.value.vendor.id, {
      name: detail.value.vendor.name,
      concurrency_budget: detail.value.vendor.concurrency_budget,
      notes: detail.value.vendor.notes || null,
      base_urls: urls
    })
    await refreshDetailEntries()
    showBaseURLDialog.value = false
    appStore.showSuccess(t('admin.upstreamHub.baseURLsSaved'))
  } catch (err) {
    appStore.showError(errorMessage(err))
  } finally {
    baseURLSaving.value = false
  }
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
  if (vendorBaseURLs.value.length === 0) {
    appStore.showError(t('admin.upstreamHub.baseURLPoolRequired'))
    return
  }
  saving.value = true
  try {
    // 地址池由厂家统一维护；最快者由服务端在保存时实测决定，客户端绕不过去。
    const baseURLs = vendorBaseURLs.value
    for (const row of entryRows.value) {
      if (!row.name.trim()) {
        throw new Error(t('admin.upstreamHub.fillRequiredFields'))
      }
      if (row.accountId) {
        await upstreamHubAPI.updateEntry(row.accountId, {
          name: row.name,
          base_url: row.base_url || baseURLs[0] || '',
          api_key: row.api_key,
          concurrency: Number(row.concurrency) || 1,
          priority: Number(row.priority) || 50,
          base_urls: baseURLs
        })
        await applyRowGroups(row.accountId, row)
        if (row.modelsDirty) {
          await upstreamHubAPI.setEntryModels(row.accountId, row.models)
        }
      } else {
        if (!row.api_key.trim()) {
          throw new Error(t('admin.upstreamHub.fillRequiredFields'))
        }
        const created = await upstreamHubAPI.createEntry(detail.value.vendor.id, {
          name: row.name,
          platform: row.platform,
          base_url: row.base_url || baseURLs[0] || '',
          api_key: row.api_key,
          concurrency: Number(row.concurrency) || 1,
          priority: Number(row.priority) || 50,
          base_urls: baseURLs
        })
        await applyRowGroups(created.account_id, row)
        if (row.models.length > 0) {
          await upstreamHubAPI.setEntryModels(created.account_id, row.models)
        }
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
