import { Component, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute, Router, Params } from '@angular/router';
import {
  NzNotificationModule,
  NzNotificationService,
} from 'ng-zorro-antd/notification';
import { Subscription } from 'rxjs';
import { RpcAccount, RpcSubmissionSnippet } from '../../dataaccess/api';
import {
  AccountService,
  UnauthenticatedError,
  PermissionDeniedError,
  Role,
  AccountNotFoundError,
} from '../../logic/account.service';
import { PaginationService } from '../../logic/pagination.service';
import {
  SubmissionService,
  InvalidSubmissionListParam,
} from '../../logic/submission.service';
import { CommonModule, Location } from '@angular/common';
import { NzTableModule } from 'ng-zorro-antd/table';
import { NzTypographyModule } from 'ng-zorro-antd/typography';
import { NzSwitchModule } from 'ng-zorro-antd/switch';
import { FormsModule } from '@angular/forms';

const DEFAULT_PAGE_INDEX = 1;
const DEFAULT_PAGE_SIZE = 10;
const SUBMISSION_LIST_RELOAD_INTERVAL = 10000;

@Component({
  selector: 'app-submission-list',
  standalone: true,
  imports: [
    NzTableModule,
    CommonModule,
    NzTypographyModule,
    NzSwitchModule,
    FormsModule,
    NzNotificationModule,
  ],
  templateUrl: './submission-list.component.html',
  styleUrl: './submission-list.component.scss',
})
export class SubmissionListComponent implements OnInit, OnDestroy {
  public sessionAccount: RpcAccount | null | undefined;
  public submissionSnippetList: RpcSubmissionSnippet[] = [];
  public totalSubmissionCount = 0;
  public pageIndex = DEFAULT_PAGE_INDEX;
  public pageSize = DEFAULT_PAGE_SIZE;
  public loading = false;
  public lastLoadedTime: number | undefined;
  public autoReloadEnabled = true;

  private sessionAccountChangedSubscription: Subscription | undefined;
  private submissionListReloadInterval:
    | ReturnType<typeof setInterval>
    | undefined;

  constructor(
    private readonly accountService: AccountService,
    private readonly submissionService: SubmissionService,
    private readonly paginationService: PaginationService,
    private readonly activatedRoute: ActivatedRoute,
    private readonly router: Router,
    private readonly nzNotificationService: NzNotificationService,
    private readonly location: Location
  ) {}

  ngOnInit(): void {
    (async () => {
      this.sessionAccount = await this.accountService.getSessionAccount();
    })().then();
    this.activatedRoute.queryParams.subscribe(async (params) => {
      await this.onQueryParamsChanged(params);
    });
    this.sessionAccountChangedSubscription =
      this.accountService.sessionAccountChanged.subscribe((account) => {
        this.sessionAccount = account;
      });
    this.submissionListReloadInterval = setInterval(async () => {
      await this.loadSubmissionSnippetList();
    }, SUBMISSION_LIST_RELOAD_INTERVAL);
  }

  ngOnDestroy(): void {
    this.sessionAccountChangedSubscription?.unsubscribe();
    if (this.submissionListReloadInterval !== undefined) {
      clearInterval(this.submissionListReloadInterval);
    }
  }

  private async onQueryParamsChanged(queryParams: Params): Promise<void> {
    this.getPaginationInfoFromQueryParams(queryParams);
    await this.loadSubmissionSnippetList();
  }

  private getPaginationInfoFromQueryParams(queryParams: Params): void {
    if (queryParams['index'] !== undefined) {
      this.pageIndex = +queryParams['index'];
    } else {
      this.pageIndex = DEFAULT_PAGE_INDEX;
    }
    if (queryParams['size'] !== undefined) {
      this.pageSize = +queryParams['size'];
    } else {
      this.pageSize = DEFAULT_PAGE_SIZE;
    }
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
          await this.submissionService.getSubmissionSnippetList(
            this.paginationService.getPageOffset(this.pageIndex, this.pageSize),
            this.pageSize
          );
        this.totalSubmissionCount = totalSubmissionCount;
        this.submissionSnippetList = submissionSnippetList;
      } else {
        const { totalSubmissionCount, submissionSnippetList } =
          await this.submissionService.getAccountSubmissionSnippetList(
            sessionAccount.iD,
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

  public onAutoReloadEnabledChange(enabled: boolean): void {
    if (enabled) {
      if (this.submissionListReloadInterval !== undefined) {
        return;
      }

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
    this.navigateToPage(index, this.pageSize);
  }

  public async onPageSizeChange(size: number): Promise<void> {
    this.navigateToPage(this.pageIndex, size);
  }

  private navigateToPage(index: number, size: number): void {
    const queryParams: any = {};
    if (index !== DEFAULT_PAGE_INDEX) {
      queryParams['index'] = index;
    }
    if (size !== DEFAULT_PAGE_SIZE) {
      queryParams['size'] = size;
    }
    this.router.navigate(['/submission-list'], { queryParams });
  }
}
