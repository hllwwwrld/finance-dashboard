import { type NextRequest, NextResponse } from "next/server"

export async function POST(request: NextRequest) {
  try {
    const { login, password } = await request.json()

    if (!login || !password) {
      return NextResponse.json({ success: false, message: "Логин и пароль обязательны" }, { status: 400 })
    }

    // В реальном приложении здесь была бы проверка на существование пользователя и сохранение в БД
    // Для мока считаем, что регистрация всегда успешна
    return NextResponse.json(
      { success: true, message: "Регистрация прошла успешно" },
      { status: 200 },
    )
  } catch (error) {
    console.error("Ошибка в /api/user/register:", error)
    return NextResponse.json(
      { success: false, message: "Внутренняя ошибка сервера" },
      { status: 500 },
    )
  }
}

