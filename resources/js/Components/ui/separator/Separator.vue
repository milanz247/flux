<script setup lang="ts">
import { cn } from '@/lib/utils'

const props = withDefaults(
  defineProps<{
    orientation?: 'horizontal' | 'vertical'
    decorative?: boolean
    class?: string
  }>(),
  { orientation: 'horizontal', decorative: true },
)
</script>

<template>
  <div
    :role="decorative ? 'none' : 'separator'"
    :aria-orientation="decorative ? undefined : orientation"
    :class="
      cn(
        'shrink-0 bg-border',
        orientation === 'horizontal' ? 'h-[1px] w-full' : 'h-full w-[1px]',
        props.class,
      )
    "
  />
</template>
