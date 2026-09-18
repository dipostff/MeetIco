import { defineStore } from 'pinia'
import { schedule } from '@/api'

export const useScheduleStore = defineStore('schedule', {
  state: () => ({
    schedule: null,
    loading: false,
    error: null
  }),

  actions: {
    async fetchSchedule() {
      this.loading = true
      this.error = null
      try {
        const response = await schedule.getSchedule()
        this.schedule = response.data
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to fetch schedule'
      } finally {
        this.loading = false
      }
    },

    async updateSchedule(data) {
      this.loading = true
      this.error = null
      try {
        const response = await schedule.updateSchedule(data)
        this.schedule = response.data
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to update schedule'
        return false
      } finally {
        this.loading = false
      }
    }
  }
})
