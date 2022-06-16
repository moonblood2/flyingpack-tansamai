<template>
  <div class="sticker">
    <div class="header">
      <div class="header-logo">
        <!-- <div class="header-logo-1">
          <img
            :src="require('@/assets/logo/courier/flash_express.png')"
            height="20"
            width="100"
          />
        </div> -->
        <div class="header-logo-2">
          <b width="100">{{ parcel.sortCode }}</b>
        </div>
      </div>

      <ReferenceNoBarcode
        :id="`${id}-reference-no`"
        :reference-no="parcel.referenceNo"
        :height="40"
        :width="350"
      />
      <div style="" class="tracking">
        {{
          this.parcel.trackingCode
            ? this.parcel.trackingCode
            : "no tacking code"
        }}
      </div>
    </div>
    <div class="title">
      <div>
        <b>DST : </b>{{ parcel.dstStoreName ? parcel.dstStoreName : "" }}
      </div>
    </div>
    <div class="address1">
      <div class="p-5">
        <b>ที่อยู่ผู้จัดส่งผู้ส่ง: </b><br />{{ `${parcel.origin.name}` }}
        <b>{{ `(${parcel.origin.phoneNumber})` }}</b>
        {{
          `${parcel.origin.address}  ${parcel.origin.district}  ${parcel.origin.state}  ${parcel.origin.province}`
        }}
        <b>{{ parcel.origin.postcode }}</b>
      </div>
      <div class="code1">{{ parcel.sortingLineCode }}</div>
    </div>
    <div class="address2">
      <div class="p-5">
        <b>ชื่อที่อยู่ผู้รับ: </b><br />{{ `${parcel.destination.name}` }}
        <b>{{ `(${parcel.destination.phoneNumber})` }}</b>
        {{
          `${parcel.destination.address}  ${parcel.destination.district}  ${parcel.destination.state}  ${parcel.destination.province}`
        }}
        <b>{{ parcel.destination.postcode }}</b>
      </div>
      <div class="code2">
        <VueQRCodeComponent :text="parcel.referenceNo" size="70" error-level="L">
        </VueQRCodeComponent>
      </div>
    </div>

    <div class="price_body">
      {{
        parcel.codAmount != 0 || parcel.codAmount != null
          ? parcel.codAmount + " THB"
          : "ไม่เก็บเงินปลายทาง"
      }}
    </div>
    <div class="foot">
      <div class="foot1">
        <div class="bg">
          <img
            :src="require('@/assets/logo/aura_logo.jpeg')"
            height="100"
            width="100"
          />
        </div>
      </div>
      <div class="shipping-code">
        <img
          :src="require('@/assets/logo/courier/flash_express.png')"
          height="20"
          width="100"
        />
        <ReferenceNoBarcode
          :id="`${id}-reference-no`"
          :reference-no="parcel.referenceNo"
          :height="40"
          :width="230"
        />
        <span><b>referenceNo : </b>{{ parcel.referenceNo }}</span>
      </div>
    </div>
    <div class="product">
      <span v-for="(v, i) of parcel.items" :key="i" class="item">
        {{ v.productCode }} (<b>{{ v.quantity }}</b
        >) {{ parcel.items.length > 1 ? "," : "" }}
      </span>
    </div>
    <div class="foot-print-date">
      <div class="note-open">
        ***กรุณาถ่ายวีดีโอ ขณะเปิดกล่องพัสดุ โดยไม่ตัดต่อ
        เพื่อเป็นหลักฐานในกรณีขอเคลมสินค้า
      </div>
    </div>
  </div>
</template>

<script>
import ReferenceNoBarcode from "@/components/labels/ReferenceNoBarcode";
import { parcelType, typeOfParcel } from "@/entities/Parcel";
import JsBarcode from "jsbarcode";

// Tansamai ADD
import VueQRCodeComponent from "vue-qrcode-component";

