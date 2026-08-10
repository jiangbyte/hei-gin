<!-- Author: Charlie -->

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{ password: string }>()

interface StrengthLevel {
  label: string
  color: string
  percent: number
}

const levels: StrengthLevel[] = [
  { label: '弱', color: '#e74c3c', percent: 25 },
  { label: '较弱', color: '#f39c12', percent: 50 },
  { label: '中等', color: '#f1c40f', percent: 75 },
  { label: '强', color: '#2ecc71', percent: 100 },
]

const strength = computed(() => {
  const pwd = props.password ?? ''
  if (!pwd) return { label: '', color: '#e0e0e0', percent: 0 }

  let score = 0
  if (pwd.length >= 8) score += 1
  if (pwd.length >= 12) score += 1
  if (/[a-z]/.test(pwd) && /[A-Z]/.test(pwd)) score += 1
  if (/[0-9]/.test(pwd) && /[^A-Za-z0-9]/.test(pwd)) score += 1

  return levels[Math.min(score, levels.length - 1)]
})

const passwordPolicy = computed(() => [
  { label: '至少 8 个字符', met: (props.password?.length ?? 0) >= 8 },
  { label: '包含大写字母', met: /[A-Z]/.test(props.password ?? '') },
  { label: '包含小写字母', met: /[a-z]/.test(props.password ?? '') },
  { label: '包含数字', met: /[0-9]/.test(props.password ?? '') },
  { label: '包含特殊字符', met: /[^A-Za-z0-9]/.test(props.password ?? '') },
])
</script>

<template>
  <div
    v-if="password"
    class="password-strength-bar"
  >
    <div class="bar-track">
      <div
        class="bar-fill"
        :style="{ width: strength.percent + '%', backgroundColor: strength.color }"
      />
    </div>
    <span
      v-if="strength.label"
      class="strength-label"
      :style="{ color: strength.color }"
    >
      {{ strength.label }}
    </span>
    <div class="policy-list">
      <div
        v-for="item in passwordPolicy"
        :key="item.label"
        class="policy-item"
      >
        <span :class="['policy-icon', item.met ? 'met' : '']">
          {{ item.met ? '✓' : '○' }}
        </span>
        <span :class="['policy-text', item.met ? 'met' : '']">{{ item.label }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.password-strength-bar {
  margin-top: 4px;
}

.bar-track {
  height: 4px;
  background-color: #e8e8e8;
  border-radius: 2px;
  overflow: hidden;
  margin-bottom: 4px;
}

.bar-fill {
  height: 100%;
  border-radius: 2px;
  transition:
    width 0.3s ease,
    background-color 0.3s ease;
}

.strength-label {
  font-size: 12px;
  font-weight: 600;
  margin-bottom: 4px;
  display: inline-block;
}

.policy-list {
  display: flex;
  flex-wrap: wrap;
  gap: 4px 12px;
  margin-top: 4px;
}

.policy-item {
  display: flex;
  align-items: center;
  gap: 3px;
  font-size: 11px;
  color: #999;
}

.policy-icon {
  font-size: 10px;
  color: #ccc;
  transition: color 0.2s;
}

.policy-icon.met {
  color: #2ecc71;
}

.policy-text.met {
  color: #333;
}
</style>
