import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as authApi from '@/api/auth'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<{
    id: number
    username: string
    email: string
    role: string
    avatar_url: string
  } | null>(null)

  const isLoggedIn = ref(!!localStorage.getItem('access_token'))

  function setAuth(data: { access_token: string; refresh_token: string; user: any }) {
    localStorage.setItem('access_token', data.access_token)
    localStorage.setItem('refresh_token', data.refresh_token)
    isLoggedIn.value = true
    user.value = {
      id: data.user.id,
      username: data.user.username,
      email: data.user.email,
      role: data.user.role,
      avatar_url: (data.user as any).avatar_url || '',
    }
  }

  async function login(params: authApi.LoginParams) {
    const res = await authApi.login(params)
    const data = res.data
    setAuth(data)
    return data
  }

  async function register(params: authApi.RegisterParams) {
    const res = await authApi.register(params)
    return res.data
  }

  function logout() {
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    isLoggedIn.value = false
    user.value = null
  }

  async function fetchProfile() {
    try {
      const res = await authApi.getProfile()
      const data = res.data
      user.value = {
        id: data.id,
        username: data.username,
        email: data.email,
        role: data.role,
        avatar_url: data.avatar_url || '',
      }
    } catch (err) {
      console.error('Failed to fetch profile:', err)
    }
  }

  // Restore user state from localStorage on app init
  function restoreFromStorage() {
    const token = localStorage.getItem('access_token')
    if (token) {
      // User is logged in, profile will be fetched later
      fetchProfile()
    }
  }

  return {
    user,
    isLoggedIn,
    login,
    register,
    logout,
    fetchProfile,
    restoreFromStorage,
  }
})
