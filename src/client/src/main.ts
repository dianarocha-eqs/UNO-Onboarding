import { bootstrapApplication } from '@angular/platform-browser';
import { AppComponent } from './app/app.component';
import { provideRouter, PreloadAllModules, withPreloading, withDebugTracing, NoPreloading} from '@angular/router';
import { routes } from './app/app.routes';
import { HttpClientModule, provideHttpClient } from '@angular/common/http';
import { importProvidersFrom } from '@angular/core';
import { provideIonicAngular } from '@ionic/angular/standalone';

bootstrapApplication(AppComponent, {
  providers: [
    // provideHttpClient(),
    provideRouter(
      routes,
      withPreloading( NoPreloading), // No Preloading Strategy for now
      withDebugTracing() // Debug tracing for development
    ), provideIonicAngular({}),
  ],
}).catch((err) => console.error(err));
