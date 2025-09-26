"use client"

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  ChartConfig,
  ChartTooltip,
  ChartTooltipContent,
  ChartLegend,
  ChartLegendContent,
} from "@/components/ui/chart";
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, PieChart, Pie, Cell, Legend, ResponsiveContainer, Tooltip, PieLabelRenderProps } from "recharts";
import { Users, DollarSign, RefreshCw, BarChart3, Download, Target, TrendingUp, Calendar } from "lucide-react";
import { Badge } from "@/components/ui/badge";

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
}

interface EmployeeSalesReport {
  employee_id: number
  employee_name: string
  number_of_sales: number
  total_sales_amount: number
}

interface SourceLeadReport {
  lead_date: string
  contact_name: string
  contact_phone: string
  contact_email: string
  lead_source: string
  assigned_employee: string
  lead_status: string
}

interface SourceSalesReport {
  source_name: string
  number_of_sales: number
  total_sales_amount: number
}

interface DealsPipelineReport {
  stage_name: string
  deal_count: number
  total_value: number
  avg_days_in_stage: number
}

interface ReportsChartsProps {
  employeeLeadReport: EmployeeLeadReport | null;
  employeeSalesReport: EmployeeSalesReport[];
  sourceLeadReport: SourceLeadReport[];
  sourceSalesReport: SourceSalesReport[];
  dealsPipelineReport: DealsPipelineReport[];
  hasRole: (roleId: number) => boolean;
  ROLE_RECEPTION: number;
  ROLE_SALES_AGENT: number;
}

export function ReportsCharts({
  employeeLeadReport,
  employeeSalesReport,
  sourceLeadReport,
  sourceSalesReport,
  dealsPipelineReport,
  hasRole,
  ROLE_RECEPTION,
  ROLE_SALES_AGENT,
}: ReportsChartsProps) {

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
  }))

  const salesBySourceData = sourceSalesReport.map((row) => ({
    name: row.source_name,
    value: row.total_sales_amount,
    sales: row.number_of_sales,
  }))

  const pipelineData = dealsPipelineReport.map((row) => ({
    stage: row.stage_name,
    deals: row.deal_count,
    value: row.total_value / 1000,
    avgDays: row.avg_days_in_stage,
  }))

  const COLORS = ["#0891b2", "#f97316", "#10b981", "#8b5cf6", "#ef4444"]

  const employeeLeadChartConfig = {
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
  } satisfies ChartConfig

  const salesBySourceChartConfig = salesBySourceData.reduce((acc, item, index) => {
    acc[item.name] = {
      label: item.name,
      color: COLORS[index % COLORS.length],
    }
    return acc
  }, {} as ChartConfig)

  return (
    <>
      {hasRole(ROLE_RECEPTION) && (
        <>
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
                    <Tooltip />
                    <Bar dataKey="new" fill="#0891b2" name="New Leads" />
                    <Bar dataKey="qualified" fill="#10b981" name="Qualified" />
                    <Bar dataKey="converted" fill="#8b5cf6" name="Converted" />
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
    </>
  );
}
