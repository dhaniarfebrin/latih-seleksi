import { call } from '@/shared/api/client'
import type { Dashboard } from './types'

export const getDashboard = () => call<Dashboard>('GET', '/v1/me/dashboard')
