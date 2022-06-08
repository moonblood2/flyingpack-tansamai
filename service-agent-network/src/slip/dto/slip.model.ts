import {CreatedAtModel} from "../../common/dto/createdAt.model";

export interface SlipDto extends CreatedAtModel{
    userId: string;
    message: string;
}