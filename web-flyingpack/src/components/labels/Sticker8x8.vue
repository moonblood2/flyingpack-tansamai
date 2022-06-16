<template>
  <Sticker8x8ShippopFlash
    :id="id"
    :cod-description="codDescription"
    :parcel="parcel"
    :show-provider-logo="true"
    :show-reference-no-barcode="false"
    class="sticker8x8"
  />
  <!-- <Sticker8x8ShippopFlash
    v-if="parcel.providerCode === provider.PROVIDER_SHIPPOP
      && parcel.courierCode === provider.SHIPPOP_COURIER_CODE_FLESH"
    :id="id"
    :cod-description="codDescription"
    :parcel="parcel"
    :show-provider-logo="true"
    :show-reference-no-barcode="false"
    class="sticker8x8"
  />
  <Sticker8x8ShippopKerry
    v-else-if="parcel.providerCode === provider.PROVIDER_SHIPPOP
      && parcel.courierCode === provider.SHIPPOP_COURIER_CODE_KERRY_PICKUP"
    :id="id"
    :cod-description="codDescription"
    :parcel="parcel"
    :show-provider-logo="true"
    class="sticker8x8"
  />
  <Sticker8x8ShippopEms
    v-else-if="parcel.providerCode === provider.PROVIDER_SHIPPOP
      && parcel.courierCode === provider.SHIPPOP_COURIER_CODE_EMS"
    :id="id"
    :cod-description="codDescription"
    :parcel="parcel"
    :show-provider-logo="true"
    :show-reference-no-barcode="false"
    class="sticker8x8"
  />
  <Sticker8x8ShippoLine
    v-else-if="parcel.providerCode === provider.PROVIDER_SHIPPOP
    && parcel.courierCode === provider.SHIPPOP_COURIER_CODE_MESSENGER"
    :id="id"
    :cod-description="codDescription"
    :parcel="parcel"
    :show-provider-logo="true"
    :show-reference-no-barcode="false"
    class="sticker8x8"
  />

  <h1
    v-else
    style="text-align: center"
  >ไม่มีให้แสดง</h1> -->
</template>

<script>
import Sticker8x8ShippopFlash from "@/components/labels/Sticker8x8ShippopFlash";
// import Sticker8x8ShippopKerry from "@/components/labels/Sticker8x8ShippopKerry";
// import Sticker8x8ShippopEms from "@/components/labels/Sticker8x8ShippopEms";
// import Sticker8x8ShippoLine from "@/components/labels/Sticker8x8ShippoLine";

import provider from "@/entities/Provider";

export default {
  name: "Sticker8x8",
  components: {
    Sticker8x8ShippopFlash,
    // Sticker8x8ShippopKerry,
    // Sticker8x8ShippopEms,
    // Sticker8x8ShippoLine,
  },
  data() {
    return {
      provider: provider,
    };
  },
  props: {
    parcel: Object,
    id: String,
  },
  computed: {
    codDescription: function() {
      if (this.parcel.codAmount > 0) {
        return `COD ${this.parcel.codAmount} บาท`;
      }
      return `ไม่เก็บเงินปลายทาง`;
    },
  },
};
</script>

<style scoped>
@media screen {
  .sticker8x8 {
    margin: 1mm auto;
  }
}

@media print {
  @page {
    size: portrait !important;
  }
}
</style>
