import code, {HttpStatusCode} from "../../common/common.code";

import express from 'express';

import debug from 'debug';

const log: debug.IDebugger = debug('app:cod-controller');

class CodController {
    private static instance: CodController;

    static getInstance(): CodController {
        if (!CodController.instance) {
            CodController.instance = new CodController();
        }
        return CodController.instance;
    }

    async getCodReport(req: express.Request, res: express.Response) {
        res.status(200).send({"status": 1})
    }

}

export default CodController.getInstance();