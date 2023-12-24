import { Pipe, PipeTransform } from '@angular/core';
import { UnitService } from '../../logic/unit.service';

@Pipe({
  name: 'memory',
  standalone: true,
})
export class MemoryPipe implements PipeTransform {
  constructor(private readonly unitService: UnitService) {}

  transform(input: number): string {
    const { value, unit } = this.unitService.memoryLimitToValueAndUnit(input);
    return `${value} ${unit}`;
  }
}
