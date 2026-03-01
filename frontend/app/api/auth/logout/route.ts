import { NextResponse } from "next/server"

export async function POST() {
  const response = NextResponse.json({ message: "Выход выполнен" }, { status: 200 })

  // Удаляем cookie
  response.cookies.delete("auth_token")

  return response
}
