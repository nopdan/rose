<script setup lang="ts">
import { ref, reactive } from "vue";
import {
  NModal,
  NForm,
  NFormItem,
  NInput,
  NSelect,
  NSwitch,
  NButton,
  NSpace,
  NText,
} from "naive-ui";
import type { CustomFormatConfig, CustomFieldConfig } from "../api";

const props = defineProps<{
  show: boolean;
}>();

const emit = defineEmits<{
  "update:show": [value: boolean];
  confirmed: [config: CustomFormatConfig];
}>();

const form = reactive({
  kind: "pinyin",
  encoding: "UTF-8",
  fields: [
    {
      type: "word",
      pinyinSeparator: "",
      pinyinPrefix: "",
      pinyinSuffix: "",
      literal: "",
    },
    {
      type: "tab",
      pinyinSeparator: "",
      pinyinPrefix: "",
      pinyinSuffix: "",
      literal: "",
    },
    {
      type: "pinyin",
      pinyinSeparator: "'",
      pinyinPrefix: "",
      pinyinSuffix: "",
      literal: "",
    },
  ] as CustomFieldConfig[],
  sortByCode: false,
  commentPrefix: "#",
  startMarker: "",
});

const errorMsg = ref("");

const kindOptions = [
  { label: "拼音", value: "pinyin" },
  { label: "五笔", value: "wubi" },
];

const encodingOptions = [
  { label: "UTF-8", value: "UTF-8" },
  { label: "UTF-16LE", value: "UTF-16LE" },
  { label: "GB18030", value: "GB18030" },
  { label: "GBK", value: "GBK" },
];

const fieldTypeOptions = [
  { label: "词条 (word)", value: "word" },
  { label: "拼音 (pinyin)", value: "pinyin" },
  { label: "编码 (code)", value: "code" },
  { label: "词频 (frequency)", value: "frequency" },
  { label: "排序 (rank)", value: "rank" },
  { label: "制表符 (tab)", value: "tab" },
  { label: "空格 (space)", value: "space" },
  { label: "字面量 (literal)", value: "literal" },
];

function addField() {
  form.fields.push({
    type: "tab",
    pinyinSeparator: "",
    pinyinPrefix: "",
    pinyinSuffix: "",
    literal: "",
  });
}

function removeField(index: number) {
  form.fields.splice(index, 1);
}

function moveField(index: number, direction: -1 | 1) {
  const target = index + direction;
  if (target < 0 || target >= form.fields.length) return;
  const temp = form.fields[index];
  form.fields[index] = form.fields[target];
  form.fields[target] = temp;
}

function submit() {
  if (form.fields.length === 0) {
    errorMsg.value = "请至少添加一个字段";
    return;
  }
  errorMsg.value = "";
  const config: CustomFormatConfig = {
    kind: form.kind,
    encoding: form.encoding,
    fields: form.fields.map((f) => ({ ...f })),
    sortByCode: form.sortByCode,
    commentPrefix: form.commentPrefix,
    startMarker: form.startMarker,
  };
  emit("confirmed", config);
  emit("update:show", false);
}

function getPreviewLine(): string {
  const parts: string[] = [];
  for (const f of form.fields) {
    switch (f.type) {
      case "word":
        parts.push("你好");
        break;
      case "pinyin":
        const sep = f.pinyinSeparator || "";
        parts.push(f.pinyinPrefix + "ni" + sep + "hao" + f.pinyinSuffix);
        break;
      case "code":
        parts.push("abcd");
        break;
      case "frequency":
        parts.push("100");
        break;
      case "rank":
        parts.push("1");
        break;
      case "tab":
        parts.push("\t");
        break;
      case "space":
        parts.push(" ");
        break;
      case "literal":
        parts.push(f.literal);
        break;
    }
  }
  return parts.join("");
}
</script>

<template>
  <n-modal
    :show="show"
    @update:show="emit('update:show', $event)"
    preset="card"
    title="自定义纯文本格式"
    style="width: 640px; max-width: 90vw"
    :mask-closable="false"
  >
    <n-form label-placement="left" label-width="100">
      <n-form-item label="类型">
        <n-select v-model:value="form.kind" :options="kindOptions" />
      </n-form-item>
      <n-form-item label="编码">
        <n-select v-model:value="form.encoding" :options="encodingOptions" />
      </n-form-item>
      <n-form-item label="注释前缀">
        <n-input v-model:value="form.commentPrefix" placeholder="如: #" style="width: 100px" />
      </n-form-item>
      <n-form-item label="词库开始标志">
        <n-input v-model:value="form.startMarker" placeholder="如: ###BEGIN###" style="width: 200px" />
      </n-form-item>
      <n-form-item label="按编码排序">
        <n-switch v-model:value="form.sortByCode" />
      </n-form-item>

      <div style="margin-bottom: 12px; font-weight: 500">字段配置</div>
      <div v-for="(field, index) in form.fields" :key="index" style="display: flex; gap: 8px; margin-bottom: 8px; align-items: flex-start">
        <n-select
          v-model:value="field.type"
          :options="fieldTypeOptions"
          style="width: 160px"
          size="small"
        />
        <template v-if="field.type === 'pinyin'">
          <n-input v-model:value="field.pinyinSeparator" placeholder="分隔符" style="width: 70px" size="small" />
          <n-input v-model:value="field.pinyinPrefix" placeholder="前缀" style="width: 60px" size="small" />
          <n-input v-model:value="field.pinyinSuffix" placeholder="后缀" style="width: 60px" size="small" />
        </template>
        <template v-if="field.type === 'literal'">
          <n-input v-model:value="field.literal" placeholder="字面量" style="width: 100px" size="small" />
        </template>
        <n-button size="small" @click="moveField(index, -1)" :disabled="index === 0">↑</n-button>
        <n-button size="small" @click="moveField(index, 1)" :disabled="index === form.fields.length - 1">↓</n-button>
        <n-button size="small" type="error" @click="removeField(index)">删</n-button>
      </div>
      <n-button dashed size="small" @click="addField" style="width: 100%">+ 添加字段</n-button>

      <div style="margin-top: 16px; padding: 8px 12px; background: #f5f5f5; border-radius: 4px; font-family: monospace; font-size: 13px">
        <n-text depth="3">预览: </n-text>
        <span style="white-space: pre">{{ getPreviewLine() }}</span>
      </div>
    </n-form>

    <div v-if="errorMsg" style="color: #d03050; margin-top: 8px">{{ errorMsg }}</div>

    <template #footer>
      <n-space justify="end">
        <n-button @click="emit('update:show', false)">取消</n-button>
        <n-button type="primary" @click="submit">确定</n-button>
      </n-space>
    </template>
  </n-modal>
</template>
