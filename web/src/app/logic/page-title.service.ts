import { EventEmitter, Injectable } from '@angular/core';
import { Title } from '@angular/platform-browser';

@Injectable({
  providedIn: 'root',
})
export class PageTitleService {
  public titleChanged = new EventEmitter<string>();

  constructor(private readonly title: Title) {}

  public setTitle(title: string): void {
    this.title.setTitle(title);
    this.titleChanged.emit(title);
  }
}
