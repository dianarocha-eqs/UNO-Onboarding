import { NgIf } from '@angular/common';
import { Component, ViewEncapsulation } from '@angular/core';
import { NavigationEnd, Router, RouterModule } from '@angular/router';
import { SidebarService } from '../sidebar.service';

@Component({
  selector: 'app-top-bar',
  imports: [RouterModule, NgIf],
  templateUrl: './top-bar.component.html',
  styleUrl: './top-bar.component.css'
})
export class TopBarComponent {

  currentRoute: string = '';
  pageTitle: string = '';
  mobileMenuOpen = false;

  constructor(private router: Router, private sidebarService: SidebarService) {
    // Listen for route changes
    this.router.events.subscribe(event => {
      if (event instanceof NavigationEnd) {
        this.currentRoute = event.url;
        this.updatePageTitle(event.url);
      }
    });
  }

  // Set page title based on route
  updatePageTitle(url: string) {
    const routesMap: { [key: string]: string } = {
      '/dashboard': 'Dashboard',
      '/sensors': 'Sensors',
      '/users': 'Users',
      '/settings': 'Settings',
      '/logout': 'Logout'
    };

    this.pageTitle = routesMap[url];
  }


  // Go back to the previous page
  goBack() {
    this.router.navigate(['..']); // Navigates back
  }


  toggleMobileMenu() {
    this.mobileMenuOpen = !this.mobileMenuOpen;
  }


  toggleSidebar() {
    this.sidebarService.toggleSidebar();
  }

  
}
