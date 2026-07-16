<script setup lang="ts">
import { DropdownMenuItem, type DropdownMenuItemProps } from 'reka-ui'
import { cn } from '@/lib/utils'

defineProps<DropdownMenuItemProps>()
</script>

<template>
  <DropdownMenuItem
    v-bind="$props"
    :class="
      cn(
        'relative flex cursor-pointer select-none items-center gap-2 rounded-sm px-3 py-2 text-sm outline-none transition-colors hover:bg-accent hover:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50',
        $attrs.class as string,
      )
    "
  >
    <slot />
  </DropdownMenuItem>
</template>
