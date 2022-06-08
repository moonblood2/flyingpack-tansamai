import serviceConfig from "../../config/service.config";
import header from "../common.header";

import express from "express";
import httpStatus from "http-status-codes";
import jwt from "jsonwebtoken";

import debug from "debug";

const log: debug.IDebugger = debug("common-authen-middleware:");

class AuthenMiddleware {
    private static instance: AuthenMiddleware;

    static getInstance(): AuthenMiddleware {
        if (!AuthenMiddleware.instance) {
            AuthenMiddleware.instance = new AuthenMiddleware();
        }
        return AuthenMiddleware.instance;
    }

    async verifyToken(req: express.Request, res: express.Response, next: express.NextFunction) {
        let token: string = <string> req.headers["authorization"];
        if (token) {
            token = token.replace("Bearer ", "");
            try {
                const decoded: any = jwt.verify(token, serviceConfig.JWT_SIGNING_KEY, {
                    algorithms: ["HS256"],
                });
                if (decoded.uid) {
                    req.headers[header.USER_ID] = decoded.uid;
                    req.headers[header.TOKEN] = token;
                    next();
                } else {
                    res.status(httpStatus.UNAUTHORIZED).send({
                        "status": false,
                        "message": "invalid token"
                    })
                }
            } catch (error) {
                res.status(httpStatus.UNAUTHORIZED).send({
                    "status": false,
                    "message": error.message
                })
            }
        } else {
            res.status(httpStatus.UNAUTHORIZED).send({
                "status": false,
                "message": "invalid token"
            })
        }
    }
}

export default AuthenMiddleware.getInstance();