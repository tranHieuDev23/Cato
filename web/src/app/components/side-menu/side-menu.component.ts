import { CommonModule } from '@angular/common';
import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { NzIconModule } from 'ng-zorro-antd/icon';
import { NzMenuModule } from 'ng-zorro-antd/menu';
import { AccountService } from '../../logic/account.service';
import { Router, RouterModule } from '@angular/router';
import { Subscription } from 'rxjs';
import { RpcAccount, RpcGetServerInfoResponse } from '../../dataaccess/api';
import { NzToolTipModule } from 'ng-zorro-antd/tooltip';
import { ServerService } from '../../logic/server.service';
import {
  NzNotificationModule,
  NzNotificationService,
} from 'ng-zorro-antd/notification';

@Component({
  selector: 'app-side-menu',
  standalone: true,
  imports: [
    NzMenuModule,
    NzIconModule,
    CommonModule,
    RouterModule,
    NzToolTipModule,
    NzNotificationModule,
  ],
  templateUrl: './side-menu.component.html',
  styleUrl: './side-menu.component.scss',
})
export class SideMenuComponent implements OnInit, OnDestroy {
  @Input() public collapsed = false;
  public sessionAccount: RpcAccount | null | undefined;
  public serverInfo: RpcGetServerInfoResponse | undefined;

  private sessionAccountChangedSubscription: Subscription;
  private serverInfoChangedSubscription: Subscription;

  constructor(
    private readonly accountService: AccountService,
    private readonly serverService: ServerService,
    private readonly notificationService: NzNotificationService,
    private readonly router: Router
  ) {
    this.sessionAccountChangedSubscription =
      this.accountService.sessionAccountChanged.subscribe((account) => {
        this.sessionAccount = account;
      });
    this.serverInfoChangedSubscription =
      this.serverService.serverInfoChanged.subscribe((serverInfo) => {
        this.serverInfo = serverInfo;
      });
  }

  public async onLogOutClicked(): Promise<void> {
    await this.accountService.deleteSession();
    this.router.navigateByUrl('/login');
  }

  ngOnInit(): void {
    (async () => {
      this.sessionAccount = await this.accountService.getSessionAccount();
      try {
        this.serverInfo = await this.serverService.getServerInfo();
      } catch (e) {
        this.notificationService.error('Failed to get server information', '');
        return;
      }
    })().then();
  }

  ngOnDestroy(): void {
    this.sessionAccountChangedSubscription.unsubscribe();
    this.serverInfoChangedSubscription.unsubscribe();
  }
}
