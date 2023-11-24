import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';

import {CustomersListComponent} from "./customers-list/customers-list.component";
import {CustomerDetailsComponent} from "./customer-details/customer-details.component";
import {OrderDetailsComponent} from "./order-details/order-details.component";
import {OrderListComponent} from "./order-list/order-list.component";
import {ProductDetailsComponent} from "./product-details/product-details.component";
import {ProductListComponent} from "./product-list/product-list.component";
import {DiscountListComponent} from "./discount-list/discount-list.component";
import {DiscountDetailsComponent} from "./discount-details/discount-details.component";

const routes: Routes = [
  {path: 'customers', component: CustomersListComponent},
  {path: 'customers/:id', component: CustomerDetailsComponent},

  {path: 'products', component: ProductListComponent},
  {path: 'products/:id', component: ProductDetailsComponent},

  {path: 'discounts', component: DiscountListComponent},
  {path: 'discounts/:id', component: DiscountDetailsComponent},

  {path: 'orders', component: OrderListComponent},
  {path: 'orders/:id', component: OrderDetailsComponent},

  {path: '', redirectTo: '/orders', pathMatch: 'full'}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {
}
