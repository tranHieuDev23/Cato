import { CommonModule, Location } from '@angular/common';
import {
  Component,
  Input,
  OnInit,
  QueryList,
  TemplateRef,
  ViewChild,
  ViewChildren,
} from '@angular/core';
import { Router } from '@angular/router';
import {
  NzNotificationModule,
  NzNotificationService,
} from 'ng-zorro-antd/notification';
import {
  UnauthenticatedError,
  PermissionDeniedError,
  AccountNotFoundError,
  AccountService,
} from '../../../logic/account.service';
import { PaginationService } from '../../../logic/pagination.service';
import { ProblemNotFoundError } from '../../../logic/problem.service';
import {
  InvalidTestCaseListParam,
  TestCaseNotFoundError,
  TestCaseService,
} from '../../../logic/test-case.service';
import { NzTableModule } from 'ng-zorro-antd/table';
import { NzButtonModule } from 'ng-zorro-antd/button';
import { NzIconModule } from 'ng-zorro-antd/icon';
import { RpcAccount } from '../../../dataaccess/api';
import { CodemirrorComponent, CodemirrorModule } from '@ctrl/ngx-codemirror';
import { FormsModule } from '@angular/forms';
import { NzTypographyModule } from 'ng-zorro-antd/typography';
import { NzCheckboxModule } from 'ng-zorro-antd/checkbox';
import { NzModalModule, NzModalService } from 'ng-zorro-antd/modal';

export interface TestCaseListItem {
  id: number;
  input: string;
  output: string;
  isHidden: boolean;
  loading: boolean;
  isSnippet: boolean;
}

@Component({
  selector: 'app-test-case-list',
  standalone: true,
  imports: [
    NzTableModule,
    CommonModule,
    NzNotificationModule,
    NzButtonModule,
    NzIconModule,
    CodemirrorModule,
    FormsModule,
    NzTypographyModule,
    NzModalModule,
    NzCheckboxModule,
  ],
  templateUrl: './test-case-list.component.html',
  styleUrl: './test-case-list.component.scss',
})
export class TestCaseListComponent implements OnInit {
  @ViewChildren(CodemirrorComponent)
  codeMirrorComponentList!: QueryList<CodemirrorComponent>;

  @ViewChild('editTestCaseModal') editTestCaseModal:
    | TemplateRef<any>
    | undefined;

  @Input() public problemID = 0;

  public sessionAccount: RpcAccount | undefined;

  public totalTestCaseCount = 0;
  public testCaseList: TestCaseListItem[] = [];
  public pageIndex = 1;
  public pageSize = 10;
  public loading = false;

  public modalTestCaseInput = '';
  public modalTestCaseOutput = '';
  public isModalTestCaseHidden = true;

  constructor(
    private readonly accountService: AccountService,
    private readonly testCaseService: TestCaseService,
    private readonly paginationService: PaginationService,
    private readonly router: Router,
    private readonly notificationService: NzNotificationService,
    private readonly modalService: NzModalService,
    private readonly location: Location
  ) {}

  ngOnInit(): void {
    this.loadTestCaseSnippetList().then();
  }

  private async loadTestCaseSnippetList(): Promise<void> {
    try {
      this.loading = true;

      const sessionAccount = await this.accountService.getSessionAccount();
      if (sessionAccount === null) {
        this.notificationService.error(
          'Failed to get session account',
          'Not logged in'
        );
        this.router.navigateByUrl('/login');
        return;
      }
      this.sessionAccount = sessionAccount;

      const { totalTestCaseCount, testCaseSnippetList } =
        await this.testCaseService.getProblemTestCaseSnippetList(
          this.problemID,
          this.paginationService.getPageOffset(this.pageIndex, this.pageSize),
          this.pageSize
        );
      this.totalTestCaseCount = totalTestCaseCount;
      this.testCaseList = testCaseSnippetList.map((item) => {
        return {
          id: item.iD,
          input: item.input,
          output: item.output,
          isHidden: item.isHidden,
          loading: false,
          isSnippet: true,
        };
      });
    } catch (e) {
      if (e instanceof UnauthenticatedError) {
        this.notificationService.error(
          'Failed to load test case list',
          'Not logged in'
        );
        this.router.navigateByUrl('/login');
        return;
      }

      if (e instanceof PermissionDeniedError) {
        this.notificationService.error(
          'Failed to load test case list',
          'Permission denied'
        );
        this.location.back();
        return;
      }

      if (e instanceof AccountNotFoundError) {
        this.notificationService.error(
          'Failed to load test case list',
          'Account not found'
        );
        this.location.back();
        return;
      }

      if (e instanceof ProblemNotFoundError) {
        this.notificationService.error(
          'Failed to load test case list',
          'Problem not found'
        );
        this.location.back();
        return;
      }

      if (e instanceof InvalidTestCaseListParam) {
        this.notificationService.error(
          'Failed to load test case list',
          'Invalid page index/size'
        );
        this.location.back();
        return;
      }

      this.notificationService.error(
        'Failed to load test case list',
        'Unknown error'
      );
      this.location.back();
    } finally {
      this.loading = false;
    }
  }

