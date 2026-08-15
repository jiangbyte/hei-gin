/** Author: Charlie */

import { defineStore } from 'pinia'
import { myNoticeApi } from '@/api'
import { wireInt } from '@/utils/wire'

export const useMessageUnreadStore = defineStore('message-unread-store', {
  state: () => ({
    unreadTotal: 0,
  }),
  actions: {
    setUnreadTotal(total: number) {
      this.unreadTotal = Math.max(0, total)
    },
    notifyRead(count = 1) {
      this.unreadTotal = Math.max(0, this.unreadTotal - Math.max(1, count))
    },
    notifyReadAll() {
      this.unreadTotal = 0
    },
    async refresh() {
      try {
        const res = await myNoticeApi.unreadCount()
        const raw = res.data
        const total =
          typeof raw === 'string' ? wireInt(raw) : typeof raw === 'number' ? raw : Number(raw ?? 0)
        this.unreadTotal = Number.isFinite(total) ? Math.max(0, total) : 0
      } catch {
        /* 忽略 */
      }
    },
  },
})
