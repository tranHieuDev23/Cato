import { Injectable } from '@angular/core';
import {
  Configuration,
  DefaultApi,
  RpcCreateAccountRequest,
  RpcCreateAccountResponse,
  RpcCreateProblemRequest,
  RpcCreateProblemResponse,
  RpcCreateSessionRequest,
  RpcCreateSessionResponse,
  RpcCreateSubmissionRequest,
  RpcCreateSubmissionResponse,
  RpcCreateTestCaseListRequest,
  RpcCreateTestCaseRequest,
  RpcCreateTestCaseResponse,
  RpcDeleteProblemRequest,
  RpcDeleteSubmissionRequest,
  RpcDeleteTestCaseRequest,
  RpcGetAccountListRequest,
  RpcGetAccountListResponse,
  RpcGetAccountProblemSnippetListRequest,
  RpcGetAccountProblemSnippetListResponse,
  RpcGetAccountProblemSubmissionSnippetListRequest,
  RpcGetAccountProblemSubmissionSnippetListResponse,
  RpcGetAccountRequest,
  RpcGetAccountResponse,
  RpcGetAccountSubmissionSnippetListRequest,
  RpcGetAccountSubmissionSnippetListResponse,
  RpcGetProblemRequest,
  RpcGetProblemResponse,
  RpcGetProblemSnippetListRequest,
  RpcGetProblemSnippetListResponse,
  RpcGetProblemSubmissionSnippetListRequest,
  RpcGetProblemSubmissionSnippetListResponse,
  RpcGetProblemTestCaseSnippetListRequest,
  RpcGetProblemTestCaseSnippetListResponse,
  RpcGetServerInfoResponse,
  RpcGetSessionResponse,
  RpcGetSubmissionRequest,
  RpcGetSubmissionResponse,
  RpcGetSubmissionSnippetListRequest,
  RpcGetSubmissionSnippetListResponse,
  RpcGetTestCaseRequest,
  RpcGetTestCaseResponse,
  RpcSetting,
  RpcUpdateAccountRequest,
  RpcUpdateAccountResponse,
  RpcUpdateProblemRequest,
  RpcUpdateProblemResponse,
  RpcUpdateSettingResponse,
  RpcUpdateTestCaseRequest,
  RpcUpdateTestCaseResponse,
  instanceOfRpcError,
} from './api';

const jsonRPCVersion = '2.0';
const clientID = 'cato-judge';

export enum ErrorCode {
  OK = 1,
  Canceled = 2,
  Unknown = 3,
  InvalidArgument = 4,
  DeadlineExceeded = 5,
  NotFound = 6,
  AlreadyExists = 7,
  PermissionDenied = 8,
  ResourceExhausted = 9,
  FailedPrecondition = 10,
  Aborted = 11,
  OutOfRange = 12,
  Unimplemented = 13,
  Internal = 14,
  Unavailable = 15,
  DataLoss = 16,
  Unauthenticated = 17,

  JRPCErrorInvalidParams = -32602,
  JRPCErrorInternal = -32602,
}

@Injectable({
  providedIn: 'root',
})
export class ApiService {
  private readonly api: DefaultApi;

  constructor() {
    this.api = new DefaultApi(
      new Configuration({
        basePath: '/api',
      })
    );
  }

  public async getServerInfo(): Promise<RpcGetServerInfoResponse> {
    const { error, result } = await this.api.getServerInfo({
      requestBodyOfTheGetServerInfoMethod: {
        jsonrpc: jsonRPCVersion,
        id: clientID,
        method: 'get_server_info',
      },
    });
    if (error) {
      throw error;
    }

    if (!result) {
      throw new Error('No response received');
    }

    return result;
  }

  public async createAccount(
    request: RpcCreateAccountRequest
  ): Promise<RpcCreateAccountResponse> {
    const { error, result } = await this.api.createAccount({
      requestBodyOfTheCreateAccountMethod: {
        jsonrpc: jsonRPCVersion,
        id: clientID,
        method: 'create_account',
        params: request,
      },
    });
    if (error) {
      throw error;
    }

    if (!result) {
      throw new Error('No response received');
    }

    return result;
  }

