import { EventEmitter, Injectable } from '@angular/core';
import { ApiService } from '../dataaccess';
import {
  RpcError,
  RpcGetServerInfoResponse,
  RpcSetting,
} from '../dataaccess/api';
import { ErrorCode } from '../dataaccess/api.service';
import { PermissionDeniedError, UnauthenticatedError } from './account.service';

@Injectable({
  providedIn: 'root',
})
export class ServerService {
  private serverInfo: RpcGetServerInfoResponse | undefined;

  public readonly serverInfoChanged =
    new EventEmitter<RpcGetServerInfoResponse>();

  constructor(private readonly api: ApiService) {}

  public async getServerInfo(): Promise<RpcGetServerInfoResponse> {
    if (this.serverInfo === undefined) {
      this.serverInfo = await this.api.getServerInfo();
    }

    return this.serverInfo;
  }

  public async updateSetting(setting: RpcSetting): Promise<RpcSetting> {
    try {
      const response = await this.api.updateSetting(setting);
      const serverInfo = await this.getServerInfo();
      serverInfo.setting = response.setting;
      this.serverInfo = serverInfo;
      this.serverInfoChanged.emit(serverInfo);
      return response.setting;
    } catch (e) {
      if (!this.api.isRpcError(e)) {
        throw e;
      }

      const apiError = e as RpcError;

      if (apiError.code === ErrorCode.Unauthenticated) {
        throw new UnauthenticatedError();
      }

      if (apiError.code === ErrorCode.PermissionDenied) {
        throw new PermissionDeniedError();
      }

      throw apiError;
    }
  }
}
