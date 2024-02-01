<script setup lang="ts">
import { RouterView } from 'vue-router'
import ToolBar from './components/ToolBar.vue'
</script>

<template>
  <div class="app-content gap-5">
    <ToolBar />
    <RouterView />
  </div>
</template>

<style scoped>
.app-content {
  display: flex;
  flex-direction: column;
  width: 100%;
}
</style>
