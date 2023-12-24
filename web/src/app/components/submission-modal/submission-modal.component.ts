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

export interface SubmissionModalData {
  submission: RpcSubmission;
}

@Component({
  selector: 'app-submission-modal',
  standalone: true,
  imports: [
    CodemirrorModule,
    NzDescriptionsModule,
    LanguagePipe,
    CommonModule,
    SubmissionStatusPipe,
    FormsModule,
  ],
  templateUrl: './submission-modal.component.html',
  styleUrl: './submission-modal.component.scss',
})
export class SubmissionModalComponent implements OnInit {
  @Input() public submission: RpcSubmission | undefined;

  private readonly nzModalData: SubmissionModalData | null =
    inject(NZ_MODAL_DATA);

  constructor(readonly codeMirrorService: CodeMirrorService) {}

  ngOnInit(): void {
    if (!this.nzModalData) {
      return;
    }
    this.submission = this.nzModalData.submission;
  }
}
