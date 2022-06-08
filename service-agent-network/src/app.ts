import express from 'express';
import * as http from 'http';
import * as bodyparser from 'body-parser';
import * as dotenv from "dotenv";

//Load env form .env
dotenv.config();

import * as winston from 'winston';
import * as expressWinston from 'express-winston';
import cors from 'cors';
import debug from 'debug';

//Routes
import {CommonRoutesConfig} from './common/common.routes.config';
import {OrderRoute} from "./order/order.route.config";
import {UserApiRoute} from "./userApi/userApi.route.config";
import {SpOrderParcelRoute} from "./spOrderParcel/spOrderParcel.route.config";
import {ProductRoute} from "./product/product.route.config";
import {CodRoute} from "./cod/cod.route.config";
import {SerialNumberRoute} from "./serialNumber/serialNumber.route.config";
import {SlipRoute} from "./slip/slip.route.config";

const app: express.Application = express();
const server: http.Server = http.createServer(app);
const port: Number = process.env.PORT ? parseInt(process.env.PORT) : 3000;
const routes: Array<CommonRoutesConfig> = [];
const debugLog: debug.IDebugger = debug('app');

app.use(bodyparser.json());
app.use(cors());
app.use(expressWinston.logger({
    transports: [
        new winston.transports.Console()
    ],
    format: winston.format.combine(
        winston.format.colorize(),
        winston.format.json()
    )
}));

//Add a new route here
routes.push(new OrderRoute(app));
routes.push(new UserApiRoute(app));
routes.push(new SpOrderParcelRoute(app));
routes.push(new ProductRoute(app));
routes.push(new CodRoute(app));
routes.push(new SerialNumberRoute(app));
routes.push(new SlipRoute(app));

app.get('/', (req: express.Request, res: express.Response) => {
    res.status(200).send(`Server up and running!`)
});

server.listen(port, () => {
    debugLog(`Server running at http://localhost:${port}`);

    routes.forEach((route: CommonRoutesConfig) => {
        debugLog(`Routes configured for ${route.getName()}`);
    });
});