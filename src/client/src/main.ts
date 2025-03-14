import { bootstrapApplication } from '@angular/platform-browser';
import { AppComponent } from './app/app.component';
import { provideRouter, PreloadAllModules, withPreloading, withDebugTracing, NoPreloading} from '@angular/router';
import { routes } from './app/app.routes';
import { provideHttpClient } from '@angular/common/http';

bootstrapApplication(AppComponent, {
  providers: [
    provideHttpClient(),
    provideRouter(
      routes,
      withPreloading( NoPreloading), // No Preloading Strategy for now
      withDebugTracing() // Debug tracing for development
    ),
  ],
}).catch((err) => console.error(err));
