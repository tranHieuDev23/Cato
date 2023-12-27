import {
  Component,
  EventEmitter,
  Input,
  OnDestroy,
  OnInit,
  Output,
  inject,
} from '@angular/core';
import { FormsModule } from '@angular/forms';
import { NgeMonacoModule } from '@cisstech/nge/monaco';
import { NzCheckboxModule } from 'ng-zorro-antd/checkbox';
import { NZ_MODAL_DATA } from 'ng-zorro-antd/modal';
import { NzTypographyModule } from 'ng-zorro-antd/typography';

export interface TestCaseEditorModalData {
  input?: string;
  output?: string;
  isHidden?: boolean | null;
}

@Component({
  selector: 'app-test-case-editor-modal',
  standalone: true,
  imports: [NgeMonacoModule, FormsModule, NzTypographyModule, NzCheckboxModule],
  templateUrl: './test-case-editor-modal.component.html',
  styleUrl: './test-case-editor-modal.component.scss',
})
export class TestCaseEditorModalComponent implements OnInit, OnDestroy {
  @Input() public input = '';
  @Input() public output = '';
  @Input() public isHidden: boolean | null = true;

  @Output() public inputChange = new EventEmitter<string>();
  @Output() public outputChange = new EventEmitter<string>();
  @Output() public isHiddenChange = new EventEmitter<boolean>();

  private readonly nzModalData: TestCaseEditorModalData | null = inject(
    NZ_MODAL_DATA,
    { optional: true }
  );

  private editInputOnChange: monaco.IDisposable | undefined;
  private editOutputOnChange: monaco.IDisposable | undefined;

  ngOnInit(): void {
    if (!this.nzModalData) {
      return;
    }
    this.input = this.nzModalData.input || '';
    this.output = this.nzModalData.output || '';
    this.isHidden = this.nzModalData.isHidden || true;
  }

  ngOnDestroy(): void {
    this.editInputOnChange?.dispose();
    this.editOutputOnChange?.dispose();
  }

  public onEditTestCaseInputEditorReady(editor: monaco.editor.IEditor): void {
    if (this.editInputOnChange) {
      this.editInputOnChange.dispose();
    }

    editor.updateOptions({ minimap: { enabled: false } });
    const editorModel = monaco.editor.createModel(this.input);
    this.editInputOnChange = editorModel.onDidChangeContent(() => {
      this.input = editorModel.getValue();
      this.inputChange.emit(this.input);
    });

    editor.setModel(editorModel);
  }

  public onEditTestCaseOutputEditorReady(editor: monaco.editor.IEditor): void {
    if (this.editOutputOnChange) {
      this.editOutputOnChange.dispose();
    }

    editor.updateOptions({ minimap: { enabled: false } });
    const editorModel = monaco.editor.createModel(this.output);
    this.editOutputOnChange = editorModel.onDidChangeContent(() => {
      this.output = editorModel.getValue();
      this.outputChange.emit(this.output);
    });

    editor.setModel(editorModel);
  }
}
