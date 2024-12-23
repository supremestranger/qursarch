// scripts/survey_detail.js

/**
 * Загружает детали опроса и отображает его.
 */
async function loadSurveyDetails() {
    const urlParams = new URLSearchParams(window.location.search);
    const surveyID = urlParams.get('id');

    if (!surveyID) {
        showError('Не указан ID опроса.');
        return;
    }

    try {
        const survey = await httpRequest('GET', `/api/surveys/${surveyID}`, null);
        const detailsContainer = document.getElementById('survey-details');

        let htmlContent = `
            <h2>${survey.Title}</h2>
            <p>${survey.Description || 'Без описания'}</p>
            <p>Создан: ${new Date(survey.CreatedAt).toLocaleString()}</p>
            <h3>Вопросы</h3>
            <ul>
        `;

        survey.Questions.forEach(q => {
            htmlContent += `<li>${q.QuestionText} (${q.QuestionType})</li>`;
        });

        htmlContent += `</ul>`;

        // Кнопки для редактирования и аналитики
        htmlContent += `
            <button onclick="editSurvey(${survey.SurveyID})">Редактировать Опрос</button>
            <button onclick="viewAnalytics(${survey.SurveyID})">Посмотреть Аналитику</button>
        `;

        detailsContainer.innerHTML = htmlContent;
    } catch (err) {
        showError(`Ошибка при загрузке деталей опроса: ${err.message}`);
    }
}

/**
 * Переход к редактированию опроса.
 * @param {number} surveyID - ID опроса.
 */
function editSurvey(surveyID) {
    window.location.href = `edit_survey.html?id=${encodeURIComponent(surveyID)}`;
}

/**
 * Переход к странице аналитики опроса.
 * @param {number} surveyID - ID опроса.
 */
function viewAnalytics(surveyID) {
    window.location.href = `analytics.html?id=${encodeURIComponent(surveyID)}`;
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

    // Загрузка деталей опроса при загрузке страницы
    loadSurveyDetails();
});
