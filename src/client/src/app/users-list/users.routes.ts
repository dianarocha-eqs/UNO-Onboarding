import { Routes } from '@angular/router';

export const USERS_ROUTES: Routes = [
  {
    path: 'create',
    loadComponent: () =>
      import('../user-add/user-add.component').then((m) => m.UserAddComponent),
  },
  {
    path: 'list',
    loadComponent: () =>
      import('./users-list.component').then((m) => m.UserListingComponent),
  }
];
