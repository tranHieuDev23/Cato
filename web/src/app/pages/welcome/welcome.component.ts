import { Component, OnInit } from '@angular/core';
import { PageTitleService } from '../../logic/page-title.service';

@Component({
  selector: 'app-welcome',
  standalone: true,
  templateUrl: './welcome.component.html',
  styleUrls: ['./welcome.component.scss'],
})
export class WelcomeComponent implements OnInit {
  constructor(private readonly pageTitleService: PageTitleService) {}

  ngOnInit(): void {
    this.pageTitleService.setTitle('Cato');
  }
}
