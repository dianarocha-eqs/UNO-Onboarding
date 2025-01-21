import { Routes } from '@angular/router';
export const routes: Routes = [
  {
    path: '',
    pathMatch: 'full',
    redirectTo: 'home', // Default route
  },
  {
    path: 'login',
    loadComponent: () =>
      import('./login-page/login-page.component').then((m) => m.LoginPageComponent),
  },
  {
    path: 'home',
    loadComponent: () =>
      import('./home/home.component').then((m) => m.HomeComponent),
  },
  {
    path: 'dashboard',
    loadComponent: () =>
      import('./dashboard/dashboard.component').then((m) => m.DashboardComponent),
  },
  {
    path: 'sensors',
    loadComponent: () => import('./sensors-list/sensors-list.component').then((m) => m.SensorsListComponent),
  },
  {
    path: 'sensors/:id', 
    loadComponent: () => 
      import('./sensor-details/sensor-details.component').then(m => m.SensorDetailsComponent),
  },
  {
    path: 'sensors/:id/edit', 
    loadComponent: () => 
      import('./sensor-edit/sensor-edit.component').then(m => m.SensorEditComponent),
  },
  {
    path: '**',
    redirectTo: 'login',
  },
];
    
  
