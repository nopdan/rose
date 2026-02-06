<script setup lang="ts">
import { ref } from "vue";
import {
  NUpload,
  NUploadDragger,
  NIcon,
  NText,
  NTag,
  type UploadFileInfo,
} from "naive-ui";
import { CloudUploadOutline } from "@vicons/ionicons5";
import { uploadFile, type UploadResult } from "../api";

const emit = defineEmits<{
  uploaded: [result: UploadResult];
}>();

const fileInfo = ref<UploadResult | null>(null);
const uploading = ref(false);
const errorMsg = ref("");

function formatSize(bytes: number): string {
  if (bytes < 1024) return bytes + " B";
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + " KB";
  return (bytes / (1024 * 1024)).toFixed(1) + " MB";
}

async function handleUpload({ file }: { file: UploadFileInfo }) {
  if (!file.file) return;
  uploading.value = true;
  errorMsg.value = "";
  try {
    const result = await uploadFile(file.file);
    fileInfo.value = result;
    emit("uploaded", result);
  } catch (e: any) {
    errorMsg.value = e.message || "上传失败";
  } finally {
    uploading.value = false;
  }
}
</script>

<template>
  <div>
    <n-upload
      :custom-request="({ file }) => handleUpload({ file })"
      :show-file-list="false"
      :multiple="false"
    >
      <n-upload-dragger>
        <div style="padding: 20px 0">
          <n-icon size="48" :depth="3">
            <CloudUploadOutline />
          </n-icon>
          <n-text style="display: block; margin-top: 8px; font-size: 16px">
            点击或拖拽文件到此处上传
          </n-text>
          <n-text depth="3" style="font-size: 12px">
            支持各种输入法词库文件
          </n-text>
        </div>
      </n-upload-dragger>
    </n-upload>

    <div v-if="fileInfo" style="margin-top: 12px; display: flex; align-items: center; gap: 8px">
      <n-tag type="success" size="medium">已上传</n-tag>
      <span>{{ fileInfo.filename }}</span>
      <n-text depth="3">({{ formatSize(fileInfo.size) }})</n-text>
    </div>

    <div v-if="uploading" style="margin-top: 8px">
      <n-text depth="3">上传中...</n-text>
    </div>

    <div v-if="errorMsg" style="margin-top: 8px; color: #d03050">
      {{ errorMsg }}
    </div>
  </div>
</template>
