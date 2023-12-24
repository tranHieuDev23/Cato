import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { CodemirrorModule } from '@ctrl/ngx-codemirror';
import { NzButtonModule } from 'ng-zorro-antd/button';
import { NzIconModule } from 'ng-zorro-antd/icon';
import { NzMenuModule } from 'ng-zorro-antd/menu';
import { NzSelectModule } from 'ng-zorro-antd/select';
import { NzUploadFile, NzUploadModule } from 'ng-zorro-antd/upload';

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
  ],
  templateUrl: './code-editor.component.html',
  styleUrl: './code-editor.component.scss',
})
export class CodeEditorComponent {
  @Input() public content = '';
  @Output() public contentChange = new EventEmitter<string>();

  @Input() public language = 'cpp';
  @Output() public languageChange = new EventEmitter<string>();

  @Output() public submitClicked = new EventEmitter<void>();

  public editorMode = 'text/x-c++src';

  public onLoadFile = (file: NzUploadFile): boolean => {
    const fileReader = new FileReader();
    fileReader.onload = (event) => {
      this.content = `${event.target?.result || ''}`;
    };
    fileReader.readAsText(file as any);
    return false;
  };

  public onSubmissionLanguageChange(language: string): void {
    if (language === 'cpp') {
      this.editorMode = 'text/x-c++src';
    }
    if (language === 'java') {
      this.editorMode = 'text/x-java';
    }
    if (language === 'python') {
      this.editorMode = 'text/x-python';
    }
    this.languageChange.emit(language);
  }

  public onContentChange(content: string): void {
    this.contentChange.emit(content);
  }

  public onSubmitClicked(): void {
    this.submitClicked.emit();
  }
}
