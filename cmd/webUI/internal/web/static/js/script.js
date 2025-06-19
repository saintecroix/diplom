// Глобальные переменные
let currentUser = null;
let jwtToken = localStorage.getItem('jwtToken');

// DOM элементы
const loginPage = document.getElementById('login-page');
const app = document.getElementById('app');
const loginForm = document.getElementById('login-form');
const usernameInput = document.getElementById('username');
const passwordInput = document.getElementById('password');
const logoutBtn = document.getElementById('logout-btn');
const usernameDisplay = document.getElementById('username-display');
const pageContent = document.getElementById('page-content');
const modalOverlay = document.getElementById('modal-overlay');
const modalTitle = document.getElementById('modal-title');
const modalBody = document.getElementById('modal-body');
const modalClose = document.getElementById('modal-close');
const modalCancel = document.getElementById('modal-cancel');
const modalConfirm = document.getElementById('modal-confirm');
const loader = document.getElementById('loader');

// Инициализация приложения
document.addEventListener('DOMContentLoaded', () => {
    // Проверка аутентификации
    checkAuth();
    
    // Настройка навигации
    setupNavigation();
    
    // Настройка модального окна
    setupModal();
    
    // Настройка обработчиков событий
    setupEventListeners();
});

// Проверка аутентификации
function checkAuth() {
    if (jwtToken) {
        // Здесь должна быть проверка JWT на бэкенде
        // Для демонстрации просто скрываем страницу входа
        loginPage.classList.add('hidden');
        app.classList.remove('hidden');
        currentUser = { username: 'Администратор' };
        usernameDisplay.textContent = currentUser.username;
        loadPage('home');
    } else {
        loginPage.classList.remove('hidden');
        app.classList.add('hidden');
    }
}

// Настройка навигации
function setupNavigation() {
    const navLinks = document.querySelectorAll('.main-nav a');
    
    navLinks.forEach(link => {
        link.addEventListener('click', (e) => {
            e.preventDefault();
            const page = link.getAttribute('data-page');
            
            // Удаляем активный класс у всех ссылок
            navLinks.forEach(l => l.classList.remove('active'));
            // Добавляем активный класс текущей ссылке
            link.classList.add('active');
            
            loadPage(page);
        });
    });
}

// Загрузка страницы
function loadPage(page) {
    showLoader();
    
    // Имитация загрузки с сервера
    setTimeout(() => {
        let content = '';
        
        switch (page) {
            case 'home':
                content = getHomePage();
                break;
            case 'wagon-search':
                content = getWagonSearchPage();
                break;
            case 'data-input':
                content = getDataInputPage();
                break;
            case 'admin':
                content = getAdminPage();
                break;
            case 'transportations':
                content = getTransportationsPage();
                break;
            default:
                content = getHomePage();
        }
        
        pageContent.innerHTML = content;
        hideLoader();
        
        // Инициализация специфичных для страницы обработчиков событий
        initPageSpecificHandlers(page);
    }, 500);
}

// Страница "Главная"
function getHomePage() {
    return `
        <div class="page-home">
            <h1 class="page-title">Главная</h1>
            <div class="card">
                <h2 class="card-title">Последние операции</h2>
                <div class="table-container">
                    <table>
                        <thead>
                            <tr>
                                <th>ID</th>
                                <th>Дата</th>
                                <th>Операция</th>
                                <th>Статус</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr>
                                <td>1</td>
                                <td>2023-06-15</td>
                                <td>Перемещение вагонов</td>
                                <td><span class="badge success">Завершено</span></td>
                            </tr>
                            <tr>
                                <td>2</td>
                                <td>2023-06-14</td>
                                <td>Создание группы</td>
                                <td><span class="badge success">Завершено</span></td>
                            </tr>
                            <tr>
                                <td>3</td>
                                <td>2023-06-13</td>
                                <td>Импорт данных</td>
                                <td><span class="badge warning">В процессе</span></td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    `;
}

