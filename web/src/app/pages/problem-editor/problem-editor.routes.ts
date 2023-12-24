import { Routes } from '@angular/router';
import { ProblemEditorComponent } from './problem-editor.component';

export const PROBLEM_EDITOR_ROUTES: Routes = [
  { path: '', component: ProblemEditorComponent },
  { path: ':id', component: ProblemEditorComponent },
];
