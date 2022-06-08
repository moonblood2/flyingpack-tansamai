import {UserIdModel} from "../../common/dto/userId.model";
import {CreatedAtModel} from "../../common/dto/createdAt.model";

export interface ProductModel extends UserIdModel, CreatedAtModel {
    anProductId?: string;
    productCode?: string;
    name?: string;
}