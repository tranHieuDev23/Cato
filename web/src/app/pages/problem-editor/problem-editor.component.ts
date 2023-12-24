import { Component, OnInit } from '@angular/core';
import {
  AbstractControl,
  FormBuilder,
  FormGroup,
  ReactiveFormsModule,
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
  ProblemNotFoundError,
  ProblemService,
} from '../../logic/problem.service';
import {
  NzNotificationModule,
  NzNotificationService,
} from 'ng-zorro-antd/notification';
import { CommonModule, Location } from '@angular/common';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { EditableRichTextComponent } from '../../components/editable-rich-text/editable-rich-text.component';
import { NzInputModule } from 'ng-zorro-antd/input';
import { NzSelectModule } from 'ng-zorro-antd/select';
import { NzGridModule } from 'ng-zorro-antd/grid';
import { NzInputNumberModule } from 'ng-zorro-antd/input-number';
import { NzButtonModule } from 'ng-zorro-antd/button';
import { PageTitleService } from '../../logic/page-title.service';
import { UnitService } from '../../logic/unit.service';

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
export class ProblemEditorComponent implements OnInit {
  public sessionAccount: RpcAccount | undefined;
  public exampleList: RpcProblemExample[] = [];

  private problemID: number | undefined;
  public formGroup: FormGroup;
  public saving = false;

  constructor(
    private readonly accountService: AccountService,
    private readonly problemService: ProblemService,
    private readonly notificationService: NzNotificationService,
    private readonly location: Location,
    private readonly router: Router,
    readonly formBuilder: FormBuilder,
    private readonly pageTitleService: PageTitleService,
    private readonly activatedRoute: ActivatedRoute,
    private readonly unitService: UnitService
  ) {
    this.formGroup = formBuilder.group(
      {
        displayName: ['', [Validators.required, this.displayNameValidator()]],
        description: ['', [Validators.required, Validators.maxLength(5000)]],
        timeLimit: [1, [Validators.required]],
        timeUnit: ['second', [Validators.required]],
        memoryLimit: [1, [Validators.required]],
        memoryUnit: ['GB', [Validators.required]],
      },
      {
        validators: [
          this.valueAndUnitValidator('timeLimit', 'timeUnit', 1, 10000),
          this.valueAndUnitValidator(
            'memoryLimit',
            'memoryUnit',
            1,
            8 * 1024 * 1024
          ),
        ],
      }
    );
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
      this.pageTitleService.setTitle('Problem Editor');
    })().then();
    this.activatedRoute.params.subscribe(async (params) => {
      await this.onParamsChanged(params);
    });
  }

  private async onParamsChanged(params: Params): Promise<void> {
    if (!params['id']) {
      this.formGroup.reset({
        displayName: '',
        description: '',
        timeLimit: 1,
        timeUnit: 'second',
        memoryLimit: 1,
        memoryUnit: 'GB',
      });
      return;
    }

    this.problemID = +params['id'];
    try {
      const problem = await this.problemService.getProblem(this.problemID);
      const { value: timeLimit, unit: timeUnit } =
        this.unitService.timeLimitToValueAndUnit(
          problem.timeLimitInMillisecond
        );
      const { value: memoryLimit, unit: memoryUnit } =
        this.unitService.memoryLimitToValueAndUnit(problem.memoryLimitInByte);
      this.formGroup.reset({
        displayName: problem.displayName,
        description: problem.description,
        timeLimit: timeLimit,
        timeUnit: timeUnit,
        memoryLimit: memoryLimit,
        memoryUnit: memoryUnit,
      });
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

      if (e instanceof ProblemNotFoundError) {
        this.notificationService.error(
          'Failed to load problem',
          'Problem cannot be found'
        );
        this.location.back();
        return;
      }

      this.notificationService.error('Failed to load problem', 'Unknown error');
      this.location.back();
    }
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
      timeUnit,
      memoryLimit,
      memoryUnit,
    } = this.formGroup.value;
    try {
      this.saving = true;
      if (this.problemID === undefined) {
        const problem = await this.problemService.createProblem(
          displayName,
          description,
          this.unitService.timeValueAndUnitToLimit(timeLimit, timeUnit),
          this.unitService.memoryValueAndUnitToLimit(memoryLimit, memoryUnit)
        );

        this.notificationService.success('Problem saved successfully!', '');
        this.router.navigateByUrl(`/problem/${problem.iD}`);
      } else {
        await this.problemService.updateProblem(
          this.problemID,
          displayName,
          description,
          this.unitService.timeValueAndUnitToLimit(timeLimit, timeUnit),
          this.unitService.memoryValueAndUnitToLimit(memoryLimit, memoryUnit)
        );
        this.notificationService.success('Problem updated successfully!', '');
        this.router.navigateByUrl(`/problem/${this.problemID}`);
      }
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

      if (e instanceof ProblemNotFoundError) {
        this.notificationService.error(
          'Failed to save problem',
          'Problem not found'
        );
        this.location.back();
        return;
      }

      this.notificationService.error('Failed to save problem', 'Unknown error');
    } finally {
      this.saving = false;
    }
  }
}
