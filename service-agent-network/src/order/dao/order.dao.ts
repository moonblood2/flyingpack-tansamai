import postgresService from "../../common/service/postgres.service";

import {OrderDto} from "../dto/order.model";
import {OrderProductModel} from "../dto/orderProduct.model";
import {OrderSummaryModel} from "../dto/orderSummary.model";
import {OrderProductSummaryModel} from "../dto/orderProductSummary.model";
import {BankAccountModel} from "../dto/bankAccount.model";
import {AnOrderItemModel} from "../dto/anOrderItem.model";
import {SerialNumberModel} from "../../serialNumber/dto/serialNumber.model";

import {ListOrderParam} from "../param/list.order.param";
import {DateParam} from "../param/date.param";
import {DateTypeEnum} from "../param/dateType.param";

import {AnOrderFulfillmentStatusEnum, AnOrderFulfillmentStatusToString} from "../dto/anOrderFulfillmentStatus.enum";

import code from "../../common/common.code";
import {CommonError, NewCommonError} from "../../common/common.error";

import fiCodeToFiName from "../../constant/fiCode";

import {v4 as uuidv4} from 'uuid';
import debug from "debug";
import {AnCourierCodeEnum} from "../dto/anCourierCode.enum";

const log: debug.IDebugger = debug('app:order-dao')

class OrderDao {
    private static instance: OrderDao;

    static getInstance(): OrderDao {
        if (!OrderDao.instance) {
            OrderDao.instance = new OrderDao();
        }
        return OrderDao.instance;
    }

    async getAnOrderIdByReferenceNo(referenceNo: string) {
        let anOrderId: string = "";
        let err = NewCommonError();

        let queryText: string = `SELECT * FROM public.order WHERE reference_no=$1`;
        try {
            const {rows, rowCount} = await postgresService.getClient().query(queryText, [referenceNo]);
            if (rowCount > 0) {
                anOrderId = rows[0]['an_order_id'];
            }
        } catch (err) {
            log(err)
            err = NewCommonError(code.ERR_INTERNAL);
            return {err};
        }
        return {anOrderId, err};
    }

    async addOrder(order: OrderDto, orderProducts: Array<OrderProductModel>) {
        order.anOrderId = uuidv4();

        let err: CommonError = NewCommonError();
        console.log("bedore query : ",order.trackingCode);
        try {
            await postgresService.getClient().query("BEGIN");

            /**Insert order*/
            let queryText = `INSERT INTO public.order(
            an_order_id, user_id, sp_order_parcel_id, reference_no, des_name, 
            des_phone_number, des_address, des_subdistrict, des_district, des_province, 
            des_postcode, courier_code, cod_amount, fulfillment_status,tracking_code,sort_code,line_code,sorting_line_code,dst_store_name) 
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17 ,$18 ,$19 );`;
            let values = [
                order.anOrderId,
                order.userId,
                order.spOrderParcelId,
                order.referenceNo,
                order.desName,
                order.desPhoneNumber,
                order.desAddress,
                order.desSubdistrict,
                order.desDistrict,
                order.desProvince,
                order.desPostcode,
                order.courierCode,
                order.codAmount ? order.codAmount : null,
                AnOrderFulfillmentStatusEnum.NOT_PACKED,
                order.trackingCode,
                order.sortCode,
                order.lineCode,
                order.sortingLineCode,
                order.dstStoreName
            ];
            console.log('queryText : ',queryText);
            console.log('values : ',values);
            await postgresService.getClient().query(queryText, values);

            /**Insert each order_product*/
            for (const p of orderProducts) {
                p.anOrderId = order.anOrderId;
                p.anOrderProductId = uuidv4();
                const queryText = `INSERT INTO public.order_product(an_order_product_id, an_order_id, an_product_id, quantity) VALUES ($1, $2, $3, $4);`;
                const values = [p.anOrderProductId, p.anOrderId, p.anProductId, p.quantity];

                await postgresService.getClient().query(queryText, values);
            }

            /**Insert bank account*/
            if (order.codAmount > 0) {
                let bankAccount = <BankAccountModel>order.bankAccount;
                queryText = `INSERT INTO public.bank_account(bank, account_no, account_name, email, fi_code, an_order_id) VALUES ($1, $2, $3, $4, $5, $6)`;
                values = [bankAccount.bank, bankAccount.accountNo, bankAccount.accountName, bankAccount.email, bankAccount.fiCode, order.anOrderId];
                await postgresService.getClient().query(queryText, values);
            }

            await postgresService.getClient().query("COMMIT");
        } catch (e) {
            await postgresService.getClient().query("ROLLBACK");
            log(e);
            err = NewCommonError(code.ERR_INTERNAL);
            if (e.code === "23505" && e.detail.includes(order.referenceNo)) {
                err = NewCommonError(code.INPUT_REFERENCE_NO_ALREADY_EXIST);
            }
        }
        return {
            order,
            err,
        }
    }

