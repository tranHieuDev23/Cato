import { Component, Input, OnInit, inject } from '@angular/core';
import { NzDescriptionsModule } from 'ng-zorro-antd/descriptions';
import { RpcSubmission } from '../../dataaccess/api';
import { LanguagePipe } from '../utils/language.pipe';
import { CommonModule } from '@angular/common';
import { SubmissionStatusPipe } from '../utils/submission-status.pipe';
import { FormsModule } from '@angular/forms';
import { NZ_MODAL_DATA } from 'ng-zorro-antd/modal';
import { NzButtonModule } from 'ng-zorro-antd/button';
import { NzIconModule } from 'ng-zorro-antd/icon';
import copyToClipboard from 'copy-to-clipboard';
import {
  NzNotificationModule,
  NzNotificationService,
} from 'ng-zorro-antd/notification';
import { SubmissionStatusColorPipe } from '../utils/submission-status-color.pipe';
import { NzTagModule } from 'ng-zorro-antd/tag';
import { NgeMonacoModule } from '@cisstech/nge/monaco';

export interface SubmissionModalData {
  submission: RpcSubmission;
}

@Component({
  selector: 'app-submission-modal',
  standalone: true,
  templateUrl: './submission-modal.component.html',
  styleUrl: './submission-modal.component.scss',
  imports: [
    NgeMonacoModule,
    NzDescriptionsModule,
    LanguagePipe,
    CommonModule,
    SubmissionStatusPipe,
    FormsModule,
    NzButtonModule,
    NzIconModule,
    NzNotificationModule,
    SubmissionStatusPipe,
    SubmissionStatusColorPipe,
    NzTagModule,
  ],
})
export class SubmissionModalComponent implements OnInit {
  @Input() public submission: RpcSubmission | undefined;

  public editorHeight = 0;

  private readonly nzModalData: SubmissionModalData | null =
    inject(NZ_MODAL_DATA);

  constructor(private readonly notificationService: NzNotificationService) {}

  ngOnInit(): void {
    if (!this.nzModalData) {
      return;
    }
    this.submission = this.nzModalData.submission;
  }

  public onMonacoEditorReady(editor: monaco.editor.IEditor): void {
    if (!this.submission) {
      return;
    }

    editor.updateOptions({
      scrollBeyondLastLine: false,
      readOnly: true,
      minimap: { enabled: false },
    });
    const editorModel = monaco.editor.createModel(this.submission.content);
    monaco.editor.setModelLanguage(editorModel, this.submission.language);

    editor.setModel(editorModel);
    editor.layout({
      height: (editor as monaco.editor.ICodeEditor).getContentHeight(),
      width: (editor as monaco.editor.ICodeEditor).getContentWidth(),
    });
  }

  public onCopyClicked(): void {
    if (!this.submission) {
      return;
    }
    copyToClipboard(this.submission.content);
    this.notificationService.success('Code copied to clipboard', '');
  }
}
