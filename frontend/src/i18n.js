import i18next from 'i18next';
import LanguageDetector from 'i18next-browser-languagedetector';
import {initReactI18next} from 'react-i18next';
import common_en from "./translations/en/common.json";
import common_ru from "./translations/ru/common.json";
import common_ua from "./translations/ua/common.json";

i18next
// load translation using xhr -> see /public/locales
// learn more: https://github.com/i18next/i18next-xhr-backend
// .use(Backend)
// detect user language
// learn more: https://github.com/i18next/i18next-browser-languageDetector
    .use(LanguageDetector)
    // pass the i18n instance to react-i18next.
    .use(initReactI18next)
    // init i18next
    // for all options read: https://www.i18next.com/overview/configuration-options
    .init({
        fallbackLng: 'en',
        debug: true,

        interpolation: {
            escapeValue: false, // not needed for react as it escapes by default
        },
        resources: {
            en: {common: common_en},
            ru: {common: common_ru},
            ua: {common: common_ua},
            uk: {common: common_ua},
        },
    });

export default i18next;
