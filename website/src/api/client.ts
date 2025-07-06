import { UserApi } from '@/api/api'
import { Configuration } from '@/api/configuration'

// Get authentication token from cookies
const getAuthToken = (): string | null => {
  // Since the backend sets the login session in cookies, 
  // the browser will automatically send it with requests
  // For now, we don't need to manually handle the token
  // as the backend middleware will check the cookie
  return null
}

// Create API client with configuration
const createApiClient = () => {
  const token = getAuthToken()
  const config = new Configuration({
    basePath: window.location.origin,
    accessToken: token || undefined,
  })
  
  return {
    userApi: new UserApi(config),
  }
}

export const apiClient = createApiClient()
