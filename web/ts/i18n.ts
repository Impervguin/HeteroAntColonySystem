type Translations = { [key: string]: string };

let currentLanguage: string = 'en';
let translations: Translations = {};

export async function initI18n(): Promise<void> {
  const storedLang = localStorage.getItem('language');
  if (storedLang && (storedLang === 'en' || storedLang === 'ru')) {
    currentLanguage = storedLang;
  } else {
    // Detect browser language
    const browserLang = navigator.language.split('-')[0];
    currentLanguage = browserLang === 'ru' ? 'ru' : 'en';
  }
  await loadTranslations(currentLanguage);
}

export async function setLanguage(lang: string): Promise<void> {
  if (lang !== 'en' && lang !== 'ru') return;
  currentLanguage = lang;
  localStorage.setItem('language', lang);
  await loadTranslations(lang);
  // Trigger re-render or update UI elements that use translations
  updateUI();
}

export function getCurrentLanguage(): string {
  return currentLanguage;
}

export function t(key: string, placeholders?: { [key: string]: string | number }): string {
  let text = translations[key] || key;
  if (placeholders) {
    for (const [placeholder, value] of Object.entries(placeholders)) {
      text = text.replace(new RegExp(`\\{${placeholder}\\}`, 'g'), String(value));
    }
  }
  return text;
}

async function loadTranslations(lang: string): Promise<void> {
  try {
    const response = await fetch(`/static/locales/${lang}.json`);
    translations = await response.json();
  } catch (error) {
    console.error(`Failed to load translations for ${lang}:`, error);
    // Fallback to English
    translations = {};
  }
}

function updateUI(): void {
  // Update elements that have data-i18n attribute
  document.querySelectorAll('[data-i18n]').forEach(el => {
    const key = el.getAttribute('data-i18n');
    if (key) {
      el.textContent = t(key);
    }
  });

  // Update placeholders
  document.querySelectorAll('[data-i18n-placeholder]').forEach(el => {
    const key = el.getAttribute('data-i18n-placeholder');
    if (key && el instanceof HTMLInputElement) {
      el.placeholder = t(key);
    }
  });

  // Update select options
  document.querySelectorAll('[data-i18n-option]').forEach(el => {
    const key = el.getAttribute('data-i18n-option');
    if (key && el instanceof HTMLOptionElement) {
      el.textContent = t(key);
    }
  });

  // Update labels
  document.querySelectorAll('label[data-i18n]').forEach(label => {
    const key = label.getAttribute('data-i18n');
    if (key) {
      label.textContent = t(key);
    }
  });
}