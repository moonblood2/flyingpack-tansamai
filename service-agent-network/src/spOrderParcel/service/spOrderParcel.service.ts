import { OrderDto } from "../../order/dto/order.model";
import {
  AnOrderFulfillmentStatusEnum,
  AnOrderFulfillmentStatusToString,
} from "../../order/dto/anOrderFulfillmentStatus.enum";
import { anCourierCodeToSpCourierCode } from "../../order/dto/anCourierCode.enum";

import { makeOrder } from "../../utils/api/shipping.api";

import { NewCommonError } from "../../common/common.error";
import code from "../../common/common.code";

import debug from "debug";

const log: debug.IDebugger = debug("app:sp-order-parcel-service");

class SpOrderParcelService {
  private static instance: SpOrderParcelService;

  private constructor() {}

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

    let results: Array<any> = [];
    let resultsResp: Array<any> = [];
    for (let o of orders) {
      console.log("newAnBankAccount before");
      console.log(o);
      // remove items property unused
      const newItems = o.items?.map(
        ({
          anOrderProductId,
          anProductId,
          name,
          imgUrl,
          serialRegex,
          roboticSKU,
          serialNumbers,
          ...rest
        }) => {
          return rest;
        }
      );
      // remove items property unused

      const newAnBankAccount = {
        fiCode: o.bankAccount?.fiCode,
        bank: o.bankAccount?.bank,
        accountName: o.bankAccount?.accountName,
        accountNo: o.bankAccount?.accountNo,
        email: o.bankAccount?.email,
      };
      //prepare orderBody
      const order = {
        referenceNo: o.referenceNo,
        status: "shipping",
        trackingCode: o.trackingCode,
        weight: 1,
        width: 1,
        length: 1,
        height: 1,
        dimension: 1,
        items: newItems,
        bankAccount: newAnBankAccount,
      };
      //Request.
      const res: any = await makeOrder({data:order});
      resultsResp.push(res);
    }

    // console.log(`test :${res.data.status} : ${res.data.message}`);
    try {
      //Map response (spOrderParcelId, trackingCode, flashSortingCode) to orders.

      for (let i = 0; i < orders.length; i++) {
        let resData = resultsResp[i].data;
        if (resData.message === "success") {
          results.push(resData);
          orders[i].fulfillmentStatus = AnOrderFulfillmentStatusEnum.PACKED;
          orders[i].fulfillmentStatusString =
            AnOrderFulfillmentStatusToString[
              AnOrderFulfillmentStatusEnum.PACKED
            ];
        }
      }
      // for (let i = 0; i < orders.length; i++) {
      //   results.push(resData);
      //   if (anParcels[i]["status"]) {
      //     orders[i].trackingCodes = [anParcels[i]["tracking_code"]];
      //     orders[i].spOrderParcelId = anParcels[i]["order_parcel_id"];
      //     orders[i].fulfillmentStatus = AnOrderFulfillmentStatusEnum.PACKED;
      //     orders[i].fulfillmentStatusString =
      //       AnOrderFulfillmentStatusToString[
      //         AnOrderFulfillmentStatusEnum.PACKED
      //       ];
      //     if (anParcels[i]["shippop_flash_sorting_code"]) {
      //       let spFlash = anParcels[i]["shippop_flash_sorting_code"];
      //       orders[i].spOrderParcelShippopFlash = {
      //         sortCode: spFlash["sort_code"],
      //         dstCode: spFlash["dst_code"],
      //         sortingLineCode: spFlash["sorting_line_code"],
      //       };
      //     }
      //   }
      // }
    } catch (error) {
      log(error);
      err = NewCommonError(code.ERR_INTERNAL);
    }
    //Return orders.
    return { orders, results, err };
  }
}

export default SpOrderParcelService.getInstance();
