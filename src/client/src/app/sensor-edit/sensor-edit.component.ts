import { Component, OnInit } from '@angular/core';
import { RouterModule} from '@angular/router';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-sensor-edit',
  standalone: true,
  imports: [RouterModule, RouterModule, CommonModule, FormsModule],
  templateUrl: './sensor-edit.component.html',
  styleUrl: './sensor-edit.component.css'
})
export class SensorEditComponent  {
}