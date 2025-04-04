<template>
  <a-layout style="min-height: 100vh">
    <!-- Navbar -->
    <a-layout-header style="background: #001529; padding: 0;">
      <div style="text-align: center; color: white; font-size: 20px; line-height: 64px;">
        Product Management
      </div>
    </a-layout-header>

    <!-- Content -->
    <a-layout style="padding: 0 20px 20px">
      <a-layout-content
        style="padding: 24px; min-height: 280px; background: #fff; margin-top: 30px;" >

        <!-- ใช้ <router-view /> ที่นี่ -->
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script>
export default {
  name: 'App',
  mounted() {
    this.$nextTick(() => {
      console.log("Connect Successfully");
    });
  }
};
</script>


