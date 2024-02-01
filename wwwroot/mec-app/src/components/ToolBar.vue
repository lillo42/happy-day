<template>
  <header class="tool-bar shadow-3">
    <img alt="Vue logo" class="logo" src="@/assets/logo.svg" width="24" height="24" />

    <RouterLink to="/products">Products</RouterLink>
  </header>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'
</script>

<style scoped>
.tool-bar {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  padding: 0 1rem;
  width: 100%;
}
</style>
