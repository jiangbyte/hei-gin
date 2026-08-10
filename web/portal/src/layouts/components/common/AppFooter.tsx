/** Author: Charlie */

export function AppFooter() {
  return (
    <footer className="muted-text py-6 text-center text-sm">
      {import.meta.env.VITE_APP_TITLE || 'HEI'} Portal
    </footer>
  )
}
