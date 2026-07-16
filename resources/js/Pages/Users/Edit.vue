<script setup lang="ts">
// Props arrive automatically from req.View("Users/Edit", dto.UserEditDTO{...})
import { ArrowLeft } from 'lucide-vue-next'
import DashboardLayout from '@/Layouts/DashboardLayout.vue'
import { Link, useForm } from '@/flux'
import { Button } from '@/Components/ui/button'
import { Input } from '@/Components/ui/input'
import { Label } from '@/Components/ui/label'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/Components/ui/card'
import InputError from '@/Components/InputError.vue'
import type { UserDTO } from '@/types'

const props = defineProps<{ user: UserDTO }>()

const form = useForm({
  name: props.user.name,
  email: props.user.email,
  password: '',
})

function submit() {
  void form.put(`/users/${props.user.id}`)
}
</script>

<template>
  <DashboardLayout title="Edit User">
    <div class="mx-auto max-w-2xl">
      <Link href="/users" class="mb-4 inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground">
        <ArrowLeft class="h-4 w-4" /> Back to users
      </Link>

      <Card>
        <CardHeader>
          <CardTitle>Edit {{ user.name }}</CardTitle>
          <CardDescription>Leave the password blank to keep the current one.</CardDescription>
        </CardHeader>
        <CardContent>
          <form class="space-y-5" @submit.prevent="submit">
            <div class="space-y-2">
              <Label for="name">Name</Label>
              <Input id="name" v-model="form.data.name" />
              <InputError :message="form.error('name')" />
            </div>

            <div class="space-y-2">
              <Label for="email">Email</Label>
              <Input id="email" v-model="form.data.email" type="email" />
              <InputError :message="form.error('email')" />
            </div>

            <div class="space-y-2">
              <Label for="password">New password (optional)</Label>
              <Input id="password" v-model="form.data.password" type="password" autocomplete="new-password" />
              <InputError :message="form.error('password')" />
            </div>

            <div class="flex justify-end gap-3">
              <Link href="/users"><Button variant="outline">Cancel</Button></Link>
              <Button type="submit" :disabled="form.processing">
                {{ form.processing ? 'Saving…' : 'Save Changes' }}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  </DashboardLayout>
</template>
