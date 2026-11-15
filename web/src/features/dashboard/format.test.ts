import { describe, expect, it } from 'vitest'
import { causeText, fmtDuration, fmtPct, queueReason } from './format'

describe('format', () => {
  it('formats percent and handles zero answers', () => {
    expect(fmtPct(2, 3)).toBe('67%')
    expect(fmtPct(0, 0)).toBe('-')
  })
  it('formats durations', () => {
    expect(fmtDuration(12400)).toBe('12 dtk')
    expect(fmtDuration(65000)).toBe('1:05')
  })
  it('words causes as estimates', () => {
    expect(causeText('ceroboh', 2)).toBe('Terlalu cepat (mungkin ceroboh): 2 soal')
    expect(causeText('belum_jelas', 1)).toContain('Belum bisa diperkirakan')
  })
  it('words queue reasons', () => {
    expect(queueReason({ pattern_id: 'p', reason: 'akurasi_rendah', accuracy_pct: 40 })).toBe('Akurasi 40%, perlu diulang')
    expect(queueReason({ pattern_id: 'p', reason: 'belum_dicoba', accuracy_pct: 0 })).toBe('Belum dicoba')
  })
})
