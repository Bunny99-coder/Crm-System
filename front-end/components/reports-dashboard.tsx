"use client"

import { cn } from "@/lib/utils";
import {
  Table,
  TableHeader,
  TableRow,
  TableHead,
  TableBody,
  TableCell,
  TableCaption,
} from "@/components/ui/table";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  ChartConfig,
  ChartTooltip,
  ChartTooltipContent,
  ChartLegend,
  ChartLegendContent,
} from "@/components/ui/chart";
import { Users, DollarSign, RefreshCw, BarChart3, Download, Target, TrendingUp, Calendar } from "lucide-react";

import { DateRangePicker } from "@/components/ui/date-range-picker";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useAuth, ROLE_RECEPTION, ROLE_SALES_AGENT } from "@/lib/auth";

import { addDays, format, subDays } from "date-fns";
import dynamic from "next/dynamic";
import { useState, useEffect } from "react";
import { api, DealsPipelineReportRow, EmployeeLeadReport, EmployeeSalesReportRow, SourceLeadReportRow, SourceSalesReportRow } from "@/lib/api";
import { DateRange } from "react-day-picker";
import { z } from "zod";
import { ResponsiveContainer, BarChart, CartesianGrid, XAxis, YAxis, Tooltip, Bar, PieChart, Pie, Cell, PieLabelRenderProps } from "recharts";

const Chart = dynamic(() => import("react-apexcharts"), { ssr: false });

const ChartContainer = dynamic(() => import("@/components/ui/chart").then((mod) => mod.ChartContainer), {
  ssr: false,
})

const ReportsCharts = dynamic(() => import("@/components/reports-charts").then((mod) => mod.ReportsCharts), {
  ssr: false,
});











export function ReportsDashboard() {
  const [error, setError] = useState("")
  const [selectedPeriod, setSelectedPeriod] = useState("month")
  const [refreshing, setRefreshing] = useState(false)

  const [employeeLeadReport, setEmployeeLeadReport] = useState<EmployeeLeadReport | null>(null)
  const [employeeSalesReport, setEmployeeSalesReport] = useState<EmployeeSalesReportRow[]>([])
  const [sourceLeadReport, setSourceLeadReport] = useState<SourceLeadReportRow[]>([])
  const [sourceSalesReport, setSourceSalesReport] = useState<SourceSalesReportRow[]>([])
  const [dealsPipelineReport, setDealsPipelineReport] = useState<DealsPipelineReportRow[]>([])

  const { user, hasRole, loading: isAuthLoading } = useAuth() // Destructure loading from useAuth

  useEffect(() => {
    if (!isAuthLoading) { // Only load data once authentication state is known
      loadReportsData()
    }
  }, [selectedPeriod, isAuthLoading]) // Add isAuthLoading to dependency array

  const loadReportsData = async () => {
    try {
      // setIsLoading(true) // No longer needed here, controlled by isAuthLoading
      setError("")

      if (hasRole(ROLE_RECEPTION)) {
        // Load all reports in parallel using new structured endpoints
        const [employeeLeadData, employeeSalesData, sourceLeadData, sourceSalesData, pipelineData] = await Promise.all([
          api.getEmployeeLeadReport(),
          api.getEmployeeSalesReport(),
          api.getSourceLeadReport(),
          api.getSourceSalesReport(),
          api.getDealsPipelineReport(selectedPeriod),
        ])

        setEmployeeLeadReport(employeeLeadData)
        setEmployeeSalesReport(employeeSalesData)
        setSourceLeadReport(sourceLeadData)
        setSourceSalesReport(sourceSalesData)
        setDealsPipelineReport(pipelineData)
      } else if (hasRole(ROLE_SALES_AGENT)) {
        const [mySalesData, pipelineData] = await Promise.all([
          api.getMySalesReport(),
          api.getDealsPipelineReport(selectedPeriod),
        ])
        setEmployeeSalesReport(mySalesData)
        setDealsPipelineReport(pipelineData)
      }
    } catch (err) {
      console.error("Failed to load reports:", err)
      setError("Failed to load reports data. Please try again.")

      // Mock data for demonstration (updated to match new interfaces)
      if (hasRole(ROLE_RECEPTION)) {
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
        })

        setEmployeeSalesReport([
          {
            employee_id: 1,
            employee_name: "John Smith",
            number_of_sales: 5,
            total_sales_amount: 1250000,
          },
          {
            employee_id: 2,
            employee_name: "Sarah Johnson",
            number_of_sales: 8,
            total_sales_amount: 2100000,
          },
          {
            employee_id: 3,
            employee_name: "Mike Brown",
            number_of_sales: 3,
            total_sales_amount: 750000,
          },
        ])
      }
    } finally {
      setRefreshing(false)
    }
  }

  const handleRefresh = async () => {
    setRefreshing(true)
    await loadReportsData()
  }

  const leadsByEmployeeData = (employeeLeadReport?.rows || []).map((row) => ({
    name: row.employee_name,
    new: row.counts.new,
    qualified: row.counts.qualified,
    converted: row.counts.converted,
  }))

  const salesBySourceData = sourceSalesReport
    .filter(row => row.source_name)
    .map((row) => ({
      name: row.source_name,
      value: row.total_sales_amount,
    }))

  const COLORS = ["#0088FE", "#00C49F", "#FFBB28", "#FF8042"];

