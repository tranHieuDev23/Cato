import { Pipe, PipeTransform } from '@angular/core';
import { RpcSubmission, RpcSubmissionSnippet } from '../../dataaccess/api';
import {
  SubmissionResult,
  SubmissionStatus,
} from '../../logic/submission.service';

@Pipe({
  name: 'submissionStatusColor',
  standalone: true,
})
export class SubmissionStatusColorPipe implements PipeTransform {
  public transform(value: RpcSubmission | RpcSubmissionSnippet): string {
    if (value.status === SubmissionStatus.Submitted) {
      return 'blue';
    }
    if (value.status === SubmissionStatus.Executing) {
      return 'gold';
    }
    if (value.status === SubmissionStatus.Finished) {
      if (value.result === SubmissionResult.OK) {
        return 'green';
      }
    }
    return 'red';
  }
}
