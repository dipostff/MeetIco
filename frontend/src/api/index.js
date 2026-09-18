import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json'
  }
})

// Add auth token to requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Auth
export const auth = {
  register: (data) => api.post('/auth/register', data),
  login: (data) => api.post('/auth/login', data),
  changePassword: (data) => api.post('/auth/change-password', data)
}

// Profile
export const profile = {
  getMe: () => api.get('/me'),
  updateMe: (data) => api.patch('/me', data)
}

// Schedule
export const schedule = {
  getSchedule: () => api.get('/schedule'),
  updateSchedule: (data) => api.patch('/schedule', data)
}

// Event Types
export const eventTypes = {
  getEventTypes: () => api.get('/event-types'),
  createEventType: (data) => api.post('/event-types', data),
  updateEventType: (id, data) => api.patch(`/event-types/${id}`, data),
  deleteEventType: (id) => api.delete(`/event-types/${id}`)
}

// Bookings
export const bookings = {
  getBookings: () => api.get('/bookings')
}

// Public
export const publicApi = {
  getPublicInfo: (username, slug) => api.get(`/public/${username}/${slug}`),
  getPublicSlots: (username, slug, params) => api.get(`/public/${username}/${slug}/slots`, { params }),
  createBooking: (username, slug, data) => api.post(`/public/${username}/${slug}/book`, data),
  getCancelInfo: (token) => api.get(`/public/cancel/${token}`),
  cancelBooking: (token) => api.post(`/public/cancel/${token}`)
}

export default api
