import { Component, OnInit } from '@angular/core';
import { RouterModule } from '@angular/router';

import { SensorService } from '../sensor.service';
import { sensor } from '../sensor';

@Component({
  selector: 'app-sensors-list',
  standalone: true,
  imports: [RouterModule],
  templateUrl: './sensors-list.component.html',
  styleUrl: './sensors-list.component.css'
})
export class SensorsListComponent implements OnInit {

  sensors: sensor[] = [];

  constructor(private sensorService: SensorService) {}

  ngOnInit(): void {
    this.getSensors(); 
  }

  getSensors(): void {
    this.sensorService.getSensors().subscribe(sensors => this.sensors = sensors);
  }
}
