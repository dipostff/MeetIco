<template>
  <div class="settings">
    <OnboardingBanner />
    
    <h1>Настройки</h1>
    
    <div class="settings-sections">
      <!-- Public Profile -->
      <div class="settings-section">
        <h2>Публичный профиль</h2>
        <form @submit.prevent="saveProfile" class="settings-form">
          <div class="form-group">
            <label class="label">Фото (URL)</label>
            <input 
              v-model="profileForm.photo_url" 
              type="url" 
              class="input" 
              placeholder="https://example.com/photo.jpg"
            />
          </div>
          <div class="form-group">
            <label class="label">Имя</label>
            <input 
              v-model="profileForm.display_name" 
              type="text" 
              class="input" 
              placeholder="Иван Иванов"
              required
            />
          </div>
          <div class="form-group">
            <label class="label">Username</label>
            <input 
              v-model="profileForm.username" 
              type="text" 
              class="input" 
              placeholder="ivanov"
              required
              minlength="3"
              maxlength="30"
            />
            <span class="hint">Только буквы, цифры, _ и -. 3-30 символов</span>
          </div>
          <div class="form-group">
            <label class="label">О себе</label>
            <textarea 
              v-model="profileForm.bio" 
              class="input" 
              placeholder="Кратко о себе"
              rows="3"
            ></textarea>
          </div>
          <button type="submit" class="btn btn-primary" :disabled="authStore.loading">
            {{ authStore.loading ? 'Сохранение...' : 'Сохранить профиль' }}
          </button>
        </form>
      </div>

      <!-- Notifications -->
      <div class="settings-section">
        <h2>Уведомления</h2>
        <form @submit.prevent="saveNotifications" class="settings-form">
          <div class="form-group">
            <label class="label">Telegram Chat ID</label>
            <input 
              v-model="notificationsForm.telegram_chat_id" 
              type="text" 
              class="input" 
              placeholder="123456789"
            />
            <span class="hint">Напишите @userinfobot в Telegram, чтобы узнать ваш Chat ID</span>
          </div>
          <button type="submit" class="btn btn-primary" :disabled="authStore.loading">
            {{ authStore.loading ? 'Сохранение...' : 'Сохранить уведомления' }}
          </button>
        </form>
      </div>

      <!-- Schedule -->
      <div class="settings-section">
        <h2>Расписание</h2>
        <form @submit.prevent="saveSchedule" class="settings-form">
          <div class="form-group">
            <label class="label">Часовой пояс</label>
            <select v-model="scheduleForm.timezone" class="input" required>
              <option value="Europe/Moscow">Москва (Europe/Moscow)</option>
              <option value="Europe/Kiev">Киев (Europe/Kiev)</option>
              <option value="Europe/Berlin">Берлин (Europe/Berlin)</option>
              <option value="Europe/London">Лондон (Europe/London)</option>
              <option value="America/New_York">Нью-Йорк (America/New_York)</option>
              <option value="America/Los_Angeles">Лос-Анджелес (America/Los_Angeles)</option>
              <option value="Asia/Tokyo">Токио (Asia/Tokyo)</option>
            </select>
          </div>
          
          <div class="schedule-table">
            <div v-for="day in days" :key="day.key" class="schedule-row">
              <div class="day-label">{{ day.label }}</div>
              <div class="day-intervals">
                <div 
                  v-for="(interval, index) in scheduleForm.weekly_hours[day.key]" 
                  :key="index"
                  class="interval-inputs"
                >
                  <input 
                    v-model="interval[0]" 
                    type="time" 
                    class="input input-sm"
                  />
                  <span>-</span>
                  <input 
                    v-model="interval[1]" 
                    type="time" 
                    class="input input-sm"
                  />
                  <button 
                    type="button" 
                    class="btn btn-ghost btn-sm"
                    @click="removeInterval(day.key, index)"
                  >
                    ✕
                  </button>
                </div>
                <button 
                  type="button" 
                  class="btn btn-secondary btn-sm"
                  @click="addInterval(day.key)"
                >
                  + Добавить интервал
                </button>
              </div>
            </div>
          </div>
          
          <button type="submit" class="btn btn-primary" :disabled="scheduleStore.loading">
            {{ scheduleStore.loading ? 'Сохранение...' : 'Сохранить расписание' }}
          </button>
        </form>
      </div>

      <!-- Security -->
      <div class="settings-section">
        <h2>Безопасность</h2>
        <form @submit.prevent="changePassword" class="settings-form">
          <div class="form-group">
            <label class="label">Текущий пароль</label>
            <input 
              v-model="passwordForm.old_password" 
              type="password" 
              class="input" 
              placeholder="••••••••"
              required
            />
          </div>
          <div class="form-group">
            <label class="label">Новый пароль</label>
            <input 
              v-model="passwordForm.new_password" 
              type="password" 
              class="input" 
              placeholder="Минимум 8 символов"
              required
              minlength="8"
            />
          </div>
          <div class="form-group">
            <label class="label">Подтвердите новый пароль</label>
            <input 
              v-model="passwordForm.confirm_password" 
              type="password" 
              class="input" 
              placeholder="••••••••"
              required
            />
          </div>
          <button type="submit" class="btn btn-primary" :disabled="authStore.loading">
            {{ authStore.loading ? 'Изменение...' : 'Изменить пароль' }}
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useScheduleStore } from '@/stores/schedule'
import OnboardingBanner from '@/components/OnboardingBanner.vue'

