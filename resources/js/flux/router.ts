// The Flux client router — the browser half of the Inertia protocol.
//
// Every navigation is an XHR carrying the X-Inertia header; the Go server
// answers with a page object ({component, props, url, version}) and the
// client swaps the Vue page component in place, updating history. Forms
// submit the same way, so validation errors arrive as JSON (422) and
// successful writes arrive as a redirect the browser follows transparently.

import axios, { AxiosError } from 'axios'
import type { ErrorBag, Page, VisitOptions } from './types'

type SwapFn = (page: Page, options: { preserveScroll: boolean }) => void

let currentPage: Page | null = null
let swap: SwapFn | null = null

function setPage(page: Page, options: { preserveScroll: boolean; replace: boolean }) {
  currentPage = page
  const state = { page }
  if (options.replace || page.url === window.location.pathname + window.location.search) {
    window.history.replaceState(state, '', page.url)
  } else {
    window.history.pushState(state, '', page.url)
  }
  swap?.(page, { preserveScroll: options.preserveScroll })
  if (!options.preserveScroll) window.scrollTo(0, 0)
}

/** Wire up the initial page and the component-swap callback (called once from app.ts). */
export function initRouter(initial: Page, swapFn: SwapFn) {
  currentPage = initial
  swap = swapFn
  window.history.replaceState({ page: initial }, '', initial.url)

  window.addEventListener('popstate', (event) => {
    const page: Page | undefined = event.state?.page
    if (page) {
      currentPage = page
      swap?.(page, { preserveScroll: false })
    } else {
      // No state (e.g. hash change) — do a fresh visit of the current URL.
      void router.visit(window.location.pathname + window.location.search, { replace: true })
    }
  })
}

async function visit(url: string, options: VisitOptions = {}): Promise<void> {
  const { method = 'get', data, preserveScroll = false, replace = false } = options

  try {
    const response = await axios.request<Page>({
      url,
      method,
      data: method === 'get' ? undefined : data,
      params: method === 'get' ? data : undefined,
      headers: {
        'X-Inertia': 'true',
        'X-Inertia-Version': currentPage?.version ?? '',
        Accept: 'application/json',
      },
      // Axios follows the 303 redirect after form posts automatically; the
      // final response is the Inertia page object of the target URL.
      validateStatus: (status) => status >= 200 && status < 300,
    })

    if (response.headers['x-inertia'] !== 'true' || !response.data?.component) {
      // Not an Inertia response (e.g. session expired to a non-Inertia
      // endpoint) — fall back to a full browser navigation.
      window.location.href = url
      return
    }

    const page = response.data
    // The request may have been redirected — trust the page object's URL.
    setPage(page, { preserveScroll, replace })
    options.onSuccess?.(page)
  } catch (error) {
    if (handleErrorResponse(error as AxiosError, options)) return
    throw error
  } finally {
    options.onFinish?.()
  }
}

/**
 * Returns true when the error was handled (validation, auth, stale assets).
 */
function handleErrorResponse(error: AxiosError, options: VisitOptions): boolean {
  const response = error.response
  if (!response) return false

  // 409 + X-Inertia-Location → asset version changed, do a hard visit.
  if (response.status === 409 && response.headers['x-inertia-location']) {
    window.location.href = response.headers['x-inertia-location'] as string
    return true
  }

  // 422 → validation error bag for useForm.
  if (response.status === 422) {
    const errors = ((response.data as { errors?: ErrorBag })?.errors ?? {}) as ErrorBag
    const message = (response.data as { message?: string })?.message
    if (Object.keys(errors).length === 0 && message) {
      errors._error = [message]
    }
    options.onError?.(errors)
    return true
  }

  // 401 → session expired; go to login.
  if (response.status === 401) {
    const redirect = (response.data as { redirect?: string })?.redirect ?? '/login'
    window.location.href = redirect
    return true
  }

  return false
}

/** Laravel/Inertia-style router facade. */
export const router = {
  visit,
  get: (url: string, data?: Record<string, unknown>, options: VisitOptions = {}) =>
    visit(url, { ...options, method: 'get', data }),
  post: (url: string, data?: Record<string, unknown>, options: VisitOptions = {}) =>
    visit(url, { ...options, method: 'post', data }),
  put: (url: string, data?: Record<string, unknown>, options: VisitOptions = {}) =>
    visit(url, { ...options, method: 'put', data }),
  patch: (url: string, data?: Record<string, unknown>, options: VisitOptions = {}) =>
    visit(url, { ...options, method: 'patch', data }),
  delete: (url: string, options: VisitOptions = {}) =>
    visit(url, { ...options, method: 'delete' }),
  reload: (options: VisitOptions = {}) =>
    visit(window.location.pathname + window.location.search, { ...options, replace: true, preserveScroll: true }),
}
