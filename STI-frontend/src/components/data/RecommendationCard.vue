<script setup lang="ts">
import { computed } from 'vue'
import {Icon} from '@iconify/vue'

const props = defineProps<{
  stockRecommendation: StockRecommendation,
  order: number
}>()

const change_value = computed(()=> (props.stockRecommendation.target_to - props.stockRecommendation.target_from))
const change_percent = computed(()=>{
  if (props.stockRecommendation.target_from === 0) return 0
  return ( change_value.value / props.stockRecommendation.target_from) * 100
})

defineEmits(['getTickerID'])
</script>

<template>
  <div :class="[`${change_value<0?'glass-red':(change_value>0?'glass-green':'glass-yellow')} ${change_value<0?'text-red-700':(change_value>0?'text-green-700':'text-amber-600')}`,'mb-2 p-2 rounded-lg mx-auto w-full max-w-[230px]'] ">
    <div class="flex justify-between">
      <div class="flex items-center h-full w-[60px] text-left">
        <p class="text-3xl">#{{order+1}}</p>
      </div>
      <div class="grow">
        <div class="flex items-center justify-between">
          <p class="font-bold cursor-pointer underline" @click="$emit('getTickerID',stockRecommendation.id)">{{stockRecommendation.ticker}}</p>
          <p>$ {{stockRecommendation?.target_to?.toFixed(2) }} </p>
        </div>
        <div :class="['flex items-center justify-between']">
          <Icon :icon="`mdi:${change_value<0?'arrow-down':(change_value>0?'arrow-up':'arrow-vertical-collapse')}`" class="mr-1"/>
          <p class="mr-1">% {{ Math.abs(change_percent)?.toFixed(2) }}</p>
          <p>($ {{ change_value?.toFixed(2) }})</p>
        </div>
      </div>
    </div>
  </div>

</template>

<style scoped></style>