    static async getAnOrderIdBySerialNumber(serialNumber: string) {
        let orderId: string = "";
        let exist: boolean = false;
        let err = null;

        const queryText = `
            SELECT aop.an_order_id FROM public.an_order_product_serial_number aops
            LEFT JOIN public.order_product aop
            ON aops.an_order_product_id=aop.an_order_product_id
            WHERE 
                ($1=serial_number_start) OR
                ($1=serial_number_end) OR 
                ($1>serial_number_start AND $1<serial_number_end) 
            LIMIT 1;
            `;
        try {
            const {rows, rowCount} = await postgresService.getClient().query(queryText, [serialNumber]);
            if (rowCount > 0) {
                exist = true
                orderId = rows[0]["an_order_id"];
            }
        } catch (err) {
            err = NewCommonError(code.ERR_INTERNAL);
            return {orderId, err, exist};
        }

        return {orderId, err, exist}
    }

    static async getOrder(by: string, key: string) {
        let order: OrderDto = {} as OrderDto;
        let err = NewCommonError();

        let andWhere = `AND ${by}=$1`;
        if (by === "sp_order_parcel_id" || by === 'an_order_id') {
            andWhere = `AND (sp_order_parcel_id=$1 OR an_order_id=$1)`;
        } else if (by === "text") {
            let {orderId, exist: serialNumberExist, err} = await OrderDao.getAnOrderIdBySerialNumber(key);
            if (err != null) {
                return {order, err};
            }

            if (serialNumberExist) {
                key = orderId;
                andWhere = `AND (an_order_id=$1::uuid)`;
            } else {
                andWhere = `AND (reference_no=$1 OR tracking_code=$1)`;
            }
        }

        const queryText = `
WITH order_cte AS(
    SELECT * 
    FROM public.order
    WHERE deleted_at=0
    ${andWhere}
    LIMIT 1
)
SELECT o.*,
op.an_order_product_id, p.an_product_id, p.product_code, p.name, p.img_url, op.quantity, p.serial_regex, p.robotic_sku,
aops.an_order_product_serial_number_id, aops.serial_number_start, aops.serial_number_end, aops.created_at as aops_created_at, aops.deleted_at as aops_deleted_at
FROM order_cte o
LEFT JOIN public.order_product op
ON o.an_order_id=op.an_order_id
LEFT JOIN public.product p
ON op.an_product_id=p.an_product_id
LEFT JOIN public.an_order_product_serial_number aops
ON op.an_order_product_id=aops.an_order_product_id
WHERE aops.deleted_at=0 OR aops.deleted_at IS NULL
ORDER BY p.product_code, aops.an_order_product_serial_number_id`;
        log(queryText)
        const values = [key];
        try {
            const {rows, rowCount} = await postgresService.getClient().query(queryText, values);
            if (rowCount > 0) {
                const row = rows[0];
                order = {
                    anOrderId: row["an_order_id"],
                    spOrderParcelId: row["sp_order_parcel_id"],
                    referenceNo: row["reference_no"],
                    desName: row["des_name"],
                    desPhoneNumber: row["des_phone_number"],
                    desAddress: row["des_address"],
                    desSubdistrict: row["des_subdistrict"],
                    desDistrict: row["des_district"],
                    desProvince: row["des_province"],
                    desPostcode: row["des_postcode"],
                    courierCode: row["courier_code"],
                    codAmount: row["cod_amount"] ? parseFloat(row["cod_amount"]) : 0,
                    createdAt: row["created_at_tmz"],
                    trackingCode: row["tracking_code"],
                    fulfillmentStatus: row["fulfillment_status"],
                    items: new Array<AnOrderItemModel>(),
                }
                let productIdToItemsIndex: Map<string, number> = new Map<string, number>();
                //Mapping row to object
                for (const row of rows) {
                    // Push item of each order
                    if (!productIdToItemsIndex.has(row["an_product_id"])) {
                        productIdToItemsIndex.set(row["an_product_id"], (<any[]>order.items).length);
                        (<AnOrderItemModel[]>order.items).push({
                            anOrderProductId: row['an_order_product_id'],
                            anProductId: row["an_product_id"],
                            productCode: row["product_code"],
                            name: row["name"],
                            imgUrl: row["img_url"],
                            quantity: row["quantity"],
                            serialRegex: row["serial_regex"],
                            roboticSKU: row["robotic_sku"],
                            serialNumbers: new Array<SerialNumberModel>(),
                        });
                    }
                    // Push serialNumber of each item.
                    if (row['an_order_product_serial_number_id']) {
                        (<AnOrderItemModel[]>order.items)[<number>productIdToItemsIndex.get(row['an_product_id'])].serialNumbers.push({
                            anOrderProductSerialNumberId: row['an_order_product_serial_number_id'],
                            serialNumberStart: row['serial_number_start'],
                            serialNumberEnd: row['serial_number_end'],
                            createdAt: row['aops_created_at'],
                            deletedAt: row['aops_deleted_at'],
                        });
                    }
                }
            }
        } catch (e) {
            log(e);
            err = NewCommonError(code.ERR_INTERNAL);
        }
        return {order, err}
    }

