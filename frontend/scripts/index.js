// scripts/index.js

// Функция для проверки аутентификации администратора
async function checkAuth() {
    try {
        const res = await httpRequest('GET', '/api/check_auth', null);
        if (res.authenticated) {
            document.getElementById('logout-button').style.display = 'inline-block';
            document.getElementById('protected-content').style.display = 'block';
            document.getElementById('auth-forms').style.display = 'none';
        }
    } catch (err) {
        console.log('Администратор не авторизован');
        document.getElementById('logout-button').style.display = 'none';
        document.getElementById('protected-content').style.display = 'none';
        document.getElementById('auth-forms').style.display = 'block';
    }
}

// Обработка формы входа
document.addEventListener('DOMContentLoaded', function() {
    const loginForm = document.getElementById('login-form');
    if (loginForm) {
        loginForm.addEventListener('submit', async function(e) {
            e.preventDefault();
            const login = document.getElementById('login').value.trim();
            const password = document.getElementById('password').value.trim();

            if (login === '' || password === '') {
                showError('Пожалуйста, заполните все поля для входа.');
                return;
            }

            try {
                const res = await httpRequest('POST', '/api/login', { login, password });
                showSuccess('Успешный вход');
                checkAuth();
            } catch (err) {
                showError(`Ошибка при входе: ${err.message}`);
            }
        });
    }

    // Обработка формы регистрации
    const registerForm = document.getElementById('register-form');
    if (registerForm) {
        registerForm.addEventListener('submit', async function(e) {
            e.preventDefault();
            const login = document.getElementById('reg-login').value.trim();
            const password = document.getElementById('reg-password').value.trim();

            if (login === '' || password === '') {
                showError('Пожалуйста, заполните все поля для регистрации.');
                return;
            }

            try {
                const res = await httpRequest('POST', '/api/register', { login, password });
                showSuccess('Успешная регистрация');
                checkAuth();
            } catch (err) {
                showError(`Ошибка при регистрации: ${err.message}`);
            }
        });
    }

    // Обработчик выхода
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

    // Проверка аутентификации при загрузке страницы
    checkAuth();
});

// Навигационные функции

/**
 * Функция для перехода на страницу создания опроса.
 */
function createSurvey() {
    window.location.href = 'create_survey.html';
}

/**
 * Функция для перехода на страницу просмотра опросов.
 */
function viewSurveys() {
    window.location.href = 'view_surveys.html';
}

/**
 * Функция для перехода на страницу аналитики опроса.
 */
function viewAnalytics() {
    const surveyID = prompt("Введите ID опроса для аналитики:");
    if (surveyID === null || surveyID.trim() === "") {
        showError('Пожалуйста, введите ID опроса.');
        return;
    }
    window.location.href = `analytics.html?id=${encodeURIComponent(surveyID.trim())}`;
}
