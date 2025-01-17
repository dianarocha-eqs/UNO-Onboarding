import { Component } from '@angular/core';
import { SensorService } from '../sensor.service';
import { sensor } from '../sensor';

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
    color: 'Green', // Default color
    category: 'Temperature', // Default category
    description: '',
    visibility: 'private' // Default visibility
  };

  constructor(private sensorService: SensorService) {}

  addSensor(): void {
    // Ensure that the form is not empty or invalid
    if (!this.sensor.name.trim() || !this.sensor.category.trim()) {
      return; 
    }

    // Call the addSensor method from the service
    this.sensorService.addSensor(this.sensor).subscribe({
      next: (addedSensor) => {
        console.log('Sensor added successfully:', addedSensor);
      },
      error: (error) => {
        console.error('Error adding sensor:', error);
      }
    });
  }
}
