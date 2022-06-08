import {SpOrderWeightModel} from "../dto/spOrderWeight.model";
import {SpOrderStatusModel} from "../dto/spOrderStatusModel";
import {SerialNumberModel} from "../../serialNumber/dto/serialNumber.model";
import {anOrderStatusMap} from "../../order/dto/anOrderStatusMap";

import spOrderParcelService from "../service/spOrderParcel.service";
import orderService from "../../order/service/order.service";
import userApiService from "../../userApi/service/userApi.service";
import serialNumberService from "../../serialNumber/service/serialNumber.service";

import hookService from "../../hook/service/hook.service";

import code, {HttpStatusCode} from "../../common/common.code";
import header from "../../common/common.header";

import express from "express";

import debug from "debug";

const log: debug.IDebugger = debug("app:sp-order-parcel-controller");

class SpOrderParcelController {
    private static instance: SpOrderParcelController;

    private constructor() {
    }

    static getInstance(): SpOrderParcelController {
        if (!SpOrderParcelController.instance) {
            SpOrderParcelController.instance = new SpOrderParcelController();
        }
        return SpOrderParcelController.instance;
    }

    async makeOrder(req: express.Request, res: express.Response) {
        const token = <string>req.headers[header.TOKEN];
        const userId = <string>req.headers[header.USER_ID];

        const {orders, results, err} = await spOrderParcelService.makeOrder(token, req.body.origin, req.body.orders);

        if (err.code === code.SUCCESS) {
            const {err} = await orderService.updateOrders(orders);

            if (err.code === code.SUCCESS) {
                const {url, err} = await userApiService.getHookUrlOrder(userId);

                if (err.code === code.SUCCESS) {
                    for (const o of orders) {
                        hookService.sendOrder({
                            referenceNo: o.referenceNo,
                            status: anOrderStatusMap.FULFILLMENT,
                            trackingCode: o.trackingCodes,
                        }, url);
                    }

                    res.status(200).send({
                        ...err,
                        data: {
                            orders: orders,
                            results: results,
                        },
                    });
                } else {
                    res.status(HttpStatusCode[<number>err.code]).send({...err});
                }
            } else {
                res.status(HttpStatusCode[<number>err.code]).send({...err});
            }
        } else {
            res.status(HttpStatusCode[<number>err.code]).send({...err});
        }
    }

    /**
     * First time when FlyingPack callback (hook) after order has packed and sent to provider.
     * It contains status, weight, width, length and height.
     */
    async receiveOrderWeight(req: express.Request, res: express.Response) {
        const {body} = req;
        let anOrderItems: Array<any> = [];
        if (body["an_order_items"]) {
            for (const e of body["an_order_items"]) {
                let serialNumberCount = 0;
                let serialNumbers: Array<string> = [];
                if (e["serial_numbers"]) {
                    for (const s of e["serial_numbers"]) {
                        if (s || s !== "") {
                            serialNumbers.push(s);
                            serialNumberCount++;
                        }
                    }
                }
                anOrderItems.push({
                    productCode: e["product_code"],
                    serialNumbers: serialNumbers,
                });
            }
        }
        const resource: SpOrderWeightModel = {
            userId: body["user_id"],
            spOrderParcelId: body["sp_order_parcel_id"],
            anOrderId: body["an_order_id"],
            anOrderItems: anOrderItems,
            status: body["status"],
            weight: body["weight"],
            width: body["width"],
            length: body["length"],
            height: body["height"],
        }
        //Update sp_order_parcel_id. Use updateSpOrderParcelId from order service.
        const {err} = await orderService.updateSpOrderParcelIdAndTrackingCode(<string>resource.anOrderId, resource.spOrderParcelId, <string>resource.trackingCode);
        if (err.code === code.SUCCESS) {
            //Get url_webhook_order by user_id. Use getHookUrlOrder from userApi service.
            const {url, err} = await userApiService.getHookUrlOrder(<string>resource.userId);
            if (err.code === code.SUCCESS) {
                //Get reference_no by sp_order_parcel_id
                const {referenceNo, err} = await orderService.getReferenceNoByAnOrderId(<string>resource.anOrderId);
                if (err.code === code.SUCCESS) {
                    const {
                        orderProducts,
                        err
                    } = await orderService.getOrderProductsByAnOrderId(<string>resource.anOrderId, <string>resource.userId);
                    if (err.code === code.SUCCESS) {
                        //Map orderProducts to orderProductSerialNumbers, (key=productCode, value=anOrderProductId)
                        let mapAnOrderProductId: Map<string, string> = new Map<string, string>();
                        for (const o of orderProducts) {
                            if (!mapAnOrderProductId.has(o.productCode)) {
                                mapAnOrderProductId.set(o.productCode, o.anOrderProductId);
                            }
                        }
                        //Make orderProductSerialNumbers
                        let orderProductSerialNumbers: Array<SerialNumberModel> = [];
                        for (const item of resource.anOrderItems) {
                            //Check if have serialNumbers
                            if (item.serialNumbers) {
                                for (const s of item.serialNumbers) {
                                    orderProductSerialNumbers.push({
                                        anOrderProductId: <string>mapAnOrderProductId.get(<string>item.productCode),
                                        serialNumberStart: s,
                                    })
                                }
                            }
                        }
                        //Save serial number to Database
                        const {err} = await serialNumberService.updateSerialNumbers(orderProductSerialNumbers);
                        //Send order weight details to client. Use sendOrderWeight from hook service.
                        hookService.sendOrder({
                            referenceNo: referenceNo,
                            status: resource.status,
                            weight: resource.weight,
                            width: resource.width,
                            length: resource.length,
                            height: resource.height,
                            dimension: <number>resource.width + <number>resource.length + <number>resource.height,
                        }, url);
                        res.status(HttpStatusCode[<number>err.code]).send({...err});
                    } else {
                        log(err);
                        res.status(HttpStatusCode[<number>err.code]).send({...err});
                    }
                } else {
                    log(err);
                    res.status(HttpStatusCode[<number>err.code]).send({...err});
                }
            } else {
                log(err);
                res.status(HttpStatusCode[<number>err.code]).send({...err});
            }
        } else {
            log(err);
            res.status(HttpStatusCode[<number>err.code]).send({...err});
        }
    }

