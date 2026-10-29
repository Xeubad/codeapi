<template>
  <AppLayout>
    <TablePageLayout>
      <template #filters>
        <div class="flex flex-col justify-between gap-4 lg:flex-row lg:items-start">
          <div class="flex flex-1 flex-wrap items-center gap-3">
            <div class="relative w-full sm:w-64">
              <Icon
                name="search"
                size="md"
                class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 dark:text-gray-500"
              />
              <input
                v-model="unitSearch"
                type="text"
                :placeholder="t('admin.aggregation.searchUnits')"
                class="input pl-10"
              />
            </div>
            <span class="text-xs text-gray-500 dark:text-gray-400">
              {{ t('admin.aggregation.unitsTotal', { count: filteredUnits.length }) }}
            </span>
          </div>
          <div class="flex w-full flex-shrink-0 flex-wrap items-center justify-end gap-3 lg:w-auto">
            <button class="btn btn-secondary" :disabled="loading" @click="loadAll">
              <Icon name="refresh" size="md" :class="loading ? 'animate-spin' : ''" />
            </button>
            <button class="btn btn-primary" @click="openCreateUnit">
              <Icon name="plus" size="md" class="mr-2" />
              {{ t('admin.aggregation.createUnit') }}
            </button>
          </div>
        </div>
      </template>

      <template #table>
        <table class="table">
          <thead>
            <tr>
              <th>{{ t('admin.aggregation.unitName') }}</th>
              <th>{{ t('admin.aggregation.entryCount') }}</th>
              <th>{{ t('admin.aggregation.routeCount') }}</th>
              <th>RPM</th>
              <th>{{ t('admin.aggregation.status') }}</th>
              <th class="text-right">{{ t('common.actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="unit in filteredUnits" :key="unit.id">
              <td>
                <div class="font-medium text-gray-900 dark:text-white">{{ unit.name }}</div>
                <div v-if="unit.description" class="max-w-[28rem] truncate text-xs text-gray-500">
                  {{ unit.description }}
                </div>
              </td>
              <td>{{ attachedCountFor(unit.id) }}</td>
              <td>{{ routesCountByGroup[unit.id] ?? '—' }}</td>
              <td>{{ unit.rpm_limit || '∞' }}</td>
              <td>
                <span class="badge" :class="unit.status === 'active' ? 'badge-primary' : 'badge-gray'">
                  {{ unit.status }}
                </span>
              </td>
              <td class="text-right">
                <button class="btn btn-sm btn-secondary" @click="openUnit(unit.id)">
                  {{ t('admin.aggregation.manage') }}
                </button>
              </td>
            </tr>
            <tr v-if="filteredUnits.length === 0">
              <td colspan="6">
                <EmptyState
                  :title="t('admin.aggregation.noUnits')"
                  :description="t('admin.aggregation.noUnitsHint')"
                />
              </td>
            </tr>
          </tbody>
        </table>
      </template>
    </TablePageLayout>

    <!-- 单元管理：沿用其他管理页的弹窗形态，不再切整页，因此没有「返回列表」 -->
    <BaseDialog
      :show="showUnitDialog"
      :title="selectedUnit?.name || t('admin.aggregation.manageUnit')"
      width="full"
      @close="closeUnit"
    >
      <div class="space-y-5">
        <div class="flex flex-wrap items-start justify-between gap-3">
          <div class="min-w-0">
            <div class="flex flex-wrap items-center gap-2">
              <span class="badge badge-primary">{{ t('admin.aggregation.compositeBadge') }}</span>
              <span v-if="selectedUnit?.description" class="truncate text-xs text-gray-500 dark:text-gray-400">
                {{ selectedUnit.description }}
              </span>
            </div>
            <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ t('admin.aggregation.unitDesc') }}</p>
          </div>
          <div class="flex flex-wrap items-center gap-2">
            <button class="btn btn-secondary btn-sm" @click="goToGroups">
              {{ t('admin.aggregation.configureInGroups') }}
            </button>
            <button class="btn btn-secondary btn-sm" :disabled="loading" @click="loadAll">
              {{ t('admin.aggregation.refresh') }}
            </button>
          </div>
        </div>

        <div class="tabs inline-flex w-full max-w-xl flex-wrap sm:w-auto" role="tablist">
          <button
            type="button"
            role="tab"
            class="tab flex-1 sm:flex-none"
            :class="detailTab === 'supply' ? 'tab-active' : ''"
            :aria-selected="detailTab === 'supply'"
            @click="detailTab = 'supply'"
          >
            {{ t('admin.aggregation.tabSupply') }}
          </button>
          <button
            type="button"
            role="tab"
            class="tab flex-1 sm:flex-none"
            :class="detailTab === 'routes' ? 'tab-active' : ''"
            :aria-selected="detailTab === 'routes'"
            @click="detailTab = 'routes'"
          >
            {{ t('admin.aggregation.tabRoutes') }}
          </button>
          <button
            type="button"
            role="tab"
            class="tab flex-1 sm:flex-none"
            :class="detailTab === 'chain' ? 'tab-active' : ''"
            :aria-selected="detailTab === 'chain'"
            @click="detailTab = 'chain'"
          >
            {{ t('admin.aggregation.tabChain') }}
          </button>
        </div>

        <div class="grid gap-4 sm:grid-cols-3">
          <div class="rounded-xl border border-gray-200 px-4 py-3 dark:border-dark-700">
            <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('admin.aggregation.statEntries') }}</p>
            <p class="mt-1 text-xl font-semibold text-gray-900 dark:text-white">{{ attachedCount }}</p>
          </div>
          <div class="rounded-xl border border-gray-200 px-4 py-3 dark:border-dark-700">
            <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('admin.aggregation.statRoutes') }}</p>
            <p class="mt-1 text-xl font-semibold text-gray-900 dark:text-white">{{ routes.length }}</p>
          </div>
          <div class="rounded-xl border border-gray-200 px-4 py-3 dark:border-dark-700">
            <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('admin.aggregation.statKeys') }}</p>
            <p class="mt-1 text-xl font-semibold text-gray-900 dark:text-white">{{ activeKeys }}</p>
            <p class="text-xs text-gray-500 dark:text-gray-400">
              {{ t('admin.aggregation.statKeysHint', { count: keysTotal }) }}
            </p>
          </div>
        </div>

        <!-- 供给汇入 -->
        <section v-if="detailTab === 'supply'" class="space-y-3">
          <div class="flex flex-wrap items-end justify-between gap-2">
            <div>
              <h2 class="text-base font-semibold text-gray-900 dark:text-white">{{ t('admin.aggregation.supply') }}</h2>
              <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('admin.aggregation.supplyHint') }}</p>
            </div>
            <span class="text-xs text-gray-500 dark:text-gray-400">
              {{ t('admin.aggregation.attachedCount', { count: attachedCount }) }}
            </span>
          </div>

          <div
            v-if="entries.length === 0"
            class="rounded-xl border border-dashed border-gray-300 px-4 py-6 text-center text-sm text-gray-500 dark:border-dark-600 dark:text-dark-300"
          >
            {{ t('admin.aggregation.noEntries') }}
          </div>

          <div v-else class="space-y-3">
            <div
              v-for="vendor in vendorGroups"
              :key="vendor.id"
              class="overflow-hidden rounded-xl border border-gray-200 dark:border-dark-700"
            >
              <button
                type="button"
                class="flex w-full flex-wrap items-center gap-x-3 gap-y-1 px-4 py-3 text-left hover:bg-gray-50 dark:hover:bg-dark-700/40"
                @click="toggleVendor(vendor.id)"
              >
                <Icon
                  name="chevronDown"
                  size="sm"
                  class="text-gray-400"
                  :class="isVendorOpen(vendor.id) ? '' : '-rotate-90'"
                />
                <span class="text-sm font-medium text-gray-900 dark:text-white">{{ vendor.name }}</span>
                <span
                  v-if="vendor.base_urls[0]"
                  class="max-w-[16rem] truncate font-mono text-xs text-gray-400 dark:text-dark-400"
                  :title="vendor.base_urls[0]"
                >
                  {{ vendor.base_urls[0] }}
                </span>
                <span class="ml-auto flex flex-wrap items-center gap-3 text-xs text-gray-500 dark:text-gray-400">
                  <span :class="vendor.over_budget ? 'text-red-600 dark:text-red-400' : ''">
                    {{ t('admin.aggregation.vendorBudget', {
                      used: vendor.allocated_concurrency,
                      total: vendor.concurrency_budget || '∞'
                    }) }}
                  </span>
                  <span class="badge" :class="attachedInVendor(vendor) > 0 ? 'badge-primary' : 'badge-gray'">
                    {{ t('admin.aggregation.vendorAttached', {
                      attached: attachedInVendor(vendor),
                      total: vendor.entries.length
                    }) }}
                  </span>
                </span>
              </button>

              <div v-if="isVendorOpen(vendor.id)" class="border-t border-gray-100 dark:border-dark-700">
                <div v-if="vendor.entries.length === 0" class="px-4 py-3 text-xs text-gray-500 dark:text-gray-400">
                  {{ t('admin.aggregation.noEntriesForVendor') }}
                </div>
                <div v-else class="table-container">
                  <table class="table">
                    <thead>
                      <tr>
                        <th class="w-12">{{ t('admin.aggregation.attach') }}</th>
                        <th>{{ t('admin.aggregation.entry') }}</th>
                        <th>{{ t('admin.aggregation.platform') }}</th>
                        <th>{{ t('admin.aggregation.concurrency') }}</th>
                        <th>{{ t('admin.aggregation.priority') }}</th>
                        <th>{{ t('admin.aggregation.models') }}</th>
                        <th>{{ t('admin.aggregation.health') }}</th>
                        <th class="text-right">{{ t('admin.aggregation.enabledHere') }}</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="entry in vendorEntries(vendor)" :key="entry.account_id">
                        <td>
                          <input
                            type="checkbox"
                            class="cursor-pointer"
                            :checked="isAttached(entry)"
                            :aria-label="entry.name"
                            @change="toggleAttach(entry)"
                          />
                        </td>
                        <td>
                          <div class="font-medium text-gray-900 dark:text-white">{{ entry.name }}</div>
                        </td>
                        <td><span class="badge badge-primary">{{ platformLabel(entry.platform) }}</span></td>
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
                            {{ t('admin.aggregation.modelsCount', { count: entry.models?.length || 0 }) }}
                          </span>
                        </td>
                        <td>
                          <div class="flex flex-wrap items-center gap-1">
                            <span v-if="entry.status !== 'active'" class="badge badge-danger">{{ entry.status }}</span>
                            <span v-if="!entry.schedulable" class="badge badge-gray">
                              {{ t('admin.aggregation.entryUnschedulable') }}
                            </span>
                            <span
                              v-if="(entry.models?.length || 0) === 0"
                              class="badge badge-warning"
                            >
                              {{ t('admin.aggregation.entryNoModels') }}
                            </span>
                          </div>
                        </td>
                        <td class="text-right">
                          <span
                            :title="isAttached(entry) ? '' : t('admin.aggregation.attachFirst')"
                            class="inline-flex"
                          >
                            <Toggle
                              :model-value="isEnabledHere(entry)"
                              :disabled="!isAttached(entry)"
                              @update:model-value="toggleEntryEnabled(entry)"
                            />
                          </span>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
          </div>
        </section>

        <!-- 归一模型 -->
        <section v-else-if="detailTab === 'routes'" class="space-y-4">
          <div>
            <h2 class="text-base font-semibold text-gray-900 dark:text-white">{{ t('admin.aggregation.modelRoutes') }}</h2>
            <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('admin.aggregation.modelRoutesHint') }}</p>
          </div>

          <div class="table-container">
            <table class="table">
              <thead>
                <tr>
                  <th>{{ t('admin.aggregation.publicModel') }}</th>
                  <th>{{ t('admin.aggregation.matchType') }}</th>
                  <th>{{ t('admin.aggregation.endpoint') }}</th>
                  <th>{{ t('admin.aggregation.targetPlatform') }}</th>
                  <th>{{ t('admin.aggregation.upstreamModel') }}</th>
                  <th>{{ t('admin.aggregation.enabled') }}</th>
                  <th class="text-right">{{ t('common.actions') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="route in routes" :key="route.id">
                  <td class="font-mono text-xs">{{ route.public_model }}</td>
                  <td>{{ matchTypeLabel(route.match_type) }}</td>
                  <td>{{ endpointLabel(route.endpoint) }}</td>
                  <td><span class="badge badge-primary">{{ platformLabel(route.target_platform) }}</span></td>
                  <td class="font-mono text-xs">{{ route.upstream_model || route.public_model }}</td>
                  <td>
                    <Toggle :model-value="route.enabled" @update:model-value="() => toggleRouteEnabled(route)" />
                  </td>
                  <td class="whitespace-nowrap text-right">
                    <button class="btn btn-sm btn-secondary" @click="startEditRoute(route)">
                      {{ t('admin.aggregation.edit') }}
                    </button>
                    <button class="btn btn-sm btn-ghost ml-1 text-red-600" @click="removeRoute(route)">
                      {{ t('admin.aggregation.delete') }}
                    </button>
                  </td>
                </tr>
                <tr v-if="routes.length === 0">
                  <td colspan="7" class="py-4 text-center text-xs text-gray-500">
                    {{ t('admin.aggregation.noRoutes') }}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <div class="rounded-xl border border-gray-200 p-4 dark:border-dark-700">
            <div class="mb-3 flex items-center justify-between">
              <h3 class="text-sm font-semibold text-gray-900 dark:text-white">
                {{ editingRouteId ? t('admin.aggregation.editRoute') : t('admin.aggregation.addRoute') }}
              </h3>
              <button
                v-if="editingRouteId"
                type="button"
                class="text-xs text-gray-500 hover:underline dark:text-gray-400"
                @click="resetRouteForm"
              >
                {{ t('admin.aggregation.cancelEdit') }}
              </button>
            </div>
            <div class="grid gap-3 sm:grid-cols-2">
              <div>
                <label class="input-label">{{ t('admin.aggregation.publicModel') }}</label>
                <input
                  v-model="routeForm.public_model"
                  class="input font-mono"
                  type="text"
                  :placeholder="t('admin.aggregation.publicModelPlaceholder')"
                />
              </div>
              <div>
                <label class="input-label">{{ t('admin.aggregation.targetPlatform') }}</label>
                <Select v-model="routeForm.target_platform" :options="routePlatformOptions" />
              </div>
              <div>
                <label class="input-label">{{ t('admin.aggregation.upstreamModel') }}</label>
                <input
                  v-model="routeForm.upstream_model"
                  class="input font-mono"
                  type="text"
                  list="aggregation-upstream-models"
                  :placeholder="t('admin.aggregation.upstreamModelPlaceholder')"
                />
                <datalist id="aggregation-upstream-models">
                  <option v-for="model in upstreamModelCandidates" :key="model" :value="model" />
                </datalist>
              </div>
              <div>
                <label class="input-label">{{ t('admin.aggregation.matchType') }}</label>
                <Select v-model="routeForm.match_type" :options="matchTypeOptions" />
              </div>
              <div>
                <label class="input-label">{{ t('admin.aggregation.endpoint') }}</label>
                <Select v-model="routeForm.endpoint" :options="endpointOptions" />
              </div>
              <div>
                <label class="input-label">{{ t('admin.aggregation.notes') }}</label>
                <input
                  v-model="routeForm.notes"
                  class="input"
                  type="text"
                  :placeholder="t('admin.aggregation.notesPlaceholder')"
                />
              </div>
            </div>
            <p v-if="selectedPlatformSupply === 0" class="mt-2 text-xs text-amber-600 dark:text-amber-400">
              {{ t('admin.aggregation.platformNoSupplyHint') }}
            </p>
            <p class="mt-2 text-xs text-gray-500 dark:text-gray-400">{{ t('admin.aggregation.routeCandidatesHint') }}</p>
            <div class="mt-3 flex justify-end">
              <button
                class="btn btn-primary btn-sm"
                :disabled="saving || !routeForm.public_model.trim()"
                @click="saveRoute"
              >
                {{ editingRouteId ? t('admin.aggregation.save') : t('admin.aggregation.addRoute') }}
              </button>
            </div>
          </div>
        </section>

        <!-- 链路视图 -->
        <section v-else class="space-y-4">
          <div class="flex flex-wrap items-center justify-between gap-2">
            <div>
              <h2 class="text-base font-semibold text-gray-900 dark:text-white">{{ t('admin.aggregation.chainView') }}</h2>
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
              <h3 class="text-xs font-semibold uppercase text-gray-500 dark:text-gray-400">{{ t('admin.aggregation.vendors') }}</h3>
              <div
                v-for="vendor in chainVendors"
                :key="vendor.id"
                class="rounded-lg border border-gray-200 px-2 py-1.5 text-sm dark:border-dark-600"
              >
                <div class="truncate font-medium text-gray-900 dark:text-white">{{ vendor.name }}</div>
                <div class="text-xs text-gray-500 dark:text-gray-400">
                  {{ t('admin.aggregation.vendorBudget', { used: vendor.allocated_concurrency, total: vendor.concurrency_budget || '∞' }) }}
                </div>
              </div>
              <p v-if="chainVendors.length === 0" class="text-xs text-gray-400 dark:text-dark-400">
                {{ t('admin.aggregation.noChainVendors') }}
              </p>
            </div>
            <div class="space-y-2">
              <h3 class="text-xs font-semibold uppercase text-gray-500 dark:text-gray-400">{{ t('admin.aggregation.entries') }}</h3>
              <div
                v-for="entry in chainEntries"
                :key="entry.account_id"
                class="rounded-lg border px-2 py-1.5 text-sm"
                :class="
                  isEnabledHere(entry)
                    ? 'border-primary-300 bg-primary-50 text-primary-800 dark:border-primary-700 dark:bg-primary-900/20 dark:text-primary-300'
                    : 'border-dashed border-gray-300 text-gray-500 dark:border-dark-600 dark:text-dark-300'
                "
              >
                <div class="truncate">{{ entry.name }}</div>
                <div class="text-xs opacity-80">
                  {{ platformLabel(entry.platform) }} ·
                  {{ isEnabledHere(entry) ? t('admin.aggregation.modelsCount', { count: entry.models?.length || 0 }) : t('admin.aggregation.disabledHere') }}
                </div>
              </div>
              <p v-if="chainEntries.length === 0" class="text-xs text-gray-400 dark:text-dark-400">
                {{ t('admin.aggregation.noChainEntries') }}
              </p>
            </div>
            <div class="space-y-2">
              <h3 class="text-xs font-semibold uppercase text-gray-500 dark:text-gray-400">{{ t('admin.aggregation.group') }}</h3>
              <div
                v-if="selectedUnit"
                class="rounded-lg border border-primary-300 bg-primary-50 px-2 py-1.5 text-sm text-primary-800 dark:border-primary-700 dark:bg-primary-900/20 dark:text-primary-300"
              >
                <div class="truncate font-medium">{{ selectedUnit.name }}</div>
                <div class="text-xs opacity-80">
                  {{ t('admin.aggregation.attachedCount', { count: attachedCount }) }} · RPM {{ selectedUnit.rpm_limit || '∞' }}
                </div>
                <div class="text-xs opacity-80">
                  {{ t('admin.aggregation.routeCount') }} {{ routes.length }}
                </div>
              </div>
            </div>
            <div class="space-y-2">
              <h3 class="text-xs font-semibold uppercase text-gray-500 dark:text-gray-400">{{ t('admin.aggregation.customers') }}</h3>
              <div class="max-h-[28rem] space-y-2 overflow-y-auto pr-1">
              <div
                v-if="keys.length === 0"
                class="rounded-lg border border-dashed border-gray-300 px-2 py-1.5 text-xs text-gray-500 dark:border-dark-600 dark:text-dark-300"
              >
                {{ t('admin.aggregation.noKeys') }}
              </div>
              <div
                v-for="key in keys"
                :key="key.id"
                class="rounded-lg border border-gray-200 px-2 py-1.5 text-sm dark:border-dark-600"
              >
                <div class="flex items-center justify-between gap-2">
                  <span class="truncate font-medium text-gray-900 dark:text-white">
                    {{ key.name || maskKey(key.key) }}
                  </span>
                  <span class="badge" :class="key.status === 'active' ? 'badge-success' : 'badge-gray'">
                    {{ keyStatusLabel(key.status) }}
                  </span>
                </div>
                <div class="truncate text-xs text-gray-500 dark:text-gray-400">{{ key.user?.email || '—' }}</div>
                <div class="font-mono text-xs text-gray-400 dark:text-dark-400">{{ maskKey(key.key) }}</div>
              </div>
              <p v-if="keysTotal > keys.length" class="text-xs text-gray-400 dark:text-dark-400">
                {{ t('admin.aggregation.keysTruncated', { shown: keys.length, total: keysTotal }) }}
              </p>
              </div>
            </div>
          </div>
        </section>
      </div>

      <template #footer>
        <button class="btn btn-secondary" @click="closeUnit">{{ t('admin.aggregation.close') }}</button>
      </template>
    </BaseDialog>

    <BaseDialog
      :show="showCreateDialog"
      :title="t('admin.aggregation.createUnit')"
      width="normal"
      @close="showCreateDialog = false"
    >
      <div class="space-y-3">
        <div>
          <label class="input-label">{{ t('admin.aggregation.unitName') }}</label>
          <input v-model="createForm.name" class="input" type="text" :placeholder="t('admin.aggregation.unitNamePlaceholder')" />
        </div>
        <div>
          <label class="input-label">{{ t('admin.aggregation.unitDescription') }}</label>
          <input
            v-model="createForm.description"
            class="input"
            type="text"
            :placeholder="t('admin.aggregation.unitDescriptionPlaceholder')"
          />
        </div>
        <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('admin.aggregation.createUnitHint') }}</p>
      </div>
      <template #footer>
        <button class="btn btn-secondary" @click="showCreateDialog = false">{{ t('admin.aggregation.cancel') }}</button>
        <button class="btn btn-primary" :disabled="creating || !createForm.name.trim()" @click="createUnit">
          {{ t('admin.aggregation.createUnit') }}
        </button>
      </template>
    </BaseDialog>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { extractApiErrorMessage } from '@/utils/apiError'
