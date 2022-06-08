import {OrderStatusModel} from "../../common/dto/orderStatus.model";
import {OrderCodStatusModel} from "../../common/dto/orderCodStatus.model";
import {SpOrderParcelIdModel} from "../../order/dto/spOrderParcelId.model";

export interface SpOrderStatusModel extends
    SpOrderParcelIdModel,
    OrderStatusModel,
    OrderCodStatusModel {
    trackingCode?: string;
    statusCompletedDate?: Date;
    codTransferredDate?: Date;
}