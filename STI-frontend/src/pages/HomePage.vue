<script setup lang="ts">
import StockRecommendations from '@/components/data/StockRecommendations.vue'
import StockTable from '@/components/data/StockTable.vue'
</script>

<template>
  <section class="max-w-[1920px] mx-auto px-2 md:px-5" id="main">
    <div class="flex flex-col min-[1354px]:flex-row">
      <div class="grow mb-5 min-[1354px]:mr-5">
        <div class="w-full glass-table shadow-md rounded-md overflow-hidden">
          <StockTable />
        </div>
      </div>
      <div class="min-w-[250px] bg-card mb-2">
        <StockRecommendations />
      </div>
    </div>
  </section>
</template>

<style scoped></style>
