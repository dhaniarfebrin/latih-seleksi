export interface TrendPoint {
  /** Whole percent, 0 to 100. */
  value: number
  /** ISO date of the session or exam. */
  date: string
  /** Full sentence for screen readers and the hover title. */
  detail: string
}

const MIN_POINTS = 3

// Day number, plus the short month on the first and last columns only so labels stay narrow.
function labels(points: TrendPoint[]) {
  return points.map((p, i) => {
    const d = new Date(p.date)
    return i === 0 || i === points.length - 1
      ? `${d.getDate()} ${d.toLocaleDateString('id-ID', { month: 'short' })}`
      : String(d.getDate())
  })
}

// Oldest to newest, one column per point. Fewer than MIN_POINTS points say nothing about a trend, so the caller shows `empty` instead.
export function TrendBars({ points, empty }: { points: TrendPoint[]; empty: string }) {
  if (points.length < MIN_POINTS) return <p className="muted">{empty}</p>
  const text = labels(points)
  return (
    <ol className="trend">
      {points.map((p, i) => (
        <li key={`${p.date}-${i}`} title={p.detail} className={i === points.length - 1 ? 'last' : undefined}>
          <span className="num">{p.value}</span>
          <span className="trend-col" aria-hidden="true">
            <i style={{ height: `${Math.max(p.value, 2)}%` }} />
          </span>
          <span className="trend-lbl">{text[i]}</span>
          <span className="sr-only">{p.detail}</span>
        </li>
      ))}
    </ol>
  )
}
