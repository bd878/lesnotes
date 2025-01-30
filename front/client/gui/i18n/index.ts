const i18n = {
  "en": {
    logout: "Logout",
    login: "Login",
    register: "Register",
    loading: "Loading...",
    login_form_header: "Login:",
    register_form_header: "Register:",
    not_found: "Not found",
    auth_process: "Authenticating...",
    not_authed: "Not authed",
    messages_auth_fallback: "Please, log in first",
    index_intro: "Example project to send messages",
    messages_header: "Messages:",
    error_occured: "Error occured:",
    loading_messages_error: "Error occured while loading messages",
    name_required_err: "Name required",
    pass_required_err: "Password required",
    msg_required_err: "Message required",
    msg_send_text: "Send",
    load_more: "Load more",
    bad_status_code: "Bad status code",
    token_expired_error: "Token expired",
    bad_response: "Bad response",
    cannot_parse_response: "Cannot parse the response"
  },
  "ru": {
    logout: "Выйти",
    login: "Войти",
    register: "Зарегистрироваться",
    loading: "Загрузка...",
    login_form_header: "Войти:",
    register_form_header: "Зарегистрироваться:",
    not_found: "Не найдено",
    auth_process: "Загрузка...",
    not_authed: "Не аутентифицирован",
    messages_auth_fallback: "Пройдите аутентификацию",
    index_intro: "Пример проекта по отправке сообщений",
    messages_header: "Сообщения:",
    error_occured: "Возникла ошибка:",
    loading_messages_error: "Возникла ошибка во время загрузки сообщений",
    name_required_err: "Необходимо ввести имя пользователя",
    pass_required_err: "Необходимо ввести пароль",
    msg_required_err: "Необходимо ввести сообщение",
    msg_send_text: "Отправить",
    load_more: "Загрузить ещё",
    bad_status_code: "Ошибка запроса",
    token_expired_error: "Истекло время жизни токена",
    bad_response: "Невалидный ответ от сервера",
    cannot_parse_response: "Не получается разобрать ответ от сервера",
  }
}

const locale = "ru";

export default function(key) {
  return i18n[locale][key]
}
