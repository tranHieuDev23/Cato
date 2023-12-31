import {
  ChangeDetectorRef,
  Component,
  HostListener,
  OnDestroy,
  OnInit,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, RouterOutlet } from '@angular/router';
import { NzIconModule } from 'ng-zorro-antd/icon';
import { NzLayoutModule } from 'ng-zorro-antd/layout';
import { SideMenuComponent } from './components/side-menu/side-menu.component';
import { AccountService } from './logic/account.service';
import { NzPageHeaderModule } from 'ng-zorro-antd/page-header';
import { PageTitleService } from './logic/page-title.service';
import { Subscription } from 'rxjs';
import { NzSpaceModule } from 'ng-zorro-antd/space';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [
    CommonModule,
    RouterOutlet,
    NzIconModule,
    NzLayoutModule,
    SideMenuComponent,
    RouterModule,
    NzPageHeaderModule,
    NzSpaceModule,
  ],
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent implements OnInit, OnDestroy {
  public collapsed = false;
  public pageTitle = 'Cato';
  public windowWidth = 0;

  private titleChangedSubscription: Subscription;

  constructor(
    private readonly accountService: AccountService,
    private readonly pageTitleService: PageTitleService,
    private readonly changeDetector: ChangeDetectorRef
  ) {
    this.windowWidth = window.innerWidth;
    this.titleChangedSubscription =
      this.pageTitleService.titleChanged.subscribe((title) => {
        this.pageTitle = title;
        this.changeDetector.detectChanges();
      });
  }

  ngOnInit(): void {
    this.accountService.getSessionAccount().then();
  }

  ngOnDestroy(): void {
    this.titleChangedSubscription.unsubscribe();
  }

  @HostListener('window: resize')
  public onWindowResize(): void {
    this.windowWidth = window.innerWidth;
  }
}