import { maskApiKey } from '@/utils/maskApiKey'
import upstreamHubAPI, { type UpstreamEntry, type UpstreamVendor } from '@/api/admin/upstreamHub'
import groupsAPI, { type GroupAPIKey } from '@/api/admin/groups'
import { CONCRETE_PLATFORM_OPTIONS } from '@/constants/platforms'
import { platformDisplayName } from '@/constants/platformCatalog'
import type {
  AdminGroup,
  CompositeModelRoute,
  CompositeModelRouteInput,
  CompositeRouteEndpoint,
  CompositeRouteMatchType
} from '@/types'
import AppLayout from '@/components/layout/AppLayout.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import Select from '@/components/common/Select.vue'
import Toggle from '@/components/common/Toggle.vue'
import Icon from '@/components/icons/Icon.vue'

interface VendorGroup extends UpstreamVendor {
  entries: UpstreamEntry[]
}

const KEY_PAGE_SIZE = 100

const { t } = useI18n()
const appStore = useAppStore()
const route = useRoute()
const router = useRouter()

const loading = ref(false)
const saving = ref(false)
const creating = ref(false)
const groups = ref<AdminGroup[]>([])
const vendors = ref<UpstreamVendor[]>([])
const entries = ref<UpstreamEntry[]>([])
const routes = ref<CompositeModelRoute[]>([])
const keys = ref<GroupAPIKey[]>([])
const keysTotal = ref(0)
const routesCountByGroup = ref<Record<number, number>>({})
const selectedUnitId = ref<number | null>(null)
const showUnitDialog = ref(false)
const detailTab = ref<'supply' | 'routes' | 'chain'>('supply')
const unitSearch = ref('')
const expandedVendors = ref<Record<number, boolean>>({})

