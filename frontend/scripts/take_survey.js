// scripts/take_survey.js

let survey = {};

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
        const fetchedSurvey = await httpRequest('GET', `/api/surveys/${surveyID}`, null);
        survey = fetchedSurvey;
        displaySurvey(fetchedSurvey);
    } catch (err) {
        showError(`Ошибка при загрузке опроса: ${err.message}`);
    }
}

/**
 * Отображает опрос на странице.
 * @param {Object} survey - Объект опроса.
 */
function displaySurvey(survey) {
    document.getElementById('survey-title').innerText = survey.title;

    const surveyForm = document.getElementById('survey-form');

    survey.questions.forEach((question, index) => {
        const questionDiv = document.createElement('div');
        questionDiv.className = 'question';
        questionDiv.id = `question-${index + 1}`;

        const questionLabel = document.createElement('label');
        questionLabel.innerText = `${index + 1}. ${question.question_text}`;
        questionLabel.htmlFor = `question-input-${index + 1}`;
        questionDiv.appendChild(questionLabel);

        if (question.question_type === 'single_choice') {
            question.options.forEach((option, optIndex) => {
                const optionDiv = document.createElement('div');
                optionDiv.className = 'option';

                const optionInput = document.createElement('input');
                optionInput.type = 'radio';
                optionInput.id = `question-${index + 1}-option-${optIndex + 1}`;
                optionInput.name = `question-${index + 1}`;
                optionInput.value = option.option_id; // Используется OptionID

                const optionLabel = document.createElement('label');
                optionLabel.htmlFor = `question-${index + 1}-option-${optIndex + 1}`;
                optionLabel.innerText = option.option_text;

                optionDiv.appendChild(optionInput);
                optionDiv.appendChild(optionLabel);
                questionDiv.appendChild(optionDiv);
            });
        } else if (question.question_type === 'multiple_choice') {
            question.options.forEach((option, optIndex) => {
                const optionDiv = document.createElement('div');
                optionDiv.className = 'option';

                const optionInput = document.createElement('input');
                optionInput.type = 'checkbox';
                optionInput.id = `question-${index + 1}-option-${optIndex + 1}`;
                optionInput.name = `question-${index + 1}`;
                optionInput.value = option.option_id; // Используется OptionID

                const optionLabel = document.createElement('label');
                optionLabel.htmlFor = `question-${index + 1}-option-${optIndex + 1}`;
                optionLabel.innerText = option.option_text;

                optionDiv.appendChild(optionInput);
                optionDiv.appendChild(optionLabel);
                questionDiv.appendChild(optionDiv);
            });
        } else if (question.question_type === 'free_text') {
            const textarea = document.createElement('textarea');
            textarea.id = `question-${index + 1}-input`;
            textarea.name = `question-${index + 1}`;
            textarea.required = true;
            questionDiv.appendChild(textarea);
        }

        // Вставка вопроса перед кнопкой отправки
        surveyForm.insertBefore(questionDiv, surveyForm.lastElementChild);
    });
}

/**
 * Обрабатывает отправку формы опроса.
 */
async function submitSurvey(event) {
    event.preventDefault();

    const urlParams = new URLSearchParams(window.location.search);
    const surveyID = urlParams.get('id');

    if (!surveyID) {
        showError('Не указан ID опроса.');
        return;
    }

    const formData = new FormData(event.target);
    const answers = [];

    // Сбор ответов
    survey.questions.forEach((question, index) => {
        const key = `question-${index + 1}`;
        if (question.question_type === 'free_text') {
            const answerText = formData.get(key).trim();
            if (answerText === "") {
                showError(`Пожалуйста, ответьте на вопрос ${index + 1}.`);
                throw new Error(`Пустой ответ на вопрос ${index + 1}.`);
            }
            answers.push({
                question_id: question.question_id,
                answer_text: answerText,
                selected_options: []
            });
        } else if (question.question_type === 'single_choice') {
            const selectedOption = formData.get(key);
            if (!selectedOption) {
                showError(`Пожалуйста, выберите вариант ответа для вопроса ${index + 1}.`);
                throw new Error(`Не выбран вариант ответа для вопроса ${index + 1}.`);
            }
            answers.push({
                question_id: question.question_id,
                answer_text: "",
                selected_options: [parseInt(selectedOption, 10)]
            });
        } else if (question.question_type === 'multiple_choice') {
            const selectedOptions = formData.getAll(key).map(v => parseInt(v, 10));
            if (selectedOptions.length === 0) {
                showError(`Пожалуйста, выберите хотя бы один вариант ответа для вопроса ${index + 1}.`);
                throw new Error(`Не выбраны варианты ответов для вопроса ${index + 1}.`);
            }
            answers.push({
                question_id: question.question_id,
                answer_text: "",
                selected_options: selectedOptions
            });
        }
    });

    const surveyResult = {
        survey_id: parseInt(surveyID, 10),
        answers: answers
    };

    try {
        const res = await httpRequest('POST', `/api/surveys/${surveyID}/submit`, surveyResult);
        showSuccess('Спасибо за прохождение опроса!');
        window.location.href = 'index.html';
    } catch (err) {
        showError(`Ошибка при отправке ответов: ${err.message}`);
    }
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

    // Обработка отправки формы опроса
    const surveyForm = document.getElementById('survey-form');
    if (surveyForm) {
        surveyForm.addEventListener('submit', submitSurvey);
    }

    // Проверка аутентификации для отображения кнопки выхода
    const token = localStorage.getItem('authToken');
    if (token) {
        document.getElementById('logout-button').style.display = 'inline-block';
    }
});
