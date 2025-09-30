import { useState, useEffect } from "react"
import { DollarSign, TrendingUp, CheckCircle } from "lucide-react"
import { api, type Deal } from "@/lib/api"
import { useAuth } from "@/lib/auth"

export function SalesAgentDealStats() {
  const [deals, setDeals] = useState<Deal[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const { user, loading } = useAuth() // Destructure loading

  console.log("SalesAgentDealStats component rendered. User:", user);

  useEffect(() => {
    console.log("useEffect in SalesAgentDealStats running. User:", user);
    if (loading) { // Wait for auth to load
      console.log("Auth is still loading in SalesAgentDealStats.");
      return;
    }
    const fetchDeals = async () => {
      if (!user?.id) {
        setError("User not authenticated.")
        setIsLoading(false)
        return
      }
      console.log("User ID is valid, attempting to fetch deals.");
      try {
        setIsLoading(true)
        const allDeals = await api.getDeals()
        const userDeals = allDeals.filter(deal => deal.created_by && deal.created_by.Int64 === user.id)
        setDeals(userDeals)
      } catch (err) {
        setError("Failed to load deals.")
        console.error("Failed to fetch deals:", err)
      } finally {
        setIsLoading(false)
      }
    }
    fetchDeals()
  }, [user, loading]) // Add loading to dependency array

  const pendingDeals = deals.filter((d) => d.deal_status === "Pending").length
  const wonDeals = deals.filter((d) => d.deal_status === "Closed-Won").length
  const lostDeals = deals.filter((d) => d.deal_status === "Closed-Lost").length
  const totalValueWonDeals = deals
    .filter((d) => d.deal_status === "Closed-Won")
    .reduce((sum, deal) => sum + deal.deal_amount, 0)

  if (isLoading) {
    return <div className="p-4 text-center text-muted-foreground">Loading deal statistics...</div>
  }

  if (error) {
    return <div className="p-4 text-center text-destructive">Error: {error}</div>
  }

  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      <div className="rounded-lg border bg-card p-6">
        <div className="flex items-center gap-2">
          <TrendingUp className="h-5 w-5 text-cyan-600" />
          <h3 className="font-semibold">My Total Deals</h3>
        </div>
        <p className="text-2xl font-bold text-cyan-600 mt-2">{deals.length}</p>
      </div>
      <div className="rounded-lg border bg-card p-6">
        <div className="flex items-center gap-2">
          <DollarSign className="h-5 w-5 text-yellow-600" />
          <h3 className="font-semibold">Pending Deals</h3>
        </div>
        <p className="text-2xl font-bold text-yellow-600 mt-2">{pendingDeals}</p>
      </div>
      <div className="rounded-lg border bg-card p-6">
        <div className="flex items-center gap-2">
          <CheckCircle className="h-5 w-5 text-green-600" />
          <h3 className="font-semibold">Won Deals</h3>
        </div>
        <p className="text-2xl font-bold text-green-600 mt-2">{wonDeals}</p>
      </div>
      <div className="rounded-lg border bg-card p-6">
        <div className="flex items-center gap-2">
          <DollarSign className="h-5 w-5 text-purple-600" />
          <h3 className="font-semibold">Total Value Won</h3>
        </div>
        <p className="text-2xl font-bold text-purple-600 mt-2">${totalValueWonDeals.toLocaleString()}</p>
      </div>
    </div>
  )
}
