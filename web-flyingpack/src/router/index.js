import Vue from 'vue';
import VueRouter from 'vue-router';

import store from "@/store";

import Login from '@/views/Login';
import Pos from '@/views/Pos';
import ReportOrderParcel from '@/views/ReportOrderParcel';
import ReportOrderProduct from '@/views/ReportOrderProduct';
import ReportFulfillment from "@/views/ReportFulfillment";
import ReportDupe from "@/views/ReportDupe";
import Label from '@/views/Label';
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
    label: '/label',
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
    if (!to.path.includes(paths.label) && !to.path.includes(paths.fulfillmentPackingSlip)) {
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
