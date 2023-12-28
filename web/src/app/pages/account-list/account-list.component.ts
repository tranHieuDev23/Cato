import {
  Component,
  OnDestroy,
  OnInit,
  TemplateRef,
  ViewChild,
} from '@angular/core';
import { ActivatedRoute, Router, Params } from '@angular/router';
import { NzButtonModule } from 'ng-zorro-antd/button';
import { NzNotificationService } from 'ng-zorro-antd/notification';
import { NzTableModule } from 'ng-zorro-antd/table';
import { RpcAccount, RpcGetServerInfoResponse } from '../../dataaccess/api';
import {
  AccountService,
  UnauthenticatedError,
  PermissionDeniedError,
  InvalidAccountListParam,
  Role,
  InvalidAccountInfoError,
} from '../../logic/account.service';
import { PageTitleService } from '../../logic/page-title.service';
import { PaginationService } from '../../logic/pagination.service';
import { CommonModule, Location } from '@angular/common';
import { NzIconModule } from 'ng-zorro-antd/icon';
import { RolePipe } from '../../components/utils/role.pipe';
import { NzFormModule } from 'ng-zorro-antd/form';
import {
  AbstractControl,
  FormBuilder,
  FormGroup,
  ReactiveFormsModule,
  ValidatorFn,
  Validators,
} from '@angular/forms';
import { ConfirmedValidator } from '../../components/utils/confirmed-validator';
import { NzModalModule, NzModalService } from 'ng-zorro-antd/modal';
import { NzRadioModule } from 'ng-zorro-antd/radio';
import { NzInputModule } from 'ng-zorro-antd/input';
import { NzSpaceModule } from 'ng-zorro-antd/space';
import { NzToolTipModule } from 'ng-zorro-antd/tooltip';
import { Subscription } from 'rxjs';
import { ServerService } from '../../logic/server.service';

const DEFAULT_PAGE_INDEX = 1;
const DEFAULT_PAGE_SIZE = 10;

@Component({
  selector: 'app-account-list',
  standalone: true,
  imports: [
    NzTableModule,
    NzButtonModule,
    CommonModule,
    NzIconModule,
    RolePipe,
    NzFormModule,
    ReactiveFormsModule,
    NzModalModule,
    NzRadioModule,
    NzInputModule,
    NzSpaceModule,
    NzToolTipModule,
  ],
  templateUrl: './account-list.component.html',
  styleUrl: './account-list.component.scss',
})
export class AccountListComponent implements OnInit, OnDestroy {
  @ViewChild('createAccountModal') createAccountModal:
    | TemplateRef<any>
    | undefined;
  @ViewChild('editAccountModal') editAccountModal: TemplateRef<any> | undefined;

  public sessionAccount: RpcAccount | undefined;
  public accountList: RpcAccount[] = [];
  public totalAccountCount = 0;
  public pageIndex = DEFAULT_PAGE_INDEX;
  public pageSize = DEFAULT_PAGE_SIZE;
  public loading = false;

  public createAccountForm: FormGroup;
  public editAccountForm: FormGroup;

  public serverInfo: RpcGetServerInfoResponse | undefined;

  private queryParamsSubscription: Subscription | undefined;
  private serverInfoChangedSubscription: Subscription;

  constructor(
    private readonly accountService: AccountService,
    private readonly paginationService: PaginationService,
    private readonly activatedRoute: ActivatedRoute,
    private readonly router: Router,
    private readonly notificationService: NzNotificationService,
    private readonly modalService: NzModalService,
    private readonly location: Location,
    private readonly pageTitleService: PageTitleService,
    readonly formBuilder: FormBuilder,
    private readonly serverService: ServerService
  ) {
    this.createAccountForm = formBuilder.group(
      {
        accountName: ['', [Validators.required, this.accountNameValidator()]],
        displayName: ['', [Validators.required, this.displayNameValidator()]],
        password: ['', [this.passwordValidator()]],
        passwordConfirm: ['', []],
        role: [Role.Contestant, []],
      },
      {
        validators: [ConfirmedValidator('password', 'passwordConfirm')],
      }
    );

    this.editAccountForm = formBuilder.group(
      {
        displayName: ['', [Validators.required, this.displayNameValidator()]],
        password: ['', [this.passwordValidator()]],
        passwordConfirm: ['', []],
        role: [Role.Contestant, []],
      },
      {
        validators: [ConfirmedValidator('password', 'passwordConfirm')],
      }
    );

    this.serverInfoChangedSubscription =
      this.serverService.serverInfoChanged.subscribe((serverInfo) => {
        this.serverInfo = serverInfo;
      });
  }

  private accountNameValidator(): ValidatorFn {
    return (control: AbstractControl): { [k: string]: boolean } | null => {
      const accountName: string = control.value;
      return this.accountService.isValidAccountName(accountName);
    };
  }

  private displayNameValidator(): ValidatorFn {
    return (control: AbstractControl): { [k: string]: boolean } | null => {
      const displayName: string = control.value;
      return this.accountService.isValidDisplayName(displayName);
    };
  }

