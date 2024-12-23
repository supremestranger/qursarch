// scripts/utils.js

/**
 * Выполняет HTTP-запрос с использованием Fetch API.
 * @param {string} method - HTTP метод (GET, POST, PUT, DELETE).
 * @param {string} url - URL для запроса.
 * @param {Object} data - Тело запроса в формате JSON.
 * @returns {Promise<Object>} - Данные ответа в формате JSON.
 */
async function httpRequest(method, url, data) {
    const options = {
        method: method,
        headers: {
            'Content-Type': 'application/json'
        },
        credentials: 'include' // Включает куки в запрос
    };

    if (data) {
        options.body = JSON.stringify(data);
    }

    const response = await fetch(url, options);

    if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Ошибка запроса');
    }

    return response.json();
}

/**
 * Отображает сообщение об ошибке пользователю.
 * @param {string} message - Сообщение об ошибке.
 */
function showError(message) {
    alert(`Ошибка: ${message}`);
}

/**
 * Отображает сообщение об успехе пользователю.
 * @param {string} message - Сообщение об успехе.
 */
function showSuccess(message) {
    alert(`Успех: ${message}`);
}
