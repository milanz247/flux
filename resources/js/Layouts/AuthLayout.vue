<script setup lang="ts">
// Centered card layout for guest pages (login, register, password flows).
import { Zap } from 'lucide-vue-next'
import { useAppName } from '@/flux'

defineProps<{ title: string; description?: string }>()

const appName = useAppName()
</script>

<template>
  <div class="flex min-h-screen flex-col items-center justify-center bg-muted/40 px-4 py-12">
    <div class="mb-8 flex items-center gap-2">
      <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-primary text-primary-foreground">
        <Zap class="h-5 w-5" />
      </div>
      <span class="text-2xl font-bold tracking-tight">{{ appName }}</span>
    </div>

    <div class="w-full max-w-md rounded-xl border bg-card p-8 text-card-foreground shadow-sm">
      <div class="mb-6 space-y-1 text-center">
        <h1 class="text-xl font-semibold tracking-tight">{{ title }}</h1>
        <p v-if="description" class="text-sm text-muted-foreground">{{ description }}</p>
      </div>
      <slot />
    </div>

    <div class="mt-6 text-center text-sm text-muted-foreground">
      <slot name="footer" />
    </div>
  </div>
</template>
