<template>
  <div class="booking-page">
    <div class="booking-card">
      <!-- Recruiter Info -->
      <div class="recruiter-info">
        <div class="avatar">{{ recruiterInitials }}</div>
        <div class="recruiter-details">
          <h2>{{ publicInfo.recruiter?.display_name || 'Recruiter' }}</h2>
          <p class="bio">{{ publicInfo.recruiter?.bio || '' }}</p>
          <div class="meta">
            <span class="meta-item">⏱ {{ publicInfo.event_type?.duration_min || 30 }} мин</span>
            <span class="meta-item">📍 {{ publicInfo.event_type?.location_type || 'Zoom' }}</span>
          </div>
        </div>
      </div>

      <!-- Date Selection -->
      <div v-if="!bookingSuccess" class="date-selection">
        <h3>Выберите дату</h3>
        <div class="dates-scroll">
          <button 
            v-for="date in availableDates" 
            :key="date.date"
            class="date-btn"
            :class="{ 
              active: selectedDate === date.date,
              disabled: !date.hasSlots 
            }"
            @click="selectDate(date.date)"
            :disabled="!date.hasSlots"
          >
            <div class="date-day">{{ date.day }}</div>
            <div class="date-number">{{ date.number }}</div>
          </button>
        </div>
      </div>

      <!-- Time Slots -->
      <div v-if="!bookingSuccess && selectedDate" class="time-slots">
        <h3>Доступное время — {{ selectedDateFormatted }}</h3>
        <div v-if="loadingSlots" class="slots-loading">
          <div class="skeleton" style="height: 40px; width: 100px; margin-right: var(--sp-2);"></div>
          <div class="skeleton" style="height: 40px; width: 100px; margin-right: var(--sp-2);"></div>
          <div class="skeleton" style="height: 40px; width: 100px;"></div>
        </div>
        <div v-else class="slots-grid">
          <button 
            v-for="slot in slots" 
            :key="slot"
            class="slot-btn"
            :class="{ active: selectedSlot === slot }"
            @click="selectedSlot = slot"
          >
            {{ slot }}
          </button>
        </div>
        <div v-if="slots.length === 0 && !loadingSlots" class="no-slots">
          Нет доступных слотов на эту дату
        </div>
      </div>

      <!-- Booking Form -->
      <div v-if="!bookingSuccess && selectedSlot" class="booking-form">
        <h3>Ваши данные</h3>
        <div class="form-group">
          <label class="label">Имя</label>
          <input 
            v-model="candidateName" 
            type="text" 
            class="input" 
            placeholder="Ваше имя"
            required
          />
        </div>
        <div class="form-group">
          <label class="label">Email</label>
          <input 
            v-model="candidateEmail" 
            type="email" 
            class="input" 
            placeholder="your@email.com"
            required
          />
        </div>
        <button 
          class="btn btn-primary" 
          @click="handleBooking"
          :disabled="booking || !candidateName || !candidateEmail"
        >
          {{ booking ? 'Бронирование...' : 'Забронировать' }}
        </button>
      </div>

      <!-- Success State -->
      <div v-if="bookingSuccess" class="success-state">
        <div class="success-icon">✓</div>
        <h3>Встреча забронирована!</h3>
        <p>Письмо отправлено на {{ candidateEmail }}</p>
        <a 
          :href="calendarUrl" 
          target="_blank" 
          class="btn btn-secondary"
        >
          Добавить в Google Calendar
        </a>
        <router-link :to="`/cancel/${cancelToken}`" class="cancel-link">
          Отменить встречу
        </router-link>
      </div>

      <!-- Error -->
      <div v-if="error" class="error">
        {{ error }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { publicApi } from '@/api'

const route = useRoute()
const username = route.params.username
const slug = route.params.slug

const publicInfo = ref({ recruiter: null, event_type: null })
const loadingInfo = ref(true)
const availableDates = ref([])
const selectedDate = ref(null)
const slots = ref([])
const loadingSlots = ref(false)
const selectedSlot = ref(null)
const candidateName = ref('')
const candidateEmail = ref('')
const booking = ref(false)
const bookingSuccess = ref(false)
const cancelToken = ref('')
const error = ref('')

const recruiterInitials = computed(() => {
  const name = publicInfo.value.recruiter?.display_name || 'R'
  return name.charAt(0).toUpperCase()
})

const selectedDateFormatted = computed(() => {
  if (!selectedDate.value) return ''
  const date = new Date(selectedDate.value)
  return date.toLocaleDateString('ru-RU', { weekday: 'long', day: 'numeric', month: 'long' })
})

const calendarUrl = computed(() => {
  if (!bookingSuccess.value || !publicInfo.value.event_type) return ''
  const startsAt = new Date()
  const endsAt = new Date(startsAt.getTime() + (publicInfo.value.event_type.duration_min * 60 * 1000))
  
  return `https://calendar.google.com/calendar/render?action=TEMPLATE&text=${encodeURIComponent(publicInfo.value.event_type.title)}&dates=${formatDateForCalendar(startsAt)}/${formatDateForCalendar(endsAt)}&details=${encodeURIComponent('Встреча с ' + publicInfo.value.recruiter.display_name)}`
})

const formatDateForCalendar = (date) => {
  return date.toISOString().replace(/[-:]/g, '').split('.')[0] + 'Z'
}

const generateDates = () => {
  const dates = []
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  
  for (let i = 0; i < 14; i++) {
    const date = new Date(today)
    date.setDate(today.getDate() + i)
    
    dates.push({
      date: formatLocalDate(date),
      day: date.toLocaleDateString('ru-RU', { weekday: 'short' }),
      number: date.getDate(),
      hasSlots: true
    })
  }
  
  availableDates.value = dates
}

const formatLocalDate = (date) => {
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  return `${y}-${m}-${d}`
}

const selectDate = async (date) => {
  selectedDate.value = date
  selectedSlot.value = null
  loadingSlots.value = true
  error.value = ''
  
  try {
    const response = await publicApi.getPublicSlots(username, slug, {
      date: date,
      tz: Intl.DateTimeFormat().resolvedOptions().timeZone
    })
    
    slots.value = Array.isArray(response.data?.slots) ? response.data.slots : []
    
    const dateIndex = availableDates.value.findIndex(d => d.date === date)
    if (dateIndex !== -1) {
      availableDates.value[dateIndex].hasSlots = slots.value.length > 0
    }
  } catch (err) {
    slots.value = []
    if (err.response?.status === 429) {
      error.value = 'Слишком много запросов, подождите немного'
    } else {
      error.value = 'Не удалось загрузить слоты'
    }
  } finally {
    loadingSlots.value = false
  }
}

const handleBooking = async () => {
  booking.value = true
  error.value = ''
  
  try {
    const response = await publicApi.createBooking(username, slug, {
      candidate_name: candidateName.value,
      candidate_email: candidateEmail.value,
      slot_time: selectedSlot.value,
      date: selectedDate.value,
      timezone: Intl.DateTimeFormat().resolvedOptions().timeZone
    })
    
    cancelToken.value = response.data.cancel_token
    bookingSuccess.value = true
  } catch (err) {
    error.value = err.response?.data?.error || 'Не удалось забронировать встречу'
  } finally {
    booking.value = false
  }
}

onMounted(async () => {
  try {
    const response = await publicApi.getPublicInfo(username, slug)
    publicInfo.value = response.data
    generateDates()
    
    // Select first date by default
    if (availableDates.value.length > 0) {
      selectDate(availableDates.value[0].date)
    }
  } catch (err) {
    error.value = 'Не удалось загрузить информацию'
  } finally {
    loadingInfo.value = false
  }
})
</script>

<style scoped>
.booking-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-bg);
  padding: var(--sp-4);
}

