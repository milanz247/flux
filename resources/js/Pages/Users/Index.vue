<script setup lang="ts">
// Props arrive automatically from req.View("Users/Index", dto.UserIndexDTO{...})
import { ref } from 'vue'
import { Plus, Search, Pencil, Trash2, Eye } from 'lucide-vue-next'
import DashboardLayout from '@/Layouts/DashboardLayout.vue'
import { Link, router } from '@/flux'
import { Button } from '@/Components/ui/button'
import { Input } from '@/Components/ui/input'
import { Badge } from '@/Components/ui/badge'
import { Card, CardContent } from '@/Components/ui/card'
import {
  Table, TableBody, TableCell, TableHead, TableHeader, TableRow,
} from '@/Components/ui/table'
import type { PaginationDTO, UserDTO } from '@/types'

const props = defineProps<{
  users: UserDTO[]
  pagination: PaginationDTO
  search: string
}>()

const searchTerm = ref(props.search)

function submitSearch() {
  void router.get('/users', { search: searchTerm.value }, { preserveScroll: true })
}

function goToPage(page: number) {
  void router.get('/users', { page, search: props.search }, { preserveScroll: true })
}

function destroy(user: UserDTO) {
  if (!confirm(`Delete ${user.name}? This cannot be undone.`)) return
  void router.delete(`/users/${user.id}`, { preserveScroll: true })
}
</script>

<template>
  <DashboardLayout title="Users">
    <div class="mb-4 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <form class="relative w-full max-w-sm" @submit.prevent="submitSearch">
        <Search class="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
        <Input
          v-model="searchTerm"
          class="pl-9"
          placeholder="Search by name or email…"
          @keyup.enter="submitSearch"
        />
      </form>
      <Link href="/users/create">
        <Button><Plus class="h-4 w-4" /> New User</Button>
      </Link>
    </div>

    <Card>
      <CardContent class="p-0">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead class="pl-4">Name</TableHead>
              <TableHead>Email</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Created</TableHead>
              <TableHead class="pr-4 text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-if="users.length === 0">
              <TableCell colspan="5" class="py-10 text-center text-muted-foreground">
                No users found.
              </TableCell>
            </TableRow>
            <TableRow v-for="user in users" :key="user.id">
              <TableCell class="pl-4 font-medium">{{ user.name }}</TableCell>
              <TableCell class="text-muted-foreground">{{ user.email }}</TableCell>
              <TableCell>
                <Badge :variant="user.verified ? 'success' : 'secondary'">
                  {{ user.verified ? 'Verified' : 'Pending' }}
                </Badge>
              </TableCell>
              <TableCell class="text-muted-foreground">{{ user.createdAt }}</TableCell>
              <TableCell class="pr-4">
                <div class="flex justify-end gap-1">
                  <Link :href="`/users/${user.id}`">
                    <Button variant="ghost" size="icon"><Eye class="h-4 w-4" /></Button>
                  </Link>
                  <Link :href="`/users/${user.id}/edit`">
                    <Button variant="ghost" size="icon"><Pencil class="h-4 w-4" /></Button>
                  </Link>
                  <Button variant="ghost" size="icon" class="text-destructive" @click="destroy(user)">
                    <Trash2 class="h-4 w-4" />
                  </Button>
                </div>
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </CardContent>
    </Card>

    <!-- Pagination -->
    <div
      v-if="pagination.totalPages > 1"
      class="mt-4 flex items-center justify-between text-sm text-muted-foreground"
    >
      <span>
        Page {{ pagination.page }} of {{ pagination.totalPages }} — {{ pagination.total }} users
      </span>
      <div class="flex gap-2">
        <Button
          variant="outline"
          size="sm"
          :disabled="pagination.page <= 1"
          @click="goToPage(pagination.page - 1)"
        >
          Previous
        </Button>
        <Button
          variant="outline"
          size="sm"
          :disabled="pagination.page >= pagination.totalPages"
          @click="goToPage(pagination.page + 1)"
        >
          Next
        </Button>
      </div>
    </div>
  </DashboardLayout>
</template>
