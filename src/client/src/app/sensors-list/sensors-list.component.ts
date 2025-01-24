import { Component, OnInit } from '@angular/core';
import { RouterModule} from '@angular/router';
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
