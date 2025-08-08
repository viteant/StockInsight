<script setup lang="ts">
import StockDataTable from '@/components/StockDataTable.vue'
</script>

<template>
  <section class="container mx-auto px-5" id="main">
    <div class="grid grid-cols-12 mt-5 gap-5">
      <div class="col-span-10 px-5">
        <div class="w-full bg-white shadow-md rounded-md  overflow-hidden">
          <div class="bg-[#5a5bf3] p-2">
            <h2 class="text-white text-xl">Stocks</h2>
          </div>

          <StockDataTable/>
        </div>
      </div>
      <div class="col-span-2">
        <h1>Recomendaciones</h1>
      </div>
    </div>
  </section>
</template>

<style scoped></style>
