"use client"

import { useState } from "react"
import { Calendar, Clock } from "lucide-react"
import { Card } from "@/components/ui/card"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"

interface Payment {
  id: string
  name: string
  amount: number
  dueDate: number
  category: string
  color: string
}

interface UpcomingPaymentsProps {
  payments: Payment[]
}

export default function UpcomingPayments({ payments }: UpcomingPaymentsProps) {
  const [isDialogOpen, setIsDialogOpen] = useState(false)
  const today = new Date()
  const currentDay = today.getDate()

  const upcoming = payments
    .map((p) => ({
      ...p,
      daysUntil: p.dueDate >= currentDay ? p.dueDate - currentDay : p.dueDate + (31 - currentDay),
    }))
    .sort((a, b) => a.daysUntil - b.daysUntil)

  // todo УДАЛИТЬ то, что выше, кол-вод ней до платежа должно прихдить с бека
  // при расчете должны учитывать весокосные года, кол-во дней в месяце и тд

  const upcomingPreview = upcoming.slice(0, 4)

  return (
    <>
      <Card 
        className="bg-card p-6 shadow-none border border-border/50 cursor-pointer hover:border-primary/50 transition-colors"
        onClick={() => setIsDialogOpen(true)}
      >
        <div className="flex items-center gap-3 mb-6">
          <Calendar className="h-5 w-5 text-primary" />
          <h2 className="text-lg font-semibold text-foreground">Предстоящие платежи</h2>
        </div>

        <div className="space-y-3">
          {upcomingPreview.length === 0 ? (
            <p className="text-center text-muted-foreground py-8">Нет предстоящих платежей</p>
          ) : (
            upcomingPreview.map((payment) => (
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

      <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
        <DialogContent className="max-w-2xl max-h-[80vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-3">
              <Calendar className="h-5 w-5 text-primary" />
              Все предстоящие платежи
            </DialogTitle>
          </DialogHeader>

          <div className="space-y-3 mt-4">
            {upcoming.length === 0 ? (
              <p className="text-center text-muted-foreground py-8">Нет предстоящих платежей</p>
            ) : (
              upcoming.map((payment) => (
                <div
                  key={payment.id}
                  className="flex items-center justify-between p-4 rounded-lg bg-muted hover:bg-muted/80 transition-colors"
                >
                  <div className="flex items-center gap-3 flex-1">
                    <div className="h-10 w-10 rounded-lg flex-shrink-0" style={{ backgroundColor: payment.color }} />
                    <div className="min-w-0 flex-1">
                      <p className="font-medium text-foreground">{payment.name}</p>
                      <div className="flex items-center gap-2 mt-1">
                        <p className="text-xs text-muted-foreground">
                          {payment.daysUntil === 0
                            ? "Сегодня"
                            : payment.daysUntil === 1
                              ? "Завтра"
                              : `Через ${payment.daysUntil} дн.`}
                        </p>
                        <span className="text-xs text-muted-foreground">•</span>
                        <p className="text-xs text-muted-foreground">{payment.dueDate} число</p>
                        <span className="text-xs text-muted-foreground">•</span>
                        <p className="text-xs text-muted-foreground">{payment.category}</p>
                      </div>
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    <p className="font-semibold text-foreground whitespace-nowrap text-lg">
                      {payment.amount.toLocaleString()} ₽
                    </p>
                    {payment.daysUntil <= 3 && <Clock className="h-5 w-5 text-orange-500" />}
                  </div>
                </div>
              ))
            )}
          </div>
        </DialogContent>
      </Dialog>
    </>
  )
}
