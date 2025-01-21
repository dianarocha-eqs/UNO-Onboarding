import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse, HttpHeaders } from '@angular/common/http';
import { Observable, of, throwError } from 'rxjs';
import { catchError, tap } from 'rxjs/operators';
import { sensor } from './interfaces/sensor'; // Assuming sensor interface is defined elsewhere

@Injectable({
  providedIn: 'root'
})
export class SensorService {
  private baseurl = 'http://localhost:8080/api/sensors'; // URL to your backend API for sensors

  constructor(private http: HttpClient) {}

   // GET
   getSensors(page: number = 1, size: number = 10, sort: string = 'name'): Observable<sensor[]> {
    return this.http.get<sensor[]>(`${this.baseurl}?page=${page}&size=${size}&sort=${sort}`).pipe(
      catchError(this.handleError<sensor[]>('getSensors', []))
    );
  }
  

   // GET
  getSensorById(id: number): Observable<sensor> {
    return this.http.get<sensor>(`${this.baseurl}/${id}`).pipe(
      catchError(this.handleError<sensor>('getSensorById', {} as sensor))
    );
  }

   // GET
   searchSensors(term: string): Observable<sensor[]> {
    if (!term.trim()) {
      return of([]);
    }
    return this.http.get<sensor[]>(`${this.baseurl}/search?term=${term}`).pipe(
      catchError(this.handleError<sensor[]>('searchSensors', []))
    );
  }

  // ---------------------------- //
  // POST
  addSensor(sensor: Omit<sensor, 'id'>): Observable<sensor> {
    return this.http.post<sensor>(`${this.baseurl}`, sensor).pipe(
      catchError(this.handleError<sensor>('addSensor', {} as sensor)) 
    );
  }

  // deleteSensor(id: number): Observable<sensor> {
  //   return this.http.delete<sensor>(`${this.baseurl}/${id}`).pipe(
  //     catchError(this.handleError<sensor>('deleteSensor', {} as sensor))
  //   );
  // }

  // PUT
  updateSensor(id: number, sensor: sensor): Observable<sensor> {
    return this.http.put<sensor>(`${this.baseurl}/${id}`, sensor, this.getHttpOptions()).pipe(
      catchError(this.handleError<sensor>('updateSensor', {} as sensor))
    );
  }

  private handleError<T>(operation = 'operation', result?: T) {
    return (error: HttpErrorResponse): Observable<T> => {
      console.error(`${operation} failed: ${error.message}`);
      // Show an alert or log to see what's happening
      alert(`Error: ${operation} failed: ${error.message}`);  // This will pop up an alert if the API fails
      return of(result as T);  // Fallback result (empty array in case of error)
    };
  }
  private getHttpOptions() {
    const headers = new HttpHeaders({ 'Content-Type': 'application/json' });
    return { headers };
  }
}
