import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function formatDaysUntil(daysUntil: number): string {
  if (daysUntil === 0) return 'сегодня'
  if (daysUntil === 1) return 'завтра'
  return `через ${daysUntil} дней`
}
