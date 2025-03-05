import { HttpClient, HttpErrorResponse, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { User } from './interfaces/user'
import { catchError, Observable, of } from 'rxjs';
import { environment } from '../environments/environment';


/**
 * Service for handling user-related API requests.
 */
@Injectable({
  providedIn: 'root'
})
export class UserService {

  /** API endpoint for users . */
  private baseurl = `${environment.apiUrl}/v1/users`;

  /**
   * Constructs the UserService.
   * @param http The Angular HttpClient for making HTTP requests.
   */
  constructor(private http: HttpClient) {}

  /**
   * Fetches a user by their ID.
   * @param id The unique identifier of the user.
   * @returns An `Observable` containing the requested `User` object.
   */
  getUser(id: string): Observable<User> {
    return this.http.get<User>(`${this.baseurl}/${id}`).pipe(
      catchError(this.handleError<User>('getUser', {} as User))
    );
  }

  /**
   * Handles HTTP errors gracefully.
   * @template T The expected return type of the failed operation.
   * @param operation The name of the operation that failed.
   * @param result Optional fallback result in case of failure.
   * @returns An `Observable` containing a fallback result.
   */

  private handleError<T>(operation = 'operation', result?: T) {
    return (error: HttpErrorResponse): Observable<T> => {
      alert(`Error: ${operation} failed: ${error.message}`); // Alert for API failure
      return of(result as T);
    };
  }

}