const getStatusBadge = (status: string) => {
  switch (status.toLowerCase()) {
    case 'new':
      return <Badge variant="secondary">New</Badge>;
    case 'contacted':
      return <Badge variant="outline">Contacted</Badge>;
    case 'qualified':
      return <Badge className="bg-blue-500 text-white">Qualified</Badge>;
    case 'converted':
      return <Badge className="bg-green-500 text-white">Converted</Badge>;
    case 'lost':
      return <Badge variant="destructive">Lost</Badge>;
    default:
      return <Badge>{status}</Badge>;
  }
};

const employeeLeadChartConfig: ChartConfig = {
  new: {
    label: "New Leads",
    color: "#0891b2",
  },
  qualified: {
    label: "Qualified",
    color: "#10b981",
  },
  converted: {
    label: "Converted",
    color: "#8b5cf6",
  },
};





  const salesByEmployeeData = employeeSalesReport.map((row) => ({
    name: row.employee_name,
    sales: row.number_of_sales,
    amount: row.total_sales_amount / 1000, // Convert to thousands
  }))



  const pipelineData = (dealsPipelineReport?.rows || []).map((row) => ({
    stage: row.stage_name,
    deals: row.deal_count,
    value: (row.total_amount ?? 0) / 1000,
    avgDays: row.avg_days_in_stage,
  }))





  const salesBySourceChartConfig = salesBySourceData.reduce((acc, item, index) => {
    acc[item.name] = {
      label: item.name,
      color: COLORS[index % COLORS.length],
    }
    return acc
  }, {} as ChartConfig)

  if (isAuthLoading) {
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

      {hasRole(ROLE_RECEPTION) && (
        <>
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
                  {/* Removed summary usage */}
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
                  {/* Removed summary usage */}
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
                  {/* Removed avg deal size usage */}
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
                  {/* Removed summary usage */}
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
              
                <ChartContainer config={employeeLeadChartConfig} className="aspect-video h-[300px]">
                  <BarChart data={leadsByEmployeeData}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="name" />
                    <YAxis />
                    <Tooltip content={<ChartTooltipContent />} />
                    <Bar dataKey="new" fill="var(--color-new)" radius={4} />
                    <Bar dataKey="qualified" fill="var(--color-qualified)" radius={4} />
                    <Bar dataKey="converted" fill="var(--color-converted)" radius={4} />
                    <ChartLegend />
                  </BarChart>
                </ChartContainer>
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
                <ChartContainer config={salesBySourceChartConfig} className="aspect-video h-[300px]">
                  <PieChart>
                    <Pie
                      data={salesBySourceData}
                      cx="50%"
                      cy="50%"
                      labelLine={false}
                      label={(props: PieLabelRenderProps) => {
                        const { name, value } = props.payload as { name: string; value: number };
                        const totalSales = salesBySourceData.reduce((sum, item) => sum + item.value, 0);
                        const calculatedPercent = totalSales > 0 ? (value / totalSales) * 100 : 0;
                        return `${name} ${calculatedPercent.toFixed(0)}%`;
                      }}
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
                </ChartContainer>
              </CardContent>
            </Card>
          </div>
        </>
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

      {hasRole(ROLE_RECEPTION) && (
        <>
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
                      {/* Removed Conv. Rate as it's not available from API */}
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {sourceSalesReport.map((row, index) => (
                      <TableRow key={index}>
                        <TableCell className="font-medium">{row.source_name}</TableCell>
                        <TableCell>{row.number_of_sales}</TableCell>
                        <TableCell>${row.total_sales_amount.toLocaleString()}</TableCell>
                        {/* Removed conversion_rate usage */}
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
                    {/* Removed Deal Value as it's not available from API */}
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
                      {/* Removed deal_value usage */}
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </CardContent>
          </Card>
        </>
      )}
    </div>
  )
}