import orderService from "../service/order.service";
import productService from "../../product/service/product.service";
import {DateParam, parseDateParam} from "../param/date.param";
import {ListOrderParam, parseListOrderParam} from "../param/list.order.param";
import {AnOrderFulfillmentStatusEnum} from "../dto/anOrderFulfillmentStatus.enum";
import hookService from "../../hook/service/hook.service";
import userApiService from "../../userApi/service/userApi.service";
import {anOrderStatusMap} from "../dto/anOrderStatusMap";

import header from "../../common/common.header";
import code, {HttpStatusCode} from "../../common/common.code";
import {NewCommonError} from "../../common/common.error";

import express from "express";

import debug from "debug";

const log: debug.IDebugger = debug("app:open-api-order-controller");

class OrderController {
    private static instance: OrderController;

    static getInstance(): OrderController {
        if (!OrderController.instance) {
            OrderController.instance = new OrderController();
        }
        return OrderController.instance;
    }

    async createOrder(req: express.Request, res: express.Response) {
        /**
         * Aggregate items, item that have the same ProductCode, sum it up.
         * **/
        
        let reqItemsTmp = req.body.items;
        let mapProductCodeQuantity: Map<string, number> = new Map<string, number>();
        for (const i of reqItemsTmp) {
            if (!mapProductCodeQuantity.has(i.productCode)) {
                mapProductCodeQuantity.set(i.productCode, i.quantity);
            } else {
                mapProductCodeQuantity.set(i.productCode, mapProductCodeQuantity.get(i.productCode) + i.quantity);
            }
        }
        let reqItems: Array<any> = [];
        for (const [k, v] of mapProductCodeQuantity) {
            reqItems.push({
                productCode: k,
                quantity: v,
            });
        }
        let reqProductCodes: Array<string> = [];
        for (const e of reqItems) {
            reqProductCodes.push(e.productCode);
        }
        /**Get products by productCode and userId*/
        let {products, err} = await productService.getProductsByProductCodes(reqProductCodes, req.body.userId);
        /**mapProductCodeAnproductId key=productCode, value=anProductId*/
        let mapProductCodeAnproductId: Map<string, string> = new Map<string, string>()
        for (const p of products) {
            mapProductCodeAnproductId.set(<string>p.productCode, <string>p.anProductId);
        }
        if (err.code === code.SUCCESS) {
            /**existingProductCodes is a list of productCode that exist in database.*/
            let existingProductCodes: Array<string> = [];
            for (const p of products) {
                existingProductCodes.push(<string>p.productCode);
            }
            /**Keep exist status of each requested productCode in results. key=productCode, value=exist status*/
            let results: Map<string, boolean> = new Map<string, boolean>();
            let exist = true;
            for (let i = 0; i < reqItems.length; i++) {
                const hasProductCode = existingProductCodes.includes(reqProductCodes[i]);
                results.set(<string>reqItems[i].productCode, hasProductCode);
                exist = exist && hasProductCode;
                if (hasProductCode) {
                    /**Add product id to requested items.*/
                    reqItems[i].anProductId = mapProductCodeAnproductId.get(reqItems[i].productCode);
                }
            }
            /**Check if all requested productCodes do exist.*/
            if (exist) {
                req.body.sortCode = req.body.trackingDetail.sortCode
                req.body.lineCode=req.body.trackingDetail.lineCode
                req.body.sortingLineCode=req.body.trackingDetail.sortingLineCode
                req.body.dstStoreName=req.body.trackingDetail.dstStoreName
                log("req.body: %O", req.body)
                log("reqItems: %O", reqItems)
                let {order, err} = await orderService.create(req.body, reqItems);
                if (err.code === code.SUCCESS) {
                    res.status(HttpStatusCode[err.code]).send({
                        ...err,
                        data: {
                            referenceNo: order.referenceNo,
                        },
                    });
                } else {
                    res.status(HttpStatusCode[<number>err.code]).send({...err});
                }
            } else {
                /**Make error message of "do not exist". Select only an item that doesn't exist productCode.*/
                let key = [];
                for (let [k, v] of results) {
                    if (!v) {
                        key.push(k);
                    }
                }
                err = NewCommonError(code.INPUT_PRODUCT_CODE_DOES_NOT_EXIST, `productCode ${key.join(", ")} do not exist`);
                res.status(HttpStatusCode[<number>err.code]).send({...err});
            }
        } else {
            res.status(HttpStatusCode[<number>err.code]).send({...err});
        }
    }

    async cancelOrderByReferenceNo(req: express.Request, res: express.Response) {
        let {anOrderId, err} = await orderService.getAnOrderIdByReferenceNo(req.body.referenceNo)
        if (err.code === code.SUCCESS) {
            let {err} = await orderService.updateOrderFulfillmentStatus(<string>anOrderId, AnOrderFulfillmentStatusEnum.CANCEL);
            if (err.code === code.SUCCESS) {
                res.status(HttpStatusCode[err.code]).send({
                    ...err,
                });
            } else {
                res.status(HttpStatusCode[<number>err.code]).send({
                    ...err,
                });
            }
        } else {
            res.status(HttpStatusCode[<number>err.code]).send({
                ...err,
            });
        }
    }