  public async getAccountList(
    request: RpcGetAccountListRequest
  ): Promise<RpcGetAccountListResponse> {
    const { error, result } = await this.api.getAccountList({
      requestBodyOfTheGetAccountListMethod: {
        jsonrpc: jsonRPCVersion,
        id: clientID,
        method: 'get_account_list',
        params: request,
      },
    });
    if (error) {
      throw error;
    }

    if (!result) {
      throw new Error('No response received');
    }

    return result;
  }

  public async getAccount(
    request: RpcGetAccountRequest
  ): Promise<RpcGetAccountResponse> {
    const { error, result } = await this.api.getAccount({
      requestBodyOfTheGetAccountMethod: {
        jsonrpc: jsonRPCVersion,
        id: clientID,
        method: 'get_account',
        params: request,
      },
    });
    if (error) {
      throw error;
    }

    if (!result) {
      throw new Error('No response received');
    }

    return result;
  }

  public async updateAccount(
    request: RpcUpdateAccountRequest
  ): Promise<RpcUpdateAccountResponse> {
    const { error, result } = await this.api.updateAccount({
      requestBodyOfTheUpdateAccountMethod: {
        jsonrpc: jsonRPCVersion,
        id: clientID,
        method: 'update_account',
        params: request,
      },
    });
    if (error) {
      throw error;
    }

    if (!result) {
      throw new Error('No response received');
    }

    return result;
  }

  public async createSession(
    request: RpcCreateSessionRequest
  ): Promise<RpcCreateSessionResponse> {
    const { error, result } = await this.api.createSession({
      requestBodyOfTheCreateSessionMethod: {
        jsonrpc: jsonRPCVersion,
        id: clientID,
        method: 'create_session',
        params: request,
      },
    });
    if (error) {
      throw error;
    }

    if (!result) {
      throw new Error('No response received');
    }

    return result;
  }

  public async deleteSession(): Promise<void> {
    const { error, result } = await this.api.deleteSession({
      requestBodyOfTheDeleteSessionMethod: {
        jsonrpc: jsonRPCVersion,
        id: clientID,
        method: 'delete_session',
        params: {},
      },
    });
    if (error) {
      throw error;
    }

    if (!result) {
      throw new Error('No response received');
    }
  }

  public async createProblem(
    request: RpcCreateProblemRequest
  ): Promise<RpcCreateProblemResponse> {
    const { error, result } = await this.api.createProblem({
      requestBodyOfTheCreateProblemMethod: {
        jsonrpc: jsonRPCVersion,
        id: clientID,
        method: 'create_problem',
        params: request,
      },
    });
    if (error) {
      throw error;
    }

    if (!result) {
      throw new Error('No response received');
    }

    return result;
  }

  public async getProblemSnippetList(
    request: RpcGetProblemSnippetListRequest
  ): Promise<RpcGetProblemSnippetListResponse> {
    const { error, result } = await this.api.getProblemSnippetList({
      requestBodyOfTheGetProblemSnippetListMethod: {
        jsonrpc: jsonRPCVersion,
        id: clientID,
        method: 'get_problem_snippet_list',
        params: request,
      },
    });
    if (error) {
      throw error;
    }

    if (!result) {
      throw new Error('No response received');
    }

    return result;
  }

  public async getProblem(
    request: RpcGetProblemRequest
  ): Promise<RpcGetProblemResponse> {
    const { error, result } = await this.api.getProblem({
      requestBodyOfTheGetProblemMethod: {
        jsonrpc: jsonRPCVersion,
        id: clientID,
        method: 'get_problem',
        params: request,
      },
    });
    if (error) {
      throw error;
    }

    if (!result) {
      throw new Error('No response received');
    }

    return result;
  }

  public async updateProblem(
    request: RpcUpdateProblemRequest
  ): Promise<RpcUpdateProblemResponse> {
    const { error, result } = await this.api.updateProblem({
      requestBodyOfTheUpdateProblemMethod: {
        jsonrpc: jsonRPCVersion,
        id: clientID,
        method: 'update_problem',
        params: request,
      },
    });
    if (error) {
      throw error;
    }

    if (!result) {
      throw new Error('No response received');
    }

    return result;
  }

  public async deleteProblem(request: RpcDeleteProblemRequest): Promise<void> {
    const { error, result } = await this.api.deleteProblem({
      requestBodyOfTheDeleteProblemMethod: {
        jsonrpc: jsonRPCVersion,
        id: clientID,
        method: 'delete_problem',
        params: request,
      },
    });
    if (error) {
      throw error;
    }

    if (!result) {
      throw new Error('No response received');
    }
  }

