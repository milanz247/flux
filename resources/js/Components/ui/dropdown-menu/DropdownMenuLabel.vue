<script setup lang="ts">
import { DropdownMenuLabel } from 'reka-ui'
import { cn } from '@/lib/utils'
</script>

<template>
  <DropdownMenuLabel :class="cn('px-3 py-2 text-sm font-medium', $attrs.class as string)">
    <slot />
  </DropdownMenuLabel>
</template>
