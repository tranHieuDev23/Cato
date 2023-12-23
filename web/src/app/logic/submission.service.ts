import { Injectable } from '@angular/core';
import { ApiService } from '../dataaccess';
import { RpcSubmissionSnippet, RpcError } from '../dataaccess/api';
import { ErrorCode } from '../dataaccess/api.service';
import {
  UnauthenticatedError,
  PermissionDeniedError,
  AccountService,
  AccountNotFoundError,
} from './account.service';

export class InvalidSubmissionListParam extends Error {
  constructor() {
    super('Invalid submission list parameters');
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
}
