import { Pipe, PipeTransform } from '@angular/core';

const KB_IN_BYTE = 1024;
const MB_IN_BYTE = KB_IN_BYTE * 1024;
const GB_IN_BYTE = MB_IN_BYTE * 1024;

@Pipe({
  name: 'memory',
  standalone: true,
})
export class MemoryPipe implements PipeTransform {
  transform(value: number): string {
    if (value % GB_IN_BYTE === 0) {
      return `${value / GB_IN_BYTE} GB`;
    }
    if (value % MB_IN_BYTE === 0) {
      return `${value / MB_IN_BYTE} MB`;
    }
    if (value % KB_IN_BYTE === 0) {
      return `${value / KB_IN_BYTE} kB`;
    }

    return `${value} byte`;
  }
}
