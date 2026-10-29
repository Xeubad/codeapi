import { apiClient } from '../client'

export interface UpstreamVendor {
  id: number
  name: string
  concurrency_budget: number
  notes: string
  base_urls: string[]
  entry_count: number
  allocated_concurrency: number
  referenced_group_count: number
  over_budget: boolean
  created_at: string
  updated_at: string
}

export interface UpstreamEntryGroupRef {
  group_id: number
  group_name: string
  platform: string
  priority: number
  enabled: boolean
}

export interface UpstreamEntry {
  account_id: number
  vendor_id: number
  name: string
  platform: string
  base_url: string
  has_api_key: boolean
  concurrency: number
  priority: number
  schedulable: boolean
  status: string
  groups: UpstreamEntryGroupRef[]
  models: string[]
}

export interface UpstreamEntryModelProbe {
  models: string[]
  existing: string[]
  /** upstream = 上游实时返回；models.dev = 公共目录兜底候选 */
  source?: string
  warning?: string
}

export interface UpstreamVendorDetail {
  vendor: UpstreamVendor
  entries: UpstreamEntry[]
  warnings: string[]
}

export interface UpstreamGroupOption {
  id: number
  name: string
  platform: string
}

export interface UpstreamVendorRequest {
  name: string
  concurrency_budget: number
  notes?: string | null
  base_urls?: string[]
}

export interface UpstreamBaseURLLatency {
  url: string
  ok: boolean
  latency_ms: number
  status_code?: number
  error?: string
}

export interface AggregationUnit {
  id: number
  name: string
  description: string
  status: string
  rpm_limit: number
  entry_count: number
  route_count: number
}

/** 聚合调度列表：composite 单元 + 已挂载条目数/路由数，一次取全。 */
export async function listAggregationUnits(): Promise<AggregationUnit[]> {
  const { data } = await apiClient.get<AggregationUnit[]>('/admin/aggregation/units')
  return data
}

export interface UpstreamEntryRequest {
  name: string
  platform?: string
  base_url: string
  api_key?: string
  concurrency?: number
  priority?: number
  schedulable?: boolean
  /** 厂家候选地址；>=2 个时服务端在保存前实测并只固化最快的一个 */
  base_urls?: string[]
}

export async function listVendors(): Promise<UpstreamVendor[]> {
  const { data } = await apiClient.get<UpstreamVendor[]>('/admin/upstream-vendors')
  return data
}

export async function listAllEntries(): Promise<UpstreamEntry[]> {
  const { data } = await apiClient.get<UpstreamEntry[]>('/admin/upstream-entries')
  return data
}

export async function getVendor(id: number): Promise<UpstreamVendorDetail> {
  const { data } = await apiClient.get<UpstreamVendorDetail>(`/admin/upstream-vendors/${id}`)
  return data
}

export async function createVendor(payload: UpstreamVendorRequest): Promise<UpstreamVendor> {
  const { data } = await apiClient.post<UpstreamVendor>('/admin/upstream-vendors', payload)
  return data
}

export async function updateVendor(id: number, payload: UpstreamVendorRequest): Promise<UpstreamVendor> {
  const { data } = await apiClient.put<UpstreamVendor>(`/admin/upstream-vendors/${id}`, payload)
  return data
}

export async function deleteVendor(id: number): Promise<void> {
  await apiClient.delete(`/admin/upstream-vendors/${id}`)
}

export async function createEntry(vendorId: number, payload: UpstreamEntryRequest): Promise<UpstreamEntry> {
  const { data } = await apiClient.post<UpstreamEntry>(`/admin/upstream-vendors/${vendorId}/entries`, payload)
  return data
}

/** 把账号管理里已有的转发账号纳入指定厂家（只写归属，不动账号其它字段）。 */
export async function importEntry(vendorId: number, accountId: number): Promise<UpstreamEntry> {
  const { data } = await apiClient.post<UpstreamEntry>(`/admin/upstream-vendors/${vendorId}/entries/import`, {
    account_id: accountId
  })
  return data
}

export async function updateEntry(accountId: number, payload: UpstreamEntryRequest): Promise<UpstreamEntry> {
  const { data } = await apiClient.put<UpstreamEntry>(`/admin/upstream-entries/${accountId}`, payload)
  return data
}

export interface UpstreamEntryDuplicateRequest {
  name: string
  platform: string
  concurrency?: number
  priority?: number
}

export async function duplicateEntry(accountId: number, payload: UpstreamEntryDuplicateRequest): Promise<UpstreamEntry> {
  const { data } = await apiClient.post<UpstreamEntry>(`/admin/upstream-entries/${accountId}/duplicate`, payload)
  return data
}

export async function deleteEntry(accountId: number): Promise<void> {
  await apiClient.delete(`/admin/upstream-entries/${accountId}`)
}

export async function setEntrySchedulable(accountId: number, schedulable: boolean): Promise<UpstreamEntry> {
  const { data } = await apiClient.put<UpstreamEntry>(`/admin/upstream-entries/${accountId}/schedulable`, { schedulable })
  return data
}

export async function setEntryGroups(accountId: number, groupIds: number[]): Promise<UpstreamEntryGroupRef[]> {
  const { data } = await apiClient.put<UpstreamEntryGroupRef[]>(`/admin/upstream-entries/${accountId}/groups`, {
    group_ids: groupIds
  })
  return data
}

export async function setEntryGroupEnabled(accountId: number, groupId: number, enabled: boolean): Promise<void> {
  await apiClient.put(`/admin/upstream-entries/${accountId}/groups/${groupId}/enabled`, { enabled })
}

export async function probeEntryModels(accountId: number): Promise<UpstreamEntryModelProbe> {
  const { data } = await apiClient.post<UpstreamEntryModelProbe>(`/admin/upstream-entries/${accountId}/probe-models`)
  return data
}

export async function setEntryModels(accountId: number, models: string[]): Promise<UpstreamEntry> {
  const { data } = await apiClient.put<UpstreamEntry>(`/admin/upstream-entries/${accountId}/models`, { models })
  return data
}

export interface MeasureBaseURLsRequest {
  account_id?: number
  platform?: string
  api_key?: string
  base_urls: string[]
  /** 本次测速的出口（来自 IP 管理的已有代理），不落库 */
  proxy_id?: number
}

export async function measureBaseURLs(payload: MeasureBaseURLsRequest): Promise<UpstreamBaseURLLatency[]> {
  const { data } = await apiClient.post<UpstreamBaseURLLatency[]>('/admin/upstream-base-urls/latency', payload)
  return data
}

export async function listGroups(): Promise<UpstreamGroupOption[]> {
  const { data } = await apiClient.get<UpstreamGroupOption[]>('/admin/upstream-groups')
  return data
}

const upstreamHubAPI = {
  listVendors,
  listAllEntries,
  getVendor,
  createVendor,
  updateVendor,
  deleteVendor,
  createEntry,
  importEntry,
  updateEntry,
  duplicateEntry,
  deleteEntry,
  setEntrySchedulable,
  setEntryGroups,
  setEntryGroupEnabled,
  probeEntryModels,
  setEntryModels,
  measureBaseURLs,
  listGroups,
  listAggregationUnits
}

export default upstreamHubAPI
