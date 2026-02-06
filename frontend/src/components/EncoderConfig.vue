<script setup lang="ts">
import { ref, computed, watch } from "vue";
import {
  NCard,
  NSelect,
  NSwitch,
  NFormItem,
  NUpload,
  NButton,
  NTag,
  NText,
  NIcon,
} from "naive-ui";
import { CloudUploadOutline } from "@vicons/ionicons5";
import { uploadFile, type EncoderConfig, type FormatInfo } from "../api";

const props = defineProps<{
  outputFormat: FormatInfo | null;
}>();

const emit = defineEmits<{
  "update:config": [config: EncoderConfig | null];
}>();

const wubiSchema = ref("86");
const useAABC = ref(true);
const codeTableFileId = ref("");
const codeTableFilename = ref("");

const schemaOptions = [
  { label: "不修改", value: "none" },
  { label: "86五笔", value: "86" },
  { label: "98五笔", value: "98" },
  { label: "新世纪", value: "06" },
  { label: "自定义码表", value: "custom" },
];

const needEncoder = computed(() => {
  if (!props.outputFormat) return false;
  return props.outputFormat.kind === 2; // 五笔需要 encoder
});

// 是否显示编码器配置面板
const showPanel = computed(() => needEncoder.value);

// 是否选择了"不修改"
const isNone = computed(() => wubiSchema.value === "none");

watch(
  [wubiSchema, useAABC, codeTableFileId, () => props.outputFormat],
  () => {
    emitConfig();
  },
  { immediate: true },
);

function emitConfig() {
  if (!props.outputFormat) {
    emit("update:config", null);
    return;
  }
  // 拼音格式不需要编码器配置
  if (props.outputFormat.kind === 1) {
    emit("update:config", {
      type: "pinyin",
      schema: "",
      codeTableFileId: "",
      useAABC: false,
    });
    return;
  }
  // 五笔格式
  if (props.outputFormat.kind === 2) {
    if (isNone.value) {
      emit("update:config", {
        type: "none",
        schema: "",
        codeTableFileId: "",
        useAABC: false,
      });
      return;
    }
    emit("update:config", {
      type: "wubi",
      schema: wubiSchema.value === "custom" ? "custom" : wubiSchema.value,
      codeTableFileId:
        wubiSchema.value === "custom" ? codeTableFileId.value : "",
      useAABC: useAABC.value,
    });
    return;
  }
  emit("update:config", null);
}

async function handleCodeTableUpload({ file }: any) {
  if (!file.file) return;
  try {
    const result = await uploadFile(file.file);
    codeTableFileId.value = result.id;
    codeTableFilename.value = result.filename;
  } catch (e: any) {
    console.error("码表上传失败", e);
  }
}
</script>

<template>
  <div v-if="showPanel">
    <n-card title="编码器配置" size="small" style="margin-top: 12px">
      <n-form-item label="五笔方案" label-placement="left" label-width="80">
        <n-select v-model:value="wubiSchema" :options="schemaOptions" style="width: 200px" />
      </n-form-item>

      <template v-if="!isNone">
        <div v-if="wubiSchema === 'custom'" style="margin-bottom: 12px">
          <n-form-item label="上传码表" label-placement="left" label-width="80">
            <div>
              <n-upload
                :custom-request="({ file }) => handleCodeTableUpload({ file })"
                :show-file-list="false"
              >
                <n-button size="small">
                  <template #icon><n-icon><CloudUploadOutline /></n-icon></template>
                  选择码表文件
                </n-button>
              </n-upload>
              <div v-if="codeTableFilename" style="margin-top: 4px">
                <n-tag size="small" type="success">{{ codeTableFilename }}</n-tag>
              </div>
              <n-text depth="3" style="font-size: 12px; display: block; margin-top: 4px">
                格式: 每行「字\t编码」，如「我\twqiy」
              </n-text>
            </div>
          </n-form-item>
        </div>

        <n-form-item label="组词规则" label-placement="left" label-width="80">
          <n-switch v-model:value="useAABC" />
          <n-text depth="3" style="margin-left: 8px; font-size: 12px">
            三字词：{{ useAABC ? 'AABC（前二次一末一）' : 'ABCC（首一次一末二）' }}
          </n-text>
        </n-form-item>
      </template>
    </n-card>
  </div>
</template>
