// Pinia store mirroring the auth state the Go server shares with every page
// (framework/response.go injects props.auth). Components can read the user
// from anywhere without prop drilling.

import { defineStore } from 'pinia'
import type { AuthUser, Page, SharedProps } from '@/flux/types'

interface AuthState {
  user: AuthUser | null
  appName: string
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    user: null,
    appName: 'Flux',
  }),

  getters: {
    isAuthenticated: (state) => state.user !== null,
    isVerified: (state) => Boolean(state.user?.emailVerifiedAt),
    initials(state): string {
      if (!state.user) return '?'
      return state.user.name
        .split(' ')
        .map((part) => part[0])
        .slice(0, 2)
        .join('')
        .toUpperCase()
    },
  },

  actions: {
    /** Called by app.ts on every page swap. */
    sync(page: Page) {
      const shared = page.props as unknown as SharedProps
      this.user = shared.auth?.user ?? null
      if (shared.appName) this.appName = shared.appName
    },
  },
})
