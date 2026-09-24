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

    async cancelBooking(booking) {
      if (!booking?.id) {
        this.error = 'Нет id встречи'
        return false
      }
      try {
        await bookings.cancelBooking(booking.id)
        this.bookings = this.bookings.filter(b => b.id !== booking.id)
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Не удалось отменить встречу'
        return false
      }
    },

    clearBookings() {
      this.bookings = []
    }
  }
})
