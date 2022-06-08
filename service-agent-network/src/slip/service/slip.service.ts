import slipDao from "../dao/slip.dao";

import debug from "debug";

const log: debug.IDebugger = debug("app:slip-service");

class SlipService {
    private static instance: SlipService;

    static getInstance(): SlipService {
        if (!SlipService.instance) {
            SlipService.instance = new SlipService();
        }
        return SlipService.instance;
    }

    async getSlip(userId: string) {
        return await slipDao.getSlip(userId);
    }
}

export default SlipService.getInstance();