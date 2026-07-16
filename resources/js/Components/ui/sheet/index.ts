import { cva, type VariantProps } from 'class-variance-authority'

export {
  DialogRoot as Sheet,
  DialogTrigger as SheetTrigger,
  DialogClose as SheetClose,
  DialogPortal as SheetPortal,
} from 'reka-ui'
export { default as SheetContent } from './SheetContent.vue'
export { default as SheetOverlay } from './SheetOverlay.vue'
export { default as SheetHeader } from './SheetHeader.vue'
export { default as SheetTitle } from './SheetTitle.vue'
export { default as SheetDescription } from './SheetDescription.vue'

export const sheetVariants = cva(
  'fixed z-50 gap-4 border bg-background p-6 shadow-lg transition ease-in-out',
  {
    variants: {
      side: {
        top: 'inset-x-0 top-0 border-b',
        bottom: 'inset-x-0 bottom-0 border-t',
        left: 'inset-y-0 left-0 h-full w-3/4 border-r sm:max-w-sm',
        right: 'inset-y-0 right-0 h-full w-3/4 border-l sm:max-w-sm',
      },
    },
    defaultVariants: { side: 'right' },
  },
)

export type SheetVariants = VariantProps<typeof sheetVariants>
