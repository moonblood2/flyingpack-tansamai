import {CreatedAtModel} from "../../common/dto/createdAt.model";

export interface SpOrderParcelModel extends CreatedAtModel{
    providerCode: number;
    price: number;
    paymentMethod: number;
    trackingCode: string;
    status: string;
    codAmount: number;
    weight: number;
    width: number;
    length: number;
    height: number;
}