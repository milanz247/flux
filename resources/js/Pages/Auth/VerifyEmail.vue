<script setup lang="ts">
// Props arrive automatically from req.View("Auth/VerifyEmail", dto.VerifyEmailPageDTO{...})
import { MailCheck } from 'lucide-vue-next'
import AuthLayout from '@/Layouts/AuthLayout.vue'
import { Link, router, useAuth } from '@/flux'
import { Button } from '@/Components/ui/button'

defineProps<{ status?: string }>()

const user = useAuth()

function resend() {
  void router.post('/verify-email/resend')
}
</script>

<template>
  <AuthLayout
    title="Check your email"
    description="We've sent a verification link to your email address."
  >
    <div class="flex flex-col items-center gap-5 text-center">
      <div class="flex h-16 w-16 items-center justify-center rounded-full bg-secondary">
        <MailCheck class="h-7 w-7 text-muted-foreground" />
      </div>

      <p class="max-w-xs text-sm text-muted-foreground">
        Click the link in the email to verify your account.
        <span v-if="user">
          Signed in as <strong class="text-foreground">{{ user.email }}</strong>.
        </span>
      </p>

      <div
        v-if="status === 'sent'"
        class="w-full rounded-lg border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-emerald-800 dark:border-emerald-800 dark:bg-emerald-950/50 dark:text-emerald-300"
      >
        A fresh verification link has been sent.
      </div>

      <Button v-if="user" id="resend-verify" class="w-full" @click="resend">
        Resend verification email
      </Button>
    </div>

    <template #footer>
      <Link href="/dashboard" class="font-medium text-foreground underline-offset-4 hover:underline">
        Go to dashboard
      </Link>
    </template>
  </AuthLayout>
</template>
