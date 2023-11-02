import bodyParser from 'koa-bodyparser';

export default bodyParser({
  jsonLimit: '56kb',
});
