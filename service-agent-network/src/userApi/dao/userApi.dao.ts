import {UserApiDto} from '../dto/userApi.model'
import postgresService from "../../common/service/postgres.service";
import code from "../../common/common.code";
import {NewCommonError} from "../../common/common.error";

import debug from 'debug';

const log: debug.IDebugger = debug('user-api-dao');

class UserDao {
    private static instance: UserDao;

    static getInstance(): UserDao {
        if (!UserDao.instance) {
            UserDao.instance = new UserDao();
        }
        return UserDao.instance;
    }

    async upsertApiKey(userApi: UserApiDto) {
        const queryText = `INSERT INTO public.user_api(user_id, api_key_prefix, api_key_hash)
         VALUES ($1, $2, $3)
         ON CONFLICT (user_id)
         DO UPDATE SET api_key_prefix=$4, api_key_hash=$5;`;
        const values = [userApi.userId, userApi.apiKeyPrefix, userApi.apiKeyHash, userApi.apiKeyPrefix, userApi.apiKeyHash];
        try {
            await postgresService.getClient().query(queryText, values);
        } catch (e) {
            log(e);
            throw new Error(e);
        }
    }

    async getByApiKeyPrefix(apiKeyPrefix: string) {
        const queryText = `SELECT * FROM public.user_api WHERE deleted_at=0 AND api_key_prefix=$1`;
        const values = [apiKeyPrefix];
        let userApi;
        try {
            const {rows, rowCount} = await postgresService.getClient().query(queryText, values);
            if (rowCount > 0) {
                const row = rows[0];
                userApi = {
                    userId: row["user_id"],
                    apiKeyPrefix: row["api_key_prefix"],
                    apiKeyHash: row["api_key_hash"],
                    urlWebHookOrder: row["url_web_hook_order"],
                    createdAt: row["created_at"],
                    deletedAt: row["deleted_at"],
                }
            } else {
                userApi = null;
            }
        } catch (e) {
            log(e);
            throw new Error(e);
        }
        return userApi;
    }

    async getHookUrlOrder(userId: string) {
        const queryText = `SELECT hook_url_order FROM public.user_api WHERE user_id=$1`;
        const values = [userId];
        let url = "";
        let err = NewCommonError();
        try {
            const {rows, rowCount} = await postgresService.getClient().query(queryText, values);
            log(rows);
            if (rowCount > 0) {
                url = rows[0]["hook_url_order"];
            } else {
                log("Not found user_id %s", userId);
                err = NewCommonError(code.ERR_INTERNAL);
            }
        } catch (e) {
            log(e);
            err = NewCommonError(code.ERR_INTERNAL);
        }
        return {url, err};
    }

    async getHookUrlOrderBySpOrderParcelId(spOrderParcelId: string) {
        const queryText = `
            SELECT hook_url_order 
            FROM public.user_api 
            WHERE user_id=(
                SELECT user_id 
                FROM public.order 
                WHERE sp_order_parcel_id=$1 
                LIMIT 1
            )`;
        const values = [spOrderParcelId];
        let url = "";
        let err = NewCommonError();
        try {
            const {rows, rowCount} = await postgresService.getClient().query(queryText, values);
            if (rowCount > 0) {
                url = rows[0]["hook_url_order"];
            } else {
                log("Not found hook_url_order by sp_order_parcel_id %s", spOrderParcelId);
                err = NewCommonError(code.ERR_INTERNAL);
            }
        } catch (e) {
            log(e);
            err = NewCommonError(code.ERR_INTERNAL);
        }
        return {url, err};
    }
}

export default UserDao.getInstance();