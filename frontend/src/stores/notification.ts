import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as notificationApi from '@/api/notification'
import type { Notification } from '@/api/notification'

export const useNotificationStore = defineStore('notification', () => {
  const unreadCount = ref(0)
  const notifications = ref<Notification[]>([])
  const total = ref(0)
  let pollTimer: ReturnType<typeof setInterval> | null = null

  async function fetchUnreadCount() {
    try {
      const res = await notificationApi.getUnreadCount()
      unreadCount.value = res.data.count
    } catch (err) {
      console.error('Failed to fetch unread count:', err)
    }
  }

  async function fetchNotifications(page = 1, pageSize = 20) {
    try {
      const res = await notificationApi.getNotifications({ page, page_size: pageSize })
      const ndata = res.data
      notifications.value = ndata.list
      total.value = ndata.total
    } catch (err) {
      console.error('Failed to fetch notifications:', err)
    }
  }

  async function markAsRead(id: number | string) {
    try {
      await notificationApi.markAsRead(id)
      await fetchUnreadCount()
    } catch (err) {
      console.error('Failed to mark notification as read:', err)
    }
  }

  async function markAllAsRead() {
    try {
      await notificationApi.markAllAsRead()
      await fetchUnreadCount()
      await fetchNotifications()
    } catch (err) {
      console.error('Failed to mark all as read:', err)
    }
  }

  function startPolling(intervalMs = 30000) {
    stopPolling()
    pollTimer = setInterval(() => {
      fetchUnreadCount()
    }, intervalMs)
  }

  function stopPolling() {
    if (pollTimer) {
      clearInterval(pollTimer)
      pollTimer = null
    }
  }

  return {
    unreadCount,
    notifications,
    total,
    fetchUnreadCount,
    fetchNotifications,
    markAsRead,
    markAllAsRead,
    startPolling,
    stopPolling,
  }
})
