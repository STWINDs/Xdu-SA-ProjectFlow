<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import * as authApi from '@/api/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const username = ref('')
const password = ref('')
const captchaId = ref('')
const captchaImage = ref('')
const captchaAnswer = ref('')
const error = ref('')
const loading = ref(false)

async function fetchCaptcha() {
  try {
    const res = await authApi.getCaptcha()
    captchaId.value = res.data.captcha_id
    captchaImage.value = res.data.captcha_image
  } catch (err) {
    console.error('Failed to fetch captcha:', err)
  }
}

async function handleLogin() {
  if (!username.value || !password.value) {
    error.value = 'Please enter username and password'
    return
  }
  if (!captchaAnswer.value) {
    error.value = 'Please enter the captcha'
    return
  }

  loading.value = true
  error.value = ''
  try {
    await authStore.login({
      username: username.value,
      password: password.value,
      captcha_id: captchaId.value,
      captcha_answer: captchaAnswer.value,
    })
    const redirect = (route.query.redirect as string) || '/'
    router.push(redirect)
  } catch (err: any) {
    const msg = err.response?.data?.message || err.message || 'Login failed'
    error.value = msg
    fetchCaptcha()
    captchaAnswer.value = ''
  } finally {
    loading.value = false
  }
}

onMounted(fetchCaptcha)
</script>

<template>
  <div class="auth-page">
    <div class="auth-card card">
      <div class="auth-logo">
        <h1>ProjectFlow</h1>
        <p class="auth-subtitle">Sign in to your account</p>
      </div>

      <form @submit.prevent="handleLogin" class="auth-form">
        <div v-if="error" class="form-error mb-2" style="text-align:center;font-size:13px;">{{ error }}</div>

        <div class="form-group">
          <label class="form-label">Username</label>
          <input
            v-model="username"
            type="text"
            class="form-input"
            placeholder="Enter your username"
            autocomplete="username"
          />
        </div>

        <div class="form-group">
          <label class="form-label">Password</label>
          <input
            v-model="password"
            type="password"
            class="form-input"
            placeholder="Enter your password"
            autocomplete="current-password"
          />
        </div>

        <div class="form-group">
          <label class="form-label">Captcha</label>
          <div class="captcha-row">
            <input
              v-model="captchaAnswer"
              type="text"
              class="form-input"
              placeholder="Enter captcha"
              style="flex:1"
            />
            <img
              v-if="captchaImage"
              :src="'data:image/png;base64,' + captchaImage"
              alt="Captcha"
              class="captcha-img"
              @click="fetchCaptcha"
              title="Click to refresh"
            />
          </div>
        </div>

        <button type="submit" class="btn btn-primary btn-block mt-2" :disabled="loading">
          <span v-if="loading" class="spinner" style="width:18px;height:18px;border-width:2px;"></span>
          <span v-else>Sign In</span>
        </button>
      </form>

      <div class="auth-footer">
        Don't have an account?
        <router-link to="/register">Create one</router-link>
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

.captcha-row {
  display: flex;
  gap: 12px;
  align-items: center;
}

.captcha-img {
  height: 48px;
  width: 120px;
  border-radius: 8px;
  border: 1px solid #dadce0;
  cursor: pointer;
  object-fit: contain;
  background: #f1f3f4;
}

.auth-footer {
  text-align: center;
  font-size: 13px;
  color: #5f6368;
}
</style>
