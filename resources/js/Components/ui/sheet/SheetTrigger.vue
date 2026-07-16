<script setup lang="ts">
import { inject, ref, type Ref } from 'vue'

const open = inject<Ref<boolean>>('sheet-open', ref(false))
const setOpen = inject<(value: boolean) => void>('sheet-set-open', () => {})
</script>

<template>
  <div @click="setOpen(!open)">
    <slot />
  </div>
</template>
