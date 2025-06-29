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
            // Удаляем отладочный вывод в продакшене
            // console.log('Response status:', xhr.status, 'Response:', xhr.responseText);

            try {
                const response = JSON.parse(xhr.responseText);

                if (xhr.status >= 200 && xhr.status < 300) {
                    // Успешный ответ (2xx)
                    if (response.status === "success") {
                        resolve(response);
                    } else {
                        // Обработка бизнес-ошибок (успешный HTTP но ошибка в логике)
                        const error = new Error(response.message || "Unknown server error");
                        error.response = response;
                        error.isBusinessError = true;
                        reject(error);
                    }
                } else {
                    // Ошибки HTTP (4xx, 5xx)
                    const error = new Error(response.message || `Server error: ${xhr.status}`);
                    error.response = response;
                    error.statusCode = xhr.status;
                    reject(error);
                }
            } catch (e) {
                // Ошибка парсинга JSON
                const error = new Error(`Invalid server response: ${xhr.responseText.substring(0, 100)}...`);
                error.responseText = xhr.responseText;
                error.statusCode = xhr.status;
                reject(error);
            }
        });

        xhr.addEventListener('error', () => {
            reject(new Error('Network error: Failed to send request'));
        });

        xhr.addEventListener('abort', () => {
            reject(new Error('Request aborted'));
        });

        xhr.timeout = 300000; // 5 минут
        xhr.ontimeout = () => {
            reject(new Error('Request timed out'));
        };

        xhr.open('POST', `${API_BASE}/upload`);
        xhr.setRequestHeader('Accept', 'application/json'); // Явно запрашиваем JSON
        xhr.send(formData);
    });
}

export function handleUploadError(error, statusDiv) {
    let errorMessage = 'Unknown error';

    if (error.message.includes('Network error')) {
        errorMessage = 'Сетевая ошибка. Проверьте подключение к интернету';
    } else if (error.message.includes('timed out')) {
        errorMessage = 'Превышено время ожидания ответа сервера';
    } else if (error.statusCode === 413) {
        errorMessage = 'Файл слишком большой. Максимальный размер 10 МБ';
    } else if (error.statusCode === 415) {
        errorMessage = 'Неподдерживаемый формат файла';
    } else if (error.isBusinessError) {
        errorMessage = error.message;
    } else {
        errorMessage = `Ошибка сервера: ${error.message}`;
    }

    return `<strong>Ошибка загрузки:</strong><br>${errorMessage}`;
}