import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {ClusterComponent} from './cluster.component';
import {ClusterListComponent} from './cluster-list/cluster-list.component';
import {CoreModule} from '../core/core.module';
import {ClusterService} from './cluster.service';
import {ClusterDetailComponent} from './cluster-detail/cluster-detail.component';
import {ClusterCreateComponent} from './cluster-create/cluster-create.component';
import {TipModule} from '../tip/tip.module';

@NgModule({
  declarations: [ClusterComponent, ClusterListComponent, ClusterDetailComponent, ClusterCreateComponent],
  imports: [
    CommonModule,
    TipModule,
    CoreModule
  ],
  providers: [ClusterService]
})
export class ClusterModule {
}
