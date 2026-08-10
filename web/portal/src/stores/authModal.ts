/** Author: Charlie */

import { create } from 'zustand'

export type AuthModalMode = 'login' | 'register'

type AuthModalState = {
  mode: AuthModalMode | null
  redirect?: string
  open: (mode: AuthModalMode, redirect?: string | null) => void
  switchMode: (mode: AuthModalMode) => void
  close: () => void
}

export const useAuthModalStore = create<AuthModalState>((set) => ({
  mode: null,
  redirect: undefined,
  open: (mode, redirect) =>
    set({
      mode,
      redirect: redirect?.trim() || undefined,
    }),
  switchMode: (mode) => set({ mode }),
  close: () => set({ mode: null, redirect: undefined }),
}))
