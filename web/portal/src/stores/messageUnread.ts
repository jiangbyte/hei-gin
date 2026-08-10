/** Author: Charlie */

import { create } from 'zustand'
import { messageApi } from '@/api'
import { wireInt } from '@/utils/wire'

type MessageUnreadState = {
  unreadTotal: number
  refresh: () => Promise<void>
  setUnreadTotal: (total: number) => void
  notifyRead: (count?: number) => void
  notifyReadAll: () => void
}

export const useMessageUnreadStore = create<MessageUnreadState>((set, get) => ({
  unreadTotal: 0,
  setUnreadTotal: (total) => set({ unreadTotal: Math.max(0, total) }),
  notifyRead: (count = 1) =>
    set({ unreadTotal: Math.max(0, get().unreadTotal - Math.max(1, count)) }),
  notifyReadAll: () => set({ unreadTotal: 0 }),
  refresh: async () => {
    try {
      const res = await messageApi.unreadCount()
      const raw = res.data
      const total =
        typeof raw === 'string' ? wireInt(raw) : typeof raw === 'number' ? raw : Number(raw ?? 0)
      set({ unreadTotal: Number.isFinite(total) ? Math.max(0, total) : 0 })
    } catch {
      /* ignore */
    }
  },
}))
