import { bootstrapApplication } from '@angular/platform-browser';
import { appConfig } from './app/app.config';
import { AppComponent } from './app/app.component';
import 'codemirror/mode/clike/clike';
import 'codemirror/mode/python/python';
import 'codemirror/mode/markdown/markdown';

bootstrapApplication(AppComponent, appConfig).catch((err) =>
  console.error(err)
);
