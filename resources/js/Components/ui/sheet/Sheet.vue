<script setup lang="ts">
import { ref, provide } from 'vue'

const open = ref(false)

provide('sheet-open', open)

function setOpen(value: boolean) {
  open.value = value
}

provide('sheet-set-open', setOpen)
</script>

<template>
  <slot />
</template>
