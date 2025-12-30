import SignupBuilder from '../builders/signupBuilder'

async function signup(ctx) {
	console.log("--> signup")

	const builder = new SignupBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

	await builder.addSettings()
	await builder.addUsername()
	await builder.addPassword()
	await builder.addSubmit()
	await builder.addFooter()
	await builder.addSidebar()

	ctx.body = await builder.build(ctx.state.error)
	ctx.status = 200;

	console.log("<-- signup")
}

export default signup;
