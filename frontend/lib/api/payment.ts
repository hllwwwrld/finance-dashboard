import { handleApiResponse } from './utils'

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || '/api'

export interface Payment {
    id: string
    name: string
    amount: number
    dueDay: number // день месяца
    daysUntil: number // дни до платежа
    category: string
    color: string
}

export interface FetchPaymentsResponse {
    payments: Payment[]
    totalExpenses: number
}

export interface DeletePaymentRequest {
    id: string
}

// Получить список всех платежей
export async function fetchPayments(): Promise<FetchPaymentsResponse> {
    try {
        const response = await fetch(`${API_BASE_URL}/payments/list`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
            // Для работы с credentials (cookies, авторизация)
            credentials: 'include',
        })

        handleApiResponse(response)

        if (!response.ok) {
            throw new Error(`Ошибка загрузки платежей: ${response.statusText}`)
        }

        const data = await response.json()
        return data
    } catch (error) {
        console.error('Ошибка при загрузке платежей:', error)
        throw error
    }
}

// Создать новый платеж
export async function createPayment(payment: Omit<Payment, 'id' | 'daysUntil'>): Promise<Payment> {
    try {
        const response = await fetch(`${API_BASE_URL}/payments/create`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
            body: JSON.stringify(payment),
        })

        handleApiResponse(response)

        if (!response.ok) {
            throw new Error(`Ошибка создания платежа: ${response.statusText}`)
        }

        const data = await response.json()
        return data
    } catch (error) {
        console.error('Ошибка при создании платежа:', error)
        throw error
    }
}

// Удалить платеж
export async function deletePayment(request: DeletePaymentRequest): Promise<void> {
    try {
        const response = await fetch(`${API_BASE_URL}/api/payments/delete`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
            body: JSON.stringify(request),
        })

        handleApiResponse(response)

        if (!response.ok) {
            throw new Error(`Ошибка удаления платежа: ${response.statusText}`)
        }
    } catch (error) {
        console.error('Ошибка при удалении платежа:', error)
        throw error
    }
}