<script setup lang="ts">
import Converter from "./components/Converter.vue";
</script>

<template>
  <Suspense>
    <n-dialog-provider>
      <Converter msg="Vite + Vue" />
    </n-dialog-provider>
  </Suspense>
</template>

<style scoped></style>