    async getOrderBySpOrderParcelId(spOrderParcelId: string) {
        return await OrderDao.getOrder('sp_order_parcel_id', spOrderParcelId);
    }

    async getOrderByAnOrderId(anOrderId: string) {
        return await OrderDao.getOrder('an_order_id', anOrderId);
    }

    async getOrderByKey(key: string) {
        return await OrderDao.getOrder('text', key);
    }

    async getOrders(params: ListOrderParam, isSorting: boolean, userId: string) {
        const offset = (params.page - 1) * params.perPage;
        const limit = params.perPage;

        let orders: Array<OrderDto> = [];
        let err = NewCommonError();
        let totalItem: number = 0;

        try {
            /**
             * Param:
             * $1=user_id
             * $2=limit, $3=offset
             * $4=startDate, $5=endDate, $6=timezone,
             * $7=reference_no
             * $8=fulfillment_status
             * $9=courier_code
             * $10-$N=an_product_id
             *
             * CTE Note:
             * order_product_count_cte, count product of each order.
             * filter_order_cte, filter by where clause.
             * order_sorted_cte, sort order by number of product of each order, quantity, product_code.
             * order_distinct_limit_cte, pagination.
             * result_cte, select only required columns.
             */

            let values = [];
            let whereDeletedAt = `WHERE o.deleted_at=0`;

            let andWhereUser = `AND o.user_id=$1`;
            values.push(userId);

            let limitOffset = `LIMIT $2 OFFSET $3`;
            values.push(limit, offset);

            let dateType = "created_at";
            if (params.dateType === DateTypeEnum.STATUS_COMPLETED_DATE) {
                dateType = "status_completed_date";
            } else if (params.dateType === DateTypeEnum.COD_TRANSFERRED_DATE) {
                dateType = "cod_transferred_date";
            }
            let andWhereCreatedAt = `AND o.${dateType} AT TIME ZONE $6 BETWEEN $4::timestamp AND $5::timestamp`;
            values.push(params.startDate, params.endDate, <string>params.timeZone);

            let andWhereKeyWord = `AND (
                o.reference_no LIKE CONCAT('%', $7::text, '%') OR 
                o.tracking_code=$7 OR
                o.des_name LIKE CONCAT('%', $7::text, '%') OR
                ba.account_no LIKE CONCAT('%', $7::text, '%') OR
                ba.account_name LIKE CONCAT('%', $7::text, '%')
            )`;
            values.push(params.keyWord);

            let andWhereSpOrderParcelId = ``; //Check booking status
            if (params.fulfillmentStatus) {
                andWhereSpOrderParcelId = `AND o.fulfillment_status = $8`;
                values.push(params.fulfillmentStatus);
            }

            let andWhereCourierCode = ``;
            if (params.courierCode in AnCourierCodeEnum) {
                andWhereCourierCode = `AND o.courier_code=$${values.length + 1}`;
                values.push(params.courierCode);
            }

            let andWhereAnProductId = ``;
            if (params.productIds.length > 0) {
                const N: number = values.length + 1;
                let productIdPlaceholderArr: Array<string> = [];
                for (let i = 0; i < params.productIds.length; i++) {
                    productIdPlaceholderArr.push(`$${N + i}`);
                }
                andWhereAnProductId = `AND op.an_product_id IN (${productIdPlaceholderArr.join(", ")})`;
                values.push(...params.productIds);
            }

            let andWhereCod = ``;
            if (params.isCod) {
                andWhereCod = `AND (o.cod_amount > 0 AND o.cod_amount IS NOT NULL)`;
            }

            let queryText = `
WITH 
order_product_count_cte AS 
(
    SELECT o.an_order_id, COUNT(op.an_product_id) AS product_count
    FROM public.order o
    LEFT JOIN public.order_product AS op
    ON o.an_order_id = op.an_order_id
    GROUP BY o.an_order_id
),
filter_order_cte AS
(
    SELECT o.an_order_id, opc.product_count
    FROM public.order o
    LEFT JOIN order_product_count_cte opc
    ON o.an_order_id = opc.an_order_id
    LEFT JOIN public.order_product op
    ON o.an_order_id = op.an_order_id
    LEFT JOIN public.bank_account ba
    ON o.an_order_id = ba.an_order_id
    ${whereDeletedAt}
    ${andWhereUser}
    ${andWhereSpOrderParcelId}
    ${andWhereCreatedAt}
    ${andWhereKeyWord}
    ${andWhereCourierCode}
    ${andWhereAnProductId}
    ${andWhereCod}
),
order_sorted_cte AS 
(
    SELECT ROW_NUMBER() OVER(
        ORDER BY fo.product_count ASC, op.quantity ASC, p.product_code ASC, o.created_at DESC
    ) AS row_number,
    ROW_NUMBER() OVER(
        PARTITION BY o.an_order_id 
        ORDER BY fo.product_count ASC, op.quantity ASC, p.product_code ASC, o.created_at DESC
    ) AS row_number_of_group,
    o.an_order_id, fo.product_count
    FROM filter_order_cte fo
    LEFT JOIN public.order AS o
    ON fo.an_order_id = o.an_order_id
    LEFT JOIN public.order_product op
    ON o.an_order_id = op.an_order_id
    LEFT JOIN public.product p
    ON op.an_product_id = p.an_product_id
    ORDER BY fo.product_count ASC, op.quantity ASC, p.product_code ASC, o.created_at DESC
),
order_distinct_limit_cte AS
(
    SELECT *,
    (
        SELECT COUNT(DISTINCT an_order_id) 
        FROM order_sorted_cte
        WHERE row_number_of_group=1
    ) AS total_item
    FROM order_sorted_cte
    WHERE row_number_of_group=1
    ${limitOffset}
),
result_cte AS
(
    SELECT order_distinct_limit_cte.row_number, o.an_order_id, o.sp_order_parcel_id, o.reference_no, 
    o.des_name, o.des_phone_number, o.des_address, o.des_subdistrict, o.des_district, o.des_province, o.des_postcode, 
    o.courier_code, o.cod_amount, o.tracking_code,o.sort_code,o.line_code,o.sorting_line_code,o.dst_store_name,
    o.fulfillment_status, o.shipping_status, o.cod_status,
    op.an_order_product_id, p.an_product_id, p.product_code, p.name, p.img_url, p.serial_regex, p.robotic_sku,
    aops.an_order_product_serial_number_id, aops.serial_number_start, aops.serial_number_end, aops.created_at as aops_created_at, aops.deleted_at as aops_deleted_at,
    ba.bank, ba.account_name, ba.account_no, ba.email, ba.fi_code,
    o.created_at AT TIME ZONE $6 AS created_at_tmz, o.status_completed_date AT TIME ZONE $6 status_completed_date_tmz, o.cod_transferred_date AT TIME ZONE $6 cod_transferred_date_tmz,
    o.jna_cod_transferred_date AT TIME ZONE $6 jna_cod_transferred_date_tmz,
    op.quantity, total_item
    FROM order_distinct_limit_cte
    LEFT JOIN public.order o
    ON order_distinct_limit_cte.an_order_id = o.an_order_id
    LEFT JOIN public.order_product op
    ON o.an_order_id = op.an_order_id
    LEFT JOIN public.product p
    ON op.an_product_id = p.an_product_id
    LEFT JOIN public.bank_account ba
    ON o.an_order_id = ba.an_order_id
    LEFT JOIN public.an_order_product_serial_number aops
    ON op.an_order_product_id=aops.an_order_product_id
    ORDER BY order_distinct_limit_cte.row_number, p.product_code, aops.an_order_product_serial_number_id
)

SELECT * FROM result_cte`;

            if (!isSorting) {
                queryText = `
WITH 
filter_order_cte AS
(
    SELECT o.an_order_id
    FROM public.order o
    LEFT JOIN public.order_product op
    ON o.an_order_id = op.an_order_id
    LEFT JOIN public.bank_account ba
    ON o.an_order_id = ba.an_order_id
    ${whereDeletedAt}
    ${andWhereUser}
    ${andWhereSpOrderParcelId}
    ${andWhereCreatedAt}
    ${andWhereKeyWord}
    ${andWhereCourierCode}
    ${andWhereAnProductId}
    ${andWhereCod}
    GROUP BY o.an_order_id
),
limit_order_cte AS
(
    SELECT * FROM filter_order_cte
    ${limitOffset}
),
result_cte AS
(
    SELECT o.an_order_id, o.sp_order_parcel_id, o.reference_no, 
    o.des_name, o.des_phone_number, o.des_address, o.des_subdistrict, o.des_district, o.des_province, o.des_postcode, 
    o.courier_code, o.cod_amount, o.tracking_code,o.sort_code,o.line_code,o.sorting_line_code,o.dst_store_name,
    o.fulfillment_status, o.shipping_status, o.cod_status,
    op.an_order_product_id, p.an_product_id, p.product_code, p.name, p.img_url, p.serial_regex, p.robotic_sku,
    aops.an_order_product_serial_number_id, aops.serial_number_start, aops.serial_number_end, aops.created_at as aops_created_at, aops.deleted_at as aops_deleted_at,
    ba.bank, ba.account_name, ba.account_no, ba.email, ba.fi_code,
    o.created_at AT TIME ZONE $6 AS created_at_tmz, o.status_completed_date AT TIME ZONE $6 status_completed_date_tmz, o.cod_transferred_date AT TIME ZONE $6 cod_transferred_date_tmz,
    o.jna_cod_transferred_date AT TIME ZONE $6 jna_cod_transferred_date_tmz,
    op.quantity
    FROM limit_order_cte
    LEFT JOIN public.order o
    ON limit_order_cte.an_order_id = o.an_order_id
    LEFT JOIN public.order_product op
    ON limit_order_cte.an_order_id = op.an_order_id
    LEFT JOIN public.product p
    ON op.an_product_id = p.an_product_id
    LEFT JOIN public.bank_account ba
    ON o.an_order_id = ba.an_order_id
    LEFT JOIN public.an_order_product_serial_number aops
    ON op.an_order_product_id=aops.an_order_product_id
    ORDER BY o.created_at DESC
)
SELECT * FROM result_cte;
`;
            }

            const {rows, rowCount} = await postgresService.getClient().query(queryText, values);
            //Mapping orders. orders (reference_no) -> items (product_id) -> serialNumbers.
            //Each order has unique reference_no and contain item list (items).
            //Each items has unique product_id and contain serial numbers list (serialNumbers).
            if (rowCount > 0) {
                //key is referenceNo and value is index of orders.
                let refNoToOrdersIndex: any = {};
                let refNoAndProductIdToItemsIndex: any = {};
                for (const row of rows) {
                    log(row);
                    const refNo = row["reference_no"];
                    const proId = row["an_product_id"];
                    //Push order. Check if not has refNo.
                    if (!refNoToOrdersIndex.hasOwnProperty(refNo)) {
                        refNoToOrdersIndex[refNo] = orders.length;
                        refNoAndProductIdToItemsIndex[refNo] = {};
                        orders.push({
                            anOrderId: row["an_order_id"],
                            spOrderParcelId: row["sp_order_parcel_id"],
                            referenceNo: row["reference_no"],
                            desName: row["des_name"],
                            desPhoneNumber: row["des_phone_number"],
                            desAddress: row["des_address"],
                            desSubdistrict: row["des_subdistrict"],
                            desDistrict: row["des_district"],
                            desProvince: row["des_province"],
                            desPostcode: row["des_postcode"],
                            courierCode: row["courier_code"],
                            codAmount: row["cod_amount"] ? parseFloat(row["cod_amount"]) : 0,
                            createdAt: row["created_at_tmz"],
                            codTransferredDate: row['cod_transferred_date_tmz'],
                            jnaCodTransferredDate: row['jna_cod_transferred_date_tmz'],
                            statusCompletedDate: row['status_completed_date_tmz'],
                            trackingCode: row["tracking_code"],
                            sortCode:row["sort_code"],
                            lineCode:row["line_code"],
                            sortingLineCode:row["sorting_line_code"],
                            dstStoreName:row["dst_store_name"],
                            fulfillmentStatus: row["fulfillment_status"],
                            fulfillmentStatusString: <string>AnOrderFulfillmentStatusToString[<AnOrderFulfillmentStatusEnum>row["fulfillment_status"]],
                            shippingStatus: row["shipping_status"],
                            codStatus: row["cod_status"],
                            items: new Array<AnOrderItemModel>(),
                            itemQuantitySum: 0,
                            bankAccount: row["cod_amount"] > 0 ? {
                                bank: row["bank"],
                                accountName: row["account_name"],
                                accountNo: row["account_no"],
                                email: row["email"],
                                fiCode: row["fi_code"],
                                fiName: (<any>fiCodeToFiName)[row["fi_code"]],
                            } : undefined,
                        });
                    }
                    //Push item.
                    if (refNoAndProductIdToItemsIndex.hasOwnProperty(refNo) && !refNoAndProductIdToItemsIndex[refNo].hasOwnProperty(proId)) {
                        //Sum item quantity of each order.
                        (<number>orders[<number>refNoToOrdersIndex[refNo]].itemQuantitySum) += Number(row["quantity"]);
                        refNoAndProductIdToItemsIndex[refNo][proId] = (<AnOrderItemModel[]>orders[<number>refNoToOrdersIndex[refNo]].items).length;
                        (<AnOrderItemModel[]>orders[<number>refNoToOrdersIndex[refNo]].items).push({
                            anOrderProductId: row['an_order_product_id'],
                            anProductId: row["an_product_id"],
                            productCode: row["product_code"],
                            name: row["name"],
                            imgUrl: row["img_url"],
                            quantity: row["quantity"],
                            serialRegex: row["serial_regex"],
                            roboticSKU: row["robotic_sku"],
                            serialNumbers: new Array<SerialNumberModel>(),
                        });
                    }
                    //Push serial number.
                    (<AnOrderItemModel[]>orders[<number>refNoToOrdersIndex[refNo]].items)[<number>refNoAndProductIdToItemsIndex[refNo][proId]].serialNumbers.push({
                        anOrderProductSerialNumberId: row['an_order_product_serial_number_id'],
                        serialNumberStart: row['serial_number_start'],
                        serialNumberEnd: row['serial_number_end'],
                        createdAt: row['aops_created_at'],
                        deletedAt: row['aops_deleted_at'],
                    });
                }
                if (isSorting) {
                    totalItem = parseInt(rows[0]["total_item"]);
                } else {
                    totalItem = orders.length;
                }
            }
        } catch (e) {
            log(e);
            err = NewCommonError(code.ERR_INTERNAL);
        }

        return {orders, totalItem, err};
    }

