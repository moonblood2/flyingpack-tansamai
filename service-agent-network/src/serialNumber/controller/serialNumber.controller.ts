import serialNumberService from "../service/serialNumber.service";
import hookService from "../../hook/service/hook.service";

import code, {HttpStatusCode} from "../../common/common.code";

import express from "express";

import debug from "debug";

const log: debug.IDebugger = debug("app:open-api-order-controller");

class SerialNumberController {
    private static instance: SerialNumberController;

    static getInstance(): SerialNumberController {
        if (!SerialNumberController.instance) {
            SerialNumberController.instance = new SerialNumberController();
        }
        return SerialNumberController.instance;
    }

    async createSerialNumber(req: express.Request, res: express.Response) {
        const userId = <string> req.headers["user-id"];

        let {err, items} = await serialNumberService.createSerialNumbers(req.body.items);
        if (err.code === code.SUCCESS) {
            await hookService.sendSerialNumber(userId, <string> req.body.anOrderId);
            res.status(HttpStatusCode[<number>err.code]).send({
                ...err,
                data: {
                    items,
                }
            });
        } else {
            res.status(HttpStatusCode[<number>err.code]).send({...err});
        }
    }

    async updateSerialNumber(req: express.Request, res: express.Response) {
        const userId = <string> req.headers["user-id"];

        let {err, serialNumbers} = await serialNumberService.updateSerialNumbers(req.body.serialNumbers);
        if (err.code === code.SUCCESS) {
            await hookService.sendSerialNumber(userId, <string> req.body.anOrderId);
            res.status(HttpStatusCode[<number>err.code]).send({
                ...err,
                data: {
                    serialNumbers,
                }
            });
        } else {
            res.status(HttpStatusCode[<number>err.code]).send({...err});
        }
    }

    async deleteSerialNumber(req: express.Request, res: express.Response) {
        const userId = <string> req.headers["user-id"];

        let {err} = await serialNumberService.deleteSerialNumbers(req.body.serialNumbers);
        if (err.code === code.SUCCESS) {
            res.status(HttpStatusCode[<number>err.code]).send({
                ...err,
            });
        } else {
            res.status(HttpStatusCode[<number>err.code]).send({...err});
        }
    }
}

export default SerialNumberController.getInstance();