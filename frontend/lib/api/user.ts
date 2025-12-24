const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || '/api'

// Интерфейс для профиля пользователя
export interface UserProfile {
    monthlyIncome: number
  }
  
  // Получить профиль пользователя (включая доход)
  export async function fetchUserProfile(): Promise<UserProfile> {
    try {
      const response = await fetch(`${API_BASE_URL}/user/profile`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      })
  
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
  export async function updateMonthlyIncome(monthlyIncome: number): Promise<UserProfile> {
    try {
      const response = await fetch(`${API_BASE_URL}/user/profile`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ monthlyIncome }),
      })
  
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
  