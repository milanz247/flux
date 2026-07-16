<script setup lang="ts">
// Props arrive automatically from req.View("Auth/Login", dto.LoginPageDTO{...})
import AuthLayout from '@/Layouts/AuthLayout.vue'
import { Link, useForm } from '@/flux'
import { Button } from '@/Components/ui/button'
import { Input } from '@/Components/ui/input'
import { Label } from '@/Components/ui/label'
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
  <AuthLayout title="Welcome back" description="Sign in to your account to continue.">
    <div
      v-if="status === 'password-reset'"
      class="mb-4 rounded-md bg-emerald-50 px-3 py-2 text-sm text-emerald-800 dark:bg-emerald-950 dark:text-emerald-300"
    >
      Your password has been reset. Sign in with your new password.
    </div>

    <form class="space-y-4" @submit.prevent="submit">
      <div class="space-y-2">
        <Label for="email">Email</Label>
        <Input id="email" v-model="form.data.email" type="email" placeholder="you@example.com" autocomplete="email" />
        <InputError :message="form.error('email')" />
      </div>

      <div class="space-y-2">
        <div class="flex items-center justify-between">
          <Label for="password">Password</Label>
          <Link href="/forgot-password" class="text-xs text-muted-foreground hover:text-foreground">
            Forgot password?
          </Link>
        </div>
        <Input id="password" v-model="form.data.password" type="password" autocomplete="current-password" />
        <InputError :message="form.error('password')" />
        <InputError :message="form.error('_error')" />
      </div>

      <Button type="submit" class="w-full" :disabled="form.processing">
        {{ form.processing ? 'Signing in…' : 'Sign in' }}
      </Button>
    </form>

    <template #footer>
      Don't have an account?
      <Link href="/register" class="font-medium text-foreground hover:underline">Sign up</Link>
    </template>
  </AuthLayout>
</template>
