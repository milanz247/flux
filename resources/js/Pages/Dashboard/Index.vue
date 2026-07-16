<script setup lang="ts">
// Props arrive automatically from req.View("Dashboard/Index", dto.DashboardDTO{...})
import { Users as UsersIcon, ArrowRight } from 'lucide-vue-next'
import DashboardLayout from '@/Layouts/DashboardLayout.vue'
import { Link } from '@/flux'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/Components/ui/card'
import { Badge } from '@/Components/ui/badge'
import { Button } from '@/Components/ui/button'
import type { DashboardStatDTO, UserDTO } from '@/types'

defineProps<{
  stats: DashboardStatDTO[]
  recentUsers: UserDTO[]
}>()
</script>

<template>
  <DashboardLayout title="Dashboard">
    <!-- KPI tiles -->
    <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <Card v-for="stat in stats" :key="stat.label">
        <CardHeader class="pb-2">
          <CardDescription>{{ stat.label }}</CardDescription>
          <CardTitle class="text-3xl">{{ stat.value }}</CardTitle>
        </CardHeader>
        <CardContent>
          <p class="text-xs text-muted-foreground">{{ stat.hint }}</p>
        </CardContent>
      </Card>
    </div>

    <!-- Recent users -->
    <Card class="mt-6">
      <CardHeader class="flex-row items-center justify-between space-y-0">
        <div class="space-y-1.5">
          <CardTitle class="flex items-center gap-2">
            <UsersIcon class="h-4 w-4" /> Recent Users
          </CardTitle>
          <CardDescription>The latest accounts created in your application.</CardDescription>
        </div>
        <Link href="/users">
          <Button variant="outline" size="sm">
            View all <ArrowRight class="h-4 w-4" />
          </Button>
        </Link>
      </CardHeader>
      <CardContent>
        <div v-if="recentUsers.length === 0" class="py-8 text-center text-sm text-muted-foreground">
          No users yet — run <code class="rounded bg-muted px-1">flux seed</code> to add demo data.
        </div>
        <ul v-else class="divide-y">
          <li v-for="user in recentUsers" :key="user.id" class="flex items-center gap-4 py-3">
            <span
              class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-secondary text-sm font-semibold text-secondary-foreground"
            >
              {{ user.name.slice(0, 1).toUpperCase() }}
            </span>
            <div class="min-w-0 flex-1">
              <p class="truncate text-sm font-medium">{{ user.name }}</p>
              <p class="truncate text-xs text-muted-foreground">{{ user.email }}</p>
            </div>
            <Badge :variant="user.verified ? 'success' : 'secondary'">
              {{ user.verified ? 'Verified' : 'Pending' }}
            </Badge>
            <span class="hidden text-xs text-muted-foreground sm:block">{{ user.createdAt }}</span>
          </li>
        </ul>
      </CardContent>
    </Card>
  </DashboardLayout>
</template>
