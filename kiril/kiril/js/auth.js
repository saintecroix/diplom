export function getToken() {
  return localStorage.getItem('jwtToken');
}

export function setToken(token) {
  localStorage.setItem('jwtToken', token);
}

export function clearToken() {
  localStorage.removeItem('jwtToken');
}

export function isAuthenticated() {
  return !!getToken();
}