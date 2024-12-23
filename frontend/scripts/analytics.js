// scripts/analytics.js

/**
 * Загружает данные аналитики и отображает графики.
 */
async function loadAnalytics() {
    const urlParams = new URLSearchParams(window.location.search);
    const surveyID = urlParams.get('id');

    if (!surveyID) {
        showError('Не указан ID опроса.');
        return;
    }

    try {
        // Получение деталей опроса для отображения заголовка
        const survey = await httpRequest('GET', `/api/surveys/${surveyID}`, null);
        document.getElementById('survey-title').innerText = `Аналитика Опроса: ${survey.Title}`;

        // Получение данных для тепловой карты
        const heatmapData = await httpRequest('GET', `/api/surveys/${surveyID}/analytics?type=heatmap`, null);
        renderHeatmap(heatmapData);

        // Получение данных для пай-чарта
        const pieChartData = await httpRequest('GET', `/api/surveys/${surveyID}/analytics?type=pie_chart`, null);
        renderPieChart(pieChartData);
    } catch (err) {
        showError(`Ошибка при загрузке аналитики: ${err.message}`);
    }
}

/**
 * Рендерит тепловую карту с использованием Chart.js.
 * @param {Object} data - Данные тепловой карты.
 */
function renderHeatmap(data) {
    const ctx = document.getElementById('heatmap-chart').getContext('2d');
    const labels = Object.keys(data);
    const counts = Object.values(data);

    new Chart(ctx, {
        type: 'bar',
        data: {
            labels: labels,
            datasets: [{
                label: 'Количество Ответов',
                data: counts,
                backgroundColor: 'rgba(54, 162, 235, 0.6)',
                borderColor: 'rgba(54, 162, 235, 1)',
                borderWidth: 1
            }]
        },
        options: {
            plugins: {
                title: {
                    display: true,
                    text: 'Тепловая Карта Ответов по Вопросам'
                }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    precision: 0
                }
            }
        }
    });
}

/**
 * Рендерит пай-чарты с использованием Chart.js.
 * @param {Array} data - Массив объектов с вопросами и распределением ответов.
 */
function renderPieChart(data) {
    const ctx = document.getElementById('pie-chart').getContext('2d');

    data.forEach((item, index) => {
        const labels = Object.keys(item.Options);
        const counts = Object.values(item.Options);

        const chartCanvas = document.createElement('canvas');
        chartCanvas.id = `pie-chart-${index}`;
        document.getElementById('charts').appendChild(chartCanvas);

        new Chart(chartCanvas, {
            type: 'pie',
            data: {
                labels: labels,
                datasets: [{
                    label: item.Question,
                    data: counts,
                    backgroundColor: generateColors(labels.length),
                    borderWidth: 1
                }]
            },
            options: {
                plugins: {
                    title: {
                        display: true,
                        text: `Распределение Ответов для Вопроса: ${item.Question}`
                    }
                }
            }
        });
    });
}

/**
 * Генерирует массив случайных цветов.
 * @param {number} count - Количество цветов.
 * @returns {Array} - Массив строк с цветами в формате rgba.
 */
function generateColors(count) {
    const colors = [];
    for (let i = 0; i < count; i++) {
        const r = Math.floor(Math.random() * 255);
        const g = Math.floor(Math.random() * 255);
        const b = Math.floor(Math.random() * 255);
        colors.push(`rgba(${r}, ${g}, ${b}, 0.6)`);
    }
    return colors;
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

    // Загрузка аналитики при загрузке страницы
    loadAnalytics();
});
