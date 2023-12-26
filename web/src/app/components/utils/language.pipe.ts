import { Pipe, PipeTransform } from '@angular/core';
import { ServerService } from '../../logic/server.service';
import { NzNotificationService } from 'ng-zorro-antd/notification';

@Pipe({
  name: 'language',
  standalone: true,
})
export class LanguagePipe implements PipeTransform {
  constructor(
    private readonly serverService: ServerService,
    private readonly notificationService: NzNotificationService
  ) {}

  public async transform(language: string | undefined): Promise<string> {
    if (language === undefined) {
      return 'Unknown';
    }

    try {
      const serverInfo = await this.serverService.getServiceInfo();
      for (const item of serverInfo.supportedLanguageList) {
        if (item.value === language) {
          return item.name;
        }
      }

      return language;
    } catch (e) {
      this.notificationService.error('Failed to get server information', '');
      return '';
    }
  }
}
