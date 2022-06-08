export interface PaginationResource {
    previousPage: number|null;
    currentPage: number;
    nextPage: number|null;
    firstPage: number;
    lastPage: number;
    isFirstPage: boolean;
    isLastPage: boolean;
    totalPage: number;
    totalItem: number;
}