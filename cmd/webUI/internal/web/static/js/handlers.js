import { showModal, hideModal, showNotification, markInputAsError, clearInputErrors } from './dom.js';
import { getUserEditForm } from './templates.js';
import { loadPage } from './router.js';



// Другие функции...

export function initPageHandlers(page) {
  switch (page) {
    case 'home':
      // Обработчики для главной страницы
      break;
    case 'wagon-search':
      initWagonHandlers();
      break;
    case 'data-input':
      initDataInputHandlers();
      break;
    case 'admin':
      initAdminHandlers();
      break;
    case 'transportations':
      initTransportationsHandlers();
      break;
  }
}

// Полная реализация initAdminHandlers
function initAdminHandlers() {
  console.log('Инициализация обработчиков для страницы администрирования');
  
  // Кнопка добавления пользователя
  const addUserBtn = document.getElementById('add-user-btn');
  if (addUserBtn) {
    addUserBtn.addEventListener('click', () => {
      showAddUserForm();
    });
  }

  // Кнопки редактирования пользователей
  const editUserBtns = document.querySelectorAll('.edit-user-btn');
  editUserBtns.forEach(btn => {
    btn.addEventListener('click', function() {
      const userId = this.getAttribute('data-id');
      console.log('Редактирование пользователя с ID:', userId);
      showEditUserForm(userId);
    });
  });

  // Кнопки удаления пользователей
  const deleteUserBtns = document.querySelectorAll('.delete-user-btn');
  deleteUserBtns.forEach(btn => {
    btn.addEventListener('click', function() {
      const userId = this.getAttribute('data-id');
      showDeleteUserConfirmation(userId);
    });
  });

  // Установка обработчиков для динамически созданных форм
  document.addEventListener('submit', function(e) {
    if (e.target.id === 'edit-user-form') {
      e.preventDefault();
      handleEditUser();
    } else if (e.target.id === 'add-user-form') {
      e.preventDefault();
      handleAddUser();
    }
  });

  // Обработчики для кнопок отмены в модальных окнах
  document.addEventListener('click', function(e) {
    if (e.target.id === 'cancel-edit-user' || 
        e.target.id === 'cancel-add-user' || 
        e.target.id === 'cancel-delete-user') {
      hideModal();
    }
  });
}

// Функция показа формы добавления пользователя
function showAddUserForm() {
  const user = { role: 'viewer' }; // Значения по умолчанию
  showModal('Добавление пользователя', getUserEditForm(user));
}

// Функция показа формы редактирования пользователя
function showEditUserForm(userId) {
  console.log('Получение данных пользователя с ID:', userId);
  
  // В реальном приложении здесь был бы API-запрос
  // Для демонстрации используем моковые данные
  const user = {
    id: userId,
    username: 'user' + userId,
    fullName: 'Пользователь ' + userId,
    role: userId % 2 === 0 ? 'admin' : 'operator'
  };
  
  showModal('Редактирование пользователя', getUserEditForm(user));
}

// Функция обработки формы редактирования пользователя
function handleEditUser() {
  clearInputErrors();
  
  const userId = document.getElementById('edit-user-id').value;
  const username = document.getElementById('edit-username').value;
  const fullName = document.getElementById('edit-full-name').value;
  const role = document.getElementById('edit-role').value;
  const password = document.getElementById('edit-password').value;
  const confirmPassword = document.getElementById('edit-confirm-password').value;
  
  console.log('Обработка редактирования пользователя:', {
    userId, username, fullName, role, 
    password: password ? '****' : '',
    confirmPassword: confirmPassword ? '****' : ''
  });
  
  // Валидация
  let hasErrors = false;
  
  if (!username) {
    markInputAsError(document.getElementById('edit-username'), 'Введите имя пользователя');
    hasErrors = true;
  }
  
  if (!fullName) {
    markInputAsError(document.getElementById('edit-full-name'), 'Введите ФИО');
    hasErrors = true;
  }
  
  if (password && password !== confirmPassword) {
    markInputAsError(document.getElementById('edit-confirm-password'), 'Пароли не совпадают');
    hasErrors = true;
  }
  
  if (hasErrors) return;
  
  // В реальном приложении здесь был бы API-запрос
  // Для демонстрации просто закрываем модальное окно
  
  hideModal();
  showNotification('Пользователь успешно обновлен', 'success');
  
  // Перезагрузка страницы для обновления списка пользователей
  setTimeout(() => {
    loadPage('admin');
  }, 1000);
}

// Функция для показа подтверждения удаления пользователя
function showDeleteUserConfirmation(userId) {
  const content = `
    <p>Вы действительно хотите удалить пользователя?</p>
    <div class="form-actions">
      <button id="confirm-delete-user" class="btn btn-danger" data-id="${userId}">Удалить</button>
      <button id="cancel-delete-user" class="btn btn-secondary">Отмена</button>
    </div>
  `;
  
  showModal('Удаление пользователя', content);
  
  // Добавляем обработчик для кнопки подтверждения
  const confirmBtn = document.getElementById('confirm-delete-user');
  if (confirmBtn) {
    confirmBtn.addEventListener('click', function() {
      const userId = this.getAttribute('data-id');
      handleDeleteUser(userId);
    });
  }
}

// Функция обработки удаления пользователя
function handleDeleteUser(userId) {
  console.log('Удаление пользователя с ID:', userId);
  
  // В реальном приложении здесь был бы API-запрос
  
  hideModal();
  showNotification('Пользователь успешно удален', 'success');
  
  // Перезагрузка страницы для обновления списка пользователей
  setTimeout(() => {
    loadPage('admin');
  }, 1000);
}

// Функция обработки добавления пользователя
function handleAddUser() {
  clearInputErrors();
  
  const username = document.getElementById('edit-username').value;
  const fullName = document.getElementById('edit-full-name').value;
  const role = document.getElementById('edit-role').value;
  const password = document.getElementById('edit-password').value;
  const confirmPassword = document.getElementById('edit-confirm-password').value;
  
  // Валидация
  let hasErrors = false;
  
  if (!username) {
    markInputAsError(document.getElementById('edit-username'), 'Введите имя пользователя');
    hasErrors = true;
  }
  
  if (!fullName) {
    markInputAsError(document.getElementById('edit-full-name'), 'Введите ФИО');
    hasErrors = true;
  }
  
  if (!password) {
    markInputAsError(document.getElementById('edit-password'), 'Введите пароль');
    hasErrors = true;
  } else if (password !== confirmPassword) {
    markInputAsError(document.getElementById('edit-confirm-password'), 'Пароли не совпадают');
    hasErrors = true;
  }
  
  if (hasErrors) return;
  
  // В реальном приложении здесь был бы API-запрос
  
  hideModal();
  showNotification('Пользователь успешно добавлен', 'success');
  
  // Перезагрузка страницы для обновления списка пользователей
  setTimeout(() => {
    loadPage('admin');
  }, 1000);
}