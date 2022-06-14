import {AnOrderIdModel} from "./anOrderId.model";
import {CreatedAtModel} from "../../common/dto/createdAt.model";
import { AnTrackingDetailModel } from "./anTrackingDetail.model";

export interface OrderProductModel extends AnOrderIdModel, CreatedAtModel,AnTrackingDetailModel{
    anOrderProductId?: string;
    anProductId?: string;
    anTrackingDetailId?:string;
    quantity?: number;
}