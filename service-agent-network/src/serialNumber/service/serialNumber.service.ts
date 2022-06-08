import {SerialNumberModel} from "../dto/serialNumber.model";
import {AnOrderItemModel} from "../../order/dto/anOrderItem.model";
import serialNumberDao from "../dao/serialNumber.dao";

class SerialNumberService {
    private static instance: SerialNumberService;

    private constructor() {
    }

    static getInstance(): SerialNumberService {
        if (!SerialNumberService.instance) {
            SerialNumberService.instance = new SerialNumberService();
        }
        return SerialNumberService.instance;
    }

    async getSerialNumbersByAnOrderId(anOrderId: string) {
        return await serialNumberDao.getSerialNumbersByAnOrderId(anOrderId);
    }

    async createSerialNumbers(items: Array<AnOrderItemModel>) {
        return await serialNumberDao.createSerialNumbers(items);
    }

    async updateSerialNumbers(serialNumbers: Array<SerialNumberModel>) {
        return await serialNumberDao.updateSerialNumbers(serialNumbers);
    }

    async deleteSerialNumbers(serialNumbers: Array<SerialNumberModel>) {
        return await serialNumberDao.deleteSerialNumbers(serialNumbers);
    }
}

export default SerialNumberService.getInstance();