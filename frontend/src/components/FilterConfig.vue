<script setup lang="ts">
import { reactive, watch } from "vue";
import {
  NCard,
  NFormItem,
  NInputNumber,
  NSwitch,
  NDynamicInput,
  NText,
} from "naive-ui";
import type { FilterConfig } from "../api";

const emit = defineEmits<{
  "update:config": [config: FilterConfig];
}>();

const form = reactive<FilterConfig>({
  minLength: 0,
  maxLength: 0,
  minFrequency: 0,
  filterEnglish: false,
  filterNumber: false,
  customRules: [],
});

function reset() {
  form.minLength = 0;
  form.maxLength = 0;
  form.minFrequency = 0;
  form.filterEnglish = false;
  form.filterNumber = false;
  form.customRules = [];
}

defineExpose({ reset });

watch(
  form,
  () => {
    emit("update:config", { ...form });
  },
  { deep: true, immediate: true },
);
</script>

<template>
  <n-card title="过滤选项" size="small" style="margin-top: 12px">
    <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 8px 16px">
      <n-form-item label="最小词长" label-placement="left" label-width="80">
        <n-input-number v-model:value="form.minLength" :min="0" size="small" style="width: 100px" />
      </n-form-item>
      <n-form-item label="最大词长" label-placement="left" label-width="80">
        <n-input-number v-model:value="form.maxLength" :min="0" size="small" style="width: 100px" />
      </n-form-item>
      <n-form-item label="最小词频" label-placement="left" label-width="80">
        <n-input-number v-model:value="form.minFrequency" :min="0" size="small" style="width: 100px" />
      </n-form-item>
    </div>

    <div style="display: flex; gap: 24px; margin-top: 8px">
      <n-form-item label="包含英文的词条" label-placement="left" label-width="120" style="margin-bottom: 0">
        <n-switch v-model:value="form.filterEnglish" size="small" />
      </n-form-item>
      <n-form-item label="包含数字的词条" label-placement="left" label-width="120" style="margin-bottom: 0">
        <n-switch v-model:value="form.filterNumber" size="small" />
      </n-form-item>
    </div>

    <div style="margin-top: 12px">
      <n-text depth="3" style="font-size: 12px; margin-bottom: 4px; display: block">
        自定义正则过滤规则（匹配到的词条会被过滤）
      </n-text>
      <n-dynamic-input
        v-model:value="form.customRules"
        placeholder="输入正则表达式"
        :min="0"
      />
    </div>
  </n-card>
</template>
