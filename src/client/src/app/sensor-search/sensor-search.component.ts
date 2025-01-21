import { Component, OnInit } from '@angular/core';
import { Observable, Subject } from 'rxjs';
import { debounceTime, distinctUntilChanged, switchMap } from 'rxjs/operators';
import { SensorService } from '../sensor.service';
import { sensor } from '../interfaces/sensor';
import { RouterModule } from '@angular/router';


@Component({
  selector: 'app-sensor-search',
  standalone: true,
  imports: [RouterModule],
  templateUrl: './sensor-search.component.html',
  styleUrls: ['./sensor-search.component.css'],
})
export class SensorSearchComponent implements OnInit {
  sensors$!: Observable<sensor[]>; 
  private searchTerms = new Subject<string>(); 

  constructor(private sensorService: SensorService) {}

  search(term: string): void {
    this.searchTerms.next(term); 
  }

  ngOnInit(): void {
    this.sensors$ = this.searchTerms.pipe(
      // Wait 100ms after each keystroke before considering the term
      debounceTime(100),

      // Ignore new term if same as previous term
      distinctUntilChanged(),

      // Switch to new search observable each time the term changes
      switchMap((term: string) => this.sensorService.searchSensors(term)), // Use your search service
    );
  }
}
