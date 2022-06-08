export class AnSerialNumber {
    constructor({
                    anOrderProductSerialNumberId,
                    anOrderProductId,
                    serialNumberStart,
                    serialNumberEnd,
                }) {
        this.anOrderProductSerialNumberId = anOrderProductSerialNumberId || "";
        this.anOrderProductId = anOrderProductId || "";
        this.serialNumberStart = serialNumberStart || "";
        this.serialNumberEnd = serialNumberEnd || "";
    }

    clone() {
        return new AnSerialNumber({
            anOrderProductSerialNumberId: this.anOrderProductSerialNumberId,
            anOrderProductId: this.anOrderProductId,
            serialNumberStart: this.serialNumberStart,
            serialNumberEnd: this.serialNumberEnd,
        })
    }
}