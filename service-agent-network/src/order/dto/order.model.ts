import {UserIdModel} from "../../common/dto/userId.model";
import {AnOrderIdModel} from "./anOrderId.model";
import {SpOrderParcelIdModel} from "./spOrderParcelId.model";
import {RefNoModel} from "./refNo.model";
import {CreatedAtModel} from "../../common/dto/createdAt.model";
import {SpOrderParcelModel} from "./spOrderParcel.model";
import {SpOrderParcelShippopModel} from "./spOrderParcelShippop.model";
import {SpOrderParcelShippopFlashModel} from "./spOrderParcelShippopFlash.model";
import {AnOrderPriceModel} from "./anOrderPrice.model";
import {BankAccountModel} from "./bankAccount.model";
import {AnOrderFulfillmentStatusEnum} from "./anOrderFulfillmentStatus.enum";
import {AnOrderItemModel} from "./anOrderItem.model";
import {AnCourierCodeEnum} from "./anCourierCode.enum";

export interface OrderDto
    extends UserIdModel,
        AnOrderIdModel,
        SpOrderParcelIdModel,
        RefNoModel,
        CreatedAtModel {
    desName: string;
    desPhoneNumber: string;
    desAddress: string;
    desSubdistrict: string;
    desDistrict: string;
    desProvince: string;
    desPostcode: string;
    courierCode: AnCourierCodeEnum;
    codAmount: number;
    trackingCodes: Array<string>;
    fulfillmentStatus: AnOrderFulfillmentStatusEnum;
    fulfillmentStatusString?: string;
    shippingStatus?: string;
    codStatus?: string;
    statusCompletedDate?: Date;
    codTransferredDate?: Date;
    jnaCodTransferredDate?: Date;
    itemQuantitySum?: number;
    items?: Array<AnOrderItemModel>;
    spOrderParcel?: SpOrderParcelModel;
    spOrderParcelShippop?: SpOrderParcelShippopModel;
    spOrderParcelShippopFlash?: SpOrderParcelShippopFlashModel;
    anOrderPrice?: AnOrderPriceModel;
    bankAccount?: BankAccountModel;
}