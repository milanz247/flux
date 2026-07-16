<script setup lang="ts">
import { DropdownMenuSeparator } from 'reka-ui'
import { cn } from '@/lib/utils'
</script>

<template>
  <DropdownMenuSeparator :class="cn('-mx-1 my-1 h-px bg-muted', $attrs.class as string)" />
</template>
