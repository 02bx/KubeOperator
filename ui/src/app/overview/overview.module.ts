import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {OverviewComponent} from './overview.component';
import {CoreModule} from '../core/core.module';
import {DescribeComponent} from './describe/describe.component';
import {AppsComponent} from './apps/apps.component';
import {SharedModule} from '../shared/shared.module';

@NgModule({
  declarations: [OverviewComponent, DescribeComponent, AppsComponent],
  imports: [
    CommonModule,
    CoreModule,
    SharedModule
  ]
})
export class OverviewModule {
}
