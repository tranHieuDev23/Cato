import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { CodemirrorModule } from '@ctrl/ngx-codemirror';
import { NzButtonModule } from 'ng-zorro-antd/button';
import { NzIconModule } from 'ng-zorro-antd/icon';
import { NzMenuModule } from 'ng-zorro-antd/menu';
import { NzSelectModule } from 'ng-zorro-antd/select';
import { NzUploadFile, NzUploadModule } from 'ng-zorro-antd/upload';
import { CodeMirrorService } from '../../../logic/code-mirror.service';
import { ServerService } from '../../../logic/server.service';
import {
  NzNotificationModule,
  NzNotificationService,
} from 'ng-zorro-antd/notification';

export interface LanguageOption {
  value: string;
  name: string;
}

@Component({
  selector: 'app-code-editor',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    CodemirrorModule,
    NzMenuModule,
    NzSelectModule,
    NzUploadModule,
    NzButtonModule,
    NzIconModule,
    NzNotificationModule,
  ],
  templateUrl: './code-editor.component.html',
  styleUrl: './code-editor.component.scss',
})
export class CodeEditorComponent implements OnInit {
  @Input() public content = '';
  @Output() public contentChange = new EventEmitter<string>();

  @Input() public language = 'cpp';
  @Output() public languageChange = new EventEmitter<string>();

  public languageOptionList: LanguageOption[] = [];

  @Output() public submitClicked = new EventEmitter<void>();

  public editorMode = 'text/x-c++src';

  constructor(
    private readonly codeMirrorService: CodeMirrorService,
    private readonly serverService: ServerService,
    private readonly notificationService: NzNotificationService
  ) {}

  ngOnInit(): void {
    (async () => {
      try {
        const serverInfo = await this.serverService.getServiceInfo();
        this.languageOptionList = serverInfo.supportedLanguageList.map(
          (item) => {
            return { value: item.value, name: item.name };
          }
        );
      } catch (e) {
        this.notificationService.error('Failed to get server information', '');
        return;
      }
    })().then();
  }

  public onLoadFile = (file: NzUploadFile): boolean => {
    const fileReader = new FileReader();
    this.updateLanguageFromFileName(file.name);

    fileReader.onload = (event) => {
      this.content = `${event.target?.result || ''}`;
    };

    fileReader.readAsText(file as any);
    return false;
  };

  private updateLanguageFromFileName(fileName: string): void {
    if (fileName.endsWith('.cpp')) {
      this.language = 'cpp';
      this.onSubmissionLanguageChange(this.language);
    }
    if (fileName.endsWith('.c')) {
      this.language = 'c';
      this.onSubmissionLanguageChange(this.language);
    }
    if (fileName.endsWith('.java')) {
      this.language = 'java';
      this.onSubmissionLanguageChange(this.language);
    }
    if (fileName.endsWith('.py')) {
      this.language = 'python';
      this.onSubmissionLanguageChange(this.language);
    }
  }

  public onSubmissionLanguageChange(language: string): void {
    this.editorMode =
      this.codeMirrorService.submissionLanguageToCodeMirrorMode(language);
    this.languageChange.emit(language);
  }

  public onContentChange(content: string): void {
    this.contentChange.emit(content);
  }

  public onSubmitClicked(): void {
    this.submitClicked.emit();
  }
}
