import { handleApiResponse } from './utils'

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || '/api'

export interface UserCredentials {
    login: string
    password: string
}

export interface LoginResponse {
    success: boolean
    message?: string
}

export interface UpdateMonthlyIncomeRequest {
    income: number
}

export async function register(credentials: UserCredentials): Promise<LoginResponse> {
    try {
        const response = await fetch(`${API_BASE_URL}/user/register`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(credentials),
            credentials: "include", // Включаем cookies
        })

        const data = await response.json()

        if (!response.ok || data.success === false) {
            return {
                success: false,
                message: data.message || "Ошибка регистрации",
            }
        }

        return {
            success: true,
        }
    } catch (error) {
        return {
            success: false,
            message: "Ошибка соединения с сервером",
        }
    }
}

export async function login(credentials: UserCredentials): Promise<LoginResponse> {
    try {
        const response = await fetch(`${API_BASE_URL}/user/login`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(credentials),
            credentials: "include", // Включаем cookies
        })

        const data = await response.json()

        if (!response.ok || data.success === false) {
            return {
                success: false,
                message: data.message || "Ошибка авторизации",
            }
        }

        return {
            success: true,
        }
    } catch (error) {
        return {
            success: false,
            message: "Ошибка соединения с сервером",
        }
    }
}

export async function logout(): Promise<void> {
    await fetch(`${API_BASE_URL}/user/logout`, {
        method: "POST",
        credentials: "include",
    })
}

// Интерфейс для профиля пользователя
export interface UserProfile {
    monthlyIncome: number
}

// Получить профиль пользователя (включая доход)
export async function fetchUserProfile(): Promise<UserProfile> {
    try {
        const response = await fetch(`${API_BASE_URL}/user/profile/fetch`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
        })

        handleApiResponse(response)

        if (!response.ok) {
            throw new Error(`Ошибка загрузки профиля: ${response.statusText}`)
        }

        const data = await response.json()
        return data
    } catch (error) {
        console.error('Ошибка при загрузке профиля:', error)
        throw error
    }
}

// Обновить доход пользователя
export async function updateMonthlyIncome(monthlyIncome: UpdateMonthlyIncomeRequest): Promise<UserProfile> {
    try {
        const response = await fetch(`${API_BASE_URL}/user/profile/update`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
            body: JSON.stringify(monthlyIncome),
        })

        handleApiResponse(response)

        if (!response.ok) {
            throw new Error(`Ошибка обновления дохода: ${response.statusText}`)
        }

        const data = await response.json()
        return data
    } catch (error) {
        console.error('Ошибка при обновлении дохода:', error)
        throw error
    }
}