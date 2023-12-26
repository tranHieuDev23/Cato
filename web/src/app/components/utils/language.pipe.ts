import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'language',
  standalone: true,
})
export class LanguagePipe implements PipeTransform {
  public transform(language: unknown): string {
    if (language === 'c') {
      return 'C';
    }
    if (language === 'cpp') {
      return 'C++';
    }
    if (language === 'java') {
      return 'Java';
    }
    if (language === 'python') {
      return 'Python';
    }
    return 'Unknown';
  }
}
