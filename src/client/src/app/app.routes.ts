import { Routes } from '@angular/router';
export const routes: Routes = [
  {
    path: '',
    pathMatch: 'full',
    redirectTo: 'home', // Default route to authentication
  },
  {
    path: 'auth',
    loadChildren: () =>
      import('./auth/auth.routes').then((m) => m.AUTH_ROUTES),
  },
  {
    path: 'home',
    loadChildren: () =>
      import('./home/home.routes').then((m) => m.HOME_ROUTES),
  },
  {
    path: 'dashboard',
    loadChildren: () =>
      import('./dashboard/dashboard.routes').then((m) => m.DASHBOARD_ROUTES),
  },
  {
    path: 'sensors',
    loadChildren: () =>
      import('./sensors-list/sensors-list.routes').then((m) => m.SENSORS_ROUTES),
  },
  {
    path: 'users',
    loadChildren: () =>
      import('./user-add/user-add.routes').then((m) => m.USERS_ROUTES),
  },
  {
    path: '**',
    redirectTo: 'home',
  },
];
    
  
