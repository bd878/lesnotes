async function status(ctx) {
	ctx.log.info("--> status")

	ctx.body = 'ok\n';
	ctx.status = 200;

	ctx.log.info("<-- status")
}

export default status;
