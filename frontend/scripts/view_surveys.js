// scripts/view_surveys.js

/**
 * Загружает список опросов из бэкенда и отображает их.
 */
async function loadSurveys() {
    try {
        // Отправка GET-запроса для получения списка опросов
        const response = await httpRequest('GET', '/api/surveys', null);
        const surveysContainer = document.getElementById('surveys-container');

        // Очистка контейнера перед отображением новых опросов
        surveysContainer.innerHTML = '';

        if (Array.isArray(response)) {
            if (response.length === 0) {
                surveysContainer.innerHTML = '<p>Нет доступных опросов.</p>';
                return;
            }

            response.forEach(survey => {
                const surveyDiv = document.createElement('div');
                surveyDiv.className = 'survey';

                // Создание HTML-контента для каждого опроса
                surveyDiv.innerHTML = `
                    <h3>${sanitizeHTML(survey.title)}</h3>
                    <p>${sanitizeHTML(survey.description || 'Без описания')}</p>
                    <p>Создан: ${new Date(survey.created_at).toLocaleString()}</p>
                    <button onclick="viewSurvey(${survey.survey_id})">Просмотреть</button>
                `;

                surveysContainer.appendChild(surveyDiv);
            });
        } else {
            surveysContainer.innerHTML = '<p>Не удалось загрузить опросы.</p>';
        }
    } catch (err) {
        showError(`Ошибка при загрузке опросов: ${err.message}`);
    }
}

/**
 * Переход к просмотру конкретного опроса.
 * @param {number} surveyID - ID опроса.
 */
function viewSurvey(surveyID) {
    window.location.href = `survey_detail.html?id=${encodeURIComponent(surveyID)}`;
}

/**
 * Функция для безопасного вывода HTML-содержимого, предотвращающая XSS атаки.
 * @param {string} str - Входная строка.
 * @returns {string} - Очищенная строка.
 */
function sanitizeHTML(str) {
    const temp = document.createElement('div');
    temp.textContent = str;
    return temp.innerHTML;
}

// Обработка выхода
document.addEventListener('DOMContentLoaded', function() {
    const logoutButton = document.getElementById('logout-button');
    if (logoutButton) {
        logoutButton.addEventListener('click', async function() {
            try {
                const res = await httpRequest('POST', '/api/logout', null);
                showSuccess(res.message);
                window.location.href = 'index.html';
            } catch (err) {
                showError(`Ошибка при выходе: ${err.message}`);
            }
        });
    }

    // Загрузка опросов при загрузке страницы
    loadSurveys();
});
