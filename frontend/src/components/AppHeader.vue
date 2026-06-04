<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import NotificationBell from './NotificationBell.vue'

defineProps<{
  username?: string
  avatarUrl?: string
}>()

const authStore = useAuthStore()
const router = useRouter()

function handleLogout() {
  authStore.logout()
  router.push('/login')
}
</script>

<template>
  <header class="app-header">
    <div class="header-left">
      <router-link to="/" class="logo">ProjectFlow</router-link>
    </div>

    <div class="header-center">
      <div class="search-box">
        <svg class="search-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#5f6368" stroke-width="2" stroke-linecap="round">
          <circle cx="11" cy="11" r="8" />
          <line x1="21" y1="21" x2="16.65" y2="16.65" />
        </svg>
        <input type="text" placeholder="Search tasks..." class="search-input" />
      </div>
    </div>

    <div class="header-right">
      <NotificationBell />
      <div class="user-menu">
        <div class="avatar" :title="authStore.user?.username || username">
          {{ (authStore.user?.username || username || 'U')[0].toUpperCase() }}
        </div>
        <span class="username-text">{{ authStore.user?.username || username }}</span>
        <button class="btn btn-flat btn-sm" @click="handleLogout">Sign out</button>
      </div>
    </div>
  </header>
</template>

<style scoped>
.app-header {
  height: 64px;
  background: #fff;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
  display: flex;
  align-items: center;
  padding: 0 24px;
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-left {
  display: flex;
  align-items: center;
  flex-shrink: 0;
}

.logo {
  font-size: 20px;
  font-weight: 600;
  color: #1a73e8;
  letter-spacing: -0.3px;
}

.logo:hover {
  text-decoration: none;
}

.header-center {
  flex: 1;
  display: flex;
  justify-content: center;
  padding: 0 32px;
}

.search-box {
  position: relative;
  width: 100%;
  max-width: 480px;
}

.search-icon {
  position: absolute;
  left: 14px;
  top: 50%;
  transform: translateY(-50%);
  pointer-events: none;
}

.search-input {
  width: 100%;
  height: 44px;
  padding: 0 16px 0 44px;
  font-size: 14px;
  font-family: inherit;
  color: #202124;
  background: #f1f3f4;
  border: none;
  border-radius: 22px;
  outline: none;
  transition: background 0.2s, box-shadow 0.2s;
}

.search-input:focus {
  background: #fff;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.15);
}

.search-input::placeholder {
  color: #9aa0a6;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-shrink: 0;
}

.user-menu {
  display: flex;
  align-items: center;
  gap: 8px;
}

.username-text {
  font-size: 14px;
  font-weight: 500;
  color: #202124;
}

@media (max-width: 768px) {
  .header-center {
    display: none;
  }
  .username-text {
    display: none;
  }
}
</style>
