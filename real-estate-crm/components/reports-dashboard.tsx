"use client"

import { useState, useEffect } from "react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, PieChart, Pie, Cell } from "recharts"
import { TrendingUp, Users, DollarSign, Target, BarChart3, Download, Calendar, RefreshCw } from "lucide-react"
import { api } from "@/lib/api"

interface EmployeeLeadReport {
  rows: Array<{
    employee_id: number
    employee_name: string
    counts: {
      new: number
      contacted: number
      qualified: number
      converted: number
      lost: number
    }
  }>
  total: {
    new: number
    contacted: number
    qualified: number
    converted: number
    lost: number
  }
  summary: {
    total_employees: number
    avg_conversion_rate: number
    top_performer: string
  }
}

interface EmployeeSalesReport {
  employee_name: string
  number_of_sales: number
  total_sales_amount: number
  avg_deal_size: number
  conversion_rate: number
}

interface SourceLeadReport {
  lead_date: string
  contact_name: string
  contact_phone: string
  contact_email: string
  lead_source: string
  assigned_employee: string
  lead_status: string
  deal_value?: number
}

interface SourceSalesReport {
  source_name: string
  number_of_sales: number
  total_sales_amount: number
  avg_deal_size: number
  lead_count: number
  conversion_rate: number
}

interface DealsPipelineReport {
  stage_name: string
  deal_count: number
  total_value: number
  avg_days_in_stage: number
}

