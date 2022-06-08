import {CommonRoutesConfig} from "../common/common.routes.config";
import userApiController from "./controller/userApi.controller";

import express from "express";
import debug from "debug";

const log: debug.IDebugger = debug("userApi-route-config:");

export class UserApiRoute extends CommonRoutesConfig {
    constructor(app: express.Application) {
        super(app, 'UserApiRoute');
    }

    configureRoutes(): express.Application {
        this.app.route("/user-api/api-key")
            .post(userApiController.generateNewApiKey)

        return this.app;
    }
}