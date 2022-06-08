import Joi from 'joi'
import {ContactInfo} from './ContactInfo'
import {validateObject, VError} from "@/utils/joi-validation";

export class Sender extends ContactInfo {
    constructor({
                    id,
                    senderType,
                    nationalIdNumber,
                    taxIdNumber,
                    passportNumber,
                    birthDate,
                    name,
                    phoneNumber,
                    address,
                    district,
                    state,
                    province,
                    postcode
                }) {
        super({
            name,
            phoneNumber,
            address,
            district,
            state,
            province,
            postcode
        });

        this.id = id || 0
        this.senderType = senderType || 1
        this.nationalIdNumber = nationalIdNumber || ""
        this.taxIdNumber = taxIdNumber || ""
        this.passportNumber = passportNumber || ""
        this.birthDate = birthDate || "1999-07-26"

        this.valid = false
        this.error = {
            ...this.error,
            senderType: new VError(),
            nationalIdNumber: new VError(),
            taxIdNumber: new VError(),
            passportNumber: new VError(),
            birthDate: new VError(),
        }
    }

    clone() {
        return new Sender({
            id: this.id,
            senderType: this.senderType,
            nationalIdNumber: this.nationalIdNumber,
            taxIdNumber: this.taxIdNumber,
            passportNumber: this.passportNumber,
            birthDate: this.birthDate,
            name: this.nationalIdNumber,
            phoneNumber: this.phoneNumber,
            address: this.address,
            district: this.district,
            state: this.state,
            province: this.$props,
            postcode: this.postcode,
        });
    }

    validate() {
        const contactValidate = super.validate();
        const {valid, error} = validateObject({
            senderType: Joi.number().min(1).max(2).required(),
            nationalIdNumber: Joi.string().length(13).required(),
            birthDate: Joi.date().iso().required(),
        }, {
            senderType: this.senderType,
            nationalIdNumber: this.nationalIdNumber,
            birthDate: this.birthDate,
        })
        this.valid = contactValidate.valid && valid
        this.error = {...error, ...contactValidate.error}
        return {valid: this.valid, error: this.error}
    }

    toPayload() {
        return {
            "sender_type": this.senderType,
            "national_id_number": this.nationalIdNumber,
            "passport_number": this.passportNumber,
            "tax_id_number": this.taxIdNumber,
            "birth_date": this.birthDate,
            "name": this.name,
            "phone_number": this.phoneNumber,
            "address": this.address,
            "district": this.district,
            "state": this.state,
            "province": this.province,
            "postcode": this.postcode,
        }
    }
}