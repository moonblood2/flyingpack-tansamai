import {CommonRoutesConfig} from "../common/common.routes.config";
import apiKeyMiddleware from "../common/middleware/apikey.middleware";
import authenMiddleware from "../common/middleware/authen.middleware";
import orderMiddleware from "./middleware/order.middleware";
import orderController from "./controller/order.controller";

import express from "express";

export class OrderRoute extends CommonRoutesConfig {
    constructor(app: express.Application) {
        super(app, 'OrderRoute');
    }

    configureRoutes(): express.Application {
        this.app.route("/open-api/order")
            .post([
                apiKeyMiddleware.verifyApiKey,
                orderMiddleware.validateOrder,
                orderController.createOrder,
            ])
        this.app.route("/open-api/cancel_order")
            .post([
                apiKeyMiddleware.verifyApiKey,
                orderMiddleware.validateCancelOrder,
                orderController.cancelOrderByReferenceNo,
            ])

        this.app.route("/closed-api/order")
            .get(
                authenMiddleware.verifyToken,
                orderController.listOrders,
            )
            .put(
                authenMiddleware.verifyToken,
                orderController.updateOrders,
            )

        this.app.route("/closed-api/order/summary")
            .get(
                authenMiddleware.verifyToken,
                orderController.getOrderSummary,
            )

        this.app.route("/closed-api/order/reference-no")
            .get(
                authenMiddleware.verifyToken,
                orderController.getOrderByReferenceNo,
            )

        this.app.route("/closed-api/order/price/user-id/:userId")
            .get(
                authenMiddleware.verifyToken,
                orderController.listOrdersWithPrice,
            )

        this.app.route("/closed-api/order/cancel/an-order-id/:anOrderId")
            .put(
                authenMiddleware.verifyToken,
                orderController.cancelOrder,
            )

        this.app.route("/closed-api/order/cod/user-id/:userId")
            .get(
                authenMiddleware.verifyToken,
                orderController.getOrderCod,
            )

        return this.app;
    }
}