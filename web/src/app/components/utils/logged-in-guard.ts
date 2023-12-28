import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, Router, UrlTree } from '@angular/router';
import { AccountService, Role } from '../../logic/account.service';
import { RpcAccount, RpcGetServerInfoResponse } from '../../dataaccess/api';
import { ServerService } from '../../logic/server.service';

@Injectable()
export class LoggedInGuard {
  constructor(
    private readonly accountService: AccountService,
    private readonly serverService: ServerService,
    private readonly router: Router
  ) {}

  public async canActivate(
    route: ActivatedRouteSnapshot
  ): Promise<boolean | UrlTree> {
    const sessionAccount = await this.accountService.getSessionAccount();
    if (sessionAccount === null) {
      return this.router.parseUrl('/login');
    }

    const serverInfo = await this.serverService.getServerInfo();
    return this.isUserAuthorized(route, sessionAccount, serverInfo);
  }

  private isUserAuthorized(
    route: ActivatedRouteSnapshot,
    sessionAccount: RpcAccount,
    serverInfo: RpcGetServerInfoResponse
  ): boolean {
    if (route.url.length === 0) {
      return true;
    }

    switch (route.url[0].path) {
      case 'profile':
        return true;
      case 'account-list':
        return sessionAccount.role === Role.Admin;
      case 'problem-list':
        return (
          sessionAccount.role === Role.Admin ||
          sessionAccount.role == Role.ProblemSetter ||
          sessionAccount.role == Role.Contestant
        );
      case 'submission-list':
        return (
          sessionAccount.role === Role.Admin ||
          sessionAccount.role == Role.ProblemSetter ||
          sessionAccount.role == Role.Contestant
        );
      case 'problem':
        return (
          sessionAccount.role === Role.Admin ||
          sessionAccount.role == Role.ProblemSetter ||
          sessionAccount.role == Role.Contestant
        );
      case 'problem-editor':
        if (route.url.length === 1) {
          // Create problem path
          return (
            !serverInfo.setting.problem.disableProblemCreation &&
            (sessionAccount.role === Role.Admin ||
              sessionAccount.role == Role.ProblemSetter)
          );
        }
        // Update problem path
        return (
          !serverInfo.setting.problem.disableProblemUpdate &&
          (sessionAccount.role === Role.Admin ||
            sessionAccount.role == Role.ProblemSetter)
        );
      case 'settings':
        return sessionAccount.role === Role.Admin;
      default:
        return true;
    }
  }
}
