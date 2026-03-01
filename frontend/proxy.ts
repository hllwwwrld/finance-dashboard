import { type NextRequest, NextResponse } from "next/server"

export function proxy(request: NextRequest) {
  const token = request.cookies.get("auth_token")?.value
  const pathname = request.nextUrl.pathname

  // Страницы, доступные без авторизации
  const isAuthPage = pathname === "/login" || pathname === "/register"

  // Если нет токена и это не страница логина/регистрации - редирект на логин
  if (!token && !isAuthPage) {
    return NextResponse.redirect(new URL("/login", request.url))
  }

  // Если есть токен и пользователь на странице логина/регистрации - редирект на дашборд
  // if (token && isAuthPage) {
  //   return NextResponse.redirect(new URL("/dashboard", request.url))
  // }

  return NextResponse.next()
}

export const config = {
  matcher: ["/((?!api|_next/static|_next/image|favicon.ico|icon.*|apple-icon.png).*)"],
}