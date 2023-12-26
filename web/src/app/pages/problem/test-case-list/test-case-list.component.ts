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
import { RpcAccount, RpcProblem } from '../../../dataaccess/api';
import { CodemirrorComponent, CodemirrorModule } from '@ctrl/ngx-codemirror';
import { FormsModule } from '@angular/forms';
import { NzTypographyModule } from 'ng-zorro-antd/typography';
import { NzCheckboxModule } from 'ng-zorro-antd/checkbox';
import { NzModalModule, NzModalRef, NzModalService } from 'ng-zorro-antd/modal';
import { NzToolTipModule } from 'ng-zorro-antd/tooltip';
import copyToClipboard from 'copy-to-clipboard';
import { EllipsisPipe } from '../../../components/utils/ellipsis.pipe';
import { NzUploadFile, NzUploadModule } from 'ng-zorro-antd/upload';

export interface TestCaseListItem {
  uuid: string;
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
    NzToolTipModule,
    EllipsisPipe,
    NzUploadModule,
  ],
  templateUrl: './test-case-list.component.html',
  styleUrl: './test-case-list.component.scss',
})
export class TestCaseListComponent implements OnInit {
  @ViewChildren(CodemirrorComponent)
  codeMirrorComponentList!: QueryList<CodemirrorComponent>;

  @ViewChild('expandTestCaseModal') expandTestCaseModal:
    | TemplateRef<any>
    | undefined;
  @ViewChild('expandTestCaseModalFooter') expandTestCaseModalFooter:
    | TemplateRef<any>
    | undefined;
  @ViewChild('editTestCaseModal') editTestCaseModal:
    | TemplateRef<any>
    | undefined;
  @ViewChild('uploadTestCaseModal') uploadTestCaseModal:
    | TemplateRef<any>
    | undefined;

  @Input() public problem: RpcProblem | undefined;

  public sessionAccount: RpcAccount | undefined;

  public totalTestCaseCount = 0;
  public testCaseList: TestCaseListItem[] = [];
  public pageIndex = 1;
  public pageSize = 10;
  public loading = false;

  public expandTestCaseModalInput = '';
  public expandTestCaseModalOutput = '';

  public editTestCaseModalInput = '';
  public editTestCaseModalOutput = '';
  public editModalTestCaseIsHidden = true;

  private uploadTestCaseModalRef: NzModalRef | undefined;

  constructor(
    private readonly accountService: AccountService,
    private readonly testCaseService: TestCaseService,
    public readonly paginationService: PaginationService,
    private readonly router: Router,
    private readonly notificationService: NzNotificationService,
    private readonly modalService: NzModalService,
    private readonly location: Location
  ) {}

  ngOnInit(): void {
    this.loadTestCaseSnippetList().then();
  }

