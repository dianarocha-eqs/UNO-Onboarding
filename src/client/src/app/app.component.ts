import { Component } from '@angular/core';
import { RouterOutlet, RouterModule} from '@angular/router';

@Component({
  selector: 'app-root',
  standalone : true,
  imports: [RouterOutlet],
  template: `
    <div class="app-container">
      <router-outlet></router-outlet>
    </div>
  `,
})
export class AppComponent {
}
