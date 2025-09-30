"use client"

import React from 'react'

export default function ReceptionDashboard() {
  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold mb-4">Reception Dashboard</h1>
      <p>Welcome, Receptionist! Here you can manage contacts, properties, and general inquiries.</p>
      {/* Add reception specific widgets and components here */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mt-4">
        <div className="bg-white p-4 rounded-lg shadow">
          <h2 className="text-lg font-semibold">Manage Contacts</h2>
          <p>View, add, and update customer contacts.</p>
        </div>
        <div className="bg-white p-4 rounded-lg shadow">
          <h2 className="text-lg font-semibold">View Properties</h2>
          <p>Browse available properties.</p>
        </div>
        <div className="bg-white p-4 rounded-lg shadow">
          <h2 className="text-lg font-semibold">Communication Logs</h2>
          <p>Review all communication records.</p>
        </div>
      </div>
    </div>
  )
}
