export class SpOrderParcelShippopFlash {
    constructor({sortCode, dstCode, sortingLineCode}) {
        this.sortCode = sortCode || "";
        this.dstCode = dstCode || "";
        this.sortingLineCode = sortingLineCode || "";
    }
}