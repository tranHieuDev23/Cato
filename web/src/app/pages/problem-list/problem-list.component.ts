import { Component, OnInit } from '@angular/core';
import { RpcAccount, RpcProblemSnippet } from '../../dataaccess/api';
import { NzTableModule } from 'ng-zorro-antd/table';
import { ActivatedRoute, Params, Router, RouterModule } from '@angular/router';
import { CommonModule, Location } from '@angular/common';
import {
  InvalidProblemListParam,
  ProblemService,
} from '../../logic/problem.service';
import {
  NzNotificationModule,
  NzNotificationService,
} from 'ng-zorro-antd/notification';
import { PaginationService } from '../../logic/pagination.service';
import {
  AccountService,
  PermissionDeniedError,
  UnauthenticatedError,
} from '../../logic/account.service';
import { NzButtonModule } from 'ng-zorro-antd/button';
import { NzTypographyModule } from 'ng-zorro-antd/typography';
import { PageTitleService } from '../../logic/page-title.service';

const DEFAULT_PAGE_INDEX = 1;
const DEFAULT_PAGE_SIZE = 10;

@Component({
  selector: 'app-problem-list',
  standalone: true,
  imports: [
    NzTypographyModule,
    NzTableModule,
    RouterModule,
    CommonModule,
    NzNotificationModule,
    NzButtonModule,
  ],
  templateUrl: './problem-list.component.html',
  styleUrl: './problem-list.component.scss',
})
export class ProblemListComponent implements OnInit {
  public sessionAccount: RpcAccount | undefined;
  public problemSnippetList: RpcProblemSnippet[] = [];
  public totalProblemCount = 0;
  public pageIndex = DEFAULT_PAGE_INDEX;
  public pageSize = DEFAULT_PAGE_SIZE;
  public loading = false;

  constructor(
    private readonly accountService: AccountService,
    private readonly problemService: ProblemService,
    private readonly paginationService: PaginationService,
    private readonly activatedRoute: ActivatedRoute,
    private readonly router: Router,
    private readonly notificationService: NzNotificationService,
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
      this.pageTitleService.setTitle('Problems');
    })().then();
    this.activatedRoute.queryParams.subscribe(async (params) => {
      await this.onQueryParamsChanged(params);
    });
  }

  private async onQueryParamsChanged(queryParams: Params): Promise<void> {
    this.getPaginationInfoFromQueryParams(queryParams);
    await this.loadProblemSnippetList();
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

  private async loadProblemSnippetList(): Promise<void> {
    try {
      this.loading = true;
      const { totalProblemCount, problemSnippetList } =
        await this.problemService.getProblemSnippetList(
          this.paginationService.getPageOffset(this.pageIndex, this.pageSize),
          this.pageSize
        );
      this.totalProblemCount = totalProblemCount;
      this.problemSnippetList = problemSnippetList;
    } catch (e) {
      if (e instanceof UnauthenticatedError) {
        this.notificationService.error(
          'Failed to load problem list',
          'Not logged in'
        );
        this.router.navigateByUrl('/login');
        return;
      }

      if (e instanceof PermissionDeniedError) {
        this.notificationService.error(
          'Failed to load problem list',
          'Permission denied'
        );
        this.location.back();
        return;
      }

      if (e instanceof InvalidProblemListParam) {
        this.notificationService.error(
          'Failed to load problem list',
          'Invalid page index/size'
        );
        this.location.back();
        return;
      }

      this.notificationService.error(
        'Failed to load problem list',
        'Unknown error'
      );
      this.location.back();
    } finally {
      this.loading = false;
    }
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
    this.router.navigate(['/problem-list'], { queryParams });
  }
}
