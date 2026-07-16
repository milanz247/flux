<script setup lang="ts">
// Props arrive automatically from req.View("Auth/VerifyEmail", dto.VerifyEmailPageDTO{...})
import { MailCheck } from 'lucide-vue-next'
import AuthLayout from '@/Layouts/AuthLayout.vue'
import { Link, router, useAuth } from '@/flux'
import { Button } from '@/Components/ui/button'
import { Alert, AlertDescription } from '@/Components/ui/alert'

defineProps<{ status?: string }>()

const user = useAuth()

function resend() {
  void router.post('/verify-email/resend')
}
</script>

<template>
  <AuthLayout
    title="Verify your email"
    description="We sent a verification link to your email address."
  >
    <div class="flex flex-col items-center gap-4 text-center">
      <div class="flex h-14 w-14 items-center justify-center rounded-full bg-secondary">
        <MailCheck class="h-6 w-6" />
      </div>

      <p class="text-sm text-muted-foreground">
        Click the link in the email to activate your account.
        <span v-if="user">Signed in as <strong>{{ user.email }}</strong>.</span>
        <span class="mt-1 block text-xs">(MAIL_DRIVER=log writes the link to storage/logs/mail.log)</span>
      </p>

      <Alert v-if="status === 'sent'" variant="success" class="w-full">
        <AlertDescription>A fresh verification link has been sent.</AlertDescription>
      </Alert>

      <Button v-if="user" class="w-full" @click="resend">Resend verification email</Button>
    </div>

    <template #footer>
      <Link href="/dashboard" class="font-medium text-foreground hover:underline">Go to dashboard</Link>
    </template>
  </AuthLayout>
</template>