.booking-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--r-lg);
  padding: var(--sp-8);
  width: 100%;
  max-width: 600px;
  box-shadow: var(--shadow-md);
}

.recruiter-info {
  display: flex;
  gap: var(--sp-4);
  margin-bottom: var(--sp-8);
  padding-bottom: var(--sp-6);
  border-bottom: 1px solid var(--color-border);
}

.avatar {
  width: 60px;
  height: 60px;
  border-radius: var(--r-full);
  background: var(--color-accent);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--text-xl);
  font-weight: 600;
}

.recruiter-details {
  flex: 1;
}

.recruiter-details h2 {
  font-size: var(--text-xl);
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: var(--sp-1);
}

.bio {
  color: var(--color-text-secondary);
  font-size: var(--text-sm);
  margin-bottom: var(--sp-2);
}

.meta {
  display: flex;
  gap: var(--sp-4);
}

.meta-item {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
}

.date-selection h3,
.time-slots h3,
.booking-form h3 {
  font-size: var(--text-lg);
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: var(--sp-4);
}

.dates-scroll {
  display: flex;
  gap: var(--sp-2);
  overflow-x: auto;
  padding-bottom: var(--sp-2);
  margin-bottom: var(--sp-6);
}

.date-btn {
  min-width: 60px;
  padding: var(--sp-3);
  border: 1px solid var(--color-border);
  border-radius: var(--r-md);
  background: var(--color-surface);
  cursor: pointer;
  transition: all var(--t-fast);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--sp-1);
}

