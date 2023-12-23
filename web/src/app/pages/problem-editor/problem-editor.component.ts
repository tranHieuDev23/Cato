import { Component, OnDestroy, OnInit } from '@angular/core';
import {
  AbstractControl,
  FormBuilder,
  FormGroup,
  ReactiveFormsModule,
  UntypedFormGroup,
  ValidatorFn,
  Validators,
} from '@angular/forms';
import { NzFormModule } from 'ng-zorro-antd/form';
import { NzTypographyModule } from 'ng-zorro-antd/typography';
import { RpcAccount, RpcProblemExample } from '../../dataaccess/api';
import {
  AccountService,
  PermissionDeniedError,
  UnauthenticatedError,
} from '../../logic/account.service';
import {
  InvalidProblemInfo,
  ProblemService,
} from '../../logic/problem.service';
import {
  NzNotificationModule,
  NzNotificationService,
} from 'ng-zorro-antd/notification';
import { CommonModule, Location } from '@angular/common';
import { Router } from '@angular/router';
import { Subscription } from 'rxjs';
import { EditableRichTextComponent } from '../../components/editable-rich-text/editable-rich-text.component';
import { NzInputModule } from 'ng-zorro-antd/input';
import { NzSelectModule } from 'ng-zorro-antd/select';
import { NzGridModule } from 'ng-zorro-antd/grid';
import { NzInputNumberModule } from 'ng-zorro-antd/input-number';
import { NzButtonModule } from 'ng-zorro-antd/button';

const KB_IN_BYTE = 1024;
const MB_IN_BYTE = KB_IN_BYTE * 1024;
const GB_IN_BYTE = MB_IN_BYTE * 1024;

@Component({
  selector: 'app-problem-editor',
  standalone: true,
  imports: [
    CommonModule,
    NzFormModule,
    ReactiveFormsModule,
    NzTypographyModule,
    NzNotificationModule,
    EditableRichTextComponent,
    NzInputModule,
    NzGridModule,
    NzSelectModule,
    NzInputNumberModule,
    NzButtonModule,
  ],
  templateUrl: './problem-editor.component.html',
  styleUrl: './problem-editor.component.scss',
})
export class ProblemEditorComponent implements OnInit, OnDestroy {
  public sessionAccount: RpcAccount | null | undefined;
  public exampleList: RpcProblemExample[] = [];

  public formGroup: FormGroup;
  public saving = false;

  private sessionAccountChangedSubscription: Subscription | undefined;

  constructor(
    private readonly accountService: AccountService,
    private readonly problemService: ProblemService,
    private readonly notificationService: NzNotificationService,
    private readonly location: Location,
    private readonly router: Router,
    readonly formBuilder: FormBuilder
  ) {
    this.formGroup = formBuilder.group(
      {
        displayName: ['', [Validators.required, this.displayNameValidator()]],
        description: ['', [Validators.required, Validators.maxLength(5000)]],
        timeLimit: [1, [Validators.required]],
        timeLimitUnitInMillisecond: [1000, [Validators.required]],
        memoryLimit: [1, [Validators.required]],
        memoryLimitUnitInByte: [GB_IN_BYTE, [Validators.required]],
      },
      {
        validators: [
          this.valueAndUnitValidator(
            'timeLimit',
            'timeLimitUnitInMillisecond',
            1,
            10000
          ),
          this.valueAndUnitValidator(
            'memoryLimit',
            'memoryLimitUnitInByte',
            1,
            8 * GB_IN_BYTE
          ),
        ],
      }
    );
    this.formGroup.reset({
      displayName: '',
      description: '',
      timeLimit: 1,
      timeLimitUnitInMillisecond: 1000,
      memoryLimit: 1,
      memoryLimitUnitInByte: GB_IN_BYTE,
    });
  }

  ngOnInit(): void {
    (async () => {
      this.sessionAccount = await this.accountService.getSessionAccount();
    })().then();
    this.sessionAccountChangedSubscription =
      this.accountService.sessionAccountChanged.subscribe((account) => {
        this.sessionAccount = account;
      });
  }

  ngOnDestroy(): void {
    this.sessionAccountChangedSubscription?.unsubscribe();
  }

  private displayNameValidator(): ValidatorFn {
    return (control: AbstractControl): { [k: string]: boolean } | null => {
      const displayName: string = control.value;
      return this.problemService.isValidDisplayName(displayName);
    };
  }

  private valueAndUnitValidator(
    valueName: string,
    unitName: string,
    minValue: number,
    maxValue: number
  ): ValidatorFn {
    return (control: AbstractControl): { [k: string]: boolean } | null => {
      const valueControl = control.get(valueName);
      const unitControl = control.get(unitName);
      if (valueControl?.errors || unitControl?.errors) {
        return null;
      }

      const value = +valueControl?.getRawValue();
      const unit = +unitControl?.getRawValue();
      const valueWithUnit = value * unit;
      if (valueWithUnit < minValue) {
        return { minimum: true };
      }

      if (valueWithUnit > maxValue) {
        return { maximum: true };
      }

      return null;
    };
  }

  public async onSaveClicked(): Promise<void> {
    const {
      displayName,
      description,
      timeLimit,
      timeLimitUnitInMillisecond,
      memoryLimit,
      memoryLimitUnitInByte,
    } = this.formGroup.value;
    try {
      this.saving = true;
      const problem = await this.problemService.createProblem(
        displayName,
        description,
        timeLimit * timeLimitUnitInMillisecond,
        memoryLimit * memoryLimitUnitInByte
      );
      this.notificationService.success('Problem saved successfully!', '');
      this.router.navigateByUrl('/problem-list');
    } catch (e) {
      if (e instanceof UnauthenticatedError) {
        this.notificationService.error(
          'Failed to save problem',
          'Not logged in'
        );
        this.router.navigateByUrl('/login');
        return;
      }

      if (e instanceof PermissionDeniedError) {
        this.notificationService.error(
          'Failed to save problem',
          'Permission denied'
        );
        this.location.back();
        return;
      }

      if (e instanceof InvalidProblemInfo) {
        this.notificationService.error(
          'Failed to save problem',
          'Invalid page index/size'
        );
        this.location.back();
        return;
      }

      this.notificationService.error('Failed to save problem', 'Unknown error');
      this.location.back();
    } finally {
      this.saving = false;
    }
  }
}
