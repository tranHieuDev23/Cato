import { Routes } from '@angular/router';
import { LoggedOutGuard } from './components/utils/logged-out-guard';
import { LoggedInGuard } from './components/utils/logged-in-guard';

export const routes: Routes = [
  { path: '', pathMatch: 'full', redirectTo: '/login' },
  {
    path: 'login',
    loadChildren: () =>
      import('./pages/login/login.routes').then((m) => m.LOGIN_ROUTES),
    canActivate: [LoggedOutGuard],
  },
  {
    path: 'profile',
    loadChildren: () =>
      import('./pages/profile/profile.routes').then((m) => m.PROFILE_ROUTES),
    canActivate: [LoggedInGuard],
  },
  {
    path: 'account-list',
    loadChildren: () =>
      import('./pages/account-list/account-list.routes').then(
        (m) => m.ACCOUNT_LIST_ROUTES
      ),
    canActivate: [LoggedInGuard],
  },
  {
    path: 'problem-list',
    loadChildren: () =>
      import('./pages/problem-list/problem-list.routes').then(
        (m) => m.PROBLEM_LIST_ROUTES
      ),
    canActivate: [LoggedInGuard],
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
    canActivate: [LoggedInGuard],
  },
  {
    path: 'problem-editor',
    loadChildren: () =>
      import('./pages/problem-editor/problem-editor.routes').then(
        (m) => m.PROBLEM_EDITOR_ROUTES
      ),
    canActivate: [LoggedInGuard],
  },
  {
    path: 'settings',
    loadChildren: () =>
      import('./pages/setting/setting.routes').then((m) => m.SETTING_ROUTES),
    canActivate: [LoggedInGuard],
  },
];
