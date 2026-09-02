async function xxx(ctx) {
	ctx.log.info("--> xxx")

	ctx.body = "<html>Pas de template</html>";
	ctx.status = 500;

	ctx.log.info("<-- xxx")
}

export default xxx;
