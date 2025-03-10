import { Routes } from '@angular/router';

export const USERS_ROUTES: Routes = [
  {
    path: 'create',
    loadComponent: () =>
      import('./user-add.component').then((m) => m.UserAddComponent),
  }
];
