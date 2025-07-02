import {getToken} from "./auth";

const API_BASE = '/api';

export async function uploadExcelFile(file, progressCallback) {
    return new Promise((resolve, reject) => {
        const xhr = new XMLHttpRequest();
        const formData = new FormData();
        formData.append('file', file);

        xhr.upload.addEventListener('progress', (event) => {
            if (event.lengthComputable && progressCallback) {
                const percent = Math.round((event.loaded / event.total) * 100);
                progressCallback(percent);
            }
        });

        xhr.addEventListener('load', () => {
            try {
                const response = xhr.responseText ? JSON.parse(xhr.responseText) : {};

                if (xhr.status >= 200 && xhr.status < 300) {
                    if (response.status === "success") {
                        resolve({
                            ...response,
                            fileName: file.name,
                            fileSize: file.size
                        });
                    } else {
                        // Обработка бизнес-ошибок из gRPC
                        const error = new Error(response.message || "Ошибка обработки файла");
                        error.type = "business";
                        error.details = response.error || "Неизвестная ошибка сервера"; // Используем поле error из gRPC
                        error.jobId = response.job_id; // Добавляем jobId для отслеживания
                        reject(error);
                    }
                } else {
                    const error = new Error(
                        response.message ||
                        `Ошибка сервера (${xhr.status}: ${xhr.statusText})`
                    );
                    error.type = "http";
                    error.status = xhr.status;
                    reject(error);
                }
            } catch (e) {
                const error = new Error("Неверный формат ответа сервера");
                error.type = "parse";
                error.originalError = e;
                reject(error);
            }
        });

        xhr.addEventListener('error', () => {
            reject(new Error('Сетевая ошибка при отправке запроса'));
        });

        xhr.addEventListener('abort', () => {
            reject(new Error('Запрос был отменён'));
        });

        xhr.timeout = 300000;
        xhr.ontimeout = () => {
            reject(new Error('Превышено время ожидания ответа'));
        };

        xhr.open('POST', `${API_BASE}/upload`);
        xhr.setRequestHeader('Accept', 'application/json');
        xhr.send(formData);
    });
}

export function handleUploadError(error) {
    let message = 'Произошла ошибка при загрузке файла';
    let details = '';

    switch (error.type) {
        case 'business':
            message = 'Ошибка обработки файла';
            details = error.details ? `<br><small>${error.details}</small>` : '';
            if (error.jobId) {
                details += `<br><small>ID задачи: ${error.jobId}</small>`;
            }
            break;
        case 'http':
            message = `Ошибка сервера (${error.status || 'неизвестный статус'})`;
            details = error.message ? `<br><small>${error.message}</small>` : '';
            break;
        case 'parse':
            message = 'Ошибка обработки ответа сервера';
            details = '<br><small>Сервер вернул некорректные данные</small>';
            break;
        default:
            details = error.message ? `<br><small>${error.message}</small>` : '';
    }

    return {
        title: message,
        details: details,
        isError: true
    };
}

async function getAllUsers() {
    const userToken = getToken(); // Получаем токен из хранилища

    try {
        const response = await fetch('/getUsers', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${userToken}`
            }
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.message || `Ошибка сервера: ${response.status}`);
        }

        const data = await response.json();

        // Валидация структуры ответа
        if (!data.users || !Array.isArray(data.users)) {
            throw new Error('Некорректный формат данных: отсутствует массив users');
        }

        // Проверяем структуру каждого пользователя
        const isValidUser = (user) => {
            return user &&
                typeof user.id === 'number' &&
                typeof user.name === 'string' &&
                typeof user.mail === 'string' &&
                typeof user.role === 'string';
        };

        if (!data.users.every(isValidUser)) {
            throw new Error('Некорректная структура данных пользователей');
        }

        return data;

    } catch (error) {
        console.error('Ошибка при получении пользователей:', error);
        throw new Error(`Не удалось загрузить пользователей: ${error.message}`);
    }
}