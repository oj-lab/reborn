import { createContext } from 'react'
import type { UserpbUser } from '@/api/api'

export interface AuthContextType {
  user: UserpbUser | null
  loading: boolean
  login: () => void
  logout: () => void
  fetchUser: () => Promise<void>
  isAuthenticated: boolean
}

export const AuthContext = createContext<AuthContextType | undefined>(undefined)
