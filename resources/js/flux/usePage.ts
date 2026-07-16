// Reactive access to the current page object — the Flux equivalent of
// Inertia's usePage(). Shared props (auth, appName) are always present.

import { computed, reactive } from 'vue'
import type { Page, SharedProps } from './types'

interface PageState {
  page: Page | null
}

const state = reactive<PageState>({ page: null })

/** Called by app.ts whenever the server sends a new page. */
export function setCurrentPage(page: Page) {
  state.page = page
}

/** The current page (component, props, url, version). */
export function usePage<P = Record<string, unknown>>() {
  return computed(() => state.page as Page<P> | null)
}

/** Shared props shortcut: the authenticated user (null for guests). */
export function useAuth() {
  return computed(() => (state.page?.props as SharedProps | undefined)?.auth.user ?? null)
}

/** Shared props shortcut: the application name. */
export function useAppName() {
  return computed(() => (state.page?.props as SharedProps | undefined)?.appName ?? 'Flux')
}
