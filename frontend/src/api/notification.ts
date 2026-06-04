import api from './index'

export interface Notification {
  id: number
  type: string
  title: string
  content: string
  is_read: boolean
  created_at: string
}

export interface NotificationListResponse {
  list: Notification[]
  total: number
  page: number
  page_size: number
}

export function getNotifications(params?: { page?: number; page_size?: number }) {
  return api.get<NotificationListResponse>('/api/notifications', { params })
}

export function getUnreadCount() {
  return api.get<{ count: number }>('/api/notifications/unread-count')
}

export function markAsRead(id: number | string) {
  return api.put(`/api/notifications/${id}/read`)
}

export function markAllAsRead() {
  return api.put('/api/notifications/read-all')
}
