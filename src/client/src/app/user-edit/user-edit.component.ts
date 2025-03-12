import { Component } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { UserService } from '../user.service';
import { User } from '../interfaces/user';

import { MatInputModule } from '@angular/material/input';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatSelectModule } from '@angular/material/select';
import { MatOptionModule } from '@angular/material/core';
import { MatButtonModule } from '@angular/material/button';
import { NgClass, NgIf } from '@angular/common';

@Component({
  selector: 'app-user-edit',
  standalone: true,
  imports: [RouterModule, NgIf, NgClass, ReactiveFormsModule, MatFormFieldModule,MatInputModule,MatSelectModule,MatOptionModule,MatButtonModule],
  templateUrl: './user-edit.component.html',
  styleUrl: './user-edit.component.css'
})
export class UserEditComponent {

  // user : User | undefined;
  // userForm: FormGroup;
  // selectedImage: string | ArrayBuffer | null = null;
  // errorMessage: string | null = null;

  // constructor(
  //   private fb: FormBuilder,
  //   private userService: UserService,
  //   private router: Router
  // ) {
  //   // Initialize the user form with validation
  //   this.userForm = this.fb.group({
  //     name: [this.user?.name, [Validators.required]],
  //     email: [this.user?.email, [Validators.required, Validators.email]],
  //     phone: [this.user?.phone, [Validators.required, this.phoneValidator]],
  //     picture: [this.user?.picture], // Optional field
  //     role: [this.user?.role] // Optional field 
  //   });
  // }

  // /**
  //  * Retrieves the error message for the email input field.
  //  * The error message is shown if the email is required or has an invalid format 
  //  * after the field has been touched.
  //  * 
  //  * @returns string | null - The error message for the email field or null if valid.
  //  */
  // get emailError(): string | null {
  //   const email = this.userForm.get('email');
  //   if (email?.hasError('required') && email?.touched) {
  //     return 'Email is required';
  //   } else if (email?.hasError('email') && email?.touched) {
  //     return 'Invalid email format';
  //   }
  //   return null;
  // }

  // /**
  //  * Retrieves the appropriate error message for the phone field.
  //  * @returns A string containing the error message, or null if no error.
  //  */
  // get phoneError(): string | null {
  //   const phone = this.userForm.get('phone');
  //   if (phone?.hasError('required') && phone?.touched) {
  //     return 'Phone number is required';
  //   } else if (phone?.hasError('invalidPhone') && phone?.touched) {
  //     return 'Invalid phone number format';
  //   }
  //   return null;
  // }

  // /**
  //  * Validates the phone number format.
  //  * @param control The form control containing the phone number input.
  //  * @returns An error object if invalid, otherwise null.
  //  */
  // phoneValidator(control: any) {
  //   const phoneRegExp = /^\+?[1-9]\d{1,14}$/;
  //   if (control.value && !phoneRegExp.test(control.value)) {
  //     return { invalidPhone: true };
  //   }
  //   return null;
  // }

  //   /**
  //  * Handles file selection for the profile picture.
  //  * @param event The file input change event.
  //  */
  // onImageSelected(event: Event): void {
  //   const input = event.target as HTMLInputElement;
  //   if (input.files && input.files[0]) {
  //     const reader = new FileReader();
  //     reader.onload = () => {
  //       this.selectedImage = reader.result; 
  //     };
  //     reader.readAsDataURL(input.files[0]);
  //   }
  // }

  //   /**
  //  * Resets the form and clears selected image.
  //  */
  // onCancel(): void {
  //   this.userForm.reset();
  //   this.selectedImage = null; 
  // }

  // /**
  //  * Submits the user form data to the server.
  //  */
  // onSubmit(): void {
  //   if (this.userForm.invalid) {
  //     return;
  //   }

  //   const user = this.userForm.value;
  //   const token = localStorage.getItem('Authorization'); 
  //   const role = localStorage.getItem('Role'); 
  //   const uuid = localStorage.getItem('UUID'); 
    
  //   user.role = role === '1';

  //   if (!token) {
  //     this.errorMessage = "Authorization token is missing from localStorage.";
  //     return; 
  //   }

  //   if (!role) {
  //     this.errorMessage = "Role is missing from localStorage.";
  //     return; 
  //   }

  //   if (!uuid) {
  //     this.errorMessage = "UUID is missing from localStorage.";
  //     return; 
  //   }

  //   this.userService.editUser(user, token).subscribe({
  //     next: (response) => {
  //       this.errorMessage = null;
  //       console.log('User edited successfully!', response);
  //       this.userForm.reset();
  //       Object.keys(this.userForm.controls).forEach((key) => {
  //         this.userForm.controls[key].setErrors(null); 
  //         this.userForm.controls[key].markAsUntouched();
  //       });
  
  //       this.selectedImage = null;
  //     },
  //     error: (err) => {
  //       this.errorMessage = err.error?.error || 'Request failed';
  //     },
  //   });
  // }
  
}