const authStore = useAuthStore()
const scheduleStore = useScheduleStore()

const profileForm = reactive({
  display_name: '',
  username: '',
  bio: '',
  photo_url: ''
})

const notificationsForm = reactive({
  telegram_chat_id: ''
})

const scheduleForm = reactive({
  timezone: 'Europe/Moscow',
  weekly_hours: {
    mon: [],
    tue: [],
    wed: [],
    thu: [],
    fri: [],
    sat: [],
    sun: []
  }
})

const passwordForm = reactive({
  old_password: '',
  new_password: '',
  confirm_password: ''
})

const days = [
  { key: 'mon', label: 'Понедельник' },
  { key: 'tue', label: 'Вторник' },
  { key: 'wed', label: 'Среда' },
  { key: 'thu', label: 'Четверг' },
  { key: 'fri', label: 'Пятница' },
  { key: 'sat', label: 'Суббота' },
  { key: 'sun', label: 'Воскресенье' }
]

const addInterval = (day) => {
  if (!scheduleForm.weekly_hours[day]) {
    scheduleForm.weekly_hours[day] = []
  }
  scheduleForm.weekly_hours[day].push(['09:00', '18:00'])
}

const removeInterval = (day, index) => {
  scheduleForm.weekly_hours[day].splice(index, 1)
}

const saveProfile = async () => {
  const username = profileForm.username?.trim() || ''
  if (!/^[a-z0-9_-]{3,30}$/.test(username)) {
    alert('Username: 3–30 символов, только латиница, цифры, дефис и подчёркивание')
    return
  }
  profileForm.username = username

  const success = await authStore.updateProfile(profileForm)
  if (success) {
    alert('Профиль сохранен!')
  } else {
    alert('Не удалось сохранить профиль')
  }
}

const saveNotifications = async () => {
  const success = await authStore.updateProfile(notificationsForm)
  if (success) {
    alert('Уведомления сохранены!')
  } else {
    alert('Не удалось сохранить уведомления')
  }
}

const saveSchedule = async () => {
  const success = await scheduleStore.updateSchedule(scheduleForm)
  if (success) {
    alert('Расписание сохранено!')
  } else {
    alert('Не удалось сохранить расписание')
  }
}

const changePassword = async () => {
  if (passwordForm.new_password !== passwordForm.confirm_password) {
    alert('Пароли не совпадают')
    return
  }
  
  if (passwordForm.new_password.length < 8) {
    alert('Пароль должен быть минимум 8 символов')
    return
  }
  
  const success = await authStore.changePassword({
    old_password: passwordForm.old_password,
    new_password: passwordForm.new_password
  })
  
  if (success) {
    alert('Пароль изменен!')
    passwordForm.old_password = ''
    passwordForm.new_password = ''
    passwordForm.confirm_password = ''
  } else {
    alert('Не удалось изменить пароль')
  }
}

onMounted(async () => {
  await authStore.fetchUser()
  
  if (authStore.user) {
    profileForm.display_name = authStore.user.display_name || ''
    profileForm.username = authStore.user.username || ''
    profileForm.bio = authStore.user.bio || ''
    profileForm.photo_url = authStore.user.photo_url || ''
    notificationsForm.telegram_chat_id = authStore.user.telegram_chat_id || ''
  }
  
  await scheduleStore.fetchSchedule()
  
  if (scheduleStore.schedule) {
    scheduleForm.timezone = scheduleStore.schedule.timezone || 'Europe/Moscow'
    scheduleForm.weekly_hours = scheduleStore.schedule.weekly_hours || {
      mon: [],
      tue: [],
      wed: [],
      thu: [],
      fri: [],
      sat: [],
      sun: []
    }
  }
})
</script>

<style scoped>
.settings {
  max-width: 800px;
}

h1 {
  font-size: var(--text-2xl);
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: var(--sp-8);
}

.settings-sections {
  display: flex;
  flex-direction: column;
  gap: var(--sp-8);
}

.settings-section {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--r-lg);
  padding: var(--sp-8);
}

.settings-section h2 {
  font-size: var(--text-xl);
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: var(--sp-6);
}

.settings-form {
  display: flex;
  flex-direction: column;
  gap: var(--sp-4);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--sp-2);
}

.hint {
  font-size: var(--text-sm);
  color: var(--color-text-muted);
}

.input-sm {
  padding: var(--sp-2) var(--sp-3);
  font-size: var(--text-sm);
}

.schedule-table {
  display: flex;
  flex-direction: column;
  gap: var(--sp-3);
  margin-bottom: var(--sp-4);
}

.schedule-row {
  display: flex;
  gap: var(--sp-4);
  align-items: flex-start;
}

.day-label {
  min-width: 120px;
  font-weight: 500;
  color: var(--color-text-primary);
  padding-top: var(--sp-2);
}

.day-intervals {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: var(--sp-2);
}

.interval-inputs {
  display: flex;
  align-items: center;
  gap: var(--sp-2);
}

.interval-inputs span {
  color: var(--color-text-secondary);
}

.btn-sm {
  padding: var(--sp-2) var(--sp-3);
  font-size: var(--text-sm);
}
</style>
