<script setup lang="ts">
// Professional dashboard shell: fixed sidebar navigation, topbar with the
// current page title and a user menu, and a content slot.
import { ref } from 'vue'
import { LogOut, Menu, MailWarning, PanelLeftClose, PanelLeftOpen } from 'lucide-vue-next'
import { router, useAuth } from '@/flux'
import { Button } from '@/Components/ui/button'
import { Alert, AlertDescription } from '@/Components/ui/alert'
import { Avatar, AvatarFallback } from '@/Components/ui/avatar'
import { Sheet, SheetContent } from '@/Components/ui/sheet'
import {
  DropdownMenu, DropdownMenuTrigger, DropdownMenuContent,
  DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuItem,
} from '@/Components/ui/dropdown-menu'
import { cn } from '@/lib/utils'
import SidebarNav from '@/Components/SidebarNav.vue'

defineProps<{ title?: string }>()

const user = useAuth()
const sidebarOpen = ref(false)
const sidebarCollapsed = ref(localStorage.getItem('flux-sidebar-collapsed') === 'true')

function toggleSidebarCollapsed() {
  sidebarCollapsed.value = !sidebarCollapsed.value
  localStorage.setItem('flux-sidebar-collapsed', String(sidebarCollapsed.value))
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
    <!-- Desktop sidebar -->
    <aside
      :class="
        cn(
          'fixed inset-y-0 left-0 z-40 hidden border-r bg-background transition-[width] duration-200 lg:block',
          sidebarCollapsed ? 'w-16' : 'w-64',
        )
      "
    >
      <SidebarNav :collapsed="sidebarCollapsed" />
    </aside>

    <!-- Mobile sidebar -->
    <Sheet v-model:open="sidebarOpen">
      <SheetContent side="left" class="w-64 p-0">
        <SidebarNav />
      </SheetContent>
    </Sheet>

    <!-- Main column -->
    <div :class="cn('transition-[padding] duration-200', sidebarCollapsed ? 'lg:pl-16' : 'lg:pl-64')">
      <!-- Topbar -->
      <header class="sticky top-0 z-20 flex h-16 items-center gap-4 border-b bg-background/95 px-4 backdrop-blur sm:px-6">
        <Button variant="ghost" size="icon" class="lg:hidden" @click="sidebarOpen = true">
          <Menu class="h-5 w-5" />
        </Button>

        <Button
          variant="ghost"
          size="icon"
          class="hidden lg:inline-flex"
          :title="sidebarCollapsed ? 'Expand sidebar' : 'Collapse sidebar'"
          @click="toggleSidebarCollapsed"
        >
          <PanelLeftOpen v-if="sidebarCollapsed" class="h-5 w-5" />
          <PanelLeftClose v-else class="h-5 w-5" />
        </Button>

        <h1 class="flex-1 text-lg font-semibold">{{ title }}</h1>

        <!-- User menu -->
        <DropdownMenu>
          <DropdownMenuTrigger as-child>
            <button class="flex items-center gap-2 rounded-md px-2 py-1.5 text-sm hover:bg-accent">
              <Avatar class="h-8 w-8">
                <AvatarFallback>{{ user?.name?.slice(0, 1).toUpperCase() ?? '?' }}</AvatarFallback>
              </Avatar>
              <span class="hidden sm:block">{{ user?.name }}</span>
            </button>
          </DropdownMenuTrigger>
          <DropdownMenuContent class="w-56">
            <DropdownMenuLabel>
              <p class="text-sm font-medium">{{ user?.name }}</p>
              <p class="truncate text-xs font-normal text-muted-foreground">{{ user?.email }}</p>
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuItem @click="logout">
              <LogOut class="h-4 w-4" />
              Log out
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </header>

      <!-- Unverified email banner -->
      <div v-if="user && !user.emailVerifiedAt" class="border-b px-4 py-2.5 sm:px-6">
        <Alert variant="default" class="flex items-center gap-3 border-amber-200 bg-amber-50 text-amber-900 dark:border-amber-900 dark:bg-amber-950 dark:text-amber-200">
          <span class="shrink-0"><MailWarning class="h-4 w-4" /></span>
          <AlertDescription class="flex flex-1 items-center justify-between gap-3">
            <span>Your email address is not verified. Check your inbox for the link.</span>
            <Button variant="outline" size="sm" @click="resendVerification">Resend email</Button>
          </AlertDescription>
        </Alert>
      </div>

      <!-- Page content -->
      <main class="p-4 sm:p-6 lg:p-8">
        <slot />
      </main>
    </div>
  </div>
</template>