export function ReportsDashboard() {
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState("")
  const [selectedPeriod, setSelectedPeriod] = useState("month")
  const [refreshing, setRefreshing] = useState(false)

  const [employeeLeadReport, setEmployeeLeadReport] = useState<EmployeeLeadReport | null>(null)
  const [employeeSalesReport, setEmployeeSalesReport] = useState<EmployeeSalesReport[]>([])
  const [sourceLeadReport, setSourceLeadReport] = useState<SourceLeadReport[]>([])
  const [sourceSalesReport, setSourceSalesReport] = useState<SourceSalesReport[]>([])
  const [dealsPipelineReport, setDealsPipelineReport] = useState<DealsPipelineReport[]>([])

  useEffect(() => {
    loadReportsData()
  }, [selectedPeriod])

  const loadReportsData = async () => {
    try {
      setIsLoading(true)
      setError("")

      // Load all reports in parallel using new structured endpoints
      const [employeeLeadData, employeeSalesData, sourceLeadData, sourceSalesData, pipelineData] = await Promise.all([
        api.get(`/reports/employee-leads?period=${selectedPeriod}`),
        api.get(`/reports/employee-sales?period=${selectedPeriod}`),
        api.get(`/reports/source-leads?period=${selectedPeriod}`),
        api.get(`/reports/source-sales?period=${selectedPeriod}`),
        api.get(`/reports/deals-pipeline?period=${selectedPeriod}`),
      ])

      setEmployeeLeadReport(employeeLeadData)
      setEmployeeSalesReport(employeeSalesData)
      setSourceLeadReport(sourceLeadData)
      setSourceSalesReport(sourceSalesData)
      setDealsPipelineReport(pipelineData)
    } catch (err) {
      console.error("Failed to load reports:", err)
      setError("Failed to load reports data. Please try again.")

      setEmployeeLeadReport({
        rows: [
          {
            employee_id: 1,
            employee_name: "John Smith",
            counts: { new: 15, contacted: 12, qualified: 8, converted: 5, lost: 2 },
          },
          {
            employee_id: 2,
            employee_name: "Sarah Johnson",
            counts: { new: 20, contacted: 18, qualified: 12, converted: 8, lost: 4 },
          },
          {
            employee_id: 3,
            employee_name: "Mike Brown",
            counts: { new: 10, contacted: 8, qualified: 6, converted: 3, lost: 3 },
          },
        ],
        total: { new: 45, contacted: 38, qualified: 26, converted: 16, lost: 9 },
        summary: { total_employees: 3, avg_conversion_rate: 35.6, top_performer: "Sarah Johnson" },
      })

      setEmployeeSalesReport([
        {
          employee_name: "John Smith",
          number_of_sales: 5,
          total_sales_amount: 1250000,
          avg_deal_size: 250000,
          conversion_rate: 33.3,
        },
        {
          employee_name: "Sarah Johnson",
          number_of_sales: 8,
          total_sales_amount: 2100000,
          avg_deal_size: 262500,
          conversion_rate: 40.0,
        },
        {
          employee_name: "Mike Brown",
          number_of_sales: 3,
          total_sales_amount: 750000,
          avg_deal_size: 250000,
          conversion_rate: 30.0,
        },
      ])
    } finally {
      setIsLoading(false)
      setRefreshing(false)
    }
  }

  const handleRefresh = async () => {
    setRefreshing(true)
    await loadReportsData()
  }

  const getStatusBadge = (status: string) => {
    const statusColors: { [key: string]: string } = {
      New: "bg-blue-50 text-blue-700 border-blue-200",
      Contacted: "bg-yellow-50 text-yellow-700 border-yellow-200",
      Qualified: "bg-green-50 text-green-700 border-green-200",
      Converted: "bg-purple-50 text-purple-700 border-purple-200",
      Lost: "bg-red-50 text-red-700 border-red-200",
    }

    return (
      <Badge variant="outline" className={statusColors[status] || "bg-gray-50 text-gray-700 border-gray-200"}>
        {status}
      </Badge>
    )
  }

  const leadsByEmployeeData =
    employeeLeadReport?.rows.map((row) => ({
      name: row.employee_name,
      new: row.counts.new,
      qualified: row.counts.qualified,
      converted: row.counts.converted,
      total: row.counts.new + row.counts.contacted + row.counts.qualified + row.counts.converted + row.counts.lost,
    })) || []

  const salesByEmployeeData = employeeSalesReport.map((row) => ({
    name: row.employee_name,
    sales: row.number_of_sales,
    amount: row.total_sales_amount / 1000, // Convert to thousands
    avgDeal: row.avg_deal_size / 1000,
    conversionRate: row.conversion_rate,
  }))

  const salesBySourceData = sourceSalesReport.map((row) => ({
    name: row.source_name,
    value: row.total_sales_amount,
    sales: row.number_of_sales,
    leads: row.lead_count,
    conversionRate: row.conversion_rate,
  }))

  const pipelineData = dealsPipelineReport.map((row) => ({
    stage: row.stage_name,
    deals: row.deal_count,
    value: row.total_value / 1000,
    avgDays: row.avg_days_in_stage,
  }))

  const COLORS = ["#0891b2", "#f97316", "#10b981", "#8b5cf6", "#ef4444"]

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-muted-foreground">Loading reports...</div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="space-y-1">
          <h1 className="text-3xl font-bold text-balance text-foreground">Reports Dashboard</h1>
          <p className="text-muted-foreground text-pretty">
            Analyze your business performance with comprehensive reports and insights
          </p>
        </div>
        <div className="flex items-center gap-4">
          <Select value={selectedPeriod} onValueChange={setSelectedPeriod}>
            <SelectTrigger className="w-40">
              <SelectValue placeholder="Select period" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="week">This Week</SelectItem>
              <SelectItem value="month">This Month</SelectItem>
              <SelectItem value="quarter">This Quarter</SelectItem>
              <SelectItem value="year">This Year</SelectItem>
            </SelectContent>
          </Select>
          <Button variant="outline" onClick={handleRefresh} disabled={refreshing} className="gap-2 bg-transparent">
            <RefreshCw className={`h-4 w-4 ${refreshing ? "animate-spin" : ""}`} />
            Refresh
          </Button>
          <Button variant="outline" className="gap-2 bg-transparent">
            <Download className="h-4 w-4" />
            Export
          </Button>
        </div>
      </div>

      {error && (
        <Alert variant="destructive">
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      {/* Key Metrics */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Leads</CardTitle>
            <Users className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-cyan-600">{employeeLeadReport?.total.new || 0}</div>
            <p className="text-xs text-muted-foreground">
              {employeeLeadReport?.summary.total_employees || 0} employees active
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Conversions</CardTitle>
            <Target className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">{employeeLeadReport?.total.converted || 0}</div>
            <p className="text-xs text-muted-foreground">
              {employeeLeadReport?.summary.avg_conversion_rate.toFixed(1) || 0}% avg conversion rate
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Sales</CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-purple-600">
              ${employeeSalesReport.reduce((sum, emp) => sum + emp.total_sales_amount, 0).toLocaleString()}
            </div>
            <p className="text-xs text-muted-foreground">
              $
              {(
                employeeSalesReport.reduce((sum, emp) => sum + emp.avg_deal_size, 0) /
                Math.max(employeeSalesReport.length, 1)
              ).toLocaleString()}{" "}
              avg deal size
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Top Performer</CardTitle>
            <TrendingUp className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-orange-600">
              {employeeLeadReport?.summary.top_performer || "N/A"}
            </div>
            <p className="text-xs text-muted-foreground">Leading in conversions</p>
          </CardContent>
        </Card>
      </div>

      {/* Charts Section */}
      <div className="grid gap-6 md:grid-cols-2">
        {/* Employee Lead Performance */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <BarChart3 className="h-5 w-5" />
              Employee Lead Performance
            </CardTitle>
            <CardDescription>Lead conversion by team member</CardDescription>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={leadsByEmployeeData}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="name" />
                <YAxis />
                <Tooltip />
                <Bar dataKey="new" fill="#0891b2" name="New Leads" />
                <Bar dataKey="qualified" fill="#10b981" name="Qualified" />
                <Bar dataKey="converted" fill="#8b5cf6" name="Converted" />
              </BarChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>

        {/* Sales by Source */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Target className="h-5 w-5" />
              Sales by Source
            </CardTitle>
            <CardDescription>Revenue distribution by lead source</CardDescription>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={salesBySourceData}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                  outerRadius={80}
                  fill="#8884d8"
                  dataKey="value"
                >
                  {salesBySourceData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip formatter={(value: number) => [`$${value.toLocaleString()}`, "Revenue"]} />
              </PieChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>
      </div>

      {pipelineData.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <TrendingUp className="h-5 w-5" />
              Deals Pipeline
            </CardTitle>
            <CardDescription>Deal distribution across pipeline stages</CardDescription>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={pipelineData}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="stage" />
                <YAxis yAxisId="left" />
                <YAxis yAxisId="right" orientation="right" />
                <Tooltip />
                <Bar yAxisId="left" dataKey="deals" fill="#0891b2" name="Deal Count" />
                <Bar yAxisId="right" dataKey="value" fill="#f97316" name="Value (K)" />
              </BarChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>
      )}

      {/* Employee Sales Performance */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <DollarSign className="h-5 w-5" />
            Employee Sales Performance
          </CardTitle>
          <CardDescription>Sales volume and revenue by team member</CardDescription>
        </CardHeader>
        <CardContent>
          <ResponsiveContainer width="100%" height={300}>
            <BarChart data={salesByEmployeeData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="name" />
              <YAxis yAxisId="left" />
              <YAxis yAxisId="right" orientation="right" />
              <Tooltip />
              <Bar yAxisId="left" dataKey="sales" fill="#0891b2" name="Number of Sales" />
              <Bar yAxisId="right" dataKey="amount" fill="#f97316" name="Revenue (K)" />
            </BarChart>
          </ResponsiveContainer>
        </CardContent>
      </Card>

      {/* Reports Tables */}
      <div className="grid gap-6 md:grid-cols-2">
        {/* Employee Lead Report */}
        <Card>
          <CardHeader>
            <CardTitle>Employee Lead Report</CardTitle>
            <CardDescription>Detailed lead breakdown by employee</CardDescription>
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Employee</TableHead>
                  <TableHead>New</TableHead>
                  <TableHead>Qualified</TableHead>
                  <TableHead>Converted</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {employeeLeadReport?.rows.map((row) => (
                  <TableRow key={row.employee_id}>
                    <TableCell className="font-medium">{row.employee_name}</TableCell>
                    <TableCell>{row.counts.new}</TableCell>
                    <TableCell>{row.counts.qualified}</TableCell>
                    <TableCell>{row.counts.converted}</TableCell>
                  </TableRow>
                )) || []}
              </TableBody>
            </Table>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Source Sales Report</CardTitle>
            <CardDescription>Sales performance by lead source</CardDescription>
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Source</TableHead>
                  <TableHead>Sales</TableHead>
                  <TableHead>Revenue</TableHead>
                  <TableHead>Conv. Rate</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {sourceSalesReport.map((row, index) => (
                  <TableRow key={index}>
                    <TableCell className="font-medium">{row.source_name}</TableCell>
                    <TableCell>{row.number_of_sales}</TableCell>
                    <TableCell>${row.total_sales_amount.toLocaleString()}</TableCell>
                    <TableCell>{row.conversion_rate.toFixed(1)}%</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </div>

      {/* Recent Lead Activity */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Calendar className="h-5 w-5" />
            Recent Lead Activity
          </CardTitle>
          <CardDescription>Latest lead interactions and status updates</CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Date</TableHead>
                <TableHead>Contact</TableHead>
                <TableHead>Source</TableHead>
                <TableHead>Assigned To</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Deal Value</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {sourceLeadReport.map((row, index) => (
                <TableRow key={index}>
                  <TableCell>{new Date(row.lead_date).toLocaleDateString()}</TableCell>
                  <TableCell className="font-medium">
                    <div>
                      <div>{row.contact_name}</div>
                      <div className="text-sm text-muted-foreground">{row.contact_email}</div>
                    </div>
                  </TableCell>
                  <TableCell>{row.lead_source}</TableCell>
                  <TableCell>{row.assigned_employee}</TableCell>
                  <TableCell>{getStatusBadge(row.lead_status)}</TableCell>
                  <TableCell>{row.deal_value ? `$${row.deal_value.toLocaleString()}` : "N/A"}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  )
}
