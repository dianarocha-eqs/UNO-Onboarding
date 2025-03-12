import { Component, OnInit } from '@angular/core';
import { User } from '../interfaces/user';
import { FormControl, FormGroup, ReactiveFormsModule } from '@angular/forms';
import { UserService } from '../user.service';
import { RouterModule } from '@angular/router';
import { NgIf , NgFor} from '@angular/common';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatOptionModule } from '@angular/material/core';
import { MatButtonModule } from '@angular/material/button';
import { MatListModule} from '@angular/material/list';

@Component({
  selector: 'app-users-list',
  imports: [RouterModule, NgIf, ReactiveFormsModule, NgFor, MatFormFieldModule,MatInputModule,MatSelectModule,MatOptionModule,MatButtonModule, MatInputModule, MatListModule],
  templateUrl: './users-list.component.html',
  styleUrl: './users-list.component.css'
})

export class UserListingComponent implements OnInit {
  users: User[] = [];
  searchControl = new FormControl('');
  errorMessage: string | null = null;
  sortOrder: number = 0;  // Default to no order (0)
  
  userForm: FormGroup = new FormGroup({
    role: new FormControl(null),
  });

  constructor(private userService: UserService) {}

  ngOnInit(): void {
    this.fetchUsers();
  }

  /**
   * Fetches users from the API based on search input and sort order.
   */
  fetchUsers(): void {

    const user = this.userForm.value;
    let token = localStorage.getItem('Authorization'); 
    let role = localStorage.getItem('Role'); 
    
    user.role = role === '1';

    if (!token) {
      this.errorMessage = "Authorization token is missing from localStorage.";
      return; 
    }

    if (!role) {
      this.errorMessage = "Role is missing from localStorage.";
      return; 
    }
    const searchQuery = this.searchControl.value || '';
    const sortValue = this.sortOrder;

    const body = { search: searchQuery, sort: sortValue };

    this.userService.listUsers(token!, body, user).subscribe({
      next: (users) => {
        this.users = users;
        this.errorMessage = null;
      },
      error: (err) => {
        this.errorMessage = err.error?.error || 'Failed to fetch users';
      },
    });
  }

  /**
   * Triggers user search and updates the list based on sort order.
   */
  onSearch(): void {
    this.fetchUsers();
  }

  /**
   * Triggered when sort option changes.
   */
  onSortChange(sortValue: number): void {
    this.sortOrder = sortValue;
    this.fetchUsers();
  }
}