const editingRouteId = ref<number | null>(null)
const routeForm = reactive<CompositeModelRouteInput>({
  public_model: '',
  target_platform: CONCRETE_PLATFORM_OPTIONS[0]?.value || 'openai',
  upstream_model: '',
  match_type: 'exact',
  endpoint: 'any',
  priority: 100,
  enabled: true,
  notes: ''
})

const units = computed(() => groups.value.filter((group) => group.platform === 'composite'))
const filteredUnits = computed(() => {
  const keyword = unitSearch.value.trim().toLowerCase()
  if (!keyword) return units.value
  return units.value.filter(
    (unit) =>
      unit.name.toLowerCase().includes(keyword) ||
      (unit.description || '').toLowerCase().includes(keyword)
  )
})
const selectedUnit = computed(() => units.value.find((group) => group.id === selectedUnitId.value) || null)
const vendorGroups = computed<VendorGroup[]>(() =>
  vendors.value.map((vendor) => ({
    ...vendor,
    entries: entries.value.filter((entry) => entry.vendor_id === vendor.id)
  }))
)
const attachedCount = computed(() => entries.value.filter((entry) => isAttached(entry)).length)
const chainVendors = computed(() => vendorGroups.value.filter((vendor) => attachedInVendor(vendor) > 0))
const chainEntries = computed(() =>
  entries.value
    .filter((entry) => isAttached(entry))
    .sort((a, b) => {
      const enabledDelta = Number(isEnabledHere(b)) - Number(isEnabledHere(a))
      return enabledDelta || a.name.localeCompare(b.name)
    })
)
const activeKeys = computed(() => keys.value.filter((key) => key.status === 'active').length)
const platformSupply = computed(() => {
  const counts: Record<string, number> = {}
  for (const entry of entries.value) {
    if (isAttached(entry)) counts[entry.platform] = (counts[entry.platform] || 0) + 1
  }
  return counts
})
const selectedPlatformSupply = computed(() => platformSupply.value[routeForm.target_platform] || 0)
const upstreamModelCandidates = computed(() => {
  const models = new Set<string>()
  for (const entry of entries.value) {
    if (!isAttached(entry) || entry.platform !== routeForm.target_platform) continue
    for (const model of entry.models || []) models.add(model)
  }
  return [...models].sort((a, b) => a.localeCompare(b))
})

