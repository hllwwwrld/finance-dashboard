"use client"

import {useState} from "react"
import {Check, TrendingDown, Wallet} from "lucide-react"
import {Card} from "@/components/ui/card"
import {Button} from "@/components/ui/button"
import {cn} from "@/lib/utils"
import UpcomingPayments from "./upcoming-payments"
import SpendingChart from "./spending-chart"

interface Payment {
    id: string
    name: string
    amount: number
    dueDay: number // день месяца
    daysUntil: number // дни до платежа
    category: string
    color: string
}

interface DashboardProps {
    totalExpenses: number
    monthlyIncome: number
    remaining: number
    payments: Payment[]
    onIncomeChange: (income: number) => void
    onOpenPayments: () => void
}

export default function Dashboard({
                                      totalExpenses = 0,
                                      monthlyIncome = 0,
                                      remaining = 0,
                                      payments = [],
                                      onIncomeChange,
                                      onOpenPayments,
                                  }: DashboardProps) {
    const [isEditingIncome, setIsEditingIncome] = useState(false)
    const [incomeInput, setIncomeInput] = useState(monthlyIncome?.toString() || '0')

    const startEditingIncome = () => {
        setIncomeInput(monthlyIncome?.toString() || "0")
        setIsEditingIncome(true)
    }

    const handleSaveIncome = () => {
        const newIncome = Number.parseFloat(incomeInput)
        if (!isNaN(newIncome) && newIncome > 0) {
            onIncomeChange(newIncome)
            setIsEditingIncome(false)
        }
    }

    const handleCancelEditIncome = () => {
        setIncomeInput(monthlyIncome?.toString() || "0")
        setIsEditingIncome(false)
    }

    const expensePercentage = monthlyIncome > 0 ? (totalExpenses / monthlyIncome) * 100 : 0

    return (
        <div className="space-y-6">
            {/* Финансовый обзор */}
            <div>
                <h2 className="text-2xl font-semibold text-foreground mb-1">Финансовый обзор</h2>
                <p className="text-sm text-muted-foreground">
                    {new Date().toLocaleDateString("ru-RU", {year: "numeric", month: "long", day: "numeric"})}
                </p>
            </div>

            {/* Статистика */}
            <div className="grid gap-4 md:grid-cols-3">
                {/* Доход */}
                <Card
                    role={isEditingIncome ? undefined : "button"}
                    tabIndex={isEditingIncome ? undefined : 0}
                    onClick={() => !isEditingIncome && startEditingIncome()}
                    onKeyDown={(e) => {
                        if (!isEditingIncome && (e.key === "Enter" || e.key === " ")) {
                            e.preventDefault()
                            startEditingIncome()
                        }
                    }}
                    className={cn(
                        "bg-card p-6 shadow-none border border-border/50",
                        !isEditingIncome && "cursor-pointer hover:border-primary/50 hover:bg-muted/30 transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                    )}
                >
                    <div className="flex items-start justify-between mb-4">
                        <div>
                            <p className="text-sm font-medium text-muted-foreground mb-1">Доход</p>
                            {isEditingIncome ? (
                                <div
                                    className="flex gap-2"
                                    onClick={(e) => e.stopPropagation()}
                                    onKeyDown={(e) => e.stopPropagation()}
                                >
                                    <input
                                        type="number"
                                        value={incomeInput}
                                        onChange={(e) => setIncomeInput(e.target.value)}
                                        onKeyDown={(e) => {
                                            if (e.key === "Enter") handleSaveIncome()
                                            if (e.key === "Escape") handleCancelEditIncome()
                                        }}
                                        className="bg-muted rounded px-2 py-1 text-foreground font-bold text-2xl w-32"
                                        autoFocus
                                    />
                                    <Button
                                        onClick={handleSaveIncome}
                                        size="icon"
                                        className="hover:text-primary"
                                        aria-label="Сохранить доход"
                                    >
                                        <Check className="size-5" strokeWidth={2} />
                                    </Button>
                                </div>
                            ) : (
                                <p className="text-3xl font-bold text-foreground">
                                    {monthlyIncome.toLocaleString()} ₽
                                </p>
                            )}
                        </div>
                        <Wallet className="h-8 w-8 text-primary pointer-events-none"/>
                    </div>
                </Card>

                {/* Расходы */}
                <Card className="bg-card p-6 shadow-none border border-border/50">
                    <div className="flex items-start justify-between mb-4">
                        <div>
                            <p className="text-sm font-medium text-muted-foreground mb-1">Расходы</p>
                            <p className="text-3xl font-bold text-destructive">{totalExpenses.toLocaleString()} ₽</p>
                            <p className="text-xs text-muted-foreground mt-2">{expensePercentage.toFixed(1)}% от
                                дохода</p>
                        </div>
                        <TrendingDown className="h-8 w-8 text-destructive"/>
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
                <SpendingChart payments={payments}/>

                {/* Предстоящие платежи */}
                <UpcomingPayments payments={payments} onOpen={onOpenPayments}/>
            </div>
        </div>
    )
}