    /**
     * Receive order status from FlyingPack.
     */
    async receiveOrderStatus(req: express.Request, res: express.Response) {
        const {body} = req;
        const resource: SpOrderStatusModel = {
            spOrderParcelId: <string>body["sp_order_parcel_id"],
            status: <string>body["status"],
            codStatus: <string>body["cod_status"],
            statusCompletedDate: <Date>body["status_completed_date"],
            codTransferredDate: <Date>body["cod_transferred_date"],
        }

        //Get reference_no by sp_order_parcel_id
        let {order, err} = await orderService.getOrderBySpOrderParcelId(<string>resource.spOrderParcelId);
        if (err.code === code.SUCCESS) {
            order.shippingStatus = resource.status;
            order.codStatus = resource.codStatus;
            order.statusCompletedDate = resource.statusCompletedDate;
            order.codTransferredDate = resource.codTransferredDate;

            let {err} = await orderService.updateOrders([order]);
            if (err.code === code.SUCCESS) {
                //Send order weight details to client. Use sendOrderWeight from hook service.
                //Get url_webhook_order by sp_order_parcel_id. Use getHookUrlOrder from userApi service.
                //Ignore codStatus that equal 'transferred' that should wait for confirmation from JNA.
                //If an Admin update the jna_transferred_date then it will be send to a client.
                const {
                    url,
                    err
                } = await userApiService.getHookUrlOrderBySpOrderParcelId(<string>resource.spOrderParcelId);
                let status = undefined;
                if (order.shippingStatus) {
                    status = order.shippingStatus;
                } else if (order.codStatus && order.codStatus !== anOrderStatusMap.TRANSFERRED) { // Ignore transferred status
                    status = order.codStatus;
                }
                if (err.code === code.SUCCESS) {
                    hookService.sendOrderStatus({
                        referenceNo: order.referenceNo,
                        status: status,
                    }, url);
                    res.status(HttpStatusCode[<number>err.code]).send({...err});
                } else {
                    log(err);
                    res.status(HttpStatusCode[<number>err.code]).send({...err});
                }
            } else {
                log(err);
                res.status(HttpStatusCode[<number>err.code]).send({...err});
            }
        } else {
            log(err);
            res.status(HttpStatusCode[<number>err.code]).send({...err});
        }
    }
}

export default SpOrderParcelController.getInstance();