  private passwordValidator(): ValidatorFn {
    return (control: AbstractControl): { [k: string]: boolean } | null => {
      const password: string = control.value;
      return this.accountService.isValidPassword(password);
    };
  }

  ngOnInit(): void {
    (async () => {
      try {
        this.serverInfo = await this.serverService.getServerInfo();
      } catch (e) {
        this.notificationService.error('Failed to get server information', '');
        return;
      }

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
      this.pageTitleService.setTitle('Accounts');
    })().then();
    this.queryParamsSubscription = this.activatedRoute.queryParams.subscribe(
      async (params) => {
        await this.onQueryParamsChanged(params);
      }
    );
    this.serverInfoChangedSubscription =
      this.serverService.serverInfoChanged.subscribe((serverInfo) => {
        this.serverInfo = serverInfo;
      });
  }

  ngOnDestroy(): void {
    this.queryParamsSubscription?.unsubscribe();
    this.serverInfoChangedSubscription.unsubscribe();
  }

  private async onQueryParamsChanged(queryParams: Params): Promise<void> {
    this.getPaginationInfoFromQueryParams(queryParams);
    await this.loadAccountList();
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

  private async loadAccountList(): Promise<void> {
    try {
      this.loading = true;
      const { totalAccountCount, accountList } =
        await this.accountService.getAccountList(
          this.paginationService.getPageOffset(this.pageIndex, this.pageSize),
          this.pageSize
        );
      this.totalAccountCount = totalAccountCount;
      this.accountList = accountList;
    } catch (e) {
      if (e instanceof UnauthenticatedError) {
        this.notificationService.error(
          'Failed to load account list',
          'Not logged in'
        );
        this.router.navigateByUrl('/login');
        return;
      }

      if (e instanceof PermissionDeniedError) {
        this.notificationService.error(
          'Failed to load account list',
          'Permission denied'
        );
        this.location.back();
        return;
      }

      if (e instanceof InvalidAccountListParam) {
        this.notificationService.error(
          'Failed to load account list',
          'Invalid page index/size'
        );
        this.location.back();
        return;
      }

      this.notificationService.error(
        'Failed to load account list',
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
    this.router.navigate(['/account-list'], { queryParams });
  }

  public onAccountEditClicked(account: RpcAccount): void {
    this.editAccountForm.reset({
      displayName: account.displayName,
      password: '',
      passwordConfirm: '',
      role: account.role,
    });
    this.modalService.create({
      nzContent: this.editAccountModal,
      nzCloseIcon: '',
      nzOnOk: async () => {
        let { displayName, password, role } = this.editAccountForm.value;
        if (displayName === '' || displayName === account.displayName) {
          displayName = undefined;
        }
        if (password === '') {
          password = undefined;
        }
        if (role === '' || role === account.role) {
          role = undefined;
        }

        try {
          await this.accountService.updateAccount(
            account.iD,
            displayName,
            role,
            password
          );
          this.notificationService.success('Updated profile successfully', '');
          await this.loadAccountList();
        } catch (e) {
          if (e instanceof UnauthenticatedError) {
            this.notificationService.error(
              'Failed to update profile',
              'Not logged in'
            );
            this.router.navigateByUrl('/login');
            return;
          }
          if (e instanceof PermissionDeniedError) {
            this.notificationService.error(
              'Failed to update profile',
              'Permission denied'
            );
            return;
          }
          if (e instanceof InvalidAccountInfoError) {
            this.notificationService.error(
              'Failed to update profile',
              'Invalid account information'
            );
            return;
          }
          this.notificationService.error(
            'Failed to update profile',
            'Unknown error'
          );
        }
      },
    });
  }

  public onCreateAccountClicked(): void {
    this.createAccountForm.reset({
      accountName: '',
      displayName: '',
      password: '',
      passwordConfirm: '',
      role: Role.Contestant,
    });
    this.modalService.create({
      nzContent: this.createAccountModal,
      nzCloseIcon: '',
      nzOnOk: async () => {
        const { accountName, displayName, password, role } =
          this.createAccountForm.value;
        try {
          await this.accountService.createAccount(
            accountName,
            displayName,
            role,
            password
          );
          this.notificationService.success('Created profile successfully', '');
          await this.loadAccountList();
        } catch (e) {
          if (e instanceof UnauthenticatedError) {
            this.notificationService.error(
              'Failed to create profile',
              'Not logged in'
            );
            this.router.navigateByUrl('/login');
            return;
          }
          if (e instanceof PermissionDeniedError) {
            this.notificationService.error(
              'Failed to create profile',
              'Permission denied'
            );
            return;
          }
          if (e instanceof InvalidAccountInfoError) {
            this.notificationService.error(
              'Failed to create profile',
              'Invalid account information'
            );
            return;
          }
          this.notificationService.error(
            'Failed to create profile',
            'Unknown error'
          );
        }
      },
    });
  }
}
