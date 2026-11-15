import { call } from '@/shared/api/client'

// Deletes the caller's attempts and user row on the server (idempotent).
export const deleteMyData = () => call<void>('DELETE', '/v1/me')
