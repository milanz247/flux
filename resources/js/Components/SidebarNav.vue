<script setup lang="ts">
import { LayoutDashboard } from 'lucide-vue-next'
import { Link, useAppName, usePage } from '@/flux'
import { cn } from '@/lib/utils'
import AppLogo from '@/Components/AppLogo.vue'

withDefaults(defineProps<{ collapsed?: boolean }>(), { collapsed: false })

const appName = useAppName()
const page = usePage()

const navigation = [
  { name: 'Dashboard', href: '/dashboard', icon: LayoutDashboard },
]

function isActive(href: string): boolean {
  const url = page.value?.url ?? ''
  return url === href || url.startsWith(href + '/') || url.startsWith(href + '?')
}
</script>

<template>
  <div :class="cn('flex h-16 items-center gap-2 border-b', collapsed ? 'justify-center px-2' : 'px-6')">
    <AppLogo :size="32" class="shrink-0" />
    <span v-if="!collapsed" class="truncate text-lg font-semibold tracking-tight">{{ appName }}</span>
  </div>

  <nav :class="cn('space-y-1 p-4', collapsed && 'px-2')">
    <Link
      v-for="item in navigation"
      :key="item.href"
      :href="item.href"
      :title="collapsed ? item.name : undefined"
      :class="
        cn(
          'flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium transition-colors',
          collapsed && 'justify-center px-0',
          isActive(item.href)
            ? 'bg-primary text-primary-foreground'
            : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground',
        )
      "
    >
      <component :is="item.icon" class="h-4 w-4 shrink-0" />
      <span v-if="!collapsed">{{ item.name }}</span>
    </Link>
  </nav>
</template>
