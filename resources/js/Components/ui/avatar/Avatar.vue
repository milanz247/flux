<script setup lang="ts">
import { AvatarRoot } from 'reka-ui'
import { cn } from '@/lib/utils'
</script>

<template>
  <AvatarRoot
    :class="cn('relative flex h-10 w-10 shrink-0 overflow-hidden rounded-full', $attrs.class as string)"
  >
    <slot />
  </AvatarRoot>
</template>