function platformLabel(platform: string): string {
  return platformDisplayName(platform) || platform
}

function attachedCountFor(groupId: number): number {
  return entries.value.filter((entry) => (entry.groups || []).some((group) => group.group_id === groupId)).length
}

function groupRef(entry: UpstreamEntry) {
  return entry.groups?.find((group) => group.group_id === selectedUnitId.value)
}

function isAttached(entry: UpstreamEntry): boolean {
  return !!groupRef(entry)
}

function isEnabledHere(entry: UpstreamEntry): boolean {
  return groupRef(entry)?.enabled ?? false
}

function attachedInVendor(vendor: VendorGroup): number {
  return vendor.entries.filter((entry) => isAttached(entry)).length
}

function vendorEntries(vendor: VendorGroup): UpstreamEntry[] {
  return [...vendor.entries].sort((a, b) => {
    const attachedDelta = Number(isAttached(b)) - Number(isAttached(a))
    return attachedDelta || a.name.localeCompare(b.name)
  })
}

function isVendorOpen(vendorId: number): boolean {
  return expandedVendors.value[vendorId] === true
}

function toggleVendor(vendorId: number) {
  expandedVendors.value = { ...expandedVendors.value, [vendorId]: !isVendorOpen(vendorId) }
}

function maskKey(key: string): string {
  return maskApiKey(key)
}

