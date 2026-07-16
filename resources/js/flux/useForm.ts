// Form helper mirroring Inertia's useForm(): reactive fields, a processing
// flag, and a validation error bag filled automatically from Go's 422
// responses (req.Validate on the server side).
//
//   const form = useForm({ name: '', email: '' })
//   form.post('/users')            // errors land in form.errors.email etc.

import { reactive } from 'vue'
import { router } from './router'
import type { ErrorBag, VisitMethod, VisitOptions } from './types'

type FormOptions = Omit<VisitOptions, 'method' | 'data' | 'onError'>

export interface FluxForm<T extends Record<string, unknown>> {
  data: T
  errors: Partial<Record<keyof T | '_error' | '_body', string>>
  processing: boolean
  /** First message for a field ('' when clean). */
  error(field: keyof T | '_error'): string
  clearErrors(): void
  reset(): void
  submit(method: VisitMethod, url: string, options?: FormOptions): Promise<void>
  get(url: string, options?: FormOptions): Promise<void>
  post(url: string, options?: FormOptions): Promise<void>
  put(url: string, options?: FormOptions): Promise<void>
  patch(url: string, options?: FormOptions): Promise<void>
  delete(url: string, options?: FormOptions): Promise<void>
}

export function useForm<T extends Record<string, unknown>>(initial: T): FluxForm<T> {
  const defaults = { ...initial }

  const form = reactive({
    data: { ...initial },
    errors: {} as Record<string, string>,
    processing: false,

    error(field: string): string {
      return form.errors[field] ?? ''
    },

    clearErrors() {
      form.errors = {}
    },

    reset() {
      Object.assign(form.data, defaults)
      form.clearErrors()
    },

    async submit(method: VisitMethod, url: string, options: FormOptions = {}) {
      if (form.processing) return
      form.processing = true
      form.clearErrors()

      await router.visit(url, {
        ...options,
        method,
        data: { ...form.data },
        onError: (bag: ErrorBag) => {
          const flat: Record<string, string> = {}
          for (const [field, messages] of Object.entries(bag)) {
            flat[field] = messages[0] ?? 'Invalid value.'
          }
          form.errors = flat
        },
        onFinish: () => {
          form.processing = false
          options.onFinish?.()
        },
      })
    },

    get: (url: string, options?: FormOptions) => form.submit('get', url, options),
    post: (url: string, options?: FormOptions) => form.submit('post', url, options),
    put: (url: string, options?: FormOptions) => form.submit('put', url, options),
    patch: (url: string, options?: FormOptions) => form.submit('patch', url, options),
    delete: (url: string, options?: FormOptions) => form.submit('delete', url, options),
  })

  return form as unknown as FluxForm<T>
}
