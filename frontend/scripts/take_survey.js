// scripts/take_survey.js

/**
 * Загружает опрос по переданному ID и отображает его.
 */
async function loadSurvey() {
    const urlParams = new URLSearchParams(window.location.search);
    const surveyID = urlParams.get('id');

    if (!surveyID) {
        showError('Не указан ID опроса.');
        return;
    }

    try {
        const survey = await httpRequest('GET', `/api/public_surveys/${surveyID}`, null);
        displaySurvey(survey);
    } catch (err) {
        showError(`Ошибка при загрузке опроса: ${err.message}`);
    }
}

/**
 * Отображает опрос на странице.
 * @param {Object} survey - Объект опроса.
 */
function displaySurvey(survey) {
    const surveyContainer = document.getElementById('survey-container');
    if (!surveyContainer) return;

    let htmlContent = `
        <h2>${sanitizeHTML(survey.title)}</h2>
        <p>${sanitizeHTML(survey.description || 'Без описания')}</p>
        <form id="submit-survey-form">
    `;

    survey.questions.forEach((question, index) => {
        htmlContent += `
            <div class="question">
                <h3>${index + 1}. ${sanitizeHTML(question.question_text)}</h3>
        `;

        if (question.question_type === 'single_choice') {
            question.options.forEach(option => {
                htmlContent += `
                    <div class="option">
                        <input type="radio" id="option_${question.question_id}_${option.option_id}" name="question_${question.question_id}" value="${sanitizeHTML(option.option_id)}" required>
                        <label for="option_${question.question_id}_${option.option_id}">${sanitizeHTML(option.option_text)}</label>
                    </div>
                `;
            });
        } else if (question.question_type === 'multiple_choice') {
            question.options.forEach(option => {
                htmlContent += `
                    <div class="option">
                        <input type="checkbox" id="option_${question.question_id}_${option.option_id}" name="question_${question.question_id}" value="${sanitizeHTML(option.option_id)}">
                        <label for="option_${question.question_id}_${option.option_id}">${sanitizeHTML(option.option_text)}</label>
                    </div>
                `;
            });
        } else if (question.question_type === 'free_text') {
            htmlContent += `
                <textarea id="question_${question.question_id}" name="question_${question.question_id}" required></textarea>
            `;
        }

        htmlContent += `
            </div>
        `;
    });

    htmlContent += `
            <button type="submit">Отправить Ответы</button>
        </form>
    `;

    surveyContainer.innerHTML = htmlContent;

    // Добавление обработчика отправки формы
    const submitSurveyForm = document.getElementById('submit-survey-form');
    if (submitSurveyForm) {
        submitSurveyForm.addEventListener('submit', submitSurvey);
    }
}

/**
 * Обрабатывает отправку опроса и сохраняет результаты на сервере.
 * @param {Event} event - Событие отправки формы.
 */
async function submitSurvey(event) {
    event.preventDefault();

    const urlParams = new URLSearchParams(window.location.search);
    const surveyID = urlParams.get('id');

    if (!surveyID) {
        showError('Не указан ID опроса.');
        return;
    }

    // Сбор ответов
    const surveyData = {
        answers: [],
        user_id: Math.ceil(Math.random() * 999)
    };

    const form = event.target;
    const formData = new FormData(form);

    // Получение уникальных QuestionIDs
    const questionIDs = new Set();
    for (let key of formData.keys()) {
        const match = key.match(/^question_(\d+)$/);
        if (match) {
            questionIDs.add(match[1]);
        }
    }

    // Сбор ответов по каждому вопросу
    questionIDs.forEach(qID => {
        const questionType = getQuestionType(qID);
        if (questionType === 'single_choice') {
            const selectedOptions = formData.getAll(`question_${qID}`);
            selectedOptions2 = {}
            selectedOptions.forEach(element => {
                console.log(element);
                selectedOptions["option_text"] = element
            });
            if (selectedOptions) {
                surveyData.answers.push({
                    question_id: parseInt(qID),
                    answer_text: "",
                    selected_options: selectedOptions2
                });
            }
        } else if (questionType === 'multiple_choice') {
            const selectedOptions = formData.getAll(`question_${qID}`);
            selectedOptions2 = {}
            selectedOptions.forEach(element => {
                selectedOptions["option_text"] = element
            });
            surveyData.answers.push({
                question_id: parseInt(qID),
                answer_text: "",
                selected_options: selectedOptions2
            });
        } else if (questionType === 'free_text') {
            const answerText = formData.get(`question_${qID}`);
            if (answerText) {
                surveyData.answers.push({
                    question_id: parseInt(qID),
                    answer_text: answerText.toString(),
                    selected_options: ""
                });
            }
        }
    });

    // Отправка данных на сервер
    try {
        console.log(surveyData)
        const res = await httpRequest('POST', `/api/public_surveys/${surveyID}/submit`, surveyData);
        showSuccess(res.message);
        window.location.href = 'index.html';
    } catch (err) {
        showError(`Ошибка при отправке опроса: ${err.message}`);
    }
}

/**
 * Определяет тип вопроса по его ID.
 * @param {string} questionID - ID вопроса.
 * @returns {string} - Тип вопроса.
 */
function getQuestionType(questionID) {
    // Здесь необходимо реализовать логику для определения типа вопроса.
    // Например, можно хранить типы вопросов в data-атрибутах или получить их с сервера.
    // В данном примере предположим, что типы вопросов известны заранее.
    // Для полноценной реализации рекомендуется получать информацию о типах вопросов с сервера.
    return 'single_choice'; // Замените на реальную логику
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

    // Загрузка опроса при загрузке страницы
    loadSurvey();
});
