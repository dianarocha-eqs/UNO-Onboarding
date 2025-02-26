import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable, tap } from 'rxjs';
import { User } from './interfaces/user';

export interface LoginResponse {
  token: string;
  user: User;
}

@Injectable({
  providedIn: 'root'
})
export class AuthService {

  private apiUrl = 'http://localhost:8080/v1/auth'; // URL for login, signup, etc.

  httpOptions = {
    headers: new HttpHeaders({ 'Content-Type': 'application/json' })
  };

  constructor(
    private http: HttpClient
  ) {}

  // Authenticate the user
  login(email: string, password: string): Observable<LoginResponse> {
    const body = JSON.stringify({ email, password });
    return this.http.post<LoginResponse>(`${this.apiUrl}/login`, body,  this.httpOptions).pipe(
      tap((response) => {
        localStorage.setItem('token', response.token);
        localStorage.setItem('user', JSON.stringify(response.user));
      })
    );
  }
}
