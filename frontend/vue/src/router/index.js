import { createRouter, createWebHistory } from 'vue-router';
import ProductManagement from '../views/ProductManagement.vue'; 

const routes = [
  {
    path: '/',
    redirect: '/product-management' 
  },
  {
    path: '/product-management',
    name: 'ProductManagement',
    component: ProductManagement, 
  },
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
});

export default router;
