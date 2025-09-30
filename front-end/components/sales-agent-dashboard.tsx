"use client"

import React from 'react'
import { SalesAgentDealStats } from "./sales-agent-deal-stats"
import { SalesAgentTasks } from "./sales-agent-tasks"

export default function SalesAgentDashboard() {
  console.log("SalesAgentDashboard component rendered.");
  return (
    <div className="p-4 space-y-6">
      <h1 className="text-3xl font-bold text-balance text-foreground">Sales Agent Dashboard</h1>
      <SalesAgentDealStats />
      <SalesAgentTasks />
    </div>
  )
}
