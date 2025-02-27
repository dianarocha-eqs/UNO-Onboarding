import { Component, OnInit } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import { AuthService } from '../../auth.service';
import { CommonModule, NgIf } from '@angular/common';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
 
/**
 * Login Page Component
 * Handles user authentication via the login form.
 */
@Component({
  selector: 'app-login-page',
  standalone: true,
  imports: [RouterModule, NgIf, ReactiveFormsModule, CommonModule],
  templateUrl: './login-page.component.html',
  styleUrl: './login-page.component.css'
})
export class LoginPageComponent {

  /** Holds the login form */
  loginForm: FormGroup;

  /**
   * Constructs the LoginPageComponent.
   * @param authService Service for authentication.
   * @param router Angular Router for navigation.
   * @param fb 
   */
  constructor(private authService: AuthService, private router: Router, private fb: FormBuilder) {
    // Validation of the email and password
    this.loginForm = this.fb.group({
      email: ['', [Validators.required, Validators.email]],  
      password: ['', [Validators.required, Validators.minLength(12)]]
    });
  }


  /**
   * Attempts to log in a user using the provided credentials.
   * @param email The user's email address.
   * @param password The user's password.
   */
  login(): void {

    // Mark all fields as touched on submit to trigger validation
    // if empty fields, sets trigger for required fields
    this.loginForm.markAllAsTouched();


    // Prevent submission if form is invalid
    if (this.loginForm.invalid) {
      return; 
    }

    // Extract values from the form and pass them to the authentication service
    const { email, password } = this.loginForm.value;
    this.authService.login(email, password).subscribe({
      next: (response) => {
        console.log('Login successful', response);
        this.router.navigate(['/home']);
      },
      error: (error) => {
        console.error('Login failed', error);
      }
    });
  }

   /**
   * Retrieves the error message for the email input field.
   * The error message is shown if the email is required or has an invalid format 
   * after the field has been touched.
   * 
   * @returns string | null - The error message for the email field or null if valid.
   */
  get emailError(): string | null {
    const email = this.loginForm.get('email');
    if (email?.hasError('required') && email?.touched) {
      return 'Email is required';
    } else if (email?.hasError('email') && email?.touched) {
      return 'Invalid email format';
    }
    return null;
  }

   /**
   * Retrieves the error message for the password input field.
   * The error message is shown if the password is required or does not meet 
   * the minimum length after the field has been touched.
   * 
   * @returns string | null - The error message for the password field or null if valid.
   */
  get passwordError(): string | null {
    const password = this.loginForm.get('password');
    if (password?.hasError('required') && password?.touched) {
      return 'Password is required';
    } else if (password?.hasError('minlength') && password?.touched) {
      return 'Password must be at least 12 characters';
    }
    return null;
  }
}
