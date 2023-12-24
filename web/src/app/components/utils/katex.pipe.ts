import { Pipe, PipeTransform } from '@angular/core';
import renderMathInElement from 'katex/contrib/auto-render';

@Pipe({
  name: 'katex',
  standalone: true,
})
export class KatexPipe implements PipeTransform {
  transform(value: string): string {
    const valueElement = document.createElement('div');
    valueElement.innerHTML = value;
    renderMathInElement(valueElement, { throwOnError: false });
    return valueElement.innerHTML;
  }
}
