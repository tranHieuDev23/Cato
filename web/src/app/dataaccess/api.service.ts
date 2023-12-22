import { Injectable } from '@angular/core';
import {
  Configuration,
  DefaultApi,
  RpcCreateUserRequest,
  RpcCreateUserResponse,
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

  public async createUser(
    request: RpcCreateUserRequest
  ): Promise<RpcCreateUserResponse> {
    try {
      const { error, result } = await this.api.createUser({
        requestBodyOfTheCreateUserMethod: {
          jsonrpc: jsonRPCVersion,
          id: clientID,
          method: 'create_user',
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
    } catch (e) {
      throw e;
    }
  }
}
