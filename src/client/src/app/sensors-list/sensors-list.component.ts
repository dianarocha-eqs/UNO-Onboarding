import { Component, OnInit } from '@angular/core';
import { RouterModule} from '@angular/router';

import { SensorService } from '../sensor.service';
import { sensor } from '../interfaces/sensor';
import { take } from 'rxjs';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-sensors-list',
  standalone: true,
  imports: [RouterModule, CommonModule],
  templateUrl: './sensors-list.component.html',
  styleUrl: './sensors-list.component.css'
})

export class SensorsListComponent {
}

// THIS IS FUTURE WORK. IGNORE !!!!!!!!!!!!
// export class SensorsListComponent implements OnInit {

//   sensors: sensor[] = [];

//   constructor(private sensorService: SensorService) {}

//   ngOnInit(): void {
//     this.getSensors();  
//   }

//   getSensors(): void {
//     this.sensorService.getSensors()
//       .pipe(take(1))  // Make sure to use take(1) for single emission
//       .subscribe({
//         next: (sensors) => {
//           this.sensors = sensors;
//           console.log(this.sensors)
//         },
//         error: (err) => console.error('Error loading sensors:', err)  // Log the error if any
//       });
//   }
// }
