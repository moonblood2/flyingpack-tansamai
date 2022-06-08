import {CommonRoutesConfig} from "../common/common.routes.config";
import authenMiddleware from "../common/middleware/authen.middleware";

import slipController from "./controller/slip.controller";

import express from "express";
import debug from "debug";

const log: debug.IDebugger = debug("slip-route-config:");

export class SlipRoute extends CommonRoutesConfig {
    constructor(app: express.Application) {
        super(app, 'SlipRoute');
    }

    configureRoutes(): express.Application {
        this.app.route("/closed-api/slip")
            .get(
                authenMiddleware.verifyToken,
                slipController.getSlip,
            )

        return this.app;
    }
}