function keyStatusLabel(status: string): string {
  const key = `admin.aggregation.keyStatus.${status}`
  const label = t(key)
  return label === key ? status : label
}

function errorMessage(err: unknown): string {
  return extractApiErrorMessage(err, t('admin.aggregation.loadFailed'))
}

async function loadUnits() {
  groups.value = await groupsAPI.getAllIncludingInactive()
  const counts: Record<number, number> = {}
  await Promise.all(
    units.value.map(async (unit) => {
      try {
        counts[unit.id] = (await groupsAPI.listCompositeRoutes(unit.id)).length
      } catch {
        counts[unit.id] = 0
      }
    })
  )
  routesCountByGroup.value = counts
}

async function loadUnitData() {
  if (!selectedUnitId.value) {
    routes.value = []
    keys.value = []
    keysTotal.value = 0
    return
  }
  try {
    const [routeList, keyPage] = await Promise.all([
      groupsAPI.listCompositeRoutes(selectedUnitId.value),
      groupsAPI.getGroupApiKeys(selectedUnitId.value, 1, KEY_PAGE_SIZE)
    ])
    routes.value = routeList
    keys.value = keyPage.items
    keysTotal.value = keyPage.total
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
    await loadUnits()
    if (showUnitDialog.value) await loadUnitData()
  } catch (err) {
    appStore.showError(errorMessage(err))
  } finally {
    loading.value = false
  }
}

