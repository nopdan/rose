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
    ConvDict(inputPath.value, srcFormat.value, targetFormat.value,true)
  }
}

let srcFormat = ref("");
const srcOpts = ref([
  {
    label: "纯汉字",
    value: "word_only",
  },
  {
    label: "搜狗细胞词库.scel",
    value: "sogou_scel",
  },
  {
    label: "搜狗拼音备份词库.bin",
    value: "sogou_bin",
  },
  {
    label: "qq拼音v6以下词库.qpyd",
    value: "qq_qpyd",
  },
  {
    label: "百度分类词库.bdict(bcd)",
    value: "baidu_bdict",
  },
  {
    label: "华宇紫光拼音词库.uwl",
    value: "ziguang_uwl",
  },
  {
    label: "微软拼音用户自定义短语.dat",
    value: "mspy_dat",
  },
  {
    label: "微软拼音自学习词汇.dat",
    value: "mspy_udl",
  },
  {
    label: "搜狗拼音.txt",
    value: "sogou",
  },
  {
    label: "qq 拼音.txt",
    value: "qq",
  },
  {
    label: "百度拼音.txt",
    value: "baidu",
  },
  {
    label: "谷歌拼音.txt",
    value: "google",
  },
  {
    label: "拼音加加.txt",
    value: "pyjj",
  },
]);

let targetFormat = ref("");
const targetOpts = ref([
  {
    label: "搜狗拼音.txt",
    value: "sogou",
  },
  {
    label: "qq 拼音.txt",
    value: "qq",
  },
  {
    label: "百度拼音.txt",
    value: "baidu",
  },
  {
    label: "谷歌拼音.txt",
    value: "google",
  },
  {
    label: "拼音加加.txt",
    value: "pyjj",
  },
  {
    label: "纯汉字",
    value: "word_only",
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
