Анализирую проект на хардкод и бизнес-логику, которую стоит вынести на бекенд. Просматриваю основные файлы:


[5 tools called]

Ищу дополнительные места с хардкодом и бизнес-логикой:


[2 tools called]

Нашёл все места с хардкодом и бизнес-логикой, которую стоит вынести на бекенд. Отчёт ниже.

## Детальный анализ хардкода и бизнес-логики

---

## 1. Хардкод: месячный доход (monthlyIncome)

### Место в коде:
```typescript:15:15:app/page.tsx
const [monthlyIncome, setMonthlyIncome] = useState(50000)
```

Проблема: значение дохода захардкожено на фронтенде.

Что нужно сделать:
- Вынести в бекенд (API: `GET /api/user/profile`, `PUT /api/user/profile`)
- Хранить в базе данных пользователя
- Загружать при инициализации приложения

---

## 2. Хардкод: категории платежей

### Место в коде:
```typescript:21:21:components/payment-form.tsx
const CATEGORIES = ["Подписки", "Коммунальные", "Транспорт", "Еда", "Развлечения", "Прочее"]
```

Проблема: список категорий статичен на фронтенде.

Что нужно сделать:
- API: `GET /api/categories`
- Хранить в БД, позволить добавлять/редактировать
- Возвращать с локализацией, иконками и настройками пользователя

---

## 3. Хардкод: цвета для платежей

### Место в коде:
```typescript:20:20:components/payment-form.tsx
const COLORS = ["#E50914", "#1DB954", "#000000", "#4F46E5", "#F59E0B", "#EF4444", "#06B6D4", "#8B5CF6"]
```

И также:
```typescript:38:38:components/spending-chart.tsx
const CHART_COLORS = ["#4F46E5", "#06B6D4", "#F59E0B", "#EF4444", "#8B5CF6", "#EC4899"]
```

Проблема: цвета захардкожены в двух местах.

Что нужно сделать:
- API: `GET /api/config/colors` или возвращать с категориями
- Хранить в конфигурации или разрешать пользователю настраивать

---

## 4. Бизнес-логика: расчёт дней до платежа

### Место в коде:
```typescript:77:81:app/page.tsx
const getDaysUntilPayment = (dueDate: number) => {
  const today = new Date()
  const currentDay = today.getDate()
  return dueDate >= currentDay ? dueDate - currentDay : dueDate + (31 - currentDay)
}
```

И также:
```typescript:34:34:components/upcoming-payments.tsx
daysUntil: p.dueDate >= currentDay ? p.dueDate - currentDay : p.dueDate + (31 - currentDay),
```

Проблема:
- Дублирование логики
- Хардкод: `31` (дней в месяце)
- Не учитывает разные месяцы и високосные годы
- Логика на клиенте

Что нужно сделать:
- Вынести на бекенд
- API: вычислять `daysUntil` при возврате платежей
- Использовать корректные даты (учитывать месяцы и годы)

---

## 5. Бизнес-логика: проверка "платёж на этой неделе"

### Место в коде:
```typescript:83:86:app/page.tsx
const isUpcomingThisWeek = (dueDate: number) => {
  const daysUntil = getDaysUntilPayment(dueDate)
  return daysUntil <= 7 && daysUntil >= 0
}
```

Проблема: бизнес-правило на фронтенде, дублирование.

Что нужно сделать:
- Вынести на бекенд
- API: возвращать флаг `isUpcomingThisWeek` в данных платежа
- Настраивать порог (например, 7 дней) через конфигурацию

---

## 6. Бизнес-логика: расчёт общей суммы расходов

### Место в коде:
```typescript:88:88:app/page.tsx
const totalExpenses = payments.reduce((sum, p) => sum + p.amount, 0)
```

И также:
```typescript:45:45:components/dashboard.tsx
const expensePercentage = (totalExpenses / monthlyIncome) * 100
```

Проблема: расчёты на клиенте, дублируются.

Что нужно сделать:
- API: `GET /api/statistics` с агрегатами
- Возвращать: `totalExpenses`, `expensePercentage`, `remaining`
- Кэшировать для производительности

---

## 7. Бизнес-логика: группировка по категориям для графика

