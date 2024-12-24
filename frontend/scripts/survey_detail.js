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
            <h2>${survey.title}</h2>
            <p>${survey.description || 'Без описания'}</p>
            <p>Создан: ${new Date(survey.created_at).toLocaleString()}</p>
            <h3>Вопросы</h3>
            <ul>
        `;

        console.log(survey)
        survey.questions.forEach(q => {
            htmlContent += `<li>${q.question_text} (${q.question_type})</li>`;
        });

        htmlContent += `</ul>`;

        // Кнопки для редактирования и аналитики
        htmlContent += `
            <button onclick="editSurvey(${survey.survey_id})">Редактировать Опрос</button>
            <button onclick="viewAnalytics(${survey.survey_id})">Посмотреть Аналитику</button>
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
