import productService from '../service/product.service';

import code, {HttpStatusCode} from "../../common/common.code";

import express from 'express';
import httpStatus from "http-status-codes";

import debug from 'debug';

const log: debug.IDebugger = debug('app:product-controller');

class ProductController {
    private static instance: ProductController;

    static getInstance(): ProductController {
        if (!ProductController.instance) {
            ProductController.instance = new ProductController();
        }
        return ProductController.instance;
    }

    async listProducts(req: express.Request, res: express.Response) {
        const userId: string = <string>req.headers["user-id"];
        const {products, err} = await productService.getProducts(userId);

        if (err.code === code.SUCCESS) {
            res.status(httpStatus.OK).send({
                ...err,
                data: products,
            });
        } else {
            res.status(HttpStatusCode[<number>err.code]).send({...err});
        }
    }

    async listProductsByUserId(req: express.Request, res: express.Response) {
        const userId: string = <string>req.params.userId;
        const {products, err} = await productService.getProducts(userId);

        if (err.code === code.SUCCESS) {
            res.status(httpStatus.OK).send({
                ...err,
                data: products,
            });
        } else {
            res.status(HttpStatusCode[<number>err.code]).send({...err});
        }
    }
}

export default ProductController.getInstance();