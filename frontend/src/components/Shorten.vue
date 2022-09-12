<template>
  <n-form label-placement="left" label-width="auto">
    <n-form-item label="码表">
      <n-input
        readonly
        type="text"
        placeholder="文件路径..."
        :value="inputPath"
        style="margin-right: 1em"
      ></n-input>
      <n-button @click="selectInputFile" type="primary" ghost> 选择 </n-button>
    </n-form-item>
    <n-form-item label="规则">
      <n-input
        type="text"
        placeholder="生成规则..."
        v-model:value="rule"
        style="margin-right: 1em"
      >
      </n-input>
      <n-button class="button" type="primary" @click="shorten"> 保存 </n-button>
    </n-form-item>
  </n-form>

  <div>
    <span style="color: red">仅支持多多格式码表</span><br />
    出简不出全规则，逗号，冒号分隔，默认 1，n 无限<br />
    1:0, 2:3, 3:2, 6:n -> <br />
    无 1 简，2 码 3 重，3 码 2 重，4 码 1 重，5 码 1 重，6 码无限重
  </div>
</template>

<script lang="ts" setup>
import { SelectFile, Shorten } from '../../wailsjs/go/main/App';

const inputPath = ref('');
function selectInputFile() {
  SelectFile().then((result) => {
    if (result != '') {
      inputPath.value = result;
    }
  });
}

const rule = ref('4:n');
function shorten() {
  if (inputPath.value && rule.value) {
    Shorten(inputPath.value, rule.value);
  }
}
</script>
