import {OrderDto} from "../dto/order.model";
import orderDao from "../dao/order.dao";

import {OrderProductModel} from "../dto/orderProduct.model";
import {SpOrderParcelModel} from "../dto/spOrderParcel.model";
import {AnOrderPriceModel} from "../dto/anOrderPrice.model";
import {AnOrderFulfillmentStatusEnum} from "../dto/anOrderFulfillmentStatus.enum";

import {ListOrderParam} from "../param/list.order.param";
import {DateParam} from "../param/date.param";

import {PaginationResource} from "../../common/resource/pagination.resource";

import {getOrderParcelByIds} from "../../utils/api/shipping.api";

import code from "../../common/common.code";
import {NewCommonError} from "../../common/common.error";

import debug, {IDebugger} from "debug";

const log: IDebugger = debug("app:order-service");

class OrderService {
    private static instance: OrderService;

    static getInstance(): OrderService {
        if (!OrderService.instance) {
            OrderService.instance = new OrderService();
        }
        return OrderService.instance;
    }

    async create(order: OrderDto, items: any) {
        let orderProducts: Array<OrderProductModel> = [];
        for (const e of items) {
            orderProducts.push({
                anProductId: e.anProductId,
                quantity: e.quantity,
            });
        }
        return await orderDao.addOrder(order, orderProducts);
    }

    async list(params: ListOrderParam, isSorting: boolean, userId: string, token: string) {
        let data: any = {};
        let {orders, totalItem, err} = await orderDao.getOrders(params, isSorting, userId);

        if (err.code !== code.SUCCESS) {
            return {err: err};
        }

        //Aggregate sp_order_parcel_id to request order data.
        let spOrderParcelIds: Array<string> = [];
        for (const o of orders) {
            if (o.spOrderParcelId) {
                spOrderParcelIds.push(o.spOrderParcelId);
            }
        }

        if (spOrderParcelIds.length > 0) {
            try {
                const res: any = await getOrderParcelByIds(token, userId, spOrderParcelIds);
                if (res.data) {
                    for (let i = 0; i < orders.length; i++) {
                        if (orders[i].spOrderParcelId && res.data[orders[i].spOrderParcelId]) {
                            const spOrder = res.data[orders[i].spOrderParcelId];
                            orders[i].spOrderParcel = {
                                providerCode: spOrder["order_parcel"]["provider_code"],
                                price: spOrder["order_parcel"]["price"],
                                paymentMethod: spOrder["order_parcel"]["payment_method"],
                                trackingCode: spOrder["order_parcel"]["tracking_code"],
                                status: spOrder["order_parcel"]["status"],
                                codAmount: spOrder["order_parcel"]["cod_amount"],
                                weight: spOrder["order_parcel"]["weight"],
                                width: spOrder["order_parcel"]["width"],
                                length: spOrder["order_parcel"]["length"],
                                height: spOrder["order_parcel"]["height"],
                            }
                            if (spOrder["order_parcel_shippop"]) {
                                orders[i].spOrderParcelShippop = {
                                    id: spOrder["order_parcel_shippop"]["id"],
                                    purchaseId: spOrder["order_parcel_shippop"]["purchase_id"],
                                    status: spOrder["order_parcel_shippop"]["status"],
                                    courierCode: spOrder["order_parcel_shippop"]["courier_code"],
                                    courierTrackingCode: spOrder["order_parcel_shippop"]["courier_tracking_code"],
                                    trackingCode: spOrder["order_parcel_shippop"]["tracking_code"],
                                    codAmount: spOrder["order_parcel_shippop"]["cod_amount"],
                                }
                            }
                            if (spOrder["order_parcel_shippop_flash"]) {
                                orders[i].spOrderParcelShippopFlash = {
                                    sortCode: spOrder["order_parcel_shippop_flash"]["sort_code"],
                                    dstCode: spOrder["order_parcel_shippop_flash"]["dst_code"],
                                    sortingLineCode: spOrder["order_parcel_shippop_flash"]["sorting_line_code"],
                                }
                            }
                        }
                    }
                } else {
                    log(res.data)
                    return {err: NewCommonError(code.ERR_INTERNAL)}
                }

            } catch (error) {
                log(error);
                return {err: NewCommonError(code.ERR_INTERNAL)}
            }
        }
        let currentPage = params.page;
        let totalPage = Math.ceil(totalItem / params.perPage);
        const pagination: PaginationResource = {
            previousPage: currentPage - 1 > 0 ? currentPage - 1 : null,
            currentPage: currentPage,
            nextPage: currentPage + 1 <= totalPage ? currentPage + 1 : null,
            firstPage: 1,
            lastPage: totalPage,
            isFirstPage: currentPage === 1,
            isLastPage: currentPage === totalPage,
            totalPage: totalPage,
            totalItem: totalItem,
        }
        data = {
            ...pagination,
            orders,
        }
        return {data, err};
    }

