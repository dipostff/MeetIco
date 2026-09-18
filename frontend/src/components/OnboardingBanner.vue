<template>
  <div class="onboarding-banner">
    <div class="steps">
      <div class="step" :class="{ active: currentStep >= 1, completed: currentStep > 1 }">
        <span class="step-number">1</span>
        <span class="step-text">Создайте аккаунт</span>
      </div>
      <div class="step" :class="{ active: currentStep >= 2, completed: currentStep > 2 }">
        <span class="step-number">2</span>
        <span class="step-text">Настройте профиль</span>
      </div>
      <div class="step" :class="{ active: currentStep >= 3, completed: currentStep > 3 }">
        <span class="step-number">3</span>
        <span class="step-text">Создайте тип встречи</span>
      </div>
      <div class="step" :class="{ active: currentStep >= 4, completed: currentStep > 4 }">
        <span class="step-number">4</span>
        <span class="step-text">Поделитесь ссылкой</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useEventTypesStore } from '@/stores/eventTypes'

const authStore = useAuthStore()
const eventTypesStore = useEventTypesStore()

const currentStep = computed(() => {
  if (!authStore.user?.username) return 1
  if (eventTypesStore.eventTypes.length === 0) return 2
  return 3
})
</script>

<style scoped>
.onboarding-banner {
  background: var(--color-accent-light);
  border: 1px solid var(--color-accent);
  border-radius: var(--r-lg);
  padding: var(--sp-6);
  margin-bottom: var(--sp-8);
}

.steps {
  display: flex;
  justify-content: space-between;
  gap: var(--sp-4);
}

.step {
  display: flex;
  align-items: center;
  gap: var(--sp-2);
  opacity: 0.5;
  transition: opacity var(--t-fast);
}

.step.active {
  opacity: 1;
}

.step.completed {
  opacity: 1;
}

.step-number {
  width: 24px;
  height: 24px;
  border-radius: var(--r-full);
  background: var(--color-border);
  color: var(--color-text-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--text-sm);
  font-weight: 600;
}

.step.active .step-number {
  background: var(--color-accent);
  color: white;
}

.step.completed .step-number {
  background: var(--color-success);
  color: white;
}

.step-text {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
}

.step.active .step-text {
  color: var(--color-text-primary);
  font-weight: 500;
}
</style>
