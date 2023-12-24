import { Component, Input, OnInit } from '@angular/core';
import { RpcSubmissionSnippet } from '../../../dataaccess/api';
import { NzTableModule } from 'ng-zorro-antd/table';
import { CommonModule, Location } from '@angular/common';
import {
  NzNotificationModule,
  NzNotificationService,
} from 'ng-zorro-antd/notification';
import { Router } from '@angular/router';
import {
  AccountNotFoundError,
  AccountService,
  PermissionDeniedError,
  Role,
  UnauthenticatedError,
} from '../../../logic/account.service';
import { PaginationService } from '../../../logic/pagination.service';
import {
  InvalidSubmissionListParam,
  SubmissionService,
} from '../../../logic/submission.service';
import { NzSwitchModule } from 'ng-zorro-antd/switch';
import { FormsModule } from '@angular/forms';
import { ProblemNotFoundError } from '../../../logic/problem.service';
import { LanguagePipe } from '../../../components/utils/language.pipe';
import { SubmissionStatusPipe } from '../../../components/utils/submission-status.pipe';

const SUBMISSION_LIST_RELOAD_INTERVAL = 10000;

@Component({
  selector: 'app-submission-list',
  standalone: true,
  templateUrl: './submission-list.component.html',
  styleUrl: './submission-list.component.scss',
  imports: [
    NzTableModule,
    CommonModule,
    NzNotificationModule,
    NzSwitchModule,
    FormsModule,
    LanguagePipe,
    SubmissionStatusPipe,
  ],
})
export class SubmissionListComponent implements OnInit {
  @Input() public problemID = 0;

  public totalSubmissionCount = 0;
  public submissionSnippetList: RpcSubmissionSnippet[] = [];
  public pageIndex = 1;
  public pageSize = 10;
  public loading = false;
  public lastLoadedTime: number | undefined;
  public autoReloadEnabled = true;

  private submissionListReloadInterval:
    | ReturnType<typeof setInterval>
    | undefined;

  constructor(
    private readonly accountService: AccountService,
    private readonly submissionService: SubmissionService,
    private readonly paginationService: PaginationService,
    private readonly router: Router,
    private readonly nzNotificationService: NzNotificationService,
    private readonly location: Location
  ) {}

  ngOnInit(): void {
    this.loadSubmissionSnippetList().then();
    this.submissionListReloadInterval = setInterval(async () => {
      await this.loadSubmissionSnippetList();
    }, SUBMISSION_LIST_RELOAD_INTERVAL);
  }

  private async loadSubmissionSnippetList(): Promise<void> {
    try {
      this.loading = true;
      const sessionAccount = await this.accountService.getSessionAccount();
      if (sessionAccount === null) {
        this.nzNotificationService.error(
          'Failed to load submission list',
          'Not logged in'
        );
        this.router.navigateByUrl('/login');
        return;
      }

      if (
        sessionAccount.role === Role.Admin ||
        sessionAccount.role === Role.ProblemSetter
      ) {
        const { totalSubmissionCount, submissionSnippetList } =
          await this.submissionService.getProblemSubmissionSnippetList(
            this.problemID,
            this.paginationService.getPageOffset(this.pageIndex, this.pageSize),
            this.pageSize
          );
        this.totalSubmissionCount = totalSubmissionCount;
        this.submissionSnippetList = submissionSnippetList;
      } else {
        const { totalSubmissionCount, submissionSnippetList } =
          await this.submissionService.getProblemSubmissionSnippetList(
            this.problemID,
            this.paginationService.getPageOffset(this.pageIndex, this.pageSize),
            this.pageSize
          );
        this.totalSubmissionCount = totalSubmissionCount;
        this.submissionSnippetList = submissionSnippetList;
      }

      this.lastLoadedTime = Date.now();
    } catch (e) {
      if (e instanceof UnauthenticatedError) {
        this.nzNotificationService.error(
          'Failed to load submission list',
          'Not logged in'
        );
        this.router.navigateByUrl('/login');
        return;
      }

      if (e instanceof PermissionDeniedError) {
        this.nzNotificationService.error(
          'Failed to load submission list',
          'Permission denied'
        );
        this.location.back();
        return;
      }

      if (e instanceof AccountNotFoundError) {
        this.nzNotificationService.error(
          'Failed to load submission list',
          'Account not found'
        );
        this.location.back();
        return;
      }

      if (e instanceof ProblemNotFoundError) {
        this.nzNotificationService.error(
          'Failed to load submission list',
          'Problem not found'
        );
        this.location.back();
        return;
      }

      if (e instanceof InvalidSubmissionListParam) {
        this.nzNotificationService.error(
          'Failed to load submission list',
          'Invalid page index/size'
        );
        this.location.back();
        return;
      }

      this.nzNotificationService.error(
        'Failed to load submission list',
        'Unknown error'
      );
      this.location.back();
    } finally {
      this.loading = false;
    }
  }

  public async onAutoReloadEnabledChange(enabled: boolean): Promise<void> {
    if (enabled) {
      if (this.submissionListReloadInterval !== undefined) {
        return;
      }

      await this.loadSubmissionSnippetList();
      this.submissionListReloadInterval = setInterval(async () => {
        await this.loadSubmissionSnippetList();
      }, SUBMISSION_LIST_RELOAD_INTERVAL);
      return;
    }

    if (this.submissionListReloadInterval === undefined) {
      return;
    }
    clearInterval(this.submissionListReloadInterval);
    this.submissionListReloadInterval = undefined;
  }

  public async onPageIndexChange(index: number): Promise<void> {
    this.pageIndex = index;
    await this.loadSubmissionSnippetList();
  }

  public async onPageSizeChange(size: number): Promise<void> {
    this.pageSize = size;
    await this.loadSubmissionSnippetList();
  }
}
