async function xxx(ctx) {
	console.log("--> xxx")

	ctx.body = "<html>Pas de template</html>";
	ctx.status = 500;

	console.log("<-- xxx")
}

export default xxx;
