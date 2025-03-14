import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../environments/environment';
import { User } from './interfaces/user';

@Injectable({
  providedIn: 'root',
})
export class UserService {
  private apiUrl = `${environment.apiUrl}/v1/users`; // Use environment variable

  constructor(private http: HttpClient) {}

  /**
   * Sends a request to create a new user.
   * @param user The user object containing name, email, phone, picture, and role.
   * @param token Authentication token required for API access.
   * @returns An Observable containing the server response with a user UUID.
   */
  addUser(user: User, token: string): Observable<{ uuid: string }> {
    let headers = new HttpHeaders();
    headers = headers.set('Authorization', token);
    headers = headers.set('Role', user.role.toString());
    // Ensure you're passing the correct data
    return this.http.post<{ uuid: string }>(`${this.apiUrl}/create`, user, { headers });
  }
  
}