.date-btn:hover:not(.disabled) {
  background: var(--color-surface-2);
}

.date-btn.active {
  background: var(--color-accent);
  color: white;
  border-color: var(--color-accent);
}

.date-btn.disabled {
  opacity: 0.5;
  cursor: not-allowed;
  pointer-events: none;
}

.date-day {
  font-size: var(--text-xs);
  text-transform: capitalize;
}

.date-number {
  font-size: var(--text-lg);
  font-weight: 600;
}

.slots-loading {
  display: flex;
  gap: var(--sp-2);
}

.slots-grid {
  display: flex;
  flex-wrap: wrap;
  gap: var(--sp-2);
  margin-bottom: var(--sp-6);
}

.slot-btn {
  padding: var(--sp-3) var(--sp-4);
  border: 1px solid var(--color-border);
  border-radius: var(--r-md);
  background: var(--color-surface);
  cursor: pointer;
  transition: all var(--t-fast);
  font-size: var(--text-sm);
}

.slot-btn:hover {
  background: var(--color-surface-2);
}

.slot-btn.active {
  background: var(--color-accent);
  color: white;
  border-color: var(--color-accent);
}

.no-slots {
  color: var(--color-text-secondary);
  font-size: var(--text-sm);
  padding: var(--sp-4);
  text-align: center;
}

.booking-form {
  margin-top: var(--sp-6);
  padding-top: var(--sp-6);
  border-top: 1px solid var(--color-border);
}

.form-group {
  margin-bottom: var(--sp-4);
}

.btn {
  width: 100%;
  justify-content: center;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.success-state {
  text-align: center;
  padding: var(--sp-8);
}

.success-icon {
  width: 60px;
  height: 60px;
  border-radius: var(--r-full);
  background: var(--color-success);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--text-2xl);
  font-weight: 600;
  margin: 0 auto var(--sp-4);
  animation: scaleIn 0.4s ease;
}

@keyframes scaleIn {
  from {
    transform: scale(0);
  }
  to {
    transform: scale(1);
  }
}

.success-state h3 {
  font-size: var(--text-xl);
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: var(--sp-2);
}

.success-state p {
  color: var(--color-text-secondary);
  margin-bottom: var(--sp-6);
}

.cancel-link {
  display: block;
  margin-top: var(--sp-4);
  color: var(--color-text-secondary);
  text-decoration: none;
  font-size: var(--text-sm);
}

.cancel-link:hover {
  color: var(--color-text-primary);
}

.error {
  padding: var(--sp-3);
  background: #FEE2E2;
  color: var(--color-error);
  border-radius: var(--r-md);
  font-size: var(--text-sm);
  margin-top: var(--sp-4);
}
</style>
