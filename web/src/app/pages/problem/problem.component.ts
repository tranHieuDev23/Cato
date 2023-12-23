import { Component } from '@angular/core';
import { NzTypographyModule } from 'ng-zorro-antd/typography';
import {
  RpcAccount,
  RpcProblem,
  RpcSubmission,
  RpcSubmissionSnippet,
} from '../../dataaccess/api';
import { CommonModule, Location } from '@angular/common';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { Subscription } from 'rxjs';
import {
  AccountService,
  PermissionDeniedError,
  UnauthenticatedError,
} from '../../logic/account.service';
import { ProblemService } from '../../logic/problem.service';
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
  ],
  templateUrl: './problem.component.html',
  styleUrl: './problem.component.scss',
})
export class ProblemComponent {
  public sessionAccount: RpcAccount | null | undefined;
  public problem: RpcProblem | undefined;

  public submissionContent = '';
  public submissionLanguage = 'cpp';
  public submissionEditorMode = 'text/x-c++src';

  public submissionStatusVisible = false;
  public totalSubmissionCount = 0;
  public submissionSnippetList: RpcSubmissionSnippet[] = [];
  public pageIndex = 1;
  public pageSize = 5;
  public loadingSubmissionList = false;

  private sessionAccountChangedSubscription: Subscription | undefined;

  constructor(
    private readonly accountService: AccountService,
    private readonly problemService: ProblemService,
    private readonly activatedRoute: ActivatedRoute,
    private readonly router: Router,
    private readonly location: Location,
    private readonly notificationService: NzNotificationService
  ) {}

  ngOnInit(): void {
    (async () => {
      this.sessionAccount = await this.accountService.getSessionAccount();
    })().then();
    this.activatedRoute.params.subscribe(async (params) => {
      await this.onParamsChanged(params);
    });
    this.sessionAccountChangedSubscription =
      this.accountService.sessionAccountChanged.subscribe((account) => {
        this.sessionAccount = account;
      });
  }

  ngOnDestroy(): void {
    this.sessionAccountChangedSubscription?.unsubscribe();
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

      this.notificationService.error('Failed to load problem', 'Unknown error');
      this.location.back();
    }
  }

  public onSubmissionLanguageChange(language: string): void {
    if (language === 'cpp') {
      this.submissionEditorMode = 'text/x-c++src';
    }
    if (language === 'python') {
      this.submissionEditorMode = 'text/x-python';
    }
  }

  public async onPageIndexChange(index: number): Promise<void> {
    this.pageIndex = index;
  }

  public async onPageSizeChange(size: number): Promise<void> {
    this.pageSize = size;
  }
}
