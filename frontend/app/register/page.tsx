"use client"

import type React from "react"

import { useState } from "react"
import { useRouter } from "next/navigation"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { register } from "@/lib/api/user"

export default function RegisterPage() {
  const router = useRouter()
  const [loginValue, setLoginValue] = useState("")
  const [password, setPassword] = useState("")
  const [error, setError] = useState("")
  const [isLoading, setIsLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError("")
    setIsLoading(true)

    const result = await register({ login: loginValue, password })

    if (result.success) {
      router.push("/login")
      router.refresh()
    } else {
      setError(result.message || "Ошибка регистрации")
    }

    setIsLoading(false)
  }

  return (
    <div className="min-h-screen bg-background flex items-center justify-center p-4">
      <div className="w-full max-w-sm">
        <div className="bg-card rounded-2xl shadow-sm border border-border p-8">
          <h1 className="text-2xl font-semibold text-foreground mb-2">Регистрация</h1>
          <p className="text-sm text-muted-foreground mb-6">Создайте новый аккаунт</p>

          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <label htmlFor="login" className="text-sm font-medium text-foreground">
                Логин
              </label>
              <Input
                id="login"
                type="text"
                value={loginValue}
                onChange={(e) => setLoginValue(e.target.value)}
                placeholder="Придумайте логин"
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
                placeholder="Придумайте пароль"
                required
                disabled={isLoading}
                className="h-11"
              />
            </div>

            {error && (
              <div className="text-sm text-red-500 bg-red-50 dark:bg-red-950/20 p-3 rounded-lg">
                {error}
              </div>
            )}

            <Button type="submit" className="w-full h-11" disabled={isLoading}>
              {isLoading ? "Регистрация..." : "Зарегистрироваться"}
            </Button>
          </form>

          <div className="mt-4 text-xs text-muted-foreground text-center">
            <span>Уже есть аккаунт? </span>
            <a href="/login" className="underline hover:text-foreground">
              Войти
            </a>
          </div>
        </div>
      </div>
    </div>
  )
}

