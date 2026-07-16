<script setup lang="ts">
// Props arrive automatically from req.View("Auth/ForgotPassword", dto.ForgotPasswordPageDTO{...})
import AuthLayout from '@/Layouts/AuthLayout.vue'
import { Link, useForm } from '@/flux'
import { Button } from '@/Components/ui/button'
import { Input } from '@/Components/ui/input'
import { Label } from '@/Components/ui/label'
import { Alert, AlertDescription } from '@/Components/ui/alert'
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
    description="Enter your email and we'll send you a reset link."
  >
    <Alert v-if="status === 'sent'" variant="success" class="mb-4">
      <AlertDescription>
        If that address exists, a reset link is on its way.
        <span class="block text-xs opacity-80">(MAIL_DRIVER=log writes it to storage/logs/mail.log)</span>
      </AlertDescription>
    </Alert>

    <form class="space-y-4" @submit.prevent="submit">
      <div class="space-y-2">
        <Label for="email">Email</Label>
        <Input id="email" v-model="form.data.email" type="email" placeholder="you@example.com" autocomplete="email" />
        <InputError :message="form.error('email')" />
      </div>

      <Button type="submit" class="w-full" :disabled="form.processing">
        {{ form.processing ? 'Sending…' : 'Send reset link' }}
      </Button>
    </form>

    <template #footer>
      Remembered it?
      <Link href="/login" class="font-medium text-foreground hover:underline">Back to sign in</Link>
    </template>
  </AuthLayout>
</template>
