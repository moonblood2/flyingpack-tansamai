import {SerialNumberModel} from "../dto/serialNumber.model";
import {AnOrderItemModel} from "../../order/dto/anOrderItem.model";

import postgresService from "../../common/service/postgres.service";

import {NewCommonError} from "../../common/common.error";
import code from "../../common/common.code";

import {v4 as uuidv4} from 'uuid';
import debug from "debug";

const log: debug.IDebugger = debug("app:serial-number-dao");

class SerialNumberDao {
    private static instance: SerialNumberDao;

    private constructor() {
    }

    static getInstance(): SerialNumberDao {
        if (!SerialNumberDao.instance) {
            SerialNumberDao.instance = new SerialNumberDao();
        }
        return SerialNumberDao.instance;
    }

    async getSerialNumbersByAnOrderId(anOrderId: string) {
        let err = NewCommonError();
        let serialNumbers: Array<SerialNumberModel> = [];

        const queryText = `
WITH order_cte AS(
    SELECT * 
    FROM public.order
    WHERE deleted_at=0 AND an_order_id=$1
    LIMIT 1
)
SELECT aops.an_order_product_serial_number_id, aops.serial_number_start, aops.serial_number_end, aops.created_at as aops_created_at, aops.deleted_at as aops_deleted_at
FROM order_cte o
LEFT JOIN public.order_product op
ON o.an_order_id=op.an_order_id
LEFT JOIN public.product p
ON op.an_product_id=p.an_product_id
LEFT JOIN public.an_order_product_serial_number aops
ON op.an_order_product_id=aops.an_order_product_id
WHERE aops.deleted_at=0 OR aops.deleted_at IS NULL
ORDER BY p.product_code, aops.an_order_product_serial_number_id
        `;
        const values = [anOrderId];
        try {
            const {rows, rowCount} = await postgresService.getClient().query(queryText, values);
            if (rowCount > 0) {
                //Mapping row to object
                for (const row of rows) {
                    serialNumbers.push({
                        anOrderProductSerialNumberId: row['an_order_product_serial_number_id'],
                        serialNumberStart: row['serial_number_start'],
                        serialNumberEnd: row['serial_number_end'],
                        createdAt: row['aops_created_at'],
                        deletedAt: row['aops_deleted_at'],
                    });
                }
            }
        } catch (e) {
            log(e);
            err = NewCommonError(code.ERR_INTERNAL);
        }
        return {serialNumbers, err}
    }

    async createSerialNumbers(items: Array<AnOrderItemModel>) {
        let err = NewCommonError();

        try {
            let n = 0;
            let queryText = `
            INSERT INTO public.an_order_product_serial_number(
                an_order_product_serial_number_id, an_order_product_id, serial_number_start, serial_number_end)
                VALUES`;
            let valueQueryText = []; // For concatenating with VALUES
            let values = [];

            for (const i of items) {
                for (const s of i.serialNumbers) {
                    //assign id
                    s.anOrderProductId = i.anOrderProductId;
                    s.anOrderProductSerialNumberId = uuidv4();
                    values.push(
                        s.anOrderProductSerialNumberId,
                        s.anOrderProductId,
                        s.serialNumberStart ? s.serialNumberStart: '',
                        s.serialNumberEnd ? s.serialNumberEnd: '',
                    );
                    valueQueryText.push(`($${n + 1}, $${n + 2}, $${n + 3}, $${n + 4})`);
                    n += 4;
                }
            }

            queryText = `${queryText} ${valueQueryText.join(', ')}`;
            if (n > 0) {
                await postgresService.getClient().query(queryText, values);
            }
        } catch (e) {
            log(e);
            err = NewCommonError(code.ERR_INTERNAL);
        }

        return {items, err};
    }

    async updateSerialNumbers(serialNumbers: Array<SerialNumberModel>) {
        let err = NewCommonError();
        try {
            let n = 0;
            let valueQueryText = []; // For concatenating with VALUES
            let values = [];
            for (const o of serialNumbers) {
                if (o.anOrderProductSerialNumberId) {
                    valueQueryText.push(`($${n + 1}, $${n + 2}, $${n + 3})`);
                    values.push(
                        o.anOrderProductSerialNumberId,
                        o.serialNumberStart ? o.serialNumberStart: '',
                        o.serialNumberEnd ? o.serialNumberEnd: '',
                    )
                    n += 3;
                }
            }
            let queryText = `
UPDATE public.an_order_product_serial_number
SET serial_number_start=tmp.serial_number_start, serial_number_end=tmp.serial_number_end 
FROM (VALUES ${valueQueryText.join(", ")}) AS tmp (id, serial_number_start, serial_number_end)
WHERE an_order_product_serial_number_id=tmp.id::uuid
`;
            if (n > 0) {
                await postgresService.getClient().query(queryText, values);
            }
        } catch (e) {
            log(e);
            err = NewCommonError(code.ERR_INTERNAL);
        }
        return {err, serialNumbers};
    }

    async deleteSerialNumbers(serialNumbers: Array<SerialNumberModel>) {
        let err = NewCommonError();
        try {
            let n = 0;
            let valueQueryText = []; // For concatenating with VALUES
            let values = [];
            for (const o of serialNumbers) {
                if (o.anOrderProductSerialNumberId) {
                    valueQueryText.push(`($${n + 1})`);
                    values.push(
                        o.anOrderProductSerialNumberId,
                    );
                    n += 1;
                }
            }
            let queryText = `
UPDATE public.an_order_product_serial_number 
SET deleted_at=extract(epoch from now()) 
FROM (VALUES ${valueQueryText.join(", ")}) AS tmp (id)
WHERE an_order_product_serial_number_id=tmp.id::uuid         
            `;
            if (n > 0) {
                await postgresService.getClient().query(queryText, values);
            }
        } catch (e) {
            log(e);
            err = NewCommonError(code.ERR_INTERNAL);
        }
        return {err};
    }
}

export default SerialNumberDao.getInstance();