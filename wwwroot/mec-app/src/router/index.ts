import { createRouter, createWebHistory } from 'vue-router'

const DiscountListView = () => import('../views/DiscountListView.vue')
const DiscountView = () => import('../views/DiscountView.vue')

const ProductListView = () => import('../views/ProductListView.vue')
const ProductView = () => import('../views/ProductView.vue')

const CustomerListView = () => import('../views/CustomerListView.vue')
const CustomerView = () => import('../views/CustomerView.vue')

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/products',
      name: 'products',
      component: ProductListView
    },
    {
      path: '/products/:id',
      name: 'product',
      component: ProductView
    },
    {
      path: '/discounts',
      name: 'discounts',
      component: DiscountListView
    },
    {
      path: '/discounts/:id',
      name: 'discount',
      component: DiscountView
    },
    {
      path: '/customers',
      name: 'customers',
      component: CustomerListView
    },
    {
      path: '/customers/:id',
      name: 'customer',
      component: CustomerView
    }
    // {
    //   path: '/about',
    //   name: 'about',
    //   // route level code-splitting
    //   // this generates a separate chunk (About.[hash].js) for this route
    //   // which is lazy-loaded when the route is visited.
    //   component: () => import('../views/AboutView.vue')
    // }
  ]
})

export default router
