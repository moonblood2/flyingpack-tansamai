import {AnOrderIdModel} from "./anOrderId.model";
import {CreatedAtModel} from "../../common/dto/createdAt.model";

export interface OrderProductModel extends AnOrderIdModel, CreatedAtModel {
    anOrderProductId?: string;
    anProductId?: string;
    quantity?: number;
}