import { Component } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { UserService } from '../user.service';
import { Router, RouterModule } from '@angular/router';
import { NgIf } from '@angular/common';

import { MatInputModule } from '@angular/material/input';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatSelectModule } from '@angular/material/select';
import { MatOptionModule } from '@angular/material/core';
import { MatButtonModule } from '@angular/material/button';

@Component({
  selector: 'app-user-add',
  standalone: true,
  imports: [RouterModule, NgIf, ReactiveFormsModule, MatFormFieldModule,MatInputModule,MatSelectModule,MatOptionModule,MatButtonModule],
  templateUrl: './user-add.component.html',
  styleUrls: ['./user-add.component.css'],
})
export class UserAddComponent {
  userForm: FormGroup;
  selectedImage: string | ArrayBuffer | null = null;
  errorMessage: string | null = null;

  constructor(
    private fb: FormBuilder,
    private userService: UserService,
    private router: Router
  ) {
    // Initialize the user form with validation
    this.userForm = this.fb.group({
      name: ['', [Validators.required]],
      email: ['', [Validators.required, Validators.email]],
      phone: ['', [Validators.required, this.phoneValidator]],
      picture: [''], // Optional field
      role: [''] // Optional field 
    });
  }

  /**
   * Retrieves the error message for the email input field.
   * The error message is shown if the email is required or has an invalid format 
   * after the field has been touched.
   * 
   * @returns string | null - The error message for the email field or null if valid.
   */
  get emailError(): string | null {
    const email = this.userForm.get('email');
    if (email?.hasError('required') && email?.touched) {
      return 'Email is required';
    } else if (email?.hasError('email') && email?.touched) {
      return 'Invalid email format';
    }
    return null;
  }

  /**
   * Retrieves the appropriate error message for the phone field.
   * @returns A string containing the error message, or null if no error.
   */
  get phoneError(): string | null {
    const phone = this.userForm.get('phone');
    if (phone?.hasError('required') && phone?.touched) {
      return 'Phone number is required';
    } else if (phone?.hasError('invalidPhone') && phone?.touched) {
      return 'Invalid phone number format';
    }
    return null;
  }

  /**
   * Validates the phone number format.
   * @param control The form control containing the phone number input.
   * @returns An error object if invalid, otherwise null.
   */
  phoneValidator(control: any) {
    const phoneRegExp = /^\+?[1-9]\d{1,14}$/;
    if (control.value && !phoneRegExp.test(control.value)) {
      return { invalidPhone: true };
    }
    return null;
  }

   /**
   * Handles file selection for the profile picture.
   * @param event The file input change event.
   */
  onImageSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    if (input.files && input.files[0]) {
      const reader = new FileReader();
      reader.onload = () => {
        this.selectedImage = reader.result; 
      };
      reader.readAsDataURL(input.files[0]);
    }
  }

   /**
   * Resets the form and clears selected image.
   */
  onCancel(): void {
    this.userForm.reset();
    this.selectedImage = null; 
  }

  /**
   * Submits the user form data to the server.
   */
  onSubmit(): void {
    if (this.userForm.invalid) {
      return;
    }

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

    this.userService.addUser(user, token).subscribe({
      next: (response) => {
        this.errorMessage = null;
        console.log('User created successfully!', response);
        this.userForm.reset();
        Object.keys(this.userForm.controls).forEach((key) => {
          this.userForm.controls[key].setErrors(null);  
          this.userForm.controls[key].markAsUntouched();
        });
  
        this.selectedImage = null;
      },
      error: (err) => {
        this.errorMessage = err.error?.error || 'Request failed';
      }
    });
  }
}
