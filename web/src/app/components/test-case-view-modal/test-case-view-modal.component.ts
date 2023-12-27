import { Component, Input, OnInit, inject } from '@angular/core';
import { NZ_MODAL_DATA } from 'ng-zorro-antd/modal';
import { FormsModule } from '@angular/forms';
import { NzTypographyModule } from 'ng-zorro-antd/typography';
import { NgeMonacoModule } from '@cisstech/nge/monaco';

export interface TestCaseViewModalData {
  input?: string;
  output?: string;
}

@Component({
  selector: 'app-test-case-view-modal',
  standalone: true,
  imports: [NgeMonacoModule, FormsModule, NzTypographyModule],
  templateUrl: './test-case-view-modal.component.html',
  styleUrl: './test-case-view-modal.component.scss',
})
export class TestCaseViewModalComponent implements OnInit {
  @Input() public input = '';
  @Input() public output = '';

  private readonly nzModalData: TestCaseViewModalData | null = inject(
    NZ_MODAL_DATA,
    { optional: true }
  );

  ngOnInit(): void {
    if (!this.nzModalData) {
      return;
    }
    this.input = this.nzModalData.input || '';
    this.output = this.nzModalData.output || '';
  }

  public onEditTestCaseInputEditorReady(editor: monaco.editor.IEditor): void {
    editor.updateOptions({ minimap: { enabled: false } });
    const editorModel = monaco.editor.createModel(this.input);
    editor.setModel(editorModel);
  }

  public onEditTestCaseOutputEditorReady(editor: monaco.editor.IEditor): void {
    editor.updateOptions({ minimap: { enabled: false } });
    const editorModel = monaco.editor.createModel(this.output);
    editor.setModel(editorModel);
  }
}
