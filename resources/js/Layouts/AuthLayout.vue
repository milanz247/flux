<script setup lang="ts">
// Split-screen authentication layout — branded gradient panel on the left,
// clean form card on the right. Matches the Laravel Vue starter kit aesthetic.
import { Zap } from 'lucide-vue-next'
import { useAppName } from '@/flux'

defineProps<{ title: string; description?: string }>()

const appName = useAppName()
</script>

<template>
  <div class="flex min-h-screen">
    <!-- Left: branded gradient panel (hidden on mobile) -->
    <div class="auth-gradient relative hidden flex-col justify-between overflow-hidden p-10 lg:flex lg:w-1/2">
      <!-- Floating decorative elements -->
      <div class="absolute inset-0 overflow-hidden">
        <div class="float-dot absolute left-[10%] top-[20%] h-64 w-64 rounded-full bg-white/5" />
        <div class="float-dot absolute right-[15%] top-[60%] h-48 w-48 rounded-full bg-white/5" />
        <div class="float-dot absolute left-[60%] top-[10%] h-32 w-32 rounded-full bg-white/5" />
        <div class="pulse-ring absolute bottom-[20%] left-[30%] h-80 w-80 rounded-full border border-white/10" />
      </div>

      <!-- Top: Logo -->
      <div class="relative z-10 flex items-center gap-3">
        <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-white/10 text-white backdrop-blur-sm">
          <Zap class="h-5 w-5" />
        </div>
        <span class="text-xl font-bold tracking-tight text-white">{{ appName }}</span>
      </div>

      <!-- Center: Tagline -->
      <div class="relative z-10 max-w-md">
        <blockquote class="space-y-4">
          <p class="text-lg font-medium leading-relaxed text-white/90">
            "A beautifully crafted full-stack framework that brings the elegance of Laravel and Inertia to the Go ecosystem."
          </p>
          <footer class="text-sm text-white/60">
            Built with Go, Vue 3, and shadcn-vue
          </footer>
        </blockquote>
      </div>

      <!-- Bottom: subtle pattern grid -->
      <div class="relative z-10">
        <p class="text-xs text-white/30">© {{ new Date().getFullYear() }} {{ appName }}. All rights reserved.</p>
      </div>
    </div>

    <!-- Right: form panel -->
    <div class="flex w-full flex-col justify-center px-4 py-12 lg:w-1/2 lg:px-8">
      <!-- Mobile logo (shown only on small screens) -->
      <div class="mb-8 flex items-center justify-center gap-2 lg:hidden">
        <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-primary text-primary-foreground">
          <Zap class="h-5 w-5" />
        </div>
        <span class="text-2xl font-bold tracking-tight">{{ appName }}</span>
      </div>

      <div class="mx-auto w-full max-w-sm">
        <!-- Title & description -->
        <div class="mb-8 space-y-2 text-center lg:text-left">
          <h1 class="text-2xl font-semibold tracking-tight">{{ title }}</h1>
          <p v-if="description" class="text-sm text-muted-foreground">{{ description }}</p>
        </div>

        <!-- Form slot -->
        <slot />

        <!-- Footer slot (sign up / sign in links) -->
        <div class="mt-6 text-center text-sm text-muted-foreground">
          <slot name="footer" />
        </div>
      </div>
    </div>
  </div>
</template>
