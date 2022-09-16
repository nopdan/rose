<template>
  <div style="display: flex; margin-bottom: 0.8em">
    <n-radio-group v-model:value="ime" name="ime">
      <n-radio-button key="pinyin" value="pinyin">拼音</n-radio-button>
      <n-radio-button key="wubi" value="wubi">五笔</n-radio-button>
    </n-radio-group>
    <n-input
      readonly
      type="text"
      placeholder="文件路径..."
      :value="inputPath"
      style="margin: 0 2%"
    ></n-input>
    <n-button @click="selectInputFile" type="primary" ghost> 选择 </n-button>
  </div>
  <div style="display: flex">
    <n-space vertical style="width: 48%">
      <div>源格式</div>
      <n-select
        placement="right"
        v-model:value="srcFormat"
        :options="srcOpts"
      />
    </n-space>
    <div style="width: 4%"></div>
    <n-space vertical style="width: 48%">
      <div>输出格式</div>
      <n-select
        placement="left"
        v-model:value="targetFormat"
        :options="targetOpts"
      />
    </n-space>
  </div>

  <n-form
    ref="formRef"
    inline
    label-width="auto"
    style="margin-top: 1em; height: 84px"
  >
    <n-form-item label="转为双拼" v-if="isPinyin" style="width: 15%">
      <n-switch v-model:value="doublePinyin"> </n-switch>
    </n-form-item>
    <n-form-item v-if="doublePinyin" label="三字词规则" style="width: 25%">
      <n-select v-model:value="dpRule" :options="dpOpts" />
    </n-form-item>
    <n-form-item v-if="doublePinyin" label="映射表" style="width: 40%">
      <n-input
        readonly
        type="text"
        placeholder="文件路径..."
        :value="dpCfg"
      ></n-input>
    </n-form-item>
    <div style="margin: auto"></div>
    <n-form-item v-if="doublePinyin">
      <n-button @click="selectDpCfg" type="primary" ghost> 选择 </n-button>
    </n-form-item>
  </n-form>

  <div style="float: right">
    <n-button
      class="button"
      type="primary"
      @click="saveFile"
      :disabled="disableBtn"
    >
      转换并保存
    </n-button>
  </div>
</template>
<script lang="ts" setup>
import { SelectFile, Convert } from '../../wailsjs/go/main/App';

const ime = ref('pinyin');
const inputPath = ref('');
const isPinyin = computed(() => {
  return ime.value === 'pinyin';
});

function selectInputFile() {
  SelectFile().then((result) => {
    if (result != '') {
      inputPath.value = result;
    }
  });
}

const srcFormat = ref('');
const targetFormat = ref('');
const srcOpts = computed(() => {
  return isPinyin.value ? srcDict : srcTable;
});
const targetOpts = computed(() => {
  return isPinyin.value && !doublePinyin.value ? targetDict : targetTable;
});

const doublePinyin = ref(false);
const dpRule = ref(0);
const dpOpts = [
  { label: 'AABC', value: 0 },
  { label: 'ABCC', value: 1 },
];
const dpCfg = ref('');
function selectDpCfg() {
  SelectFile().then((result) => {
    if (result != '') {
      dpCfg.value = result;
    }
  });
}

const disableBtn = computed(() => {
  return !(inputPath.value && srcFormat.value && targetFormat.value);
});

watch(ime, () => {
  srcFormat.value = '';
  targetFormat.value = '';
  doublePinyin.value = false;
});

watch(doublePinyin, () => {
  targetFormat.value = doublePinyin.value ? 'duoduo' : '';
});

function saveFile() {
  Convert(
    inputPath.value,
    srcFormat.value,
    targetFormat.value,
    isPinyin.value,
    doublePinyin.value,
    dpRule.value,
    dpCfg.value
  );
}

const srcDict = [
  {
    label: '纯汉字',
    value: 'word_only',
  },
  {
    label: '搜狗细胞词库.scel',
    value: 'sogou_scel',
  },
  {
    label: '搜狗拼音备份词库.bin',
    value: 'sogou_bin',
  },
  {
    label: 'qq拼音v6以下词库.qpyd',
    value: 'qq_qpyd',
  },
  {
    label: '百度分类词库.bdict(bcd)',
    value: 'baidu_bdict',
  },
  {
    label: '华宇紫光拼音词库.uwl',
    value: 'ziguang_uwl',
  },
  {
    label: '微软用户自定义短语.dat',
    value: 'mspy_dat',
  },
  {
    label: '微软拼音自学习词汇.dat',
    value: 'mspy_udl',
  },
  {
    label: '搜狗拼音.txt',
    value: 'sogou',
  },
  {
    label: 'qq 拼音.txt',
    value: 'qq',
  },
  {
    label: '百度拼音.txt',
    value: 'baidu',
  },
  {
    label: '谷歌拼音.txt',
    value: 'google',
  },
  {
    label: '拼音加加.txt',
    value: 'pyjj',
  },
];

const targetDict = [
  {
    label: '微软用户自定义短语.dat',
    value: 'mspy_dat',
  },
  {
    label: '搜狗拼音.txt',
    value: 'sogou',
  },
  {
    label: 'qq 拼音.txt',
    value: 'qq',
  },
  {
    label: '百度拼音.txt',
    value: 'baidu',
  },
  {
    label: '谷歌拼音.txt',
    value: 'google',
  },
  {
    label: '拼音加加.txt',
    value: 'pyjj',
  },
  {
    label: '纯汉字',
    value: 'word_only',
  },
];

const srcTable = [
  {
    label: '微软用户自定义短语.dat',
    value: 'msudp_dat',
  },
  {
    label: '微软五笔.lex',
    value: 'mswb_dat',
  },
  {
    label: '百度手机自定义方案.def',
    value: 'baidu_def',
  },
  {
    label: '极点.mb',
    value: 'jidian_mb',
  },
  {
    label: 'fcitx4.mb',
    value: 'fcitx4_mb',
  },
  {
    label: '多多.txt',
    value: 'duoduo',
  },
  {
    label: '冰凌.txt',
    value: 'bingling',
  },
  {
    label: '极点.txt',
    value: 'jidian',
  },
];

const targetTable = [
  {
    label: '微软用户自定义短语.dat',
    value: 'msudp_dat',
  },
  {
    label: '百度手机自定义方案.def',
    value: 'baidu_def',
  },
  {
    label: '多多.txt',
    value: 'duoduo',
  },
  {
    label: '冰凌.txt',
    value: 'bingling',
  },
  {
    label: '极点.txt',
    value: 'jidian',
  },
];
</script>
