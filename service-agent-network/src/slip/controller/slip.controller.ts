import slipService from "../service/slip.service";

import express from "express";
import httpStatus from "http-status-codes";

import debug from "debug";
import code, {HttpStatusCode} from "../../common/common.code";

const log: debug.IDebugger = debug("app:user-api-controller");

class SlipController {
    private static instance: SlipController;

    static getInstance(): SlipController {
        if (!SlipController.instance) {
            SlipController.instance = new SlipController();
        }
        return SlipController.instance;
    }

    async getSlip(req: express.Request, res: express.Response) {
        const userId: string = <string>req.headers["user-id"];

        let {slip, err} = await slipService.getSlip(userId);
        if (err.code === code.SUCCESS) {
            res.status(httpStatus.OK).send({
                ...err,
                data: {
                    slip: slip
                },
            });
        } else {
            res.status(HttpStatusCode[<number>err.code]).send({...err});
        }
    }
}

export default SlipController.getInstance();