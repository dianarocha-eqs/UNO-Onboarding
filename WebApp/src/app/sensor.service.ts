import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable, of } from 'rxjs';
import { catchError, tap } from 'rxjs/operators';
import { sensor } from './sensor'; // Assuming sensor interface is defined elsewhere

@Injectable({
  providedIn: 'root'
})
export class SensorService {
  private baseurl = 'api/sensors'; // URL to your backend API for sensors

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
    return (error: any): Observable<T> => {
      // Log the error and optionally show an error message
      console.error(error); // Optionally, integrate a logging service
      return of(result as T); // Let the app continue by returning a fallback result
    };
  }
  private getHttpOptions() {
    const headers = new HttpHeaders({ 'Content-Type': 'application/json' });
    return { headers };
  }
}