// Страница "Поиск вагонов"
function getWagonSearchPage() {
    return `
        <div class="page-wagon-search">
            <h1 class="page-title">Поиск вагонов</h1>
            
            <div class="card">
                <h2 class="card-title">Создание группы</h2>
                <form id="create-group-form">
                    <div class="form-group">
                        <label for="group-name">Наименование группы</label>
                        <input type="text" id="group-name" name="group-name" required>
                    </div>
                    <div class="form-group">
                        <label for="wagon-numbers">Номера вагонов (через запятую)</label>
                        <textarea id="wagon-numbers" name="wagon-numbers" rows="4" required></textarea>
                    </div>
                    <button type="submit" class="btn-primary">Создать группу</button>
                </form>
            </div>
            
            <div class="card">
                <h2 class="card-title">Существующие группы</h2>
                <div class="table-container">
                    <table id="groups-table">
                        <thead>
                            <tr>
                                <th><input type="checkbox" id="select-all-groups"></th>
                                <th>Наименование</th>
                                <th>Количество вагонов</th>
                                <th>Действия</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr>
                                <td><input type="checkbox" class="group-checkbox"></td>
                                <td>Группа 1</td>
                                <td>15</td>
                                <td>
                                    <button class="btn-secondary edit-group-btn">Редактировать</button>
                                </td>
                            </tr>
                            <tr>
                                <td><input type="checkbox" class="group-checkbox"></td>
                                <td>Группа 2</td>
                                <td>8</td>
                                <td>
                                    <button class="btn-secondary edit-group-btn">Редактировать</button>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
                <div class="actions" style="margin-top: 20px;">
                    <button id="delete-groups-btn" class="btn-danger">Удалить выбранные</button>
                </div>
            </div>
        </div>
    `;
}

// Страница "Ввод данных"
function getDataInputPage() {
    return `
        <div class="page-data-input">
            <h1 class="page-title">Ввод данных</h1>
            
            <div class="card">
                <h2 class="card-title">Загрузка файла</h2>
                <form id="upload-form">
                    <div class="form-group">
                        <label for="data-file">Выберите файл для загрузки</label>
                        <input type="file" id="data-file" name="data-file" required>
                    </div>
                    <div class="progress-bar hidden" id="progress-container">
                        <div class="progress" id="progress-bar"></div>
                    </div>
                    <button type="submit" class="btn-primary">Загрузить</button>
                </form>
            </div>
        </div>
    `;
}

// Страница "Администрирование"
function getAdminPage() {
    return `
        <div class="page-admin">
            <h1 class="page-title">Администрирование</h1>
            
            <div class="card">
                <h2 class="card-title">Управление пользователями</h2>
                
                <div class="actions" style="margin-bottom: 20px;">
                    <button id="create-user-btn" class="btn-primary">Создать пользователя</button>
                    <button id="edit-user-btn" class="btn-secondary">Редактировать пользователя</button>
                </div>
                
                <div class="table-container">
                    <table id="users-table">
                        <thead>
                            <tr>
                                <th>Логин</th>
                                <th>Роль</th>
                                <th>Дата создания</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr>
                                <td>admin</td>
                                <td>Администратор</td>
                                <td>2023-01-01</td>
                            </tr>
                            <tr>
                                <td>user1</td>
                                <td>Оператор</td>
                                <td>2023-05-15</td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    `;
}

// Страница "Перевозки"
function getTransportationsPage() {
    return `
        <div class="page-transportations">
            <h1 class="page-title">Перевозки</h1>
            
            <div class="card">
                <div class="search-container" style="margin-bottom: 20px;">
                    <input type="text" id="transportation-search" placeholder="Поиск по перевозкам..." class="form-group input">
                </div>
                
                <div class="table-container">
                    <table id="transportations-table">
                        <thead>
                            <tr>
                                <th>Номер</th>
                                <th>Маршрут</th>
                                <th>Дата отправления</th>
                                <th>Дата прибытия</th>
                                <th>Статус</th>
                                <th>Действия</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr>
                                <td>TR-2023-001</td>
                                <td>Москва - Санкт-Петербург</td>
                                <td>2023-06-10</td>
                                <td>2023-06-12</td>
                                <td><span class="badge success">Завершено</span></td>
                                <td>
                                    <button class="btn-secondary">Детали</button>
                                </td>
                            </tr>
                            <tr>
                                <td>TR-2023-002</td>
                                <td>Новосибирск - Владивосток</td>
                                <td>2023-06-15</td>
                                <td>2023-06-20</td>
                                <td><span class="badge warning">В пути</span></td>
                                <td>
                                    <button class="btn-secondary">Детали</button>
                                </td>
                            </tr>
                            <tr>
                                <td>TR-2023-003</td>
                                <td>Казань - Сочи</td>
                                <td>2023-06-18</td>
                                <td>2023-06-22</td>
                                <td><span class="badge info">Запланировано</span></td>
                                <td>
                                    <button class="btn-secondary">Детали</button>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    `;
}

// Настройка обработчиков событий
function setupEventListeners() {
    // Форма входа
    loginForm.addEventListener('submit', (e) => {
        e.preventDefault();
        handleLogin();
    });
    
    // Кнопка выхода
    logoutBtn.addEventListener('click', handleLogout);
    
    // Кнопка тестовых данных
    const testCredentialsBtn = document.getElementById('test-credentials-btn');
    if (testCredentialsBtn) {
        testCredentialsBtn.addEventListener('click', useTestCredentials);
    }
}

