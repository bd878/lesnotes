async function status(ctx) {
	console.log("--> status")

	ctx.body = 'ok\n';
	ctx.status = 200;

	console.log("<-- status")
}

export default status;