    async updateOrders(orders: Array<OrderDto>) {
        let err = NewCommonError();
        log(orders)
        try {
            for (const order of orders) {
                if (order.anOrderId) {
                    let values = [];
                    let queryText = "";
                    let plc = [];//placeholder
                    values.push(order.anOrderId);
                    if (order.spOrderParcelId) {
                        values.push(order.spOrderParcelId);
                        plc.push(`sp_order_parcel_id=$${values.length}`);
                    }
                    if (order.trackingCode) {
                        values.push(order.trackingCode);
                        plc.push(`tracking_code=$${values.length}`);
                    }
                    if (order.codAmount) {
                        values.push(order.codAmount);
                        plc.push(`cod_amount=$${values.length}`);
                    }
                    if (order.fulfillmentStatus >= 0) {
                        values.push(order.fulfillmentStatus);
                        plc.push(`fulfillment_status=$${values.length}`);
                    }
                    // if (order.shippingStatus) {
                    //     values.push(order.shippingStatus);
                    //     plc.push(`shipping_status=$${values.length}`);
                    // }
                    if (order.codStatus) {
                        values.push(order.codStatus);
                        plc.push(`cod_status=$${values.length}`);
                    }
                    if (order.statusCompletedDate) {
                        values.push(order.statusCompletedDate);
                        plc.push(`status_completed_date=$${values.length}`);
                    }
                    if (order.codTransferredDate) {
                        values.push(order.codTransferredDate);
                        plc.push(`cod_transferred_date=$${values.length}`);
                    }
                    if (order.jnaCodTransferredDate) {
                        values.push(order.jnaCodTransferredDate);
                        plc.push(`jna_cod_transferred_date=$${values.length}`);
                    }
                    //If have something to update.
                    if (plc.length > 0) {
                        queryText = `UPDATE public.order SET ${plc.join(", ")} WHERE an_order_id=$1;`;
                        await postgresService.getClient().query(queryText, values);
                        
                    }
                }
            }
        } catch (error) {
            log(error);
            err = NewCommonError(code.ERR_INTERNAL);
        }

        return {err, orders};
    }

