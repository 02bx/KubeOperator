import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {F5BigIpComponent} from './f5-big-ip.component';
import {CoreModule} from '../core/core.module';
import {TipModule} from '../tip/tip.module';
import {SharedModule} from '../shared/shared.module';

@NgModule({
  declarations: [F5BigIpComponent],
  imports: [
    CommonModule,
    CoreModule,
    TipModule,
    SharedModule
  ]
})
export class F5BigIpModule {
}
