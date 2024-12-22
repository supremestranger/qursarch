// scripts/edit_survey.js

/**
 * Функция для загрузки существующего опроса по его ID.
 */
async function loadSurvey() {
    const surveyID = new URLSearchParams(window.location.search).get('id') || document.getElementById('edit-survey-id-input').value.trim();
    if (surveyID === '') {
        showError('Пожалуйста, введите ID опроса.');
        return;
    }

    try {
        const survey = await httpRequest('GET', `/survey/${encodeURIComponent(surveyID)}`);
        document.getElementById('edit-survey-details').style.display = 'block';
        document.getElementById('edit-survey-title').value = survey.title;
        document.getElementById('edit-survey-description').value = survey.description;

        const container = document.getElementById('edit-questions-container');
        container.innerHTML = ''; // Очистка предыдущих вопросов

        survey.questions.forEach(function(q, index) {
            const questionDiv = document.createElement('div');
            questionDiv.className = 'question';
            questionDiv.innerHTML = `
                <h3>Вопрос ${index + 1}</h3>
                <label>Текст вопроса:</label>
                <input type="text" name="question_text" value="${q.question_text}" required>
                
                <label>Тип вопроса:</label>
                <select name="question_type" onchange="toggleOptions(this)">
                    <option value="single_choice" ${q.question_type === 'single_choice' ? 'selected' : ''}>Один вариант ответа</option>
                    <option value="multiple_choice" ${q.question_type === 'multiple_choice' ? 'selected' : ''}>Несколько вариантов ответа</option>
                    <option value="free_text" ${q.question_type === 'free_text' ? 'selected' : ''}>Свободный ответ</option>
                </select>
                
                <div class="options-container" style="${q.question_type === 'free_text' ? 'display:none;' : ''}">
                    ${q.question_type !== 'free_text' ? q.options.map(opt => `
                        <div class="answer-option">
                            <input type="text" name="option_text" value="${opt}" required>
                            <button type="button" onclick="removeOption(this)">Удалить</button>
                        </div>
                    `).join('') : ''}
                    <button type="button" onclick="addOption(this)">Добавить Вариант</button>
                </div>
                <button type="button" onclick="removeQuestion(this)">Удалить Вопрос</button>
            `;
            container.appendChild(questionDiv);
        });

        renumberQuestions();
    } catch (err) {
        showError(`Ошибка при загрузке опроса: ${err}`);
    }
}

/**
 * Функция для добавления нового вопроса в форму редактирования опроса.
 */
function addEditQuestion() {
    const container = document.getElementById('edit-questions-container');
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
        <button type="button" onclick="removeQuestion(this)">Удалить Вопрос</button>
    `;
    container.appendChild(questionDiv);
}

/**
 * Функция для удаления вопроса из формы редактирования опроса.
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
        <button type="button" onclick="removeOption(this)">Удалить</button>
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
    const container = document.getElementById('edit-questions-container');
    const questionDivs = container.querySelectorAll('.question');
    questionDivs.forEach((qDiv, index) => {
        const header = qDiv.querySelector('h3');
        if (header) {
            header.textContent = `Вопрос ${index + 1}`;
        }
    });
}

/**
 * Функция для обработки и отправки изменений в опросе.
 */
document.addEventListener('DOMContentLoaded', function() {
    const editForm = document.getElementById('edit-survey-form');
    if (editForm) {
        editForm.addEventListener('submit', async function(e) {
            e.preventDefault();
            const surveyID = new URLSearchParams(window.location.search).get('id') || document.getElementById('edit-survey-id-input').value.trim();
            const title = document.getElementById('edit-survey-title').value.trim();
            const description = document.getElementById('edit-survey-description').value.trim();

            if (surveyID === '') {
                showError('Пожалуйста, введите ID опроса.');
                return;
            }

            if (title === '') {
                showError('Пожалуйста, введите название опроса.');
                return;
            }

            const questions = [];
            const questionDivs = document.querySelectorAll('#edit-questions-container .question');
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
                const res = await httpRequest('PUT', `/survey/edit/${encodeURIComponent(surveyID)}`, surveyData);
                showSuccess('Опрос успешно обновлён');
                window.location.href = 'index.html';
            } catch (err) {
                showError(`Ошибка при редактировании опроса: ${err}`);
            }
        });
    }

    // Проверка наличия параметра ID в URL при загрузке страницы
    const urlParams = new URLSearchParams(window.location.search);
    const surveyID = urlParams.get('id');
    if (surveyID) {
        document.getElementById('edit-survey-id-input').value = surveyID;
        loadSurvey();
    }
});
