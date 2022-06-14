import {RefNoModel} from "../../order/dto/refNo.model";
import {OrderStatusModel} from "../../common/dto/orderStatus.model";
import {OrderParcelShapeModel} from "../../common/dto/orderParcelShape.model";
import {OrderSerialNo} from "../../common/dto/orderSerialNo";
import {ItemPayload} from "./item.payload";

export interface OrderPayload
    extends RefNoModel,
        OrderStatusModel,
        OrderSerialNo,
        OrderParcelShapeModel {
    dimension?: number;
    items?: Array<ItemPayload>;
    trackingCode?:string;
    dateCod?: string;
}