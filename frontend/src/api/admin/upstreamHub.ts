import { apiClient } from '../client'

export interface UpstreamVendor {
  id: number
  name: string
  concurrency_budget: number
  notes: string
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
}

export interface UpstreamEntryRequest {
  name: string
  platform?: string
  base_url: string
  api_key?: string
  concurrency?: number
  priority?: number
  schedulable?: boolean
}

export async function listVendors(): Promise<UpstreamVendor[]> {
  const { data } = await apiClient.get<UpstreamVendor[]>('/admin/upstream-vendors')
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

export async function listGroups(): Promise<UpstreamGroupOption[]> {
  const { data } = await apiClient.get<UpstreamGroupOption[]>('/admin/upstream-groups')
  return data
}

const upstreamHubAPI = {
  listVendors,
  getVendor,
  createVendor,
  updateVendor,
  deleteVendor,
  createEntry,
  updateEntry,
  duplicateEntry,
  deleteEntry,
  setEntrySchedulable,
  setEntryGroups,
  setEntryGroupEnabled,
  listGroups
}

export default upstreamHubAPI
