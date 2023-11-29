import Koa from 'koa';
import Router from '@koa/router';
import Config from 'config';
import helmet from './handlers/helmet.js';
import errors from './handlers/errors.js';
import logger from './handlers/logger.js';
import bodyParser from './handlers/bodyParser.js';
import useragent from './handlers/useragent.js';
import favicon from './handlers/favicon.js';

import assets from './routes/assets.js';
import index from './routes/index.js';

const app = new Koa();
const router = new Router();

app.use(helmet);
app.use(errors);
app.use(logger);
app.use(bodyParser);
app.use(useragent);
app.use(favicon);

router
  .get('/public/:filename', assets)
  .get('/', index)

app.use(router.routes());

const port = process.env.PORT || Config.get("port");

app.listen(port, () => {
  console.log(`App is listening on ${port} port`);
});
