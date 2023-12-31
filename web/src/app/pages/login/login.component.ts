import { Component, OnInit } from '@angular/core';
import {
  AbstractControl,
  ReactiveFormsModule,
  UntypedFormBuilder,
  UntypedFormGroup,
  ValidatorFn,
  Validators,
} from '@angular/forms';
import { NzTabsModule } from 'ng-zorro-antd/tabs';
import { NzFormModule } from 'ng-zorro-antd/form';
import { NzInputModule } from 'ng-zorro-antd/input';
import { NzButtonModule } from 'ng-zorro-antd/button';
import { NzRadioModule } from 'ng-zorro-antd/radio';
import {
  NzNotificationModule,
  NzNotificationService,
} from 'ng-zorro-antd/notification';
import { Router } from '@angular/router';
import {
  AccountLoginDisabledError,
  AccountNameTakenError,
  AccountNotFoundError,
  AccountService,
  IncorrectPasswordError,
  InvalidAccountInfoError,
  Role,
} from '../../logic/account.service';
import { ConfirmedValidator } from '../../components/utils/confirmed-validator';
import { CommonModule } from '@angular/common';
import { PageTitleService } from '../../logic/page-title.service';
import { NzGridModule } from 'ng-zorro-antd/grid';
import { RpcGetServerInfoResponse } from '../../dataaccess/api';
import { ServerService } from '../../logic/server.service';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [
    CommonModule,
    NzTabsModule,
    NzFormModule,
    ReactiveFormsModule,
    NzInputModule,
    NzNotificationModule,
    NzButtonModule,
    NzRadioModule,
    NzGridModule,
  ],
  templateUrl: './login.component.html',
  styleUrl: './login.component.scss',
})
export class LoginComponent implements OnInit {
  public loginForm: UntypedFormGroup;
  public registerForm: UntypedFormGroup;

  public serverInfo: RpcGetServerInfoResponse | undefined;

  constructor(
    private readonly accountService: AccountService,
    private readonly notificationService: NzNotificationService,
    private readonly router: Router,
    formBuilder: UntypedFormBuilder,
    private readonly pageTitleService: PageTitleService,
    private readonly serverService: ServerService
  ) {
    this.loginForm = formBuilder.group({
      accountName: ['', [Validators.required]],
      password: ['', [Validators.required]],
    });
    this.loginForm.reset({ accountName: '', password: '' });
    this.registerForm = formBuilder.group(
      {
        displayName: ['', [Validators.required, this.displayNameValidator()]],
        accountName: ['', [Validators.required, this.accountNameValidator()]],
        password: ['', [Validators.required, this.passwordValidator()]],
        passwordConfirm: ['', [Validators.required]],
        role: [''],
      },
      {
        validators: [ConfirmedValidator('password', 'passwordConfirm')],
      }
    );
    this.registerForm.reset({
      displayName: '',
      accountName: '',
      password: '',
      passwordConfirm: '',
      role: Role.Contestant,
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
    this.pageTitleService.setTitle('Cato');
    (async () => {
      try {
        this.serverInfo = await this.serverService.getServerInfo();
      } catch (e) {
        this.notificationService.error('Failed to get server information', '');
        return;
      }
    })().then();
  }

  public async onLoginClicked(): Promise<void> {
    const { accountName, password } = this.loginForm.value;
    try {
      const account = await this.accountService.createSession(
        accountName,
        password
      );
      this.notificationService.success(
        'Logged in successfully',
        `Welcome, ${account.displayName}`
      );
    } catch (e) {
      if (e instanceof AccountNotFoundError) {
        this.notificationService.error('Failed to log in', 'Account not found');
        return;
      }
      if (e instanceof IncorrectPasswordError) {
        this.notificationService.error(
          'Failed to log in',
          'Incorrect password'
        );
        return;
      }
      if (e instanceof IncorrectPasswordError) {
        this.notificationService.error(
          'Failed to log in',
          'Incorrect password'
        );
        return;
      }
      if (e instanceof AccountLoginDisabledError) {
        this.notificationService.error('Failed to log in', 'Login is disabled');
        return;
      }
      this.notificationService.error('Failed to log in', 'Unknown error');
      return;
    }
    this.router.navigateByUrl('/problem-list');
  }

  public async onRegisterClicked(): Promise<void> {
    const { displayName, accountName, password, role } =
      this.registerForm.value;
    try {
      const account = await this.accountService.createAccount(
        accountName,
        displayName,
        role,
        password
      );
      this.notificationService.success(
        'Registered successfully',
        `Welcome, ${account.displayName}`
      );
    } catch (e) {
      if (e instanceof AccountNameTakenError) {
        this.notificationService.error(
          'Failed to register',
          'Account name is already taken'
        );
        return;
      }
      if (e instanceof InvalidAccountInfoError) {
        this.notificationService.error(
          'Failed to register',
          'Invalid account information'
        );
        return;
      }
      this.notificationService.error('Failed to log in', 'Unknown error');
      return;
    }

    try {
      await this.accountService.createSession(accountName, password);
    } catch (e) {
      return;
    }

    this.router.navigateByUrl('/problem-list');
  }
}
