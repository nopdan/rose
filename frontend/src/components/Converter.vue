<template>
  <n-form
    ref="formRef"
    :model="model"
    label-placement="left"
    :label-width="100"
    :style="{
      minWidth: '640px',
      maxWidth: '800px',
    }"
  >
    <div style="height: 2rem"></div>
    <n-card title="输入" hoverable content-style="padding-bottom: unset;">
      <n-form-item label=" " path="uploadValue">
        <n-upload
          v-model:action="uploadUrl"
          :max="1"
          v-model:file-list="fileList"
        >
          <n-upload-dragger>
            <div style="margin-bottom: 10px">
              <n-icon size="48" :depth="3">
                <archive-icon />
              </n-icon>
            </div>
            <n-text style="font-size: 16px"> 点击或者拖动文件到该区域 </n-text>
          </n-upload-dragger>
        </n-upload>
      </n-form-item>
      <n-form-item label="源格式" path="inputFormat" style="margin-bottom: 0">
        <n-select
          v-model:value="model.inputFormat"
          placeholder="Select"
          :options="generalOptions()"
        />
      </n-form-item>
    </n-card>
    <div style="height: 1rem"></div>

    <n-card title="输出" hoverable content-style="padding-bottom: unset;">
      <n-form-item label="输出方案" path="schemaGroupValue">
        <n-radio-group v-model:value="model.kind" name="schemaGroup">
          <n-radio value="pinyin"> 拼音 </n-radio>
          <n-radio value="wubi"> 五笔 </n-radio>
          <n-radio value="words"> 纯词组 </n-radio>
        </n-radio-group>
      </n-form-item>
      <n-form-item label="输出格式" path="outputFormat">
        <n-select
          v-model:value="model.outputFormat"
          placeholder="Select"
          :options="outputFormatOptions"
        />
      </n-form-item>

      <n-form-item
        label="编码方案"
        path="mbGroupValue"
        v-if="model.kind === 'wubi'"
      >
        <n-radio-group v-model:value="model.schema" name="mbgroup">
          <n-radio value="original"> 原始 </n-radio>
          <n-radio value="phrase"> 拼音 </n-radio>
          <n-radio value="86"> 86 五笔 </n-radio>
          <n-radio value="98"> 98 五笔 </n-radio>
          <n-radio value="06"> 新世纪 </n-radio>
          <n-radio value="custom"> 自定义 </n-radio>
        </n-radio-group>
      </n-form-item>

      <n-form-item label=" " path="uploadValue" v-if="displayCustom">
        <n-upload
          v-model:action="uploadMbUrl"
          :max="1"
          v-model:file-list="mbList"
        >
          <n-upload-dragger>
            <div style="margin-bottom: 10px">
              <n-icon size="48" :depth="3">
                <archive-icon />
              </n-icon>
            </div>
            <n-text style="font-size: 16px"> 点击或者拖动文件到该区域 </n-text>
          </n-upload-dragger>
        </n-upload>
      </n-form-item>

      <n-form-item
        label="三字词规则"
        path="ruleGroupValue"
        v-if="displayCustom"
      >
        <n-radio-group v-model:value="model.rule" name="rulegroup">
          <n-radio value="AABC"> AABC </n-radio>
          <n-radio value="ABCC"> ABCC </n-radio>
        </n-radio-group>
      </n-form-item>
    </n-card>
    <n-card :bordered="false">
      <div style="float: right">
        <n-button
          class="button"
          type="primary"
          @click="handleClick"
          v-model:disabled="disableButton"
        >
          转换并保存
        </n-button>
      </div>
    </n-card>
  </n-form>
</template>

<script setup lang="ts">
import { defineComponent, computed, ref, watch } from "vue";
import { ArchiveOutline as ArchiveIcon } from "@vicons/ionicons5";
import axios from "axios";
import type { UploadFileInfo } from "naive-ui";
import { useDialog } from "naive-ui";
const debugMode = false;
const host = debugMode ? "http://localhost:8080" : "";
const uploadUrl = host + "/upload";
const uploadMbUrl = uploadUrl + "/mb";

const displayCustom = computed(() => {
  return model.value.kind === "wubi" && model.value.schema === "custom";
});
const formRef = ref(null);
const model = ref({
  name: "",
  inputFormat: "sgbak",
  kind: "pinyin",
  schema: "original",
  rule: "AABC",
  outputFormat: "sogou",
});

let data: { name: any; id: any; canMarshal: any; kind: any }[];
data = await fetch(host + "/list").then((res) => res.json());

const generalOptions = () => {
  let opt: any[] = [];
  data.forEach((item) => {
    opt.push({
      label: item.name,
      value: item.id,
    });
  });
  return opt;
};

console.log(generalOptions());

const pinyinFormatOptions = () => {
  let opt: any[] = [];
  data.forEach((item) => {
    if (item.kind === 1 && item.canMarshal) {
      opt.push({
        label: item.name,
        value: item.id,
      });
    }
  });
  return opt;
};

const wubiFormatOptions = () => {
  let opt: any[] = [];
  data.forEach((item) => {
    if (item.kind === 2 && item.canMarshal) {
      opt.push({
        label: item.name,
        value: item.id,
      });
    }
  });
  return opt;
};

const outputFormatOptions = computed(() => {
  if (model.value.kind === "pinyin") {
    return pinyinFormatOptions();
  } else if (model.value.kind === "wubi") {
    return wubiFormatOptions();
  } else {
    return [{ label: "纯词组", value: "words" }];
  }
});

watch(
  () => model.value.kind,
  () => {
    if (model.value.kind === "pinyin") {
      model.value.outputFormat = "sogou";
    } else if (model.value.kind === "wubi") {
      model.value.outputFormat = "def";
    } else {
      model.value.outputFormat = "words";
    }
  }
);

const fileList = ref<UploadFileInfo[]>([]);
const mbList = ref<UploadFileInfo[]>([]);
const disableButton = computed(() => {
  if (fileList.value.length === 0) {
    return true;
  }
  if (
    model.value.kind === "wubi" &&
    model.value.schema === "custom" &&
    mbList.value.length === 0
  ) {
    return true;
  }
  return false;
});

watch(
  () => fileList.value.length,
  () => {
    if (fileList.value.length === 1) {
      const fi = fileList.value[0];
      model.value.name = fi.name;
      // 根据后缀猜测格式
      const l = fi.name.split(".");
      console.log(l);
      let suffix = "";
      if (l.length >= 2) {
        suffix = l[l.length - 1];
      }
      console.log(suffix);
      generalOptions().forEach((item) => {
        if (item.value === suffix) {
          model.value.inputFormat = suffix;
        }
      });
    } else {
      model.value.name = "";
    }
  }
);

const dialog = useDialog();
function handleClick() {
  console.log(model.value);
  axios.post(host + "/api", JSON.stringify(model.value)).then((res) => {
    console.log(res);
    if (res.data !== "error") {
      dialog.success({
        title: "Success",
        content: "保存到：" + res.data,
        positiveText: "确定",
      });
    } else {
      dialog.error({
        title: "Error",
        content: "导出失败!",
        positiveText: "确定",
      });
    }
  });
}

defineComponent({
  components: {
    ArchiveIcon,
  },
});
</script>
