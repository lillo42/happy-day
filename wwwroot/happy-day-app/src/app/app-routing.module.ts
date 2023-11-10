import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { ListCustomersComponent } from "./list-customers/list-customers.component";
import { CustomerDetailsComponent } from "./customer-details/customer-details.component";

const routes: Routes = [
  {path: 'customers', component: ListCustomersComponent },
  {path: 'customers/:id', component: CustomerDetailsComponent },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
