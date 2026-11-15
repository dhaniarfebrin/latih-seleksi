// Fetch wrapper for the Go API: JSON in and out, ApiError on failure, the
// Authorization header from the signed-in user, and the server clock offset.
import { getToken } from '@/features/auth/auth'

const BASE = import.meta.env.VITE_API_URL ?? ''

export class ApiError extends Error {
  constructor(
    readonly status: number,
    message: string,
  ) {
    super(message)
  }
}

// Server clock minus client clock in ms, from the Date header of the latest
// response (1 s resolution). The exam timer uses it instead of the raw client clock.
let serverOffsetMs = 0
export const getServerOffset = () => serverOffsetMs

export async function call<T>(method: string, path: string, body?: unknown): Promise<T> {
  const headers: Record<string, string> = {}
  if (body !== undefined) headers['Content-Type'] = 'application/json'
  const token = await getToken()
  if (token) headers.Authorization = `Bearer ${token}`
  const res = await fetch(BASE + path, {
    method,
    headers,
    body: body === undefined ? undefined : JSON.stringify(body),
  })
  const date = Date.parse(res.headers.get('Date') ?? '')
  if (!Number.isNaN(date)) serverOffsetMs = date - Date.now()
  if (!res.ok) {
    let msg = `HTTP ${res.status}`
    try {
      const j = (await res.json()) as { error?: string }
      if (j.error) msg = j.error
    } catch {
      // keep the status message
    }
    if (res.status === 401) msg = 'Sesi login berakhir, silakan login lagi'
    else if (res.status === 403 && msg === 'email not verified') msg = 'Email belum diverifikasi. Cek kotak masuk lalu muat ulang status.'
    throw new ApiError(res.status, msg)
  }
  return (await res.json()) as T
}
