/**
 * Обрабатывает ответ от API и выполняет редирект на /login при 401 ошибке
 * @returns true если был выполнен редирект, false в противном случае
 */
export function handleApiResponse(response: Response): boolean {
    if (response.status === 401) {
        // Удаляем токен авторизации из cookies
        if (typeof document !== 'undefined') {
            document.cookie = 'auth_token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT'
        }
        // Выполняем редирект на страницу логина
        if (typeof window !== 'undefined') {
            window.location.href = '/login'
        }
        return true
    }
    return false
}