    async getOrderBySpOrderParcelId(spOrderParcelId: string) {
        return await orderDao.getOrderBySpOrderParcelId(spOrderParcelId);
    }

    async getOrderByAnOrderId(anOrderId: string) {
        return await orderDao.getOrderByAnOrderId(anOrderId);
    }

    async getAnOrderIdByReferenceNo(referenceNo: string) {
        return await orderDao.getAnOrderIdByReferenceNo(referenceNo);
    }

    async getOrderByKey(key: string) {
        return await orderDao.getOrderByKey(key);
    }

    async updateOrders(orders: Array<OrderDto>) {
        return await orderDao.updateOrders(orders);
    }

    async updateSpOrderParcelIdAndTrackingCode(anOrderId: string, spOrderParcelId: string, trackingCode: string) {
        return await orderDao.updateSpOrderParcelIdAndTrackingCode(anOrderId, spOrderParcelId, trackingCode);
    }

    async getReferenceNoByAnOrderId(anOrderId: string) {
        return await orderDao.getReferenceNoByAnOrderId(anOrderId);
    }

    async getReferenceNoBySpOrderParcelId(spOrderParcelId: string) {
        return await orderDao.getReferenceNoBySpOrderParcelId(spOrderParcelId);
    }

    async getOrderProductsByAnOrderId(anOrderId: string, userId: string) {
        return await orderDao.getOrderProductsByAnOrderId(anOrderId, userId);
    }

    async getOrderSummary(userId: string, params: DateParam) {
        return await orderDao.getOrderSummary(userId, params);
    }

    async getOrderProductSummary(userId: string, params: DateParam) {
        return await orderDao.getOrderProductSummary(userId, params);
    }

    calculateOrderPrice(orders: Array<OrderDto>) {
        let totalFulfillmentServiceCharge: number = 0;
        for (let i = 0; i < orders.length; i++) {
            orders[i].anOrderPrice = {
                fulfillmentServiceCharge: 0,
                parcelCost: 0,
                parcelRemoteAreaPrice: 0,
                parcelSellingPrice: 0
            }

            let {itemQuantitySum, spOrderParcel} = orders[i];
            itemQuantitySum = <number> itemQuantitySum;
            spOrderParcel = <SpOrderParcelModel> spOrderParcel;

            let fulfillmentServiceCharge = 0;
            if (itemQuantitySum === 1) {
                fulfillmentServiceCharge = 3;
            } else if (itemQuantitySum >= 2 && itemQuantitySum <= 3) {
                fulfillmentServiceCharge = 5;
            } else if (itemQuantitySum >= 4 && itemQuantitySum <= 5) {
                fulfillmentServiceCharge = 10;
            } else if (itemQuantitySum >= 6) {
                fulfillmentServiceCharge = 13;
            }

            (<AnOrderPriceModel> orders[i].anOrderPrice).fulfillmentServiceCharge = fulfillmentServiceCharge;
            totalFulfillmentServiceCharge += fulfillmentServiceCharge;

            if (orders[i].spOrderParcelId && orders[i].spOrderParcel) {
                (<AnOrderPriceModel> orders[i].anOrderPrice).parcelCost = spOrderParcel.price;
                (<AnOrderPriceModel> orders[i].anOrderPrice).parcelRemoteAreaPrice = 0;
                (<AnOrderPriceModel> orders[i].anOrderPrice).parcelSellingPrice = 0;
            }
        }
        return {orders, totalFulfillmentServiceCharge};
    }

    async updateOrderFulfillmentStatus(anOrderId: string, fulfillmentStatus: AnOrderFulfillmentStatusEnum) {
        return await orderDao.updateOrderFulfillmentStatus(anOrderId, fulfillmentStatus);
    }

    async getOrderCod(params: ListOrderParam, userId: string) {
        return await orderDao.getOrders(params, false, userId);
    }
}

export default OrderService.getInstance();