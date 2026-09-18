<template>
  <div class="register-page">
    <div class="register-card">
      <h1>Регистрация в Meetico</h1>
      <p class="subtitle">Создайте аккаунт для управления встречами</p>
      
      <form @submit.prevent="handleRegister" class="register-form">
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
          <label class="label">Имя</label>
          <input 
            v-model="displayName" 
            type="text" 
            class="input" 
            placeholder="Иван Иванов"
            required
          />
        </div>
        
        <div class="form-group">
          <label class="label">Пароль</label>
          <input 
            v-model="password" 
            type="password" 
            class="input" 
            placeholder="Минимум 8 символов"
            required
            minlength="8"
          />
        </div>
        
        <div class="form-group">
          <label class="label">Подтвердите пароль</label>
          <input 
            v-model="confirmPassword" 
            type="password" 
            class="input" 
            placeholder="••••••••"
            required
          />
        </div>
        
        <div v-if="authStore.error" class="error">
          {{ authStore.error }}
        </div>
        
        <div v-if="passwordError" class="error">
          {{ passwordError }}
        </div>
        
        <button 
          type="submit" 
          class="btn btn-primary" 
          :disabled="authStore.loading"
        >
          {{ authStore.loading ? 'Регистрация...' : 'Зарегистрироваться' }}
        </button>
      </form>
      
      <div class="footer">
        Уже есть аккаунт? 
        <router-link to="/login" class="link">Войти</router-link>
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
const displayName = ref('')
const password = ref('')
const confirmPassword = ref('')
const passwordError = ref('')

const handleRegister = async () => {
  passwordError.value = ''
  
  if (password.value !== confirmPassword.value) {
    passwordError.value = 'Пароли не совпадают'
    return
  }
  
  if (password.value.length < 8) {
    passwordError.value = 'Пароль должен быть минимум 8 символов'
    return
  }
  
  const success = await authStore.register({
    email: email.value,
    password: password.value,
    display_name: displayName.value
  })
  
  if (success) {
    router.push('/dashboard/bookings')
  }
}
</script>

<style scoped>
.register-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-bg);
  padding: var(--sp-4);
}

.register-card {
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

.register-form {
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
