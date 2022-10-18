import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import {CustomerListsComponent} from "./customer-lists/customer-lists.component";
import {ProductListComponent} from "./product-list/product-list.component";

const routes: Routes = [
  { path: 'customers', component: CustomerListsComponent },
  { path: 'products', component: ProductListComponent },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
