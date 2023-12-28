import { Injectable } from '@angular/core';
import { ApiService } from '../dataaccess';
import {
  RpcError,
  RpcProblem,
  RpcProblemExample,
  RpcProblemSnippet,
} from '../dataaccess/api';
import { ErrorCode } from '../dataaccess/api.service';
import { PermissionDeniedError, UnauthenticatedError } from './account.service';

export class ProblemCreationDisabledError extends Error {
  constructor() {
    super('Problem creation is disabled');
  }
}

export class ProblemUpdateDisabledError extends Error {
  constructor() {
    super('Account update is disabled');
  }
}

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

  public async getProblem(uuid: string): Promise<RpcProblem> {
    try {
      const response = await this.api.getProblem({ uUID: uuid });
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

      if (apiError.code == ErrorCode.NotFound) {
        throw new ProblemNotFoundError();
      }

      throw e;
    }
  }

  public async createProblem(
    displayName: string,
    description: string,
    timeLimitInMillisecond: number,
    memoryLimitInByte: number,
    exampleList: RpcProblemExample[]
  ): Promise<RpcProblem> {
    try {
      const response = await this.api.createProblem({
        displayName,
        description,
        timeLimitInMillisecond,
        memoryLimitInByte,
        exampleList,
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

      if (apiError.code === ErrorCode.Unavailable) {
        throw new ProblemCreationDisabledError();
      }

      throw e;
    }
  }

  public async updateProblem(
    uuid: string,
    displayName: string,
    description: string,
    timeLimitInMillisecond: number,
    memoryLimitInByte: number,
    exampleList: RpcProblemExample[]
  ): Promise<RpcProblem> {
    try {
      const response = await this.api.updateProblem({
        uUID: uuid,
        displayName,
        description,
        timeLimitInMillisecond,
        memoryLimitInByte,
        exampleList,
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

      if (apiError.code == ErrorCode.NotFound) {
        throw new ProblemNotFoundError();
      }

      if (apiError.code === ErrorCode.Unavailable) {
        throw new ProblemUpdateDisabledError();
      }

      throw e;
    }
  }

  public async deleteProblem(uuid: string): Promise<void> {
    try {
      await this.api.deleteProblem({ uUID: uuid });
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

      if (apiError.code == ErrorCode.NotFound) {
        throw new ProblemNotFoundError();
      }

      if (apiError.code === ErrorCode.Unavailable) {
        throw new ProblemUpdateDisabledError();
      }

      throw e;
    }
  }
}
