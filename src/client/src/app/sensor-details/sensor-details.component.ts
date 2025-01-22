import { Component, OnInit } from '@angular/core';
import {  RouterModule} from '@angular/router';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-sensor-details',
  standalone: true,
  imports: [RouterModule, CommonModule],
  templateUrl: './sensor-details.component.html',
  styleUrl: './sensor-details.component.css'
})
export class SensorDetailsComponent  {
}