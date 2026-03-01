"use client"

import type React from "react"

import { useState } from "react"
import { useRouter } from "next/navigation"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { login } from "@/lib/api/user"

export default function LoginPage() {
  const router = useRouter()
  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")
  const [error, setError] = useState("")
  const [isLoading, setIsLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError("")
    setIsLoading(true)

    const result = await login({ login: username, password })

    if (result.success) {
      router.push("/dashboard")
      router.refresh()
    } else {
      setError(result.message || "Ошибка авторизации")
    }

    setIsLoading(false)
  }

  return (
    <div className="min-h-screen bg-background flex items-center justify-center p-4">
      <div className="w-full max-w-sm">
        <div className="bg-card rounded-2xl shadow-sm border border-border p-8">
          <h1 className="text-2xl font-semibold text-foreground mb-2">Вход</h1>
          <p className="text-sm text-muted-foreground mb-6">Войдите в свой аккаунт</p>

          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <label htmlFor="username" className="text-sm font-medium text-foreground">
                Логин
              </label>
              <Input
                id="username"
                type="text"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                placeholder="Введите логин"
                required
                disabled={isLoading}
                className="h-11"
              />
            </div>

            <div className="space-y-2">
              <label htmlFor="password" className="text-sm font-medium text-foreground">
                Пароль
              </label>
              <Input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="Введите пароль"
                required
                disabled={isLoading}
                className="h-11"
              />
            </div>

            {error && <div className="text-sm text-red-500 bg-red-50 dark:bg-red-950/20 p-3 rounded-lg">{error}</div>}

            <Button type="submit" className="w-full h-11" disabled={isLoading}>
              {isLoading ? "Вход..." : "Войти"}
            </Button>
          </form>

          <div className="mt-6 text-xs text-muted-foreground text-center">
            <p>Демо аккаунты:</p>
            <p className="mt-1">admin / password123</p>
            <p>user / demo</p>
          </div>

          <div className="mt-4 text-xs text-muted-foreground text-center">
            <span>Нет аккаунта? </span>
            <a href="/register" className="underline hover:text-foreground">
              Зарегистрироваться
            </a>
          </div>
        </div>
      </div>
    </div>
  )
}
