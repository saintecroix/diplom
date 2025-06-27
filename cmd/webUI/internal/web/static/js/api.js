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
            console.log('Response status:', xhr.status, 'Response:', xhr.responseText);

            try {
                const response = JSON.parse(xhr.responseText);

                if (xhr.status >= 200 && xhr.status < 300) {
                    resolve(response);
                } else {
                    const error = new Error(response.error || `Server error: ${xhr.status}`);
                    error.response = response;
                    reject(error);
                }
            } catch (e) {
                const error = new Error(`Invalid JSON response: ${xhr.responseText.substring(0, 100)}...`);
                error.response = xhr.responseText;
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
        xhr.send(formData);
    });
}