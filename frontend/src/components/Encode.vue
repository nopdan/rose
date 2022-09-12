<template>
  <n-form label-placement="left" label-width="auto">
    <n-form-item label="单字全码表">
      <n-input
        readonly
        type="text"
        placeholder="文件路径..."
        :value="charPath"
        style="margin-right: 5%"
      ></n-input>
      <n-button @click="selectFile(true)" type="primary" ghost> 选择 </n-button>
    </n-form-item>
    <n-form-item label="词库">
      <n-input
        readonly
        type="text"
        placeholder="文件路径..."
        :value="dictPath"
        style="margin-right: 5%"
      ></n-input>
      <n-button @click="selectFile(false)" type="primary" ghost>
        选择
      </n-button>
    </n-form-item>
    <n-form-item label="造词规则">
      <n-input type="text" v-model:value="encRule"></n-input>
    </n-form-item>
  </n-form>
  <div style="display: flex; margin-top: 1em">
    <div style="color: red;margin-right: auto">仅支持多多格式码表</div>
    <n-button
      class="button"
      type="primary"
      @click="encode(true)"
      ghost
      :disabled="disableBtn"
      style="margin-right: 1em"
    >
      校验
    </n-button>
    <n-button
      class="button"
      type="primary"
      @click="encode(false)"
      :disabled="disableBtn"
    >
      编码
    </n-button>
  </div>
</template>
<script lang="ts" setup>
import { SelectFile,Encode } from '../../wailsjs/go/main/App';

const charPath = ref('');
const dictPath = ref('');
const encRule = ref('2=AaAbBaBb,3=AaAbBaCa,0=AaBaCaZa');
const disableBtn = computed(() => {
  return !(charPath.value && dictPath.value && encRule.value);
});

function selectFile(flag: boolean) {
  SelectFile().then((result) => {
    if (result != '') {
      flag ? (charPath.value = result) : (dictPath.value = result);
    }
  });
}
function encode(isCheck:boolean) {
  Encode(charPath.value,dictPath.value,encRule.value,isCheck)
}
</script>
