import { Component, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { NzButtonModule } from 'ng-zorro-antd/button';
import { NzCheckboxModule } from 'ng-zorro-antd/checkbox';
import {
  RpcSetting,
  RpcSettingFromJSON,
  RpcSettingToJSON,
} from '../../dataaccess/api';
import { CommonModule } from '@angular/common';
import { ServerService } from '../../logic/server.service';
import {
  NzNotificationModule,
  NzNotificationService,
} from 'ng-zorro-antd/notification';

import {
  UnauthenticatedError,
  PermissionDeniedError,
} from '../../logic/account.service';
import { Router } from '@angular/router';
import { PageTitleService } from '../../logic/page-title.service';

@Component({
  selector: 'app-setting',
  standalone: true,
  imports: [
    FormsModule,
    NzButtonModule,
    NzCheckboxModule,
    CommonModule,
    NzNotificationModule,
  ],
  templateUrl: './setting.component.html',
  styleUrl: './setting.component.scss',
})
export class SettingComponent implements OnInit {
  public setting: RpcSetting | undefined;
  public saving = false;

  constructor(
    private readonly serverService: ServerService,
    private readonly notificationService: NzNotificationService,
    private readonly router: Router,
    private readonly pageTitleService: PageTitleService
  ) {}

  ngOnInit(): void {
    (async () => {
      try {
        const serverInfo = await this.serverService.getServerInfo();
        this.setting = RpcSettingFromJSON(RpcSettingToJSON(serverInfo.setting));
      } catch (e) {
        this.notificationService.error('Failed to get server information', '');
        return;
      }
    })().then();
    this.pageTitleService.setTitle('Settings');
  }

  public async onSaveClicked(): Promise<void> {
    if (!this.setting) {
      return;
    }

    try {
      this.setting = await this.serverService.updateSetting(this.setting);
      this.notificationService.success('Updated settings successfully', '');
    } catch (e) {
      if (e instanceof UnauthenticatedError) {
        this.notificationService.error(
          'Failed to update settings',
          'Not logged in'
        );
        this.router.navigateByUrl('/login');
        return;
      }
      if (e instanceof PermissionDeniedError) {
        this.notificationService.error(
          'Failed to update settings',
          'Permission denied'
        );
        return;
      }
      this.notificationService.error(
        'Failed to update settings',
        'Unknown error'
      );
    }
  }
}
