import { Component, OnInit } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import { AuthService } from '../../auth.service';
import { User } from '../../interfaces/user'
import { NgIf } from '@angular/common';
 
@Component({
  selector: 'app-login-page',
  standalone: true,
  imports: [RouterModule],
  templateUrl: './login-page.component.html',
  styleUrl: './login-page.component.css'
})

export class LoginPageComponent implements OnInit {

  user: Partial<User> = {
    email: '',
    password: ''
  }
  errorMessage: string | null = null;

  constructor(private authService: AuthService, private router: Router) {}

  ngOnInit(): void {
    this.login(this.user.email!, this.user.password!);
  }

  login(email: string, password: string) {
    if (email && password) {
      this.authService.login(email, password).subscribe({
        next: (response) => {
          console.log('Login successful', response);
        },
        error: (error) => {
          console.error('Login failed', error);
        }
      });
    } else {
      console.log('Please fill in both fields');
    }
  }
}
