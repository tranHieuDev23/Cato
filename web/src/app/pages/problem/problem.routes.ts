import { Routes } from '@angular/router';
import { ProblemComponent } from './problem.component';

export const PROBLEM_ROUTES: Routes = [
  { path: ':id', component: ProblemComponent },
];
