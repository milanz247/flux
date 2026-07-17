// Core types of the Flux Inertia-style protocol.

export interface AuthUser {
  id: number
  name: string
  email: string
  emailVerifiedAt?: string
}

/** One-shot message set server-side via req.RedirectWith(url, type, message). */
export interface Flash {
  type: string
  message: string
}

/** Props the Go framework injects into every page (framework/response.go). */
export interface SharedProps {
  auth: { user: AuthUser | null }
  appName: string
  flash: Flash | null
}

/** The page object exchanged with the Go server. */
export interface Page<P = Record<string, unknown>> {
  component: string
  props: P & SharedProps
  url: string
  version: string
}

export type VisitMethod = 'get' | 'post' | 'put' | 'patch' | 'delete'

export interface VisitOptions {
  method?: VisitMethod
  data?: Record<string, unknown>
  /** Keep the scroll position (e.g. pagination inside a long page). */
  preserveScroll?: boolean
  /** Replace the history entry instead of pushing a new one. */
  replace?: boolean
  onSuccess?: (page: Page) => void
  onError?: (errors: Record<string, string[]>) => void
  onFinish?: () => void
}

/** Laravel-style validation error bag: field → messages. */
export type ErrorBag = Record<string, string[]>
