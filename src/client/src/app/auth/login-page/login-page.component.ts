import { Component, OnInit } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import { AuthService } from '../../auth.service';
import { User } from '../../interfaces/user'
 
/**
 * Login Page Component
 * Handles user authentication via the login form.
 */
@Component({
  selector: 'app-login-page',
  standalone: true,
  imports: [RouterModule],
  templateUrl: './login-page.component.html',
  styleUrl: './login-page.component.css'
})
export class LoginPageComponent {

  /** Holds user login credentials as a partial `User` object. */
  public user: Partial<User> = {
    email: '',
    password: ''
  };

  /** Stores error messages in case of login failure. */
  public errorMessage: string | null = null;

  /**
   * Constructs the LoginPageComponent.
   * @param authService Service for authentication.
   * @param router Angular Router for navigation.
   */
  constructor(private authService: AuthService, private router: Router) {}

  /**
   * Attempts to log in a user using the provided credentials.
   * @param email The user's email address.
   * @param password The user's password.
   */
  public login(email: string, password: string): void {
    if (email && password) {
      this.authService.login(email, password).subscribe({
        next: (response) => {
          console.log('Login successful', response);
           // Navigate to home page after successful login
           this.router.navigate(['/home']);
        },
        error: (error) => {
          console.error('Login failed', error);
          this.errorMessage = 'Invalid credentials. Please try again.';
        }
      });
    } else {
      console.log('Please fill in both fields');
    }
  }
}
