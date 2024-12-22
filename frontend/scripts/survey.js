// scripts/survey.js

/**
 * Функция для загрузки опроса по его ID из URL.
 */
async function loadSurvey() {
    const urlParams = new URLSearchParams(window.location.search);
    const surveyID = urlParams.get('id');
    if (!surveyID) {
        showError('Не указан ID опроса.');
        window.location.href = 'index.html';
        return;
    }

    try {
        const survey = await httpRequest('GET', `/survey/${encodeURIComponent(surveyID)}`);
        document.getElementById('survey-title').textContent = survey.title;
        document.getElementById('survey-description').textContent = survey.description;

        const container = document.getElementById('questions-container');
        container.innerHTML = '';

        survey.questions.forEach(function(q, index) {
            const questionDiv = document.createElement('div');
            questionDiv.className = 'question';
            questionDiv.innerHTML = `
                <h3>Вопрос ${index + 1}: ${q.question_text}</h3>
            `;

            if (q.question_type === 'single_choice') {
                q.options.forEach(function(opt) {
                    const optionLabel = document.createElement('label');
                    optionLabel.innerHTML = `
                        <input type="radio" name="question_${q.question_id}" value="${opt.option_text}" required>
                        ${opt.option_text}
                    `;
                    questionDiv.appendChild(optionLabel);
                    questionDiv.appendChild(document.createElement('br'));
                });
            } else if (q.question_type === 'multiple_choice') {
                q.options.forEach(function(opt) {
                    const optionLabel = document.createElement('label');
                    optionLabel.innerHTML = `
                        <input type="checkbox" name="question_${q.question_id}" value="${opt.option_text}">
                        ${opt.option_text}
                    `;
                    questionDiv.appendChild(optionLabel);
                    questionDiv.appendChild(document.createElement('br'));
                });
            } else if (q.question_type === 'free_text') {
                const textarea = document.createElement('textarea');
                textarea.name = `question_${q.question_id}`;
                textarea.rows = 3;
                questionDiv.appendChild(textarea);
            }

            container.appendChild(questionDiv);
        });
    } catch (err) {
        showError(`Ошибка при загрузке опроса: ${err}`);
    }
}

/**
 * Функция для обработки отправки ответов на опрос.
 */
document.addEventListener('DOMContentLoaded', function() {
    loadSurvey();

    const surveyForm = document.getElementById('survey-form');
    if (surveyForm) {
        surveyForm.addEventListener('submit', async function(e) {
            e.preventDefault();
            const urlParams = new URLSearchParams(window.location.search);
            const surveyID = urlParams.get('id');

            const userID = prompt('Введите ваш UserID:');
            if (!userID) {
                showError('UserID обязателен для прохождения опроса.');
                return;
            }

            const answers = {};
            const questions = document.querySelectorAll('#questions-container .question');
            let valid = true;

            for (const qDiv of questions) {
                const h3 = qDiv.querySelector('h3');
                const qIDMatch = h3.textContent.match(/Вопрос\s+\d+:\s+(.*)/);
                const questionText = qIDMatch ? qIDMatch[1] : '';
                const qID = extractQuestionID(questionText); // Функция для извлечения QuestionID

                if (qID === null) {
                    showError('Не удалось определить ID вопроса.');
                    valid = false;
                    break;
                }

                const inputs = qDiv.querySelectorAll('input, textarea');
                if (inputs.length === 0) continue;

                if (qDiv.querySelector('textarea')) {
                    const answerText = qDiv.querySelector('textarea').value.trim();
                    if (answerText === '') {
                        showError('Пожалуйста, ответьте на все вопросы.');
                        valid = false;
                        break;
                    }
                    answers[qID] = { answer_text: answerText };
                } else if (qDiv.querySelector('input[type="radio"]')) {
                    const selected = qDiv.querySelector('input[type="radio"]:checked');
                    if (!selected) {
                        showError('Пожалуйста, ответьте на все вопросы.');
                        valid = false;
                        break;
                    }
                    answers[qID] = { selected_options: [selected.value] };
                } else if (qDiv.querySelector('input[type="checkbox"]')) {
                    const selected = Array.from(qDiv.querySelectorAll('input[type="checkbox"]:checked')).map(cb => cb.value);
                    if (selected.length === 0) {
                        showError('Пожалуйста, ответьте на все вопросы.');
                        valid = false;
                        break;
                    }
                    answers[qID] = { selected_options: selected };
                }
            }

            if (!valid) return;

            const submission = {
                user_id: userID,
                answers: answers
            };

            try {
                const res = await httpRequest('POST', `/survey/${encodeURIComponent(surveyID)}/submit`, submission);
                showSuccess('Опрос успешно пройден');
                window.location.href = 'index.html';
            } catch (err) {
                showError(`Ошибка при отправке ответов: ${err}`);
            }
        });
    }
});

/**
 * Функция для извлечения QuestionID из текста вопроса.
 * В данной реализации возвращает null, так как сервер не предоставляет QuestionID.
 * Необходимо изменить API сервера для передачи QuestionID.
 * @param {string} questionText - Текст вопроса.
 * @returns {number|null} - ID вопроса или null.
 */
function extractQuestionID(questionText) {
    // В текущей реализации сервер не возвращает question_id, поэтому возвращаем null.
    // Для корректной работы необходимо обновить API сервера, чтобы возвращать question_id.
    return null;
}
