<script lang="ts" setup>
</script>

<template>
  <n-tabs
    class="card-tabs"
    default-value="convert"
    size="large"
    animated
    style="margin: 0 -4px"
    pane-style="padding-left: 4px; padding-right: 4px; box-sizing: border-box;"
  >
    <n-tab-pane name="convert" tab="转换">
      <Convert />
    </n-tab-pane>
    <n-tab-pane name="encode" tab="编码">
      <Encode />
    </n-tab-pane>
    <n-tab-pane name="shorten" tab="全码转简码">
      <Shorten />
    </n-tab-pane>
  </n-tabs>
</template>

<style>
#app {
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 2em;
}

.card-tabs {
  min-width: 450px;
  min-height: 300px;
}
</style>
