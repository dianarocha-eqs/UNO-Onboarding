import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable, tap, throwError } from 'rxjs';
import { User } from './interfaces/user';
import { environment } from '../environments/environment';

/**
 * Represents the response received after a successful login.
 */
export interface LoginResponse {
  /**
   * Authentication token used to maintain the user's session.
   */
  token: string;

  /**
   * Authenticated user's details.
   */
  user: User;
}

/**
 * Service for handling authentication-related API requests.
 */
@Injectable({
  providedIn: 'root'
})
export class AuthService {

  /** API endpoint for authentication */
  private baseurl = `${environment.apiUrl}/v1/auth`;

  /** HTTP options for requests */
  private httpOptions = {
    headers: new HttpHeaders({ 'Content-Type': 'application/json' })
  };

  /**
   * Constructs the AuthService.
   * @param http Angular HttpClient for making HTTP requests.
   */
  constructor(private http: HttpClient) {}

  /**
   * Authenticates the user by sending login credentials to the API.
   * @param email The user's email address.
   * @param password The user's password.
   * @returns An `Observable` containing the login response.
   */
  login(email: string, password: string): Observable<LoginResponse> {

    if (!this.validateEmail(email)) {
      return throwError(() => new Error('Invalid email format.'));
    }

    if (!this.validatePassword(password)) {
      return throwError(() => new Error('Password must be 12 characters.'));
    }

    
    const body = JSON.stringify({ email, password });
    return this.http.post<LoginResponse>(`${this.baseurl}/login`, body,  this.httpOptions).pipe(
      tap((response) => {
        localStorage.setItem('token', response.token);
        localStorage.setItem('user', JSON.stringify(response.user));
      })
    );
  }

  /**
   * Validates email format.
   * @param email The email string to validate.
   * @returns `true` if valid, otherwise `false`.
   */
  private validateEmail(email: string): boolean {
    const emailPattern = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
    return emailPattern.test(email);
  }

  /**
   * Validates password length (all passwords were set to be size 12).
   * @param password The password string to validate.
   * @returns `true` if valid, otherwise `false`.
   */
  private validatePassword(password: string): boolean {
    return password.length == 12;
  }

}
