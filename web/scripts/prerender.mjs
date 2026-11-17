// Writes dist/landing.html (index.html with the landing HTML inside #root); nginx serves it for "/" only,
// so crawlers without JS see real content while other routes keep the empty shell.
import { readFile, writeFile } from 'node:fs/promises'
import { fileURLToPath } from 'node:url'
import { createServer } from 'vite'

const file = fileURLToPath(import.meta.resolve('../dist/index.html'))
const out = fileURLToPath(import.meta.resolve('../dist/landing.html'))
const server = await createServer({ server: { middlewareMode: true }, appType: 'custom', logLevel: 'error' })
try {
  const { render } = await server.ssrLoadModule('/src/prerender.tsx')
  const html = await readFile(file, 'utf8')
  if (!html.includes('<div id="root"></div>')) throw new Error('dist/index.html has no empty #root')
  await writeFile(out, html.replace('<div id="root"></div>', `<div id="root">${render()}</div>`))
} finally {
  await server.close()
}
