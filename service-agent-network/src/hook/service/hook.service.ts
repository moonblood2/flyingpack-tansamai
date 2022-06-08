import {AnOrderItemModel} from "../../order/dto/anOrderItem.model";

import userApiService from "../../userApi/service/userApi.service";
import orderService from "../../order/service/order.service";

import {OrderPayload} from "../payload/order.payload";
import {OrderStatusPayload} from "../payload/orderStatus.payload";
import {ItemPayload} from "../payload/item.payload";

import code from "../../common/common.code";

import axios from "axios";

import debug from "debug";

const log: debug.IDebugger = debug("app:hook-service");

class HookService {
    private static instance: HookService;

    private constructor() {
    }

    static getInstance(): HookService {
        if (!HookService.instance) {
            HookService.instance = new HookService();
        }
        return HookService.instance;
    }

    async sendOrder(resource: OrderPayload, url: string) {
        try {
            await axios({
                method: "PUT",
                url: url,
                headers: {"api-key": "HTTP_API_KEY", "Content-Type": "application/json"},
                data: {data: resource}
            })
        } catch (e) {
            log(e);
        }
    }

    async sendOrderStatus(resource: OrderStatusPayload, url: string) {
        try {
            await axios({
                method: "PUT",
                url: url,
                headers: {"api-key": "HTTP_API_KEY", "Content-Type": "application/json"},
                data: {data: resource}
            })
        } catch (e) {
            log(e);
        }
    }

    async sendSerialNumber(userId: string, anOrderId: string) {
        let {url, err} = await userApiService.getHookUrlOrder(userId);
        if (err.code === code.SUCCESS) {
            let {err, referenceNo} = await orderService.getReferenceNoByAnOrderId(anOrderId);
            if (err.code === code.SUCCESS) {
                let {err, order} = await orderService.getOrderByAnOrderId(anOrderId);
                if (err.code === code.SUCCESS) {
                    let items: Array<ItemPayload> = [];
                    for (const i of <Array<AnOrderItemModel>>order.items) {
                        let serialNumbers: Array<string> = [];
                        let rangedSerialNumbers: Array<Array<string>> = [];
                        for (const s of i.serialNumbers) {
                            if (!s.serialNumberEnd) {
                                serialNumbers.push(s.serialNumberStart);
                            } else {
                                rangedSerialNumbers.push([s.serialNumberStart, s.serialNumberEnd]);
                            }
                        }
                        if (serialNumbers.length > 0) {
                            items.push({
                                productCode: i.productCode,
                                serialRange: false,
                                serialNumbers: serialNumbers,
                            });
                        }
                        if (rangedSerialNumbers.length > 0) {
                            items.push({
                                productCode: i.productCode,
                                serialRange: true,
                                serialNumbers: rangedSerialNumbers,
                            });
                        }
                    }
                    await this.sendOrder({
                        referenceNo: referenceNo,
                        trackingCode: undefined,
                        items: items ? items : undefined,
                    }, url);
                } else {
                    return false;
                }
            } else {
                return false;
            }
        } else {
            return false;
        }
        return true
    }
}

export default HookService.getInstance();