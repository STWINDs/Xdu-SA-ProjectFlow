import api from './index'

export interface LoginParams {
  username: string
  password: string
  captcha_id: string
  captcha_answer: string
}

export interface RegisterParams {
  username: string
  email: string
  password: string
}

export interface AuthResponse {
  access_token: string
  refresh_token: string
  user: {
    id: number
    username: string
    email: string
    role: string
  }
}

export interface CaptchaResponse {
  captcha_id: string
  captcha_image: string
}

export interface UserProfile {
  id: number
  username: string
  email: string
  role: string
  avatar_url: string
  created_at: string
}

export function login(params: LoginParams) {
  return api.post<AuthResponse>('/api/auth/login', params)
}

export function register(params: RegisterParams) {
  return api.post<UserProfile>('/api/auth/register', params)
}

export function refresh(refreshToken: string) {
  return api.post<{ access_token: string; refresh_token: string }>('/api/auth/refresh', {
    refresh_token: refreshToken,
  })
}

export function getCaptcha() {
  return api.get<CaptchaResponse>('/api/auth/captcha')
}

export function getProfile() {
  return api.get<UserProfile>('/api/auth/profile')
}
