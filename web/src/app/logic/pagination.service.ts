import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root',
})
export class PaginationService {
  constructor() {}

  public getPageOffset(pageIndex: number, pageSize: number): number {
    return (pageIndex - 1) * pageSize;
  }

  public getLastPageIndex(totalItemCount: number, pageSize: number): number {
    return Math.ceil(totalItemCount / pageSize);
  }
}
