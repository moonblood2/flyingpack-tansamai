import {UserApiDto} from "../dto/userApi.model";
import usersApiDao from "../dao/userApi.dao";
import apiKey from "../../common/common.apikey";
import apiKeyConfig from "../../config/apikey.config";

import debug from "debug";

const log: debug.IDebugger = debug("app:user-api-service");

class UserApiService {
    private static instance: UserApiService;

    static getInstance(): UserApiService {
        if (!UserApiService.instance) {
            UserApiService.instance = new UserApiService();
        }
        return UserApiService.instance;
    }

    async generateNewApiKey(resource: UserApiDto) {
        //Generate api key.
        const key = apiKey.generateKey(apiKeyConfig.API_KEY_LENGTH, apiKeyConfig.API_KEY_PREFIX_LENGTH);
        resource.apiKeyPrefix = <string> key.prefixApiKey;
        resource.apiKeyHash = key.hash;
        //Save to database
        await usersApiDao.upsertApiKey(resource);

        return key.apiKey;
    }

    async getByApiKeyPrefix(apiKeyPrefix: string) {
        return await usersApiDao.getByApiKeyPrefix(apiKeyPrefix);
    }

    async getHookUrlOrder(userId: string) {
        return await usersApiDao.getHookUrlOrder(userId);
    }

    async getHookUrlOrderBySpOrderParcelId(anOrderId: string) {
        return await usersApiDao.getHookUrlOrderBySpOrderParcelId(anOrderId);
    }
}

export default UserApiService.getInstance();