// scripts/view_surveys.js

/**
 * Загружает список опросов из бэкенда и отображает их.
 */
async function loadSurveys() {
    try {
        const response = await httpRequest('GET', '/api/surveys', null);
        const surveysContainer = document.getElementById('surveys-container');

        if (Array.isArray(response)) {
            response.forEach(survey => {
                const surveyDiv = document.createElement('div');
                surveyDiv.className = 'survey';

                surveyDiv.innerHTML = `
                    <h3>${survey.Title}</h3>
                    <p>${survey.Description || 'Без описания'}</p>
                    <p>Создан: ${new Date(survey.CreatedAt).toLocaleString()}</p>
                    <button onclick="viewSurvey(${survey.SurveyID})">Просмотреть</button>
                `;

                surveysContainer.appendChild(surveyDiv);
            });
        } else {
            surveysContainer.innerHTML = '<p>Нет доступных опросов.</p>';
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
    window.location.href = `survey_detail.html?id=${surveyID}`;
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
