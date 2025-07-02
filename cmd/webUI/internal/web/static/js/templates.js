export function getHomePage() {
  return `
      <div class="page-home">
        <h1 class="page-title">Главная</h1>
        <div class="stats-container">
          <div class="stat-card" style="background-color: var(--primary-color);">
            <span class="stat-value">578</span>
            <span class="stat-label">Всего вагонов</span>
            <i class="fas fa-train stat-icon"></i>
          </div>
          <div class="stat-card" style="background-color: var(--success-color);">
            <span class="stat-value">42</span>
            <span class="stat-label">Активных групп</span>
            <i class="fas fa-layer-group stat-icon"></i>
          </div>
          <div class="stat-card" style="background-color: var(--warning-color);">
            <span class="stat-value">128</span>
            <span class="stat-label">Перемещений за месяц</span>
            <i class="fas fa-exchange-alt stat-icon"></i>
          </div>
          <div class="stat-card" style="background-color: var(--danger-color);">
            <span class="stat-value">15</span>
            <span class="stat-label">Требуют внимания</span>
            <i class="fas fa-exclamation-circle stat-icon"></i>
          </div>
        </div>
        
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

export function getWagonSearchPage() {
  return `
      <div class="page-wagon-search">
        <h1 class="page-title">Поиск и управление вагонами</h1>

        <div class="card">
          <h2 class="card-title">Поиск вагонов</h2>
          <form id="wagon-search-form" class="form">
            <div class="form-group">
              <label for="wagon-number">Номер вагона</label>
              <input type="text" id="wagon-number" placeholder="Введите номер вагона">
            </div>
            <div class="form-group">
              <label for="wagon-status">Статус</label>
              <select id="wagon-status">
                <option value="">Все статусы</option>
                <option value="active">Активный</option>
                <option value="maintenance">На обслуживании</option>
                <option value="reserved">Зарезервирован</option>
              </select>
            </div>
            <div class="form-group">
              <label for="wagon-location">Местоположение</label>
              <input type="text" id="wagon-location" placeholder="Введите местоположение">
            </div>
            <div class="form-actions">
              <button type="submit" class="btn btn-primary">
                <i class="fas fa-search"></i> Поиск
              </button>
              <button type="reset" class="btn btn-secondary">
                <i class="fas fa-undo"></i> Сбросить
              </button>
            </div>
          </form>
        </div>

        <div class="card">
          <h2 class="card-title">Результаты поиска</h2>
          <div class="table-container">
            <table id="wagon-results-table">
              <thead>
                <tr>
                  <th>Номер</th>
                  <th>Тип</th>
                  <th>Статус</th>
                  <th>Местоположение</th>
                  <th>Последнее обновление</th>
                  <th>Действия</th>
                </tr>
              </thead>
              <tbody>
                <!-- Данные будут заполнены JavaScript -->
              </tbody>
            </table>
          </div>
        </div>
      </div>
    `;
}

export function getGroupManagePage() {
  return `
      <div class="page-group-manage">
        <h1 class="page-title">Управление группами вагонов</h1>
        
        <div class="card">
          <h2 class="card-title">Активные группы</h2>
        </div>

        <div class="card">
          <h2 class="card-title">Создание группы вагонов</h2>
          <form id="create-group-form" class="form">
            <div class="form-group">
              <label for="group-name">Название группы</label>
              <input type="text" id="group-name" placeholder="Введите название группы">
            </div>
            <div class="form-group">
              <label for="wagon-numbers">Номера вагонов</label>
              <textarea id="wagon-numbers" rows="4" placeholder="Введите номера вагонов через запятую или каждый с новой строки"></textarea>
            </div>
            <div class="form-actions">
              <button type="submit" class="btn btn-primary">
                <i class="fas fa-plus"></i> Создать группу
              </button>
            </div>
          </form>
        </div>
        
      </div>
    `;
}

export function getDataInputPage() {
  return `
    <div class="container mt-5">
      <div class="row justify-content-center">
        <div class="col-md-8">
          <div class="card shadow-sm">
            <div class="card-header bg-primary text-white">
              <h5 class="card-title mb-0">
                <i class="bi bi-cloud-upload me-2"></i>Загрузка данных о рейсах
              </h5>
            </div>
            <div class="card-body">
              <div id="uploadContainer">
                <form id="uploadForm" class="mb-4">
                  <div class="mb-3">
                    <label for="excelFile" class="form-label">Выберите файл Excel</label>
                    <div class="input-group">
                      <input class="form-control" type="file" id="excelFile" accept=".xlsx, .xls" required>
                    </div>
                    <div class="form-text">Поддерживаются файлы формата .xlsx и .xls</div>
                  </div>
                  
                  <div class="d-grid gap-2">
                    <button type="submit" class="btn btn-success btn-lg">
                      <i class="bi bi-upload me-2"></i>Загрузить данные
                    </button>
                  </div>
                </form>
                
                <div id="uploadProgress" class="progress mb-3" style="display: none; height: 25px;">
                  <div id="progressBar" class="progress-bar progress-bar-striped progress-bar-animated" 
                       role="progressbar" style="width: 0%">0%</div>
                </div>
                
                <div id="uploadStatus" class="alert" style="display: none;"></div>
              </div>
              
              <div class="mt-4">
                <h6><i class="bi bi-info-circle me-2"></i>Инструкция по загрузке:</h6>
                <ul class="small">
                  <li>Файл должен содержать данные о рейсах в установленном формате</li>
                  <li>Первая строка должна содержать заголовки столбцов</li>
                  <li>Максимальный размер файла: 10 МБ</li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  `;
}

export function getAdminPage() {
  return `
      <div class="page-admin">
        <h1 class="page-title">Администрирование</h1>
        
        <div class="card">
          <h2 class="card-title">Управление пользователями</h2>
          <div class="form-actions">
            <button id="add-user-btn" class="btn btn-primary">
              <i class="fas fa-user-plus"></i> Добавить пользователя
            </button>
          </div>
          <div class="table-container">
            <table>
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Имя пользователя</th>
                  <th>Email</th>
                  <th>Роль</th>
                  <th>Статус</th>
                  <th>Действия</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td>1</td>
                  <td>admin</td>
                  <td>admin@example.com</td>
                  <td>Администратор</td>
                  <td><span class="badge success">Активен</span></td>
                  <td>
                    <button class="btn btn-secondary btn-sm edit-user" data-id="1">
                      <i class="fas fa-edit"></i>
                    </button>
                    <button class="btn btn-danger btn-sm delete-user" data-id="1">
                      <i class="fas fa-trash"></i>
                    </button>
                  </td>
                </tr>
                <tr>
                  <td>2</td>
                  <td>operator</td>
                  <td>operator@example.com</td>
                  <td>Оператор</td>
                  <td><span class="badge success">Активен</span></td>
                  <td>
                    <button class="btn btn-secondary btn-sm edit-user" data-id="2">
                      <i class="fas fa-edit"></i>
                    </button>
                    <button class="btn btn-danger btn-sm delete-user" data-id="2">
                      <i class="fas fa-trash"></i>
                    </button>
                  </td>
                </tr>
                <tr>
                  <td>3</td>
                  <td>viewer</td>
                  <td>viewer@example.com</td>
                  <td>Наблюдатель</td>
                  <td><span class="badge warning">Заблокирован</span></td>
                  <td>
                    <button class="btn btn-secondary btn-sm edit-user" data-id="3">
                      <i class="fas fa-edit"></i>
                    </button>
                    <button class="btn btn-danger btn-sm delete-user" data-id="3">
                      <i class="fas fa-trash"></i>
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
        
<!--        <div class="card">-->
<!--          <h2 class="card-title">Системные настройки</h2>-->
<!--          <form id="settings-form" class="form">-->
<!--            <div class="form-group">-->
<!--              <label for="setting-company">Название компании</label>-->
<!--              <input type="text" id="setting-company" value="ООО 'ЖелДорТранс'">-->
<!--            </div>-->
<!--            <div class="form-group">-->
<!--              <label for="setting-timeout">Таймаут сессии (минуты)</label>-->
<!--              <input type="number" id="setting-timeout" value="30" min="5" max="120">-->
<!--            </div>-->
<!--            <div class="form-group">-->
<!--              <label for="setting-theme">Тема оформления</label>-->
<!--              <select id="setting-theme">-->
<!--                <option value="light">Светлая</option>-->
<!--                <option value="dark">Темная</option>-->
<!--                <option value="auto">Системная</option>-->
<!--              </select>-->
<!--            </div>-->
<!--            <div class="form-actions">-->
<!--              <button type="submit" class="btn btn-primary">-->
<!--                <i class="fas fa-save"></i> Сохранить настройки-->
<!--              </button>-->
<!--            </div>-->
<!--          </form>-->
<!--        </div>-->
        
        <div class="card">
          <h2 class="card-title">Журнал системы</h2>
          <div class="form-group">
            <label for="log-level">Уровень логирования</label>
            <select id="log-level">
              <option value="error">Только ошибки</option>
              <option value="warning">Предупреждения и ошибки</option>
              <option value="info" selected>Информационные сообщения</option>
              <option value="debug">Отладочная информация</option>
            </select>
          </div>
          <div class="log-container">
            <div class="log-entry error">
              <span class="log-time">2023-06-15 14:32:45</span>
              <span class="log-level">ОШИБКА</span>
              <span class="log-message">Не удалось подключиться к базе данных</span>
            </div>
            <div class="log-entry warning">
              <span class="log-time">2023-06-15 14:30:12</span>
              <span class="log-level">ПРЕДУПРЕЖДЕНИЕ</span>
              <span class="log-message">Попытка неавторизованного доступа (IP: 192.168.1.105)</span>
            </div>
            <div class="log-entry info">
              <span class="log-time">2023-06-15 14:28:30</span>
              <span class="log-level">ИНФО</span>
              <span class="log-message">Пользователь admin вошел в систему</span>
            </div>
          </div>
        </div>
      </div>
    `;
}

export function getTransportationsPage() {
  return `
      <div class="page-transportations">
        <h1 class="page-title">Перевозки</h1>
        
        <div class="card">
          <h2 class="card-title">Активные перевозки</h2>
          <div class="table-container">
            <table>
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Маршрут</th>
                  <th>Группа вагонов</th>
                  <th>Статус</th>
                  <th>Дата отправления</th>
                  <th>Ожидаемая дата прибытия</th>
                  <th>Действия</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td>1</td>
                  <td>Москва → Санкт-Петербург</td>
                  <td>Группа А (15 вагонов)</td>
                  <td><span class="badge warning">В пути</span></td>
                  <td>2023-06-10</td>
                  <td>2023-06-16</td>
                  <td>
                    <button class="btn btn-secondary btn-sm view-transportation" data-id="1">
                      <i class="fas fa-eye"></i>
                    </button>
                    <button class="btn btn-primary btn-sm update-status" data-id="1">
                      <i class="fas fa-edit"></i>
                    </button>
                  </td>
                </tr>
                <tr>
                  <td>2</td>
                  <td>Екатеринбург → Новосибирск</td>
                  <td>Группа B (8 вагонов)</td>
                  <td><span class="badge warning">В пути</span></td>
                  <td>2023-06-12</td>
                  <td>2023-06-18</td>
                  <td>
                    <button class="btn btn-secondary btn-sm view-transportation" data-id="2">
                      <i class="fas fa-eye"></i>
                    </button>
                    <button class="btn btn-primary btn-sm update-status" data-id="2">
                      <i class="fas fa-edit"></i>
                    </button>
                  </td>
                </tr>
                <tr>
                  <td>3</td>
                  <td>Краснодар → Ростов-на-Дону</td>
                  <td>Группа C (5 вагонов)</td>
                  <td><span class="badge success">Прибыл</span></td>
                  <td>2023-06-14</td>
                  <td>2023-06-15</td>
                  <td>
                    <button class="btn btn-secondary btn-sm view-transportation" data-id="3">
                      <i class="fas fa-eye"></i>
                    </button>
                    <button class="btn btn-primary btn-sm update-status" data-id="3">
                      <i class="fas fa-edit"></i>
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
        
        <div class="card">
          <h2 class="card-title">Создать новую перевозку</h2>
          <form id="transportation-form" class="form">
            <div class="form-group">
              <label for="transportation-from">Пункт отправления</label>
              <input type="text" id="transportation-from" placeholder="Введите пункт отправления">
            </div>
            <div class="form-group">
              <label for="transportation-to">Пункт назначения</label>
              <input type="text" id="transportation-to" placeholder="Введите пункт назначения">
            </div>
            <div class="form-group">
              <label for="transportation-group">Группа вагонов</label>
              <select id="transportation-group">
                <option value="">Выберите группу вагонов</option>
                <option value="1">Группа А (15 вагонов)</option>
                <option value="2">Группа B (8 вагонов)</option>
                <option value="3">Группа C (5 вагонов)</option>
                <option value="4">Группа D (12 вагонов)</option>
              </select>
            </div>
            <div class="form-group">
              <label for="transportation-departure">Дата отправления</label>
              <input type="date" id="transportation-departure">
            </div>
            <div class="form-group">
              <label for="transportation-arrival">Ожидаемая дата прибытия</label>
              <input type="date" id="transportation-arrival">
            </div>
            <div class="form-group">
              <label for="transportation-notes">Примечания</label>
              <textarea id="transportation-notes" rows="3" placeholder="Дополнительная информация"></textarea>
            </div>
            <div class="form-actions">
              <button type="submit" class="btn btn-primary">
                <i class="fas fa-plus"></i> Создать перевозку
              </button>
              <button type="reset" class="btn btn-secondary">
                <i class="fas fa-undo"></i> Сбросить
              </button>
            </div>
          </form>
        </div>
        
        <div class="card">
          <h2 class="card-title">История перевозок</h2>
          <div class="form-group">
            <label for="history-period">Период</label>
            <select id="history-period">
              <option value="7">Последние 7 дней</option>
              <option value="30" selected>Последние 30 дней</option>
              <option value="90">Последние 3 месяца</option>
              <option value="180">Последние 6 месяцев</option>
              <option value="365">Последний год</option>
            </select>
          </div>
          <div class="table-container">
            <table id="history-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Маршрут</th>
                  <th>Группа вагонов</th>
                  <th>Статус</th>
                  <th>Дата отправления</th>
                  <th>Дата прибытия</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td>4</td>
                  <td>Москва → Казань</td>
                  <td>Группа D (12 вагонов)</td>
                  <td><span class="badge success">Завершено</span></td>
                  <td>2023-05-28</td>
                  <td>2023-06-02</td>
                </tr>
                <tr>
                  <td>5</td>
                  <td>Санкт-Петербург → Мурманск</td>
                  <td>Группа E (10 вагонов)</td>
                  <td><span class="badge success">Завершено</span></td>
                  <td>2023-05-25</td>
                  <td>2023-05-31</td>
                </tr>
                <tr>
                  <td>6</td>
                  <td>Владивосток → Хабаровск</td>
                  <td>Группа F (7 вагонов)</td>
                  <td><span class="badge danger">Отменено</span></td>
                  <td>2023-05-20</td>
                  <td>-</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    `;
}

export function getUserEditForm(user) {
  return `
      <form id="edit-user-form" class="form">
        <div class="form-group">
          <label for="edit-username">Имя пользователя</label>
          <input type="text" id="edit-username" name="username" value="${user.username || ''}" required>
          <span class="error-message"></span>
        </div>
        <div class="form-group">
          <label for="edit-full-name">ФИО</label>
          <input type="text" id="edit-full-name" name="fullName" value="${user.fullName || ''}" required>
          <span class="error-message"></span>
        </div>
        <div class="form-group">
          <label for="edit-role">Роль</label>
          <select id="edit-role" name="role" required>
            <option value="admin" ${user.role === 'admin' ? 'selected' : ''}>Администратор</option>
            <option value="operator" ${user.role === 'operator' ? 'selected' : ''}>Оператор</option>
            <option value="viewer" ${user.role === 'viewer' ? 'selected' : ''}>Просмотр</option>
          </select>
          <span class="error-message"></span>
        </div>
        <div class="form-group">
          <label for="edit-password">Новый пароль (оставьте пустым, чтобы не менять)</label>
          <input type="password" id="edit-password" name="password">
          <span class="error-message"></span>
        </div>
        <div class="form-group">
          <label for="edit-confirm-password">Подтвердите пароль</label>
          <input type="password" id="edit-confirm-password" name="confirmPassword">
          <span class="error-message"></span>
        </div>
        <input type="hidden" id="edit-user-id" value="${user.id || ''}">
        <div class="form-actions">
          <button type="submit" class="btn btn-primary">Сохранить</button>
          <button type="button" id="cancel-edit-user" class="btn btn-secondary">Отмена</button>
        </div>
      </form>
    `;
}