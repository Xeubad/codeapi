import type { TrendDataPoint } from '@/types'
import { formatDateLocalInput } from '@/utils/format'

export interface DaySnapshot {
  requests: number
  tokens: number
  cost: number
  cacheHitRate: number
}

export function computeCacheHitRate(
  inputTokens: number,
  cacheReadTokens: number,
  cacheCreationTokens = 0
): number {
  const promptTokens = inputTokens + cacheReadTokens + cacheCreationTokens
  return promptTokens > 0 ? (cacheReadTokens / promptTokens) * 100 : 0
}

export function snapshotFromTrendPoint(point: TrendDataPoint | undefined): DaySnapshot | null {
  if (!point) return null
  return {
    requests: point.requests,
    tokens: point.total_tokens,
    cost: point.actual_cost,
    cacheHitRate: computeCacheHitRate(point.input_tokens, point.cache_read_tokens, point.cache_creation_tokens),
  }
}

export function findTrendPointByDate(trend: TrendDataPoint[], date: Date): TrendDataPoint | undefined {
  const key = formatDateLocalInput(date)
  return trend.find(point => point.date.slice(0, 10) === key)
}

export function previousDaySnapshot(trend: TrendDataPoint[], now = new Date()): DaySnapshot | null {
  const yesterday = new Date(now.getFullYear(), now.getMonth(), now.getDate() - 1)
  return snapshotFromTrendPoint(findTrendPointByDate(trend, yesterday))
}