    async updateSpOrderParcelIdAndtrackingCode(anOrderId: string, spOrderParcelId: string, trackingCode: string) {
        let err = NewCommonError();
        try {
            const queryText = `UPDATE public.order SET sp_order_parcel_id=$1, tracking_code=$2, fulfillment_status=$3 WHERE an_order_id=$4 RETURNING an_order_id`;
            const values = [spOrderParcelId, trackingCode, AnOrderFulfillmentStatusEnum.PACKED, anOrderId];
            const {rowCount} = await postgresService.getClient().query(queryText, values);
            if (rowCount < 0) {
                log("Not found an_order_id %s", anOrderId);
                err = NewCommonError(code.ERR_INTERNAL);
            }
        } catch (e) {
            log(e);
        }
        return {err};
    }

    async getReferenceNoByAnOrderId(anOrderId: string) {
        let referenceNo: string = "";
        let err = NewCommonError();
        try {
            const queryText = `SELECT reference_no FROM public.order WHERE an_order_id=$1`
            const values = [anOrderId];
            const {rows, rowCount} = await postgresService.getClient().query(queryText, values);
            if (rowCount > 0) {
                referenceNo = rows[0]["reference_no"];
            } else {
                log("Not found an_order_id %s", anOrderId);
                err = NewCommonError(code.ERR_INTERNAL);
            }
        } catch (e) {
            log(e);
        }
        return {referenceNo, err};
    }

