import { Routes } from '@angular/router';

export const routes: Routes = [
  { path: '', pathMatch: 'full', redirectTo: '/welcome' },
  {
    path: 'login',
    loadChildren: () =>
      import('./pages/login/login.routes').then((m) => m.LOGIN_ROUTES),
  },
  {
    path: 'welcome',
    loadChildren: () =>
      import('./pages/welcome/welcome.routes').then((m) => m.WELCOME_ROUTES),
  },
  {
    path: 'profile',
    loadChildren: () =>
      import('./pages/profile/profile.routes').then((m) => m.PROFILE_ROUTES),
  },
  {
    path: 'account-list',
    loadChildren: () =>
      import('./pages/account-list/account-list.routes').then(
        (m) => m.ACCOUNT_LIST_ROUTES
      ),
  },
  {
    path: 'problem-list',
    loadChildren: () =>
      import('./pages/problem-list/problem-list.routes').then(
        (m) => m.PROBLEM_LIST_ROUTES
      ),
  },
  {
    path: 'submission-list',
    loadChildren: () =>
      import('./pages/submission-list/submission-list.routes').then(
        (m) => m.SUBMISSION_LIST_ROUTES
      ),
  },
  {
    path: 'problem',
    loadChildren: () =>
      import('./pages/problem/problem.routes').then((m) => m.PROBLEM_ROUTES),
  },
  {
    path: 'problem-editor',
    loadChildren: () =>
      import('./pages/problem-editor/problem-editor.routes').then(
        (m) => m.PROBLEM_EDITOR_ROUTES
      ),
  },
  {
    path: 'submission',
    loadChildren: () =>
      import('./pages/submission/submission.routes').then(
        (m) => m.SUBMISSION_ROUTES
      ),
  },
];
