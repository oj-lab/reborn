import { Area, AreaChart, XAxis, YAxis } from "recharts"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart"
import { useTranslation } from 'react-i18next'
import { TrendingUp } from 'lucide-react'

// Mock data, should be fetched from API in production
const chartData = [
  { date: "2024-06-01", activeUsers: 1200 },
  { date: "2024-06-02", activeUsers: 1350 },
  { date: "2024-06-03", activeUsers: 1180 },
  { date: "2024-06-04", activeUsers: 1420 },
  { date: "2024-06-05", activeUsers: 1560 },
  { date: "2024-06-06", activeUsers: 1320 },
  { date: "2024-06-07", activeUsers: 1680 },
  { date: "2024-06-08", activeUsers: 1750 },
  { date: "2024-06-09", activeUsers: 1620 },
  { date: "2024-06-10", activeUsers: 1890 },
  { date: "2024-06-11", activeUsers: 2103 },
  { date: "2024-06-12", activeUsers: 1980 },
  { date: "2024-06-13", activeUsers: 2250 },
  { date: "2024-06-14", activeUsers: 2103 },
]

const chartConfig = {
  activeUsers: {
    label: "活跃用户",
    theme: {
      light: "var(--primary)",
      dark: "var(--primary)",
    },
  },
}

export function DailyActiveUsersChart() {
  const { t } = useTranslation()

  // Calculate trend percentage
  const currentValue = chartData[chartData.length - 1]?.activeUsers || 0
  const previousValue = chartData[chartData.length - 2]?.activeUsers || 0
  const trendPercentage = previousValue > 0 ? ((currentValue - previousValue) / previousValue * 100) : 0
  const isPositiveTrend = trendPercentage >= 0

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <div>
          <CardTitle className="flex items-center gap-2">
            <TrendingUp className="h-4 w-4" />
            {t('dashboard.dailyActiveUsers', '每日活跃用户')}
          </CardTitle>
          <CardDescription>
            {t('dashboard.dailyActiveUsersDescription', '过去14天的每日活跃用户数量趋势')}
          </CardDescription>
        </div>
        <div className="text-right">
          <div className="text-2xl font-bold">{currentValue.toLocaleString()}</div>
          <p className={`text-xs flex items-center gap-1 ${
            isPositiveTrend ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400'
          }`}>
            <span>{isPositiveTrend ? '↗' : '↘'}</span>
            {Math.abs(trendPercentage).toFixed(1)}% {t('dashboard.fromYesterday', '较昨日')}
          </p>
        </div>
      </CardHeader>
      <CardContent>
        <ChartContainer config={chartConfig} className="aspect-auto h-[250px] w-full">
          <AreaChart
            accessibilityLayer
            data={chartData}
            margin={{
              left: 12,
              right: 12,
              top: 5,
              bottom: 5,
            }}
          >
            <defs>
              <linearGradient id="fillActiveUsers" x1="0" y1="0" x2="0" y2="1">
                <stop
                  offset="5%"
                  stopColor="var(--color-activeUsers)"
                  stopOpacity={0.6}
                />
                <stop
                  offset="95%"
                  stopColor="var(--color-activeUsers)"
                  stopOpacity={0.1}
                />
              </linearGradient>
            </defs>
            <XAxis
              dataKey="date"
              tickLine={false}
              axisLine={false}
              tickMargin={8}
              tickFormatter={(value) => {
                const date = new Date(value)
                return date.toLocaleDateString("zh-CN", {
                  month: "short",
                  day: "numeric",
                })
              }}
            />
            <YAxis
              tickLine={false}
              axisLine={false}
              tickMargin={8}
              tickFormatter={(value) => `${(value / 1000).toFixed(1)}k`}
            />
            <ChartTooltip
              cursor={false}
              content={
                <ChartTooltipContent
                  labelFormatter={(value) => {
                    return new Date(value).toLocaleDateString("zh-CN", {
                      month: "short",
                      day: "numeric",
                      year: "numeric",
                    })
                  }}
                  indicator="dot"
                />
              }
            />
            <Area
              dataKey="activeUsers"
              type="monotone"
              fill="url(#fillActiveUsers)"
              fillOpacity={1}
              stroke="var(--color-activeUsers)"
              strokeWidth={2}
              dot={{
                fill: "var(--color-activeUsers)",
                strokeWidth: 2,
                r: 0,
              }}
              activeDot={{
                r: 4,
                fill: "var(--color-activeUsers)",
                strokeWidth: 2,
                stroke: "hsl(var(--background))",
              }}
            />
          </AreaChart>
        </ChartContainer>
      </CardContent>
    </Card>
  )
}
