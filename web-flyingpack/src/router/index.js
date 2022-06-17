import Vue from 'vue';
import VueRouter from 'vue-router';

import store from "@/store";

import Login from '@/views/Login';
import Pos from '@/views/Pos';
import ReportOrderParcel from '@/views/ReportOrderParcel';
import ReportOrderProduct from '@/views/ReportOrderProduct';
import ReportFulfillment from "@/views/ReportFulfillment";
import ReportDupe from "@/views/ReportDupe";

// Tansamai ADD
import Label from '@/views/Label';
import Label4x6 from '@/views/Label4x6';
import Label8x8 from '@/views/Label8x8';
import Label100x75 from '@/views/Label100x75';

import DashboardFulfillment from "@/views/DashboardFulfillment";
import OrderFulfillmentChecker from "@/views/OrderFulfillmentChecker";
import AccountingReportFulfillment from "@/views/AccountingReportFulfillment";
import AccountingReportCod from "@/views/AccountingReportCod";
import FulfillmentPackingSlip from "@/views/FulfillmentPackingSlip";

import {UserRoles} from "@/entities/User";

Vue.use(VueRouter);

export const paths = {
    login: '/login',
    pos: '/pos',
    reportOrderParcel: '/report-order-parcel',
    reportOrderProduct: '/report-order-product',
    reportFulfillment: '/report-fulfillment',
    reportDupe: '/report-dupe',
    // Tansamai ADD
    label: '/label',
    label4x6: '/Label4x6',
    label8x8: '/Label8x8',
    label100x75: '/Label100x75',
    dashBoardFulfillment: '/dashboard-fulfillment',
    orderFulfillmentChecker: '/order-fulfillment-checker',
    accountingReportFulfillment: '/accounting-report-fulfillment',
    accountingReportCod: '/accounting-report-cod',

    fulfillmentPackingSlip: '/fulfillment-packing-slip',
}

const getDefaultPath = () => {
    switch (store.state.user.profile.role) {
        case UserRoles.SHOP:
            return paths.reportOrderParcel;
        case UserRoles.AGENT_NETWORK_MEMBER:
            return paths.reportFulfillment;
        case UserRoles.ACCOUNTING:
            return paths.accountingReportFulfillment;
        default:
            return "/login";
    }
}

const includeHeader = (path) => [
    paths.login,
    paths.pos,
    paths.reportOrderParcel,
    paths.reportOrderProduct,
    paths.reportFulfillment,
    paths.reportDupe,
    paths.dashBoardFulfillment,
    paths.orderFulfillmentChecker,
    paths.accountingReportFulfillment,
    paths.accountingReportCod,
].includes(path);

const includeSidebar = (path) => [
    paths.reportOrderParcel,
    paths.reportOrderProduct,
    paths.reportFulfillment,
    paths.reportDupe,
    paths.dashBoardFulfillment,
    paths.orderFulfillmentChecker,
    paths.accountingReportFulfillment,
    paths.accountingReportCod,
].includes(path);

const checkPermission = (permittedRole) => {
    if (!store.state.user.profile.role) return false;
    return !((permittedRole & store.state.user.profile.role) === 0);
}

const ifAuthenticated = (permittedRole) => {
    return (to, from, next) => {
        if (store.getters.isAuthenticated && checkPermission(permittedRole)) {
            next();
            return;
        }
        next("/login");
    }
};

const ifNotAuthenticated = (to, from, next) => {
    if (!store.getters.isAuthenticated || getDefaultPath() === paths.login) {
        next();
        return;
    }
    next(getDefaultPath());
};

const redirectToDefaultPath = (to, from, next) => {
    next(getDefaultPath());
}

const routes = [
    {
        path: "/",
        beforeEnter: redirectToDefaultPath,
    },
    {
        path: '/login',
        component: Login,
        beforeEnter: ifNotAuthenticated
    },
    {
        path: '/pos',
        component: Pos,
        beforeEnter: ifAuthenticated(UserRoles.SHOP | UserRoles.AGENT_NETWORK_MEMBER),
    },
    {
        path: '/report-order-parcel',
        component: ReportOrderParcel,
        beforeEnter: ifAuthenticated(UserRoles.SHOP | UserRoles.AGENT_NETWORK_MEMBER),
    },
    {
        path: '/report-order-product',
        component: ReportOrderProduct,
        beforeEnter: ifAuthenticated(UserRoles.SHOP),
    },
    {
        path: '/report-fulfillment',
        component: ReportFulfillment,
        beforeEnter: ifAuthenticated(UserRoles.AGENT_NETWORK_MEMBER)
    },
    {
        path: '/report-dupe',
        component: ReportDupe,
        beforeEnter: ifAuthenticated(UserRoles.AGENT_NETWORK_MEMBER)
    },
    {
        path: '/label/:key',
        name: 'Label',
        component: Label,
        beforeEnter: ifAuthenticated(UserRoles.SHOP | UserRoles.AGENT_NETWORK_MEMBER),
    },
    {
        path: '/label4x6/:key',
        name: 'Label4x6',
        component: Label4x6,
        beforeEnter: ifAuthenticated(UserRoles.SHOP | UserRoles.AGENT_NETWORK_MEMBER),
    },
    {
        path: '/label8x8/:key',
        name: 'Label8x8',
        component: Label8x8,
        beforeEnter: ifAuthenticated(UserRoles.SHOP | UserRoles.AGENT_NETWORK_MEMBER),
    },
    {
        path: '/label100x75/:key',
        name: 'Label100x75',
        component: Label100x75,
        beforeEnter: ifAuthenticated(UserRoles.SHOP | UserRoles.AGENT_NETWORK_MEMBER),
    },
    {
        path: `${paths.fulfillmentPackingSlip}/:key`,
        name: "PackingSlip",
        component: FulfillmentPackingSlip,
        beforeEnter: ifAuthenticated(UserRoles.AGENT_NETWORK_MEMBER),
    },
    {
        path: '/dashboard-fulfillment',
        component: DashboardFulfillment,
        beforeEnter: ifAuthenticated(UserRoles.AGENT_NETWORK_MEMBER),
    },
    {
        path: '/order-fulfillment-checker',
        component: OrderFulfillmentChecker,
        beforeEnter: ifAuthenticated(UserRoles.AGENT_NETWORK_MEMBER),
    },
    {
        path: paths.accountingReportFulfillment,
        component: AccountingReportFulfillment,
        beforeEnter: ifAuthenticated(UserRoles.ACCOUNTING),
    },
    {
        path: paths.accountingReportCod,
        component: AccountingReportCod,
        beforeEnter: ifAuthenticated(UserRoles.ACCOUNTING),
    },
]

const router = new VueRouter({
    mode: 'history',
    base: process.env.BASE_URL,
    routes
});

router.afterEach((to) => {
    if (!to.path.includes(paths.label) && 
    // Tansamai ADD
    !to.path.includes(paths.label4x6) && 
    !to.path.includes(paths.label8x8) && 
    !to.path.includes(paths.label100x75) && 
    !to.path.includes(paths.fulfillmentPackingSlip)) {
        import('bootstrap/dist/css/bootstrap.min.css');
        import('bootstrap-vue/dist/bootstrap-vue.min.css');
    }
});

export {
    includeHeader,
    includeSidebar,
    checkPermission,
    getDefaultPath,
}

export default router
