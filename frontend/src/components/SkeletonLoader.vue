<template>
  <div class="skeleton-loader">
    <div class="skeleton-line" :style="{ width: width }"></div>
  </div>
</template>

<script setup>
defineProps({
  width: {
    type: String,
    default: '100%'
  }
})
</script>

<style scoped>
.skeleton-loader {
  margin-bottom: var(--sp-2);
}

.skeleton-line {
  height: 16px;
  border-radius: var(--r-md);
}
</style>
