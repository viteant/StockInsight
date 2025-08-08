<script setup lang="ts">
//@ts-expect-error: any implícito en el componente
import Vue3Datatable from '@bhplugin/vue3-datatable'
import '@bhplugin/vue3-datatable/dist/style.css'
import { onMounted, ref } from 'vue'
import { useStockStore } from '@/stores/stock'
import { storeToRefs } from 'pinia'
import { useDateFormat } from '@vueuse/core'

const stockStore = useStockStore()
const { getStocks, total, page, limit, getLoading } = storeToRefs(stockStore)

const paginationChange = (data: any) => {
  switch (data.change_type) {
    case 'pagesize':
      stockStore.setLoading(true)
      stockStore.setLimit(data.pagesize)
      break
    case 'page':
      stockStore.setLoading(true)
      stockStore.setPage(data.current_page)
      break
    case 'sort':
      console.log(data)
      stockStore.setLoading(true)
      stockStore.setFilters({
        orderBy: data.sort_column,
        orderDir: data.sort_direction,
      })
      break
    default:
      console.log(data.change_type, data)
  }
}

const cols =[
  { field: 'ticker', title: 'Ticker', isUnique: false, sort: true },
  { field: 'company', title: 'Company', sort: true },
  { field: 'brokerage', title: 'Brokerage', sort: true },
  { field: 'action', title: 'Action', sort: true, minWidth: '150px', slotMode: true },
  { field: 'rating_from', title: 'Rating From', sort: true },
  { field: 'rating_to', title: 'Rating To', sort: true },
  {
    field: 'target_from',
    title: 'Target From',
    type: 'number',
    sort: true,
    cellRenderer: (v: any) => `$ ${v.target_from.toFixed(2)}`,
  },
  {
    field: 'target_to',
    title: 'Target To',
    type: 'number',
    sort: true,
    cellRenderer: (v: any) => `$ ${v.target_to.toFixed(2)}`,
  },
  {
    field: 'created_at',
    title: 'Time',
    type: 'date',
    sort: true,
    minWidth: '100px',
    cellRenderer: (v: any) => {
      return useDateFormat(v.created_at, 'DD-MM-YYYY').value
    },
  },
]

onMounted(() => {
  stockStore.refresh()
})
</script>

<template>
  <vue3-datatable
    :columns="cols"
    :rows="getStocks"
    :loading="getLoading"
    :sortable="true"
    :is-server-mode="true"
    :columnFilter="false"
    :totalRows="total"
    :page="page"
    skin="bh-table-compact"
    :pageSize="limit"
    :stickyHeader="true"
    @change="paginationChange"
    class="alt-pagination"
  >
    <template #action="{value}">
      {{ value.action.toLowerCase().replace(/(^\s*\w|[.!?]\s*\w)/g, (c:string) => c.toUpperCase()) }}
    </template>
  </vue3-datatable>
</template>

<style scoped></style>
