<script lang="ts" setup>
import { ref } from "vue";
let value1 = ref(null);
let value2 = ref(null);
let value3 = ref(null);
</script>

<template>
  <n-space id="ipt">
    <n-input
      id="ipt1"
      v-model:value="value1"
      type="textarea"
      :autosize="{
        minRows: 8,
        maxRows: 8,
      }"
      placeholder="在这里输入要处理的词组……"
    ></n-input>
    <n-space vertical align="center">
      <n-input
        id="ipt2"
        v-model:value="value2"
        type="textarea"
        :autosize="{
          minRows: 8,
          maxRows: 8,
        }"
        placeholder="在这里输入规则……"
      ></n-input>
    </n-space>
    <n-input
      id="ipt3"
      v-model:value="value3"
      type="textarea"
      :autosize="{
        minRows: 8,
        maxRows: 8,
      }"
      placeholder="这里是输出的结果……"
    ></n-input>
  </n-space>

  <n-space justify="space-between">
    <n-space align="center">
      <n-text>&nbsp;词　库&nbsp;</n-text>
      <n-input class="filedir"></n-input>
      <n-button>选择</n-button>
    </n-space>
    <n-button class="button">合并</n-button>
  </n-space>

  <n-space justify="space-between">
    <n-space align="center">
      <n-text>单字码表</n-text>
      <n-input class="filedir"></n-input>
      <n-button> 选择 </n-button>
    </n-space>
    <n-button class="button">处理</n-button>
  </n-space>
</template>

<style>
#ipt1 {
  width: 160px;
}
#ipt2 {
  width: 160px;
}
#ipt3 {
  width: 280px;
}
</style>
