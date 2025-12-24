"use client"

import { ResponsiveContainer, PieChart, Pie, Cell } from "recharts"
import { Card } from "@/components/ui/card"
import { TrendingUp } from "lucide-react"

interface Payment {
  id: string
  name: string
  amount: number
  dueDate: number
  category: string
  color: string
}

interface SpendingChartProps {
  payments: Payment[]
}

export default function SpendingChart({ payments }: SpendingChartProps) {
  // Группировка по категориям, надо перенести на бек
  // как варик -- сделать ручку получения платежей по категориям на беке, тут можно будет упростить
  // принимать прост клиент айди и по нему вытаскивать все категории 
  const categoryData = payments.reduce(
    (acc, payment) => {
      const existing = acc.find((c) => c.name === payment.category)
      if (existing) {
        existing.value += payment.amount
      } else {
        acc.push({
          name: payment.category,
          value: payment.amount,
        })
      }
      return acc
    },
    [] as Array<{ name: string; value: number }>,
  )

  const CHART_COLORS = ["#4F46E5", "#06B6D4", "#F59E0B", "#EF4444", "#8B5CF6", "#EC4899"]

  const totalExpenses = categoryData.reduce((sum, cat) => sum + cat.value, 0)

  return (
    <Card className="bg-card p-6 shadow-none border border-border/50">
      <div className="flex items-center gap-3 mb-6">
        <TrendingUp className="h-5 w-5 text-primary" />
        <h2 className="text-lg font-semibold text-foreground">Расходы по категориям</h2>
      </div>

      {categoryData.length === 0 ? (
        <p className="text-center text-muted-foreground py-8">Нет данных для отображения</p>
      ) : (
        <div className="flex flex-col lg:flex-row items-center justify-between gap-6">
          <div className="w-full">
            <ResponsiveContainer width="100%" height={200}>
              <PieChart>
                <Pie
                  data={categoryData}
                  cx="50%"
                  cy="50%"
                  innerRadius={60}
                  outerRadius={100}
                  paddingAngle={2}
                  dataKey="value"
                >
                  {categoryData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={CHART_COLORS[index % CHART_COLORS.length]} />
                  ))}
                </Pie>
              </PieChart>
            </ResponsiveContainer>
          </div>

          <div className="w-full space-y-2">
            {categoryData.map((category, index) => (
              <div key={category.name} className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <div
                    className="h-3 w-3 rounded-full"
                    style={{
                      backgroundColor: CHART_COLORS[index % CHART_COLORS.length],
                    }}
                  />
                  <span className="text-sm font-medium text-foreground">{category.name}</span>
                </div>
                <div className="text-right">
                  <p className="text-sm font-semibold text-foreground">{category.value.toLocaleString()} ₽</p>
                  <p className="text-xs text-muted-foreground">
                    {((category.value / totalExpenses) * 100).toFixed(0)}%
                  </p>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </Card>
  )
}
