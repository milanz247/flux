<script setup lang="ts">
import { inject } from 'vue'
import { X } from 'lucide-vue-next'

const setOpen = inject<(value: boolean) => void>('sheet-set-open', () => {})
</script>

<template>
  <button
    class="absolute right-4 top-4 rounded-sm opacity-70 ring-offset-background transition-opacity hover:opacity-100 focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
    @click="setOpen(false)"
  >
    <X class="h-4 w-4" />
    <span class="sr-only">Close</span>
  </button>
</template>
