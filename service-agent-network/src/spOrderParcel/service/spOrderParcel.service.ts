import {OrderDto} from "../../order/dto/order.model";
import {
    AnOrderFulfillmentStatusEnum,
    AnOrderFulfillmentStatusToString
} from "../../order/dto/anOrderFulfillmentStatus.enum";
import {anCourierCodeToSpCourierCode} from "../../order/dto/anCourierCode.enum";

import {makeOrder} from "../../utils/api/shipping.api";

import {NewCommonError} from "../../common/common.error";
import code from "../../common/common.code";

import debug from "debug";

const log: debug.IDebugger = debug("app:sp-order-parcel-service");

class SpOrderParcelService {
    private static instance: SpOrderParcelService;

    private constructor() {
    }

    static getInstance(): SpOrderParcelService {
        if (!SpOrderParcelService.instance) {
            SpOrderParcelService.instance = new SpOrderParcelService();
        }
        return SpOrderParcelService.instance;
    }

    /**
     * makOrder
     * Build payload for request to make order with Shipping Service,
     * then populate trackingCode to input orders.
     * */
    async makeOrder(token: string, origin: any, orders: Array<OrderDto>) {
        let err = NewCommonError();
        //Map payload.
        let anOrders: Array<any> = [];
        let results: Array<any> = [];
        for (const o of orders) {
            //NOTE: At this time write hard-code provider_code=1, mean it's only use with Shippop.
            //NOTE: Write hard-code parcel_shape.
            anOrders.push({
                "provider_code": 1,
                "courier_code": anCourierCodeToSpCourierCode[o.courierCode],
                "enable_cod": o.codAmount > 0,
                "cod_amount": o.codAmount,
                "origin": {
                    "name": origin.name,
                    "address": origin.address,
                    "district": origin.district,
                    "state": origin.state,
                    "province": origin.province,
                    "postcode": origin.postcode,
                    "phone_number": origin.phone_number
                },
                "destination": {
                    "name": o.desName,
                    "address": o.desAddress,
                    "district": o.desSubdistrict,
                    "state": o.desDistrict,
                    "province": o.desProvince,
                    "postcode": o.desPostcode,
                    "phone_number": o.desPhoneNumber,
                },
                "parcel_shape": {
                    "weight": 1,
                    "width": 1,
                    "length": 1,
                    "height": 1
                }
            })
        }
        //Request.
        try {
            const res: any = await makeOrder(token, anOrders);
            //Map response (spOrderParcelId, trackingCode, flashSortingCode) to orders.
            let resData = res.data;
            if (resData['an_parcels']) {
                let anParcels = resData['an_parcels'];
                for (let i = 0; i < orders.length; i++) {
                    results.push({
                        status: anParcels[i]['status'],
                        message: anParcels[i]['message'],
                    })
                    if (anParcels[i]['status']) {
                        orders[i].trackingCodes = [anParcels[i]['tracking_code']];
                        orders[i].spOrderParcelId = anParcels[i]['order_parcel_id'];
                        orders[i].fulfillmentStatus = AnOrderFulfillmentStatusEnum.PACKED;
                        orders[i].fulfillmentStatusString = AnOrderFulfillmentStatusToString[AnOrderFulfillmentStatusEnum.PACKED];
                        if (anParcels[i]['shippop_flash_sorting_code']) {
                            let spFlash = anParcels[i]['shippop_flash_sorting_code'];
                            orders[i].spOrderParcelShippopFlash = {
                                sortCode: spFlash['sort_code'],
                                dstCode: spFlash['dst_code'],
                                sortingLineCode: spFlash['sorting_line_code'],
                            }
                        }
                    }
                }
            }
        } catch (error) {
            log(error);
            err = NewCommonError(code.ERR_INTERNAL);
        }
        //Return orders.
        return {orders, results, err};
    }
}

export default SpOrderParcelService.getInstance();