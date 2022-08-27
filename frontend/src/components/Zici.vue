<script lang="ts" setup>
import { ref } from "vue";
import { SelectFile, ConvDict } from "../../wailsjs/go/main/App";

let inputPath = ref("");
function selectInputFile() {
  SelectFile().then((result) => {
    if (result != "") {
      inputPath.value = result;
    }
  });
}

function saveFile() {
  if (inputPath.value && srcFormat.value && targetFormat.value) {
    ConvDict(inputPath.value, srcFormat.value, targetFormat.value,false)
  }
}

let srcFormat = ref("");
const srcOpts = ref([
  {
    label: "百度手机自定义方案.def",
    value: "baidu_def",
  },
  {
    label: "极点.mb",
    value: "jidian_mb",
  },
  {
    label: "fcitx4.mb",
    value: "fcitx4_mb",
  },
  {
    label: "多多.txt",
    value: "duoduo",
  },
  {
    label: "冰凌.txt",
    value: "bingling",
  },
  {
    label: "极点.txt",
    value: "jidian",
  },
]);

let targetFormat = ref("");
const targetOpts = ref([
  {
    label: "百度手机自定义方案.def",
    value: "baidu_def",
  },
  {
    label: "多多.txt",
    value: "duoduo",
  },
  {
    label: "冰凌.txt",
    value: "bingling",
  },
  {
    label: "极点.txt",
    value: "jidian",
  },
]);
</script>

<template>
  <n-space vertical>
    <n-space justify="space-between" style="margin-right: 10px">
      <n-input
        readonly
        type="text"
        placeholder="文件路径..."
        :value="inputPath"
        style="min-width: 300px"
      ></n-input>
      <n-button @click="selectInputFile"> 选择词库 </n-button>
    </n-space>

    <n-space
      ><n-space vertical class="select">
        <n-text>源格式</n-text>
        <n-space vertical>
          <n-select v-model:value="srcFormat" :options="srcOpts" /> </n-space
      ></n-space>
      <n-space vertical class="select">
        <n-text>目标格式</n-text>
        <n-space vertical>
          <n-select v-model:value="targetFormat" :options="targetOpts" />
        </n-space>
      </n-space>
    </n-space>
    <p></p>
    <n-space justify="end" style="margin-right: 10px">
      <n-button class="button" type="primary" @click="saveFile">
        转换并保存
      </n-button>
    </n-space>
  </n-space>
</template>

<style>
.select {
  width: 210px;
}
</style>
