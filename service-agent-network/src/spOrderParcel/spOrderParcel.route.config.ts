import {CommonRoutesConfig} from "../common/common.routes.config";
import authenMiddleware from "../common/middleware/authen.middleware";
import spOrderParcelController from "./controller/spOrderParcel.controller";

import express from "express";

export class SpOrderParcelRoute extends CommonRoutesConfig {
    constructor(app: express.Application) {
        super(app, 'SpOrderParcelRoute');
    }

    configureRoutes(): express.Application {
        this.app.route('/closed-api/shipping/order')
            .post([
                authenMiddleware.verifyToken,
                spOrderParcelController.makeOrder,
            ])

        this.app.route("/closed/hook/order-weight")
            .post([
                spOrderParcelController.receiveOrderWeight
            ])

        this.app.route("/closed/hook/order-status")
            .post([
                spOrderParcelController.receiveOrderStatus,
            ])

        return this.app;
    }
}