<script setup lang="ts">
import { AvatarImage } from 'reka-ui'

defineProps<{ src: string; alt?: string }>()
</script>

<template>
  <AvatarImage :src="src" :alt="alt" class="aspect-square h-full w-full" />
</template>
