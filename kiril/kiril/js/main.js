import { elements, showModal, showApp, showLogin } from './dom.js';
import { isAuthenticated } from './auth.js';
import { setupRouter, loadPage } from './router.js';
import { setupEvents } from './events.js';

document.addEventListener('DOMContentLoaded', () => {
  // Настройка функциональности модального окна
  const modalClose = document.getElementById('modal-close');
  const modalCancel = document.getElementById('modal-cancel');
  const modalOverlay = document.getElementById('modal-overlay');
  
  if (modalClose) {
    modalClose.addEventListener('click', () => {
      modalOverlay.classList.add('hidden');
    });
  }
  
  if (modalCancel) {
    modalCancel.addEventListener('click', () => {
      modalOverlay.classList.add('hidden');
    });
  }
  
  setupEvents();
  setupRouter();

  // Проверка авторизации при загрузке страницы
  if (isAuthenticated()) {
    showApp('Администратор');
    loadPage('home');
  } else {
    showLogin();
  }
});