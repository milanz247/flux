<script setup lang="ts">
// Props arrive automatically from req.View("Dashboard/Index", dto.DashboardDTO{})
// The authenticated user is available via shared props (useAuth).
import { Zap } from 'lucide-vue-next'
import DashboardLayout from '@/Layouts/DashboardLayout.vue'
import { useAuth } from '@/flux'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/Components/ui/card'

const user = useAuth()
</script>

<template>
  <DashboardLayout title="Dashboard">
    <div class="mx-auto max-w-4xl">
      <!-- Welcome card -->
      <Card class="overflow-hidden">
        <div class="relative">
          <!-- Subtle gradient top border -->
          <div class="absolute inset-x-0 top-0 h-1 bg-gradient-to-r from-primary/60 via-primary to-primary/60" />
          <CardHeader class="pt-8">
            <div class="flex items-center gap-4">
              <div class="flex h-12 w-12 items-center justify-center rounded-xl bg-primary/10">
                <Zap class="h-6 w-6 text-primary" />
              </div>
              <div>
                <CardTitle class="text-xl">
                  Welcome back, {{ user?.name?.split(' ')[0] }}!
                </CardTitle>
                <CardDescription>
                  You're logged in to your Flux dashboard.
                </CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent>
            <p class="text-sm text-muted-foreground">
              This is your application dashboard. Flux gives you a beautiful starting point for
              building full-stack Go applications with Vue and shadcn-vue — everything you need
              to ship fast.
            </p>
          </CardContent>
        </div>
      </Card>
    </div>
  </DashboardLayout>
</template>
