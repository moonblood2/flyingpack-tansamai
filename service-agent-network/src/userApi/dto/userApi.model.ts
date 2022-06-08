import {CreatedAtModel} from "../../common/dto/createdAt.model";

export interface UserApiDto extends CreatedAtModel{
    userId: string;
    apiKeyPrefix: string;
    apiKeyHash: string;
    urlWebHookOrder: string;
}