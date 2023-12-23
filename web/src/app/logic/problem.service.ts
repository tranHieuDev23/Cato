import { Injectable } from '@angular/core';
import { ApiService } from '../dataaccess';
import { RpcError, RpcProblem, RpcProblemSnippet } from '../dataaccess/api';
import { ErrorCode } from '../dataaccess/api.service';
import { PermissionDeniedError, UnauthenticatedError } from './account.service';

export class InvalidProblemListParam extends Error {
  constructor() {
    super('Invalid problem list parameters');
  }
}

export class InvalidProblemInfo extends Error {
  constructor() {
    super('Invalid problem info');
  }
}

export class ProblemNotFoundError extends Error {
  constructor() {
    super('Problem not found');
  }
}

@Injectable({
  providedIn: 'root',
})
export class ProblemService {
  constructor(private readonly api: ApiService) {}

  public isValidDisplayName(
    displayName: string
  ): { [k: string]: boolean } | null {
    if (displayName.length > 256) {
      return { error: true, maxLength: true };
    }
    return null;
  }

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

  public async getProblem(id: number): Promise<RpcProblem> {
    try {
      const response = await this.api.getProblem({ iD: id });
      return response.problem;
    } catch (e) {
      if (!this.api.isRpcError(e)) {
        throw e;
      }

      const apiError = e as RpcError;
      if (apiError.code == ErrorCode.Unauthenticated) {
        throw new UnauthenticatedError();
      }

      if (apiError.code == ErrorCode.PermissionDenied) {
        throw new PermissionDeniedError();
      }

      throw e;
    }
  }

  public async createProblem(
    displayName: string,
    description: string,
    timeLimitInMillisecond: number,
    memoryLimitInByte: number
  ): Promise<RpcProblem> {
    try {
      const response = await this.api.createProblem({
        displayName,
        description,
        timeLimitInMillisecond,
        memoryLimitInByte,
        exampleList: [],
      });
      return response.problem;
    } catch (e) {
      if (!this.api.isRpcError(e)) {
        throw e;
      }

      const apiError = e as RpcError;
      if (apiError.code == ErrorCode.JRPCErrorInvalidParams) {
        throw new InvalidProblemInfo();
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
