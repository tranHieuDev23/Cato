import { Component } from '@angular/core';
import { NzTypographyModule } from 'ng-zorro-antd/typography';
import { RpcAccount, RpcProblem } from '../../dataaccess/api';
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
  ],
  templateUrl: './problem.component.html',
  styleUrl: './problem.component.scss',
})
export class ProblemComponent {
  public sessionAccount: RpcAccount | null | undefined;
  public problem: RpcProblem | undefined;

  public submissionContent = '';

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
}