import "@/styles/common.css";
export default {
  name: "Sticker4x6",
  props: {
    id: String,
    showProviderLogo: Boolean,
    showReferenceNoBarcode: Boolean,
    parcel: Object,
    codDescription: String,
    timestamp: String,
  },
  components: {
    ReferenceNoBarcode,
    VueQRCodeComponent,
  },
  data() {
    return {
      parcelType: parcelType,
    };
  },
  mounted() {
    JsBarcode(`#${this.id}-reference-no`, this.parcel.trackingCode, {
      format: "CODE128",
      width: 10,
      height: 100,
      margin: 0,
      displayValue: false,
    });
  },
  methods: {
    typeOfParcel: typeOfParcel,
    getDatePrint() {
      const today = new Date();
      const date =
        today.getFullYear() +
        543 +
        "-" +
        (today.getMonth() + 1 <= 9
          ? "0" + (today.getMonth() + 1)
          : today.getMonth() + 1) +
        "-" +
        (today.getDate() <= 9 ? "0" + today.getDate() : today.getDate());
      const time = today.getHours() + ":" + today.getMinutes();
      const timestamp = date + " " + time;
      return timestamp;
    },
    logTest() {
      console.log("parcelType: ", parcelType);
    },
  },
};
</script>

<style scoped>
.sticker {
  width: 4in;
  height: 6in;
  padding: 2mm;
  margin-bottom: 2px;
  box-sizing: border-box;
  background-color: white;
  margin-left: auto;
  margin-right: auto;
  border: 1px solid black;
  page-break-after: always;
}
.header {
  /* border: 1px solid black;
  border-bottom: 0; */
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  width: 100%;
  padding: 0.5rem 0rem;
}
.header-logo {
  display: flex;
  justify-content: center;
  width: 100%;
}
.header-logo-1 {
  width: 30%;
  display: flex;
  justify-content: center;
}
.header-logo-2 {
  display: flex;
  justify-content: center;
  width: 70%;
  font-size: 25px;
  letter-spacing: 2px;
}
.title {
  border: 1px solid black;
  border-bottom: 0;
  width: 100%;
  text-align: left;
}
.title div {
  font-size: 12px;
  padding: 0.5rem 5px;
}

.tracking {
  font-size: 12px;
  text-transform: uppercase;
  padding-bottom: 2px;
}
.address1 {
  display: flex;
  flex-direction: row;
  border: 1px solid black;
  border-bottom: 0;
  width: 100%;
  height: 15%;
  font-size: 12px;
  text-align: left;
}
.address2 {
  display: flex;
  flex-direction: row;
  border: 1px solid black;
  border-bottom: 0;
  border-top: 0;
  width: 100%;
  /* height: 15%; */
  font-size: 12px;
  text-align: left;
}
.address1 div,
.address2 div {
  padding: 0.5rem 5px;
}
/* 
.code {
  display: flex;
  flex-direction: row;
  border: 1px solid black;
  border-bottom: 0;
  width: 100%;
  height: 15%;
}
.code1 {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  justify-content: center;
  align-items: center;
} */
.code1 {
  background-color: black;
  display: flex;
  justify-content: center;
  align-items: center;
  color: white;
  width: 50%;
  height: 80%;
  font-size: 50px;
}

.code2 {
  height: 50%;
}
.price_body {
  width: 100%;
  display: inline-block;
  padding: 0.5rem 1px;
  text-align: center;
  background-color: #000;
  color: #fff;
  font-size: 20px;
}
.shipping-code {
  width: 100%;
  display: inline-block;
  padding: 0.5rem 1px;
  text-align: center;
}
.shipping-code span {
  font-size: 12px;
}
.note-open {
  display: inline-block;
  width: 96%;
  font-size: 8px;
  padding: 0px 2%;
  text-align: right;
}

.foot {
  display: flex;
  flex-direction: row;
  border: 1px solid black;
  width: 100%;
  height: 20%;
}
.bg {
  display: flex;
  width: 100%;
  height: 100%;
  justify-content: center;
  align-items: center;
  border-right: 1px solid black;
  font-size: 10px;
}
.cod {
  display: flex;
  width: 100%;
  height: 100%;
  justify-content: center;
  align-items: center;
  border-right: 1px solid black;
  border-top: 1px solid black;
}

.foot1 {
  display: flex;
  flex-direction: column;
  width: 40%;
  height: 100%;
  min-width: 100px;
}
.product {
  display: inline-block;
  flex-direction: column;
  justify-content: space-between;
  width: 96%;
  height: auto;
  min-height: 50px;
  font-size: 12px;
  text-align: left;
  padding: 0.5rem 2%;
  border: 1px solid #000;
}
.foot-print-date {
  display: flex;
  justify-content: space-between;
  font-size: 8px;
}

@media print {
  @page {
    size: 5.5in 8.5in;
    size: landscape;
  }
}

@media print {
  .sticker {
    width: 50%;
    margin: auto;
    text-align: center;
    font-size: 14px;
    font-weight: bold;
  }
}
</style>
