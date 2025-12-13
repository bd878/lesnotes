import login from '../login/login'
import * as is from '../../third_party/is'
import api from '../../api';

async function validateLogin(ctx, next) {
	console.log("--> validate login")

	const response = await processLogin(ctx)

	if (response.error.error) {
		ctx.state.error = response.error.human
		return await login(ctx)
	}

	const expiresAt = new Date(Math.round(response.expiresAt / 1_000_000))
	console.log("expiresAt", expiresAt.toString())
	ctx.set({"Set-Cookie":  "token=" + response.token + "; Expires=" + expiresAt.toString() + "; HttpOnly; Path=/; Secure; Domain=" + `${DOMAIN}`})

	await next()

	console.log("<-- validate login")
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
