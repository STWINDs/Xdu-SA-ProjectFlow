<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useNotificationStore } from '@/stores/notification'

const notificationStore = useNotificationStore()
const showDropdown = ref(false)
const dropdownRef = ref<HTMLElement | null>(null)

function toggleDropdown() {
  showDropdown.value = !showDropdown.value
  if (showDropdown.value) {
    notificationStore.fetchNotifications()
  }
}

function closeDropdown(e: MouseEvent) {
  if (dropdownRef.value && !dropdownRef.value.contains(e.target as Node)) {
    showDropdown.value = false
  }
}

function handleMarkAsRead(id: number) {
  notificationStore.markAsRead(id)
}

function timeAgo(dateStr: string): string {
  const date = new Date(dateStr)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffMins = Math.floor(diffMs / 60000)
  if (diffMins < 1) return 'Just now'
  if (diffMins < 60) return `${diffMins}m ago`
  const diffHours = Math.floor(diffMins / 60)
  if (diffHours < 24) return `${diffHours}h ago`
  const diffDays = Math.floor(diffHours / 24)
  return `${diffDays}d ago`
}

onMounted(() => {
  notificationStore.startPolling()
  document.addEventListener('click', closeDropdown)
})

onUnmounted(() => {
  notificationStore.stopPolling()
  document.removeEventListener('click', closeDropdown)
})
</script>

<template>
  <div class="notification-bell" ref="dropdownRef">
    <button class="bell-btn" @click="toggleDropdown" title="Notifications">
      <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#3c4043" stroke-width="2" stroke-linecap="round">
        <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9" />
        <path d="M13.73 21a2 2 0 0 1-3.46 0" />
      </svg>
      <span v-if="notificationStore.unreadCount > 0" class="badge">
        {{ notificationStore.unreadCount > 99 ? '99+' : notificationStore.unreadCount }}
      </span>
    </button>

    <div v-if="showDropdown" class="dropdown">
      <div class="dropdown-header">
        <span>Notifications</span>
        <button
          v-if="notificationStore.notifications.length > 0"
          class="btn btn-flat btn-sm"
          @click="notificationStore.markAllAsRead()"
        >
          Mark all read
        </button>
      </div>

      <div class="dropdown-body">
        <div v-if="notificationStore.notifications.length === 0" class="empty-state">
          No notifications yet
        </div>

        <div
          v-for="n in notificationStore.notifications"
          :key="n.id"
          class="notification-item"
          :class="{ unread: !n.is_read }"
          @click="handleMarkAsRead(n.id)"
        >
          <div class="notification-dot" v-if="!n.is_read"></div>
          <div class="notification-content">
            <div class="notification-title">{{ n.title }}</div>
            <div class="notification-text">{{ n.content }}</div>
            <div class="notification-time">{{ timeAgo(n.created_at) }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.notification-bell {
  position: relative;
}

.bell-btn {
  width: 40px;
  height: 40px;
  border: none;
  background: transparent;
  border-radius: 50%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  transition: background 0.2s;
}

.bell-btn:hover {
  background: #f1f3f4;
}

.badge {
  position: absolute;
  top: 2px;
  right: 2px;
  min-width: 18px;
  height: 18px;
  padding: 0 5px;
  font-size: 10px;
  font-weight: 600;
  color: #fff;
  background: #ea4335;
  border-radius: 9px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.dropdown {
  position: absolute;
  top: 48px;
  right: 0;
  width: 360px;
  max-height: 480px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12);
  overflow: hidden;
  z-index: 200;
}

.dropdown-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 16px;
  border-bottom: 1px solid #e8eaed;
  font-size: 15px;
  font-weight: 600;
  color: #202124;
}

.dropdown-body {
  max-height: 400px;
  overflow-y: auto;
}

.empty-state {
  padding: 32px 16px;
  text-align: center;
  color: #9aa0a6;
  font-size: 14px;
}

.notification-item {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 12px 16px;
  cursor: pointer;
  transition: background 0.15s;
  border-bottom: 1px solid #f1f3f4;
}

.notification-item:hover {
  background: #f6f8fc;
}

.notification-item.unread {
  background: #f6f8fc;
}

.notification-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #1a73e8;
  margin-top: 6px;
  flex-shrink: 0;
}

.notification-content {
  flex: 1;
  min-width: 0;
}

.notification-title {
  font-size: 13px;
  font-weight: 600;
  color: #202124;
  margin-bottom: 2px;
}

.notification-text {
  font-size: 13px;
  color: #5f6368;
  line-height: 1.4;
  margin-bottom: 4px;
}

.notification-time {
  font-size: 11px;
  color: #9aa0a6;
}
</style>
