import {UserIdModel} from "../../common/dto/userId.model";
import {SpOrderParcelIdModel} from "../../order/dto/spOrderParcelId.model";
import {AnOrderIdModel} from "../../order/dto/anOrderId.model";
import {OrderStatusModel} from "../../common/dto/orderStatus.model";
import {OrderParcelShapeModel} from "../../common/dto/orderParcelShape.model";
import {AnOrderItem} from "../../common/dto/anOrderItem";

export interface SpOrderWeightModel
    extends
        UserIdModel,
        SpOrderParcelIdModel,
        AnOrderIdModel,
        OrderStatusModel,
        OrderParcelShapeModel {
    trackingCode?: string;
    anOrderItems: Array<AnOrderItem>;
}