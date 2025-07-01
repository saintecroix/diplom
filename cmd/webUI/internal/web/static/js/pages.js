import { elements, showLoader, hideLoader } from './dom.js';
import { initPageHandlers } from './handlers.js';

import {
  getHomePage,
  getWagonSearchPage,
  getGroupManagePage,
  getDataInputPage,
  getAdminPage,
  getTransportationsPage
} from './templates.js';

export function loadPage(page) {
  showLoader();
  setTimeout(() => {
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
  }, 500);
}
