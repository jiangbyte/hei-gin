<!-- Author: Charlie -->

<script setup lang="ts">
import { NCheckbox, NCheckboxGroup, NRadio } from 'naive-ui'
import { computed } from 'vue'
import { dictList } from '@/utils/dict'

const modelValue = defineModel<any>({
  default: null,
})

const props = withDefaults(
  defineProps<{
    dictCode: any
    type?: any
    placeholder?: any
    clearable?: any
    disabled?: any
    size?: any
  }>(),
  {
    type: 'select',
    placeholder: '',
    clearable: true,
    disabled: false,
    size: 'medium',
  },
)

const emit = defineEmits<{
  change: [value: any]
}>()

const options = computed(() => dictList(props.dictCode))

function handleUpdateValue(value: any) {
  modelValue.value = value
  emit('change', modelValue.value)
}
</script>

<template>
  <NSelect
    v-if="type === 'select'"
    :value="modelValue"
    :options="options"
    :placeholder="placeholder"
    :clearable="clearable"
    :disabled="disabled"
    :size="size"
    @update:value="handleUpdateValue"
  />

  <NRadioGroup
    v-else-if="type === 'radio'"
    class="flex flex-wrap gap-x-16px gap-y-8px"
    :value="modelValue"
    :disabled="disabled"
    :size="size"
    @update:value="handleUpdateValue"
  >
    <NRadio
      v-for="option in options"
      :key="option.value"
      :value="option.value"
    >
      {{ option.label }}
    </NRadio>
  </NRadioGroup>

  <NCheckboxGroup
    v-else
    class="flex flex-wrap gap-x-16px gap-y-8px"
    :value="Array.isArray(modelValue) ? modelValue : []"
    :disabled="disabled"
    @update:value="handleUpdateValue"
  >
    <NCheckbox
      v-for="option in options"
      :key="option.value"
      :value="option.value"
    >
      {{ option.label }}
    </NCheckbox>
  </NCheckboxGroup>
</template>
