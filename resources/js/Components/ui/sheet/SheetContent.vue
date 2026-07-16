<script setup lang="ts">
import { inject, ref, type Ref } from 'vue'
import { cn } from '@/lib/utils'

const props = withDefaults(
  defineProps<{
    side?: 'left' | 'right' | 'top' | 'bottom'
    class?: string
  }>(),
  { side: 'left' },
)

const open = inject<Ref<boolean>>('sheet-open', ref(false))
const setOpen = inject<(value: boolean) => void>('sheet-set-open', () => {})

const sideClasses: Record<string, string> = {
  left: 'inset-y-0 left-0 h-full w-3/4 sm:max-w-sm border-r',
  right: 'inset-y-0 right-0 h-full w-3/4 sm:max-w-sm border-l',
  top: 'inset-x-0 top-0 border-b',
  bottom: 'inset-x-0 bottom-0 border-t',
}

const slideClasses: Record<string, { from: string; to: string }> = {
  left: { from: '-translate-x-full', to: 'translate-x-0' },
  right: { from: 'translate-x-full', to: 'translate-x-0' },
  top: { from: '-translate-y-full', to: 'translate-y-0' },
  bottom: { from: 'translate-y-full', to: 'translate-y-0' },
}
</script>

<template>
  <Teleport to="body">
    <!-- Backdrop -->
    <Transition
      enter-active-class="transition-opacity duration-300"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition-opacity duration-300"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="open"
        class="fixed inset-0 z-50 bg-black/80"
        @click="setOpen(false)"
      />
    </Transition>

    <!-- Panel -->
    <Transition
      enter-active-class="transition-transform duration-300 ease-out"
      :enter-from-class="slideClasses[side].from"
      :enter-to-class="slideClasses[side].to"
      leave-active-class="transition-transform duration-300 ease-in"
      :leave-from-class="slideClasses[side].to"
      :leave-to-class="slideClasses[side].from"
    >
      <div
        v-if="open"
        :class="
          cn(
            'fixed z-50 gap-4 bg-background p-6 shadow-lg',
            sideClasses[side],
            props.class,
          )
        "
      >
        <slot />
      </div>
    </Transition>
  </Teleport>
</template>
