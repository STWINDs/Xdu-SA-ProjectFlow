<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import AppHeader from '@/components/AppHeader.vue'

const authStore = useAuthStore()
const router = useRouter()
const loading = ref(false)

const username = authStore.user?.username || 'User'

onMounted(async () => {
  if (authStore.isLoggedIn && !authStore.user) {
    await authStore.fetchProfile()
  }
})

function handleLogout() {
  authStore.logout()
  router.push('/login')
}
</script>

<template>
  <div class="settings-page">
    <AppHeader :username="username" />

    <div class="page-content">
      <h2 style="font-size:22px;font-weight:600;margin-bottom:24px">Settings</h2>

      <!-- Profile card -->
      <div class="card mb-3">
        <h3 style="font-size:16px;font-weight:600;margin-bottom:16px">Profile</h3>
        <div class="flex gap-2">
          <div class="avatar" style="width:64px;height:64px;font-size:24px">
            {{ authStore.user?.username?.[0]?.toUpperCase() || 'U' }}
          </div>
          <div>
            <div style="font-size:16px;font-weight:500;margin-bottom:4px">
              {{ authStore.user?.username || 'Unknown' }}
            </div>
            <div style="font-size:13px;color:#5f6368;margin-bottom:2px">
              {{ authStore.user?.email || 'No email' }}
            </div>
            <span class="tag tag-info">{{ authStore.user?.role || 'member' }}</span>
          </div>
        </div>
      </div>

      <!-- Account card -->
      <div class="card mb-3">
        <h3 style="font-size:16px;font-weight:600;margin-bottom:16px">Account</h3>

        <div class="settings-row flex-between">
          <div>
            <div style="font-size:14px;font-weight:500">Change Password</div>
            <div style="font-size:12px;color:#5f6368">Update your account password</div>
          </div>
          <button class="btn btn-secondary btn-sm" disabled>Coming Soon</button>
        </div>

        <div class="settings-row flex-between">
          <div>
            <div style="font-size:14px;font-weight:500">Notifications</div>
            <div style="font-size:12px;color:#5f6368">Manage notification preferences</div>
          </div>
          <button class="btn btn-secondary btn-sm" disabled>Coming Soon</button>
        </div>
      </div>

      <!-- Danger zone -->
      <div class="card" style="border:1px solid #fce8e6">
        <h3 style="font-size:16px;font-weight:600;color:#ea4335;margin-bottom:16px">Danger Zone</h3>
        <p style="font-size:13px;color:#5f6368;margin-bottom:16px">
          Sign out of your account. You can sign back in at any time.
        </p>
        <button class="btn btn-danger btn-sm" @click="handleLogout">Sign Out</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.settings-page {
  min-height: 100vh;
  background: #f6f8fc;
}

.page-content {
  max-width: 680px;
  margin: 0 auto;
  padding: 24px;
}

.settings-row {
  padding: 14px 0;
  border-bottom: 1px solid #f1f3f4;
}

.settings-row:last-child {
  border-bottom: none;
}
</style>
