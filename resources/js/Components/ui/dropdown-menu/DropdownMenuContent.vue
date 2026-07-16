<script setup lang="ts">
import { inject, ref, type Ref, onMounted, onUnmounted } from 'vue'
import { cn } from '@/lib/utils'

const props = withDefaults(
  defineProps<{
    align?: 'start' | 'center' | 'end'
    class?: string
  }>(),
  { align: 'end' },
)

const open = inject<Ref<boolean>>('dropdown-menu-open', ref(false))
const close = inject<() => void>('dropdown-menu-close', () => {})

function onClickOutside(e: MouseEvent) {
  const target = e.target as HTMLElement
  if (!target.closest('[data-dropdown-content]') && !target.closest('[data-dropdown-trigger]')) {
    close()
  }
}

onMounted(() => {
  document.addEventListener('click', onClickOutside, { capture: true })
})

onUnmounted(() => {
  document.removeEventListener('click', onClickOutside, { capture: true })
})
</script>

<template>
  <Transition
    enter-active-class="transition ease-out duration-100"
    enter-from-class="transform opacity-0 scale-95"
    enter-to-class="transform opacity-100 scale-100"
    leave-active-class="transition ease-in duration-75"
    leave-from-class="transform opacity-100 scale-100"
    leave-to-class="transform opacity-0 scale-95"
  >
    <div
      v-if="open"
      data-dropdown-content
      :class="
        cn(
          'absolute z-50 mt-2 min-w-[8rem] overflow-hidden rounded-md border bg-popover p-1 text-popover-foreground shadow-md',
          align === 'end' && 'right-0',
          align === 'start' && 'left-0',
          align === 'center' && 'left-1/2 -translate-x-1/2',
          props.class,
        )
      "
    >
      <slot />
    </div>
  </Transition>
</template>
