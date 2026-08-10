<!-- Author: Charlie -->

<script setup lang="ts"></script>

<template>
  <naive-provider>
    <div id="loading">
      <n-spin size="large">
        加载中...
      </n-spin>
    </div>
  </naive-provider>
</template>

<style scoped>
#loading {
  width: 100vw;
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  gap: 15vh;
  position: fixed;
  z-index: 1;
}
</style>
