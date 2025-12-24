"use client"

import { useState } from "react"
import { TrendingDown, Wallet } from "lucide-react"
import { Card } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import UpcomingPayments from "./upcoming-payments"
import SpendingChart from "./spending-chart"

interface Payment {
  id: string
  name: string
  amount: number
  dueDate: number
  category: string
  color: string
}

interface DashboardProps {
  totalExpenses: number
  monthlyIncome: number
  remaining: number
  payments: Payment[]
  onIncomeChange: (income: number) => void
}

export default function Dashboard({
  totalExpenses,
  monthlyIncome,
  remaining,
  payments,
  onIncomeChange,
}: DashboardProps) {
  const [isEditingIncome, setIsEditingIncome] = useState(false)
  const [incomeInput, setIncomeInput] = useState(monthlyIncome.toString())

  const handleSaveIncome = () => {
    const newIncome = Number.parseFloat(incomeInput)
    if (!isNaN(newIncome) && newIncome > 0) {
      onIncomeChange(newIncome)
      setIsEditingIncome(false)
    }
  }

  const expensePercentage = (totalExpenses / monthlyIncome) * 100

  return (
    <div className="space-y-6">
      {/* Финансовый обзор */}
      <div>
        <h2 className="text-2xl font-semibold text-foreground mb-1">Финансовый обзор</h2>
        <p className="text-sm text-muted-foreground">
          {new Date().toLocaleDateString("ru-RU", { year: "numeric", month: "long", day: "numeric" })}
        </p>
      </div>

      {/* Статистика */}
      <div className="grid gap-4 md:grid-cols-3">
        {/* Доход */}
        <Card className="bg-card p-6 shadow-none border border-border/50">
          <div className="flex items-start justify-between mb-4">
            <div>
              <p className="text-sm font-medium text-muted-foreground mb-1">Доход</p>
              {isEditingIncome ? (
                <div className="flex gap-2">
                  <input
                    type="number"
                    value={incomeInput}
                    onChange={(e) => setIncomeInput(e.target.value)}
                    className="bg-muted rounded px-2 py-1 text-foreground font-bold text-2xl w-32"
                    autoFocus
                  />
                  <Button onClick={handleSaveIncome} size="sm" className="h-8">
                    ✓
                  </Button>
                </div>
              ) : (
                <p
                  className="text-3xl font-bold text-foreground cursor-pointer hover:text-primary transition-colors"
                  onClick={() => setIsEditingIncome(true)}
                >
                  {monthlyIncome.toLocaleString()} ₽
                </p>
              )}
            </div>
            <Wallet className="h-8 w-8 text-primary" />
          </div>
        </Card>

        {/* Расходы */}
        <Card className="bg-card p-6 shadow-none border border-border/50">
          <div className="flex items-start justify-between mb-4">
            <div>
              <p className="text-sm font-medium text-muted-foreground mb-1">Расходы</p>
              <p className="text-3xl font-bold text-destructive">{totalExpenses.toLocaleString()} ₽</p>
              <p className="text-xs text-muted-foreground mt-2">{expensePercentage.toFixed(1)}% от дохода</p>
            </div>
            <TrendingDown className="h-8 w-8 text-destructive" />
          </div>
        </Card>

        {/* Остаток */}
        <Card
          className={`p-6 shadow-none border ${
            remaining > 0
              ? "bg-green-50 border-green-200 dark:bg-green-950 dark:border-green-800"
              : "bg-red-50 border-red-200 dark:bg-red-950 dark:border-red-800"
          }`}
        >
          <div className="flex items-start justify-between mb-4">
            <div>
              <p className="text-sm font-medium text-muted-foreground mb-1">Остаток</p>
              <p
                className={`text-3xl font-bold ${
                  remaining > 0 ? "text-green-600 dark:text-green-400" : "text-red-600 dark:text-red-400"
                }`}
              >
                {remaining.toLocaleString()} ₽
              </p>
            </div>
            <div
              className={`h-8 w-8 rounded-full flex items-center justify-center ${
                remaining > 0 ? "bg-green-200 dark:bg-green-800" : "bg-red-200 dark:bg-red-800"
              }`}
            >
              <span className="text-lg font-bold">{remaining > 0 ? "✓" : "!"}</span>
            </div>
          </div>
        </Card>
      </div>

      <div className="grid gap-6 lg:grid-cols-2">
        {/* График расходов */}
        <SpendingChart payments={payments} />

        {/* Предстоящие платежи */}
        <UpcomingPayments payments={payments} />
      </div>
    </div>
  )
}
