import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { ListCustomersComponent } from "./list-customers/list-customers.component";
import { CustomerDetailsComponent } from "./customer-details/customer-details.component";
import { ProductDetailsComponent } from "./product-details/product-details.component";
import { ProductListComponent } from "./product-list/product-list.component";

const routes: Routes = [
  { path: 'customers', component: ListCustomersComponent },
  { path: 'customers/:id', component: CustomerDetailsComponent },

  { path: 'products', component: ProductListComponent },
  { path: 'products/:id', component: ProductDetailsComponent },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
