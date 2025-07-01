import { elements, showLoader, hideLoader } from './dom.js';
import { initPageHandlers } from './handlers.js';
import {
  getHomePage,
  getWagonSearchPage,
  getDataInputPage,
  getAdminPage,
  getTransportationsPage,
  getGroupManagePage
} from './templates.js';

export function setupRouter() {
  const navLinks = document.querySelectorAll('.sidebar-nav a');
  navLinks.forEach(link => {
    link.addEventListener('click', (e) => {
      e.preventDefault();
      navLinks.forEach(l => l.classList.remove('active'));
      link.classList.add('active');
      loadPage(link.getAttribute('data-page'));
    });
  });
}

export function loadPage(page) {
  showLoader();
  
  let content = '';
  switch (page) {
    case 'home':
      content = getHomePage();
      break;
    case 'wagon-search':
      content = getWagonSearchPage();
      break;
    case 'group-manage':
      content = getGroupManagePage();
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

  elements.pageContent.innerHTML = content;
  hideLoader();
  initPageHandlers(page);
}