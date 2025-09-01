// Утилиты для работы с авторизацией

// Проверка авторизации
export async function checkAuth() {
  try {
    const response = await fetch('/auth/check', {
      method: 'GET',
      credentials: 'include'
    })
    return response.ok
  } catch (error) {
    console.error('Ошибка проверки авторизации:', error)
    return false
  }
}

// Получение информации о пользователе
export async function getUserInfo() {
  try {
    const response = await fetch('/auth/user', {
      method: 'GET',
      credentials: 'include'
    })
    if (response.ok) {
      return await response.json()
    }
    return null
  } catch (error) {
    console.error('Ошибка получения информации о пользователе:', error)
    return null
  }
}

// Выход из системы
export async function logout() {
  try {
    const response = await fetch('/auth/logout', {
      method: 'GET',
      credentials: 'include'
    })
    return response.ok
  } catch (error) {
    console.error('Ошибка выхода из системы:', error)
    return false
  }
}

// Проверка прав доступа
export async function checkPermission(resource, action) {
  try {
    const response = await fetch(`/auth/check-permission?resource=${resource}&action=${action}`, {
      method: 'GET',
      credentials: 'include'
    })
    return response.ok
  } catch (error) {
    console.error('Ошибка проверки прав доступа:', error)
    return false
  }
} 