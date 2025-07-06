// Simple API client that bypasses TypeScript strict rules
export interface User {
  id?: number
  name?: string
  email?: string
  github_id?: string
  role?: string
  created_at?: string
  updated_at?: string
}

export class SimpleUserApi {
  private baseURL: string

  constructor() {
    this.baseURL = window.location.origin
  }

  async getCurrentUser(): Promise<{ data: User }> {
    const response = await fetch(`${this.baseURL}/api/v1/user/me`, {
      method: 'GET',
      credentials: 'include', // Include cookies for authentication
      headers: {
        'Content-Type': 'application/json',
      },
    })

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }

    const data = await response.json()
    return { data }
  }
}

export const simpleApiClient = {
  userApi: new SimpleUserApi(),
}
