import type { Pattern } from './types'

export interface PatternGroup {
  section: string
  patterns: Pattern[]
}

// The section is the id segment before the first dot ("TIU.numerik.x" -> "TIU").
export const sectionOf = (id: string) => id.split('.')[0]

// Groups patterns by section, keeping the order of first appearance.
export function groupPatterns(patterns: Pattern[]): PatternGroup[] {
  const groups: PatternGroup[] = []
  for (const p of patterns) {
    const section = sectionOf(p.id)
    const group = groups.find((g) => g.section === section)
    if (group) group.patterns.push(p)
    else groups.push({ section, patterns: [p] })
  }
  return groups
}
