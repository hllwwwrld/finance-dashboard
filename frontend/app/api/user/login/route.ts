import { type NextRequest, NextResponse } from "next/server"

// Мок-пользователи для локального тестирования
const DEMO_USERS = {
  admin: "password123",
  user: "demo",
}

export async function POST(request: NextRequest) {
  try {
    const { login, password } = await request.json()

    if (!login || !password) {
      return NextResponse.json({ success: false, message: "Логин и пароль обязательны" }, { status: 400 })
    }

    if (DEMO_USERS[login as keyof typeof DEMO_USERS] !== password) {
      return NextResponse.json({ success: false, message: "Неверный логин или пароль" }, { status: 401 })
    }

    const response = NextResponse.json(
      { success: true, message: "Успешная авторизация" },
      { status: 200 },
    )

    response.cookies.set("auth_token", `user_${login}_${Date.now()}`, {
      httpOnly: true,
      secure: process.env.NODE_ENV === "production",
      sameSite: "lax",
      maxAge: 60 * 60 * 24 * 7,
      path: "/",
    })

    return response
  } catch (error) {
    console.error("Ошибка в /api/user/login:", error)
    return NextResponse.json(
      { success: false, message: "Внутренняя ошибка сервера" },
      { status: 500 },
    )
  }
}

