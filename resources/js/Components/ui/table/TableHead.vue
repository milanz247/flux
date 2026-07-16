<script setup lang="ts">
import { cn } from '@/lib/utils'
</script>

<template>
  <th
    :class="
      cn(
        'h-10 px-2 text-left align-middle font-medium text-muted-foreground',
        $attrs.class as string,
      )
    "
  >
    <slot />
  </th>
</template>
