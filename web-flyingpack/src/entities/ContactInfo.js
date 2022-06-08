import Joi from 'joi'
import {validateObject, VError} from "@/utils/joi-validation";

export class ContactInfo {
    constructor({name, phoneNumber, address, district, state, province, postcode}) {
        this.name = name || ""
        this.phoneNumber = phoneNumber || ""
        this.address = address || ""
        this.district = district || "" //District is sub unit of State.
        this.state = state || ""//State is sub unit of Province.
        this.province = province || ""
        this.postcode = postcode || ""

        this.valid = false
        this.error = {
            name: new VError(),
            phoneNumber: new VError(),
            address: new VError(),
            district: new VError(),
            state: new VError(),
            province: new VError(),
            postcode: new VError(),
        }
    }

    clone() {
        return new ContactInfo({
            name: this.name,
            phoneNumber: this.phoneNumber,
            address: this.address,
            district: this.district,
            state: this.state,
            province: this.province,
            postcode: this.postcode,
        })
    }

    validate() {
        const {valid, error} = validateObject({
            name: Joi.string().min(5).messages({
                'string.empty': "ห้ามเป็นค่าว่าง",
                'string.min': "อย่างน้อย 5 ตัว"
            }),
            phoneNumber: Joi.string().min(9).max(10).required(),
            address: Joi.string(),
            district: Joi.string(),
            state: Joi.string(),
            province: Joi.string(),
            postcode: Joi.string(),
        }, {
            name: this.name,
            phoneNumber: this.phoneNumber,
            address: this.address,
            district: this.district,
            state: this.district,
            province: this.province,
            postcode: this.postcode,
        })
        this.valid = valid
        this.error = error
        return {valid: this.valid, error: this.error}
    }

    toPayload() {
        return {
            "name": this.name,
            "address": this.address,
            "district": this.district,
            "state": this.state,
            "province": this.province,
            "postcode": this.postcode,
            "phone_number": this.phoneNumber
        }
    }
}