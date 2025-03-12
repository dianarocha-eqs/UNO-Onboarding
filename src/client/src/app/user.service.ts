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
    return this.http.post<{ uuid: string }>(`${this.apiUrl}/create`, user, { headers });
  }

  /**
   * Fetches a list of users based on search criteria.
   * Only accessible to administrators.
   * @param token Authentication token required for API access.
   * @param search Search query (name or email).
   * @param sort Sorting order ('asc' or 'desc').
   * @param user The user object to obtain the role.
   * @returns An Observable containing the user list response.
   */
  listUsers(token: string, body: any, user: User): Observable<User[]> {
    let headers = new HttpHeaders();
    headers = headers.set('Authorization', token);
    headers = headers.set('Role', user.role.toString());
    return this.http.post<User[]>(`${this.apiUrl}/list`, body, { headers });
  }  
  
}
