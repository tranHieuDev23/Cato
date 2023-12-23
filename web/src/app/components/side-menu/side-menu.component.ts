import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { NzIconModule } from 'ng-zorro-antd/icon';
import { NzMenuModule } from 'ng-zorro-antd/menu';
import { AccountService } from '../../logic/account/account.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-side-menu',
  standalone: true,
  imports: [NzMenuModule, NzIconModule, CommonModule],
  templateUrl: './side-menu.component.html',
  styleUrl: './side-menu.component.scss',
})
export class SideMenuComponent {
  @Input() public collapsed = false;

  constructor(
    private readonly accountService: AccountService,
    private readonly router: Router
  ) {}

  public async onLogOutClicked(): Promise<void> {
    await this.accountService.deleteSession();
    this.router.navigateByUrl('/login');
  }
}
