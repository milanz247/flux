<script setup lang="ts">
import { DialogDescription } from 'reka-ui'
import { cn } from '@/lib/utils'
</script>

<template>
  <DialogDescription :class="cn('text-sm text-muted-foreground', $attrs.class as string)">
    <slot />
  </DialogDescription>
</template>
