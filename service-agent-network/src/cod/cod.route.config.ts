import { CommonRoutesConfig } from "../common/common.routes.config";
import authenMiddleware from "../common/middleware/authen.middleware";
import codController from "./controller/cod.controller";

import express from "express";

export class CodRoute extends CommonRoutesConfig {
    constructor(app: express.Application) {
        super(app, 'CodRoute');
    }

    configureRoutes(): express.Application {
        this.app.route("/closed-api/cod/report")
            .get(
                codController.getCodReport,
            )

        return this.app;
    }
}