import userApiService from "../service/userApi.service";

import express from "express";
import httpStatus from "http-status-codes";

import debug from "debug";

const log: debug.IDebugger = debug("app:user-api-controller");

class UserApiController {
    private static instance: UserApiController;

    static getInstance(): UserApiController {
        if (!UserApiController.instance) {
            UserApiController.instance = new UserApiController();
        }
        return UserApiController.instance;
    }

    async generateNewApiKey(req: express.Request, res: express.Response) {
        const apiKey = await userApiService.generateNewApiKey(req.body);

        res.status(httpStatus.CREATED).send({
            apiKey: apiKey,
        });
    }
}

export default UserApiController.getInstance();