async function refreshEntries() {
  const [vendorList, entryList] = await Promise.all([upstreamHubAPI.listVendors(), upstreamHubAPI.listAllEntries()])
  vendors.value = vendorList
  entries.value = entryList
}

function openUnit(unitId: number) {
  selectedUnitId.value = unitId
  detailTab.value = 'supply'
  showUnitDialog.value = true
  resetRouteForm()
  const anyAttached = entries.value.some((entry) => isAttached(entry))
  const vendorState: Record<number, boolean> = {}
  for (const vendor of vendorGroups.value) {
    vendorState[vendor.id] = anyAttached ? vendor.entries.some((entry) => isAttached(entry)) : true
  }
  expandedVendors.value = vendorState
  void loadUnitData()
}

function closeUnit() {
  showUnitDialog.value = false
  selectedUnitId.value = null
  routes.value = []
  keys.value = []
  keysTotal.value = 0
  if (route.query.unit) {
    const query = { ...route.query }
    delete query.unit
    void router.replace({ query })
  }
}

function goToGroups() {
  void router.push('/admin/groups')
}

// ---------- 新建单元 ----------
const showCreateDialog = ref(false)
const createForm = reactive({ name: '', description: '' })

function openCreateUnit() {
  createForm.name = ''
  createForm.description = ''
  showCreateDialog.value = true
}

