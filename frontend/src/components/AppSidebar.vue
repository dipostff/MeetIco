<template>
  <div class="sidebar">
    <div class="logo">Meetico</div>
    <div class="user-info">
      <div class="avatar">{{ userInitials }}</div>
      <div class="user-details">
        <div class="name">{{ authStore.user?.display_name || 'User' }}</div>
        <div class="username">@{{ authStore.user?.username || 'username' }}</div>
      </div>
    </div>
    <nav class="nav">
      <router-link to="/dashboard/bookings" class="nav-item">
        📅 Встречи
      </router-link>
      <router-link to="/dashboard/event-types" class="nav-item">
        🔗 Мои ссылки
      </router-link>
      <router-link to="/dashboard/settings" class="nav-item">
        ⚙️ Настройки
      </router-link>
    </nav>
    <button @click="logout" class="logout-btn">
      Выйти
    </button>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const userInitials = computed(() => {
  const name = authStore.user?.display_name || 'U'
  return name.charAt(0).toUpperCase()
})

const logout = () => {
  authStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.sidebar {
  width: 240px;
  background: var(--color-surface);
  border-right: 1px solid var(--color-border);
  padding: var(--sp-6);
  display: flex;
  flex-direction: column;
  height: 100vh;
  position: fixed;
  left: 0;
  top: 0;
}

.logo {
  font-size: var(--text-xl);
  font-weight: 600;
  color: var(--color-accent);
  margin-bottom: var(--sp-8);
}

.user-info {
  display: flex;
  align-items: center;
  gap: var(--sp-3);
  margin-bottom: var(--sp-8);
  padding-bottom: var(--sp-6);
  border-bottom: 1px solid var(--color-border);
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: var(--r-full);
  background: var(--color-accent);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
}

.user-details {
  flex: 1;
}

.name {
  font-weight: 500;
  color: var(--color-text-primary);
}

.username {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
}

.nav {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: var(--sp-2);
}

.nav-item {
  padding: var(--sp-3) var(--sp-4);
  border-radius: var(--r-md);
  color: var(--color-text-secondary);
  text-decoration: none;
  transition: all var(--t-fast);
}

.nav-item:hover {
  background: var(--color-surface-2);
  color: var(--color-text-primary);
}

.nav-item.router-link-active {
  background: var(--color-accent-light);
  color: var(--color-accent);
}

.logout-btn {
  margin-top: var(--sp-6);
  padding: var(--sp-3) var(--sp-4);
  border: 1px solid var(--color-border);
  border-radius: var(--r-md);
  background: transparent;
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: all var(--t-fast);
}

.logout-btn:hover {
  background: var(--color-surface-2);
  color: var(--color-text-primary);
}
</style>
