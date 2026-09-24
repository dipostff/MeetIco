<template>
  <div class="bookings">
    <OnboardingBanner />
    
    <div class="header">
      <h1>Встречи</h1>
      <button class="btn btn-secondary" @click="copyLink">
        📋 Скопировать ссылку
      </button>
    </div>

    <!-- Empty State -->
    <div v-if="!bookingsStore.loading && bookingsStore.bookings.length === 0" class="empty-state">
      <div class="empty-icon">📅</div>
      <h2>Пока тихо</h2>
      <p>Поделитесь ссылкой — и здесь появятся встречи</p>
      <button class="btn btn-primary" @click="copyLink">
        Скопировать ссылку
      </button>
    </div>

    <!-- Loading State -->
    <div v-if="bookingsStore.loading" class="loading-state">
      <div class="skeleton" style="height: 80px; margin-bottom: var(--sp-4);"></div>
      <div class="skeleton" style="height: 80px; margin-bottom: var(--sp-4);"></div>
      <div class="skeleton" style="height: 80px;"></div>
    </div>

    <!-- Bookings List -->
    <div v-else class="bookings-list">
      <div v-for="group in groupedBookings" :key="group.title" class="booking-group">
        <h3 class="group-title">{{ group.title }}</h3>
        <div v-for="booking in group.bookings" :key="booking.id" class="booking-card">
          <div class="booking-time">
            <div class="time">{{ formatTime(booking.starts_at) }}</div>
            <div class="date">{{ formatDate(booking.starts_at) }}</div>
          </div>
          <div class="booking-info">
            <div class="candidate-name">{{ booking.candidate_name }}</div>
            <div class="candidate-email">{{ booking.candidate_email }}</div>
          </div>
          <button
            class="btn btn-ghost btn-sm"
            :disabled="cancellingId === booking.id"
            @click="cancelBooking(booking)"
          >
            {{ cancellingId === booking.id ? 'Отмена...' : 'Отменить' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useBookingsStore } from '@/stores/bookings'
import { useAuthStore } from '@/stores/auth'
import { useEventTypesStore } from '@/stores/eventTypes'
import OnboardingBanner from '@/components/OnboardingBanner.vue'

const router = useRouter()
const bookingsStore = useBookingsStore()
const authStore = useAuthStore()
const eventTypesStore = useEventTypesStore()

const groupedBookings = computed(() => {
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  
  const tomorrow = new Date(today)
  tomorrow.setDate(tomorrow.getDate() + 1)
  
  const groups = [
    { title: 'Сегодня', bookings: [] },
    { title: 'Завтра', bookings: [] },
    { title: 'Позже', bookings: [] }
  ]
  
  bookingsStore.bookings.forEach(booking => {
    const bookingDate = new Date(booking.starts_at)
    bookingDate.setHours(0, 0, 0, 0)
    
    if (bookingDate.getTime() === today.getTime()) {
      groups[0].bookings.push(booking)
    } else if (bookingDate.getTime() === tomorrow.getTime()) {
      groups[1].bookings.push(booking)
    } else {
      groups[2].bookings.push(booking)
    }
  })
  
  return groups.filter(g => g.bookings.length > 0)
})

const formatTime = (dateStr) => {
  const date = new Date(dateStr)
  return date.toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' })
}

const formatDate = (dateStr) => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('ru-RU', { day: 'numeric', month: 'short' })
}

const copyLink = async () => {
  if (!authStore.user) {
    await authStore.fetchUser()
  }

  const username = authStore.user?.username
  if (!username) {
    alert('Сначала укажите username в Настройках')
    router.push('/dashboard/settings')
    return
  }

  if (!eventTypesStore.eventTypes.length) {
    await eventTypesStore.fetchEventTypes()
  }

  const eventType = eventTypesStore.eventTypes.find(et => et.is_active) || eventTypesStore.eventTypes[0]
  if (!eventType) {
    alert('Сначала создайте тип встречи в разделе «Мои ссылки»')
    router.push('/dashboard/event-types')
    return
  }

  const link = `${window.location.origin}/book/${username}/${eventType.slug}`
  try {
    await navigator.clipboard.writeText(link)
    alert('Ссылка скопирована!')
  } catch {
    prompt('Скопируйте ссылку вручную:', link)
  }
}

const cancellingId = ref(null)

const cancelBooking = async (booking) => {
  if (cancellingId.value) return
  cancellingId.value = booking.id
  const ok = await bookingsStore.cancelBooking(booking)
  cancellingId.value = null
  if (!ok) {
    alert(bookingsStore.error || 'Не удалось отменить встречу')
  }
}

onMounted(async () => {
  await Promise.all([
    bookingsStore.fetchBookings(),
    eventTypesStore.fetchEventTypes(),
    authStore.user ? Promise.resolve() : authStore.fetchUser()
  ])
})
</script>

<style scoped>
.bookings {
  max-width: 800px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--sp-8);
}

.header h1 {
  font-size: var(--text-2xl);
  font-weight: 600;
  color: var(--color-text-primary);
}

.empty-state {
  text-align: center;
  padding: var(--sp-12);
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--r-lg);
}

.empty-icon {
  font-size: 48px;
  margin-bottom: var(--sp-4);
}

.empty-state h2 {
  font-size: var(--text-xl);
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: var(--sp-2);
}

.empty-state p {
  color: var(--color-text-secondary);
  margin-bottom: var(--sp-6);
}

.loading-state {
  padding: var(--sp-4);
}

.booking-group {
  margin-bottom: var(--sp-8);
}

.group-title {
  font-size: var(--text-lg);
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: var(--sp-4);
}

.booking-card {
  display: flex;
  align-items: center;
  gap: var(--sp-4);
  padding: var(--sp-4);
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--r-md);
  margin-bottom: var(--sp-2);
}

.booking-time {
  text-align: center;
  min-width: 80px;
}

.time {
  font-size: var(--text-lg);
  font-weight: 600;
  color: var(--color-text-primary);
}

.date {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
}

.booking-info {
  flex: 1;
}

.candidate-name {
  font-weight: 500;
  color: var(--color-text-primary);
  margin-bottom: var(--sp-1);
}

.candidate-email {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
}

.btn-sm {
  padding: var(--sp-2) var(--sp-3);
  font-size: var(--text-sm);
}
</style>
