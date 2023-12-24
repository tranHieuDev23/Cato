import { Pipe, PipeTransform } from '@angular/core';
import { UnitService } from '../../logic/unit.service';

@Pipe({
  name: 'timeLimit',
  standalone: true,
})
export class TimeLimitPipe implements PipeTransform {
  constructor(private readonly unitService: UnitService) {}

  transform(input: number): string {
    const { value, unit } = this.unitService.timeLimitToValueAndUnit(input);
    return `${value} ${unit}`;
  }
}
