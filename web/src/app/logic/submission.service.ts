import { Injectable } from '@angular/core';
import { ApiService } from '../dataaccess';
import {
  RpcSubmissionSnippet,
  RpcError,
  RpcSubmission,
} from '../dataaccess/api';
import { ErrorCode } from '../dataaccess/api.service';
import {
  UnauthenticatedError,
  PermissionDeniedError,
  AccountService,
  AccountNotFoundError,
} from './account.service';
import { ProblemNotFoundError } from './problem.service';

export enum SubmissionStatus {
  Submitted = 1,
  Executing = 2,
  Finished = 3,
}

export enum SubmissionResult {
  OK = 1,
  CompileError = 2,
  RuntimeError = 3,
  TimeLimitExceeded = 4,
  MemoryLimitExceed = 5,
  WrongAnswer = 6,
}

export class InvalidSubmissionListParam extends Error {
  constructor() {
    super('Invalid submission list parameters');
  }
}

export class InvalidSubmissionInfo extends Error {
  constructor() {
    super('Invalid submission information');
  }
}

export class SubmissionNotFoundError extends Error {
  constructor() {
    super('Submission not found');
  }
}

@Injectable({
  providedIn: 'root',
})
export class SubmissionService {
  constructor(
    private readonly api: ApiService,
    private readonly accountService: AccountService
  ) {}

  public async getSubmission(id: number): Promise<RpcSubmission> {
    try {
      const response = await this.api.getSubmission({ iD: id });
      return response.submission;
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
        throw new SubmissionNotFoundError();
      }

      throw e;
    }
  }

  public async getSubmissionSnippetList(
    offset: number,
    limit: number
  ): Promise<{
    totalSubmissionCount: number;
    submissionSnippetList: RpcSubmissionSnippet[];
  }> {
    try {
      const response = await this.api.getSubmissionSnippetList({
        offset,
        limit,
      });
      return {
        totalSubmissionCount: response.totalSubmissionCount,
        submissionSnippetList: response.submissionSnippetList,
      };
    } catch (e) {
      if (!this.api.isRpcError(e)) {
        throw e;
      }

      const apiError = e as RpcError;
      if (apiError.code == ErrorCode.JRPCErrorInvalidParams) {
        throw new InvalidSubmissionListParam();
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

  public async getAccountSubmissionSnippetList(
    accountID: number,
    offset: number,
    limit: number
  ): Promise<{
    totalSubmissionCount: number;
    submissionSnippetList: RpcSubmissionSnippet[];
  }> {
    try {
      const response = await this.api.getAccountSubmissionSnippetList({
        accountID: accountID,
        offset,
        limit,
      });
      return {
        totalSubmissionCount: response.totalSubmissionCount,
        submissionSnippetList: response.submissionSnippetList,
      };
    } catch (e) {
      if (!this.api.isRpcError(e)) {
        throw e;
      }

      const apiError = e as RpcError;
      if (apiError.code == ErrorCode.JRPCErrorInvalidParams) {
        throw new InvalidSubmissionListParam();
      }

      if (apiError.code == ErrorCode.Unauthenticated) {
        throw new UnauthenticatedError();
      }

      if (apiError.code == ErrorCode.PermissionDenied) {
        throw new PermissionDeniedError();
      }

      if (apiError.code == ErrorCode.NotFound) {
        throw new AccountNotFoundError();
      }

      throw e;
    }
  }

  public async getSessionAccountSubmissionSnippetList(
    offset: number,
    limit: number
  ): Promise<{
    totalSubmissionCount: number;
    submissionSnippetList: RpcSubmissionSnippet[];
  }> {
    const sessionAccount = await this.accountService.getSessionAccount();
    if (sessionAccount === null) {
      throw new UnauthenticatedError();
    }

    return this.getAccountSubmissionSnippetList(
      sessionAccount.iD,
      offset,
      limit
    );
  }

  public async getProblemSubmissionSnippetList(
    problemID: number,
    offset: number,
    limit: number
  ): Promise<{
    totalSubmissionCount: number;
    submissionSnippetList: RpcSubmissionSnippet[];
  }> {
    try {
      const response = await this.api.getProblemSubmissionSnippetList({
        problemID: problemID,
        offset,
        limit,
      });
      return {
        totalSubmissionCount: response.totalSubmissionCount,
        submissionSnippetList: response.submissionSnippetList,
      };
    } catch (e) {
      if (!this.api.isRpcError(e)) {
        throw e;
      }

      const apiError = e as RpcError;
      if (apiError.code == ErrorCode.JRPCErrorInvalidParams) {
        throw new InvalidSubmissionListParam();
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

      throw e;
    }
  }

  public async getAccountProblemSubmissionSnippetList(
    accountID: number,
    problemID: number,
    offset: number,
    limit: number
  ): Promise<{
    totalSubmissionCount: number;
    submissionSnippetList: RpcSubmissionSnippet[];
  }> {
    try {
      const response = await this.api.getAccountProblemSubmissionSnippetList({
        accountID,
        problemID,
        offset,
        limit,
      });
      return {
        totalSubmissionCount: response.totalSubmissionCount,
        submissionSnippetList: response.submissionSnippetList,
      };
    } catch (e) {
      if (!this.api.isRpcError(e)) {
        throw e;
      }

      const apiError = e as RpcError;
      if (apiError.code == ErrorCode.JRPCErrorInvalidParams) {
        throw new InvalidSubmissionListParam();
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

  public async createSubmission(
    problemID: number,
    content: string,
    language: string
  ): Promise<RpcSubmissionSnippet> {
    try {
      const response = await this.api.createSubmission({
        problemID,
        content,
        language,
      });
      return response.submissionSnippet;
    } catch (e) {
      if (!this.api.isRpcError(e)) {
        throw e;
      }

      const apiError = e as RpcError;
      if (apiError.code == ErrorCode.JRPCErrorInvalidParams) {
        throw new InvalidSubmissionInfo();
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

      throw e;
    }
  }
}
