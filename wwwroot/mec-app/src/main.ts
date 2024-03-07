import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'

import { CustomerService } from './services/customer.service'
import { DiscountService } from './services/discount.service'
import { OrderService } from './services/order.service'
import { ProductService } from './services/product.service'

import 'primevue/resources/themes/lara-light-green/theme.css'
import 'primeicons/primeicons.css'
import 'primeflex/primeflex.css'

import PrimeVue from 'primevue/config'

import AutoComplete from 'primevue/autocomplete'
import Calendar from 'primevue/calendar'
import Button from 'primevue/button'
import ConfirmDialog from 'primevue/confirmdialog'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import ColumnGroup from 'primevue/columngroup' // optional
import Dropdown from 'primevue/dropdown'
import InputGroup from 'primevue/inputgroup'
import InputGroupAddon from 'primevue/inputgroupaddon'
import InputMask from 'primevue/inputmask'
import InputNumber from 'primevue/inputnumber'
import InputText from 'primevue/inputtext'
import Paginator from 'primevue/paginator'
import Row from 'primevue/row' // optional
import Textarea from 'primevue/textarea'
import Toast from 'primevue/toast'

import ConfirmationService from 'primevue/confirmationservice'
import ToastService from 'primevue/toastservice'

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.use(PrimeVue)
app.component('AutoComplete', AutoComplete)
// eslint-disable-next-line vue/multi-word-component-names
app.component('Calendar', Calendar)
// eslint-disable-next-line vue/multi-word-component-names, vue/no-reserved-component-names
app.component('Button', Button)
app.component('ConfirmDialog', ConfirmDialog)
app.component('DataTable', DataTable)
// eslint-disable-next-line vue/multi-word-component-names
app.component('Column', Column)
app.component('ColumnGroup', ColumnGroup)
// eslint-disable-next-line vue/multi-word-component-names
app.component('Dropdown', Dropdown)
app.component('InputGroup', InputGroup)
app.component('InputGroupAddon', InputGroupAddon)
app.component('InputMask', InputMask)
app.component('InputNumber', InputNumber)
app.component('InputText', InputText)
// eslint-disable-next-line vue/multi-word-component-names
app.component('Paginator', Paginator)
// eslint-disable-next-line vue/multi-word-component-names
app.component('Row', Row)
// eslint-disable-next-line vue/multi-word-component-names, vue/no-reserved-component-names
app.component('Textarea', Textarea)
// eslint-disable-next-line vue/multi-word-component-names
app.component('Toast', Toast)

app.use(ToastService)
app.use(ConfirmationService)

app.provide('CustomerService', new CustomerService())
app.provide('DiscountService', new DiscountService())
app.provide('ProductService', new ProductService())
app.provide('OrderService', new OrderService())

app.mount('#app')
