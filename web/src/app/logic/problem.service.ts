import { Injectable } from '@angular/core';
import { ApiService } from '../dataaccess';
import { RpcError, RpcProblemSnippet } from '../dataaccess/api';
import { ErrorCode } from '../dataaccess/api.service';
import { PermissionDeniedError, UnauthenticatedError } from './account.service';

export class InvalidProblemListParam extends Error {
  constructor() {
    super('Invalid problem list parameters');
  }
}

@Injectable({
  providedIn: 'root',
})
export class ProblemService {
  constructor(private readonly api: ApiService) {}

  public async getProblemSnippetList(
    offset: number,
    limit: number
  ): Promise<{
    totalProblemCount: number;
    problemSnippetList: RpcProblemSnippet[];
  }> {
    try {
      const response = await this.api.getProblemSnippetList({ offset, limit });
      return {
        totalProblemCount: response.totalProblemCount,
        problemSnippetList: response.problemSnippetList,
      };
    } catch (e) {
      if (!this.api.isRpcError(e)) {
        throw e;
      }

      const apiError = e as RpcError;
      if (apiError.code == ErrorCode.JRPCErrorInvalidParams) {
        throw new InvalidProblemListParam();
      }

      if (apiError.code == ErrorCode.Unauthenticated) {
        throw new UnauthenticatedError();
      }

      if (apiError.code == ErrorCode.PermissionDenied) {
        throw new PermissionDeniedError();
      }

      throw e;
    }
  }
}
