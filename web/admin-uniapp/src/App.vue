<!-- Author: Charlie -->

<script setup lang="ts">
import { onLaunch } from '@dcloudio/uni-app'
import { useAuthStore } from '@/stores/auth'

onLaunch(() => {
  const authStore = useAuthStore()
  if (!authStore.isLogin) {
    setTimeout(() => {
      uni.reLaunch({ url: '/pages/auth/login/login' })
    }, 100)
  }
})
</script>
<style lang="scss">
@import 'uview-pro/index.scss';

page {
  background: var(--color-neutral-50);
  color: var(--color-neutral-900);
  font-size: var(--text-sm);
  font-weight: 400;
  line-height: 1.5;
}
</style>
