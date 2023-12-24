import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { NzFormModule } from 'ng-zorro-antd/form';
import { NzInputModule } from 'ng-zorro-antd/input';
import { NzLayoutModule } from 'ng-zorro-antd/layout';
import { NzSelectModule } from 'ng-zorro-antd/select';
import { RpcAccount } from '../../dataaccess/api';
import {
  AccountService,
  InvalidAccountInfoError,
  PermissionDeniedError,
  UnauthenticatedError,
} from '../../logic/account.service';
import { CommonModule } from '@angular/common';
import { NzButtonModule } from 'ng-zorro-antd/button';
import {
  AbstractControl,
  FormBuilder,
  FormGroup,
  ReactiveFormsModule,
  ValidatorFn,
  Validators,
} from '@angular/forms';
import { ConfirmedValidator } from '../../components/utils/confirmed-validator';
import {
  NzNotificationModule,
  NzNotificationService,
} from 'ng-zorro-antd/notification';
import { NzTypographyModule } from 'ng-zorro-antd/typography';
import { PageTitleService } from '../../logic/page-title.service';

@Component({
  selector: 'app-profile',
  standalone: true,
  imports: [
    NzFormModule,
    NzLayoutModule,
    NzInputModule,
    NzSelectModule,
    CommonModule,
    NzButtonModule,
    ReactiveFormsModule,
    NzNotificationModule,
    NzTypographyModule,
    NzInputModule,
  ],
  templateUrl: './profile.component.html',
  styleUrl: './profile.component.scss',
})
export class ProfileComponent implements OnInit {
  public sessionAccount: RpcAccount | undefined;
  public profileForm: FormGroup;

  constructor(
    private readonly accountService: AccountService,
    private readonly notificationService: NzNotificationService,
    private readonly router: Router,
    readonly formBuilder: FormBuilder,
    private readonly pageTitleService: PageTitleService
  ) {
    this.profileForm = formBuilder.group(
      {
        displayName: ['', [Validators.required, this.displayNameValidator()]],
        password: ['', [this.passwordValidator()]],
        passwordConfirm: ['', []],
      },
      {
        validators: [ConfirmedValidator('password', 'passwordConfirm')],
      }
    );
    this.profileForm.reset({
      displayName: '',
      password: '',
      passwordConfirm: '',
    });
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
      if (!password) {
        return null;
      }

      return this.accountService.isValidPassword(password);
    };
  }

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
      this.profileForm.reset({
        displayName: this.sessionAccount?.displayName,
        password: '',
        passwordConfirm: '',
      });
      this.pageTitleService.setTitle('Edit Profile');
    })().then();
  }

  public async onUpdateClicked(): Promise<void> {
    if (!this.sessionAccount) {
      return;
    }

    let { displayName, password } = this.profileForm.value;
    if (displayName === '' || displayName === this.sessionAccount.displayName) {
      displayName = undefined;
    }
    if (password === '') {
      password = undefined;
    }

    try {
      const account = await this.accountService.updateAccount(
        this.sessionAccount.iD,
        displayName,
        undefined,
        password
      );
      this.notificationService.success('Updated profile successfully', '');
      this.sessionAccount = account;
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
  }
}
