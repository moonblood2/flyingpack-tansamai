import { CommonRoutesConfig } from "../common/common.routes.config";
import authenMiddleware from "../common/middleware/authen.middleware";
import ProductController from "./controller/product.controller";

import express from "express";

export class ProductRoute extends CommonRoutesConfig {
    constructor(app: express.Application) {
        super(app, 'ProductRoute');
    }

    configureRoutes(): express.Application {
        this.app.route("/closed-api/product")
            .get(
                authenMiddleware.verifyToken,
                ProductController.listProducts,
            )

        this.app.route("/closed-api/product/user-id/:userId")
            .get(
                authenMiddleware.verifyToken,
                ProductController.listProductsByUserId,
            )
        return this.app;
    }
}