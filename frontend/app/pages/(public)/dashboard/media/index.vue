<script setup lang="ts">
import FileDrop from '~/components/FileDrop.vue';
import { refreshAccessToken } from '~/helpers/apiClient';
import type { Node } from '~/stores/interfaces';

const dropComponent = ref();
const selectedFile: Ref<File | undefined> = ref();

const selectFile = (file?: File) => (selectedFile.value = file);

//// test code
const config = useRuntimeConfig();
const mediaUrl = `${config.public.baseApi}/media`;
const isLoading = ref(false);
const fileLink: Ref<string | undefined> = ref(undefined);


const mediaStore = useMediaStore();

const submitFile = async () => {
  if (!selectedFile.value) return;
  isLoading.value = true;
  const body = new FormData();
  body.append('file', selectedFile.value);
  selectedFile.value = undefined;
  dropComponent.value.reset();
  await mediaStore
    .post(body)
    .then(r => (fileLink.value = `${mediaUrl}/${(r as Node).user_id}/${(r as Node).metadata?.transformed_path as string}`))
    .catch(e => console.error(e))
    .finally(() => (isLoading.value = false));
};

const uploadTestFile = async () => {
  const fileUrl = "https://raw.githubusercontent.com/gin-gonic/logo/master/color.png";

  try {
    const response = await fetch(fileUrl);
    const blob = await response.blob();

    const file = new File([blob], "test.png", { type: blob.type });

    selectedFile.value = file;

    await submitFile();
  } catch (error) {
    console.error("Failed to upload test file:", error);
  }
}

const refreshSession = async () => {
  await refreshAccessToken();
}

////
</script>

<template>
  Media
  <FileDrop ref="dropComponent" @select="selectFile" />
  <button @click="refreshSession">Refresh session</button>
  <button @click="uploadTestFile">Upload test file</button>
  <p>Expected</p>
  <img src="https://raw.githubusercontent.com/gin-gonic/logo/master/color.png" alt="test"
    style="width: 50px; height: 50px;">
  <p v-if="isLoading">Uploading file...</p>
  <p>Test Serve Existing Upload</p>
  <img :src="`${mediaUrl}/2891564074434561/3659081319481348.png`" alt="test"
    style="width: 50px; height: 50px; background-color: grey">
  <p v-if="fileLink">Actual</p>
  <img v-if="fileLink" :src="`${fileLink}`" alt="test" style="width: 50px; height: 50px;">
</template>