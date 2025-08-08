<script setup lang="ts">
//@ts-expect-error: any implícito en el componente
import Vue3Datatable from '@bhplugin/vue3-datatable'
import '@bhplugin/vue3-datatable/dist/style.css'
import { ref } from 'vue'

const loading = ref(false)
const cols =
  ref([
    { field: 'ticker', title: 'Ticker', isUnique: false, sort:true, filter:true },
    { field: 'company', title: 'Company', sort:true, filter:true },
    { field: 'brokerage', title: 'Brokerage', sort:true, filter:true },
    { field: 'action', title: 'Action', sort:false, filter:false },
    { field: 'rating_from', title: 'Rating From',sort:false, filter:true },
    { field: 'rating_to', title: 'Rating To',sort:false, filter:true },
    { field: 'target_from', title: 'Target From', type: 'number', sort:true, filter:true },
    { field: 'target_to', title: 'Target To', type: 'number', sort:true, filter:true },
  ]) || [];

</script>

<template>
  <vue3-datatable :columns="cols" :loading="loading" :sortable="true" :columnFilter="false">
  </vue3-datatable>
</template>

<style scoped></style>
