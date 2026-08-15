<!-- Author: Charlie -->

<template>
  <view class="crop-page">
    <view class="crop-header">
      <view class="crop-action" @click="goBack">
        <u-icon name="close" color="#fff" size="20" />
      </view>
      <text class="crop-title">裁剪头像</text>
      <view style="width: 44px" />
    </view>

    <view class="crop-body">
      <view
        v-if="s.ok"
        class="crop-box"
        :style="{ width: s.iw + 'px', height: s.ih + 'px' }"
      >
        <image class="crop-pic" :src="s.path" mode="aspectFit" />
        <view
          class="crop-sel"
          :style="selS"
          @touchstart.stop="onStart"
          @touchmove.stop="onMove"
          @touchend.stop="onEnd"
        >
          <view class="crop-grid-h" /><view class="crop-grid-v" />
          <view
            class="crop-dot tl"
            @touchstart.stop="onDStart($event, 'tl')"
            @touchmove.stop="onDMove"
          />
          <view
            class="crop-dot tr"
            @touchstart.stop="onDStart($event, 'tr')"
            @touchmove.stop="onDMove"
          />
          <view
            class="crop-dot bl"
            @touchstart.stop="onDStart($event, 'bl')"
            @touchmove.stop="onDMove"
          />
          <view
            class="crop-dot br"
            @touchstart.stop="onDStart($event, 'br')"
            @touchmove.stop="onDMove"
          />
        </view>
      </view>
    </view>

    <canvas
      canvas-id="cropCanvas"
      class="hidden-canvas"
      :style="{ width: s.iw + 'px', height: s.ih + 'px' }"
      :width="s.iw"
      :height="s.ih"
    />

    <view class="crop-foot">
      <view class="crop-btn" @click="goBack">取消</view>
      <view class="crop-btn prim" @click="confirm">确认</view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed, reactive } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { authApi } from '@/api'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const SW = uni.getSystemInfoSync().windowWidth
const WH = uni.getSystemInfoSync().windowHeight
const MAX_H = WH - 88

const s = reactive({
  path: '',
  ok: false,
  ow: 0,
  oh: 0, // 原图宽高
  iw: 0,
  ih: 0, // 容器 = 渲染图宽高 (aspectFit)
  sx: 0,
  sy: 0,
  sz: 0,
  loading: false,
})

const selS = computed(() => ({
  left: s.sx + 'px',
  top: s.sy + 'px',
  width: s.sz + 'px',
  height: s.sz + 'px',
}))

onLoad((opts: any) => {
  s.path = decodeURIComponent(opts.imagePath || '')
  uni.getImageInfo({
    src: s.path,
    success: (info) => {
      s.ow = info.width
      s.oh = info.height

      // aspectFit: 完整显示图片
      let iw = SW
      let ih = iw * (info.height / info.width)
      if (ih > MAX_H) {
        ih = MAX_H
        iw = ih * (info.width / info.height)
      }
      s.iw = Math.round(iw)
      s.ih = Math.round(ih)

      // 初始选区: 正方形，取小边的 70%，居中
      const sz = Math.round(Math.min(iw, ih) * 0.7)
      s.sx = Math.round((iw - sz) / 2)
      s.sy = Math.round((ih - sz) / 2)
      s.sz = sz
      s.ok = true
    },
    fail: () => uni.showToast({ title: '图片加载失败', icon: 'none' }),
  })
})

// --- 拖动 ---
let bx = 0,
  by = 0,
  px = 0,
  py = 0
function onStart(e: any) {
  px = e.touches[0].pageX
  py = e.touches[0].pageY
  bx = s.sx
  by = s.sy
}
function onMove(e: any) {
  let nx = bx + e.touches[0].pageX - px
  let ny = by + e.touches[0].pageY - py
  if (nx < 0) nx = 0
  if (ny < 0) ny = 0
  if (nx + s.sz > s.iw) nx = s.iw - s.sz
  if (ny + s.sz > s.ih) ny = s.ih - s.sz
  s.sx = Math.round(nx)
  s.sy = Math.round(ny)
}
function onEnd() {}

// --- 缩放 ---
let dt = '',
  dX = 0,
  dY = 0,
  dBx = 0,
  dBy = 0,
  dSz = 0,
  MIN = 40
