import { ApplicationConfig, importProvidersFrom } from '@angular/core';
import { provideRouter } from '@angular/router';

import { routes } from './app.routes';
import { provideNzIcons } from './icons-provider';
import { en_US, provideNzI18n } from 'ng-zorro-antd/i18n';
import { registerLocaleData } from '@angular/common';
import en from '@angular/common/locales/en';
import { FormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';
import { provideAnimations } from '@angular/platform-browser/animations';
import { LoggedOutGuard } from './components/utils/logged-out-guard';
import { LoggedInGuard } from './components/utils/logged-in-guard';
import { NgeMonacoModule } from '@cisstech/nge/monaco';

registerLocaleData(en);

export const appConfig: ApplicationConfig = {
  providers: [
    provideRouter(routes),
    provideNzIcons(),
    provideNzI18n(en_US),
    importProvidersFrom(FormsModule),
    importProvidersFrom(HttpClientModule),
    importProvidersFrom(
      NgeMonacoModule.forRoot({
        theming: {
          themes: ['vs-dark'],
          default: 'vs-dark',
        },
      })
    ),
    provideAnimations(),
    LoggedOutGuard,
    LoggedInGuard,
  ],
};