function useTestCredentials() {
    usernameInput.value = 'admin';
    passwordInput.value = 'admin';
    
    // Убираем подсветку ошибок, если они были
    clearInputError(usernameInput);
    clearInputError(passwordInput);
    
    // Можно автоматически отправить форму
    // handleLogin();
    // Или просто заполнить поля и ждать, когда пользователь сам нажмет "Войти"
}

// Обработка входа
function handleLogin() {
    const username = usernameInput.value.trim();
    const password = passwordInput.value.trim();
    
    // Валидация
    let isValid = true;
    
    if (!username) {
        markInputAsError(usernameInput, 'Введите логин');
        isValid = false;
    } else {
        clearInputError(usernameInput);
    }
    
    if (!password) {
        markInputAsError(passwordInput, 'Введите пароль');
        isValid = false;
    } else {
        clearInputError(passwordInput);
    }
    
    if (!isValid) return;
    
    showLoader();
    
    // Имитация запроса к серверу
    setTimeout(() => {
        // В реальном приложении здесь будет запрос к серверу
        if (username === 'admin' && password === 'admin') {
            // Имитация получения JWT токена
            jwtToken = 'demo-jwt-token';
            localStorage.setItem('jwtToken', jwtToken);
            
            checkAuth();
        } else {
            markInputAsError(usernameInput, 'Неправильный логин или пароль');
            markInputAsError(passwordInput, 'Неправильный логин или пароль');
        }
        
        hideLoader();
    }, 1000);
}

// Обработка выхода
function handleLogout() {
    jwtToken = null;
    localStorage.removeItem('jwtToken');
    currentUser = null;
    
    // Очистка формы входа
    loginForm.reset();
    clearInputError(usernameInput);
    clearInputError(passwordInput);
    
    checkAuth();
}

// Пометить поле ввода как ошибочное
function markInputAsError(inputElement, message) {
    const formGroup = inputElement.closest('.form-group');
    formGroup.classList.add('has-error');
    const errorElement = formGroup.querySelector('.error-message');
    errorElement.textContent = message;
}

// Очистить ошибку поля ввода
function clearInputError(inputElement) {
    const formGroup = inputElement.closest('.form-group');
    formGroup.classList.remove('has-error');
    const errorElement = formGroup.querySelector('.error-message');
    errorElement.textContent = '';
}

// Настройка модального окна
function setupModal() {
    modalClose.addEventListener('click', hideModal);
    modalCancel.addEventListener('click', hideModal);
    modalOverlay.addEventListener('click', (e) => {
        if (e.target === modalOverlay) {
            hideModal();
        }
    });
}

// Показать модальное окно
function showModal(title, content, confirmCallback = null) {
    modalTitle.textContent = title;
    modalBody.innerHTML = content;
    modalOverlay.classList.remove('hidden');
    
    if (confirmCallback) {
        modalConfirm.onclick = () => {
            confirmCallback();
            hideModal();
        };
        modalConfirm.classList.remove('hidden');
    } else {
        modalConfirm.classList.add('hidden');
    }
}

// Скрыть модальное окно
function hideModal() {
    modalOverlay.classList.add('hidden');
}

// Показать загрузчик
function showLoader() {
    loader.classList.remove('hidden');
}

// Скрыть загрузчик
function hideLoader() {
    loader.classList.add('hidden');
}

