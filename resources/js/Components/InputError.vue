<script setup lang="ts">
// Renders a validation message from useForm().errors under a field.
defineProps<{ message?: string }>()
</script>

<template>
  <p v-if="message" class="text-sm font-medium text-destructive">
    {{ message }}
  </p>
</template>
