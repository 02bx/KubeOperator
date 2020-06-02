import { NgModule } from '@angular/core';
import { HostComponent } from './host.component';
import { HostListComponent } from './host-list/host-list.component';
import { HostCreateComponent } from './host-create/host-create.component';
import { HostDeleteComponent } from './host-delete/host-delete.component';
import {CoreModule} from '../../core/core.module';
import { HostDetailComponent } from './host-detail/host-detail.component';
import {SharedModule} from '../../shared/shared.module';



@NgModule({
  declarations: [HostComponent, HostListComponent, HostCreateComponent, HostDeleteComponent, HostDetailComponent],
    imports: [
        CoreModule,
        SharedModule,
    ]
})
export class HostModule { }
