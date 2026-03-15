"use client"

import {useState, useEffect} from "react"
import {useRouter} from "next/navigation"
import {Trash2, ChevronLeft, ChevronRight, Clock, Loader2, RefreshCw, LogOut} from "lucide-react"
import {Button} from "@/components/ui/button"
import Dashboard from "@/components/dashboard"
import PaymentForm from "@/components/payment-form"
import {fetchPayments, createPayment, deletePayment, type Payment, DeletePaymentRequest} from "@/lib/api/payment"
import {fetchUserProfile, updateMonthlyIncome, UpdateMonthlyIncomeRequest, logout} from "@/lib/api/user"

export default function DashboardPage() {
    const router = useRouter()
    const [payments, setPayments] = useState<Payment[]>([])
    const [totalExpenses, setTotalExpenses] = useState(0)
    const [isLoading, setIsLoading] = useState(true)
    const [error, setError] = useState<string | null>(null)
    const [showForm, setShowForm] = useState(false)
    const [monthlyIncome, setMonthlyIncome] = useState<number | null>(null)
    const [paymentsOpen, setPaymentsOpen] = useState(true)
    const [sidebarOpen, setSidebarOpen] = useState(false)
    const [paymentsLoaded, setPaymentsLoaded] = useState(false)
    const [profileLoaded, setProfileLoaded] = useState(false)
    const isDataLoaded = paymentsLoaded && profileLoaded

    // Загрузка платежей из API при монтировании компонента
    useEffect(() => {
        const loadPayments: () => Promise<void> = async (): Promise<void> => {
            try {
                setIsLoading(true)
                setError(null)
                const data = await fetchPayments()
                setPayments(data.payments)
                setTotalExpenses(data.totalExpenses)
            } catch (err) {
                setError(err instanceof Error ? err.message : "Не удалось загрузить платежи")
                console.error("Ошибка загрузки платежей:", err)
            } finally {
                setIsLoading(false)
                setPaymentsLoaded(true)
            }
        }

        loadPayments()
    }, [])

    // Загрузка дохода пользователя из API при монтировании компонента
    useEffect(() => {
        const loadUserProfile = async () => {
            try {
                const profile = await fetchUserProfile()
                setMonthlyIncome(profile.monthlyIncome)
                setProfileLoaded(true)
            } catch (err) {
                console.error("Ошибка загрузки дохода:", err)
                setProfileLoaded(true)
            }
        }

        loadUserProfile()
    }, [])

    const handleAddPayment = async (payment: Omit<Payment, "id" | "daysUntil">) => {
        try {
            setError(null)
            await createPayment(payment)
            const updatePaymentsData = await fetchPayments()
            setPayments(updatePaymentsData.payments)
            setTotalExpenses(updatePaymentsData.totalExpenses)
            setShowForm(false)
        } catch (err) {
            setError(err instanceof Error ? err.message : "Не удалось создать платеж")
            console.error("Ошибка создания платежа:", err)
        }
    }

    const handleDeletePayment = async (request: DeletePaymentRequest) => {
        try {
            setError(null)
            await deletePayment(request)
            const data = await fetchPayments() // получаем свежие данные
            setPayments(data.payments)
            setTotalExpenses(data.totalExpenses)
        } catch (err) {
            setError(err instanceof Error ? err.message : "Не удалось удалить платеж")
            console.error("Ошибка удаления платежа:", err)
        }
    }

    const handleRefreshPayments = async () => {
        try {
            setIsLoading(true)
            setError(null)
            const data = await fetchPayments()
            setPayments(data.payments)
            setTotalExpenses(data.totalExpenses)
        } catch (err) {
            setError(err instanceof Error ? err.message : "Не удалось обновить платежи")
            console.error("Ошибка обновления платежей:", err)
        } finally {
            setIsLoading(false)
        }
    }

    const handleIncomeChange = async (newIncome: number) => {
        try {
            setError(null)
            const updatedProfile = await updateMonthlyIncome({income: newIncome})
            setMonthlyIncome(updatedProfile.monthlyIncome)
        } catch (err) {
            setError(err instanceof Error ? err.message : "Не удалось сохранить доход")
            console.error("Ошибка сохранения дохода:", err)
        }
    }

    const isUpcomingThisWeek = (payment: Payment) => {
        return 0 <= payment.daysUntil && payment.daysUntil <= 7
    }

    const handleLogout = async () => {
        await logout()
        router.push("/")
        router.refresh()
    }

    const computedRemaining = (monthlyIncome ?? 0) - (totalExpenses ?? 0)

    return (
        <div className="min-h-screen bg-background">
            {!sidebarOpen && (
                <div className="absolute top-4 right-4 z-50">
                    <Button onClick={handleLogout} variant="ghost" size="sm" className="gap-2">
                        <LogOut className="h-4 w-4"/>
                        Выйти
                    </Button>
                </div>
            )}

            <div className="flex h-screen flex-col lg:flex-row">
                <div
                    className={`transition-all duration-300 ${
                        paymentsOpen ? "flex-1" : "w-0"
                    } overflow-hidden lg:overflow-y-auto`}
                >
                    <div className="relative h-full">
                        <Button
                            onClick={() => setPaymentsOpen(!paymentsOpen)}
                            size="icon"
                            variant="ghost"
                            className="absolute top-6 left-6 z-10 h-9 w-9"
                        >
                            {paymentsOpen ? (
                                <ChevronLeft className="h-5 w-5"/>
                            ) : (
                                <ChevronRight className="h-5 w-5"/>
                            )}
                        </Button>

                        <div className="p-6 lg:p-8 h-full overflow-y-auto">
                            {isDataLoaded ? (
                                <Dashboard
                                    totalExpenses={totalExpenses}
                                    monthlyIncome={monthlyIncome ?? 0}
                                    remaining={computedRemaining}
                                    payments={payments}
                                    onIncomeChange={handleIncomeChange}
                                />
                            ) : (
                                <div className="flex items-center justify-center h-64">
                                    <div className="text-muted-foreground">Загрузка...</div>
                                </div>
                            )}
                        </div>
                    </div>
                </div>

                <div
                    className={`border-t border-border bg-sidebar transition-all duration-300 relative ${
                        sidebarOpen ? "p-6 lg:w-96 lg:border-t-0 lg:border-l" : "w-0 p-0 border-0 overflow-hidden"
                    } lg:overflow-y-auto`}
                >
                    {sidebarOpen && (
                        <>
                            <Button
                                onClick={() => setSidebarOpen(false)}
                                size="icon"
                                variant="ghost"
                                className="absolute top-6 left-6 z-10 h-9 w-9"
                            >
                                <ChevronRight className="h-5 w-5"/>
                            </Button>
                            <div className="flex items-center justify-between mb-6 pl-14">
                                <h2 className="text-2xl font-semibold text-foreground">Платежи</h2>
                                <div className="flex items-center gap-2">
                                    <Button
                                        onClick={handleRefreshPayments}
                                        size="icon"
                                        variant="ghost"
                                        className="h-9 w-9"
                                        disabled={isLoading}
                                        title="Обновить список платежей"
                                    >
                                        <RefreshCw className={`h-4 w-4 ${isLoading ? "animate-spin" : ""}`}/>
                                    </Button>
                                    <Button onClick={() => setShowForm(!showForm)} className="h-9 px-3">
                                        +
                                    </Button>
                                </div>
                            </div>

                            {showForm && (
                                <div className="mb-8">
                                    <PaymentForm onSubmit={handleAddPayment} onCancel={() => setShowForm(false)}/>
                                </div>
                            )}

                            {error && (
                                <div className="mb-4 p-3 rounded-lg bg-destructive/10 border border-destructive/20">
                                    <p className="text-sm text-destructive">{error}</p>
                                </div>
                            )}

                            <div className="space-y-3">
                                {isLoading ? (
                                    <div className="flex items-center justify-center py-8">
                                        <Loader2 className="h-6 w-6 animate-spin text-muted-foreground"/>
                                        <span className="ml-2 text-sm text-muted-foreground">Загрузка платежей...</span>
                                    </div>
                                ) : !Array.isArray(payments) || payments.length === 0 ? (
                                    <p className="text-center text-sidebar-foreground/50 py-8">
                                        Нет платежей. Добавьте первый!
                                    </p>
                                ) : (
                                    payments
                                        .sort((a, b) => a.dueDay - b.dueDay)
                                        .map((payment) => (
                                            <div
                                                key={payment.id}
                                                className="group flex items-center justify-between p-3 rounded-lg bg-muted hover:bg-muted/80 transition-colors"
                                            >
                                                <div className="flex items-center gap-3 flex-1 min-w-0">
                                                    <div
                                                        className="h-8 w-8 rounded-lg flex-shrink-0"
                                                        style={{backgroundColor: payment.color}}
                                                    />
                                                    <div className="min-w-0 flex-1">
                                                        <p className="font-medium text-foreground truncate">
                                                            {payment.name}
                                                        </p>
                                                        <p className="text-xs text-muted-foreground">
                                                            {payment.dueDay} число
                                                        </p>
                                                    </div>
                                                </div>
                                                <div className="flex items-center gap-2">
                                                    <p className="font-semibold text-foreground whitespace-nowrap">
                                                        {payment.amount.toLocaleString()} ₽
                                                    </p>
                                                    {isUpcomingThisWeek(payment) && (
                                                        <Clock className="h-4 w-4 text-orange-500 flex-shrink-0"/>
                                                    )}
                                                    <Button
                                                        onClick={() => handleDeletePayment({id: payment.id})}
                                                        size="icon"
                                                        variant="ghost"
                                                        className="h-6 w-6 opacity-0 transition-opacity group-hover:opacity-100 flex-shrink-0"
                                                    >
                                                        <Trash2 className="h-3 w-3 text-destructive"/>
                                                    </Button>
                                                </div>
                                            </div>
                                        ))
                                )}
                            </div>
                        </>
                    )}
                </div>

                {!sidebarOpen && (
                    <Button
                        onClick={() => setSidebarOpen(true)}
                        size="icon"
                        variant="ghost"
                        className="fixed top-1/2 right-4 z-10 h-9 w-9 -translate-y-1/2"
                    >
                        <ChevronLeft className="h-5 w-5"/>
                    </Button>
                )}
            </div>
        </div>
    )
}

