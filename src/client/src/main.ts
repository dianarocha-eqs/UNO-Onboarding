import { bootstrapApplication } from '@angular/platform-browser';
import { AppComponent } from './app/app.component';
import { provideRouter, PreloadAllModules, withPreloading, withDebugTracing} from '@angular/router';
import { routes } from './app/app.routes';
import { HttpClientModule, provideHttpClient } from '@angular/common/http';
import { importProvidersFrom } from '@angular/core';

bootstrapApplication(AppComponent, {
  providers: [
    provideHttpClient(),
    provideRouter(
      routes,
      withPreloading(PreloadAllModules), // Preloading strategy
      withDebugTracing() // Debug tracing for development
    ),
  ],
}).catch((err) => console.error(err));
