import { Pipe, PipeTransform } from '@angular/core';
import { Role } from '../../logic/account.service';

@Pipe({
  name: 'role',
  standalone: true,
})
export class RolePipe implements PipeTransform {
  transform(value: unknown): string {
    if (value === Role.Admin) {
      return 'Admin';
    }
    if (value === Role.ProblemSetter) {
      return 'Problem Setter';
    }
    if (value === Role.Contestant) {
      return 'Contestant';
    }
    if (value === Role.Worker) {
      return 'Worker';
    }
    return '';
  }
}
