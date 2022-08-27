<script lang="ts" setup>
import { ref } from "vue";
import { SelectFile, Shorten } from "../../wailsjs/go/main/App";

let inputPath = ref("");
function selectInputFile() {
  SelectFile().then((result) => {
    if (result != "") {
      inputPath.value = result;
    }
  });
}

let rule = ref("");
function shorten() {
  if (inputPath.value && rule.value) {
    Shorten(inputPath.value, rule.value);
  }
}
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

    <p></p>
    <n-space justify="space-between" style="margin-right: 10px">
      <n-input
        type="text"
        placeholder="生成规则..."
        v-model:value="rule"
        style="min-width: 300px"
      >
      </n-input>
      <n-button class="button" type="primary" @click="shorten">
        转换并保存
      </n-button>
    </n-space>
    <p>
      <span style="color: red">仅支持多多格式码表</span><br />
      出简不出全规则，逗号，冒号分隔，默认 1，n 无限<br />
      1:0, 2:3, 3:2, 6:n -> <br />
      无 1 简，2 码 3 重，3 码 2 重，4 码 1 重，5 码 1 重，6 码无限重
    </p>
  </n-space>
</template>
