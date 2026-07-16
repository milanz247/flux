<script setup lang="ts">
// SPA navigation link — the Flux equivalent of Inertia's <Link>. Renders an
// <a> but intercepts the click and routes through the Flux client, so pages
// swap without a full reload. Modified clicks (ctrl/cmd/shift, middle
// button) fall through to the browser.
import { router } from './router'
import type { VisitMethod } from './types'

const props = withDefaults(
  defineProps<{
    href: string
    method?: VisitMethod
    preserveScroll?: boolean
  }>(),
  { method: 'get', preserveScroll: false },
)

function onClick(event: MouseEvent) {
  if (event.defaultPrevented) return
  if (event.button !== 0 || event.metaKey || event.ctrlKey || event.shiftKey || event.altKey) return

  event.preventDefault()
  void router.visit(props.href, {
    method: props.method,
    preserveScroll: props.preserveScroll,
  })
}
</script>

<template>
  <a :href="href" @click="onClick">
    <slot />
  </a>
</template>