// Инициализация обработчиков событий для конкретной страницы
function initPageSpecificHandlers(page) {
    switch (page) {
        case 'wagon-search':
            initWagonSearchHandlers();
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

// Обработчики для страницы "Поиск вагонов"
function initWagonSearchHandlers() {
    // Создание группы
    const createGroupForm = document.getElementById('create-group-form');
    if (createGroupForm) {
        createGroupForm.addEventListener('submit', (e) => {
            e.preventDefault();
            
            const groupName = document.getElementById('group-name').value.trim();
            const wagonNumbers = document.getElementById('wagon-numbers').value.trim();
            
            if (!groupName || !wagonNumbers) {
                showModal('Ошибка', 'Заполните все поля формы');
                return;
            }
            
            // Формируем JSON для отправки на сервер
            const groupData = {
                name: groupName,
                wagons: wagonNumbers.split(',').map(num => num.trim()).filter(num => num)
            };
            
            console.log('Отправка данных группы на сервер:', groupData);
            
            showModal('Успех', `Группа "${groupName}" успешно создана с ${groupData.wagons.length} вагонами`);
            createGroupForm.reset();
        });
    }
    
    // Удаление групп
    const deleteGroupsBtn = document.getElementById('delete-groups-btn');
    if (deleteGroupsBtn) {
        deleteGroupsBtn.addEventListener('click', () => {
            const selectedGroups = document.querySelectorAll('.group-checkbox:checked');
            
            if (selectedGroups.length === 0) {
                showModal('Ошибка', 'Выберите хотя бы одну группу для удаления');
                return;
            }
            
            showModal(
                'Подтверждение удаления',
                `Вы уверены, что хотите удалить ${selectedGroups.length} выбранных групп?`,
                () => {
                    // Здесь будет логика удаления групп
                    console.log('Удаление выбранных групп');
                    showModal('Успех', 'Выбранные группы успешно удалены');
                }
            );
        });
    }
    
    // Редактирование групп
    const editGroupBtns = document.querySelectorAll('.edit-group-btn');
    editGroupBtns.forEach(btn => {
        btn.addEventListener('click', () => {
            const groupRow = btn.closest('tr');
            const groupName = groupRow.querySelector('td:nth-child(2)').textContent;
            
            showModal(
                `Редактирование группы "${groupName}"`,
                `
                <div class="form-group">
                    <label for="edit-group-name">Наименование группы</label>
                    <input type="text" id="edit-group-name" value="${groupName}" class="form-control">
                </div>
                <div class="form-group">
                    <label>Номера вагонов</label>
                    <div class="wagon-list" style="margin-top: 10px;">
                        <div class="wagon-item" style="display: flex; align-items: center; margin-bottom: 5px;">
                            <input type="text" value="12345678" class="form-control" style="flex-grow: 1; margin-right: 10px;">
                            <button class="btn-danger" style="padding: 5px 10px;">×</button>
                        </div>
                        <div class="wagon-item" style="display: flex; align-items: center; margin-bottom: 5px;">
                            <input type="text" value="87654321" class="form-control" style="flex-grow: 1; margin-right: 10px;">
                            <button class="btn-danger" style="padding: 5px 10px;">×</button>
                        </div>
                    </div>
                    <button id="add-wagon-btn" class="btn-secondary" style="margin-top: 10px;">Добавить вагон</button>
                </div>
                `,
                () => {
                    const newName = document.getElementById('edit-group-name').value;
                    console.log('Сохранение изменений группы:', newName);
                    showModal('Успех', `Группа "${newName}" успешно обновлена`);
                }
            );
            
            // Обработчик для кнопки добавления вагона в модальном окне
            const addWagonBtn = document.getElementById('add-wagon-btn');
            if (addWagonBtn) {
                addWagonBtn.addEventListener('click', () => {
                    const wagonList = document.querySelector('.wagon-list');
                    const newWagonItem = document.createElement('div');
                    newWagonItem.className = 'wagon-item';
                    newWagonItem.style = 'display: flex; align-items: center; margin-bottom: 5px;';
                    newWagonItem.innerHTML = `
                        <input type="text" placeholder="Номер вагона" class="form-control" style="flex-grow: 1; margin-right: 10px;">
                        <button class="btn-danger" style="padding: 5px 10px;">×</button>
                    `;
                    wagonList.appendChild(newWagonItem);
                    
                    // Обработчик для кнопки удаления вагона
                    newWagonItem.querySelector('.btn-danger').addEventListener('click', () => {
                        wagonList.removeChild(newWagonItem);
                    });
                });
            }
            
            // Обработчики для кнопок удаления вагонов
            document.querySelectorAll('.wagon-item .btn-danger').forEach(btn => {
                btn.addEventListener('click', function() {
                    this.closest('.wagon-item').remove();
                });
            });
        });
    });
}

// Обработчики для страницы "Ввод данных"
function initDataInputHandlers() {
    const uploadForm = document.getElementById('upload-form');
    const fileInput = document.getElementById('data-file');
    const progressContainer = document.getElementById('progress-container');
    const progressBar = document.getElementById('progress-bar');
    
    if (uploadForm) {
        uploadForm.addEventListener('submit', (e) => {
            e.preventDefault();
            
            if (!fileInput.files.length) {
                showModal('Ошибка', 'Выберите файл для загрузки');
                return;
            }
            
            const file = fileInput.files[0];
            console.log('Выбран файл:', file.name);
            
            // Показываем прогресс бар
            progressContainer.classList.remove('hidden');
            progressBar.style.width = '0%';
            
            // Имитация загрузки файла
            let progress = 0;
            const interval = setInterval(() => {
                progress += 5;
                progressBar.style.width = `${progress}%`;
                
                if (progress >= 100) {
                    clearInterval(interval);
                    
                    // Имитация обработки на сервере
                    setTimeout(() => {
                        showModal(
                            'Успешная загрузка',
                            `Файл "${file.name}" успешно загружен и обработан. Сервер вернул следующие данные: <pre>{"status": "success", "records_processed": 42}</pre>`
                        );
                        progressContainer.classList.add('hidden');
                        uploadForm.reset();
                    }, 500);
                }
            }, 100);
        });
    }
}

// Обработчики для страницы "Администрирование"
function initAdminHandlers() {
    // Создание пользователя
    const createUserBtn = document.getElementById('create-user-btn');
    if (createUserBtn) {
        createUserBtn.addEventListener('click', () => {
            showModal(
                'Создание пользователя',
                `
                <form id="create-user-form">
                    <div class="form-group">
                        <label for="new-user-login">Логин</label>
                        <input type="text" id="new-user-login" required>
                    </div>
                    <div class="form-group">
                        <label for="new-user-password">Пароль</label>
                        <input type="password" id="new-user-password" required>
                    </div>
                    <div class="form-group">
                        <label for="new-user-role">Роль</label>
                        <select id="new-user-role" required>
                            <option value="">Выберите роль</option>
                            <option value="admin">Администратор</option>
                            <option value="operator">Оператор</option>
                            <option value="viewer">Наблюдатель</option>
                        </select>
                    </div>
                </form>
                `,
                () => {
                    const login = document.getElementById('new-user-login').value;
                    const password = document.getElementById('new-user-password').value;
                    const role = document.getElementById('new-user-role').value;
                    
                    if (!login || !password || !role) {
                        showModal('Ошибка', 'Заполните все поля');
                        return;
                    }
                    
                    console.log('Создание пользователя:', { login, password, role });
                    showModal('Успех', `Пользователь "${login}" успешно создан`);
                }
            );
        });
    }
    
    // Редактирование пользователя
    const editUserBtn = document.getElementById('edit-user-btn');
    if (editUserBtn) {
        editUserBtn.addEventListener('click', () => {
            showModal(
                'Редактирование пользователя',
                `
                <div class="form-group">
                    <label for="select-user">Выберите пользователя</label>
                    <select id="select-user" class="form-control">
                        <option value="">Выберите пользователя</option>
                        <option value="admin">admin</option>
                        <option value="user1">user1</option>
                    </select>
                </div>
                <div id="user-edit-form" class="hidden">
                    <div class="form-group">
                        <label for="edit-user-password">Новый пароль (оставьте пустым, чтобы не менять)</label>
                        <input type="password" id="edit-user-password" class="form-control">
                    </div>
                    <div class="form-group">
                        <label for="edit-user-role">Роль</label>
                        <select id="edit-user-role" class="form-control">
                            <option value="admin">Администратор</option>
                            <option value="operator">Оператор</option>
                            <option value="viewer">Наблюдатель</option>
                        </select>
                    </div>
                </div>
                `,
                () => {
                    const selectedUser = document.getElementById('select-user').value;
                    const newPassword = document.getElementById('edit-user-password').value;
                    const newRole = document.getElementById('edit-user-role').value;
                    
                    if (!selectedUser) {
                        showModal('Ошибка', 'Выберите пользователя');
                        return;
                    }
                    
                    console.log('Обновление пользователя:', { 
                        user: selectedUser, 
                        newPassword: newPassword || 'не изменен', 
                        newRole 
                    });
                    
                    showModal('Успех', `Данные пользователя "${selectedUser}" успешно обновлены`);
                }
            );
            
            // Показать форму редактирования при выборе пользователя
            const selectUser = document.getElementById('select-user');
            if (selectUser) {
                selectUser.addEventListener('change', function() {
                    const editForm = document.getElementById('user-edit-form');
                    if (this.value) {
                        editForm.classList.remove('hidden');
                    } else {
                        editForm.classList.add('hidden');
                    }
                });
            }
        });
    }
}

// Обработчики для страницы "Перевозки"
function initTransportationsHandlers() {
    const searchInput = document.getElementById('transportation-search');
    if (searchInput) {
        searchInput.addEventListener('input', function() {
            const searchTerm = this.value.toLowerCase();
            const rows = document.querySelectorAll('#transportations-table tbody tr');
            
            rows.forEach(row => {
                const text = row.textContent.toLowerCase();
                if (text.includes(searchTerm)) {
                    row.style.display = '';
                } else {
                    row.style.display = 'none';
                }
            });
        });
    }
}