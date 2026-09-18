<template>
  <div class="event-types">
    <OnboardingBanner />
    
    <div class="header">
      <h1>Мои ссылки</h1>
      <button class="btn btn-primary" @click="showCreateModal = true">
        + Создать тип встречи
      </button>
    </div>

    <!-- Empty State -->
    <div v-if="!eventTypesStore.loading && eventTypesStore.eventTypes.length === 0" class="empty-state">
      <div class="empty-icon">🔗</div>
      <h2>Нет типов встреч</h2>
      <p>Создайте первый тип встречи, чтобы начать делиться ссылками</p>
      <button class="btn btn-primary" @click="showCreateModal = true">
        Создать тип встречи
      </button>
    </div>

    <!-- Loading State -->
    <div v-if="eventTypesStore.loading" class="loading-state">
      <div class="skeleton" style="height: 100px; margin-bottom: var(--sp-4);"></div>
      <div class="skeleton" style="height: 100px; margin-bottom: var(--sp-4);"></div>
      <div class="skeleton" style="height: 100px;"></div>
    </div>

    <!-- Event Types List -->
    <div v-else class="event-types-list">
      <div 
        v-for="eventType in eventTypesStore.eventTypes" 
        :key="eventType.id" 
        class="event-type-card"
      >
        <div class="event-type-info">
          <h3>{{ eventType.title }}</h3>
          <div class="meta">
            <span class="meta-item">⏱ {{ eventType.duration_min }} мин</span>
            <span class="meta-item">📍 {{ eventType.location_type }}</span>
            <span class="meta-item">📋 {{ eventType.category }}</span>
          </div>
          <div class="link">
            {{ publicLink(eventType) }}
            <button class="copy-btn" @click="copyLink(eventType)">
              📋
            </button>
          </div>
        </div>
        <div class="event-type-actions">
          <button class="btn btn-ghost btn-sm" @click="editEventType(eventType)">
            Редактировать
          </button>
          <button 
            class="btn btn-ghost btn-sm" 
            :class="{ 'text-danger': !eventType.is_active }"
            @click="toggleEventType(eventType)"
          >
            {{ eventType.is_active ? 'Деактивировать' : 'Активировать' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Create/Edit Modal -->
    <div v-if="showCreateModal" class="modal-overlay" @click="showCreateModal = false">
      <div class="modal" @click.stop>
        <h2>{{ editingEventType ? 'Редактировать тип встречи' : 'Создать тип встречи' }}</h2>
        <form @submit.prevent="handleSave" class="modal-form">
          <div class="form-group">
            <label class="label">Название</label>
            <input 
              v-model="formData.title" 
              type="text" 
              class="input" 
              placeholder="30 минутная встреча"
              required
            />
          </div>
          <div class="form-group">
            <label class="label">Slug (URL)</label>
            <input 
              v-model="formData.slug" 
              type="text" 
              class="input" 
              placeholder="30min"
              required
              pattern="[a-z0-9_-]+"
            />
          </div>
          <div class="form-group">
            <label class="label">Длительность (мин)</label>
            <input 
              v-model.number="formData.duration_min" 
              type="number" 
              class="input" 
              placeholder="30"
              required
              min="15"
              step="5"
            />
          </div>
          <div class="form-group">
            <label class="label">Тип локации</label>
            <select v-model="formData.location_type" class="input" required>
              <option value="zoom">Zoom</option>
              <option value="google-meet">Google Meet</option>
              <option value="phone">Телефон</option>
              <option value="in-person">Встреча в офисе</option>
            </select>
          </div>
          <div class="form-group">
            <label class="label">Значение локации</label>
            <input 
              v-model="formData.location_value" 
              type="text" 
              class="input" 
              placeholder="Zoom link или телефон"
            />
          </div>
          <div class="form-group">
            <label class="label">Категория</label>
            <select v-model="formData.category" class="input" required>
              <option value="interview">Собеседование</option>
              <option value="consultation">Консультация</option>
              <option value="mentoring">Менторство</option>
              <option value="other">Другое</option>
            </select>
          </div>
          <div class="modal-actions">
            <button type="button" class="btn btn-secondary" @click="showCreateModal = false">
              Отмена
            </button>
            <button type="submit" class="btn btn-primary" :disabled="eventTypesStore.loading">
              {{ eventTypesStore.loading ? 'Сохранение...' : 'Сохранить' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useEventTypesStore } from '@/stores/eventTypes'
import { useAuthStore } from '@/stores/auth'
import OnboardingBanner from '@/components/OnboardingBanner.vue'

const eventTypesStore = useEventTypesStore()
const authStore = useAuthStore()

const showCreateModal = ref(false)
const editingEventType = ref(null)
const formData = ref({
  title: '',
  slug: '',
  duration_min: 30,
  location_type: 'zoom',
  location_value: '',
  category: 'interview'
})

const publicLink = (eventType) => {
  const username = authStore.user?.username
  return username ? `${window.location.origin}/book/${username}/${eventType.slug}` : ''
}

const copyLink = (eventType) => {
  const link = publicLink(eventType)
  navigator.clipboard.writeText(link)
  alert('Ссылка скопирована!')
}

const editEventType = (eventType) => {
  editingEventType.value = eventType
  formData.value = {
    title: eventType.title,
    slug: eventType.slug,
    duration_min: eventType.duration_min,
    location_type: eventType.location_type,
    location_value: eventType.location_value || '',
    category: eventType.category
  }
  showCreateModal.value = true
}

const toggleEventType = async (eventType) => {
  const success = await eventTypesStore.updateEventType(eventType.id, {
    is_active: !eventType.is_active
  })
  if (!success) {
    alert('Не удалось обновить тип встречи')
  }
}

const handleSave = async () => {
  let success
  if (editingEventType.value) {
    success = await eventTypesStore.updateEventType(editingEventType.value.id, formData.value)
  } else {
    success = await eventTypesStore.createEventType(formData.value)
  }
  
  if (success) {
    showCreateModal.value = false
    resetForm()
  } else {
    alert('Не удалось сохранить тип встречи')
  }
}

const resetForm = () => {
  editingEventType.value = null
  formData.value = {
    title: '',
    slug: '',
    duration_min: 30,
    location_type: 'zoom',
    location_value: '',
    category: 'interview'
  }
}

onMounted(() => {
  eventTypesStore.fetchEventTypes()
})
</script>

<style scoped>
.event-types {
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

.event-types-list {
  display: flex;
  flex-direction: column;
  gap: var(--sp-4);
}

.event-type-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--sp-6);
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--r-lg);
}

.event-type-info {
  flex: 1;
}

.event-type-info h3 {
  font-size: var(--text-lg);
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: var(--sp-2);
}

.meta {
  display: flex;
  gap: var(--sp-4);
  margin-bottom: var(--sp-2);
}

.meta-item {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
}

.link {
  display: flex;
  align-items: center;
  gap: var(--sp-2);
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  background: var(--color-surface-2);
  padding: var(--sp-2) var(--sp-3);
  border-radius: var(--r-md);
  width: fit-content;
}

.copy-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 0;
  font-size: var(--text-sm);
}

.event-type-actions {
  display: flex;
  gap: var(--sp-2);
}

.btn-sm {
  padding: var(--sp-2) var(--sp-3);
  font-size: var(--text-sm);
}

.text-danger {
  color: var(--color-error) !important;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background: var(--color-surface);
  border-radius: var(--r-lg);
  padding: var(--sp-8);
  width: 100%;
  max-width: 500px;
  max-height: 90vh;
  overflow-y: auto;
}

.modal h2 {
  font-size: var(--text-xl);
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: var(--sp-6);
}

.modal-form {
  display: flex;
  flex-direction: column;
  gap: var(--sp-4);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--sp-2);
}

.modal-actions {
  display: flex;
  gap: var(--sp-4);
  margin-top: var(--sp-4);
}

.modal-actions .btn {
  flex: 1;
}
</style>
