<script setup lang="ts">
// Props arrive automatically from req.View("Auth/Login", dto.LoginPageDTO{...})
import AuthLayout from '@/Layouts/AuthLayout.vue'
import { Link, useForm } from '@/flux'
import { Button } from '@/Components/ui/button'
import { Input } from '@/Components/ui/input'
import { Label } from '@/Components/ui/label'
import { Checkbox } from '@/Components/ui/checkbox'
import InputError from '@/Components/InputError.vue'

defineProps<{ status?: string }>()

const form = useForm({
  email: '',
  password: '',
})

function submit() {
  void form.post('/login')
}
</script>

<template>
  <AuthLayout title="Log in to your account" description="Enter your email and password to access your dashboard.">
    <div
      v-if="status === 'password-reset'"
      class="mb-6 rounded-lg border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-emerald-800 dark:border-emerald-800 dark:bg-emerald-950/50 dark:text-emerald-300"
    >
      Your password has been reset. Sign in with your new password.
    </div>

    <form class="space-y-5" @submit.prevent="submit">
      <div class="space-y-2">
        <Label for="login-email">Email address</Label>
        <Input
          id="login-email"
          v-model="form.data.email"
          type="email"
          placeholder="name@example.com"
          autocomplete="email"
        />
        <InputError :message="form.error('email')" />
      </div>

      <div class="space-y-2">
        <div class="flex items-center justify-between">
          <Label for="login-password">Password</Label>
          <Link href="/forgot-password" class="text-xs font-medium text-muted-foreground transition-colors hover:text-foreground">
            Forgot password?
          </Link>
        </div>
        <Input
          id="login-password"
          v-model="form.data.password"
          type="password"
          placeholder="••••••••"
          autocomplete="current-password"
        />
        <InputError :message="form.error('password')" />
        <InputError :message="form.error('_error')" />
      </div>

      <Button id="login-submit" type="submit" class="w-full" :disabled="form.processing">
        {{ form.processing ? 'Signing in…' : 'Log in' }}
      </Button>
    </form>

    <template #footer>
      Don't have an account?
      <Link href="/register" class="font-medium text-foreground underline-offset-4 hover:underline">
        Sign up
      </Link>
    </template>
  </AuthLayout>
</template>
