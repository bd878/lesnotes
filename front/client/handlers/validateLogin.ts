import login from '../routes/login'
import * as is from '../third_party/is'
import api from '../api';

async function validateLogin(ctx, next) {
	ctx.log.info("--> validateLogin")

	const response = await processLogin(ctx)

	if (response.error.error) {
		ctx.log.error(response.error)
		ctx.state.error = response.error.human
		await login(ctx)
	} else {
		const expiresAt = new Date(response.expiresAt)
		ctx.log.info("expiresAt", expiresAt.toString())
		ctx.set({"Set-Cookie":  "token=" + response.token + "; Expires=" + expiresAt.toString() + "; HttpOnly; Path=/; Secure; Domain=" + `${DOMAIN}`})

		await next()
	}

	ctx.log.info("<-- validateLogin")
}

async function processLogin(ctx) {
	let form = ctx.request.body

	if (is.empty(form)) {
		// let backend validate
		form = {}
	}

	const params = new URLSearchParams(ctx.search)

	return await api.login(form.login, form.password, params.get("lang"))
}

export default validateLogin;
