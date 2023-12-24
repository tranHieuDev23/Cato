import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'language',
  standalone: true,
})
export class LanguagePipe implements PipeTransform {
  public transform(language: unknown): string {
    if (language === 'cpp') {
      return 'C++';
    }
    if (language === 'java') {
      return 'java';
    }
    if (language === 'python') {
      return 'Python';
    }
    return 'Unknown';
  }
}
