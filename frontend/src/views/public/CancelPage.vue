<template>
  <div class="cancel-page">
    <div class="cancel-card">
      <!-- Cancelled Successfully -->
      <div v-if="cancelSuccess" class="success-state">
        <div class="icon">✓</div>
        <h2>Встреча отменена</h2>
        <p>{{ cancelInfo?.recruiter_display_name || 'Организатор' }} получил уведомление</p>
      </div>

      <!-- Loading State -->
      <div v-else-if="loading" class="loading-state">
        <div class="skeleton" style="height: 60px; width: 60px; border-radius: 50%; margin: 0 auto var(--sp-4);"></div>
        <div class="skeleton" style="height: 24px; width: 200px; margin: 0 auto var(--sp-2);"></div>
        <div class="skeleton" style="height: 16px; width: 300px; margin: 0 auto var(--sp-4);"></div>
        <div class="skeleton" style="height: 40px; width: 150px; margin: 0 auto;"></div>
      </div>

      <!-- Active Booking -->
      <div v-else-if="cancelInfo && cancelInfo.status === 'confirmed'" class="active-booking">
        <div class="icon">📅</div>
        <h2>Встреча</h2>
        <div class="booking-details">
          <div class="detail">
            <span class="label">Кто:</span>
            <span class="value">{{ cancelInfo.candidate_name }}</span>
          </div>
          <div class="detail">
            <span class="label">Когда:</span>
            <span class="value">{{ formattedDate }}</span>
          </div>
          <div class="detail">
            <span class="label">Длительность:</span>
            <span class="value">{{ cancelInfo.duration_min }} мин</span>
          </div>
          <div class="detail">
            <span class="label">Тип:</span>
            <span class="value">{{ cancelInfo.event_type_title }}</span>
          </div>
          <div class="detail">
            <span class="label">С:</span>
            <span class="value">{{ cancelInfo.recruiter_display_name }}</span>
          </div>
        </div>
        <button 
          class="btn btn-danger" 
          @click="handleCancel"
          :disabled="cancelling"
        >
          {{ cancelling ? 'Отмена...' : 'Отменить встречу' }}
        </button>
      </div>

      <!-- Already Cancelled -->
      <div v-else-if="cancelInfo && cancelInfo.status === 'cancelled'" class="cancelled-state">
        <div class="icon">✓</div>
        <h2>Встреча уже отменена</h2>
        <p>Эта встреча была отменена ранее</p>
      </div>

      <!-- Error -->
      <div v-else-if="error" class="error-state">
        <div class="icon">⚠️</div>
        <h2>Ошибка</h2>
        <p>{{ error }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { publicApi } from '@/api'

const route = useRoute()
const token = route.params.token

const cancelInfo = ref(null)
const loading = ref(true)
const cancelling = ref(false)
const cancelSuccess = ref(false)
const error = ref('')

const formattedDate = computed(() => {
  if (!cancelInfo.value?.starts_at) return ''
  const date = new Date(cancelInfo.value.starts_at)
  return date.toLocaleDateString('ru-RU', { 
    day: 'numeric', 
    month: 'long', 
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
})

const handleCancel = async () => {
  cancelling.value = true
  error.value = ''
  
  try {
    await publicApi.cancelBooking(token)
    if (cancelInfo.value) {
      cancelInfo.value = { ...cancelInfo.value, status: 'cancelled' }
    }
    cancelSuccess.value = true
  } catch (err) {
    error.value = err.response?.data?.error || 'Не удалось отменить встречу'
  } finally {
    cancelling.value = false
  }
}

onMounted(async () => {
  try {
    const response = await publicApi.getCancelInfo(token)
    cancelInfo.value = response.data
  } catch (err) {
    error.value = 'Не удалось найти встречу'
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.cancel-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-bg);
  padding: var(--sp-4);
}

.cancel-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--r-lg);
  padding: var(--sp-12);
  width: 100%;
  max-width: 500px;
  box-shadow: var(--shadow-md);
  text-align: center;
}

.loading-state {
  padding: var(--sp-8);
}

.icon {
  width: 60px;
  height: 60px;
  border-radius: var(--r-full);
  background: var(--color-surface-2);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--text-2xl);
  margin: 0 auto var(--sp-4);
}

.active-booking h2,
.cancelled-state h2,
.success-state h2,
.error-state h2 {
  font-size: var(--text-xl);
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: var(--sp-4);
}

.booking-details {
  text-align: left;
  margin-bottom: var(--sp-6);
  padding: var(--sp-6);
  background: var(--color-surface-2);
  border-radius: var(--r-md);
}

.detail {
  display: flex;
  justify-content: space-between;
  padding: var(--sp-2) 0;
  border-bottom: 1px solid var(--color-border);
}

.detail:last-child {
  border-bottom: none;
}

.detail .label {
  color: var(--color-text-secondary);
  font-size: var(--text-sm);
}

.detail .value {
  color: var(--color-text-primary);
  font-weight: 500;
}

.btn {
  width: 100%;
  justify-content: center;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.cancelled-state p,
.success-state p,
.error-state p {
  color: var(--color-text-secondary);
  margin-bottom: var(--sp-4);
}

.success-state .icon {
  background: var(--color-success);
  color: white;
}

.error-state .icon {
  background: #FEE2E2;
  color: var(--color-error);
}
</style>
