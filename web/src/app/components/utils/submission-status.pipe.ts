import { Pipe, PipeTransform } from '@angular/core';
import { RpcSubmission, RpcSubmissionSnippet } from '../../dataaccess/api';

@Pipe({
  name: 'submissionStatus',
  standalone: true,
})
export class SubmissionStatusPipe implements PipeTransform {
  public transform(value: RpcSubmission | RpcSubmissionSnippet): string {
    if (value.status === 1) {
      return 'Submitted';
    }
    if (value.status === 2) {
      return 'Executing';
    }
    if (value.status === 3) {
      if (value.result === 1) {
        return 'Accepted';
      }
      if (value.result === 2) {
        return 'Compile Error';
      }
      if (value.result === 3) {
        return 'Runtime Error';
      }
      if (value.result === 4) {
        return 'Time Limit Exceeded';
      }
      if (value.result === 5) {
        return 'Memory Limit Exceeded';
      }
      if (value.result === 6) {
        return 'Wrong Answer';
      }
    }
    return 'Unknown';
  }
}
