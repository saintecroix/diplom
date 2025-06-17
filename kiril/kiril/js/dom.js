// Элементы DOM
export const elements = {
  loginPage: document.getElementById('login-page'),
  app: document.getElementById('app'),
  loginForm: document.getElementById('login-form'),
  usernameInput: document.getElementById('username'),
  passwordInput: document.getElementById('password'),
  logoutBtn: document.getElementById('logout-btn'),
  usernameDisplay: document.getElementById('username-display'),
  pageContent: document.getElementById('page-content'),
  modalOverlay: document.getElementById('modal-overlay')
};

// Функция показа модального окна
export function showModal(title, content, confirmCallback = null) {
  const modalTitle = document.getElementById('modal-title');
  const modalBody = document.getElementById('modal-body');
  const modalConfirm = document.getElementById('modal-confirm');
  
  if (modalTitle) modalTitle.textContent = title;
  if (modalBody) modalBody.innerHTML = content;
  
  if (confirmCallback && modalConfirm) {
    modalConfirm.addEventListener('click', confirmCallback);
    modalConfirm.classList.remove('hidden');
  } else if (modalConfirm) {
    modalConfirm.classList.add('hidden');
  }
  
  elements.modalOverlay.classList.remove('hidden');
}

// Отображение основного приложения
export function showApp(username) {
  if (elements.usernameDisplay) {
    elements.usernameDisplay.textContent = username;
  }
  
  elements.loginPage.classList.add('hidden');
  elements.app.classList.remove('hidden');
}

// Отображение страницы входа
export function showLogin() {
  elements.app.classList.add('hidden');
  elements.loginPage.classList.remove('hidden');
}

// Функции для работы с ошибками формы
export function markInputAsError(inputElement, message) {
  if (!inputElement) return;
  
  inputElement.classList.add('error');
  const errorSpan = inputElement.nextElementSibling;
  if (errorSpan && errorSpan.classList.contains('error-message')) {
    errorSpan.textContent = message;
    errorSpan.style.display = 'block';
  }
}

export function clearInputErrors() {
  const errorInputs = document.querySelectorAll('.error');
  errorInputs.forEach(input => {
    input.classList.remove('error');
    const errorSpan = input.nextElementSibling;
    if (errorSpan && errorSpan.classList.contains('error-message')) {
      errorSpan.textContent = '';
      errorSpan.style.display = 'none';
    }
  });
}

// Функции для работы с загрузчиком
export function showLoader() {
  const loader = document.getElementById('loader');
  if (loader) loader.classList.remove('hidden');
}

export function hideLoader() {
  const loader = document.getElementById('loader');
  if (loader) loader.classList.add('hidden');
}


export function hideModal() {
  const modalOverlay = document.getElementById('modal-overlay');
  if (modalOverlay) {
    modalOverlay.classList.add('hidden');
  }
}

export function showNotification(message, type = 'info') {
  // Создаем элемент уведомления
  const notification = document.createElement('div');
  notification.className = `notification ${type}`;
  notification.innerText = message;
  
  // Добавляем на страницу
  document.body.appendChild(notification);
  
  // Анимация появления
  setTimeout(() => {
    notification.classList.add('show');
  }, 10);
  
  // Автоматическое скрытие через 3 секунды
  setTimeout(() => {
    notification.classList.remove('show');
    setTimeout(() => {
      notification.remove();
    }, 300);
  }, 3000);
}