import {CommonRoutesConfig} from "../common/common.routes.config";
import authenMiddleware from "../common/middleware/authen.middleware";
import serialNumberController from "./controller/serialNumber.controller";

import express from "express";

export class SerialNumberRoute extends CommonRoutesConfig {
    constructor(app: express.Application) {
        super(app, 'OrderRoute');
    }

    configureRoutes(): express.Application {
        this.app.route("/closed-api/order/order-product/serial-number")
            .post(
                authenMiddleware.verifyToken,
                serialNumberController.createSerialNumber,
            )
            .put(
                authenMiddleware.verifyToken,
                serialNumberController.updateSerialNumber,
            )
            .delete(
                authenMiddleware.verifyToken,
                serialNumberController.deleteSerialNumber,
            )

        return this.app;
    }
}