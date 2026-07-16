<script setup lang="ts">
import { AvatarFallback } from 'reka-ui'
import { cn } from '@/lib/utils'
</script>

<template>
  <AvatarFallback
    :class="
      cn(
        'flex h-full w-full items-center justify-center rounded-full bg-primary text-xs font-semibold text-primary-foreground',
        $attrs.class as string,
      )
    "
  >
    <slot />
  </AvatarFallback>
</template>
