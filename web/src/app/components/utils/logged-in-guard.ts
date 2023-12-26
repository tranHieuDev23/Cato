import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, Router, UrlTree } from '@angular/router';
import { AccountService, Role } from '../../logic/account.service';
import { RpcAccount } from '../../dataaccess/api';

@Injectable()
export class LoggedInGuard {
  constructor(
    private readonly accountService: AccountService,
    private readonly router: Router
  ) {}

  public async canActivate(
    route: ActivatedRouteSnapshot
  ): Promise<boolean | UrlTree> {
    const sessionAccount = await this.accountService.getSessionAccount();
    if (sessionAccount === null) {
      return this.router.parseUrl('/login');
    }

    return this.isUserAuthorized(route, sessionAccount);
  }

  private isUserAuthorized(
    route: ActivatedRouteSnapshot,
    sessionAccount: RpcAccount
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
        return (
          sessionAccount.role === Role.Admin ||
          sessionAccount.role == Role.ProblemSetter
        );
      default:
        return true;
    }
  }
}
