import {SlipDto} from '../dto/slip.model'
import postgresService from "../../common/service/postgres.service";
import code from "../../common/common.code";
import {NewCommonError} from "../../common/common.error";

import debug from 'debug';

const log: debug.IDebugger = debug('slip-dao');

class SlipDao {
    private static instance: SlipDao;

    static getInstance(): SlipDao {
        if (!SlipDao.instance) {
            SlipDao.instance = new SlipDao();
        }
        return SlipDao.instance;
    }

    async getSlip(userId: string) {
        const queryText = `
            SELECT message 
            FROM public.slip 
            WHERE user_id=$1`;
        const values = [userId];
        let slip: SlipDto = {userId: "", message: ""};
        let err = NewCommonError();
        try {
            const {rows, rowCount} = await postgresService.getClient().query(queryText, values);
            if (rowCount > 0) {
                slip.userId = userId;
                slip.message = rows[0]['message'];
                slip.createdAt = rows[0]['created_at'];
                slip.deletedAt = rows[0]['deleted_at'];
            } else {
                log("Not found slip by user_id %s", userId);
                err = NewCommonError(code.ERR_INTERNAL);
            }
        } catch (e) {
            log(e);
            err = NewCommonError(code.ERR_INTERNAL);
        }
        return {slip, err};
    }
}

export default SlipDao.getInstance();