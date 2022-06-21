<template>
  <div class="sticker">
    <div class="header">
      <div class="header-logo">
        <div class="header-logo-1">
          <img
            :src="require('@/assets/logo/courier/flash_express.png')"
            height="20"
            width="100"
          />
        </div>
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
        <b>ผู้ส่ง: </b>{{ `${parcel.origin.name}` }}
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
        <b>ผู้รับ: </b>{{ `${parcel.destination.name}` }}
        <b>{{ `(${parcel.destination.phoneNumber})` }}</b>
        {{
          `${parcel.destination.address}  ${parcel.destination.district}  ${parcel.destination.state}  ${parcel.destination.province}`
        }}
        <b>{{ parcel.destination.postcode }}</b>
      </div>
      <div class="code2"></div>
    </div>
    <!-- <div class="code">
      <div class="code1">
        <div>{{ parcel.referenceNo }}</div>
        <div>xxx-xxx</div>
      </div>
      <div class="code2">RO2</div>
    </div> -->
    <div class="foot">
      <div class="foot1">
        <div class="bg">
          {{
            parcel.codAmount != 0 || parcel.codAmount != null
              ? parcel.codAmount + " บาท"
              : "ไม่เก็บเงินปลายทาง"
          }}
        </div>
        <div class="cod">{{ codAmount == 0 ? "" : "COD" }}</div>
      </div>
      <!-- <div>{{ codDescription }}</div> -->
      <!-- <div class="foot2">
        <div>
          {{
            parcel.codAmount !== 0 || parcel.codAmount !== null
              ? parcel.codAmount + " บาท"
              : "ไม่เก็บเงินปลายทาง"
          }}
        </div>
      </div> -->
      <div class="product">
        <span v-for="(v, i) of parcel.items" :key="i" class="item">
          {{ v.productCode }} (<b>{{ v.quantity }}</b
          >) {{ parcel.items.length > 1 ? "," : "" }}
        </span>
        <div class="note-open">***กรุณาถ่ายวิดีโอขณะแกะสินค้า</div>
      </div>
    </div>
    <div class="foot-print-date">
      <span><b>referenceNo : </b>{{ parcel.referenceNo }}</span>
      <span><b>Print : </b>{{ getDatePrint() }}</span>
    </div>
  </div>
</template>

<script>
import ReferenceNoBarcode from "@/components/labels/ReferenceNoBarcode";
import { parcelType, typeOfParcel } from "@/entities/Parcel";
import JsBarcode from "jsbarcode";

import "@/styles/common.css";
export default {
  name: "Sticker100x75",
  props: {
    id: String,
    showProviderLogo: Boolean,
    showReferenceNoBarcode: Boolean,
    parcel: Object,
    codDescription: String,
    timestamp: String,
    paperSize: {
      type: String,
      default: '3.9in 3in !important'
    }
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
/* @page {
  size: 3.9in 3in !important;
  margin: 0px !important;
  padding: 0px !important;
} */
</style>

<style scoped>
.sticker {
  width: 3.93701in;
  height: 2.95276in;
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
  padding-left: 5px;
  font-size: 12px;
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
  height: 20%;
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
  height: 20%;
  font-size: 12px;
  text-align: left;
}
.address1 div,
.address2 div {
  padding-left: 5px;
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
  width: 70%;
  height: 160%;
  font-size: 50px;
}

.code2 {
  display: flex;
  justify-content: flex-end;
  align-items: flex-end;
  color: red;
  width: 70%;
  height: 100%;
  font-size: 9px;
  padding-top: 2px;
}
.note-open {
  display: flex;
  justify-content: flex-end;
  color: red;
  font-size: 11px;
  padding: 3px 0;
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
  width: 20%;
  height: 100%;
}
.product {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  width: 100%;
  height: 100%;
  font-size: 15px;
  padding: 2px;
  text-align: left;
}
.foot-print-date {
  display: flex;
  justify-content: space-between;
  font-size: 8px;
}
</style>
