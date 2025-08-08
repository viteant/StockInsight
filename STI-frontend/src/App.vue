<script setup lang="ts">
import NavBar from '@/components/NavBar.vue'
import FooterNav from '@/components/FooterNav.vue'
</script>

<template>
  <NavBar/>
  <div class="pt-[60px] pb-[30px] min-h-screen bg-neutral-50">
    <RouterView/>
  </div>
  <FooterNav/>
</template>

<style scoped></style>
