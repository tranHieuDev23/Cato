import {
  Component,
  EventEmitter,
  forwardRef,
  Input,
  Output,
} from '@angular/core';
import {
  ControlValueAccessor,
  FormsModule,
  NG_VALUE_ACCESSOR,
} from '@angular/forms';
import { QuillModules, QuillModule } from 'ngx-quill';

@Component({
  selector: 'app-editable-rich-text',
  standalone: true,
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => EditableRichTextComponent),
      multi: true,
    },
  ],
  imports: [FormsModule, QuillModule],
  templateUrl: './editable-rich-text.component.html',
  styleUrls: ['./editable-rich-text.component.scss'],
})
export class EditableRichTextComponent implements ControlValueAccessor {
  @Input() public placeholder = '';
  @Input() public text = '';
  @Input() public readOnly = false;
  @Output() public textChange = new EventEmitter<string>();

  public quillModules: QuillModules = {
    toolbar: [
      [
        'bold',
        'italic',
        'underline',
        { list: 'ordered' },
        { list: 'bullet' },
        { script: 'sub' },
        { script: 'super' },
        { color: [] },
        { background: [] },
        { align: [] },
        'link',
      ],
    ],
  };

  private onChange: (value: unknown) => void = () => {};
  private onTouched: () => void = () => {};

  constructor() {}

  public onContentChanged(value: string): void {
    console.log(value);
    this.onChange(value);
    this.textChange.emit(value);
  }

  public writeValue(value: string): void {
    this.text = value;
    this.onChange(value);
    this.textChange.emit(value);
  }

  public registerOnChange(fn: (value: unknown) => void): void {
    this.onChange = fn;
  }

  public registerOnTouched(fn: () => void): void {
    this.onTouched = fn;
  }

  public setDisabledState(isDisabled: boolean): void {
    this.readOnly = isDisabled;
  }
}
