// typescript-eslint 8.x rejects TypeScript 7 (peer <6.1.0). The project compiles with TS 7, so ESLint alone
// is pointed at the side-by-side `typescript-6` alias. Preloaded by the `lint` script only; tsc/vite are unaffected.
const Module = require('node:module')
const resolve = Module._resolveFilename
Module._resolveFilename = function (request, ...rest) {
  return resolve.call(this, request === 'typescript' ? 'typescript-6' : request, ...rest)
}
