import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import {CustomerListsComponent} from "./customer-lists/customer-lists.component";

const routes: Routes = [
  { path: 'customers', component: CustomerListsComponent }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
