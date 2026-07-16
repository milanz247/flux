<script setup lang="ts">
// Props arrive automatically from req.View("Auth/Register", dto.RegisterPageDTO{})
import AuthLayout from '@/Layouts/AuthLayout.vue'
import { Link, useForm } from '@/flux'
import { Button } from '@/Components/ui/button'
import { Input } from '@/Components/ui/input'
import { Label } from '@/Components/ui/label'
import InputError from '@/Components/InputError.vue'

const form = useForm({
  name: '',
  email: '',
  password: '',
  passwordConfirmation: '',
})

function submit() {
  void form.post('/register')
}
</script>

<template>
  <AuthLayout title="Create an account" description="Enter your details below to create your account.">
    <form class="space-y-5" @submit.prevent="submit">
      <div class="space-y-2">
        <Label for="register-name">Full name</Label>
        <Input
          id="register-name"
          v-model="form.data.name"
          placeholder="John Doe"
          autocomplete="name"
        />
        <InputError :message="form.error('name')" />
      </div>

      <div class="space-y-2">
        <Label for="register-email">Email address</Label>
        <Input
          id="register-email"
          v-model="form.data.email"
          type="email"
          placeholder="name@example.com"
          autocomplete="email"
        />
        <InputError :message="form.error('email')" />
      </div>

      <div class="space-y-2">
        <Label for="register-password">Password</Label>
        <Input
          id="register-password"
          v-model="form.data.password"
          type="password"
          placeholder="••••••••"
          autocomplete="new-password"
        />
        <InputError :message="form.error('password')" />
      </div>

      <div class="space-y-2">
        <Label for="register-password-confirm">Confirm password</Label>
        <Input
          id="register-password-confirm"
          v-model="form.data.passwordConfirmation"
          type="password"
          placeholder="••••••••"
          autocomplete="new-password"
        />
        <InputError :message="form.error('passwordConfirmation')" />
      </div>

      <Button id="register-submit" type="submit" class="w-full" :disabled="form.processing">
        {{ form.processing ? 'Creating account…' : 'Create account' }}
      </Button>
    </form>

    <template #footer>
      Already have an account?
      <Link href="/login" class="font-medium text-foreground underline-offset-4 hover:underline">
        Log in
      </Link>
    </template>
  </AuthLayout>
</template>
