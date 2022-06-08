import {CreatedAtModel} from "../../common/dto/createdAt.model";

export interface SerialNumberModel extends CreatedAtModel {
    anOrderProductSerialNumberId?: string;
    anOrderProductId?: string;
    serialNumberStart: string;
    serialNumberEnd?: string;
}