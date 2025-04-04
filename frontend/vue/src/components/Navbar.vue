<template>
  <a-layout style="min-height: 100vh">
    <a-layout-header style="background: #001529; padding: 0;">
      <div style="text-align: center; color: white; font-size: 20px; line-height: 64px;">
        Product Management
      </div>
    </a-layout-header>

    <a-layout style="padding: 0 24px 24px">
      <a-layout-content
        style="padding: 24px; min-height: 280px; background: #fff"
      >
        <!-- ใช้ <router-view /> ที่นี่ -->
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script>
export default {
  name: 'Navbar'
}
</script>

<style scoped>
a-layout-header {
  height: 64px; 
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
}

a-layout-content {
  margin-top: 10px; 
}
</style>
