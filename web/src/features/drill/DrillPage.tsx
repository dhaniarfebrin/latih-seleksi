import { useSearchParams } from 'react-router-dom'
import { DrillStart } from './components/DrillStart'
import { Question } from './components/Question'
import { Summary } from './components/Summary'
import { useDrill } from './useDrill'

export function DrillPage() {
  const [params] = useSearchParams()
  const { s, patterns, patternsError, start, answer, submit, next, reset } = useDrill()

  const attempt = s.attempt
  const labelsOf = (patternId: string) => patterns.find((p) => p.id === patternId)?.misconceptions ?? {}
  return (
    <section className="screen narrow">
      {s.screen === 'start' && <DrillStart initialPatternId={params.get('pattern') ?? ''} patterns={patterns} loadError={patternsError} onStart={start} />}
      {s.screen === 'question' && attempt && (
        <Question
          key={s.index}
          item={attempt.items[s.index]}
          total={attempt.items.length}
          labels={labelsOf(attempt.pattern_id)}
          onAnswer={(pos, key, ms) => answer(attempt.id, pos, key, ms)}
          onNext={next}
        />
      )}
      {s.screen === 'summary' && attempt && (
        <Summary
          items={attempt.items}
          labels={labelsOf(attempt.pattern_id)}
          onSubmit={() => submit(attempt.id)}
          onAgain={reset}
        />
      )}
    </section>
  )
}
