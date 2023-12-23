import { Injectable } from '@angular/core';
import { Router, UrlTree } from '@angular/router';
import { AccountService } from '../../logic/account/account.service';

@Injectable()
export class LoggedOutGuard {
  constructor(
    private readonly accountService: AccountService,
    private readonly router: Router
  ) {}

  public async canActivate(): Promise<boolean | UrlTree> {
    const sessionAccount = await this.accountService.getSessionAccount();
    if (sessionAccount === null) {
      return true;
    }

    return this.router.parseUrl('/welcome');
  }
}
