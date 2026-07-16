<script setup lang="ts">
import { inject, type Ref } from 'vue'

const toggle = inject<() => void>('dropdown-menu-toggle', () => {})
</script>

<template>
  <div @click="toggle">
    <slot />
  </div>
</template>
