<script setup lang="ts">
// Clean dashboard shell matching the Laravel Vue starter kit: fixed sidebar
// with smooth transitions, topbar with breadcrumb and user dropdown.
import { LayoutDashboard, LogOut, Menu, Zap, Settings, MailWarning } from 'lucide-vue-next'
import { Link, router, useAuth, useAppName, usePage } from '@/flux'
import { Button } from '@/Components/ui/button'
import { Avatar, AvatarFallback } from '@/Components/ui/avatar'
import {
  DropdownMenu, DropdownMenuTrigger, DropdownMenuContent,
  DropdownMenuItem, DropdownMenuSeparator, DropdownMenuLabel,
} from '@/Components/ui/dropdown-menu'
import { Sheet, SheetTrigger, SheetContent, SheetClose } from '@/Components/ui/sheet'
import { cn } from '@/lib/utils'

defineProps<{ title?: string }>()

const user = useAuth()
const appName = useAppName()
const page = usePage()

const navigation = [
  { name: 'Dashboard', href: '/dashboard', icon: LayoutDashboard },
]

function isActive(href: string): boolean {
  const url = page.value?.url ?? ''
  return url === href || url === '/' && href === '/dashboard'
}

function logout() {
  void router.post('/logout')
}

function resendVerification() {
  void router.post('/verify-email/resend')
}

function userInitials(): string {
  if (!user.value) return '?'
  return user.value.name
    .split(' ')
    .map((part: string) => part[0])
    .slice(0, 2)
    .join('')
    .toUpperCase()
}
</script>

<template>
  <div class="min-h-screen bg-background">
    <!-- Mobile sidebar (Sheet) -->
    <Sheet>
      <div class="lg:hidden">
        <SheetTrigger>
          <Button id="mobile-menu-toggle" variant="ghost" size="icon" class="fixed left-4 top-4 z-50 lg:hidden">
            <Menu class="h-5 w-5" />
          </Button>
        </SheetTrigger>
      </div>
      <SheetContent side="left" class="w-64 p-0">
        <SheetClose />
        <div class="flex h-full flex-col">
          <!-- Mobile sidebar header -->
          <div class="flex h-16 items-center gap-2 border-b px-6">
            <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-primary text-primary-foreground">
              <Zap class="h-4 w-4" />
            </div>
            <span class="text-lg font-semibold tracking-tight">{{ appName }}</span>
          </div>
          <!-- Mobile sidebar nav -->
          <nav class="flex-1 space-y-1 p-3">
            <Link
              v-for="item in navigation"
              :key="item.href"
              :href="item.href"
              :class="
                cn(
                  'flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-colors',
                  isActive(item.href)
                    ? 'bg-sidebar-accent text-sidebar-accent-foreground'
                    : 'text-sidebar-foreground/70 hover:bg-sidebar-accent hover:text-sidebar-accent-foreground',
                )
              "
            >
              <component :is="item.icon" class="h-4 w-4" />
              {{ item.name }}
            </Link>
          </nav>
        </div>
      </SheetContent>
    </Sheet>

    <!-- Desktop sidebar -->
    <aside class="fixed inset-y-0 left-0 z-30 hidden w-64 border-r bg-sidebar lg:block">
      <!-- Sidebar header -->
      <div class="flex h-16 items-center gap-2 border-b px-6">
        <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-primary text-primary-foreground">
          <Zap class="h-4 w-4" />
        </div>
        <span class="text-lg font-semibold tracking-tight text-sidebar-foreground">{{ appName }}</span>
      </div>

      <!-- Navigation -->
      <nav class="flex-1 space-y-1 p-3">
        <Link
          v-for="item in navigation"
          :key="item.href"
          :href="item.href"
          :class="
            cn(
              'flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-colors',
              isActive(item.href)
                ? 'bg-sidebar-accent text-sidebar-accent-foreground'
                : 'text-sidebar-foreground/70 hover:bg-sidebar-accent hover:text-sidebar-accent-foreground',
            )
          "
        >
          <component :is="item.icon" class="h-4 w-4" />
          {{ item.name }}
        </Link>
      </nav>

      <!-- Sidebar footer: user info -->
      <div class="absolute bottom-0 left-0 right-0 border-t p-3">
        <div class="flex items-center gap-3 rounded-lg px-3 py-2">
          <Avatar class="h-8 w-8">
            <AvatarFallback class="bg-primary/10 text-xs text-primary">
              {{ userInitials() }}
            </AvatarFallback>
          </Avatar>
          <div class="min-w-0 flex-1">
            <p class="truncate text-sm font-medium text-sidebar-foreground">{{ user?.name }}</p>
            <p class="truncate text-xs text-sidebar-foreground/60">{{ user?.email }}</p>
          </div>
        </div>
      </div>
    </aside>

    <!-- Main content area -->
    <div class="lg:pl-64">
      <!-- Topbar -->
      <header class="sticky top-0 z-20 flex h-16 items-center gap-4 border-b bg-background/95 px-4 backdrop-blur supports-[backdrop-filter]:bg-background/60 sm:px-6 lg:px-8">
        <!-- Spacer for mobile menu button -->
        <div class="w-10 lg:hidden" />

        <!-- Page title / breadcrumb -->
        <h1 class="flex-1 text-lg font-semibold">{{ title }}</h1>

        <!-- User dropdown (topbar) -->
        <DropdownMenu>
          <DropdownMenuTrigger>
            <button
              id="user-menu-button"
              data-dropdown-trigger
              class="flex items-center gap-2 rounded-lg px-2 py-1.5 transition-colors hover:bg-accent"
            >
              <Avatar class="h-8 w-8">
                <AvatarFallback class="bg-primary/10 text-xs text-primary">
                  {{ userInitials() }}
                </AvatarFallback>
              </Avatar>
              <span class="hidden text-sm font-medium sm:block">{{ user?.name }}</span>
            </button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" class="w-56">
            <DropdownMenuLabel>
              <div>
                <p class="text-sm font-medium">{{ user?.name }}</p>
                <p class="truncate text-xs font-normal text-muted-foreground">{{ user?.email }}</p>
              </div>
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuItem id="logout-button" @click="logout">
              <LogOut class="h-4 w-4" />
              Log out
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </header>

      <!-- Unverified email banner -->
      <div
        v-if="user && !user.emailVerifiedAt"
        class="flex items-center gap-3 border-b bg-amber-50 px-6 py-2.5 text-sm text-amber-900 dark:bg-amber-950/50 dark:text-amber-200"
      >
        <MailWarning class="h-4 w-4 shrink-0" />
        <span class="flex-1">Your email address is not verified.</span>
        <Button variant="outline" size="sm" @click="resendVerification">Resend email</Button>
      </div>

      <!-- Page content -->
      <main class="p-4 sm:p-6 lg:p-8">
        <slot />
      </main>
    </div>
  </div>
</template>
