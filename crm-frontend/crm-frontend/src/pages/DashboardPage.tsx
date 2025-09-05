// Replace the entire contents of src/pages/DashboardPage.tsx with this.

import { useState, useEffect } from 'react';
import { getEmployeeLeadReport, getEmployeeSalesReport, getSourceSalesReport } from '../services/apiService';
import { Title, Paper, Text, SimpleGrid, Loader, Alert } from '@mantine/core';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer, PieChart, Pie, Cell } from 'recharts';
import RecentActivity from '../components/RecentActivity';
import axios from 'axios';

// Define the types for our report data structures
interface EmployeeLeadRow {
  employee_name: string;
  counts: { new: number; contacted: number; qualified: number; converted: number };
}
interface EmployeeLeadReport {
  rows: EmployeeLeadRow[];
}
interface EmployeeSaleRow {
  employee_name: string;
  number_of_sales: number;
  total_sales_amount: number;
}
interface SourceSaleRow {
  source_name: string;
  total_sales_amount: number;
}

const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042', '#AF19FF'];

const DashboardPage = () => {
  // Use 'null' as initial state to clearly distinguish between "loading", "error", and "no data"
  const [leadReport, setLeadReport] = useState<EmployeeLeadRow[] | null>(null);
  const [salesReport, setSalesReport] = useState<EmployeeSaleRow[] | null>(null);
  const [sourceSalesReport, setSourceSalesReport] = useState<SourceSaleRow[] | null>(null);
  const [loading, setLoading] = useState(true);
  const [pageError, setPageError] = useState<string | null>(null); // A general error for the whole page

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      
      try {
        // We use Promise.allSettled to ensure we always get a result for each promise, even if one fails
        const results = await Promise.allSettled([
          getEmployeeLeadReport(),
          getEmployeeSalesReport(),
          getSourceSalesReport(),
        ]);

        // Process lead report result
        if (results[0].status === 'fulfilled') {
          setLeadReport(results[0].value.rows || []);
        } else {
          setLeadReport(null); // Set to null on permission error or other failure
          console.error("Failed to load lead report:", results[0].reason);
        }
        
        // Process employee sales report result
        if (results[1].status === 'fulfilled') {
          setSalesReport(results[1].value || []);
        } else {
          setSalesReport(null);
          console.error("Failed to load employee sales report:", results[1].reason);
        }
        
        // Process source sales report result
        if (results[2].status === 'fulfilled') {
          setSourceSalesReport(results[2].value || []);
        } else {
          setSourceSalesReport(null);
          console.error("Failed to load source sales report:", results[2].reason);
        }
      } catch (err) {
        // This catch block is for truly unexpected errors
        setPageError("An unexpected error occurred while loading dashboard data.");
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, []);

  // Summary stat calculations are now safer, defaulting to 0 if data is null
  const totalLeads = leadReport?.reduce((sum, row) => sum + row.counts.new + row.counts.contacted, 0) || 0;
  const totalSales = salesReport?.reduce((sum, row) => sum + row.total_sales_amount, 0) || 0;
  const totalDeals = salesReport?.reduce((sum, row) => sum + row.number_of_sales, 0) || 0;

  if (loading) return <Loader size="xl" />;
  if (pageError) return <Alert color="red" title="Error">{pageError}</Alert>;

  return (
    <div>
      <Title order={2} mb="xl">Dashboard</Title>

      <SimpleGrid cols={{ base: 1, sm: 3 }} mb="xl">
        <Paper withBorder p="md" radius="md" ta="center">
          <Text size="xl" fw={700}>${totalSales.toLocaleString()}</Text>
          <Text size="sm" c="dimmed">Total Sales Volume</Text>
        </Paper>
        <Paper withBorder p="md" radius="md" ta="center">
          <Text size="xl" fw={700}>{totalDeals}</Text>
          <Text size="sm" c="dimmed">Closed Deals</Text>
        </Paper>
        <Paper withBorder p="md" radius="md" ta="center">
          <Text size="xl" fw={700}>{totalLeads}</Text>
          <Text size="sm" c="dimmed">Active Leads</Text>
        </Paper>
      </SimpleGrid>

      <SimpleGrid cols={{ base: 1, lg: 2 }} spacing="xl">
        {/* --- Conditionally Render Leads Chart --- */}
        {leadReport ? (
          <Paper withBorder p="md" radius="md">
            <Title order={4} mb="md">Leads by Employee</Title>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={leadReport} margin={{ top: 5, right: 20, left: 10, bottom: 5 }}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="employee_name" />
                <YAxis />
                <Tooltip />
                <Legend />
                <Bar dataKey="counts.new" stackId="a" fill="#8884d8" name="New" />
                <Bar dataKey="counts.contacted" stackId="a" fill="#82ca9d" name="Contacted" />
                <Bar dataKey="counts.qualified" stackId="a" fill="#ffc658" name="Qualified" />
                <Bar dataKey="counts.converted" stackId="a" fill="#ff8042" name="Converted" />
              </BarChart>
            </ResponsiveContainer>
          </Paper>
        ) : (
          <Paper withBorder p="md" radius="md" style={{ display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
            <Text c="dimmed">You do not have permission to view the Leads Report.</Text>
          </Paper>
        )}

        {/* --- Conditionally Render Sales by Source Chart --- */}
        {sourceSalesReport ? (
          <Paper withBorder p="md" radius="md">
            <Title order={4} mb="md">Sales by Lead Source</Title>
            <ResponsiveContainer width="100%" height={300}>
                <PieChart>
                    <Pie
                        data={sourceSalesReport}
                        cx="50%"
                        cy="50%"
                        labelLine={false}
                        outerRadius={110}
                        fill="#8884d8"
                        dataKey="total_sales_amount"
                        nameKey="source_name"
                        label={({ name, percent }) => (percent ? `${name} ${(percent * 100).toFixed(0)}%` : name)}
                    >
                        {sourceSalesReport.map((entry, index) => (
                            <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                        ))}
                    </Pie>
                    <Tooltip formatter={(value: number) => `$${value.toLocaleString()}`} />
                    <Legend />
                </PieChart>
            </ResponsiveContainer>
          </Paper>
        ) : (
          <Paper withBorder p="md" radius="md" style={{ display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
            <Text c="dimmed">You do not have permission to view the Sales by Source Report.</Text>
          </Paper>
        )}
      </SimpleGrid>
      
      <Paper withBorder p="md" radius="md" mt="xl">
        <Title order={4} mb="md">Recent Activity</Title>
        <RecentActivity />
      </Paper>
    </div>
  );
};

export default DashboardPage;