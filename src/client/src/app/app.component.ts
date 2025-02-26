import { Component } from '@angular/core';
import { RouterOutlet, RouterModule} from '@angular/router';
import { FormsModule } from '@angular/forms';  // <-- Import FormsModule here
import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

@Component({
  selector: 'app-root',
  standalone : true,
  imports: [RouterOutlet, RouterModule],
  template: `
    <div class="app-container">
      <router-outlet></router-outlet>
    </div>
  `,
})
export class AppComponent {
}
