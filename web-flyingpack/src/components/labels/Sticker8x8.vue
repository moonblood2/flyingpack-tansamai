<!-- <template>
  <Sticker8x8ShippopFlash
    :id="id"
    :cod-description="codDescription"
    :parcel="parcel"
    :show-provider-logo="true"
    :show-reference-no-barcode="false"
    class="sticker8x8"
  />
  <Sticker8x8ShippopFlash
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
  >ไม่มีให้แสดง</h1>
</template>

<script>
import Sticker8x8ShippopFlash from "@/components/labels/Sticker8x8ShippopFlash";
import Sticker8x8ShippopKerry from "@/components/labels/Sticker8x8ShippopKerry";
import Sticker8x8ShippopEms from "@/components/labels/Sticker8x8ShippopEms";
import Sticker8x8ShippoLine from "@/components/labels/Sticker8x8ShippoLine";

import provider from "@/entities/Provider";

export default {
  name: "Sticker8x8",
  components: {
    Sticker8x8ShippopFlash,
    Sticker8x8ShippopKerry,
    Sticker8x8ShippopEms,
    Sticker8x8ShippoLine,
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
    codDescription: function () {
      if (this.parcel.codAmount > 0) {
        return `COD ${this.parcel.codAmount} บาท`;
      }
      return `ไม่เก็บเงินปลายทาง`;
    },
  },
};
</script> -->

<!-- Tansamai ADD -->

<template>
  <div class="sticker-flash">
    <div class="logo monochrome">
      <ReferenceNoBarcode
        v-if="
          showReferenceNoBarcode &&
          typeOfParcel(parcel) === parcelType.AN_PARCEL
        "
        :id="`${id}-reference-no`"
        :height="36"
        :reference-no="parcel.referenceNo"
        :width="126"
      />
      <img
        v-else-if="
          showProviderLogo || typeOfParcel(parcel) === parcelType.PARCEL
        "
        :src="require('@/assets/logo/provider/shippop.png')"
        height="15"
        width="100"
      />
      <img
        :src="require('@/assets/logo/courier/flash_express.png')"
        height="15"
        width="100"
      />
    </div>
    <div class="barcode">
      <img :id="id" height="40" width="300" />
      <b>{{ this.parcel.trackingCode }}</b>
    </div>
    <div class="origin">
      <b>ผู้ส่ง</b>
      {{ `${parcel.origin.name}  (${parcel.origin.phoneNumber})` }}
      {{
        `${parcel.origin.address}  ${parcel.origin.district}  ${parcel.origin.state}  ${parcel.origin.province}`
      }}
      <b>{{ parcel.origin.postcode }}</b>
    </div>
    <div class="destination">
      <span class="first-line">
        <b>ผู้รับ</b>
        {{ `${parcel.destination.name}  (${parcel.destination.phoneNumber})` }}
      </span>
      {{
        `${parcel.destination.address}  ${parcel.destination.district}  ${parcel.destination.state}  ${parcel.destination.province}`
      }}
      <b>{{ parcel.destination.postcode }}</b>
    </div>
    <div v-if="parcel.spOrderParcelShippopFlash" class="sorting">
      <div class="b1">
        <b>{{ parcel.spOrderParcelShippopFlash.sortCode }}</b> <br />
        <b>{{ parcel.spOrderParcelShippopFlash.dstCode }}</b>
      </div>
      <div class="b2" style="background: black; color: white">
        {{ parcel.spOrderParcelShippopFlash.sortingLineCode }}
      </div>
    </div>
    <div class="other">
      <div class="item-list">
        <span v-for="(v, i) of parcel.items" :key="i" class="item">
          {{ v.productCode }} (<b>{{ v.quantity }}</b
          >),
        </span>
      </div>
      <div class="cod">
        {{ codDescription }}
      </div>
    </div>
  </div>
</template>

<script>
import ReferenceNoBarcode from "@/components/labels/ReferenceNoBarcode";

import { parcelType, typeOfParcel } from "@/entities/Parcel";

import JsBarcode from "jsbarcode";

import "@/styles/common.css";

export default {
  name: "Sticker8x8",
  props: {
    id: String,
    showProviderLogo: Boolean,
    showReferenceNoBarcode: Boolean,
    parcel: Object,
    codDescription: String,
  },
  components: {
    ReferenceNoBarcode,
  },
  data() {
    return {
      parcelType: parcelType,
    };
  },
  mounted() {
    JsBarcode(`#${this.id}`, this.parcel.trackingCode, {
      format: "CODE128",
      width: 10,
      height: 100,
      margin: 0,
      displayValue: false,
    });
  },
  methods: {
    typeOfParcel: typeOfParcel,
  },
};
</script>

<style scoped>
.sticker-flash {
  width: 3.2in;
  height: 3.2in;
  padding: 1mm;
  background-color: white;
  margin: 0px auto 10px auto;
  overflow: hidden;
  position: relative;
  box-sizing: border-box;
  page-break-after: always;
  border: 1px solid #666666;
}

.sticker-flash .logo {
  display: inline-block;
  width: 100%;
  flex-direction: row;
  justify-content: space-around;
  /* align-items: center; */
  /* height: 20px; */
  text-align: center;
  border-bottom: 0;
}

.sticker-flash .barcode {
  width: 90%;
  margin: 0px 5%;
  display: inline-block;
  flex-direction: column;
  align-items: center;
  /* height: 58px; */
  border-bottom: 0;
  text-align: center;
}

.sticker-flash .barcode img {
  width: 100%;
}

.sticker-flash .origin {
  height: auto;
  padding: 1mm;
  border: 1px solid #666;
  border-bottom: 0;
  font-size: 11px;
}

.sticker-flash .destination {
  height: auto;
  padding: 1mm;
  border: 1px solid #666;
  border-bottom: 0;
  font-size: 11px;
}

.sticker-flash .sorting {
  display: flex;
  height: 35px;
  border: 1px solid #666;
  border-bottom: 0;
}

.sticker-flash .sorting .b1 {
  width: 60%;
  font-size: 14px;
  text-align: center;
}

.sticker-flash .sorting .b2 {
  width: 40%;
  color: white;
  font-size: 25px;
  font-weight: bold;
  text-align: center;
  background: black;
}

.sticker-flash .destination .first-line {
  font-size: 13px;
}

.sticker-flash .other {
  display: flex;
  height: 50px;
  border: 1px solid #666;
}

.sticker-flash .other .item-list {
  width: 50%;
  font-size: 9px;
  border-right: 1px solid #666;
}

.sticker-flash .other .cod {
  width: 50%;
  margin: auto;
  text-align: center;
  font-size: 14px;
  font-weight: bold;
}
</style>
