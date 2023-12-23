import { Injectable } from '@angular/core';
import {
  Configuration,
  DefaultApi,
  RpcCreateAccountRequest,
  RpcCreateAccountResponse,
} from './api';

const jsonRPCVersion = '2.0';
const clientID = 'cato-judge';

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
}
