import renderDesktop from './desktop'
import renderMobile from './mobile'

async function mobileOrDesktop(ctx) {
	console.log("--> home")

	if (ctx.userAgent.isMobile)
		await renderMobile(ctx)
	else
		await renderDesktop(ctx)

	console.log("<-- home")
}

export default mobileOrDesktop
