<script setup lang="ts">
import type { Field } from '~/components/DataTable.vue';


const view = ref<'table' | 'list'>('table');

const headers = [
  { label: 'Name', key: 'name' },
  { label: 'Size', key: 'size' },
  { label: 'Type', key: 'type' },
  { label: 'Parent doc', key: 'parent' },
  { label: 'Date', key: 'date' },
  { label: 'Action', key: 'action' },
];

//// Test data
const mediaFiles = ref([
  {
    name: "file01.png",
    size: "500KB",
    type: "image",
    parent: "Aaaa",
    date: new Date(),
  },
  {
    name: "file02.png",
    size: "500KB",
    type: "image",
    parent: "Aaaa",
    date: new Date(),
  },
  {
    name: "file03.png",
    size: "500KB",
    type: "image",
    parent: "Aaaa",
    date: new Date(),
  },
  {
    name: "file04.png",
    size: "500KB",
    type: "image",
    parent: "Aaaa",
    date: new Date(),
  },
  {
    name: "file05.png",
    size: "500KB",
    type: "image",
    parent: "Aaaa",
    date: new Date(),
  },
  {
    name: "file06.png",
    size: "500KB",
    type: "image",
    parent: "Aaaa",
    date: new Date(),
  },
  {
    name: "file07.png",
    size: "500KB",
    type: "image",
    parent: "Aaaa",
    date: new Date(),
  },
  {
    name: "file08.png",
    size: "500KB",
    type: "image",
    parent: "Aaaa",
    date: new Date(),
  },
]);
////

const rows: ComputedRef<Field[]> = computed(() =>
  mediaFiles.value.map(res => {
    return {
      name: { content: res.name, type: 'text' },
      size: { content: res.size, type: 'text' },
      type: { content: `<tag class="color:grey">${res.type}</tag>`, type: 'html' },
      parent: { content: `<tag class="color:grey">${res.parent}</tag>`, type: 'html' },
      date: { content: res.date.toLocaleDateString(), type: 'text' },
      action: { type: 'slot', data: res },
    };
  }),
);

</script>

<template>
  <h1>Documents</h1>

  <h3>create new document</h3>
  <NuxtLink to="/dashboard/docs/new" :prefetch="false">+</NuxtLink>
  <hr>

  <DataTable v-if="view === 'table'" :headers="headers" :rows="rows">
  </DataTable>
</template>