    async listOrders(req: express.Request, res: express.Response) {
        const userId = <string>req.headers[header.USER_ID];
        const token = <string>req.headers[header.TOKEN];
        const params: ListOrderParam = parseListOrderParam(req);

        const {data, err} = await orderService.list(params, true, userId, token);
        if (err.code === code.SUCCESS) {
            res.status(HttpStatusCode[err.code]).send({
                ...err,
                data: data,
            });
        } else {
            res.status(HttpStatusCode[<number>err.code]).send({
                ...err,
            });
        }
    }

    //getOrderSummary get today order summary report.
    async getOrderSummary(req: express.Request, res: express.Response) {
        const userId = <string>req.headers["user-id"];
        const params: DateParam = parseDateParam(req);

        const {orderSum, err} = await orderService.getOrderSummary(userId, params);
        if (err.code === code.SUCCESS) {
            const {orderProductSum, err} = await orderService.getOrderProductSummary(userId, params);
            if (err.code === code.SUCCESS) {
                res.status(HttpStatusCode[err.code]).send({
                    ...err,
                    data: {
                        orderSummary: orderSum,
                        orderProductSummary: orderProductSum,
                    },
                });
            } else {
                res.status(HttpStatusCode[<number>err.code]).send({...err});
            }
        } else {
            res.status(HttpStatusCode[<number>err.code]).send({...err});
        }
    }

    async getOrderByReferenceNo(req: express.Request, res: express.Response) {
        const key = <string>req.query["referenceNo"];

        const {order, err} = await orderService.getOrderByKey(key);
        if (err.code === code.SUCCESS) {
            res.status(HttpStatusCode[err.code]).send({
                ...err,
                data: {
                    order: order,
                },
            });
        } else {
            res.status(HttpStatusCode[<number>err.code]).send({...err});
        }
    }

    async listOrdersWithPrice(req: express.Request, res: express.Response) {
        const userId = req.params.userId;
        const token = <string>req.headers[header.TOKEN];
        const params: ListOrderParam = parseListOrderParam(req);

        //Set no limit, no have pagination when query.
        params.perPage = 99999;
        params.page = 1;

        let {data, err} = await orderService.list(params, false, userId, token);
        if (err.code === code.SUCCESS) {
            //Calculate price
            const {orders, totalFulfillmentServiceCharge} = orderService.calculateOrderPrice(data.orders);
            data = {
                ...data,
                orders: orders,
                totalFulfillmentServiceCharge: totalFulfillmentServiceCharge,
            };

            res.status(HttpStatusCode[err.code]).send({
                ...err,
                data: data,
            });
        } else {
            res.status(HttpStatusCode[<number>err.code]).send({
                ...err,
            });
        }
    }

    async cancelOrder(req: express.Request, res: express.Response) {
        let {err} = await orderService.updateOrderFulfillmentStatus(req.params.anOrderId, AnOrderFulfillmentStatusEnum.CANCEL);
        if (err.code === code.SUCCESS) {
            res.status(HttpStatusCode[err.code]).send({
                ...err,
            });
        } else {
            res.status(HttpStatusCode[<number>err.code]).send({
                ...err,
            });
        }
    }

    async updateOrders(req: express.Request, res: express.Response) {
        let userId = <string>req.headers["user-id"];
        /**
         * Update database
         * */
        let reqOrders = req.body.orders;
        let agentUserId = req.body.agentUserId;

        if (agentUserId) {
            userId = agentUserId;
        }

        let {orders, err} = await orderService.updateOrders(reqOrders);
        if (err.code === code.SUCCESS) {
            /**
             * Hook updated data to customer.
             * */
            let {url, err} = await userApiService.getHookUrlOrder(userId);
            if (err.code === code.SUCCESS) {
                for (let i = 0; i < reqOrders.length; i++) {
                    let item = reqOrders[i];
                    let {
                        referenceNo,
                        err
                    } = await orderService.getReferenceNoByAnOrderId(<string>item.anOrderId);
                    if (err.code === code.SUCCESS) {
                        let trackingCode = undefined;
                        let status = undefined;
                        let transferredDate = undefined;

                        if (item.trackingCode) {
                            trackingCode = item.trackingCode;
                        }
                        if (item.jnaCodTransferredDate) {
                            status = anOrderStatusMap.TRANSFERRED;
                            transferredDate = item.jnaCodTransferredDate;
                        } else if (item.fulfillmentStatus == AnOrderFulfillmentStatusEnum.PACKED) {
                            status = anOrderStatusMap.FULFILLMENT;
                        }

                        hookService.sendOrder({
                            referenceNo: referenceNo,
                            trackingCode: trackingCode,
                            status: status,
                            dateCod: transferredDate,
                        }, url);
                    }
                }
            }
            res.status(HttpStatusCode[code.SUCCESS]).send({
                ...err,
                data: orders,
            });
        } else {
            res.status(HttpStatusCode[<number>err.code]).send({
                ...err,
            });
        }
    }

    async getOrderCod(req: express.Request, res: express.Response) {
        const userId = <string>req.params.userId;
        const params = parseListOrderParam(req);
        params.page = 1;
        params.perPage = 999999;
        params.isCod = true;
        const {orders, err} = await orderService.getOrderCod(params, userId);
        if (err.code === code.SUCCESS) {
            res.status(HttpStatusCode[err.code]).send({
                ...err,
                data: {
                    orders: orders,
                },
            });
        } else {
            res.status(HttpStatusCode[<number>err.code]).send({...err});
        }
    }
}

export default OrderController.getInstance();