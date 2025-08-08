<script setup lang="ts">
import { reactive, watchEffect } from "vue"
import { useStockStore } from '@/stores/stock'

const props = withDefaults(defineProps<{
  initialFilters?: Partial<StockFilters>
  hidePagination?: boolean
  class?: string
}>(), {
  initialFilters: () => ({}),
  hidePagination: true,
  class: ""
})

const stockStore = useStockStore()

const filters = reactive<StockFilters>({
  ticker: stockStore?.filters?.ticker ?? "",
  company: stockStore?.filters?.company ?? "",
  brokerage: stockStore?.filters?.brokerage ?? "",
  target_from: stockStore?.filters?.target_from ?? undefined,
  target_to: stockStore?.filters?.target_to ?? undefined,
  date_from: stockStore?.filters?.date_from ?? "",
  date_to: stockStore?.filters?.date_to ?? "",
})

watchEffect(() => {
  Object.assign(filters, { ...filters, ...props.initialFilters })
})



function onSubmit() {
  stockStore.setLoading(true)
  stockStore.setFilters(filters)
}

function onReset() {
  filters.ticker = ""
  filters.company = ""
  filters.brokerage = ""
  filters.target_from = null
  filters.target_to = null
  filters.date_from = ""
  filters.date_to = ""
  stockStore.setLoading(true)
  stockStore.setFilters(filters)

}
</script>

<template>
  <form
    :class="['w-full p-1 space-y-3', props.class]"
    @submit.prevent="onSubmit"
  >
    <p class="text-lg">Filters</p>

    <div class="flex flex-col space-y-1">
      <label class="text-sm text-gray-600">Ticker</label>
      <input type="text" v-model="filters.ticker" placeholder="AAPL, MSFT..."
             class="h-9 w-full rounded-md border border-gray-300 px-2 text-sm outline-none focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/40" />
    </div>

    <div class="flex flex-col space-y-1">
      <label class="text-sm text-gray-600">Company</label>
      <input type="text" v-model="filters.company" placeholder="Apple Inc."
             class="h-9 w-full rounded-md border border-gray-300 px-2 text-sm outline-none focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/40" />
    </div>

    <div class="flex flex-col space-y-1">
      <label class="text-sm text-gray-600">Brokerage</label>
      <input type="text" v-model="filters.brokerage" placeholder="Goldman, JPM..."
             class="h-9 w-full rounded-md border border-gray-300 px-2 text-sm outline-none focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/40" />
    </div>

    <div class="flex flex-col space-y-1">
      <label class="text-sm text-gray-600">Target From</label>
      <input
        type="number"
        step="0.01"
        v-model.number="filters.target_from"
        placeholder="Ex: 100.50"
        class="h-9 w-full rounded-md border border-gray-300 px-2 text-sm outline-none
           focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/40"
      />
    </div>

    <div class="flex flex-col space-y-1">
      <label class="text-sm text-gray-600">Target To</label>
      <input
        type="number"
        step="0.01"
        v-model.number="filters.target_to"
        placeholder="Ex: 150.00"
        class="h-9 w-full rounded-md border border-gray-300 px-2 text-sm outline-none
           focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/40"
      />
    </div>

    <div class="flex flex-col space-y-1">
      <label class="text-sm text-gray-600">From Date</label>
      <input type="date" v-model="filters.date_from"
             class="h-9 w-full rounded-md border border-gray-300 px-2 text-sm outline-none focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/40" />
    </div>

    <div class="flex flex-col space-y-1">
      <label class="text-sm text-gray-600">To Date</label>
      <input type="date" v-model="filters.date_to"
             class="h-9 w-full rounded-md border border-gray-300 px-2 text-sm outline-none focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/40" />
    </div>

    <div class="flex gap-2 pt-2">
      <button type="submit"
              class="flex-1 rounded-md border border-gray-900 bg-indigo-900 px-3 py-2 text-sm font-medium text-white transition-all hover:bg-indigo-700 duration-200">
        Search
      </button>
      <button type="button" @click="onReset"
              class="flex-1 rounded-md border border-gray-300 bg-white px-3 py-2 text-sm font-medium text-gray-700 transition hover:bg-gray-50">
        Clear
      </button>
    </div>
  </form>
</template>
