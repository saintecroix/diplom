const API_BASE = '/api'; // Базовый URL API

// Функция для загрузки Excel-файла
export async function uploadExcelFile(file, progressCallback) {
    const formData = new FormData();
    formData.append('file', file);

    return new Promise((resolve, reject) => {
        const xhr = new XMLHttpRequest();

        // Отслеживаем прогресс загрузки
        xhr.upload.addEventListener('progress', (event) => {
            if (event.lengthComputable) {
                const percent = Math.round((event.loaded / event.total) * 100);
                progressCallback(percent);
            }
        });

        // Обработка успешной загрузки
        xhr.addEventListener('load', () => {
            if (xhr.status >= 200 && xhr.status < 300) {
                try {
                    resolve(JSON.parse(xhr.response));
                } catch (e) {
                    reject(new Error('Ошибка разбора ответа сервера'));
                }
            } else {
                reject(new Error(`Ошибка сервера: ${xhr.status}`));
            }
        });

        // Обработка ошибок
        xhr.addEventListener('error', () => {
            reject(new Error('Сетевая ошибка'));
        });

        // Отправляем запрос
        xhr.open('POST', `${API_BASE}/upload`);
        xhr.send(formData);
    });
}