<script setup lang="ts">
import { computed } from 'vue'
import { NSelect, NTag } from 'naive-ui'
import type { FormatInfo } from '../api'
import { formatKindLabel } from '../api'

const props = defineProps<{
  formats: FormatInfo[]
  modelValue: string
  label: string
  mode: 'import' | 'export'
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
  'customSelect': []
}>()

const CUSTOM_ID = '__custom__'

const options = computed(() => {
  const list = props.formats
    .filter(f => props.mode === 'import' ? f.canImport : f.canExport)
    .map(f => ({
      label: `${f.name} (${f.id})`,
      value: f.id,
      kind: f.kind,
    }))
  list.push({
    label: '自定义格式...',
    value: CUSTOM_ID,
    kind: 0,
  })
  return list
})

function handleUpdate(value: string) {
  if (value === CUSTOM_ID) {
    emit('customSelect')
    return
  }
  emit('update:modelValue', value)
}

const selectedFormat = computed(() =>
  props.formats.find(f => f.id === props.modelValue)
)
</script>

<template>
  <div>
    <div style="margin-bottom: 6px; font-weight: 500">{{ label }}</div>
    <n-select
      :value="modelValue"
      :options="options"
      filterable
      placeholder="请选择格式"
      @update:value="handleUpdate"
      :render-label="({ label: text }: any) => text"
    />
    <div v-if="selectedFormat" style="margin-top: 4px">
      <n-tag size="small" :type="selectedFormat.kind === 1 ? 'info' : selectedFormat.kind === 2 ? 'warning' : 'default'">
        {{ formatKindLabel(selectedFormat.kind) }}
      </n-tag>
    </div>
  </div>
</template>
