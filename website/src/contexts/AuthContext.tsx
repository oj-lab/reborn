import React, { createContext, useCallback, useEffect, useState, useMemo } from 'react'
import type { ReactNode } from 'react'
import { simpleApiClient } from '@/api/simple-client'
import type { User } from '@/api/simple-client'

interface AuthContextType {
  user: User | null
  loading: boolean
  login: () => void
  logout: () => void
  fetchUser: () => Promise<void>
  isAuthenticated: boolean
}

export const AuthContext = createContext<AuthContextType | undefined>(undefined)

interface AuthProviderProps {
  children: ReactNode
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)

  const fetchUser = useCallback(async () => {
    try {
      const response = await simpleApiClient.userApi.getCurrentUser()
      setUser(response.data)
    } catch (error) {
      console.error('Failed to fetch user:', error)
      setUser(null)
    } finally {
      setLoading(false)
    }
  }, [])

  const login = useCallback(() => {
    window.location.href = '/auth/login?provider=github'
  }, [])

  const logout = useCallback(() => {
    // Clear user state
    setUser(null)
    // Redirect to logout endpoint to clear the session cookie
    window.location.href = '/auth/logout'
  }, [])

  useEffect(() => {
    fetchUser()
  }, [fetchUser])

  const value: AuthContextType = useMemo(() => ({
    user,
    loading,
    login,
    logout,
    fetchUser,
    isAuthenticated: user !== null,
  }), [user, loading, login, logout, fetchUser])

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  )
}
