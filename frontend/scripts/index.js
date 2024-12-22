// scripts/index.js

// Функция для проверки авторизации
async function checkAuth() {
    try {
        const res = await httpRequest('GET', '/check_auth', null);
        if (res.authenticated) {
            document.getElementById('logout-button').style.display = 'inline-block';
            document.getElementById('protected-content').style.display = 'block';
            document.getElementById('auth-forms').style.display = 'none';
        }
    } catch (err) {
        console.log('Пользователь не авторизован');
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
                const res = await httpRequest('POST', '/login', { login, password });
                showSuccess('Успешный вход');
                checkAuth();
            } catch (err) {
                showError('Ошибка при входе: ' + err);
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
                const res = await httpRequest('POST', '/register', { login, password });
                showSuccess('Успешная регистрация');
                checkAuth();
            } catch (err) {
                showError('Ошибка при регистрации: ' + err);
            }
        });
    }

    // Обработчик выхода
    const logoutButton = document.getElementById('logout-button');
    if (logoutButton) {
        logoutButton.addEventListener('click', async function() {
            try {
                const res = await httpRequest('POST', '/logout', null);
                showSuccess(res.message);
                // Перезагрузить страницу или обновить интерфейс
                window.location.href = 'index.html';
            } catch (err) {
                showError('Ошибка при выходе: ' + err);
            }
        });
    }

    // Проверка авторизации при загрузке страницы
    checkAuth();
});

// Навигационные функции

/**
 * Функция для перехода на страницу прохождения опроса.
 */
function takeSurvey() {
    const surveyID = prompt("Введите ID опроса:");
    if (surveyID === null || surveyID.trim() === "") {
        showError('Пожалуйста, введите ID опроса.');
        return;
    }
    window.location.href = `survey.html?id=${encodeURIComponent(surveyID.trim())}`;
}

/**
 * Функция для перехода на страницу редактирования опроса.
 */
function editSurvey() {
    const surveyID = prompt("Введите ID опроса для редактирования:");
    if (surveyID === null || surveyID.trim() === "") {
        showError('Пожалуйста, введите ID опроса.');
        return;
    }
    window.location.href = `edit_survey.html?id=${encodeURIComponent(surveyID.trim())}`;
}

/**
 * Функция для перехода на страницу аналитики опроса.
 */
function viewAnalysis() {
    const surveyID = prompt("Введите ID опроса для аналитики:");
    if (surveyID === null || surveyID.trim() === "") {
        showError('Пожалуйста, введите ID опроса.');
        return;
    }
    window.location.href = `analysis.html?id=${encodeURIComponent(surveyID.trim())}`;
}
