import { renderToString } from 'react-dom/server'
import { MemoryRouter } from 'react-router-dom'
import { describe, expect, it } from 'vitest'
import { LandingContent } from './LandingContent'

const html = renderToString(
  <MemoryRouter>
    <LandingContent />
  </MemoryRouter>,
)

describe('LandingContent', () => {
  it('has exactly one h1', () => {
    expect(html.match(/<h1/g)).toHaveLength(1)
  })
  it('shows the disclaimer and links to /login', () => {
    expect(html).toContain('bukan produk resmi')
    expect(html).toContain('href="/login"')
  })
  it('speaks to the reader as kamu and labels the sample sheet', () => {
    expect(html).not.toContain('Anda')
    expect(html).toContain('Titik lemahmu')
    expect(html).toContain('Contoh')
  })
  it('names no official body', () => {
    for (const name of ['BKN', 'KemenPANRB', 'FHCI']) expect(html).not.toContain(name)
  })
})
