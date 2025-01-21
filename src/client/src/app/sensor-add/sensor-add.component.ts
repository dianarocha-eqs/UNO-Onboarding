import { Component } from '@angular/core';
import { SensorService } from '../sensor.service';
import { sensor, sensorCategory, sensorColor, sensorVisibility } from '../interfaces/sensor';
import { take } from 'rxjs';

@Component({
  selector: 'app-sensor-add',
  standalone: true,
  imports: [],
  templateUrl: './sensor-add.component.html',
  styleUrl: './sensor-add.component.css'
})
export class SensorAddComponent {

  // Define the sensor object using the interface, without 'id' and 'files' initially.
  sensor: Omit<sensor, 'id'> = {
    name: '',
    color: sensorColor.Green, // Default color
    category: sensorCategory.Temperature, // Default category
    description: '',
    visibility: sensorVisibility.Public // Default visibility
  };

  constructor(private sensorService: SensorService) {}

  addSensor(): void {
    // Ensure that the form is not empty or invalid
    if (!this.sensor.name.trim() || !this.sensor.category.trim()) {
      return; 
    }

    // Call the addSensor method from the service
    this.sensorService.addSensor(this.sensor)
    .pipe(take(1)) // Automatically unsubscribes after the first emission
    .subscribe({
      next: (addedSensor) => {
        console.log('Sensor added successfully:', addedSensor);
      },
      error: (error) => {
        console.error('Error adding sensor:', error);
      }
    });
  }
}
