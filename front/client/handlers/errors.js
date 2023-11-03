
export default async function middleware(ctx, next) {
  try {
    await next();
  } catch (e) {
    if (e.status) {
      ctx.body = e.message;
      ctx.status = e.status;
    } else {
      ctx.body = 'Internal error';
      ctx.status = 500;
      console.error(e.message, e.stack);
    }
  }
}
