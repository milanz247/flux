<script setup lang="ts">
// Props arrive automatically from req.View("Welcome", dto.WelcomePageDTO{})
// Public marketing/landing page — the auth state (guest vs signed-in) is
// read from shared props, same as every other page.
import { Server, Layers, Blocks } from 'lucide-vue-next'
import { Link, useAppName, useAuth } from '@/flux'
import { Button } from '@/Components/ui/button'
import { Card } from '@/Components/ui/card'
import AppLogo from '@/Components/AppLogo.vue'

const user = useAuth()
const appName = useAppName()

const highlights = [
  { icon: Server, text: 'A Go + Gin backend with Laravel-style routing, DTOs and validation.' },
  { icon: Layers, text: 'Server-driven pages via a lightweight Inertia-style adapter.' },
  { icon: Blocks, text: 'shadcn-vue components styled with Tailwind, ready to extend.' },
]
</script>

<template>
  <div class="min-h-screen bg-background">
    <header class="mx-auto flex max-w-5xl items-center justify-between px-6 py-8">
      <div class="flex items-center gap-2">
        <AppLogo :size="28" />
        <span class="text-sm font-semibold tracking-tight">{{ appName }}</span>
      </div>

      <nav class="flex items-center gap-4 text-sm">
        <template v-if="user">
          <Link href="/dashboard">
            <Button size="sm">Dashboard</Button>
          </Link>
        </template>
        <template v-else>
          <Link href="/login" class="text-foreground hover:text-foreground/80">Log in</Link>
          <Link href="/register">
            <Button variant="outline" size="sm">Register</Button>
          </Link>
        </template>
      </nav>
    </header>

    <main class="mx-auto flex max-w-5xl flex-col items-center px-6 pb-20 pt-6">
      <Card class="grid w-full overflow-hidden p-0 shadow-lg md:grid-cols-2">
        <!-- Content -->
        <div class="flex flex-col justify-center gap-8 p-10">
          <div>
            <h1 class="text-lg font-semibold">Let's get started</h1>
            <p class="mt-2 text-sm text-muted-foreground">
              {{ appName }} has everything you need to ship a full-stack Go app.
              Here's what's included.
            </p>
          </div>

          <ul class="space-y-4">
            <li v-for="item in highlights" :key="item.text" class="flex items-start gap-3">
              <span class="mt-0.5 flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-muted">
                <component :is="item.icon" class="h-3.5 w-3.5 text-foreground" />
              </span>
              <span class="text-sm text-muted-foreground">{{ item.text }}</span>
            </li>
          </ul>

          <div class="flex items-center gap-3">
            <Link :href="user ? '/dashboard' : '/register'">
              <Button>{{ user ? 'Go to dashboard' : 'Get started' }}</Button>
            </Link>
            <Link v-if="!user" href="/login">
              <Button variant="ghost">Sign in</Button>
            </Link>
          </div>
        </div>

        <!-- Hero art -->
        <div class="relative order-first min-h-[240px] overflow-hidden md:order-last md:min-h-[460px]">
          <div class="absolute inset-0 bg-gradient-to-br from-indigo-500 via-violet-500 to-cyan-400" />

          <!-- Vertical accent bars -->
          <div class="absolute inset-y-0 left-0 flex w-24 opacity-70">
            <div class="h-full w-1/4 bg-black/20" />
            <div class="h-full w-1/4 bg-white/20" />
            <div class="h-full w-1/4 bg-black/10" />
            <div class="h-full w-1/4 bg-white/10" />
          </div>

          <!-- Soft glows -->
          <div class="pulse-ring absolute -right-16 -top-16 h-64 w-64 rounded-full bg-white/20 blur-2xl" />
          <div class="pulse-ring absolute -bottom-20 left-1/3 h-72 w-72 rounded-full bg-black/10 blur-3xl" />

          <!-- Brand mark -->
          <div class="relative flex h-full flex-col items-center justify-center gap-4 p-10 text-center">
            <AppLogo :size="88" />
            <span class="text-4xl font-bold tracking-tight text-white drop-shadow-sm">{{ appName }}</span>
            <span class="text-sm font-medium text-white/80">Go · Vue · shadcn-vue</span>
          </div>
        </div>
      </Card>
    </main>
  </div>
</template>
