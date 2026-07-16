<script setup lang="ts">
// Professional dashboard shell: fixed sidebar navigation, topbar with the
// current page title and a user menu, and a content slot.
import { ref } from 'vue'
import { LayoutDashboard, Users, LogOut, Menu, Zap, ChevronDown, MailWarning } from 'lucide-vue-next'
import { Link, router, useAuth, useAppName, usePage } from '@/flux'
import { Button } from '@/Components/ui/button'
import { cn } from '@/lib/utils'

defineProps<{ title?: string }>()

const user = useAuth()
const appName = useAppName()
const page = usePage()

const sidebarOpen = ref(false)
const userMenuOpen = ref(false)

const navigation = [
  { name: 'Dashboard', href: '/dashboard', icon: LayoutDashboard },
  { name: 'Users', href: '/users', icon: Users },
]

function isActive(href: string): boolean {
  const url = page.value?.url ?? ''
  return url === href || url.startsWith(href + '/') || url.startsWith(href + '?')
}

function logout() {
  void router.post('/logout')
}

function resendVerification() {
  void router.post('/verify-email/resend')
}
</script>

<template>
  <div class="min-h-screen bg-muted/40">
    <!-- Sidebar -->
    <aside
      :class="
        cn(
          'fixed inset-y-0 left-0 z-40 w-64 transform border-r bg-background transition-transform duration-200 lg:translate-x-0',
          sidebarOpen ? 'translate-x-0' : '-translate-x-full',
        )
      "
    >
      <div class="flex h-16 items-center gap-2 border-b px-6">
        <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-primary text-primary-foreground">
          <Zap class="h-4 w-4" />
        </div>
        <span class="text-lg font-semibold tracking-tight">{{ appName }}</span>
      </div>

      <nav class="space-y-1 p-4">
        <Link
          v-for="item in navigation"
          :key="item.href"
          :href="item.href"
          :class="
            cn(
              'flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium transition-colors',
              isActive(item.href)
                ? 'bg-primary text-primary-foreground'
                : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground',
            )
          "
        >
          <component :is="item.icon" class="h-4 w-4" />
          {{ item.name }}
        </Link>
      </nav>
    </aside>

    <!-- Mobile sidebar backdrop -->
    <div
      v-if="sidebarOpen"
      class="fixed inset-0 z-30 bg-black/50 lg:hidden"
      @click="sidebarOpen = false"
    />

    <!-- Main column -->
    <div class="lg:pl-64">
      <!-- Topbar -->
      <header class="sticky top-0 z-20 flex h-16 items-center gap-4 border-b bg-background/95 px-4 backdrop-blur sm:px-6">
        <Button variant="ghost" size="icon" class="lg:hidden" @click="sidebarOpen = !sidebarOpen">
          <Menu class="h-5 w-5" />
        </Button>

        <h1 class="flex-1 text-lg font-semibold">{{ title }}</h1>

        <!-- User menu -->
        <div class="relative">
          <button
            class="flex items-center gap-2 rounded-md px-2 py-1.5 text-sm hover:bg-accent"
            @click="userMenuOpen = !userMenuOpen"
          >
            <span
              class="flex h-8 w-8 items-center justify-center rounded-full bg-primary text-xs font-semibold text-primary-foreground"
            >
              {{ user?.name?.slice(0, 1).toUpperCase() ?? '?' }}
            </span>
            <span class="hidden sm:block">{{ user?.name }}</span>
            <ChevronDown class="h-4 w-4 text-muted-foreground" />
          </button>

          <div
            v-if="userMenuOpen"
            class="absolute right-0 mt-2 w-56 rounded-md border bg-popover p-1 text-popover-foreground shadow-md"
            @click="userMenuOpen = false"
          >
            <div class="border-b px-3 py-2">
              <p class="text-sm font-medium">{{ user?.name }}</p>
              <p class="truncate text-xs text-muted-foreground">{{ user?.email }}</p>
            </div>
            <button
              class="mt-1 flex w-full items-center gap-2 rounded-sm px-3 py-2 text-sm hover:bg-accent"
              @click="logout"
            >
              <LogOut class="h-4 w-4" />
              Log out
            </button>
          </div>
        </div>
      </header>

      <!-- Unverified email banner -->
      <div
        v-if="user && !user.emailVerifiedAt"
        class="flex items-center gap-3 border-b bg-amber-50 px-6 py-2.5 text-sm text-amber-900 dark:bg-amber-950 dark:text-amber-200"
      >
        <MailWarning class="h-4 w-4 shrink-0" />
        <span class="flex-1">Your email address is not verified. Check your inbox for the link.</span>
        <Button variant="outline" size="sm" @click="resendVerification">Resend email</Button>
      </div>

      <!-- Page content -->
      <main class="p-4 sm:p-6 lg:p-8">
        <slot />
      </main>
    </div>
  </div>
</template>
