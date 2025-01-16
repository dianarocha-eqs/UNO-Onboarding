import { Routes } from '@angular/router';
import { HomeComponent } from './home/home.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { LoginPageComponent } from './login-page/login-page.component';
import { SensorDetailsComponent } from './sensor-details/sensor-details.component';
import { SensorsListComponent } from './sensors-list/sensors-list.component';
import { SensorEditComponent } from './sensor-edit/sensor-edit.component';

export const routes: Routes = [
    { path: '', redirectTo: '/login', pathMatch: 'full' },
    { path: 'login', component: LoginPageComponent },
    
          // in the future we need to add authentication protection so home and the other routers are only accessed if login is successful
    { 
      path: 'home', 
      component: HomeComponent,
    },
    { path: 'dashboard', 
      component: DashboardComponent, 
      },
    { 
      path: 'sensors', 
      component: SensorsListComponent, 
    },
    { path: 'sensors/:id',
      component: SensorDetailsComponent, 
      },
    { 
      path: 'sensors/:id/edit', 
      component: SensorEditComponent, 
    }
];
    
  
