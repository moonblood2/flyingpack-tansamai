import {RefNoModel} from "../../order/dto/refNo.model";
import {OrderStatusModel} from "../../common/dto/orderStatus.model";

export interface OrderStatusPayload extends RefNoModel, OrderStatusModel {
}