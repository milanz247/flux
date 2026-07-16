<script setup lang="ts">
// Props arrive automatically from req.View("Auth/ResetPassword", dto.ResetPasswordPageDTO{...})
import AuthLayout from '@/Layouts/AuthLayout.vue'
import { useForm } from '@/flux'
import { Button } from '@/Components/ui/button'
import { Input } from '@/Components/ui/input'
import { Label } from '@/Components/ui/label'
import InputError from '@/Components/InputError.vue'

const props = defineProps<{ token: string; email: string }>()

const form = useForm({
  token: props.token,
  password: '',
  passwordConfirmation: '',
})

function submit() {
  void form.post('/reset-password')
}
</script>

<template>
  <AuthLayout title="Reset password" :description="`Choose a new password for ${email}.`">
    <form class="space-y-4" @submit.prevent="submit">
      <div class="space-y-2">
        <Label for="password">New password</Label>
        <Input id="password" v-model="form.data.password" type="password" autocomplete="new-password" />
        <InputError :message="form.error('password')" />
      </div>

      <div class="space-y-2">
        <Label for="passwordConfirmation">Confirm new password</Label>
        <Input
          id="passwordConfirmation"
          v-model="form.data.passwordConfirmation"
          type="password"
          autocomplete="new-password"
        />
        <InputError :message="form.error('passwordConfirmation')" />
        <InputError :message="form.error('token' as never)" />
        <InputError :message="form.error('_error')" />
      </div>

      <Button type="submit" class="w-full" :disabled="form.processing">
        {{ form.processing ? 'Resetting…' : 'Reset password' }}
      </Button>
    </form>
  </AuthLayout>
</template>
