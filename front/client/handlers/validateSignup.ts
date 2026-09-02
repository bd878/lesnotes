import signup from '../routes/signup'
import * as is from '../third_party/is'
import api from '../api';

async function validateSignup(ctx, next) {
	ctx.log.info("--> validateSignup")

	const response = await processSignup(ctx)

	if (response.error.error) {
		ctx.log.error(response.error)
		ctx.state.error = response.error.human
		await signup(ctx)
	} else {
		const expiresAt = new Date(response.expiresAt)
		ctx.log.info("expiresAt", expiresAt.toString())
		ctx.set({"Set-Cookie":  "token=" + response.token + "; Expires=" + expiresAt.toString() + "; HttpOnly; Path=/; Secure; Domain=" + `${DOMAIN}`})

		await next()
	}

	ctx.log.info("<-- validateSignup")
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
