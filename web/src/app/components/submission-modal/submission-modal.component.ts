import { Component, Input, OnInit, inject } from '@angular/core';
import { CodemirrorModule } from '@ctrl/ngx-codemirror';
import { NzDescriptionsModule } from 'ng-zorro-antd/descriptions';
import { RpcSubmission } from '../../dataaccess/api';
import { LanguagePipe } from '../utils/language.pipe';
import { CommonModule } from '@angular/common';
import { SubmissionStatusPipe } from '../utils/submission-status.pipe';
import { FormsModule } from '@angular/forms';
import { CodeMirrorService } from '../../logic/code-mirror.service';
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

export interface SubmissionModalData {
  submission: RpcSubmission;
}

@Component({
  selector: 'app-submission-modal',
  standalone: true,
  templateUrl: './submission-modal.component.html',
  styleUrl: './submission-modal.component.scss',
  imports: [
    CodemirrorModule,
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

  private readonly nzModalData: SubmissionModalData | null =
    inject(NZ_MODAL_DATA);

  constructor(
    readonly codeMirrorService: CodeMirrorService,
    private readonly notificationService: NzNotificationService
  ) {}

  ngOnInit(): void {
    if (!this.nzModalData) {
      return;
    }
    this.submission = this.nzModalData.submission;
  }

  public onCopyClicked(): void {
    if (!this.submission) {
      return;
    }
    copyToClipboard(this.submission.content);
    this.notificationService.success('Code copied to clipboard', '');
  }
}
