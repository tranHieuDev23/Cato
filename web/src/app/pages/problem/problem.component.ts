import { Component, OnInit, ViewChild } from '@angular/core';
import { NzTypographyModule } from 'ng-zorro-antd/typography';
import { RpcAccount, RpcProblem } from '../../dataaccess/api';
import { CommonModule, Location } from '@angular/common';
import { ActivatedRoute, Params, Router } from '@angular/router';
import {
  AccountService,
  PermissionDeniedError,
  UnauthenticatedError,
} from '../../logic/account.service';
import {
  ProblemNotFoundError,
  ProblemService,
} from '../../logic/problem.service';
import {
  NzNotificationModule,
  NzNotificationService,
} from 'ng-zorro-antd/notification';
import { CodemirrorModule } from '@ctrl/ngx-codemirror';
import { FormsModule } from '@angular/forms';
import { NzGridModule } from 'ng-zorro-antd/grid';
import { TimeLimitPipe } from '../../components/utils/time-limit.pipe';
import { MemoryPipe } from '../../components/utils/memory.pipe';
import { NzUploadModule } from 'ng-zorro-antd/upload';
import { NzButtonModule } from 'ng-zorro-antd/button';
import { NzMenuModule } from 'ng-zorro-antd/menu';
import { NzSelectModule } from 'ng-zorro-antd/select';
import { NzIconModule } from 'ng-zorro-antd/icon';
import { NzTableModule } from 'ng-zorro-antd/table';
import {
  NzTabComponent,
  NzTabSetComponent,
  NzTabsModule,
} from 'ng-zorro-antd/tabs';
import { SubmissionListComponent } from './submission-list/submission-list.component';
import { CodeEditorComponent } from './code-editor/code-editor.component';
import {
  InvalidSubmissionInfo,
  SubmissionService,
} from '../../logic/submission.service';

@Component({
  selector: 'app-problem',
  standalone: true,
  imports: [
    NzTypographyModule,
    CommonModule,
    NzNotificationModule,
    CodemirrorModule,
    FormsModule,
    NzGridModule,
    TimeLimitPipe,
    MemoryPipe,
    NzUploadModule,
    NzButtonModule,
    NzMenuModule,
    NzSelectModule,
    NzIconModule,
    NzTableModule,
    NzTabsModule,
    SubmissionListComponent,
    CodeEditorComponent,
  ],
  templateUrl: './problem.component.html',
  styleUrl: './problem.component.scss',
})
export class ProblemComponent implements OnInit {
  @ViewChild('problemTabSet') problemTabSet: NzTabSetComponent | undefined;
  @ViewChild('submissionTab') submissionTab: NzTabComponent | undefined;

  public sessionAccount: RpcAccount | undefined;
  public problem: RpcProblem | undefined;

  public submissionContent = '';
  public submissionLanguage = 'cpp';

  constructor(
    private readonly accountService: AccountService,
    private readonly problemService: ProblemService,
    private readonly submissionService: SubmissionService,
    private readonly activatedRoute: ActivatedRoute,
    private readonly router: Router,
    private readonly location: Location,
    private readonly notificationService: NzNotificationService
  ) {}

  ngOnInit(): void {
    (async () => {
      const sessionAccount = await this.accountService.getSessionAccount();
      if (sessionAccount === null) {
        this.notificationService.error(
          'Failed to load profile page',
          'Not logged in'
        );
        this.router.navigateByUrl('/login');
        return;
      }

      this.sessionAccount = sessionAccount;
    })().then();
    this.activatedRoute.params.subscribe(async (params) => {
      await this.onParamsChanged(params);
    });
  }

  private async onParamsChanged(params: Params): Promise<void> {
    if (!params['id']) {
      this.location.back();
      return;
    }

    const problemID = +params['id'];
    try {
      this.problem = await this.problemService.getProblem(problemID);
    } catch (e) {
      if (e instanceof UnauthenticatedError) {
        this.notificationService.error(
          'Failed to load problem',
          'Not logged in'
        );
        this.router.navigateByUrl('/login');
        return;
      }

      if (e instanceof PermissionDeniedError) {
        this.notificationService.error(
          'Failed to load problem',
          'Permission denied'
        );
        this.location.back();
        return;
      }

      if (e instanceof ProblemNotFoundError) {
        this.notificationService.error(
          'Failed to load problem',
          'Problem cannot be found'
        );
        this.location.back();
        return;
      }

      this.notificationService.error('Failed to load problem', 'Unknown error');
      this.location.back();
    }
  }

  public async onSubmitClicked(): Promise<void> {
    if (!this.problem) {
      return;
    }

    try {
      await this.submissionService.createSubmission(
        this.problem?.iD,
        this.submissionContent,
        this.submissionLanguage
      );
      this.notificationService.success('Solution submitted successfully', '');
      if (
        this.problemTabSet &&
        this.submissionTab &&
        this.submissionTab.position
      ) {
        this.problemTabSet.setSelectedIndex(this.submissionTab.position);
      }
    } catch (e) {
      if (e instanceof UnauthenticatedError) {
        this.notificationService.error(
          'Failed to submit solution',
          'Not logged in'
        );
        this.router.navigateByUrl('/login');
        return;
      }

      if (e instanceof PermissionDeniedError) {
        this.notificationService.error(
          'Failed to submit solution',
          'Permission denied'
        );
        this.location.back();
        return;
      }

      if (e instanceof ProblemNotFoundError) {
        this.notificationService.error(
          'Failed to load problem',
          'Problem not found'
        );
        this.location.back();
        return;
      }

      if (e instanceof InvalidSubmissionInfo) {
        this.notificationService.error(
          'Failed to submit solution',
          'Invalid submission file'
        );
        return;
      }

      this.notificationService.error(
        'Failed to submit solution',
        'Unknown error'
      );
      this.location.back();
    }
  }
}
