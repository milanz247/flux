<script setup lang="ts">
// Props arrive automatically from req.View("Users/Show", dto.UserShowDTO{...})
import { ArrowLeft, Pencil } from 'lucide-vue-next'
import DashboardLayout from '@/Layouts/DashboardLayout.vue'
import { Link } from '@/flux'
import { Button } from '@/Components/ui/button'
import { Badge } from '@/Components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle } from '@/Components/ui/card'
import type { UserDTO } from '@/types'

defineProps<{ user: UserDTO }>()
</script>

<template>
  <DashboardLayout title="User Details">
    <div class="mx-auto max-w-2xl">
      <Link href="/users" class="mb-4 inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground">
        <ArrowLeft class="h-4 w-4" /> Back to users
      </Link>

      <Card>
        <CardHeader class="flex-row items-center gap-4 space-y-0">
          <span
            class="flex h-14 w-14 items-center justify-center rounded-full bg-primary text-xl font-semibold text-primary-foreground"
          >
            {{ user.name.slice(0, 1).toUpperCase() }}
          </span>
          <div class="flex-1">
            <CardTitle class="text-xl">{{ user.name }}</CardTitle>
            <p class="text-sm text-muted-foreground">{{ user.email }}</p>
          </div>
          <Link :href="`/users/${user.id}/edit`">
            <Button variant="outline"><Pencil class="h-4 w-4" /> Edit</Button>
          </Link>
        </CardHeader>
        <CardContent>
          <dl class="grid gap-4 sm:grid-cols-3">
            <div>
              <dt class="text-xs uppercase tracking-wide text-muted-foreground">ID</dt>
              <dd class="mt-1 text-sm font-medium">#{{ user.id }}</dd>
            </div>
            <div>
              <dt class="text-xs uppercase tracking-wide text-muted-foreground">Status</dt>
              <dd class="mt-1">
                <Badge :variant="user.verified ? 'success' : 'secondary'">
                  {{ user.verified ? 'Verified' : 'Pending verification' }}
                </Badge>
              </dd>
            </div>
            <div>
              <dt class="text-xs uppercase tracking-wide text-muted-foreground">Created</dt>
              <dd class="mt-1 text-sm font-medium">{{ user.createdAt }}</dd>
            </div>
          </dl>
        </CardContent>
      </Card>
    </div>
  </DashboardLayout>
</template>