async function createUnit() {
  if (!createForm.name.trim()) return
  creating.value = true
  try {
    const group = await groupsAPI.create({
      name: createForm.name.trim(),
      description: createForm.description.trim() || undefined,
      platform: 'composite',
      rate_multiplier: 1,
      subscription_type: 'standard'
    })
    showCreateDialog.value = false
    appStore.showSuccess(t('admin.aggregation.unitCreated'))
    await loadAll()
    openUnit(group.id)
  } catch (err) {
    appStore.showError(errorMessage(err))
  } finally {
    creating.value = false
  }
}

// ---------- 供给 ----------
async function toggleAttach(entry: UpstreamEntry) {
  if (!selectedUnitId.value) return
  const otherGroups = (entry.groups || [])
    .map((group) => group.group_id)
    .filter((groupId) => groupId !== selectedUnitId.value)
  const nextGroupIds = isAttached(entry) ? otherGroups : [...otherGroups, selectedUnitId.value]
  try {
    await upstreamHubAPI.setEntryGroups(entry.account_id, nextGroupIds)
    expandedVendors.value = { ...expandedVendors.value, [entry.vendor_id]: true }
    await refreshEntries()
  } catch (err) {
    appStore.showError(errorMessage(err))
    await refreshEntries().catch(() => {})
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

// ---------- 归一模型 ----------
const routePlatformOptions = computed(() =>
  CONCRETE_PLATFORM_OPTIONS.map((platform) => {
    const supply = platformSupply.value[platform.value] || 0
    return {
      value: platform.value,
      label:
        supply > 0
          ? t('admin.aggregation.platformSupplyOption', { name: platform.label, count: supply })
          : t('admin.aggregation.platformNoSupplyOption', { name: platform.label })
    }
  })
)
const matchTypeOptions = computed(() => [
  { value: 'exact', label: t('admin.aggregation.matchExact') },
  { value: 'prefix', label: t('admin.aggregation.matchPrefix') }
])
const endpointOptions = computed(() =>
  (
    [
      'any',
      'messages',
      'count_tokens',
      'responses',
      'chat_completions',
      'embeddings',
      'images',
      'gemini'
    ] as CompositeRouteEndpoint[]
  ).map((endpoint) => ({ value: endpoint, label: t(`admin.aggregation.endpoints.${endpoint}`) }))
)

function matchTypeLabel(matchType: CompositeRouteMatchType): string {
  return matchType === 'prefix'
    ? t('admin.aggregation.matchPrefix')
    : t('admin.aggregation.matchExact')
}

function endpointLabel(endpoint: CompositeRouteEndpoint): string {
  const key = `admin.aggregation.endpoints.${endpoint}`
  const label = t(key)
  return label === key ? endpoint : label
}

function resetRouteForm() {
  editingRouteId.value = null
  routeForm.public_model = ''
  routeForm.upstream_model = ''
  routeForm.match_type = 'exact'
  routeForm.endpoint = 'any'
  routeForm.notes = ''
  routeForm.priority = 100
  routeForm.enabled = true
  routeForm.target_platform = CONCRETE_PLATFORM_OPTIONS[0]?.value || 'openai'
}

function startEditRoute(route: CompositeModelRoute) {
  editingRouteId.value = route.id
  routeForm.public_model = route.public_model
  routeForm.target_platform = route.target_platform
  routeForm.upstream_model = route.upstream_model || ''
  routeForm.match_type = route.match_type
  routeForm.endpoint = route.endpoint
  routeForm.notes = route.notes || ''
  routeForm.priority = route.priority
  routeForm.enabled = route.enabled
}

async function saveRoute() {
  if (!selectedUnitId.value) return
  if (!routeForm.public_model.trim()) {
    appStore.showError(t('admin.aggregation.routeRequired'))
    return
  }
  saving.value = true
  try {
    const payload: CompositeModelRouteInput = {
      public_model: routeForm.public_model.trim(),
      match_type: routeForm.match_type,
      target_platform: routeForm.target_platform,
      upstream_model: (routeForm.upstream_model || '').trim() || undefined,
      endpoint: routeForm.endpoint,
      priority: routeForm.priority ?? 100,
      enabled: routeForm.enabled ?? true,
      notes: routeForm.notes?.trim() || undefined
    }
    if (editingRouteId.value) {
      await groupsAPI.updateCompositeRoute(selectedUnitId.value, editingRouteId.value, payload)
      appStore.showSuccess(t('admin.aggregation.routeUpdated'))
    } else {
      await groupsAPI.createCompositeRoute(selectedUnitId.value, payload)
      appStore.showSuccess(t('admin.aggregation.routeCreated'))
    }
    resetRouteForm()
    await loadUnitData()
    await loadUnits()
  } catch (err) {
    const message = errorMessage(err)
    if (/already exists|COMPOSITE_ROUTE_EXISTS/i.test(message)) {
      appStore.showError(t('admin.aggregation.routeExists'))
    } else {
      appStore.showError(message)
    }
  } finally {
    saving.value = false
  }
}

async function toggleRouteEnabled(route: CompositeModelRoute) {
  if (!selectedUnitId.value) return
  try {
    await groupsAPI.updateCompositeRoute(selectedUnitId.value, route.id, {
      public_model: route.public_model,
      match_type: route.match_type,
      target_platform: route.target_platform,
      upstream_model: route.upstream_model,
      endpoint: route.endpoint,
      priority: route.priority,
      enabled: !route.enabled,
      notes: route.notes
    })
    await loadUnitData()
  } catch (err) {
    appStore.showError(errorMessage(err))
  }
}

async function removeRoute(route: CompositeModelRoute) {
  if (!selectedUnitId.value) return
  try {
    await groupsAPI.deleteCompositeRoute(selectedUnitId.value, route.id)
    appStore.showSuccess(t('admin.aggregation.routeDeleted'))
    if (editingRouteId.value === route.id) resetRouteForm()
    await loadUnitData()
    await loadUnits()
  } catch (err) {
    appStore.showError(errorMessage(err))
  }
}

onMounted(async () => {
  await loadAll()
  const unitParam = Number(route.query.unit)
  if (Number.isFinite(unitParam) && unitParam > 0 && units.value.some((unit) => unit.id === unitParam)) {
    openUnit(unitParam)
  }
})
</script>
