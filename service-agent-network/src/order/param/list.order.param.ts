import express from "express";
import {AnOrderFulfillmentStatusEnum, parseAnOrderFulfillmentStatusEnum} from "../dto/anOrderFulfillmentStatus.enum";
import {DateTypeEnum, parseDateTypeEnum} from "./dateType.param";

export interface ListOrderParam {
    startDate: string;
    endDate: string;
    dateType: DateTypeEnum;
    timeZone?: string;
    page: number;
    perPage: number;
    keyWord: string;
    isCod: boolean;
    fulfillmentStatus?: AnOrderFulfillmentStatusEnum;
    courierCode: number;
    productIds: Array<string>;
    quantity: number;
}

export const parseListOrderParam = (req: express.Request): ListOrderParam => {
    let startDate = <string> req.query.startDate;
    let endDate = <string> req.query.endDate;
    let dateType = parseDateTypeEnum(<string> req.query.dateType);
    let timeZone = "Asia/Bangkok";
    let page = Number(<string> req.query.page) || 1;
    let perPage = Number(<string> req.query.perPage) || 999999;
    let keyWord = <string> req.query.keyWord || "";
    let fulfillmentStatus = parseAnOrderFulfillmentStatusEnum(<string> req.query.fulfillmentStatus);
    let courierCode = Number(<string> req.query.courierCode) || 0;
    let productIdsReq = req.query.productIds ? req.query.productIds: [];
    let productIds: Array<string> = [];
    if (productIdsReq instanceof Array) {
        productIds = <Array<string>> productIdsReq;
    } else {
        productIds = [];
    }
    let quantity = Number(<string> req.query.quantity)

    return {
        startDate: startDate,
        endDate: endDate,
        dateType: dateType,
        timeZone: timeZone,
        page: page,
        perPage: perPage,
        keyWord: keyWord,
        isCod: false,
        fulfillmentStatus: fulfillmentStatus,
        courierCode: courierCode,
        productIds: productIds,
        quantity: quantity,
    };
}