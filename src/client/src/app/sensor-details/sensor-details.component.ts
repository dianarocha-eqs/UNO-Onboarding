import { Component, OnInit } from '@angular/core';
import { RouterModule, ActivatedRoute} from '@angular/router';

import { SensorService } from '../sensor.service';
import { sensor } from '../interfaces/sensor';
import { take } from 'rxjs';

@Component({
  selector: 'app-sensor-details',
  standalone: true,
  imports: [RouterModule],
  templateUrl: './sensor-details.component.html',
  styleUrl: './sensor-details.component.css'
})
export class SensorDetailsComponent implements OnInit{

  sensor: sensor | null = null;

  constructor( 
    private sensorService: SensorService,
    private route: ActivatedRoute, 
  ) {}

  ngOnInit(): void {
    this.getSensorById();
  }

  getSensorById(): void {
    const id = +this.route.snapshot.paramMap.get('id')!;
    this.sensorService.getSensorById(id)
    .pipe(take(1))
    .subscribe({
      next: (data) => {
        this.sensor = data;
      },
      error: (err) => {
        console.error( 'Error retrieving sensor data');
      }
    });
  }


}
