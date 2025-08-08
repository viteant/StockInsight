export {}

declare global {
  type Stock = {
    id: string
    ticker: string
    company: string
    brokerage: string
    action: string
    rating_from: string
    rating_to: string
    normalize_rating_from: string
    normalize_rating_to: string
    target_from: number
    target_to: number
    created_at: string // ISO date string
  }

  type StockFilters = {
    id?:string
    page?: number | null
    limit?: number | null
    orderBy?: string | null
    orderDir?: "asc" | "desc"
    ticker?: string
    company?: string
    brokerage?: string
    target_from?: number | null
    target_to?: number | null
    date_from?: string
    date_to?: string
  }

  type StockResponse = {
    items: Stock[]
    limit: number
    page: number
    total: number
    total_pages: number
  }

  type StockRecommendation = {
    id: string
    ticker: string
    company: string
    brokerage: string
    action: string
    target_from: number
    target_to: number
    normalize_rating_from: string
    normalize_rating_to: string
    weight_score: number
  }


}
