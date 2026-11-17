import { describe, expect, it } from 'vitest'
import { groupPatterns, sectionOf } from './patterns'

const p = (id: string) => ({ id, label: id, misconceptions: {} })

describe('groupPatterns', () => {
  it('groups by the segment before the first dot, in first-seen order', () => {
    const groups = groupPatterns([p('TIU.numerik.a'), p('TIU.verbal.b'), p('TWK.integritas'), p('TWK.bela-negara')])
    expect(groups.map((g) => [g.section, g.patterns.length])).toEqual([
      ['TIU', 2],
      ['TWK', 2],
    ])
  })

  it('returns no groups for no patterns and uses the whole id without a dot', () => {
    expect(groupPatterns([])).toEqual([])
    expect(sectionOf('lainnya')).toBe('lainnya')
  })
})
