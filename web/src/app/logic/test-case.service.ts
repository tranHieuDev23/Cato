import { Injectable } from '@angular/core';
import { RpcTestCase, RpcError, RpcTestCaseSnippet } from '../dataaccess/api';
import { ApiService, ErrorCode } from '../dataaccess/api.service';
import { UnauthenticatedError, PermissionDeniedError } from './account.service';
import {
  ProblemNotFoundError,
  ProblemUpdateDisabledError,
} from './problem.service';
import { Encoder } from 'bazinga64';

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

  public async getTestCase(uuid: string): Promise<RpcTestCase> {
    try {
      const response = await this.api.getTestCase({ uUID: uuid });
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
    problemUUID: string,
    offset: number,
    limit: number
  ): Promise<{
    totalTestCaseCount: number;
    testCaseSnippetList: RpcTestCaseSnippet[];
  }> {
    try {
      const response = await this.api.getProblemTestCaseSnippetList({
        problemUUID,
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
    problemUUID: string,
    input: string,
    output: string,
    isHidden: boolean
  ): Promise<RpcTestCaseSnippet> {
    try {
      const response = await this.api.createTestCase({
        problemUUID,
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

      if (apiError.code === ErrorCode.Unavailable) {
        throw new ProblemUpdateDisabledError();
      }

      throw e;
    }
  }

  public async createTestCaseList(
    problemUUID: string,
    zippedTestData: ArrayBuffer
  ): Promise<void> {
    try {
      await this.api.createTestCaseList({
        problemUUID,
        zippedTestData: Encoder.toBase64(zippedTestData).asString,
      });
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

      if (apiError.code === ErrorCode.Unavailable) {
        throw new ProblemUpdateDisabledError();
      }

      throw e;
    }
  }

  public async updateTestCase(
    uuid: string,
    input: string,
    output: string,
    isHidden: boolean
  ): Promise<RpcTestCaseSnippet> {
    try {
      const response = await this.api.updateTestCase({
        uUID: uuid,
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

      if (apiError.code === ErrorCode.Unavailable) {
        throw new ProblemUpdateDisabledError();
      }

      throw e;
    }
  }

  public async deleteTestCase(uuid: string): Promise<void> {
    try {
      await this.api.deleteTestCase({ uUID: uuid });
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

      if (apiError.code === ErrorCode.Unavailable) {
        throw new ProblemUpdateDisabledError();
      }

      throw e;
    }
  }
}
