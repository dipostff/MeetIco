import { defineStore } from 'pinia'
import { eventTypes } from '@/api'

export const useEventTypesStore = defineStore('eventTypes', {
  state: () => ({
    eventTypes: [],
    loading: false,
    error: null
  }),

  actions: {
    async fetchEventTypes() {
      this.loading = true
      this.error = null
      try {
        const response = await eventTypes.getEventTypes()
        this.eventTypes = Array.isArray(response.data) ? response.data : []
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to fetch event types'
        this.eventTypes = []
      } finally {
        this.loading = false
      }
    },

    async createEventType(data) {
      this.loading = true
      this.error = null
      try {
        const response = await eventTypes.createEventType(data)
        this.eventTypes.push(response.data)
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to create event type'
        return false
      } finally {
        this.loading = false
      }
    },

    async updateEventType(id, data) {
      this.loading = true
      this.error = null
      try {
        const response = await eventTypes.updateEventType(id, data)
        const index = this.eventTypes.findIndex(et => et.id === id)
        if (index !== -1) {
          this.eventTypes[index] = response.data
        }
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to update event type'
        return false
      } finally {
        this.loading = false
      }
    },

    async deleteEventType(id) {
      this.loading = true
      this.error = null
      try {
        await eventTypes.deleteEventType(id)
        this.eventTypes = this.eventTypes.filter(et => et.id !== id)
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to delete event type'
        return false
      } finally {
        this.loading = false
      }
    }
  }
})
