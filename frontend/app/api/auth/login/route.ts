import { type NextRequest, NextResponse } from "next/server"

// В продакшене здесь должна быть проверка в базе данных
const DEMO_USERS = {
  admin: "password123",
  user: "demo",
}

export async function POST(request: NextRequest) {
  try {
    const { username, password } = await request.json()

    // Валидация
    if (!username || !password) {
      return NextResponse.json({ message: "Логин и пароль обязательны" }, { status: 400 })
    }

    // Проверка учетных данных
    if (DEMO_USERS[username as keyof typeof DEMO_USERS] !== password) {
      return NextResponse.json({ message: "Неверный логин или пароль" }, { status: 401 })
    }

    // Создаем ответ с успешной авторизацией
    const response = NextResponse.json({ message: "Успешная авторизация" }, { status: 200 })

    // Устанавливаем cookie с токеном (в продакшене использовать JWT)
    response.cookies.set("auth_token", `user_${username}_${Date.now()}`, {
      httpOnly: true,
      secure: process.env.NODE_ENV === "production",
      sameSite: "lax",
      maxAge: 60 * 60 * 24 * 7, // 7 дней
      path: "/",
    })

    return response
  } catch (error) {
    return NextResponse.json({ message: "Внутренняя ошибка сервера" }, { status: 500 })
  }
}
