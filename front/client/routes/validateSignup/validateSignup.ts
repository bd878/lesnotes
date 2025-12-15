import signup from '../signup/signup'
import * as is from '../../third_party/is'
import api from '../../api';

async function validateSignup(ctx, next) {
	console.log("--> validate signup")

	const response = await processSignup(ctx)

	if (response.error.error) {
		ctx.state.error = response.error.human
		return await signup(ctx)
	}

	const expiresAt = new Date(Math.round(response.expiresAt / 1_000_000))
	console.log("expiresAt", expiresAt.toString())
	ctx.set({"Set-Cookie":  "token=" + response.token + "; Expires=" + expiresAt.toString() + "; HttpOnly; Path=/; Secure; Domain=" + `${DOMAIN}`})

	await next()

	console.log("<-- validate signup")
}

async function processSignup(ctx) {
	let form = ctx.request.body


	if (is.empty(form)) {
		// let backend validate
		form = {}
	}

	const params = new URLSearchParams(ctx.search)

	return await api.signup(form.login, form.password, params.get("lang"))
}

export default validateSignup
