import { elements, showModal, showApp, showLogin, markInputAsError, clearInputErrors } from './dom.js';
import { setToken, clearToken, isAuthenticated } from './auth.js';
import { loadPage } from './router.js';

export function setupEvents() {
  // Обработчик формы входа
  elements.loginForm.addEventListener('submit', (e) => {
    e.preventDefault();
    handleLogin();
  });

  // Обработчик кнопки выхода
  elements.logoutBtn.addEventListener('click', handleLogout);

  // Обработчик тестовых данных
  const testCredentialsBtn = document.getElementById('test-credentials-btn');
  if (testCredentialsBtn) {
    testCredentialsBtn.addEventListener('click', useTestCredentials);
  }
}

function handleLogin() {
  const username = document.getElementById('username').value.trim();
  const password = document.getElementById('password').value.trim();

  clearInputErrors();

  if (!username || !password) {
    if (!username) markInputAsError(document.getElementById('username'), 'Введите логин');
    if (!password) markInputAsError(document.getElementById('password'), 'Введите пароль');
    return;
  }

  // Проверка тестовых данных
  if (username === 'admin' && password === 'admin') {
    setToken('demo-token');
    showApp('Администратор');
    loadPage('home');
  } else {
    markInputAsError(document.getElementById('username'), 'Неверный логин или пароль');
    markInputAsError(document.getElementById('password'), 'Неверный логин или пароль');
  }
}

function handleLogout() {
  clearToken();
  showLogin();
}

function useTestCredentials() {
  document.getElementById('username').value = 'admin';
  document.getElementById('password').value = 'admin';
  clearInputErrors();
}