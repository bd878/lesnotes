import Koa from 'koa';
import Router from '@koa/router';

import helmet from './handlers/helmet.js';
import errors from './handlers/errors.js';
import logger from './handlers/logger.js';
import bodyParser from './handlers/bodyParser.js';
import filename from './handlers/filename.js';
import useragent from './handlers/useragent.js';
import favicon from './handlers/favicon.js';

import assets from './routes/assets.js';
import client from './routes/client.js';

const app = new Koa();
const router = new Router();

app.use(helmet);
app.use(errors);
app.use(logger);
app.use(bodyParser);
app.use(filename);
app.use(useragent);
app.use(favicon);

router
  .get('/public/', assets)
  .get('/', client);

app.use(router.routes());

export default app;
