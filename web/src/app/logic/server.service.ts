import { Injectable } from '@angular/core';
import { ApiService } from '../dataaccess';
import { RpcGetServerInfoResponse } from '../dataaccess/api';

@Injectable({
  providedIn: 'root',
})
export class ServerService {
  private serverInfo: RpcGetServerInfoResponse | undefined;

  constructor(private readonly api: ApiService) {}

  public async getServiceInfo(): Promise<RpcGetServerInfoResponse> {
    if (this.serverInfo === undefined) {
      this.serverInfo = await this.api.getServerInfo();
    }

    return this.serverInfo;
  }
}
