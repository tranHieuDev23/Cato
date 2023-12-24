import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root',
})
export class CodeMirrorService {
  public submissionLanguageToCodeMirrorMode(language: string): string {
    if (language === 'cpp') {
      return 'text/x-c++src';
    }
    if (language === 'java') {
      return 'text/x-java';
    }
    if (language === 'python') {
      return 'text/x-python';
    }
    return 'text';
  }
}
