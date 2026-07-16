<script setup lang="ts">
import { X } from 'lucide-vue-next'
import { DialogClose, DialogContent, type DialogContentEmits, type DialogContentProps } from 'reka-ui'
import { cn } from '@/lib/utils'
import SheetOverlay from './SheetOverlay.vue'
import { sheetVariants, type SheetVariants } from './index'

defineProps<DialogContentProps & { side?: SheetVariants['side'] }>()
defineEmits<DialogContentEmits>()
</script>

<template>
  <SheetOverlay />
  <DialogContent
    v-bind="$props"
    :class="cn(sheetVariants({ side }), $attrs.class as string)"
  >
    <slot />
    <DialogClose
      class="absolute right-4 top-4 rounded-sm opacity-70 ring-offset-background transition-opacity hover:opacity-100 focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
    >
      <X class="h-4 w-4" />
      <span class="sr-only">Close</span>
    </DialogClose>
  </DialogContent>
</template>
