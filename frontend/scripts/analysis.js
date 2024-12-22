// scripts/analysis.js

/**
 * Функция для обработки формы запроса аналитики.
 */
document.addEventListener('DOMContentLoaded', function() {
    const analysisForm = document.getElementById('analysis-form');
    if (analysisForm) {
        analysisForm.addEventListener('submit', async function(e) {
            e.preventDefault();
            const surveyID = document.getElementById('analysis-survey-id-input').value.trim();
            const type = document.getElementById('analysis-type').value;

            if (surveyID === '') {
                showError('Пожалуйста, введите ID опроса.');
                return;
            }

            try {
                const res = await httpRequest('GET', `/analytics/${encodeURIComponent(surveyID)}?type=${encodeURIComponent(type)}`);
                displayAnalysis(res, type);
            } catch (err) {
                showError(`Ошибка при получении аналитики: ${err}`);
            }
        });
    }
});

/**
 * Функция для отображения результатов аналитики.
 * @param {object} data - Данные аналитики.
 * @param {string} type - Тип аналитики.
 */
function displayAnalysis(data, type) {
    const resultDiv = document.getElementById('analysis-result');
    resultDiv.innerHTML = ''; // Очистка предыдущих результатов

    if (type === 'heatmap') {
        const table = document.createElement('table');
        table.innerHTML = `
            <tr>
                <th>Вопрос</th>
                <th>Количество Ответов</th>
            </tr>
        `;
        for (const [question, count] of Object.entries(data)) {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${question}</td>
                <td>${count}</td>
            `;
            table.appendChild(row);
        }
        resultDiv.appendChild(table);
    } else if (type === 'response_count') {
        const p = document.createElement('p');
        p.textContent = `Количество ответов: ${data.response_count}`;
        resultDiv.appendChild(p);
    } else if (type === 'average_score') {
        const p = document.createElement('p');
        p.textContent = `Средний балл: ${data.average_score}`;
        resultDiv.appendChild(p);
    } else {
        resultDiv.textContent = 'Неизвестный тип анализа.';
    }
}
