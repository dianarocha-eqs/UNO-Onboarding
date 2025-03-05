import { Component } from '@angular/core';
import { ActivatedRoute, Router, RouterOutlet} from '@angular/router';
import { SidebarComponent } from './sidebar/sidebar.component';
import { NgIf } from '@angular/common';
import { TopBarComponent } from './top-bar/top-bar.component';

@Component({
  selector: 'app-root',
  standalone : true,
  imports: [RouterOutlet, SidebarComponent, TopBarComponent, NgIf],
  template: `
    <!-- app.component.html -->
    <div class="app-container">
      <!-- Sidebar (Navigation Bar) -->
      <app-sidebar *ngIf="!isLoginPage"></app-sidebar> <!-- Conditional rendering -->

      <!-- Topbar (Top Bar Information) -->
      <app-top-bar *ngIf="!isLoginPage"></app-top-bar> <!-- Conditional rendering -->

      <!-- Main Content Area -->
      <div class="main-content">
        <router-outlet></router-outlet> <!-- Routed Views Render Here -->
      </div>
    </div>

  `,
})
export class AppComponent {

  isLoginPage: boolean = false;

  constructor(private router: Router, private activatedRoute: ActivatedRoute) {
    // Check if the current route is 'auth' (login page)
    this.router.events.subscribe(() => {
      this.isLoginPage = this.router.url.includes('auth');
    });
  }
  
}
