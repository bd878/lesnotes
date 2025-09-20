import renderDesktop from './desktop'
import renderMobile from './mobile'

async function mobileOrDesktop(ctx) {
	if (ctx.userAgent.isMobile)
		await renderMobile(ctx)
	else
		await renderDesktop(ctx)
}

export default mobileOrDesktop
