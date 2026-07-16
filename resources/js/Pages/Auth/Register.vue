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
  <AuthLayout title="Create an account" description="A verification link will be emailed to you.">
    <form class="space-y-4" @submit.prevent="submit">
      <div class="space-y-2">
        <Label for="name">Name</Label>
        <Input id="name" v-model="form.data.name" placeholder="Jane Doe" autocomplete="name" />
        <InputError :message="form.error('name')" />
      </div>

      <div class="space-y-2">
        <Label for="email">Email</Label>
        <Input id="email" v-model="form.data.email" type="email" placeholder="you@example.com" autocomplete="email" />
        <InputError :message="form.error('email')" />
      </div>

      <div class="space-y-2">
        <Label for="password">Password</Label>
        <Input id="password" v-model="form.data.password" type="password" autocomplete="new-password" />
        <InputError :message="form.error('password')" />
      </div>

      <div class="space-y-2">
        <Label for="passwordConfirmation">Confirm password</Label>
        <Input
          id="passwordConfirmation"
          v-model="form.data.passwordConfirmation"
          type="password"
          autocomplete="new-password"
        />
        <InputError :message="form.error('passwordConfirmation')" />
      </div>

      <Button type="submit" class="w-full" :disabled="form.processing">
        {{ form.processing ? 'Creating account…' : 'Create account' }}
      </Button>
    </form>

    <template #footer>
      Already registered?
      <Link href="/login" class="font-medium text-foreground hover:underline">Sign in</Link>
    </template>
  </AuthLayout>
</template>
