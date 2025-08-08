<script setup lang="ts">
import NavBar from '@/components/NavBar.vue'
import FooterNav from '@/components/FooterNav.vue'
</script>

<template>
  <NavBar/>
  <div class="">
    <RouterView/>
  </div>
  <FooterNav/>
</template>

<style scoped></style>
