<script setup lang="ts">
import { ref, provide } from 'vue'

const open = ref(false)

provide('dropdown-menu-open', open)

function toggle() {
  open.value = !open.value
}

function close() {
  open.value = false
}

provide('dropdown-menu-toggle', toggle)
provide('dropdown-menu-close', close)
</script>

<template>
  <div class="relative inline-block text-left">
    <slot />
  </div>
</template>
