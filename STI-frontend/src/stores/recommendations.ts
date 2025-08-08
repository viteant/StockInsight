import { defineStore } from 'pinia'
import { onMounted, ref } from 'vue'

export const useRecommendationsStore = defineStore('recommendations', () => {
  const recommendations = ref<StockRecommendation[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const apiUrl = import.meta.env.VITE_API_URL // Asegúrate que esté en tu .env

  async function fetchRecommendations() {
    loading.value = true
    error.value = null
    try {
      const res = await fetch(`${apiUrl}/api/recommendations`)
      if (!res.ok) throw new Error(`Error HTTP: ${res.status}`)
      recommendations.value = await res.json()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Error desconocido'
    } finally {
      loading.value = false
    }
  }

  onMounted(async () => {
    await fetchRecommendations()
  })

  return {
    recommendations,
    loading,
    error
  }
})
