import { Pipe, PipeTransform } from '@angular/core';
import { RpcSubmission, RpcSubmissionSnippet } from '../../dataaccess/api';
import {
  SubmissionResult,
  SubmissionStatus,
} from '../../logic/submission.service';

@Pipe({
  name: 'submissionStatus',
  standalone: true,
})
export class SubmissionStatusPipe implements PipeTransform {
  public transform(value: RpcSubmission | RpcSubmissionSnippet): string {
    if (value.status === SubmissionStatus.Submitted) {
      return 'Submitted';
    }
    if (value.status === SubmissionStatus.Executing) {
      return 'Executing';
    }
    if (value.status === SubmissionStatus.Finished) {
      if (value.result === SubmissionResult.OK) {
        return 'Accepted';
      }
      if (value.result === SubmissionResult.CompileError) {
        return 'Compile Error';
      }
      if (value.result === SubmissionResult.RuntimeError) {
        return 'Runtime Error';
      }
      if (value.result === SubmissionResult.TimeLimitExceeded) {
        return 'Time Limit Exceeded';
      }
      if (value.result === SubmissionResult.MemoryLimitExceed) {
        return 'Memory Limit Exceeded';
      }
      if (value.result === SubmissionResult.WrongAnswer) {
        return 'Wrong Answer';
      }
    }
    return 'Unknown';
  }
}
