// scripts/create_survey.js

let questionCount = 0;

/**
 * Добавляет новый вопрос в форму создания опроса.
 */
function addQuestion() {
    questionCount++;
    const container = document.getElementById('questions-container');

    const questionDiv = document.createElement('div');
    questionDiv.className = 'question';
    questionDiv.id = `question-${questionCount}`;

    questionDiv.innerHTML = `
        <h4>Вопрос ${questionCount}</h4>
        <label for="question_text_${questionCount}">Текст Вопроса:</label>
        <textarea id="question_text_${questionCount}" name="question_text" required></textarea>

        <label for="question_type_${questionCount}">Тип Вопроса:</label>
        <select id="question_type_${questionCount}" name="question_type" onchange="toggleOptions(${questionCount})" required>
            <option value="">--Выберите Тип--</option>
            <option value="single_choice">Один вариант ответа</option>
            <option value="multiple_choice">Несколько вариантов ответа</option>
            <option value="free_text">Свободная форма</option>
        </select>

        <div id="options_container_${questionCount}" style="display:none;">
            <h5>Варианты Ответов</h5>
            <div id="options_list_${questionCount}">
                <!-- Варианты ответов будут добавляться динамически -->
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
 */
function addOption(qID) {
    const optionsList = document.getElementById(`options_list_${qID}`);
    const optionCount = optionsList.childElementCount + 1;

    const optionDiv = document.createElement('div');
    optionDiv.className = 'option';
    optionDiv.id = `question_${qID}_option_${optionCount}`;

    optionDiv.innerHTML = `
        <label for="question_${qID}_option_text_${optionCount}">Вариант ${optionCount}:</label>
        <input type="text" id="question_${qID}_option_text_${optionCount}" name="question_${qID}_option_text" required>
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

// Обработка отправки формы создания опроса
const createSurveyForm = document.getElementById('create-survey-form');
if (createSurveyForm) {
    createSurveyForm.addEventListener('submit', async function(e) {
        e.preventDefault();

        const title = document.getElementById('title').value.trim();
        const description = document.getElementById('description').value.trim();

        const questions = [];
        for (let i = 1; i <= questionCount; i++) {
            const questionText = document.getElementById(`question_text_${i}`);
            const questionType = document.getElementById(`question_type_${i}`);

            if (questionText && questionType && questionText.value.trim() !== "" && questionType.value !== "") {
                const question = {
                    QuestionText: questionText.value.trim(),
                    QuestionType: questionType.value,
                    Options: []
                };

                if (questionType.value === "single_choice" || questionType.value === "multiple_choice") {
                    const optionsList = document.querySelectorAll(`#question-${i} .option input[name="question_${i}_option_text"]`);
                    optionsList.forEach(opt => {
                        if (opt.value.trim() !== "") {
                            question.Options.push({ OptionText: opt.value.trim() });
                        }
                    });

                    if (question.Options.length === 0) {
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
            Title: title,
            Description: description,
            Questions: questions
        };

        try {
            const res = await httpRequest('POST', '/api/surveys', surveyData);
            showSuccess(`Опрос создан с ID: ${res.survey_id}`);
            window.location.href = 'view_surveys.html';
        } catch (err) {
            showError(`Ошибка при создании опроса: ${err.message}`);
        }
    });
}
