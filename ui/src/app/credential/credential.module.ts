import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {CredentialListComponent} from './credential-list/credential-list.component';
import {CredentialCreateComponent} from './credential-create/credential-create.component';
import {CredentialComponent} from './credential.component';

import {CoreModule} from '../core/core.module';
import {SharedModule} from '../shared/shared.module';

@NgModule({
  declarations: [CredentialListComponent, CredentialCreateComponent, CredentialComponent],
  imports: [
    CommonModule,

    CoreModule,
    SharedModule
  ]
})
export class CredentialModule {
}
