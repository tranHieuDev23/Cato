import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterOutlet } from '@angular/router';
import { ApiService } from './dataaccess';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, RouterOutlet, FormsModule],
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss',
})
export class AppComponent {
  public message = '';
  public echoMessage = '';

  constructor(private readonly apiService: ApiService) {}

  public async onEchoClicked(): Promise<void> {
    const response = await this.apiService.echo({
      message: this.message,
    });
    this.echoMessage = response?.message || 'No message returned!';
  }
}
