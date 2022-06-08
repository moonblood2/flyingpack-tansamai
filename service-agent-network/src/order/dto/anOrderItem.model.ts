import {SerialNumberModel} from "../../serialNumber/dto/serialNumber.model";

export interface AnOrderItemModel {
    anOrderProductId: string;
    anProductId: string;
    productCode: string;
    name: string;
    imgUrl: string;
    quantity: number;
    serialRegex: string;
    roboticSKU: string;
    serialNumbers: Array<SerialNumberModel>
}