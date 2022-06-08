export class Product {
    constructor({id, name, price, quantity}) {
        this.id = id || 0
        this.name = name || ""
        this.price = price || 0.0
        this.quantity = quantity || 0
    }
    //clone
    clone() {
        return new Product({
            id: this.id,
            name: this.name,
            price: this.price,
            quantity: this.quantity,
        })
    }
    //validate
    //toPayload
    toPayload() {
        return {
            "id": this.id,
            "quantity": this.quantity,
        }
    }
}