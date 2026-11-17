import { renderToString } from 'react-dom/server'
import { MemoryRouter } from 'react-router-dom'
import { LandingContent } from '@/features/landing/LandingContent'

export function render() {
  return renderToString(
    <MemoryRouter>
      <LandingContent />
    </MemoryRouter>,
  )
}
