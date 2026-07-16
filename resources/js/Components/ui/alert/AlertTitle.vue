<script setup lang="ts">
import { cn } from '@/lib/utils'
</script>

<template>
  <h5 :class="cn('mb-1 font-medium leading-none tracking-tight', $attrs.class as string)">
    <slot />
  </h5>
</template>
