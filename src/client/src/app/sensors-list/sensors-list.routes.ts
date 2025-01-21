import { Routes } from '@angular/router';

export const SENSORS_ROUTES: Routes = [
  {
    path: '',
    loadComponent: () =>
      import('./sensors-list.component').then((m) => m.SensorsListComponent),
  },
  {
    path: ':id',
    loadComponent: () =>
      import('../sensor-details/sensor-details.component').then((m) => m.SensorDetailsComponent),
  },
  {
    path: ':id/edit',
    loadComponent: () =>
      import('../sensor-edit/sensor-edit.component').then((m) => m.SensorEditComponent),
  },
];
