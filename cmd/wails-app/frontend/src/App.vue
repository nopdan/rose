<script lang="ts" setup>
import Pinyin from "./components/Pinyin.vue";
import Zici from "./components/Zici.vue";
// import Tool from "./components/Tool.vue";
</script>

<template>
  <div id="main">
    <n-tabs
      class="card-tabs"
      default-value="pinyin"
      size="large"
      animated
      style="margin: 0 -4px"
      pane-style="padding-left: 4px; padding-right: 4px; box-sizing: border-box;"
    >
      <n-tab-pane name="pinyin" tab="拼音词库转换">
        <Pinyin />
      </n-tab-pane>

      <n-tab-pane name="zici" tab="字词码表转换"> 
        <Zici />
      </n-tab-pane>
      <!--<n-tab-pane name="tool" tab="工具">
        <Tool />
      </n-tab-pane>-->
    </n-tabs>
  </div>
</template>

<style>
#main {
  margin-left: 10px;
  width: 450px;
  height: 220px;
  /* height: fill-available; */
  /* max-width: fill-available;*/
  /* border: 1px solid; */
}
</style>
