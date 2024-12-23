// scripts/create_survey.js

let questionCount = 0;

/**
 * Добавляет новый вопрос в форму создания опроса.
 */
function addQuestion(question = null) {
    questionCount++;
    const container = document.getElementById('questions-container');

    const questionDiv = document.createElement('div');
    questionDiv.className = 'question';
    questionDiv.id = `question-${questionCount}`;

    questionDiv.innerHTML = `
        <h4>Вопрос ${questionCount}</h4>
        <label for="question_text_${questionCount}">Текст Вопроса:</label>
        <textarea id="question_text_${questionCount}" name="question_text" required>${question ? question.question_text : ''}</textarea>

        <label for="question_type_${questionCount}">Тип Вопроса:</label>
        <select id="question_type_${questionCount}" name="question_type" onchange="toggleOptions(${questionCount})" required>
            <option value="">--Выберите Тип--</option>
            <option value="single_choice" ${question && question.question_type === 'single_choice' ? 'selected' : ''}>Один вариант ответа</option>
            <option value="multiple_choice" ${question && question.question_type === 'multiple_choice' ? 'selected' : ''}>Несколько вариантов ответа</option>
            <option value="free_text" ${question && question.question_type === 'free_text' ? 'selected' : ''}>Свободная форма</option>
        </select>

        <div id="options_container_${questionCount}" style="display:${question && (question.question_type === 'single_choice' || question.question_type === 'multiple_choice') ? 'block' : 'none'};">
            <h5>Варианты Ответов</h5>
            <div id="options_list_${questionCount}">
                ${question && (question.question_type === 'single_choice' || question.question_type === 'multiple_choice') ? 
                    question.options.map((opt, idx) => `
                        <div class="option" id="question_${questionCount}_option_${idx + 1}">
                            <label for="question_${questionCount}_option_text_${idx + 1}">Вариант ${idx + 1}:</label>
                            <input type="text" id="question_${questionCount}_option_text_${idx + 1}" name="question_${questionCount}_option_text" required value="${opt.option_text}">
                            <button type="button" onclick="removeOption(${questionCount}, ${idx + 1})">Удалить Вариант</button>
                        </div>
                    `).join('') : ''}
            </div>
            <button type="button" onclick="addOption(${questionCount})">Добавить Вариант</button>
        </div>

        <button type="button" onclick="removeQuestion(${questionCount})">Удалить Вопрос</button>
    `;

    container.appendChild(questionDiv);
}

/**
 * Добавляет вариант ответа к заданному вопросу.
 * @param {number} qID - ID вопроса.
 * @param {string} [optionText] - Предварительно заполненный текст варианта (при редактировании).
 */
function addOption(qID, optionText = '') {
    const optionsList = document.getElementById(`options_list_${qID}`);
    const optionCount = optionsList.childElementCount + 1;

    const optionDiv = document.createElement('div');
    optionDiv.className = 'option';
    optionDiv.id = `question_${qID}_option_${optionCount}`;

    optionDiv.innerHTML = `
        <label for="question_${qID}_option_text_${optionCount}">Вариант ${optionCount}:</label>
        <input type="text" id="question_${qID}_option_text_${optionCount}" name="question_${qID}_option_text" required value="${optionText}">
        <button type="button" onclick="removeOption(${qID}, ${optionCount})">Удалить Вариант</button>
    `;

    optionsList.appendChild(optionDiv);
}

/**
 * Удаляет вариант ответа из вопроса.
 * @param {number} qID - ID вопроса.
 * @param {number} oID - ID варианта ответа.
 */
function removeOption(qID, oID) {
    const optionDiv = document.getElementById(`question_${qID}_option_${oID}`);
    if (optionDiv) {
        optionDiv.remove();
    }
}

/**
 * Удаляет вопрос из формы создания опроса.
 * @param {number} qID - ID вопроса.
 */
function removeQuestion(qID) {
    const questionDiv = document.getElementById(`question-${qID}`);
    if (questionDiv) {
        questionDiv.remove();
    }
}

/**
 * Показывает или скрывает контейнер с вариантами ответов в зависимости от выбранного типа вопроса.
 * @param {number} qID - ID вопроса.
 */
function toggleOptions(qID) {
    const select = document.getElementById(`question_type_${qID}`);
    const selected = select.value;
    const optionsContainer = document.getElementById(`options_container_${qID}`);

    if (selected === "single_choice" || selected === "multiple_choice") {
        optionsContainer.style.display = 'block';
    } else {
        optionsContainer.style.display = 'none';
    }
}

/**
 * Обрабатывает отправку формы создания опроса.
 */
async function submitCreateSurvey(event) {
    event.preventDefault();

    const title = document.getElementById('title').value.trim();
    const description = document.getElementById('description').value.trim();

    const questions = [];
    for (let i = 1; i <= questionCount; i++) {
        const questionText = document.getElementById(`question_text_${i}`);
        const questionType = document.getElementById(`question_type_${i}`);

        if (questionText && questionType && questionText.value.trim() !== "" && questionType.value !== "") {
            const question = {
                question_text: questionText.value.trim(),
                question_type: questionType.value,
                options: []
            };

            if (questionType.value === "single_choice" || questionType.value === "multiple_choice") {
                const optionsList = document.querySelectorAll(`#question-${i} .option input[name="question_${i}_option_text"]`);
                optionsList.forEach(opt => {
                    if (opt.value.trim() !== "") {
                        question.options.push({ option_text: opt.value.trim() });
                    }
                });

                if (question.options.length === 0) {
                    showError(`Пожалуйста, добавьте хотя бы один вариант ответа для вопроса ${i}.`);
                    return;
                }
            }

            questions.push(question);
        }
    }

    if (questions.length === 0) {
        showError("Пожалуйста, добавьте хотя бы один вопрос.");
        return;
    }

    const surveyData = {
        title: title,
        description: description,
        questions: questions
    };

    try {
        const res = await httpRequest('POST', '/api/surveys', surveyData);
        showSuccess(`Опрос создан с ID: ${res.survey_id}`);
        window.location.href = 'view_surveys.html';
    } catch (err) {
        showError(`Ошибка при создании опроса: ${err.message}`);
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

    // Обработка отправки формы создания опроса
    const createSurveyForm = document.getElementById('create-survey-form');
    if (createSurveyForm) {
        createSurveyForm.addEventListener('submit', submitCreateSurvey);
    }
}
);
