<!-- Author: Charlie -->

<template>
  <n-el
    tag="div"
    class="common-wrapper"
  >
    <n-flex
      align="center"
      :wrap="false"
      class="h-full"
    >
      <slot />
    </n-flex>
  </n-el>
</template>

<style scoped>
.common-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 32px;
  height: 32px;
  color: var(--text-color-2);
  cursor: pointer;
  border-radius: var(--border-radius);
  transition:
    background-color 0.2s var(--cubic-bezier-ease-in-out),
    color 0.2s var(--cubic-bezier-ease-in-out);
}

.common-wrapper:hover {
  color: var(--primary-color);
  background-color: var(--button-color-2-hover);
}
</style>
