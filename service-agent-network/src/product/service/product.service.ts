import {ProductModel} from "../dto/product.model";
import productDao from "../dao/product.dao";

class ProductService {
    private static instance: ProductService;

    private constructor() {
    }

    static getInstance(): ProductService {
        if (!ProductService.instance) {
            ProductService.instance = new ProductService();
        }
        return ProductService.instance;
    }

    async doProductsExist(items: any, userId: string) {
        //Convert items to products.
        let products: Array<ProductModel> = [];
        for (const e of items) {
            products.push({
                userId: userId,
                productCode: e.productCode,
            });
        }
        return await productDao.doProductsExist(products);
    }

    /**
     * getProductsByProductCodes, it's just call getProductsByProductCodes() from productDao.
     * @param productCodes
     * @param userId
     */
    async getProductsByProductCodes(productCodes: Array<string>, userId: string) {
        return await productDao.getProductsByProductCodes(productCodes, userId);
    }

    async getProducts(userId: string) {
        return await productDao.getProducts(userId);
    }
}

export default ProductService.getInstance();