function onDStart(e: any, t: string) {
  dt = t
  dX = e.touches[0].pageX
  dY = e.touches[0].pageY
  dBx = s.sx
  dBy = s.sy
  dSz = s.sz
}
function onDMove(e: any) {
  const dx = e.touches[0].pageX - dX,
    dy = e.touches[0].pageY - dY
  let sz = dSz,
    x = dBx,
    y = dBy
  if (dt === 'br') {
    sz = dSz + Math.max(dx, dy)
    if (sz < MIN) sz = MIN
    if (x + sz > s.iw) sz = s.iw - x
    if (y + sz > s.ih) sz = s.ih - y
  } else if (dt === 'tr') {
    sz = dSz - dy
    if (sz < MIN) sz = MIN
    y = dBy + dSz - sz
    if (y < 0) {
      y = 0
      sz = dBy + dSz
    }
    if (x + sz > s.iw) sz = s.iw - x
  } else if (dt === 'bl') {
    sz = dSz - dx
    if (sz < MIN) sz = MIN
    x = dBx + dSz - sz
    if (x < 0) {
      x = 0
      sz = dBx + dSz
    }
    if (y + sz > s.ih) sz = s.ih - y
  } else {
    const d = Math.max(dx, dy)
    sz = dSz - d
    if (sz < MIN) sz = MIN
    x = dBx + dSz - sz
    y = dBy + dSz - sz
    if (x < 0 || y < 0) {
      sz = Math.min(dBx + dSz, dBy + dSz)
      x = 0
      y = 0
    }
  }
  s.sx = Math.round(x)
  s.sy = Math.round(y)
  s.sz = Math.round(sz)
}

// --- 确认裁剪 ---
function confirm() {
  if (s.loading) return
  s.loading = true
  uni.showLoading({ title: '裁剪中...' })
  // canvas 缓冲 = 图片渲染尺寸，坐标 1:1 对应
  const ctx = uni.createCanvasContext('cropCanvas')
  ctx.drawImage(s.path, 0, 0, s.iw, s.ih)
  ctx.draw(true, () => {
    uni.canvasToTempFilePath({
      x: s.sx,
      y: s.sy,
      width: s.sz,
      height: s.sz,
      destWidth: 640,
      destHeight: 640,
      canvasId: 'cropCanvas',
      success: (res) => {
        uni.hideLoading()
        authApi.uploadAvatar(res.tempFilePath).then(async () => {
          await authStore.refreshUserInfo()
          uni.showToast({ title: '成功', icon: 'none' })
          setTimeout(() => goBack(), 1000)
        })
      },
      fail: () => {
        uni.hideLoading()
        uni.showToast({ title: '裁剪失败', icon: 'none' })
        s.loading = false
      },
    })
  })
}

function goBack() {
  uni.navigateBack()
}
</script>

<style lang="scss" scoped>
.crop-page {
  position: fixed;
  inset: 0;
  background: #000;
  z-index: 1000;
  display: flex;
  flex-direction: column;
}
.crop-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  height: 44px;
  color: #fff;
  flex-shrink: 0;
}
.crop-title {
  font-size: 16px;
  font-weight: 600;
}
.crop-action {
  padding: 8px;
  width: 44px;
}
.crop-body {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}
.crop-box {
  position: relative;
}
.crop-pic {
  display: block;
  width: 100%;
  height: 100%;
}
.crop-sel {
  position: absolute;
  top: 0;
  left: 0;
  z-index: 3;
  box-shadow: 0 0 0 9999px rgba(0, 0, 0, 0.5);
  border: 2px solid rgba(255, 255, 255, 0.9);
}
.crop-grid-h {
  position: absolute;
  top: 33.33%;
  left: 0;
  width: 100%;
  height: 33.33%;
  border-top: 1px dashed rgba(255, 255, 255, 0.4);
  border-bottom: 1px dashed rgba(255, 255, 255, 0.4);
}
.crop-grid-v {
  position: absolute;
  left: 33.33%;
  top: 0;
  width: 33.33%;
  height: 100%;
  border-left: 1px dashed rgba(255, 255, 255, 0.4);
  border-right: 1px dashed rgba(255, 255, 255, 0.4);
}
.crop-dot {
  position: absolute;
  width: 24px;
  height: 24px;
  z-index: 4;
}
.crop-dot.tl {
  top: -13px;
  left: -13px;
  border-top: 3px solid #fff;
  border-left: 3px solid #fff;
}
.crop-dot.tr {
  top: -13px;
  right: -13px;
  border-top: 3px solid #fff;
  border-right: 3px solid #fff;
}
.crop-dot.bl {
  bottom: -13px;
  left: -13px;
  border-bottom: 3px solid #fff;
  border-left: 3px solid #fff;
}
.crop-dot.br {
  bottom: -13px;
  right: -13px;
  border-bottom: 3px solid #fff;
  border-right: 3px solid #fff;
}
.hidden-canvas {
  position: fixed;
  top: -9999px;
  left: -9999px;
  width: 375px;
  height: 375px;
}
.crop-foot {
  display: flex;
  padding: 0 16px 28px;
  gap: 12px;
  flex-shrink: 0;
}
.crop-btn {
  flex: 1;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 22px;
  font-size: 15px;
  font-weight: 600;
  background: rgba(255, 255, 255, 0.15);
  color: #fff;
}
.crop-btn.prim {
  background: #3c96f3;
}
</style>
