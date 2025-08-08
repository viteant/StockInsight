<script setup lang="ts">
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { Icon } from '@iconify/vue'
import FilterForm from '@/components/forms/FilterForm.vue'
import StockDataTable from '@/components/data/StockDataTable.vue'
import { useStockStore } from '@/stores/stock'
import { computed, watch } from 'vue'

const stockStore = useStockStore()

const hasFilters = computed(() =>
  Object.entries({
    ...stockStore.getFilters,
    page: undefined,
    limit: undefined,
  }).some(([_, value]) => {
    return value !== null && value !== undefined && String(value).trim() !== ''
  }),
)

const clearFilters = () => {
  stockStore.setLoading(true)
  stockStore.resetFilters()
}

</script>

<template>
  <div class="bg-[#5a5bf3] p-2 flex flex-row items-center justify-between">
    <h2 class="text-white text-xl">Stocks</h2>
    <div class="flex items-center">
      <button
        v-if="hasFilters"
        class="border border-white rounded-lg text-white px-5 mr-2 cursor-pointer hover:bg-indigo-600 transition-all duration-400"
        @click.prevent="clearFilters"
      >
        Clear Filter
      </button>
      <Popover>
        <PopoverTrigger>
          <div class="flex border border-white px-5  rounded text-white text-base items-center cursor-pointer hover:bg-indigo-600 transition-all duration-400">
            <p class="mr-1">Filters</p> <Icon icon="mdi:filter-gear" class="text-white text-xl" />
          </div>
        </PopoverTrigger>
        <PopoverContent class="bg-white border-indigo-200">
          <FilterForm />
        </PopoverContent>
      </Popover>
    </div>
  </div>
  <StockDataTable />
</template>

<style scoped></style>
