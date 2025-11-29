<script setup lang="ts">
import FileDrop from '~/components/FileDrop.vue';

const dropComponent = ref();
const selectedFile: Ref<File | undefined> = ref();

const selectFile = (file?: File) => (selectedFile.value = file);
</script>

<template>
  Media
  <FileDrop ref="dropComponent" @select="selectFile" />
</template>