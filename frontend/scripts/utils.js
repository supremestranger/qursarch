// scripts/utils.js

/**
 * Функция для отправки HTTP-запросов с использованием fetch().
 * @param {string} method - HTTP метод (GET, POST, PUT, DELETE).
 * @param {string} url - URL запроса.
 * @param {object|null} data - Данные для отправки (для POST, PUT).
 * @returns {Promise<object>} - Возвращает Promise с результатом запроса.
 */
async function httpRequest(method, url, data = null) {
    const options = {
        method: method,
        headers: {}
    };

    if (data) {
        options.headers['Content-Type'] = 'application/json';
        options.body = JSON.stringify(data);
    }

    try {
        const response = await fetch(url, options);
        const responseData = await response.json();
        if (!response.ok) {
            throw new Error(responseData || 'Ошибка сети');
        }
        return responseData;
    } catch (error) {
        throw error.message || 'Ошибка сети';
    }
}

/**
 * Функция для отображения сообщений об ошибках.
 * @param {string} message - Сообщение об ошибке.
 */
function showError(message) {
    alert(message);
}

/**
 * Функция для отображения успешных сообщений.
 * @param {string} message - Успешное сообщение.
 */
function showSuccess(message) {
    alert(message);
}
