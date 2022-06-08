import commonCode, {Message} from "./common.code";

export interface CommonError {
    code?: number;
    message?: string;
}

export const NewCommonError = (code: number = commonCode.SUCCESS, message: string = Message[code]): CommonError => {
    return {
        code: code,
        message: message,
    }
}