// Пример API route для Next.js
// Этот файл создан для примера. Замените на реальный эндпоинт вашего бекенда.

import { NextRequest, NextResponse } from 'next/server'

export interface Payment {
  id: string
  name: string
  amount: number
  dueDate: number
  daysUntil: number
  category: string
  color: string
}

export interface FetchPaymentsResponse {
  payments: Payment[]
  totalExpenses: number
}

// Пример данных для тестирования (замените на запрос к реальной БД)
const mockPayments: FetchPaymentsResponse = {
  payments: [
    {
      id: "1",
      name: "Netflix",
      amount: 499,
      dueDate: 5,
      daysUntil: -1,
      category: "Подписки",
      color: "#E50914",
    },
    {
      id: "2",
      name: "Spotify",
      amount: 169,
      dueDate: 10,
      daysUntil: -1,
      category: "Подписки",
      color: "#1DB954",
    },
    {
      id: "3",
      name: "Apple Music",
      amount: 169,
      dueDate: 15,
      daysUntil: 4,
      category: "Подписки",
      color: "#000000",
    },
    {
      id: "4",
      name: "Интернет",
      amount: 1999,
      dueDate: 1,
      daysUntil: -1,
      category: "Коммунальные",
      color: "#4F46E5",
    },
  ],
  totalExpenses: 2836,
}
  

// GET /api/payments - получить список платежей
export async function GET(request: NextRequest) {
  try {
    // TODO: Замените на запрос к реальной базе данных
    // const payments = await db.payments.findMany()
    
    // Для примера возвращаем mock данные
    return NextResponse.json(mockPayments)
  } catch (error) {
    console.error('Ошибка при получении платежей:', error)
    return NextResponse.json(
      { error: 'Не удалось загрузить платежи' },
      { status: 500 }
    )
  }
}

// POST /api/payments - создать новый платеж
export async function POST(request: NextRequest) {
  try {
    const body = await request.json()
    const { name, amount, dueDate, category, color } = body

    // Валидация
    if (!name || !amount || !dueDate || !category || !color) {
      return NextResponse.json(
        { error: 'Не все обязательные поля заполнены' },
        { status: 400 }
      )
    }

    // TODO: Замените на сохранение в реальную базу данных
    // const payment = await db.payments.create({ data: { ... } })
    
    // Для примера создаем новый платеж
    const newPayment: Payment = {
      id: Date.now().toString(),
      name,
      amount: Number(amount),
      dueDate: Number(dueDate),
      daysUntil: -1,
      category,
      color,
    }

    return NextResponse.json(newPayment, { status: 201 })
  } catch (error) {
    console.error('Ошибка при создании платежа:', error)
    return NextResponse.json(
      { error: 'Не удалось создать платеж' },
      { status: 500 }
    )
  }
}
