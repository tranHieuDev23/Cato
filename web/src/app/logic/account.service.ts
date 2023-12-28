import { EventEmitter, Injectable } from '@angular/core';
import { ApiService } from '../dataaccess';
import {
  RpcAccount,
  RpcError,
  RpcUpdateAccountRequestFromJSON,
} from '../dataaccess/api';
import { ErrorCode } from '../dataaccess/api.service';

export enum Role {
  Admin = 'admin',
  ProblemSetter = 'problem_setter',
  Contestant = 'contestant',
  Worker = 'worker',
}

export class AccountCreationDisabledError extends Error {
  constructor() {
    super('Account creation is disabled');
  }
}

export class AccountUpdateDisabledError extends Error {
  constructor() {
    super('Account update is disabled');
  }
}

export class AccountLoginDisabledError extends Error {
  constructor() {
    super('Login is disabled');
  }
}

export class AccountNotFoundError extends Error {
  constructor() {
    super('Account not found');
  }
}

export class IncorrectPasswordError extends Error {
  constructor() {
    super('Account not found');
  }
}

export class InvalidAccountInfoError extends Error {
  constructor() {
    super('Invalid account information');
  }
}

export class AccountNameTakenError extends Error {
  constructor() {
    super('Account name taken');
  }
}

export class UnauthenticatedError extends Error {
  constructor() {
    super('Not logged in');
  }
}

export class PermissionDeniedError extends Error {
  constructor() {
    super('Permission denied');
  }
}

export class InvalidAccountListParam extends Error {
  constructor() {
    super('Invalid account list parameters');
  }
}

@Injectable({
  providedIn: 'root',
})
export class AccountService {
  private sessionAccount: RpcAccount | null | undefined;

  public readonly sessionAccountChanged = new EventEmitter<RpcAccount | null>();

  constructor(private readonly api: ApiService) {}

  public isValidAccountName(username: string): { [k: string]: boolean } | null {
    if (username.length < 6) {
      return { error: true, minLength: true };
    }
    if (username.length > 32) {
      return { error: true, maxLength: true };
    }
    if (!/^[a-zA-Z0-9]*$/.test(username)) {
      return { error: true, pattern: true };
    }
    return null;
  }

  public isValidDisplayName(
    displayName: string
  ): { [k: string]: boolean } | null {
    if (displayName.length > 32) {
      return { error: true, maxLength: true };
    }
    return null;
  }

  public isValidPassword(password: string): { [k: string]: boolean } | null {
    if (0 < password.length && password.length < 8) {
      return { error: true, minLength: true };
    }
    return null;
  }

  public async createAccount(
    accountName: string,
    displayName: string,
    role: string,
    password: string
  ): Promise<RpcAccount> {
    try {
      const response = await this.api.createAccount({
        accountName,
        displayName,
        role,
        password,
      });
      return response.account;
    } catch (e) {
      if (!this.api.isRpcError(e)) {
        throw e;
      }

      const apiError = e as RpcError;
      if (apiError.code === ErrorCode.JRPCErrorInvalidParams) {
        throw new InvalidAccountInfoError();
      }

      if (apiError.code === ErrorCode.AlreadyExists) {
        throw new AccountNameTakenError();
      }

      if (apiError.code === ErrorCode.Unavailable) {
        throw new AccountCreationDisabledError();
      }

      throw apiError;
    }
  }

  public async updateAccount(
    id: number,
    displayName: string | undefined,
    role: string | undefined,
    password: string | undefined
  ): Promise<RpcAccount> {
    try {
      const response = await this.api.updateAccount(
        RpcUpdateAccountRequestFromJSON({
          ID: id,
          DisplayName: displayName,
          Role: role,
          Password: password,
        })
      );
      if (response.account.iD === this.sessionAccount?.iD) {
        this.sessionAccount = response.account;
        this.sessionAccountChanged.emit(response.account);
      }
      return response.account;
    } catch (e) {
      if (!this.api.isRpcError(e)) {
        throw e;
      }

      const apiError = e as RpcError;
      if (apiError.code === ErrorCode.JRPCErrorInvalidParams) {
        throw new InvalidAccountInfoError();
      }

      if (apiError.code === ErrorCode.Unauthenticated) {
        throw new UnauthenticatedError();
      }

      if (apiError.code === ErrorCode.PermissionDenied) {
        throw new PermissionDeniedError();
      }

      if (apiError.code === ErrorCode.Unavailable) {
        throw new AccountUpdateDisabledError();
      }

      throw apiError;
    }
  }

  public async createSession(
    accountName: string,
    password: string
  ): Promise<RpcAccount> {
    try {
      const response = await this.api.createSession({
        accountName,
        password,
      });
      this.sessionAccount = response.account;
      this.sessionAccountChanged.emit(this.sessionAccount);
      return response.account;
    } catch (e) {
      if (!this.api.isRpcError(e)) {
        throw e;
      }

      const apiError = e as RpcError;
      if (apiError.code === ErrorCode.NotFound) {
        throw new AccountNotFoundError();
      }

      if (
        apiError.code === ErrorCode.Unauthenticated ||
        apiError.code === ErrorCode.JRPCErrorInvalidParams
      ) {
        throw new IncorrectPasswordError();
      }

      if (apiError.code === ErrorCode.Unavailable) {
        throw new AccountLoginDisabledError();
      }

      throw apiError;
    }
  }

  public async deleteSession(): Promise<void> {
    try {
      await this.api.deleteSession();
      this.sessionAccount = null;
    } catch (e) {
      if (!this.api.isRpcError(e)) {
        throw e;
      }

      const apiError = e as RpcError;
      if (apiError.code === ErrorCode.Unauthenticated) {
        this.sessionAccount = null;
        this.sessionAccountChanged.emit(null);
        return;
      }

      throw apiError;
    }
  }

  public async getSessionAccount(): Promise<RpcAccount | null> {
    if (this.sessionAccount == undefined) {
      try {
        const response = await this.api.getSession();
        this.sessionAccount = response.account;
        this.sessionAccountChanged.emit(this.sessionAccount);
      } catch (e) {
        if (!this.api.isRpcError(e)) {
          throw e;
        }

        const apiError = e as RpcError;
        if (apiError.code === ErrorCode.Unauthenticated) {
          this.sessionAccount = null;
          this.sessionAccountChanged.emit(null);
          return null;
        } else {
          throw apiError;
        }
      }
    }

    return this.sessionAccount;
  }

  public async getAccountList(
    offset: number,
    limit: number
  ): Promise<{
    totalAccountCount: number;
    accountList: RpcAccount[];
  }> {
    try {
      const response = await this.api.getAccountList({ offset, limit });
      return {
        totalAccountCount: response.totalAccountCount,
        accountList: response.accountList,
      };
    } catch (e) {
      if (!this.api.isRpcError(e)) {
        throw e;
      }

      const apiError = e as RpcError;
      if (apiError.code == ErrorCode.JRPCErrorInvalidParams) {
        throw new InvalidAccountListParam();
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
}
