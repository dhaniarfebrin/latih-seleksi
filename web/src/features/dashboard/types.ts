// Field names mirror api/internal/insight/insight.go.
export interface SectionStat {
  section: string
  answered: number
  correct: number
  avg_time_ms: number
}

export interface CauseCounts {
  tekanan_waktu: number
  ceroboh: number
  kurang_paham: number
  belum_jelas: number
}

export interface PatternStat {
  pattern_id: string
  section: string
  answered: number
  correct: number
  avg_time_ms: number
  causes: CauseCounts
  top_misconception_id: string
}

export interface SessionStat {
  attempt_id: string
  mode: string
  pattern_id: string
  started_at: string
  answered: number
  correct: number
}

export interface QueueItem {
  pattern_id: string
  reason: 'akurasi_rendah' | 'belum_dicoba'
  accuracy_pct: number
}

export interface Dashboard {
  sections: SectionStat[]
  patterns: PatternStat[]
  sessions: SessionStat[]
  queue: QueueItem[]
}
