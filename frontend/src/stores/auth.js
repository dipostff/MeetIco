import { defineStore } from 'pinia'
import { auth, profile } from '@/api'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    token: localStorage.getItem('token') || null,
    loading: false,
    error: null
  }),

  getters: {
    isAuthenticated: (state) => !!state.token,
    currentUser: (state) => state.user
  },

  actions: {
    async register(data) {
      this.loading = true
      this.error = null
      try {
        const response = await auth.register(data)
        this.token = response.data.token
        localStorage.setItem('token', response.data.token)
        await this.fetchUser()
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Registration failed'
        return false
      } finally {
        this.loading = false
      }
    },

    async login(data) {
      this.loading = true
      this.error = null
      try {
        const response = await auth.login(data)
        this.token = response.data.token
        localStorage.setItem('token', response.data.token)
        await this.fetchUser()
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Login failed'
        return false
      } finally {
        this.loading = false
      }
    },

    async fetchUser() {
      if (!this.token) return

      this.loading = true
      try {
        const response = await profile.getMe()
        this.user = response.data
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to fetch user'
        this.logout()
      } finally {
        this.loading = false
      }
    },

    async updateProfile(data) {
      this.loading = true
      this.error = null
      try {
        const response = await profile.updateMe(data)
        this.user = response.data
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to update profile'
        return false
      } finally {
        this.loading = false
      }
    },

    async changePassword(data) {
      this.loading = true
      this.error = null
      try {
        await auth.changePassword(data)
        return true
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to change password'
        return false
      } finally {
        this.loading = false
      }
    },

    logout() {
      this.user = null
      this.token = null
      localStorage.removeItem('token')
    }
  }
})
