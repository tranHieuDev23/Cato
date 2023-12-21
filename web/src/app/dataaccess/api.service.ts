import { Injectable } from '@angular/core';
import {
  Configuration,
  DefaultApi,
  RpcEchoRequest,
  RpcEchoResponse,
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

  public async echo(request: RpcEchoRequest): Promise<RpcEchoResponse> {
    try {
      const { error, result } = await this.api.echo({
        requestBodyOfTheEchoMethod: {
          jsonrpc: jsonRPCVersion,
          id: clientID,
          method: 'echo',
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
