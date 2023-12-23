import { UntypedFormGroup } from '@angular/forms';

export function ConfirmedValidator(controlName: string, matchingControlName: string) {
  return (formGroup: UntypedFormGroup) => {
    const control = formGroup.controls[controlName];
    const matchingControl = formGroup.controls[matchingControlName];
    if (control.errors || matchingControl.errors) {
      return;
    }
    if (control.value !== matchingControl.value) {
      matchingControl.setErrors({ confirmed: true });
    }
  };
}
