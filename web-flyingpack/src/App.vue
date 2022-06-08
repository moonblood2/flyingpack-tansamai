<template>
  <div>
    <b-navbar
        v-if="showHeader"
        fixed="top"
        style="background: #242f3d !important;"
        type="dark"
        variant="light"
    >
      <b-navbar-brand @click="onClickNavbarBrand">
        Flyingpack
      </b-navbar-brand>
      <h1 style="color: white" v-if="env.MODE === 'staging'"> เฉพาะทดสอบ </h1>
      <!-- Right aligned nav items -->
      <b-navbar-nav class="ml-auto">
        <b-nav-item-dropdown v-if="isAuthenticated" right>
          <!-- Using 'button-content' slot -->
          <template #button-content>
            <em>{{ $store.state.user.profile.name }}</em>
          </template>
          <b-dropdown-item @click="onClickLogout">Sign Out</b-dropdown-item>
        </b-nav-item-dropdown>
      </b-navbar-nav>
    </b-navbar>
    <div v-if="showSidebar" class="sidenav">
      <router-link v-if="checkPermission(UserRoles.SHOP | UserRoles.AGENT_NETWORK_MEMBER)" :to="paths.pos">POS</router-link>
      <router-link v-if="checkPermission(UserRoles.SHOP | UserRoles.AGENT_NETWORK_MEMBER)" :to="paths.reportOrderParcel">รายงานพัสดุ</router-link>
      <router-link v-if="checkPermission(UserRoles.SHOP)" :to="paths.reportOrderProduct">รายงานสินค้าอื่นๆ</router-link>
      <router-link v-if="checkPermission(UserRoles.AGENT_NETWORK_MEMBER)" :to="paths.dashBoardFulfillment">ภาพรวม fulfillment</router-link>
      <router-link v-if="checkPermission(UserRoles.AGENT_NETWORK_MEMBER)" :to="paths.reportFulfillment">รายงาน fulfillment</router-link>
      <router-link v-if="checkPermission(UserRoles.AGENT_NETWORK_MEMBER)" :to="paths.orderFulfillmentChecker">fulfillment checker</router-link>
      <router-link v-if="checkPermission(UserRoles.ACCOUNTING)" :to="paths.accountingReportFulfillment">รายงาน fulfillment</router-link>
      <router-link v-if="checkPermission(UserRoles.ACCOUNTING)" :to="paths.accountingReportCod">รายงาน COD</router-link>
    </div>
    <router-view :class="{main: showSidebar}"/>
  </div>
</template>

<script>
import Vue from 'vue'
import {mapGetters} from 'vuex'
import {BootstrapVue, BootstrapVueIcons} from 'bootstrap-vue';

import {AUTH_LOGOUT} from "@/store/actions/auth";
import {paths, includeHeader, includeSidebar, checkPermission, getDefaultPath} from "@/router";
import {UserRoles} from "@/entities/User";

import env from "@/constants/env";

Vue.use(BootstrapVue)
Vue.use(BootstrapVueIcons)

Vue.filter('numberWithCommas', function (value) {
  return value.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
})

export default {
  name: "app",
  data() {
    return {
      paths: paths,
      showHeader: includeHeader(this.$router.history.current.path),
      showSidebar: includeSidebar(this.$router.history.current.path),
      checkPermission: checkPermission,
      UserRoles: UserRoles,
      env: env,
    }
  },
  computed: {
    ...mapGetters(['isAuthenticated'])
  },
  watch: {
    $route(to) {
      this.showHeader = includeHeader(to.path);
      this.showSidebar = includeSidebar(to.path);
      if (to.path === "/") {
        this.$router.push(getDefaultPath());
      }
    }
  },
  methods: {
    onClickNavbarBrand() {
      if (this.$route.path !== getDefaultPath(this.$store.state.user.profile.role)) {
        this.$router.push(getDefaultPath(this.$store.state.user.profile.role));
      }
    },
    onClickLogout() {
      this.$store.dispatch(AUTH_LOGOUT).then(() => this.$router.push("/login"));
    },
    onClickLabel() {
      let routeData = this.$router.resolve({name: 'Label'});
      window.open(routeData.href, '_blank');
    },
  },
}
</script>

<style>
.navbar {
  height: 56px;
}

.navbar-brand:hover {
  cursor: pointer;
}

.sidenav {
  height: 100%;
  width: 200px;
  position: fixed;
  z-index: 1;
  top: 56px;
  left: 0;
  background-color: #fff;
  box-shadow: 0 0 5rem 0 rgba(12, 17, 22, 0.1) !important;
  overflow-x: hidden;
  padding-top: 20px;
}

.sidenav a {
  display: block;
  padding: 6px 8px 6px 16px;
  vertical-align: middle;
  text-decoration: none;
  font-size: 1rem;
  color: rgba(0, 0, 0, .5);
}

.sidenav a.router-link-exact-active {
  width: 100%;
  color: rgba(0, 0, 0, .8);
  border-bottom-right-radius: 25px;
  border-top-right-radius: 25px;
  background: linear-gradient(87deg, rgba(36, 47, 61, 0), rgba(36, 47, 61, 0.4)) !important;
}

.main {
  height: 100%;
  margin-left: 200px; /* Same as the width of the sidenav */
  padding-top: 80px;
  color: #212529;
}

</style>
