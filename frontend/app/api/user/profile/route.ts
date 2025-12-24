// API route для работы с профилем пользователя
// GET /api/user/profile - получить профиль (включая доход)
// PUT /api/user/profile - обновить профиль (включая доход)

import { NextRequest, NextResponse } from 'next/server'

// Интерфейс для профиля пользователя
export interface UserProfile {
  monthlyIncome: number
  // В будущем здесь могут быть другие поля профиля:
  // currency: string
  // language: string
  // etc.
}

// GET /api/user/profile - получить профиль пользователя
export async function GET(request: NextRequest) {
  try {
    // TODO: В реальном приложении здесь будет:
    // 1. Получение ID пользователя из сессии/токена
    // 2. Запрос к базе данных для получения профиля
    // const userId = await getUserIdFromSession(request)
    // const profile = await db.userProfile.findUnique({ where: { userId } })
    
    // Для примера используем mock данные
    // В реальном приложении, если профиля нет, можно вернуть значения по умолчанию
    const mockProfile: UserProfile = {
      monthlyIncome: 50000, // Значение по умолчанию
    }

    return NextResponse.json(mockProfile)
  } catch (error) {
    console.error('Ошибка при получении профиля:', error)
    return NextResponse.json(
      { error: 'Не удалось загрузить профиль' },
      { status: 500 }
    )
  }
}

// PUT /api/user/profile - обновить профиль пользователя
export async function PUT(request: NextRequest) {
  try {
    const body = await request.json()
    const { monthlyIncome } = body

    // Валидация
    if (monthlyIncome === undefined || monthlyIncome === null) {
      return NextResponse.json(
        { error: 'Поле monthlyIncome обязательно для заполнения' },
        { status: 400 }
      )
    }

    const income = Number(monthlyIncome)
    if (isNaN(income) || income <= 0) {
      return NextResponse.json(
        { error: 'Доход должен быть положительным числом' },
        { status: 400 }
      )
    }

    // TODO: В реальном приложении здесь будет:
    // 1. Получение ID пользователя из сессии/токена
    // 2. Сохранение в базу данных
    // const userId = await getUserIdFromSession(request)
    // const profile = await db.userProfile.upsert({
    //   where: { userId },
    //   update: { monthlyIncome: income },
    //   create: { userId, monthlyIncome: income }
    // })

    // Для примера возвращаем обновленный профиль
    const updatedProfile: UserProfile = {
      monthlyIncome: income,
    }

    return NextResponse.json(updatedProfile)
  } catch (error) {
    console.error('Ошибка при обновлении профиля:', error)
    return NextResponse.json(
      { error: 'Не удалось обновить профиль' },
      { status: 500 }
    )
  }
}
