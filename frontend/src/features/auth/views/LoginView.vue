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
    <!-- Glow blobs -->
    <div class="blob blob-1"></div>
    <div class="blob blob-2"></div>

    <main class="login-card fade-in">
      <!-- Logo -->
      <div style="text-align:center;margin-bottom:2rem;">
        <div style="width:52px;height:52px;background:var(--color-accent);border-radius:14px;display:inline-flex;align-items:center;justify-content:center;font-size:1.5rem;margin-bottom:1rem;box-shadow:0 0 24px var(--color-accent-glow);">
          📋
        </div>
        <h1 style="font-size:1.5rem;font-weight:700;margin:0;">KRS System</h1>
        <p style="font-size:0.875rem;color:var(--color-muted);margin-top:0.375rem;">
          Sistem Pengajuan Kartu Rencana Studi
        </p>
      </div>

      <!-- Form -->
      <form @submit.prevent="onSubmit" style="display:flex;flex-direction:column;gap:1.125rem;">
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
              style="position:absolute;right:0.75rem;top:50%;transform:translateY(-50%);background:none;border:none;cursor:pointer;color:var(--color-muted);font-size:1rem;"
            >{{ showPassword ? '🙈' : '👁' }}</button>
          </div>
        </div>

        <div v-if="error" class="alert alert-error" style="margin-top:-0.25rem;">{{ error }}</div>

        <button
          id="login-submit"
          type="submit"
          class="btn btn-primary"
          :disabled="auth.loading"
          style="width:100%;padding:0.75rem;font-size:0.9375rem;margin-top:0.25rem;"
        >
          <span v-if="auth.loading" class="spinner" style="width:16px;height:16px;"></span>
          {{ auth.loading ? 'Memproses…' : 'Masuk →' }}
        </button>
      </form>

      <!-- Dev hint -->
      <div style="margin-top:1.75rem;padding-top:1.25rem;border-top:1px solid var(--color-border);">
        <p style="font-size:0.73rem;color:var(--color-muted);margin:0 0 0.5rem;text-align:center;">
          Akun Dev (password: <code>password123</code>)
        </p>
        <div style="display:grid;grid-template-columns:1fr 1fr;gap:0.375rem;">
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
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-surface);
  position: relative;
  overflow: hidden;
  padding: 1rem;
}
.blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.15;
  pointer-events: none;
}
.blob-1 {
  width: 500px; height: 500px;
  background: var(--color-accent);
  top: -150px; left: -150px;
}
.blob-2 {
  width: 400px; height: 400px;
  background: #7c3aed;
  bottom: -100px; right: -100px;
}
.login-card {
  background: var(--color-surface-card);
  border: 1px solid var(--color-border);
  border-radius: 18px;
  padding: 2.5rem 2rem;
  width: 100%;
  max-width: 400px;
  position: relative;
  z-index: 1;
  box-shadow: 0 24px 48px rgba(0,0,0,0.4);
}
.hint-btn {
  background: var(--color-surface-elevated);
  border: 1px solid var(--color-border);
  border-radius: 6px;
  padding: 0.375rem 0.5rem;
  font-size: 0.72rem;
  color: var(--color-muted);
  cursor: pointer;
  transition: all 0.15s;
}
.hint-btn:hover { background: var(--color-surface-hover); color: var(--color-text); }
</style>
