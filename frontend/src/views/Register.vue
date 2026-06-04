<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const username = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const error = ref('')
const loading = ref(false)

async function handleRegister() {
  error.value = ''

  if (!username.value || !email.value || !password.value) {
    error.value = 'All fields are required'
    return
  }

  if (password.value !== confirmPassword.value) {
    error.value = 'Passwords do not match'
    return
  }

  if (password.value.length < 6) {
    error.value = 'Password must be at least 6 characters'
    return
  }

  loading.value = true
  try {
    await authStore.register({
      username: username.value,
      email: email.value,
      password: password.value,
    })
    router.push('/login')
  } catch (err: any) {
    const msg = err.response?.data?.message || err.message || 'Registration failed'
    error.value = msg
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="auth-page">
    <div class="auth-card card">
      <div class="auth-logo">
        <h1>ProjectFlow</h1>
        <p class="auth-subtitle">Create your account</p>
      </div>

      <form @submit.prevent="handleRegister" class="auth-form">
        <div v-if="error" class="form-error mb-2" style="text-align:center;font-size:13px;">{{ error }}</div>

        <div class="form-group">
          <label class="form-label">Username</label>
          <input
            v-model="username"
            type="text"
            class="form-input"
            placeholder="Choose a username"
          />
        </div>

        <div class="form-group">
          <label class="form-label">Email</label>
          <input
            v-model="email"
            type="email"
            class="form-input"
            placeholder="Enter your email"
          />
        </div>

        <div class="form-group">
          <label class="form-label">Password</label>
          <input
            v-model="password"
            type="password"
            class="form-input"
            placeholder="Create a password (min 6 characters)"
          />
        </div>

        <div class="form-group">
          <label class="form-label">Confirm Password</label>
          <input
            v-model="confirmPassword"
            type="password"
            class="form-input"
            placeholder="Confirm your password"
          />
        </div>

        <button type="submit" class="btn btn-primary btn-block mt-2" :disabled="loading">
          <span v-if="loading" class="spinner" style="width:18px;height:18px;border-width:2px;"></span>
          <span v-else>Create Account</span>
        </button>
      </form>

      <div class="auth-footer">
        Already have an account?
        <router-link to="/login">Sign in</router-link>
      </div>
    </div>
  </div>
</template>

<style scoped>
.auth-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
  background: #f6f8fc;
}

.auth-card {
  width: 100%;
  max-width: 420px;
  padding: 40px;
}

.auth-logo {
  text-align: center;
  margin-bottom: 32px;
}

.auth-logo h1 {
  font-size: 26px;
  font-weight: 700;
  color: #1a73e8;
  letter-spacing: -0.5px;
}

.auth-subtitle {
  font-size: 14px;
  color: #5f6368;
  margin-top: 6px;
}

.auth-form {
  margin-bottom: 20px;
}

.auth-footer {
  text-align: center;
  font-size: 13px;
  color: #5f6368;
}
</style>
