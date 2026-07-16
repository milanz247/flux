<script setup lang="ts">
import { cn } from '@/lib/utils'
</script>

<template>
  <td :class="cn('p-2 align-middle', $attrs.class as string)">
    <slot />
  </td>
</template>
