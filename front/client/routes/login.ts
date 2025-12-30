import LoginBuilder from '../builders/loginBuilder';

async function login(ctx) {
	console.log("--> login")

	const builder = new LoginBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

	await builder.addSettings()
	await builder.addUsername()
	await builder.addPassword()
	await builder.addSubmit()
	await builder.addFooter()
	await builder.addSidebar()

	ctx.body = await builder.build(ctx.state.error)
	ctx.status = 200;

	console.log("<-- login")
}

export default login;
