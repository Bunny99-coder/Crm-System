"use client"

import React from 'react'

export default function SalesAgentDashboard() {
  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold mb-4">Sales Agent Dashboard</h1>
      <p>Welcome, Sales Agent! Here you can manage your leads, deals, and tasks.</p>
      {/* Add sales agent specific widgets and components here */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mt-4">
        <div className="bg-white p-4 rounded-lg shadow">
          <h2 className="text-lg font-semibold">My Leads</h2>
          <p>View and manage your assigned leads.</p>
        </div>
        <div className="bg-white p-4 rounded-lg shadow">
          <h2 className="text-lg font-semibold">My Deals</h2>
          <p>Track the progress of your deals.</p>
        </div>
        <div className="bg-white p-4 rounded-lg shadow">
          <h2 className="text-lg font-semibold">My Tasks</h2>
          <p>See your upcoming tasks and activities.</p>
        </div>
      </div>
    </div>
  )
}
