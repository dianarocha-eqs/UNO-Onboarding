import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router, RouterModule} from '@angular/router';

import { SensorService } from '../sensor.service';
import { sensor } from '../sensor';

@Component({
  selector: 'app-sensor-edit',
  standalone: true,
  imports: [RouterModule],
  templateUrl: './sensor-edit.component.html',
  styleUrl: './sensor-edit.component.css'
})
export class SensorEditComponent implements OnInit{

  sensor: sensor | null = null;

  constructor(  
    private sensorService: SensorService,
    private route: ActivatedRoute,
    private router: Router
  ) {}

  ngOnInit(): void {
    const id = +this.route.snapshot.paramMap.get('id')!;
    this.sensorService.getSensorById(id).subscribe(sensor => {this.sensor = sensor;});

  }
  
  updateSensor(): void {
    if (this.sensor) {
      this.sensorService.updateSensor(this.sensor.id, this.sensor).subscribe({
        next: (updatedSensor) => {
          console.log('Sensor updated successfully:', updatedSensor);
          this.router.navigate([`/sensors/${updatedSensor.id}`]);
        },
        error: (error) => {
          console.error('Error updating sensor:', error);
        }
      });
    }
  }
}