  private async loadTestCaseSnippetList(): Promise<void> {
    if (!this.problem) {
      return;
    }

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
          this.problem.uUID,
          this.paginationService.getPageOffset(this.pageIndex, this.pageSize),
          this.pageSize
        );
      this.totalTestCaseCount = totalTestCaseCount;
      this.testCaseList = testCaseSnippetList.map((item) => {
        return {
          uuid: item.uUID,
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

  public async onTestCaseInputClicked(
    index: number,
    testCaseListItem: TestCaseListItem
  ): Promise<void> {
    if (testCaseListItem.loading) {
      return;
    }
    if (testCaseListItem.isSnippet) {
      await this.loadTestCaseFromSnippet(index, testCaseListItem);
    }
    copyToClipboard(this.testCaseList[index].input);
    this.notificationService.success('Input copied to clipboard', '');
  }

  public async onTestCaseOutputClicked(
    index: number,
    testCaseListItem: TestCaseListItem
  ): Promise<void> {
    if (testCaseListItem.loading) {
      return;
    }
    if (testCaseListItem.isSnippet) {
      await this.loadTestCaseFromSnippet(index, testCaseListItem);
    }
    copyToClipboard(this.testCaseList[index].output);
    this.notificationService.success('Output copied to clipboard', '');
  }

  public async onTestCaseExpandClicked(
    index: number,
    testCaseListItem: TestCaseListItem
  ): Promise<void> {
    if (testCaseListItem.loading) {
      return;
    }

    if (testCaseListItem.isSnippet) {
      await this.loadTestCaseFromSnippet(index, testCaseListItem);
    }

    const testCase = this.testCaseList[index];
    this.expandTestCaseModalInput = testCase.input;
    this.expandTestCaseModalOutput = testCase.output;
    this.modalService.create({
      nzContent: this.expandTestCaseModal,
      nzWidth: 'fit-content',
      nzFooter: this.expandTestCaseModalFooter,
    });
  }

  public onExpandTestCaseModalCopyInputClicked(): void {
    copyToClipboard(this.expandTestCaseModalInput);
    this.notificationService.success('Input copied to clipboard', '');
  }

  public onExpandTestCaseModalCopyOutputClicked(): void {
    copyToClipboard(this.expandTestCaseModalInput);
    this.notificationService.success('Output copied to clipboard', '');
  }

  public async onTestCaseEditClicked(
    index: number,
    testCaseListItem: TestCaseListItem
  ): Promise<void> {
    if (testCaseListItem.loading) {
      return;
    }

    if (testCaseListItem.isSnippet) {
      await this.loadTestCaseFromSnippet(index, testCaseListItem);
    }

    const testCase = this.testCaseList[index];
    this.editTestCaseModalInput = testCase.input;
    this.editTestCaseModalOutput = testCase.output;
    this.editModalTestCaseIsHidden = testCase.isHidden;

    this.modalService.create({
      nzContent: this.editTestCaseModal,
      nzWidth: 'fit-content',
      nzOnOk: async () => {
        try {
          await this.testCaseService.updateTestCase(
            testCase.uuid,
            this.editTestCaseModalInput,
            this.editTestCaseModalOutput,
            this.editModalTestCaseIsHidden
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
              'Test case not found'
            );
            await this.loadTestCaseSnippetList();
            return;
          }

          this.notificationService.error(
            'Failed to update test case',
            'Unknown error'
          );
        }
      },
    });
  }

  public async onTestCaseDeleteClicked(
    testCaseListItem: TestCaseListItem
  ): Promise<void> {
    if (testCaseListItem.loading) {
      return;
    }

    const testCaseID = testCaseListItem.uuid;
    this.modalService.create({
      nzContent: 'Are you sure? This action is <b>irreversible</b>',
      nzOkDanger: true,
      nzCloseIcon: '',
      nzOnOk: async () => {
        try {
          await this.testCaseService.deleteTestCase(testCaseID);
          this.notificationService.success(
            'Deleted test case successfully',
            ''
          );
          await this.loadTestCaseSnippetList();
        } catch (e) {
          if (e instanceof UnauthenticatedError) {
            this.notificationService.error(
              'Failed to delete test case',
              'Not logged in'
            );
            this.router.navigateByUrl('/login');
            return;
          }

          if (e instanceof PermissionDeniedError) {
            this.notificationService.error(
              'Failed to delete test case',
              'Permission denied'
            );
            this.location.back();
            return;
          }

          if (e instanceof ProblemNotFoundError) {
            this.notificationService.error(
              'Failed to delete test case',
              'Test case not found'
            );
            await this.loadTestCaseSnippetList();
            return;
          }

          this.notificationService.error(
            'Failed to delete test case',
            'Unknown error'
          );
        }
      },
    });
  }

  public onCreateTestCaseClicked(): void {
    if (!this.problem) {
      return;
    }

    this.editTestCaseModalInput = '';
    this.editTestCaseModalOutput = '';
    this.editModalTestCaseIsHidden = true;

    const problemID = this.problem.uUID;
    this.modalService.create({
      nzContent: this.editTestCaseModal,
      nzWidth: 'fit-content',
      nzOnOk: async () => {
        try {
          await this.testCaseService.createTestCase(
            problemID,
            this.editTestCaseModalInput,
            this.editTestCaseModalOutput,
            this.editModalTestCaseIsHidden
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
              'Problem not found'
            );
            await this.loadTestCaseSnippetList();
            return;
          }
        }
      },
    });
  }

  public onUploadZippedTestCasesClicked(): void {
    this.uploadTestCaseModalRef?.close();
    this.uploadTestCaseModalRef = this.modalService.create({
      nzContent: this.uploadTestCaseModal,
      nzWidth: 'fit-content',
      nzCloseIcon: '',
      nzOkText: null,
    });
  }

  public onLoadFile = (file: NzUploadFile): boolean => {
    if (!this.problem) {
      return false;
    }

    const problemID = this.problem.uUID;
    const fileReader = new FileReader();
    fileReader.onload = async (event) => {
      this.uploadTestCaseModalRef?.close();
      try {
        const content = event.target?.result as ArrayBuffer;
        await this.testCaseService.createTestCaseList(problemID, content);
        this.notificationService.success(
          'Uploaded zipped test cases successfully',
          ''
        );
        await this.loadTestCaseSnippetList();
      } catch (e) {
        if (e instanceof UnauthenticatedError) {
          this.notificationService.error(
            'Failed to upload zipped test cases',
            'Not logged in'
          );
          this.router.navigateByUrl('/login');
          return;
        }

        if (e instanceof PermissionDeniedError) {
          this.notificationService.error(
            'Failed to upload zipped test cases',
            'Permission denied'
          );
          this.location.back();
          return;
        }

        if (e instanceof ProblemNotFoundError) {
          this.notificationService.error(
            'Failed to upload zipped test cases',
            'Problem not found'
          );
          await this.loadTestCaseSnippetList();
          return;
        }

        this.notificationService.error(
          'Failed to upload zipped test cases',
          'Unknown error'
        );
      }
    };
    fileReader.readAsArrayBuffer(file as any);
    return false;
  };

  private async loadTestCaseFromSnippet(
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
        testCaseListItem.uuid
      );
      this.testCaseList = [...this.testCaseList];
      this.testCaseList[index] = {
        uuid: testCase.uUID,
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
          'Test case not found'
        );
        await this.loadTestCaseSnippetList();
        return;
      }

      this.notificationService.error(
        'Failed to load test case',
        'Unknown error'
      );
    } finally {
      this.testCaseList = [...this.testCaseList];
      this.testCaseList[index] = {
        ...this.testCaseList[index],
        loading: false,
      };
    }
  }
}
