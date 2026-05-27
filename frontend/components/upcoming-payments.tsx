"use client"

import { Calendar, Clock } from "lucide-react"
import { Card } from "@/components/ui/card"

interface Payment {
    id: string
    name: string
    amount: number
    dueDay: number // день месяца
    daysUntil: number // дни до платежа
    category: string
    color: string
}

interface UpcomingPaymentsProps {
    payments: Payment[]
    onOpen: () => void
}

export default function UpcomingPayments({ payments, onOpen }: UpcomingPaymentsProps) {
    const today = new Date()
    const currentDay = today.getDate()

    const upcoming = payments
        // .map((p) => ({
        //     ...p,
        //     daysUntil: p.dueDate >= currentDay ? p.dueDate - currentDay : p.dueDate + (31 - currentDay),
        // })) убрал, беру теперь с бека, но не тестил
        .filter((a) => a.daysUntil >= 0)
        .sort((a, b) => a.daysUntil - b.daysUntil)
        .slice(0, 5)

    return (
        <Card
            role="button"
            tabIndex={0}
            onClick={onOpen}
            onKeyDown={(e) => {
                if (e.key === "Enter" || e.key === " ") {
                    e.preventDefault()
                    onOpen()
                }
            }}
            className="cursor-pointer bg-card p-6 shadow-none border border-border/50 transition-colors hover:border-primary/40 hover:bg-muted/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
            aria-label="Открыть список платежей"
        >
            <div className="flex items-center gap-3 mb-6">
                <Calendar className="h-5 w-5 text-primary" />
                <h2 className="text-lg font-semibold text-foreground">Предстоящие платежи</h2>
            </div>

            <div className="space-y-3">
                {upcoming.length === 0 ? (
                    <p className="text-center text-muted-foreground py-8">Нет предстоящих платежей</p>
                ) : (
                    upcoming.map((payment) => (
                        <div
                            key={payment.id}
                            className="flex items-center justify-between p-3 rounded-lg bg-muted hover:bg-muted/80 transition-colors"
                        >
                            <div className="flex items-center gap-3 flex-1">
                                <div className="h-8 w-8 rounded-lg flex-shrink-0" style={{ backgroundColor: payment.color }} />
                                <div className="min-w-0 flex-1">
                                    <p className="font-medium text-foreground truncate">{payment.name}</p>
                                    <p className="text-xs text-muted-foreground">
                                        {payment.daysUntil === 0
                                            ? "Сегодня"
                                            : payment.daysUntil === 1
                                                ? "Завтра"
                                                : `Через ${payment.daysUntil} дн.`}
                                    </p>
                                </div>
                            </div>
                            <div className="flex items-center gap-2">
                                <p className="font-semibold text-foreground whitespace-nowrap">{payment.amount.toLocaleString()} ₽</p>
                                {payment.daysUntil <= 3 && <Clock className="h-4 w-4 text-orange-500" />}
                            </div>
                        </div>
                    ))
                )}
            </div>
        </Card>
    )
}
