import { defineStore } from 'pinia'
import { bookings } from '@/api'

export const useBookingsStore = defineStore('bookings', {
  state: () => ({
    bookings: [],
    loading: false,
    error: null
  }),

  actions: {
    async fetchBookings() {
      this.loading = true
      this.error = null
      try {
        const response = await bookings.getBookings()
        this.bookings = Array.isArray(response.data) ? response.data : []
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to fetch bookings'
        this.bookings = []
      } finally {
        this.loading = false
      }
    },

    clearBookings() {
      this.bookings = []
    }
  }
})
