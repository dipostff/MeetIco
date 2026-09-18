<template>
  <div class="login-page">
    <div class="login-card">
      <h1>Вход в Meetico</h1>
      <p class="subtitle">Войдите, чтобы управлять своими встречами</p>
      
      <form @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <label class="label">Email</label>
          <input 
            v-model="email" 
            type="email" 
            class="input" 
            placeholder="your@email.com"
            required
          />
        </div>
        
        <div class="form-group">
          <label class="label">Пароль</label>
          <input 
            v-model="password" 
            type="password" 
            class="input" 
            placeholder="••••••••"
            required
          />
        </div>
        
        <div v-if="authStore.error" class="error">
          {{ authStore.error }}
        </div>
        
        <button 
          type="submit" 
          class="btn btn-primary" 
          :disabled="authStore.loading"
        >
          {{ authStore.loading ? 'Вход...' : 'Войти' }}
        </button>
      </form>
      
      <div class="footer">
        Нет аккаунта? 
        <router-link to="/register" class="link">Зарегистрироваться</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')

const handleLogin = async () => {
  const success = await authStore.login({
    email: email.value,
    password: password.value
  })
  
  if (success) {
    router.push('/dashboard/bookings')
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-bg);
  padding: var(--sp-4);
}

.login-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--r-lg);
  padding: var(--sp-12);
  width: 100%;
  max-width: 400px;
  box-shadow: var(--shadow-md);
}

h1 {
  font-size: var(--text-2xl);
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: var(--sp-2);
}

.subtitle {
  color: var(--color-text-secondary);
  margin-bottom: var(--sp-8);
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: var(--sp-6);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--sp-2);
}

.error {
  padding: var(--sp-3);
  background: #FEE2E2;
  color: var(--color-error);
  border-radius: var(--r-md);
  font-size: var(--text-sm);
}

.btn {
  width: 100%;
  justify-content: center;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.footer {
  margin-top: var(--sp-8);
  text-align: center;
  color: var(--color-text-secondary);
  font-size: var(--text-sm);
}

.link {
  color: var(--color-accent);
  text-decoration: none;
  font-weight: 500;
}

.link:hover {
  text-decoration: underline;
}
</style>
