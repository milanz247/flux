<script setup lang="ts">
import { cn } from '@/lib/utils'
import { alertVariants, type AlertVariants } from './index'

defineProps<{ variant?: AlertVariants['variant'] }>()
</script>

<template>
  <div role="alert" :class="cn(alertVariants({ variant }), $attrs.class as string)">
    <slot />
  </div>
</template>
