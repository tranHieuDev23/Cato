import { Pipe, PipeTransform } from '@angular/core';
import { ServerService } from '../../logic/server.service';

@Pipe({
  name: 'language',
  standalone: true,
})
export class LanguagePipe implements PipeTransform {
  constructor(private readonly serverService: ServerService) {}

  public async transform(language: string | undefined): Promise<string> {
    if (language === undefined) {
      return 'Unknown';
    }

    const serverInfo = await this.serverService.getServiceInfo();
    for (const item of serverInfo.supportedLanguageList) {
      if (item.value === language) {
        return item.name;
      }
    }

    return language;
  }
}
