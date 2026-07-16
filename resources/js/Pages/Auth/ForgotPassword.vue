<script setup lang="ts">
// Props arrive automatically from req.View("Auth/ForgotPassword", dto.ForgotPasswordPageDTO{...})
import AuthLayout from '@/Layouts/AuthLayout.vue'
import { Link, useForm } from '@/flux'
import { Button } from '@/Components/ui/button'
import { Input } from '@/Components/ui/input'
import { Label } from '@/Components/ui/label'
import InputError from '@/Components/InputError.vue'

defineProps<{ status?: string }>()

const form = useForm({ email: '' })

function submit() {
  void form.post('/forgot-password')
}
</script>

<template>
  <AuthLayout
    title="Forgot your password?"
    description="Enter your email address and we'll send you a password reset link."
  >
    <div
      v-if="status === 'sent'"
      class="mb-6 rounded-lg border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-emerald-800 dark:border-emerald-800 dark:bg-emerald-950/50 dark:text-emerald-300"
    >
      If that address exists, a reset link is on its way.
    </div>

    <form class="space-y-5" @submit.prevent="submit">
      <div class="space-y-2">
        <Label for="forgot-email">Email address</Label>
        <Input
          id="forgot-email"
          v-model="form.data.email"
          type="email"
          placeholder="name@example.com"
          autocomplete="email"
        />
        <InputError :message="form.error('email')" />
      </div>

      <Button id="forgot-submit" type="submit" class="w-full" :disabled="form.processing">
        {{ form.processing ? 'Sending…' : 'Send reset link' }}
      </Button>
    </form>

    <template #footer>
      Remember your password?
      <Link href="/login" class="font-medium text-foreground underline-offset-4 hover:underline">
        Back to log in
      </Link>
    </template>
  </AuthLayout>
</template>
