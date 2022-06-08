<template>
  <div class="sticker">
    <div class="logo">
      <b-img :src="require('@/assets/logo/flyingpack4x6.png')"/>
      <!--      <b-img :src="require('@/assets/logo/rosegold4x6.png')"/>-->
      <h3>Rose Gold</h3>
    </div>
    <div class="barcode">
      <svg :id="id"></svg>
    </div>
    <div class="origin">
      <p>
        <b>ผู้ส่ง</b> {{ `${parcel.origin.name}  (${parcel.origin.phoneNumber})` }}
        <br/>
        {{ `${parcel.origin.address}  ${parcel.origin.district}  ${parcel.origin.state}  ${parcel.origin.province}` }}
        <b>{{ parcel.origin.postcode }}</b>
      </p>
    </div>
    <div class="destination">
      <p>
        <span class="first-line">
          <b>ผู้รับ</b> {{ `${parcel.destination.name}  (${parcel.destination.phoneNumber})` }}
        </span>
        <br/>
        {{
          `${parcel.destination.address}  ${parcel.destination.district}  ${parcel.destination.state}  ${parcel.destination.province}`
        }}
        <b>{{ parcel.destination.postcode }}</b>
      </p>
    </div>
    <div class="other">
      <h6>
        หมายเหตุ: {{ codDescription }}
      </h6>
    </div>
  </div>
</template>

<script>
import JsBarcode from "jsbarcode"

export default {
  name: "Sticker8x8",
  props: {
    parcel: Object,
    id: String,
  },
  mounted() {
    JsBarcode(`#${this.id}`, this.parcel.trackingCode, {
      width: 2,
      height: 98,
      fontSize: 20,
      displayValue: true,
    });
  },
  computed: {
    codDescription: function () {
      if (this.parcel.codAmount > 0) {
        return `ยอดเก็บเงินปลายทาง ${this.parcel.codAmount} บาท`
      }
      return `ไม่เก็บเงินปลายทาง`
    }
  }
}
</script>

<style scoped>
.sticker {
  width: 4in;
  height: 6in;
  padding: 2mm;
  margin-bottom: 2px;
  background-color: white;
  page-break-after: always;
  margin-left: auto;
  margin-right: auto;
}

.logo {
  display: flex;
  flex-direction: row;
  justify-content: space-around;
  align-items: center;
  height: 15%;
  padding: 10px;
  border: 1px solid black;
  border-bottom: 0;
}

.barcode {
  display: flex;
  flex-direction: column;
  align-items: center;
  height: 20%;
  border: 1px solid black;
  border-bottom: 0;
}

.origin {
  height: 20%;
  padding: 1mm;
  border: 1px solid black;
  border-bottom: 0;
  font-size: 18px;
}

.destination {
  height: 30%;
  padding: 1mm;
  border: 1px solid black;
  border-bottom: 0;
  font-size: 18px;
}

.destination .first-line {
  font-size: 20px;
}

.other {
  height: 15%;
  border: 1px solid black;
  font-size: 14px;
}

@media print {
  .sticker {
    position: absolute;
    height: 100%;
    width: 100%;
    padding: 20px;
    page-break-after: always;
  }

  .sticker .origin {
    font-size: 20px;
  }

  .sticker .destination {
    font-size: 20px;
  }

  .sticker .destination .first-line {
    font-size: 22px;
  }

  .sticker .other {
    font-size: 16px;
  }
}
</style>