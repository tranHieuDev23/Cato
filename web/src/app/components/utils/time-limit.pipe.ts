import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'timeLimit',
  standalone: true,
})
export class TimeLimitPipe implements PipeTransform {
  transform(value: number): string {
    if (value % 1000 === 0) {
      if (value === 1000) {
        return `1 second`;
      }

      return `${value / 1000} seconds`;
    }

    return `${value} ms`;
  }
}
