// src/stores/stock.ts
import { defineStore } from 'pinia'
import { ref, computed, watch, onMounted } from 'vue'

const apiUrl = import.meta.env.VITE_API_URL

// Si los tipos son globales según lo que definimos antes:
// type Stock
// type StockFilters
// type StockResponse

function debounce<T extends (...args: any[]) => any>(fn: T, wait = 300) {
  let t: number | undefined
  return (...args: Parameters<T>) => {
    window.clearTimeout(t)
    t = window.setTimeout(() => fn(...args), wait)
  }
}

export const useStockStore = defineStore('stock', () => {
  // DATA
  const stocks = ref<Stock[]>([])
  const retries=ref(0)
  const filters = ref<StockFilters>({
    id:undefined,
    page: 1,
    limit: 20,
    orderBy: '',
    orderDir: undefined,
    ticker: '',
    company: '',
    brokerage: '',
    target_from: undefined as number | undefined,
    target_to: undefined as number | undefined,
    date_from: '',
    date_to: ''
  })

  const getFilters = computed(() => filters.value??{})

  // META

  const loading = ref(false)
  const error = ref<string | null>(null)
  const total = ref(0)
  const totalPages = ref(0)
  // URL de depuración

  const lastURL = ref('')

  const getLoading = computed(() => {
    return loading.value
  })
  const page = computed(() => filters.value.page ?? 1)
  const limit = computed(() => filters.value.limit ?? 20)
  const getStocks=computed( () => {
    return stocks.value
  })

  function buildQuery(q: StockFilters) {
    const entries: [string, string][] = []
    const push = (k: string, v: unknown) => {
      if (v === '' || v === null || v === undefined) return
      entries.push([k, String(v)])
    }
    push('page', q.page)
    push('limit', q.limit)
    push('orderBy', q.orderBy)
    push('orderDir', q.orderDir)
    push('id', q.id)
    push('ticker', q.ticker?.toString().trim())
    push('company', q.company?.toString().trim())
    push('brokerage', q.brokerage?.toString().trim())
    push('target_from_min', q.target_from)
    push('target_from_max', q.target_to)
    push('date_from', q.date_from?.toString().trim())
    push('date_to', q.date_to?.toString().trim())
    return entries.map(([k, v]) => `${encodeURIComponent(k)}=${encodeURIComponent(v)}`).join('&')
  }

  async function fetchStocks({ append = false }: { append?: boolean } = {}) {
    loading.value = true
    error.value = null
    try {
      const url = `${apiUrl}/api/stocks?${buildQuery(filters.value)}`
      lastURL.value = url
      const res = await fetch(url, { headers: { Accept: 'application/json' } })
      if (!res.ok) throw new Error(`HTTP ${res.status}`)
      const data: StockResponse = await res.json()

      total.value = data.total ?? 0
      totalPages.value = data.total_pages ?? 0

      stocks.value = append ? [...stocks.value, ...(data.items || [])] : (data.items || [])
      setTimeout(()=>{
        loading.value = false
      }, 800)

    } catch (e:any) {
      error.value = e?.message ?? 'Error desconocido'
      if(retries.value<3){
        retries.value++
        setTimeout(()=>{
          fetchStocks({ append: false })
        },500)
      }else{
        retries.value = 0
        loading.value = false
      }
    }
  }

  // Helpers de filtros
  function setLoading(value: boolean) {
    loading.value = value
  }

  function setFilter<K extends keyof StockFilters>(key: K, value: StockFilters[K]) {
    // Cada vez que cambian filtros, reseteamos a la primera página
    if (key !== 'page') filters.value.page = 1;
    (filters.value as any)[key] = value as any
  }

  function setFilters(partial: Partial<StockFilters>) {
    filters.value = { ...filters.value, ...partial, page: 1 }

  }

  function resetFilters() {
    stocks.value = []
    filters.value = {
      id:undefined,
      page: 1,
      limit: 20,
      ticker: '',
      company: '',
      brokerage: '',
      target_from: undefined,
      target_to: undefined,
      date_from: '',
      date_to: '',
    }
  }

  // Paginación
  function setPage(n: number) {
    filters.value.page = Math.max(1, n)
  }
  function nextPage() {
    if (page.value < totalPages.value) filters.value.page = page.value + 1
  }
  function prevPage() {
    if (page.value > 1) filters.value.page = page.value - 1
  }
  function setLimit(n: number) {
    filters.value.limit = Math.max(1, n)
    filters.value.page = 1
  }

  // Recarga y carga incremental
  function refresh() {
    return fetchStocks({ append: false })
  }

  async function loadMore() {
    if (page.value >= totalPages.value) return
    setPage(page.value + 1)
    // Como hay watch con debounce, forzamos fetch inmediato:
    await fetchStocks({ append: true })
  }

  // Auto-fetch cuando cambian los filtros (debounced)
  const debouncedFetch = debounce( () => fetchStocks({ append: false }), 500)
  watch(filters, () => {
    debouncedFetch()
  }, { deep: true, immediate: true })

  return {
    // state
    stocks, filters, error, total, totalPages, lastURL,
    // getters
    page, limit, getStocks, getLoading, getFilters,
    // actions
    fetchStocks, refresh, loadMore,
    setFilter, setFilters, resetFilters,
    setPage, nextPage, prevPage, setLimit, setLoading
  }
})
