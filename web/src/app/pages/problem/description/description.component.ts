import {
  Component,
  ElementRef,
  Input,
  OnChanges,
  OnInit,
  TemplateRef,
  ViewChild,
} from '@angular/core';
import {
  RpcAccount,
  RpcGetServerInfoResponse,
  RpcProblem,
  RpcProblemExample,
} from '../../../dataaccess/api';
import { NzDescriptionsModule } from 'ng-zorro-antd/descriptions';
import { CommonModule, Location } from '@angular/common';
import { TimeLimitPipe } from '../../../components/utils/time-limit.pipe';
import { MemoryPipe } from '../../../components/utils/memory.pipe';
import { Router, RouterModule } from '@angular/router';
import { NzButtonModule } from 'ng-zorro-antd/button';
import { NzTableModule } from 'ng-zorro-antd/table';
import { NzTypographyModule } from 'ng-zorro-antd/typography';
import { NzModalService } from 'ng-zorro-antd/modal';
import {
  ProblemNotFoundError,
  ProblemService,
} from '../../../logic/problem.service';
import { NzNotificationService } from 'ng-zorro-antd/notification';
import {
  UnauthenticatedError,
  PermissionDeniedError,
  AccountService,
  Role,
} from '../../../logic/account.service';
import { NzIconModule } from 'ng-zorro-antd/icon';
import renderMathInElement from 'katex/contrib/auto-render';
import { EllipsisPipe } from '../../../components/utils/ellipsis.pipe';
import { NzToolTipModule } from 'ng-zorro-antd/tooltip';
import copyToClipboard from 'copy-to-clipboard';
import { TestCaseViewModalComponent } from '../../../components/test-case-view-modal/test-case-view-modal.component';
import { ServerService } from '../../../logic/server.service';

@Component({
  selector: 'app-description',
  standalone: true,
  imports: [
    CommonModule,
    NzDescriptionsModule,
    TimeLimitPipe,
    MemoryPipe,
    RouterModule,
    NzButtonModule,
    NzTableModule,
    NzTypographyModule,
    NzIconModule,
    EllipsisPipe,
    NzToolTipModule,
    TestCaseViewModalComponent,
  ],
  templateUrl: './description.component.html',
  styleUrl: './description.component.scss',
})
export class DescriptionComponent implements OnInit, OnChanges {
  @ViewChild('problemDescriptionContainer')
  problemDescriptionContainer: ElementRef | undefined;
  @ViewChild('expandTestCaseModal') expandExampleModal:
    | TemplateRef<any>
    | undefined;
  @ViewChild('expandTestCaseModalFooter') expandExampleModalFooter:
    | TemplateRef<any>
    | undefined;

  public expandExampleModalInput = '';
  public expandExampleModalOutput = '';

  @Input() public problem!: RpcProblem;

  public sessionAccount: RpcAccount | undefined;
  public serverInfo: RpcGetServerInfoResponse | undefined;

  constructor(
    private readonly modalService: NzModalService,
    private readonly problemService: ProblemService,
    private readonly notificationService: NzNotificationService,
    private readonly location: Location,
    private readonly router: Router,
    private readonly accountService: AccountService,
    private readonly serverService: ServerService
  ) {}

  ngOnInit(): void {
    setTimeout(() => this.renderKatex());
    (async () => {
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

      try {
        this.serverInfo = await this.serverService.getServerInfo();
      } catch (e) {
        this.notificationService.error('Failed to get server information', '');
        this.location.back();
        return;
      }
    })().then();
  }

  ngOnChanges(): void {
    setTimeout(() => this.renderKatex());
  }

  private renderKatex(): void {
    if (!this.problemDescriptionContainer) {
      return;
    }

    renderMathInElement(this.problemDescriptionContainer.nativeElement, {
      throwOnError: false,
    });
  }

  public canUpdateProblem(): boolean {
    if (this.serverInfo?.setting.problem.disableProblemUpdate) {
      return false;
    }
    if (!this.sessionAccount || !this.problem) {
      return false;
    }
    if (this.sessionAccount.role === Role.Admin) {
      return true;
    }
    if (
      this.sessionAccount.role === Role.ProblemSetter &&
      this.sessionAccount.iD === this.problem.author.iD
    ) {
      return true;
    }
    return false;
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

  public onExampleExpandClicked(example: RpcProblemExample): void {
    this.expandExampleModalInput = example.input;
    this.expandExampleModalOutput = example.output;
    this.modalService.create({
      nzContent: this.expandExampleModal,
      nzWidth: 'fit-content',
      nzFooter: this.expandExampleModalFooter,
    });
  }

  public onExpandTestCaseModalCopyInputClicked(): void {
    copyToClipboard(this.expandExampleModalInput);
    this.notificationService.success('Input copied to clipboard', '');
  }

  public onExpandTestCaseModalCopyOutputClicked(): void {
    copyToClipboard(this.expandExampleModalOutput);
    this.notificationService.success('Output copied to clipboard', '');
  }
}