### Место в коде:
```typescript:22:36:components/spending-chart.tsx
const categoryData = payments.reduce(
  (acc, payment) => {
    const existing = acc.find((c) => c.name === payment.category)
    if (existing) {
      existing.value += payment.amount
    } else {
      acc.push({
        name: payment.category,
        value: payment.amount,
      })
    }
    return acc
  },
  [] as Array<{ name: string; value: number }>,
)
```

Проблема: группировка и агрегация на клиенте.

Что нужно сделать:
- API: `GET /api/statistics/by-category`
- Выполнять агрегацию на бекенде
- Возвращать готовые данные для графика

---

## 8. Бизнес-логика: сортировка платежей

### Место в коде:
```typescript:175:175:app/page.tsx
.sort((a, b) => a.dueDate - b.dueDate)
```

И также:
```typescript:36:36:components/upcoming-payments.tsx
.sort((a, b) => a.daysUntil - b.daysUntil)
```

Проблема: сортировка на клиенте, разные критерии.

Что нужно сделать:
- API: параметр `?sortBy=dueDate` или `?sortBy=daysUntil`
- Сортировка на бекенде
- Единообразие и производительность

---

## 9. Хардкод: значения по умолчанию в форме

### Место в коде:
```typescript:24:30:components/payment-form.tsx
const [form, setForm] = useState({
  name: "",
  amount: "",
  dueDate: "1",
  category: "Подписки",
  color: COLORS[0],
})
```

Проблема: значения по умолчанию захардкожены.

Что нужно сделать:
- API: `GET /api/config/defaults`
- Возвращать с учётом пользовательских настроек и последних значений

---

## 10. Бизнес-логика: валидация формы

### Место в коде:
```typescript:32:50:components/payment-form.tsx
const handleSubmit = (e: React.FormEvent) => {
  e.preventDefault()
  if (form.name && form.amount) {
    // ...
  }
}
```

И также:
```typescript:74:80:app/api/payments/route.ts
if (!name || !amount || !dueDate || !category || !color) {
  return NextResponse.json(
    { error: 'Не все обязательные поля заполнены' },
    { status: 400 }
  )
}
```

Проблема: валидация дублируется на клиенте и бекенде, правила на клиенте.

Что нужно сделать:
- Валидацию на бекенде (схема Zod/Yup)
- API: возвращать ошибки валидации
- На клиенте — только базовая проверка для UX

---

## 11. Бизнес-логика: расчёт остатка

### Место в коде:
```typescript:89:89:app/page.tsx
const remaining = monthlyIncome - totalExpenses
```

Проблема: простой расчёт на клиенте.

Что нужно сделать:
- API: `GET /api/statistics`
- Возвращать вычисленный `remaining`
- Учитывать дополнительные правила, если появятся

---

## 12. Хардкод: магические числа

- `31` — дни в месяце (строки 80, 34)
- `7` — дни для "на этой неделе" (строка 85)
- `3` — дни для предупреждения (строка 75 в upcoming-payments)
- `4` — количество платежей в превью (строка 38)

Что нужно сделать:
- Вынести в конфигурацию бекенда
- API: `GET /api/config` с настройками

---

## Итоговая таблица рекомендаций

| Место | Тип | Приоритет | API эндпоинт |
|-------|-----|-----------|--------------|
| monthlyIncome | Хардкод | Высокий | `GET/PUT /api/user/profile` |
| CATEGORIES | Хардкод | Высокий | `GET /api/categories` |
| COLORS | Хардкод | Средний | `GET /api/config/colors` |
| getDaysUntilPayment | Логика | Высокий | Вычислять на бекенде |
| totalExpenses | Логика | Высокий | `GET /api/statistics` |
| categoryData (группировка) | Логика | Высокий | `GET /api/statistics/by-category` |
| Сортировка | Логика | Средний | Параметры запроса |
| Валидация | Логика | Высокий | Схема валидации на бекенде |
| Магические числа | Хардкод | Средний | `GET /api/config` |

---

## Рекомендуемая архитектура бекенда

1. `/api/user/profile` — профиль пользователя (доход, настройки)
2. `/api/categories` — категории платежей (CRUD)
3. `/api/config` — конфигурация (цвета, значения по умолчанию, пороги)
4. `/api/statistics` — статистика (расходы, остаток, проценты)
5. `/api/statistics/by-category` — группировка по категориям
6. `/api/payments` — платежи (уже есть, добавить вычисляемые поля)

Это упростит поддержку, улучшит производительность и централизует бизнес-логику.

Могу подготовить план миграции или начать с конкретного раздела.