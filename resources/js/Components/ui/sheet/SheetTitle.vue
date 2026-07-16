<script setup lang="ts">
import { DialogTitle } from 'reka-ui'
import { cn } from '@/lib/utils'
</script>

<template>
  <DialogTitle :class="cn('text-lg font-semibold text-foreground', $attrs.class as string)">
    <slot />
  </DialogTitle>
</template>
