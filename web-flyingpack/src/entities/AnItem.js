export class AnItem {
    constructor({
                    anProductId,
                    productCode,
                    name,
                    imgUrl,
                    quantity,
                    serialNumbers,
                }) {
        this.anProductId = anProductId || "";
        this.name = name || "";
        this.imgUrl = imgUrl || "";
        this.quantity = quantity || 0;
        this.productCode = productCode || "";
        this.serialNumbers = serialNumbers || [];
    }

    clone() {
        return new AnItem({
            anProductId: this.anProductId,
            name: this.name,
            imgUrl: this.imgUrl,
            quantity: this.quantity,
            productCode: this.productCode,
            serialNumbers: [...this.serialNumbers],
        })
    }
}