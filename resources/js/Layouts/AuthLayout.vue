<script setup lang="ts">
// Centered card layout for guest pages (login, register, password flows).
import { useAppName, useFlash } from '@/flux'
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/Components/ui/card'
import { Alert, AlertDescription } from '@/Components/ui/alert'
import AppLogo from '@/Components/AppLogo.vue'

defineProps<{ title: string; description?: string }>()

const appName = useAppName()
const flash = useFlash()
</script>

<template>
  <div class="flex min-h-screen flex-col items-center justify-center bg-muted/40 px-4 py-12">
    <div class="mb-8 flex items-center gap-2">
      <AppLogo :size="40" />
      <span class="text-2xl font-bold tracking-tight">{{ appName }}</span>
    </div>

    <Alert
      v-if="flash"
      :variant="flash.type === 'error' ? 'destructive' : 'default'"
      class="mb-4 w-full max-w-md"
    >
      <AlertDescription>{{ flash.message }}</AlertDescription>
    </Alert>

    <Card class="w-full max-w-md">
      <CardHeader class="text-center">
        <CardTitle class="text-xl">{{ title }}</CardTitle>
        <CardDescription v-if="description">{{ description }}</CardDescription>
      </CardHeader>
      <CardContent>
        <slot />
      </CardContent>
    </Card>

    <div class="mt-6 text-center text-sm text-muted-foreground">
      <slot name="footer" />
    </div>
  </div>
</template>
