import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'

import { ProductService } from './services/product.service'

import 'primevue/resources/themes/lara-light-green/theme.css'
import 'primeicons/primeicons.css'
import 'primeflex/primeflex.css'

import PrimeVue from 'primevue/config'

import Button from 'primevue/button'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import ColumnGroup from 'primevue/columngroup' // optional
import Dropdown from 'primevue/dropdown'
import InputGroup from 'primevue/inputgroup'
import InputGroupAddon from 'primevue/inputgroupaddon'
import InputNumber from 'primevue/inputnumber'
import InputText from 'primevue/inputtext'
import Paginator from 'primevue/paginator'
import Row from 'primevue/row' // optional
import Toast from 'primevue/toast'

import ToastService from 'primevue/toastservice'

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.use(PrimeVue)
// eslint-disable-next-line vue/multi-word-component-names, vue/no-reserved-component-names
app.component('Button', Button)
app.component('DataTable', DataTable)
// eslint-disable-next-line vue/multi-word-component-names
app.component('Column', Column)
app.component('ColumnGroup', ColumnGroup)
// eslint-disable-next-line vue/multi-word-component-names
app.component('Dropdown', Dropdown)
app.component('InputGroup', InputGroup)
app.component('InputGroupAddon', InputGroupAddon)
app.component('InputNumber', InputNumber)
app.component('InputText', InputText)
// eslint-disable-next-line vue/multi-word-component-names
app.component('Paginator', Paginator)
// eslint-disable-next-line vue/multi-word-component-names
app.component('Row', Row)
// eslint-disable-next-line vue/multi-word-component-names
app.component('Toast', Toast)

app.use(ToastService)

app.provide('ProductService', new ProductService())

app.mount('#app')
