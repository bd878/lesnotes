const i18n = {
	"en": {
		logout: "Logout",
		login: "Login",
		register: "Register",
		loading: "Loading...",
		username: "Username",
		password: "Password",
		authed: "Authed",
		search: "Search...",
		text: "Text",
		file: "File",
		send: "Send",
		delete: "Delete",
		edit: "Edit",
		publish: "Publish",
		private: "Revoke link",
		update: "Save",
		cancel: "Cancel",
		miniapp_only: "This page should be opened in Telegram only",
		text_placeholder: "Content...",
		title_placeholder: "Title..."
	},
	"ru": {
		logout: "Выйти",
		login: "Войти",
		username: "Логин",
		password: "Пароль",
		search: "Поиск...",
		text: "Текст",
		file: "Файл",
		send: "Отправить",
		attachments: "Файлы",
		register: "Зарегистрироваться",
		loading: "Загрузка...",
		delete: "Удалить",
		edit: "Редактировать",
		publish: "Опубликовать",
		private: "Скрыть от всех",
		update: "Сохранить",
		cancel: "Отмена",
		miniapp_only: "Эта страница работает только из приложения Telegram",
		text_placeholder: "Сообщение...",
		title_placeholder: "Заголовок..."
	}
}

const locale = "ru";

export default function(key) {
	return i18n[locale][key]
}
