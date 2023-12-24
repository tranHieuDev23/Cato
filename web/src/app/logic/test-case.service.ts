import { Injectable } from '@angular/core';
import { RpcTestCase, RpcError, RpcTestCaseSnippet } from '../dataaccess/api';
import { ApiService, ErrorCode } from '../dataaccess/api.service';
import { UnauthenticatedError, PermissionDeniedError } from './account.service';
import { ProblemNotFoundError } from './problem.service';

export class InvalidTestCaseListParam extends Error {
  constructor() {
    super('Invalid Test case list parameters');
  }
}

export class InvalidTestCaseInfo extends Error {
  constructor() {
    super('Invalid Test case information');
  }
}

export class TestCaseNotFoundError extends Error {
  constructor() {
    super('Test case not found');
  }
}

@Injectable({
  providedIn: 'root',
})
export class TestCaseService {
  constructor(private readonly api: ApiService) {}

  public async getTestCase(id: number): Promise<RpcTestCase> {
    try {
      const response = await this.api.getTestCase({ iD: id });
      return response.testCase;
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
        throw new TestCaseNotFoundError();
      }

      throw e;
    }
  }

  public async getProblemTestCaseSnippetList(
    problemID: number,
    offset: number,
    limit: number
  ): Promise<{
    totalTestCaseCount: number;
    testCaseSnippetList: RpcTestCaseSnippet[];
  }> {
    try {
      const response = await this.api.getProblemTestCaseSnippetList({
        problemID: problemID,
        offset,
        limit,
      });
      return {
        totalTestCaseCount: response.totalTestCaseCount,
        testCaseSnippetList: response.testCaseSnippetList,
      };
    } catch (e) {
      if (!this.api.isRpcError(e)) {
        throw e;
      }

      const apiError = e as RpcError;
      if (apiError.code == ErrorCode.JRPCErrorInvalidParams) {
        throw new InvalidTestCaseListParam();
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

  public async createTestCase(
    problemID: number,
    input: string,
    output: string,
    isHidden: boolean
  ): Promise<RpcTestCaseSnippet> {
    try {
      const response = await this.api.createTestCase({
        problemID,
        input,
        output,
        isHidden,
      });
      return response.testCaseSnippet;
    } catch (e) {
      if (!this.api.isRpcError(e)) {
        throw e;
      }

      const apiError = e as RpcError;
      if (apiError.code == ErrorCode.JRPCErrorInvalidParams) {
        throw new InvalidTestCaseInfo();
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

  public async updateTestCase(
    id: number,
    input: string,
    output: string,
    isHidden: boolean
  ): Promise<RpcTestCaseSnippet> {
    try {
      const response = await this.api.updateTestCase({
        iD: id,
        input,
        output,
        isHidden,
      });
      return response.testCaseSnippet;
    } catch (e) {
      if (!this.api.isRpcError(e)) {
        throw e;
      }

      const apiError = e as RpcError;
      if (apiError.code == ErrorCode.JRPCErrorInvalidParams) {
        throw new InvalidTestCaseInfo();
      }

      if (apiError.code == ErrorCode.Unauthenticated) {
        throw new UnauthenticatedError();
      }

      if (apiError.code == ErrorCode.PermissionDenied) {
        throw new PermissionDeniedError();
      }

      if (apiError.code == ErrorCode.NotFound) {
        throw new TestCaseNotFoundError();
      }

      throw e;
    }
  }

  public async deleteTestCase(id: number): Promise<void> {
    try {
      await this.api.deleteTestCase({ iD: id });
    } catch (e) {
      if (!this.api.isRpcError(e)) {
        throw e;
      }

      const apiError = e as RpcError;
      if (apiError.code == ErrorCode.JRPCErrorInvalidParams) {
        throw new InvalidTestCaseInfo();
      }

      if (apiError.code == ErrorCode.Unauthenticated) {
        throw new UnauthenticatedError();
      }

      if (apiError.code == ErrorCode.PermissionDenied) {
        throw new PermissionDeniedError();
      }

      if (apiError.code == ErrorCode.NotFound) {
        throw new TestCaseNotFoundError();
      }

      throw e;
    }
  }
}