  public async onPageIndexChange(index: number): Promise<void> {
    this.pageIndex = index;
    await this.loadTestCaseSnippetList();
  }

  public async onPageSizeChange(size: number): Promise<void> {
    this.pageSize = size;
    await this.loadTestCaseSnippetList();
  }

  public async expandTestCase(
    index: number,
    testCaseListItem: TestCaseListItem
  ): Promise<void> {
    if (testCaseListItem.loading || !testCaseListItem.isSnippet) {
      return;
    }

    this.testCaseList = [...this.testCaseList];
    this.testCaseList[index].loading = true;

    try {
      const testCase = await this.testCaseService.getTestCase(
        testCaseListItem.id
      );
      this.testCaseList = [...this.testCaseList];
      this.testCaseList[index] = {
        id: testCase.iD,
        input: testCase.input,
        output: testCase.output,
        isHidden: testCase.isHidden,
        loading: false,
        isSnippet: false,
      };
    } catch (e) {
      if (e instanceof UnauthenticatedError) {
        this.notificationService.error(
          'Failed to load test case',
          'Not logged in'
        );
        this.router.navigateByUrl('/login');
        return;
      }

      if (e instanceof PermissionDeniedError) {
        this.notificationService.error(
          'Failed to load test case',
          'Permission denied'
        );
        this.location.back();
        return;
      }

      if (e instanceof TestCaseNotFoundError) {
        this.notificationService.error(
          'Failed to load test case',
          'TestCase not found'
        );
        await this.loadTestCaseSnippetList();
        return;
      }
    } finally {
      this.testCaseList = [...this.testCaseList];
      this.testCaseList[index] = {
        ...this.testCaseList[index],
        loading: false,
      };
    }
  }

  public async onTestCaseEditClicked(
    index: number,
    testCaseListItem: TestCaseListItem
  ): Promise<void> {
    if (testCaseListItem.loading) {
      return;
    }

    if (testCaseListItem.isSnippet) {
      await this.expandTestCase(index, testCaseListItem);
    }

    const testCase = this.testCaseList[index];
    this.modalTestCaseInput = testCase.input;
    this.modalTestCaseOutput = testCase.output;
    this.isModalTestCaseHidden = testCase.isHidden;

    this.modalService.create({
      nzContent: this.editTestCaseModal,
      nzWidth: 'fit-content',
      nzOnOk: async () => {
        try {
          await this.testCaseService.updateTestCase(
            testCase.id,
            this.modalTestCaseInput,
            this.modalTestCaseOutput,
            this.isModalTestCaseHidden
          );
          this.notificationService.success(
            'Updated test case successfully',
            ''
          );
          await this.loadTestCaseSnippetList();
        } catch (e) {
          if (e instanceof UnauthenticatedError) {
            this.notificationService.error(
              'Failed to update test case',
              'Not logged in'
            );
            this.router.navigateByUrl('/login');
            return;
          }

          if (e instanceof PermissionDeniedError) {
            this.notificationService.error(
              'Failed to update test case',
              'Permission denied'
            );
            this.location.back();
            return;
          }

          if (e instanceof ProblemNotFoundError) {
            this.notificationService.error(
              'Failed to update test case',
              'TestCase not found'
            );
            await this.loadTestCaseSnippetList();
            return;
          }
        }
      },
    });
  }

  public onCreateTestCaseClicked(): void {
    this.modalTestCaseInput = '';
    this.modalTestCaseOutput = '';
    this.isModalTestCaseHidden = true;
    this.modalService.create({
      nzContent: this.editTestCaseModal,
      nzWidth: 'fit-content',
      nzOnOk: async () => {
        try {
          await this.testCaseService.createTestCase(
            this.problemID,
            this.modalTestCaseInput,
            this.modalTestCaseOutput,
            this.isModalTestCaseHidden
          );
          this.notificationService.success(
            'Created test case successfully',
            ''
          );
          await this.loadTestCaseSnippetList();
        } catch (e) {
          if (e instanceof UnauthenticatedError) {
            this.notificationService.error(
              'Failed to create test case',
              'Not logged in'
            );
            this.router.navigateByUrl('/login');
            return;
          }

          if (e instanceof PermissionDeniedError) {
            this.notificationService.error(
              'Failed to create test case',
              'Permission denied'
            );
            this.location.back();
            return;
          }

          if (e instanceof ProblemNotFoundError) {
            this.notificationService.error(
              'Failed to create test case',
              'TestCase not found'
            );
            await this.loadTestCaseSnippetList();
            return;
          }
        }
      },
    });
  }
}