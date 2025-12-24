"use client"

import type React from "react"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"

interface PaymentFormProps {
  onSubmit: (payment: {
    name: string
    amount: number
    dueDate: number
    category: string
    color: string
  }) => void
  onCancel: () => void
}

const COLORS = ["#E50914", "#1DB954", "#000000", "#4F46E5", "#F59E0B", "#EF4444", "#06B6D4", "#8B5CF6"]
const CATEGORIES = ["Подписки", "Коммунальные", "Транспорт", "Еда", "Развлечения", "Прочее"]

export default function PaymentForm({ onSubmit, onCancel }: PaymentFormProps) {
  const [form, setForm] = useState({
    name: "",
    amount: "",
    dueDate: "1",
    category: "Подписки",
    color: COLORS[0],
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (form.name && form.amount) {
      onSubmit({
        name: form.name,
        amount: Number.parseFloat(form.amount),
        dueDate: Number.parseInt(form.dueDate),
        category: form.category,
        color: form.color,
      })
      setForm({
        name: "",
        amount: "",
        dueDate: "1",
        category: "Подписки",
        color: COLORS[0],
      })
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4 bg-muted rounded-xl p-4">
      <div>
        <label className="text-sm font-medium text-foreground block mb-2">Название</label>
        <Input
          placeholder="Например: Netflix"
          value={form.name}
          onChange={(e) => setForm({ ...form, name: e.target.value })}
          className="bg-card border-border"
        />
      </div>

      <div>
        <label className="text-sm font-medium text-foreground block mb-2">Сумма (₽)</label>
        <Input
          type="number"
          placeholder="0"
          value={form.amount}
          onChange={(e) => setForm({ ...form, amount: e.target.value })}
          className="bg-card border-border"
        />
      </div>

      <div>
        <label className="text-sm font-medium text-foreground block mb-2">Число месяца</label>
        <Input
          type="number"
          min="1"
          max="31"
          value={form.dueDate}
          onChange={(e) => setForm({ ...form, dueDate: e.target.value })}
          className="bg-card border-border"
        />
      </div>

      <div>
        <label className="text-sm font-medium text-foreground block mb-2">Категория</label>
        <select
          value={form.category}
          onChange={(e) => setForm({ ...form, category: e.target.value })}
          className="w-full rounded-lg bg-card border border-border p-2 text-foreground text-sm"
        >
          {CATEGORIES.map((cat) => (
            <option key={cat} value={cat}>
              {cat}
            </option>
          ))}
        </select>
      </div>

      <div>
        <label className="text-sm font-medium text-foreground block mb-2">Цвет</label>
        <div className="flex gap-2 flex-wrap">
          {COLORS.map((color) => (
            <button
              key={color}
              type="button"
              className={`h-8 w-8 rounded-lg transition-transform ${
                form.color === color ? "ring-2 ring-primary ring-offset-2" : ""
              }`}
              style={{ backgroundColor: color }}
              onClick={() => setForm({ ...form, color })}
            />
          ))}
        </div>
      </div>

      <div className="flex gap-2 pt-4">
        <Button type="submit" className="flex-1">
          Добавить
        </Button>
        <Button type="button" variant="secondary" onClick={onCancel} className="flex-1">
          Отмена
        </Button>
      </div>
    </form>
  )
}
