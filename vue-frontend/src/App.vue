<!--Key Roles and Purposes of App.vue
1.Root Component: It acts as the root component in the component tree, meaning all other 
components will be children of App.vue.
2.Layout Structure: It often contains the basic layout structure, such as headers, 
footers, navigation bars, and the main content area where routed components are displayed.
3.Router Integration: It includes the <router-view /> component, which is the 
placeholder for rendering the matched components based on the current route.
4.Global Styles: It can include global styles that apply to the entire application. -->

<script setup>
  import { RouterLink, RouterView } from 'vue-router'
  import Header from './components/HeaderComponent.vue'
  import { ref, provide } from 'vue'
  //for profile.vue 
  const usernameProfile = ref(' ')
  provide('usernameProfile', usernameProfile)


</script>

<template>
  <Header />
  <RouterView />
</template> 
