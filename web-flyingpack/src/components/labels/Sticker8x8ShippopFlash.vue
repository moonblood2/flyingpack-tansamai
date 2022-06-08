<template>
  <div class="sticker-flash">
    <div class="logo monochrome">
      <ReferenceNoBarcode
          v-if="showReferenceNoBarcode && typeOfParcel(parcel) === parcelType.AN_PARCEL"
          :id="`${id}-reference-no`"
          :height="36"
          :reference-no="parcel.referenceNo"
          :width="126"
      />
      <img v-else-if="showProviderLogo || typeOfParcel(parcel) === parcelType.PARCEL"
           :src="require('@/assets/logo/provider/shippop.png')"
           height="15" width="100"
      />
      <img :src="require('@/assets/logo/courier/flash_express.png')" height="15" width="100"/>
    </div>
    <div class="barcode">
      <img :id="id" height="40" width="300"/>
      <b>{{ this.parcel.trackingCode }}</b>
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
    <div v-if="parcel.spOrderParcelShippopFlash" class="sorting">
      <div class="b1">
        <b>{{ parcel.spOrderParcelShippopFlash.sortCode }}</b> <br>
        <b>{{ parcel.spOrderParcelShippopFlash.dstCode }}</b>
      </div>
      <div class="b2" style="background: black; color: white">
        {{ parcel.spOrderParcelShippopFlash.sortingLineCode }}
      </div>
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
import ReferenceNoBarcode from "@/components/labels/ReferenceNoBarcode";

import {parcelType, typeOfParcel} from "@/entities/Parcel";

import JsBarcode from "jsbarcode";

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
  }
}
</script>

<style scoped>
.sticker-flash {
  width: 3.2in;
  height: 3.2in;
  padding: 1mm;
  background-color: white;
  margin-bottom: 10px;
  overflow: hidden;
  position: relative;
  box-sizing: border-box;
  page-break-after: always;
  border: 1px solid #666666;
}

.sticker-flash .logo {
  display: flex;
  flex-direction: row;
  justify-content: space-around;
  align-items: center;
  height: 20px;
  border-bottom: 0;
}

.sticker-flash .barcode {
  display: flex;
  flex-direction: column;
  align-items: center;
  height: 58px;
  border-bottom: 0;
}

.sticker-flash .origin {
  height: 40px;
  padding: 1mm;
  border: 1px solid #666;
  border-bottom: 0;
  font-size: 10px;
}

.sticker-flash .destination {
  height: 70px;
  padding: 1mm;
  border: 1px solid #666;
  border-bottom: 0;
  font-size: 12px;
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

.sticker-flash .other .item-list .item {

}

.sticker-flash .other .cod {
  width: 50%;
  margin: auto;
  text-align: center;
  font-size: 14px;
  font-weight: bold;
}

@media screen {
  .sticker-flash .barcode {
    margin-bottom: 2mm;
  }
  .sticker-flash .sorting .b1 {
    width: 60%;
    font-size: 12px;
    text-align: center;
  }
}

@media print {
  .sticker-flash .sorting .b2 {
    color: white;
    background: black;
    -webkit-print-color-adjust: exact;
  }
}
</style>