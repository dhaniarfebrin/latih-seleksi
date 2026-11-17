import { call } from '@/shared/api/client'

// Deletes the caller's attempts and user row on the server (idempotent).
export const deleteMyData = () => call<void>('DELETE', '/v1/me')

// Everything the server holds about the caller, as parsed JSON.
export const exportMyData = () => call<unknown>('GET', '/v1/me/export')