  public async createTestCase(
    request: RpcCreateTestCaseRequest
  ): Promise<RpcCreateTestCaseResponse> {
    const { error, result } = await this.api.createTestCase({
      requestBodyOfTheCreateTestCaseMethod: {
        jsonrpc: jsonRPCVersion,
        id: clientID,
        method: 'create_test_case',
        params: request,
      },
    });
    if (error) {
      throw error;
    }

    if (!result) {
      throw new Error('No response received');
    }

    return result;
  }

  public async createTestCaseList(
    request: RpcCreateTestCaseListRequest
  ): Promise<void> {
    const { error, result } = await this.api.createTestCaseList({
      requestBodyOfTheCreateTestCaseListMethod: {
        jsonrpc: jsonRPCVersion,
        id: clientID,
        method: 'create_test_case_list',
        params: request,
      },
    });
    if (error) {
      throw error;
    }

    if (!result) {
      throw new Error('No response received');
    }
  }

  public async getProblemTestCaseSnippetList(
    request: RpcGetProblemTestCaseSnippetListRequest
  ): Promise<RpcGetProblemTestCaseSnippetListResponse> {
    const { error, result } = await this.api.getProblemTestCaseSnippetList({
      requestBodyOfTheGetProblemTestCaseSnippetListMethod: {
        jsonrpc: jsonRPCVersion,
        id: clientID,
        method: 'get_problem_test_case_snippet_list',
        params: request,
      },
    });
    if (error) {
      throw error;
    }

    if (!result) {
      throw new Error('No response received');
    }

    return result;
  }

  public async getTestCase(
    request: RpcGetTestCaseRequest
  ): Promise<RpcGetTestCaseResponse> {
    const { error, result } = await this.api.getTestCase({
      requestBodyOfTheGetTestCaseMethod: {
        jsonrpc: jsonRPCVersion,
        id: clientID,
        method: 'get_test_case',
        params: request,
      },
    });
    if (error) {
      throw error;
    }

    if (!result) {
      throw new Error('No response received');
    }

    return result;
  }

  public async updateTestCase(
    request: RpcUpdateTestCaseRequest
  ): Promise<RpcUpdateTestCaseResponse> {
    const { error, result } = await this.api.updateTestCase({
      requestBodyOfTheUpdateTestCaseMethod: {
        jsonrpc: jsonRPCVersion,
        id: clientID,
        method: 'update_test_case',
        params: request,
      },
    });
    if (error) {
      throw error;
    }

    if (!result) {
      throw new Error('No response received');
    }

    return result;
  }

  public async deleteTestCase(
    request: RpcDeleteTestCaseRequest
  ): Promise<void> {
    const { error, result } = await this.api.deleteTestCase({
      requestBodyOfTheDeleteTestCaseMethod: {
        jsonrpc: jsonRPCVersion,
        id: clientID,
        method: 'delete_test_case',
        params: request,
      },
    });
    if (error) {
      throw error;
    }

    if (!result) {
      throw new Error('No response received');
    }
  }

  public async getAccountProblemSnippetList(
    request: RpcGetAccountProblemSnippetListRequest
  ): Promise<RpcGetAccountProblemSnippetListResponse> {
    const { error, result } = await this.api.getAccountProblemSnippetList({
      requestBodyOfTheGetAccountProblemSnippetListMethod: {
        jsonrpc: jsonRPCVersion,
        id: clientID,
        method: 'get_account_problem_snippet_list',
        params: request,
      },
    });
    if (error) {
      throw error;
    }

    if (!result) {
      throw new Error('No response received');
    }

    return result;
  }

  public async createSubmission(
    request: RpcCreateSubmissionRequest
  ): Promise<RpcCreateSubmissionResponse> {
    const { error, result } = await this.api.createSubmission({
      requestBodyOfTheCreateSubmissionMethod: {
        jsonrpc: jsonRPCVersion,
        id: clientID,
        method: 'create_submission',
        params: request,
      },
    });
    if (error) {
      throw error;
    }

    if (!result) {
      throw new Error('No response received');
    }

    return result;
  }

