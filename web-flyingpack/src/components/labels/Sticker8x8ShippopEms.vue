<template>
  <div class="sticker-ems">
    <div style="display: flex; justify-content: flex-end; align-items: center; height: 10px">
      <img v-if="showProviderLogo || typeOfParcel(parcel) === parcelType.PARCEL"
           :src="require('@/assets/logo/provider/shippop.png')"
           height="6" width="33"
      />
    </div>
    <div class="header">
      <div class="logo">
        <img :src="require('@/assets/logo/courier/thailandpost_ems.png')" height="14" width="120"/>
      </div>
      <div class="text">
        <div>บริการจัดส่งสินค้า (e-Parcel)</div>
        <div>ใบอนุญาตสำหรับลูกค้าธุรกิจ เลขที่ รล.135/2564</div>
        <div>ชำระค่าฝากส่งตามที่ ปณท กำหนด</div>
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
      <div class="postcode">
        {{ parcel.destination.postcode }}
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
        <ReferenceNoBarcode
            v-if="showReferenceNoBarcode && typeOfParcel(parcel) === parcelType.AN_PARCEL"
            :id="`${id}-reference-no`"
            :height="36"
            :reference-no="parcel.referenceNo" :width="126"
        />
      </div>
    </div>
    <div class="barcode">
      <img :id="id" height="40" width="280"/>
      <b>{{ this.parcel.trackingCode }}</b>
    </div>
  </div>
</template>

<script>
import ReferenceNoBarcode from "@/components/labels/ReferenceNoBarcode";

import JsBarcode from "jsbarcode";
import {parcelType, typeOfParcel} from "@/entities/Parcel";

export default {
  name: "Sticker8x8FLash",
  props: {
    id: String,
    parcel: Object,
    showProviderLogo: Boolean,
    showReferenceNoBarcode: Boolean,
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
      width: 4,
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
.sticker-ems {
  width: 3.2in;
  height: 3.2in;
  padding: 0 2mm;
  background-color: white;
  margin-bottom: 10px;
  overflow: hidden;
  position: relative;
  box-sizing: border-box;
  page-break-after: always;
  border: 1px solid #666666;
}

.sticker-ems .header {
  display: flex;
  flex-direction: row;
  align-items: center;
  height: 35px;
  border: 1px solid #666;
  border-bottom: 0;
}

.sticker-ems .header .logo {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 50%;
}

.sticker-ems .header .text {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 50%;
  text-align: center;
  font-size: 7px;
}

.sticker-ems .origin {
  height: 30px;
  padding: 1mm;
  border: 1px solid #666;
  border-bottom: 0;
  font-size: 10px;
}

.sticker-ems .destination {
  position: relative;
  height: 70px;
  padding: 1mm;
  border: 1px solid #666;
  border-bottom: 0;
  font-size: 12px;
}

.sticker-ems .destination .first-line {
  font-size: 13px;
}

.sticker-ems .destination .postcode {
  position: absolute;
  display: flex;
  justify-content: center;
  align-items: center;
  width: 50px;
  height: 22px;
  bottom: 0;
  right: 0;
  font-size: 15px;
  font-weight: bold;
  border-width: 1px 0 0 1px;
  border-color: #666;
  border-style: solid;
}

.sticker-ems .other {
  display: flex;
  height: 70px;
  border: 1px solid #666;
}

.sticker-ems .other .item-list {
  width: 50%;
  font-size: 9px;
  border-right: 1px solid #666;
}

.sticker-ems .other .item-list .item {

}

.sticker-ems .barcode {
  margin-top: 10px;
  display: flex;
  flex-direction: column;
  align-items: center;
  height: 58px;
  border-bottom: 0;
}

.sticker-ems .other .cod {
  width: 50%;
  margin: auto;
  text-align: center;
  font-size: 14px;
  font-weight: bold;
}

@media screen {
  .sticker-ems {
    margin: 5mm auto;
  }
}
</style>