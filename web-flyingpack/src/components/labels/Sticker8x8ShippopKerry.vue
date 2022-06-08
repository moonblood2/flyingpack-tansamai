<template>
  <div class="sticker-kerry">
    <div class="header">
      <div class="logo monochrome">
        <img v-if="showProviderLogo || typeOfParcel(parcel) === parcelType.PARCEL"
             :src="require('@/assets/logo/provider/shippop.png')"
             height="21" width="120"
        />
        <ReferenceNoBarcode
            v-else-if="showReferenceNoBarcode &&typeOfParcel(parcel) === parcelType.AN_PARCEL"
            :id="`${id}-reference-no`"
            :reference-no="parcel.referenceNo"
            :height="36" :width="126"
        />
        <!-- Use shippop_k.png for temporary in the promotion season, will back to use kerry_express.png on 2021/08/31 -->
        <img class="monochrome" :src="require('@/assets/logo/courier/shippop_k.png')" height="36" width="126"/>
      </div>
      <div class="qrcode">
        <canvas :id="id" height="108" width="108"></canvas>
        <b>{{ this.parcel.trackingCode }}</b>
      </div>
    </div>

    <div class="origin">
      <b>ผู้ส่ง</b> {{ `${parcel.origin.name}  (${parcel.origin.phoneNumber})` }}
      {{ `${parcel.origin.address}  ${parcel.origin.district}  ${parcel.origin.state}  ${parcel.origin.province}` }}
      <b>{{ parcel.origin.postcode }}</b>
    </div>
    <div class="destination">
        <span class="first-line">
          <b>ผู้รับ</b> {{ `${parcel.destination.name}  (${parcel.destination.phoneNumber})` }}
        </span>
      {{
        `${parcel.destination.address}  ${parcel.destination.district}  ${parcel.destination.state}  ${parcel.destination.province}`
      }}
      <b>{{ parcel.destination.postcode }}</b>
    </div>
    <div class="other">
      <div class="item-list">
          <span v-for="(v, i) of parcel.items" :key="i" class="item">
            {{ v.productCode }} (<b>{{ v.quantity }}</b>),
          </span>
      </div>
      <div class="cod">
        {{ codDescription }}
      </div>
    </div>
  </div>
</template>

<script>
import QRCode from "qrcode";
import {parcelType, typeOfParcel} from "@/entities/Parcel";
import ReferenceNoBarcode from "@/components/labels/ReferenceNoBarcode";

import "@/styles/common.css";

export default {
  name: "Sticker8x8FLash",
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
    }
  },
  mounted() {
    QRCode.toCanvas(document.getElementById(this.id), this.parcel.trackingCode, {
      width: 100,
      margin: 0,
    });
  },
  methods: {
    typeOfParcel: typeOfParcel,
  }
}
</script>

<style scoped>
.sticker-kerry {
  width: 3.2in;
  height: 3.2in;
  padding: 2mm;
  background-color: white;
  margin-bottom: 10px;
  overflow: hidden;
  position: relative;
  box-sizing: border-box;
  page-break-after: always;
  border: 1px solid #666666;
}

.sticker-kerry .header {
  display: flex;
  height: 130px;
  border-bottom: 0;
}

.sticker-kerry .header .logo {
  display: flex;
  flex-direction: column;
  justify-content: space-around;
  align-items: center;
  width: 180px;
}

.sticker-kerry .header .qrcode {
  display: flex;
  flex-direction: column;
  justify-content: space-around;
  align-items: center;
  width: 120px;
  font-size: 14px;
  text-align: center;
}

.sticker-kerry .origin {
  height: 40px;
  padding: 1mm;
  border: 1px solid #666;
  border-bottom: 0;
  font-size: 10px;
}

.sticker-kerry .destination {
  height: 70px;
  padding: 1mm;
  border: 1px solid #666;
  border-bottom: 0;
  font-size: 12px;
}

.sticker-kerry .destination .first-line {
  font-size: 13px;
}

.sticker-kerry .other {
  display: flex;
  height: 50px;
  border: 1px solid #666;
}

.sticker-kerry .other .item-list {
  width: 50%;
  font-size: 9px;
  border-right: 1px solid #666;
}

.sticker-kerry .other .item-list .item {

}

.sticker-kerry .other .cod {
  width: 50%;
  margin: auto;
  text-align: center;
  font-size: 14px;
  font-weight: bold;
}
</style>