import postgresService from "../../common/service/postgres.service";
import {ProductModel} from "../dto/product.model";
import {CommonError, NewCommonError} from "../../common/common.error";
import code from "../../common/common.code";

import debug from 'debug';

const log: debug.IDebugger = debug('app:product-dao');

class ProductDao {
    private static instance: ProductDao;

    private constructor() {
    }

    static getInstance(): ProductDao {
        if (!ProductDao.instance) {
            ProductDao.instance = new ProductDao();
        }
        return ProductDao.instance;
    }

    private static makeUniqueProduct(products: Array<ProductModel>): Array<ProductModel> {
        //Make unique product.
        let uniqueProduct: Array<ProductModel> = [];
        let map = new Map();
        for (const p of products) {
            if (map.get(p.productCode) === undefined) {
                map.set(`${p.userId},${p.productCode}`, true);
                uniqueProduct.push(p);
            }
        }
        return uniqueProduct;
    }

    private static makeUniqueProductCodes(productCodes: Array<string>): Array<string> {
        let uniqueProductCodes: Array<string> = [];
        let map: Map<string, boolean> = new Map();
        for (const p of productCodes) {
            if (map.get(p) === undefined) {
                map.set(p, true);
                uniqueProductCodes.push(p);
            }
        }
        return uniqueProductCodes;
    }

    async doProductsExist(products: Array<ProductModel>) {
        //Make unique product.
        let uniqueProduct = ProductDao.makeUniqueProduct(products);
        let result: Map<string, boolean> = new Map();
        let err: CommonError = NewCommonError();
        try {
            await postgresService.getClient().query('BEGIN');
            for (const p of uniqueProduct) {
                const queryText = `SELECT EXISTS (SELECT 1 FROM public.product p WHERE p.product_code=$1 AND user_id=$2 AND deleted_at=0)`;
                const values = [p.productCode, p.userId];
                const {rows} = await postgresService.getClient().query(queryText, values);
                result.set(<string>p.productCode, <boolean>rows[0]["exists"]);
            }
            await postgresService.getClient().query('COMMIT');
        } catch (e) {
            await postgresService.getClient().query('ROLLBACK');
            log(e);
            err = NewCommonError(code.ERR_INTERNAL);
        }
        return {result, err};
    }

    /**
     * getProducts return products from database sorted by an_product_id.
     * @param productCodes
     * @param userId
     */
    async getProductsByProductCodes(productCodes: Array<string>, userId: string) {
        /**Prepared statement**/
        let queryText = `SELECT * FROM public.product WHERE deleted_at=0 AND user_id=$1 `;
        /**$1 is user_id, so next is 2.*/
        let n: number = 2;
        let params: Array<string> = [];
        for (let i = 0; i < productCodes.length; i++, n++) {
            params.push(`$${n}`);
        }
        queryText += `AND product_code IN (${params.join(", ")}) ORDER BY product_code`;
        const values = [userId, ...productCodes];

        let products: Array<ProductModel> = [];
        let err: CommonError = NewCommonError();
        try {
            const {rows} = await postgresService.getClient().query(queryText, values);
            for (const row of rows) {
                products.push({
                    anProductId: row["an_product_id"],
                    userId: row["user_id"],
                    productCode: row["product_code"],
                    name: row["name"],
                    createdAt: row["created_at"],
                    deletedAt: row["deleted_at"],
                });
            }
        } catch (e) {
            log(e);
            err = NewCommonError(code.ERR_INTERNAL);
        }
        return {products, err};
    }

    async getProducts(userId: string) {
        let products: Array<ProductModel> = [];
        let err = NewCommonError();

        const queryText = `SELECT * FROM public.product WHERE deleted_at=0 AND user_id=$1`;
        const values = [userId];

        try {
            const {rows} = await postgresService.getClient().query(queryText, values);
            for (const row of rows) {
                products.push({
                    anProductId: row["an_product_id"],
                    productCode: row["product_code"],
                    name: row["name"],
                    createdAt: row["created_at"],
                });
            }
        } catch (e) {
            log(e);
            err = NewCommonError(code.ERR_INTERNAL);
        }

        return {products, err};
    }
}

export default ProductDao.getInstance();