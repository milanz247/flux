<script setup lang="ts">
import { DropdownMenuContent, DropdownMenuPortal, type DropdownMenuContentProps } from 'reka-ui'
import { cn } from '@/lib/utils'

withDefaults(
  defineProps<DropdownMenuContentProps>(),
  { sideOffset: 4, align: 'end' },
)
</script>

<template>
  <DropdownMenuPortal>
    <DropdownMenuContent
      v-bind="$props"
      :class="
        cn(
          'z-50 min-w-[10rem] overflow-hidden rounded-md border bg-popover p-1 text-popover-foreground shadow-md',
          $attrs.class as string,
        )
      "
    >
      <slot />
    </DropdownMenuContent>
  </DropdownMenuPortal>
</template>
