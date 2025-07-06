import React, { useEffect } from 'react'
import { useAuth } from '@/hooks/useAuth'

// This component handles the post-login redirect and refreshes user data
const AuthCallback: React.FC = () => {
  const { fetchUser } = useAuth()

  useEffect(() => {
    // Check if we just came back from OAuth login
    const urlParams = new URLSearchParams(window.location.search)
    const justLoggedIn = urlParams.has('auth_success') || 
                         window.location.pathname === '/auth/callback'

    if (justLoggedIn) {
      // Clear the URL parameters
      window.history.replaceState({}, '', window.location.pathname)
      // Refresh user data
      fetchUser()
    }
  }, [fetchUser])

  return null
}

export default AuthCallback
