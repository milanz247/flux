<script setup lang="ts">
import { Separator } from 'reka-ui'
import { cn } from '@/lib/utils'

withDefaults(
  defineProps<{ orientation?: 'horizontal' | 'vertical'; decorative?: boolean }>(),
  { orientation: 'horizontal', decorative: true },
)
</script>

<template>
  <Separator
    :orientation="orientation"
    :decorative="decorative"
    :class="
      cn(
        'shrink-0 bg-border',
        orientation === 'horizontal' ? 'h-px w-full' : 'h-full w-px',
        $attrs.class as string,
      )
    "
  />
</template>
