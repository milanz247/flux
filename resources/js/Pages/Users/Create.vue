<script setup lang="ts">
// Props arrive automatically from req.View("Users/Create", dto.UserCreateDTO{})
import { ArrowLeft } from 'lucide-vue-next'
import DashboardLayout from '@/Layouts/DashboardLayout.vue'
import { Link, useForm } from '@/flux'
import { Button } from '@/Components/ui/button'
import { Input } from '@/Components/ui/input'
import { Label } from '@/Components/ui/label'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/Components/ui/card'
import InputError from '@/Components/InputError.vue'

const form = useForm({
  name: '',
  email: '',
  password: '',
})

function submit() {
  void form.post('/users')
}
</script>

<template>
  <DashboardLayout title="New User">
    <div class="mx-auto max-w-2xl">
      <Link href="/users" class="mb-4 inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground">
        <ArrowLeft class="h-4 w-4" /> Back to users
      </Link>

      <Card>
        <CardHeader>
          <CardTitle>Create User</CardTitle>
          <CardDescription>Add a new account. Validation runs server-side via the CreateUserDTO.</CardDescription>
        </CardHeader>
        <CardContent>
          <form class="space-y-5" @submit.prevent="submit">
            <div class="space-y-2">
              <Label for="name">Name</Label>
              <Input id="name" v-model="form.data.name" placeholder="Jane Doe" />
              <InputError :message="form.error('name')" />
            </div>

            <div class="space-y-2">
              <Label for="email">Email</Label>
              <Input id="email" v-model="form.data.email" type="email" placeholder="jane@example.com" />
              <InputError :message="form.error('email')" />
            </div>

            <div class="space-y-2">
              <Label for="password">Password</Label>
              <Input id="password" v-model="form.data.password" type="password" autocomplete="new-password" />
              <InputError :message="form.error('password')" />
            </div>

            <div class="flex justify-end gap-3">
              <Link href="/users"><Button variant="outline">Cancel</Button></Link>
              <Button type="submit" :disabled="form.processing">
                {{ form.processing ? 'Creating…' : 'Create User' }}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  </DashboardLayout>
</template>