    async getReferenceNoBySpOrderParcelId(spOrderParcelId: string) {
        let referenceNo: string = "";
        let err = NewCommonError();
        try {
            const queryText = `SELECT reference_no FROM public.order WHERE sp_order_parcel_id=$1`
            const values = [spOrderParcelId];
            const {rows, rowCount} = await postgresService.getClient().query(queryText, values);
            if (rowCount > 0) {
                referenceNo = rows[0]["reference_no"];
            } else {
                log("Not found sp_order_parcel_id %s", spOrderParcelId);
                err = NewCommonError(code.ERR_INTERNAL);
            }
        } catch (e) {
            log(e);
        }
        return {referenceNo, err};
    }

    async getOrderProductsByAnOrderId(anOrderId: string, userId: string) {
        let orderProducts: Array<any> = [];
        let err = NewCommonError();
        try {
            const queryText = `
SELECT o.user_id, op.an_order_id, op.an_order_product_id, op.an_product_id, op.created_at, p.product_code, p.name, op.quantity
FROM public.order o
INNER JOIN public.order_product op
ON o.an_order_id=op.an_order_id
INNER JOIN public.product p
ON op.an_product_id=p.an_product_id
WHERE o.user_id=$2 AND op.an_order_id=$1`
            const values = [anOrderId, userId];
            const {rows} = await postgresService.getClient().query(queryText, values);
            for (const row of rows) {
                orderProducts.push({
                    userId: row["user_id"],
                    anOrderId: row["an_order_id"],
                    anOrderProductId: row["an_order_product_id"],
                    anProductId: row["an_product_id"],
                    createdAt: row["created_at"],
                    productCode: row["product_code"],
                    name: row["name"],
                    quantity: row["quantity"],
                });
            }
        } catch (e) {
            log(e);
            err = NewCommonError(code.ERR_INTERNAL);
        }

        return {orderProducts, err};
    }