  public async getSubmissionSnippetList(
    request: RpcGetSubmissionSnippetListRequest
  ): Promise<RpcGetSubmissionSnippetListResponse> {
    const { error, result } = await this.api.getSubmissionSnippetList({
      requestBodyOfTheGetSubmissionSnippetListMethod: {
        jsonrpc: jsonRPCVersion,
        id: clientID,
        method: 'get_submission_snippet_list',
        params: request,
      },
    });
    if (error) {
      throw error;
    }

    if (!result) {
      throw new Error('No response received');
    }

    return result;
  }

  public async getSubmission(
    request: RpcGetSubmissionRequest
  ): Promise<RpcGetSubmissionResponse> {
    const { error, result } = await this.api.getSubmission({
      requestBodyOfTheGetSubmissionMethod: {
        jsonrpc: jsonRPCVersion,
        id: clientID,
        method: 'get_submission',
        params: request,
      },
    });
    if (error) {
      throw error;
    }

    if (!result) {
      throw new Error('No response received');
    }

    return result;
  }

  public async deleteSubmission(
    request: RpcDeleteSubmissionRequest
  ): Promise<void> {
    const { error, result } = await this.api.deleteSubmission({
      requestBodyOfTheDeleteSubmissionMethod: {
        jsonrpc: jsonRPCVersion,
        id: clientID,
        method: 'delete_submission',
        params: request,
      },
    });
    if (error) {
      throw error;
    }

    if (!result) {
      throw new Error('No response received');
    }
  }

  public async getAccountSubmissionSnippetList(
    request: RpcGetAccountSubmissionSnippetListRequest
  ): Promise<RpcGetAccountSubmissionSnippetListResponse> {
    const { error, result } = await this.api.getAccountSubmissionSnippetList({
      requestBodyOfTheGetAccountSubmissionSnippetListMethod: {
        jsonrpc: jsonRPCVersion,
        id: clientID,
        method: 'get_account_submission_snippet_list',
        params: request,
      },
    });
    if (error) {
      throw error;
    }

    if (!result) {
      throw new Error('No response received');
    }

    return result;
  }

  public async getProblemSubmissionSnippetList(
    request: RpcGetProblemSubmissionSnippetListRequest
  ): Promise<RpcGetProblemSubmissionSnippetListResponse> {
    const { error, result } = await this.api.getProblemSubmissionSnippetList({
      requestBodyOfTheGetProblemSubmissionSnippetListMethod: {
        jsonrpc: jsonRPCVersion,
        id: clientID,
        method: 'get_problem_submission_snippet_list',
        params: request,
      },
    });
    if (error) {
      throw error;
    }

    if (!result) {
      throw new Error('No response received');
    }

    return result;
  }

  public async getAccountProblemSubmissionSnippetList(
    request: RpcGetAccountProblemSubmissionSnippetListRequest
  ): Promise<RpcGetAccountProblemSubmissionSnippetListResponse> {
    const { error, result } =
      await this.api.getAccountProblemSubmissionSnippetList({
        requestBodyOfTheGetAccountProblemSubmissionSnippetListMethod: {
          jsonrpc: jsonRPCVersion,
          id: clientID,
          method: 'get_account_problem_submission_snippet_list',
          params: request,
        },
      });
    if (error) {
      throw error;
    }

    if (!result) {
      throw new Error('No response received');
    }

    return result;
  }

  public async getSession(): Promise<RpcGetSessionResponse> {
    const { error, result } = await this.api.getSession({
      requestBodyOfTheGetSessionMethod: {
        jsonrpc: jsonRPCVersion,
        id: clientID,
        method: 'get_session',
        params: {},
      },
    });
    if (error) {
      throw error;
    }

    if (!result) {
      throw new Error('No response received');
    }

    return result;
  }

  public async updateSetting(
    setting: RpcSetting
  ): Promise<RpcUpdateSettingResponse> {
    const { error, result } = await this.api.updateSetting({
      requestBodyOfTheUpdateSettingMethod: {
        jsonrpc: jsonRPCVersion,
        id: clientID,
        method: 'update_setting',
        params: {
          setting: setting,
        },
      },
    });
    if (error) {
      throw error;
    }

    if (!result) {
      throw new Error('No response received');
    }

    return result;
  }

  public isRpcError(e: unknown): boolean {
    if (!(e instanceof Object)) {
      return false;
    }

    return instanceOfRpcError(e);
  }
}
