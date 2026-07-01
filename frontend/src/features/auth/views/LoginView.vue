<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/features/auth/stores/auth.store'

const router = useRouter()
const auth = useAuthStore()

const email = ref('')
const password = ref('')
const error = ref('')
const showPassword = ref(false)

async function onSubmit() {
  error.value = ''
  try {
    await auth.login({ email: email.value, password: password.value })
    // Redirect berdasarkan role
    switch (auth.role) {
      case 'MAHASISWA': await router.push('/mahasiswa/krs'); break
      case 'DOSEN':     await router.push('/dosen/persetujuan'); break
      case 'STAFF':     await router.push('/staff/persetujuan'); break
      case 'ADMIN':     await router.push('/admin/dashboard'); break
      default:          await router.push('/')
    }
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Login gagal, periksa email & password'
  }
}
</script>

<template>
  <div class="login-bg">
    <main class="login-card fade-in">
      <!-- Logo & Title -->
      <div style="text-align:center; margin-bottom:2.5rem;">
        <div class="login-logo-box">
          <svg width="22" height="22" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M16.5 8.25V6a2.25 2.25 0 0 0-2.25-2.25H6A2.25 2.25 0 0 0 3.75 6v12A2.25 2.25 0 0 0 6 20.25h12A2.25 2.25 0 0 0 20.25 18V9.75A2.25 2.25 0 0 0 18 7.5h-1.5ZM13.5 3.75v3c0 .621.504 1.125 1.125 1.125h3m-9 3.75h4.5m-4.5 3h7.5"/>
          </svg>
        </div>
        <h1 style="font-size:1.5rem; font-weight:700; margin:0; text-transform:uppercase; letter-spacing:0.05em;">KRS System</h1>
        <p style="font-size:0.8125rem; color:var(--color-text-muted); margin-top:0.375rem; letter-spacing:0.02em;">
          Sistem Pengajuan Kartu Rencana Studi
        </p>
      </div>

      <!-- Form -->
      <form @submit.prevent="onSubmit" style="display:flex; flex-direction:column; gap:1.25rem;">
        <div class="form-field">
          <label class="form-label">Email</label>
          <input
            v-model="email"
            id="login-email"
            type="email"
            required
            autocomplete="email"
            class="form-input"
            placeholder="email@kampus.ac.id"
          />
        </div>

        <div class="form-field">
          <label class="form-label">Password</label>
          <div style="position:relative;">
            <input
              v-model="password"
              id="login-password"
              :type="showPassword ? 'text' : 'password'"
              required
              autocomplete="current-password"
              class="form-input"
              placeholder="••••••••"
              style="padding-right:3rem;"
            />
            <button
              type="button"
              @click="showPassword = !showPassword"
              style="position:absolute; right:0.875rem; top:50%; transform:translateY(-50%); background:none; border:none; cursor:pointer; color:var(--color-text-muted); display:flex; align-items:center;"
            >
              <svg v-if="showPassword" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" d="M3.98 8.223A10.477 10.477 0 0 0 1.934 12C3.226 16.338 7.244 19.5 12 19.5c.993 0 1.953-.138 2.863-.395M21 21l-3.486-3.486m0 0a9 9 0 0 0-12-12L3 3m14.514 14.514a3 3 0 0 1-4.243-4.243m0 0L9 9m11.12 3a10.43 10.43 0 0 0-1.896-3.78"/>
              </svg>
              <svg v-else width="18" height="18" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" d="M2.036 12.322a1.012 1.012 0 0 1 0-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178Z"/>
                <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"/>
              </svg>
            </button>
          </div>
        </div>

        <div v-if="error" class="alert alert-error" style="font-size: 0.8125rem;">
          <svg width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24" style="flex-shrink:0;">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m9-.75a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 3.75h.008v.008H12v-.008Z"/>
          </svg>
          {{ error }}
        </div>

        <button
          id="login-submit"
          type="submit"
          class="btn btn-primary"
          :disabled="auth.loading"
          style="width:100%; padding:0.75rem;"
        >
          <span v-if="auth.loading" class="spinner" style="margin-right:0.5rem;"></span>
          {{ auth.loading ? 'Memproses' : 'Masuk' }}
        </button>
      </form>

      <!-- Dev hint -->
      <div style="margin-top:2.25rem; padding-top:1.5rem; border-top:1.5px solid var(--color-border-subtle);">
        <p style="font-size:0.75rem; color:var(--color-text-muted); margin:0 0 0.75rem; text-align:center; font-family:var(--font-mono);">
          Akun Dev (password: password123)
        </p>
        <div style="display:grid; grid-template-columns:1fr 1fr; gap:0.5rem;">
          <button class="hint-btn" @click="email='admin@kampus.ac.id';password='password123'">Admin</button>
          <button class="hint-btn" @click="email='andi@student.ac.id';password='password123'">Mahasiswa</button>
          <button class="hint-btn" @click="email='ahmad@kampus.ac.id';password='password123'">Dosen</button>
          <button class="hint-btn" @click="email='staff@kampus.ac.id';password='password123'">Staff</button>
        </div>
      </div>
    </main>
  </div>
</template>

<style scoped>
.login-bg {
  min-height: 100dvh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-bg);
  position: relative;
  overflow: hidden;
  padding: 1.5rem;
}
.login-card {
  background: var(--color-surface);
  border: 2px solid var(--color-border);
  border-radius: var(--border-radius);
  padding: 3rem 2.5rem;
  width: 100%;
  max-width: 420px;
  position: relative;
  z-index: 1;
  box-shadow: 0 10px 30px rgba(0,0,0,0.08);
}
.login-logo-box {
  width: 48px;
  height: 48px;
  background: var(--color-accent);
  color: #ffffff;
  border-radius: var(--border-radius);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 1.25rem;
}
.hint-btn {
  background: var(--color-surface);
  border: 1.5px solid var(--color-border);
  border-radius: var(--border-radius);
  padding: 0.5rem;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--color-text);
  cursor: pointer;
  transition: all 150ms ease;
}
.hint-btn:hover {
  background: var(--color-surface-hover);
  border-color: var(--color-accent);
  color: var(--color-accent);
}
</style>
