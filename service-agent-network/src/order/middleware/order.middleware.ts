import express from 'express';
import code, {HttpStatusCode} from "../../common/common.code";
import {NewCommonError} from "../../common/common.error";

class OrderMiddleware {
    async validateOrder(req: express.Request, res: express.Response, next: express.NextFunction) {
        let valid: boolean;
        /**Body is required*/
        valid = !!req.body;
        /**Validate destination*/
        valid = valid && !!req.body['desPhoneNumber'];
        valid = valid && req.body['desPhoneNumber'].length <= 20;
        if (!valid) {
            let err = NewCommonError(code.INPUT_DES_PHONE_NUMBER_INVALID);
            res.status(HttpStatusCode[<number>err.code]).send({...err});
            return;
        }

        /**Validate COD amount*/
        //valid = valid && Number.isInteger(req.body['codAmount']);
        req.body['codAmount'] = <number> req.body['codAmount'];
        valid = valid && req.body['codAmount'] >= 0;
        if (!valid) {
            let err = NewCommonError(code.INPUT_COD_AMOUNT_INVALID);
            res.status(HttpStatusCode[<number>err.code]).send({...err});
            return;
        }
        if (req.body['items']) {
            for (let i = 0; i < req.body['items'].length; i++) {
                req.body['items'][i].quantity = parseInt(req.body['items'][i].quantity);
            }
        }


        /**Bank account is required if COD more than zero*/
        if (req.body['codAmount'] > 0) {
            if (!req.body['bankAccount']) {
                let err = NewCommonError(code.INPUT_BANK_ACCOUNT_INVALID);
                res.status(HttpStatusCode[<number>err.code]).send({...err});
                return;
            }
            if (!req.body['bankAccount']['bank']) {
                let err = NewCommonError(code.INPUT_BANK_INVALID);
                res.status(HttpStatusCode[<number>err.code]).send({...err});
                return;
            }
            if (!req.body['bankAccount']['accountName']) {
                let err = NewCommonError(code.INPUT_ACCOUNT_NAME_INVALID);
                res.status(HttpStatusCode[<number>err.code]).send({...err});
                return;
            }
            if (!req.body['bankAccount']['accountNo']) {
                let err = NewCommonError(code.INPUT_ACCOUNT_NO_INVALID);
                res.status(HttpStatusCode[<number>err.code]).send({...err});
                return;
            }
            if (!req.body['bankAccount']['email']) {
                let err = NewCommonError(code.INPUT_EMAIL_INVALID);
                res.status(HttpStatusCode[<number>err.code]).send({...err});
                return;
            }
            if (!req.body['bankAccount']['fiCode']) {
                let err = NewCommonError(code.INPUT_FI_CODE_INVALID);
                res.status(HttpStatusCode[<number>err.code]).send({...err});
                return;
            }
        }
        //Remove unnecessary char
        req.body['desPhoneNumber'] = req.body['desPhoneNumber'].trim()
        req.body['desAddress'] = req.body['desAddress'].replace(/\n/g, ' ').trim()
        req.body['desSubdistrict'] = req.body['desSubdistrict'].trim()
        req.body['desDistrict'] = req.body['desDistrict'].trim()
        req.body['desProvince'] = req.body['desProvince'].trim()
        req.body['desPostcode'] = req.body['desPostcode'].trim()

        next();
    }
    async validateCancelOrder(req: express.Request, res: express.Response, next: express.NextFunction) {
        let valid: boolean;
        /**Body is required*/
        valid = !!req.body;
        /**Validate referenceNo*/
        valid = valid && !!req.body['referenceNo'];
        if (!valid) {
            let err = NewCommonError(code.INPUT_REFERENCE_NO_INVALID);
            res.status(HttpStatusCode[<number>err.code]).send({...err});
            return;
        }

        next();
    }

}

export default new OrderMiddleware();