    async getOrderSummary(userId: string, params: DateParam) {
        let orderSum: Array<OrderSummaryModel> = [];
        let err = NewCommonError();
        /**
         * $1=user_id
         * $2=timeZone
         * $3=startDate
         * $4=endDate
         * */
        try {
            const queryText = `
WITH order_summary_cte AS
(
    SELECT o.an_order_id, o.fulfillment_status,
    STRING_AGG(p.product_code || CONCAT('(', op.quantity, ')'), ', ' ORDER BY p.product_code) AS labels, 
    COUNT(op.an_product_id) AS product_count, 
    o.created_at
    FROM public.order AS o
    LEFT JOIN public.order_product AS op
    ON o.an_order_id = op.an_order_id
    LEFT JOIN public.product AS p
    ON op.an_product_id = p.an_product_id
    WHERE o.deleted_at=0
    AND o.created_at at time zone $2 BETWEEN $3 AND $4
    AND o.user_id=$1
    GROUP BY o.an_order_id
)
SELECT labels, COUNT(labels),
COUNT(an_order_id) FILTER (WHERE fulfillment_status = 1) AS packed_count,
COUNT(an_order_id) FILTER (WHERE fulfillment_status = 2) AS not_packed_count,
COUNT(an_order_id) FILTER (WHERE fulfillment_status = 3) AS cancel_count
FROM order_summary_cte
GROUP BY labels
ORDER BY labels`;
            const values = [userId, params.timeZone, params.startDate, params.endDate];
            const {rows} = await postgresService.getClient().query(queryText, values);
            for (const row of rows) {
                if (row['labels'] !== null) {
                    orderSum.push({
                        labels: row["labels"],
                        count: <number>row["count"],
                        packedCount: <number>row['packed_count'],
                        notPackedCount: <number>row['not_packed_count'],
                        cancelCount: <number>row['cancel_count'],
                    });
                }
            }
        } catch (e) {
            err = NewCommonError(code.ERR_INTERNAL);
            log(e);
        }

        return {orderSum, err};
    }

