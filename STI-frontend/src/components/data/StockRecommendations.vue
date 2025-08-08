<script setup lang="ts">
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Icon } from '@iconify/vue'
import { useRecommendationsStore } from '@/stores/recommendations'
import { computed, onMounted } from 'vue'
import RecommendationCard from '@/components/data/RecommendationCard.vue'
import { useStockStore } from '@/stores/stock'

const recommendationsStore = useRecommendationsStore()
const stockStore = useStockStore()
const stockRecommendations = computed(() => recommendationsStore.recommendations??[])



const filterByID=(id:string) => {
  stockStore.setLoading(true)
  stockStore.setFilters({...stockStore.filters, id: id})
}

</script>

<template>
  <Tabs default-value="buy" class="w-full not-md:mb-2">
    <TabsList class="w-full dark">
      <TabsTrigger value="buy"> Buy <Icon icon="gg:trending-down" /> </TabsTrigger>
      <TabsTrigger value="sell"> Sell <Icon icon="gg:trending" /> </TabsTrigger>
    </TabsList>
    <TabsContent value="buy">
      <div class="px-2 rounded-2xl flex flex-row min-[1354px]:flex-col gap-x-1 flex-wrap w-full justify-between">
        <p class="text-white text-xl mb-1 px-2 w-full">Best to Invest</p>
        <RecommendationCard
          :order="index"
        v-if="stockRecommendations.length > 0"
        v-for="(value, index) in stockRecommendations.filter(
          (r: StockRecommendation) => r.normalize_rating_to == 'buy',
        )"
          @getTickerID="filterByID"
        :stockRecommendation="value"
        :key="index"
      />
      </div>
    </TabsContent>
    <TabsContent value="sell">
      <div class="p-2 rounded-2xl flex flex-row min-[1354px]:flex-col flex-wrap w-full justify-between">
        <p class="text-white text-xl mb-1 px-2 w-full">Best to Sell</p>
        <RecommendationCard
          :order="index"
        v-if="stockRecommendations.length > 0"
        v-for="(recommendation, index) in stockRecommendations.filter(
          (r: StockRecommendation) => r.normalize_rating_to == 'sell',
        )"
          @getTickerID="filterByID"
        :stockRecommendation="recommendation"
        :key="index"
      />
      </div>
    </TabsContent>
  </Tabs>
</template>

<style scoped></style>
