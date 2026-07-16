import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

/** shadcn's class combiner: conditional classes + Tailwind conflict merging. */
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}
