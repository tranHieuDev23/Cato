import { Component, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute, Router, Params, RouterModule } from '@angular/router';
import {
  NzNotificationModule,
  NzNotificationService,
} from 'ng-zorro-antd/notification';
import { RpcAccount, RpcSubmissionSnippet } from '../../dataaccess/api';
import {
  AccountService,
  UnauthenticatedError,
  PermissionDeniedError,
  AccountNotFoundError,
} from '../../logic/account.service';
import { PaginationService } from '../../logic/pagination.service';
import {
  SubmissionService,
  InvalidSubmissionListParam,
  SubmissionNotFoundError,
} from '../../logic/submission.service';
import { CommonModule, Location } from '@angular/common';
import { NzTableModule } from 'ng-zorro-antd/table';
import { NzTypographyModule } from 'ng-zorro-antd/typography';
import { NzSwitchModule } from 'ng-zorro-antd/switch';
import { FormsModule } from '@angular/forms';
import { LanguagePipe } from '../../components/utils/language.pipe';
import { SubmissionStatusPipe } from '../../components/utils/submission-status.pipe';
import { PageTitleService } from '../../logic/page-title.service';
import { SubmissionModalComponent } from '../../components/submission-modal/submission-modal.component';
import { NzModalModule, NzModalService } from 'ng-zorro-antd/modal';
import { NzTagModule } from 'ng-zorro-antd/tag';
import { SubmissionStatusColorPipe } from '../../components/utils/submission-status-color.pipe';
import { Subscription } from 'rxjs';

const DEFAULT_PAGE_INDEX = 1;
const DEFAULT_PAGE_SIZE = 10;
const SUBMISSION_LIST_RELOAD_INTERVAL = 10000;

@Component({
  selector: 'app-submission-list',
  standalone: true,
  templateUrl: './submission-list.component.html',
  styleUrl: './submission-list.component.scss',
  imports: [
    NzTableModule,
    CommonModule,
    NzTypographyModule,
    NzSwitchModule,
    FormsModule,
    NzNotificationModule,
    LanguagePipe,
    RouterModule,
    SubmissionStatusPipe,
    NzModalModule,
    SubmissionStatusColorPipe,
    NzTagModule,
  ],
})
export class SubmissionListComponent implements OnInit, OnDestroy {
  public sessionAccount: RpcAccount | undefined;
  public submissionSnippetList: RpcSubmissionSnippet[] = [];
  public totalSubmissionCount = 0;
  public pageIndex = DEFAULT_PAGE_INDEX;
  public pageSize = DEFAULT_PAGE_SIZE;
  public loading = false;
  public lastLoadedTime: number | undefined;
  public autoReloadEnabled = true;

  private queryParamsSubscription: Subscription | undefined;
  private submissionListReloadInterval:
    | ReturnType<typeof setInterval>
    | undefined;

  constructor(
    private readonly accountService: AccountService,
    private readonly submissionService: SubmissionService,
    private readonly paginationService: PaginationService,
    private readonly activatedRoute: ActivatedRoute,
    private readonly router: Router,
    private readonly notificationService: NzNotificationService,
    private readonly modalService: NzModalService,
    private readonly location: Location,
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
      this.pageTitleService.setTitle('Submissions');
    })().then();
    this.queryParamsSubscription = this.activatedRoute.queryParams.subscribe(
      async (params) => {
        await this.onQueryParamsChanged(params);
      }
    );
    this.submissionListReloadInterval = setInterval(async () => {
      await this.loadSubmissionSnippetList();
    }, SUBMISSION_LIST_RELOAD_INTERVAL);
  }

  ngOnDestroy(): void {
    this.queryParamsSubscription?.unsubscribe();
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
      const { totalSubmissionCount, submissionSnippetList } =
        await this.submissionService.getSubmissionSnippetList(
          this.paginationService.getPageOffset(this.pageIndex, this.pageSize),
          this.pageSize
        );
      this.totalSubmissionCount = totalSubmissionCount;
      this.submissionSnippetList = submissionSnippetList;
      this.lastLoadedTime = Date.now();
    } catch (e) {
      if (e instanceof UnauthenticatedError) {
        this.notificationService.error(
          'Failed to load submission list',
          'Not logged in'
        );
        this.router.navigateByUrl('/login');
        return;
      }

      if (e instanceof PermissionDeniedError) {
        this.notificationService.error(
          'Failed to load submission list',
          'Permission denied'
        );
        this.location.back();
        return;
      }

      if (e instanceof AccountNotFoundError) {
        this.notificationService.error(
          'Failed to load submission list',
          'Account not found'
        );
        this.location.back();
        return;
      }

      if (e instanceof InvalidSubmissionListParam) {
        this.notificationService.error(
          'Failed to load submission list',
          'Invalid page index/size'
        );
        this.location.back();
        return;
      }

      this.notificationService.error(
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

  public async onSubmissionSnippetClicked(
    submissionSnippet: RpcSubmissionSnippet
  ): Promise<void> {
    try {
      const submission = await this.submissionService.getSubmission(
        submissionSnippet.iD
      );
      this.modalService.create({
        nzContent: SubmissionModalComponent,
        nzData: { submission },
        nzWidth: 'fit-content',
        nzTitle: `Submission #${submission.iD}`,
        nzFooter: null,
      });
    } catch (e) {
      if (e instanceof UnauthenticatedError) {
        this.notificationService.error(
          'Failed to load submission',
          'Not logged in'
        );
        this.router.navigateByUrl('/login');
        return;
      }

      if (e instanceof PermissionDeniedError) {
        this.notificationService.error(
          'Failed to load submission',
          'Permission denied'
        );
        this.location.back();
        return;
      }

      if (e instanceof SubmissionNotFoundError) {
        this.notificationService.error(
          'Failed to load submission',
          'Submission not found'
        );
        await this.loadSubmissionSnippetList();
        return;
      }
    }
  }
}
