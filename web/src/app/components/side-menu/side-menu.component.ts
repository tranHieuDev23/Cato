import { CommonModule } from '@angular/common';
import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { NzIconModule } from 'ng-zorro-antd/icon';
import { NzMenuModule } from 'ng-zorro-antd/menu';
import { AccountService } from '../../logic/account/account.service';
import { Router } from '@angular/router';
import { Subscription } from 'rxjs';
import { RpcAccount } from '../../dataaccess/api';

@Component({
  selector: 'app-side-menu',
  standalone: true,
  imports: [NzMenuModule, NzIconModule, CommonModule],
  templateUrl: './side-menu.component.html',
  styleUrl: './side-menu.component.scss',
})
export class SideMenuComponent implements OnInit, OnDestroy {
  @Input() public collapsed = false;
  public sessionAccount: RpcAccount | null | undefined;

  private sessionAccountChangedSubscription: Subscription | undefined;

  constructor(
    private readonly accountService: AccountService,
    private readonly router: Router
  ) {}

  public async onLogOutClicked(): Promise<void> {
    await this.accountService.deleteSession();
    this.router.navigateByUrl('/login');
  }

  ngOnInit(): void {
    this.sessionAccountChangedSubscription =
      this.accountService.sessionAccountChanged.subscribe((account) => {
        this.sessionAccount = account;
      });
  }

  ngOnDestroy(): void {
    this.sessionAccountChangedSubscription?.unsubscribe();
  }
}
