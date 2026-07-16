<script setup lang="ts">
import { inject } from 'vue'
import { cn } from '@/lib/utils'

const props = defineProps<{ class?: string }>()
const close = inject<() => void>('dropdown-menu-close', () => {})
</script>

<template>
  <button
    :class="
      cn(
        'relative flex w-full cursor-default select-none items-center gap-2 rounded-sm px-2 py-1.5 text-sm outline-none transition-colors hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50',
        props.class,
      )
    "
    @click="close()"
  >
    <slot />
  </button>
</template>
