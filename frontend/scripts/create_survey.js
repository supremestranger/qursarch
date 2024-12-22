// scripts/create_survey.js

/**
 * Функция для добавления нового вопроса в форму создания опроса.
 */
function addQuestion() {
    const container = document.getElementById('questions-container');
    const questionCount = container.children.length + 1;

    const questionDiv = document.createElement('div');
    questionDiv.className = 'question';
    questionDiv.innerHTML = `
        <h3>Вопрос ${questionCount}</h3>
        <label>Текст вопроса:</label>
        <input type="text" name="question_text" required>
        
        <label>Тип вопроса:</label>
        <select name="question_type" onchange="toggleOptions(this)">
            <option value="single_choice">Один вариант ответа</option>
            <option value="multiple_choice">Несколько вариантов ответа</option>
            <option value="free_text">Свободный ответ</option>
        </select>
        
        <div class="options-container">
            <button type="button" onclick="addOption(this)">Добавить Вариант</button>
        </div>
        <button type="button" class="negative-btn" onclick="removeQuestion(this)">Удалить Вопрос</button>
    `;
    container.appendChild(questionDiv);
}

/**
 * Функция для удаления вопроса из формы создания опроса.
 * @param {HTMLElement} button - Кнопка удаления вопроса.
 */
function removeQuestion(button) {
    const questionDiv = button.parentElement;
    questionDiv.remove();
    renumberQuestions();
}

/**
 * Функция для добавления нового варианта ответа к вопросу.
 * @param {HTMLElement} button - Кнопка добавления варианта.
 */
function addOption(button) {
    const optionsContainer = button.parentElement;
    const optionDiv = document.createElement('div');
    optionDiv.className = 'answer-option';
    optionDiv.innerHTML = `
        <input type="text" name="option_text" placeholder="Текст варианта" required>
        <button type="button" class="negative-btn" onclick="removeOption(this)">Удалить</button>
    `;
    optionsContainer.insertBefore(optionDiv, button);
}

/**
 * Функция для удаления варианта ответа из вопроса.
 * @param {HTMLElement} button - Кнопка удаления варианта.
 */
function removeOption(button) {
    const optionDiv = button.parentElement;
    optionDiv.remove();
}

/**
 * Функция для отображения или скрытия контейнера опций ответов в зависимости от типа вопроса.
 * @param {HTMLElement} select - Выпадающий список выбора типа вопроса.
 */
function toggleOptions(select) {
    const questionDiv = select.parentElement.parentElement;
    const optionsContainer = questionDiv.querySelector('.options-container');
    if (select.value === 'free_text') {
        optionsContainer.style.display = 'none';
    } else {
        optionsContainer.style.display = 'block';
    }
}

/**
 * Функция для перенумерации вопросов после удаления.
 */
function renumberQuestions() {
    const container = document.getElementById('questions-container');
    const questionDivs = container.querySelectorAll('.question');
    questionDivs.forEach((qDiv, index) => {
        const header = qDiv.querySelector('h3');
        if (header) {
            header.textContent = `Вопрос ${index + 1}`;
        }
    });
}

// Обработка формы создания опроса
document.addEventListener('DOMContentLoaded', function() {
    const createForm = document.getElementById('create-survey-form');
    if (createForm) {
        createForm.addEventListener('submit', async function(e) {
            e.preventDefault();
            const title = document.getElementById('survey-title').value.trim();
            const description = document.getElementById('survey-description').value.trim();

            if (title === '') {
                showError('Пожалуйста, введите название опроса.');
                return;
            }

            const questions = [];
            const questionDivs = document.querySelectorAll('#questions-container .question');
            for (const qDiv of questionDivs) {
                const questionText = qDiv.querySelector('input[name="question_text"]').value.trim();
                const questionType = qDiv.querySelector('select[name="question_type"]').value;
                const options = [];

                if (questionType === 'single_choice' || questionType === 'multiple_choice') {
                    const optionInputs = qDiv.querySelectorAll('input[name="option_text"]');
                    for (const optInput of optionInputs) {
                        const optText = optInput.value.trim();
                        if (optText !== '') {
                            options.push(optText);
                        }
                    }
                    if (options.length === 0) {
                        showError(`Пожалуйста, добавьте хотя бы один вариант ответа для вопроса: "${questionText}".`);
                        return;
                    }
                }

                questions.push({
                    question_text: questionText,
                    question_type: questionType,
                    options: options
                });
            }

            if (questions.length === 0) {
                showError('Пожалуйста, добавьте хотя бы один вопрос в опрос.');
                return;
            }

            const surveyData = {
                title: title,
                description: description,
                questions: questions
            };

            try {
                const res = await httpRequest('POST', '/survey/create', surveyData);
                showSuccess(`Опрос успешно создан с ID ${res}`);
                window.location.href = 'index.html';
            } catch (err) {
                showError(`Ошибка при создании опроса: ${err}`);
            }
        });
    }
});