    async getOrderProductSummary(userId: string, params: DateParam) {
        let orderProductSum: Array<OrderProductSummaryModel> = [];
        let err = NewCommonError();

        /**
         * $1=user_id
         * $2=timeZone
         * $3=startDate
         * $4=endDate
         * */
        try {
            const queryText = `
SELECT p.product_code, SUM(op.quantity) FROM public.order o
LEFT JOIN public.order_product op
ON o.an_order_id=op.an_order_id
LEFT JOIN public.product p
ON op.an_product_id=p.an_product_id
WHERE o.deleted_at=0
AND o.created_at at time zone $2 BETWEEN $3 AND $4
AND o.user_id=$1
GROUP BY p.product_code
ORDER BY p.product_code`;

            const values = [userId, params.timeZone, params.startDate, params.endDate];
            const {rows} = await postgresService.getClient().query(queryText, values);
            for (const row of rows) {
                if (row['product_code'] !== null) {
                    orderProductSum.push({
                        productCode: row["product_code"],
                        sum: row["sum"],
                    });
                }
            }
        } catch (e) {
            err = NewCommonError(code.ERR_INTERNAL);
            log(e);
        }

        return {orderProductSum, err};
    }

    async updateOrderFulfillmentStatus(anOrderId: string, fulfillmentStatus: AnOrderFulfillmentStatusEnum) {
        let err = NewCommonError();
        const queryText = `UPDATE public.order SET fulfillment_status=$1 WHERE an_order_id=$2`;
        const values = [fulfillmentStatus, anOrderId];

        try {
            await postgresService.getClient().query(queryText, values);
        } catch (error) {
            log(error);
            err = NewCommonError(code.ERR_INTERNAL);
        }

        return {err};
    }
}

export default OrderDao.getInstance();