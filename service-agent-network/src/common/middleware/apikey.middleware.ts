import apiKey, {Key} from "../common.apikey";
import apiKeyConfig from "../../config/apikey.config";
import userApiService from "../../userApi/service/userApi.service";

import express from "express";
import httpStatus from "http-status-codes";

import debug from "debug";
import {user} from "ts-postgres/dist/src/defaults";

const log: debug.IDebugger = debug("common-apikey-middleware:");

class ApiKeyMiddleware {
    private static instance: ApiKeyMiddleware;

    static getInstance(): ApiKeyMiddleware {
        if (!ApiKeyMiddleware.instance) {
            ApiKeyMiddleware.instance = new ApiKeyMiddleware();
        }
        return ApiKeyMiddleware.instance;
    }

    async verifyApiKey(req: express.Request, res: express.Response, next: express.NextFunction) {
        let valid = true;
        //Check empty apiKey
        const reqApiKey = req.body.apiKey;
        if (reqApiKey) {
            //Get apiKeyHash by userApi service, getByApiKeyPrefix
            const prefix = apiKey.retrieveApiKeyPrefix(reqApiKey, apiKeyConfig.API_KEY_PREFIX_LENGTH);
            const userApi = await userApiService.getByApiKeyPrefix(prefix);
            //If have requested api key in database.
            if (userApi) {
                const key: Key = {
                    apiKey: reqApiKey,
                    hash: userApi.apiKeyHash,
                };
                if (apiKey.verifyKey(key)) {
                    //Pass userId to the next handler.
                    req.body.userId = userApi.userId;
                    next();
                } else {
                    res.status(httpStatus.UNAUTHORIZED).send({
                        "status": false,
                        "message": "apiKey invalid"
                    })
                }
            } else {
                res.status(httpStatus.UNAUTHORIZED).send({
                    "status": false,
                    "message": "apiKey invalid"
                })
            }
        } else {
            res.status(httpStatus.UNAUTHORIZED).send({
                "status": false,
                "message": "apiKey invalid"
            })
        }
    }
}

export default ApiKeyMiddleware.getInstance();