import { HttpClient, HttpErrorResponse, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { User } from './interfaces/user'
import { catchError, Observable, of } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class UserService {

  private baseurl = 'http://localhost:8080/v1/users'; // URL to backend API users
  
  constructor(private http: HttpClient) {}

    // GET
  getUser(id: string): Observable<User> {
    return this.http.get<User>(`${this.baseurl}/${id}`).pipe(
      catchError(this.handleError<User>('getUser', {} as User))
    );
  }

  private handleError<T>(operation = 'operation', result?: T) {
    return (error: HttpErrorResponse): Observable<T> => {
      console.error(`${operation} failed: ${error.message}`);
      // Show an alert or log to see what's happening
      alert(`Error: ${operation} failed: ${error.message}`);  // This will pop up an alert if the API fails
      return of(result as T); 
    };
  }
  private getHttpOptions() {
    const headers = new HttpHeaders({ 'Content-Type': 'application/json' });
    return { headers };
  }
}
