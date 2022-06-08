import Joi from "joi";

// validateObject validates all field in object by set abortEarly to false.
export const validateObject = (scheme, value) => {
    let valid = true
    let newError = {}
    for (const key in scheme) {
        newError[key] = new VError(true, "")
    }
    const { error } = Joi.object(scheme).validate(value, {abortEarly: false})
    if (error !== undefined) {
        //Mapping error from Joi to each property of original object.
        for (const detail of error.details) {
            newError[detail.context.key] = new VError(false, detail.message)
            valid = false
        }
    }
    return {valid: valid, error: newError}
}

export class VError {
    constructor(valid, message) {
        this.valid = valid === true || valid === false ? valid : null
        this.message = message || ""
    }
}