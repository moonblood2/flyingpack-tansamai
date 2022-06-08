import Vue from 'vue'
import Vuex from 'vuex'

import user from "./modules/user"
import auth from "./modules/auth"
import label from './modules/label';
import packingSlip from './modules/packing-slip';
import pos from "./modules/pos";
import slip from "./modules/slip";

Vue.use(Vuex)

export default new Vuex.Store({
    modules: {
        user,
        auth,
        label,
        packingSlip,
        pos,
        slip,
    }
})
