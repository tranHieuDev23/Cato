import {
  Component,
  ElementRef,
  OnDestroy,
  OnInit,
  ViewChild,
} from '@angular/core';
import { NzTypographyModule } from 'ng-zorro-antd/typography';
import { RpcAccount, RpcProblem } from '../../dataaccess/api';
import { CommonModule, Location } from '@angular/common';
import { ActivatedRoute, Params, Router, RouterModule } from '@angular/router';
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
import { NzDescriptionsModule } from 'ng-zorro-antd/descriptions';
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
import { PageTitleService } from '../../logic/page-title.service';
import { TestCaseListComponent } from './test-case-list/test-case-list.component';
import { NzModalModule, NzModalService } from 'ng-zorro-antd/modal';
import { Subscription } from 'rxjs';
import renderMathInElement from 'katex/contrib/auto-render';

@Component({
  selector: 'app-problem',
  standalone: true,
  imports: [
    NzTypographyModule,
    CommonModule,
    NzNotificationModule,
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
    NzDescriptionsModule,
    TestCaseListComponent,
    RouterModule,
    NzModalModule,
  ],
  templateUrl: './problem.component.html',
  styleUrl: './problem.component.scss',
})
export class ProblemComponent implements OnInit, OnDestroy {
  @ViewChild('problemDescriptionContainer') problemDescriptionContainer:
    | ElementRef
    | undefined;
  @ViewChild('problemTabSet') problemTabSet: NzTabSetComponent | undefined;
  @ViewChild('submissionTab') submissionTab: NzTabComponent | undefined;

  public sessionAccount: RpcAccount | undefined;
  public problem: RpcProblem | undefined;

  public submissionContent = '';
  public submissionLanguage = 'cpp';

  private queryParamsSubscription: Subscription | undefined;

  constructor(
    private readonly accountService: AccountService,
    private readonly problemService: ProblemService,
    private readonly submissionService: SubmissionService,
    private readonly activatedRoute: ActivatedRoute,
    private readonly router: Router,
    private readonly location: Location,
    private readonly notificationService: NzNotificationService,
    private readonly modalService: NzModalService,
    private readonly pageTitleService: PageTitleService
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
    this.queryParamsSubscription = this.activatedRoute.params.subscribe(
      async (params) => {
        await this.onParamsChanged(params);
      }
    );
  }

  ngOnDestroy(): void {
    this.queryParamsSubscription?.unsubscribe();
  }

  private async onParamsChanged(params: Params): Promise<void> {
    if (!params['uuid']) {
      this.location.back();
      return;
    }

    const problemUUID = `${params['uuid']}`;
    try {
      this.problem = await this.problemService.getProblem(problemUUID);
      this.pageTitleService.setTitle(this.problem.displayName);
      setTimeout(() => {
        if (!this.problemDescriptionContainer) {
          return;
        }
        renderMathInElement(this.problemDescriptionContainer.nativeElement, {
          throwOnError: false,
        });
      }, 0);
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

  public onDeleteClicked(): void {
    if (!this.problem) {
      return;
    }

    const problemUUID = this.problem.uUID;
    this.modalService.create({
      nzContent: 'Are you sure? This action is <b>irreversible</b>',
      nzOkDanger: true,
      nzOnOk: async () => {
        try {
          await this.problemService.deleteProblem(problemUUID);
          this.notificationService.success('Problem deleted successfully', '');
          this.location.back();
        } catch (e) {
          if (e instanceof UnauthenticatedError) {
            this.notificationService.error(
              'Failed to delete problem',
              'Not logged in'
            );
            this.router.navigateByUrl('/login');
            return;
          }

          if (e instanceof PermissionDeniedError) {
            this.notificationService.error(
              'Failed to delete problem',
              'Permission denied'
            );
            this.location.back();
            return;
          }

          if (e instanceof ProblemNotFoundError) {
            this.notificationService.error(
              'Failed to delete problem',
              'Problem not found'
            );
            this.location.back();
            return;
          }

          this.notificationService.error(
            'Failed to delete problem',
            'Unknown error'
          );
        }
      },
    });
  }

  public async onSubmitClicked(): Promise<void> {
    if (!this.problem) {
      return;
    }

    try {
      await this.submissionService.createSubmission(
        this.problem.uUID,
        this.submissionContent,
        this.submissionLanguage
      );
      this.notificationService.success('Submission submitted successfully', '');
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
          'Failed to submit submission',
          'Not logged in'
        );
        this.router.navigateByUrl('/login');
        return;
      }

      if (e instanceof PermissionDeniedError) {
        this.notificationService.error(
          'Failed to submit submission',
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
          'Failed to submit submission',
          'Invalid submission file'
        );
        return;
      }

      this.notificationService.error(
        'Failed to submit submission',
        'Unknown error'
      );
    }
  }
}
