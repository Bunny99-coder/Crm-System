import { BarChart, Bar, XAxis, YAxis, CartesianGrid, PieChart, Pie, Cell, Legend, ResponsiveContainer, Tooltip, PieLabelRenderProps } from "recharts";
import { ChartConfig, ChartTooltip, ChartTooltipContent, ChartLegend, ChartLegendContent } from "@/components/ui/chart";

interface ClientSideRechartsProps {
  type: "bar" | "pie";
  data: any[];
  config?: ChartConfig;
  // Add other props as needed for specific chart types
  barKeys?: { dataKey: string; fill: string; name: string; yAxisId?: string }[];
  pieDataKey?: string;
  pieCx?: string;
  pieCy?: string;
  pieOuterRadius?: number;
  pieLabel?: (props: PieLabelRenderProps) => string;
  pieColors?: string[];
}

export function ClientSideRecharts({ type, data, config, barKeys, pieDataKey, pieCx, pieCy, pieOuterRadius, pieLabel, pieColors }: ClientSideRechartsProps) {
  if (type === "bar") {
    return (
      <ResponsiveContainer width="100%" height={300}>
        <BarChart data={data}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="name" />
          <YAxis />
          <Tooltip />
          {barKeys?.map((key, index) => (
            <Bar key={index} dataKey={key.dataKey} fill={key.fill} name={key.name} yAxisId={key.yAxisId} />
          ))}
        </BarChart>
      </ResponsiveContainer>
    );
  } else if (type === "pie") {
    return (
      <ResponsiveContainer width="100%" height={300}>
        <PieChart>
          <Pie
            data={data}
            cx={pieCx}
            cy={pieCy}
            labelLine={false}
            label={pieLabel}
            outerRadius={pieOuterRadius}
            fill="#8884d8"
            dataKey={pieDataKey}
          >
            {data.map((entry, index) => (
              <Cell key={`cell-${index}`} fill={pieColors?.[index % pieColors.length]} />
            ))}
          </Pie>
          <Tooltip />
        </PieChart>
      </ResponsiveContainer>
    );
  }
  return null;
}
