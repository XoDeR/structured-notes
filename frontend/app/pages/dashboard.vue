<script setup lang="ts">
useNodeStore().fetch();
console.log("nodes fetched");
</script>

<template>
  <NuxtPage />
</template>