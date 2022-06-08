<template>
  <div class="slip">
    <div class="message" v-html="message">
    </div>
    <div class="ref-no">
      <b> ReferenceNo. </b> {{ item.referenceNo }}
    </div>
    <div class="items">
      <div v-for="(e, i) in item.items" :key="i">
        <b style="font-size: 10px">
          {{ e.productCode }} - {{ e.quantity }}
        </b>
        <div style="font-size: 10px">
          <span v-for="(s, j) in e.serialNumbers" :key="j">
            {{ s.serialNumberStart}} {{ s.serialNumberEnd ? `- ${s.serialNumberEnd}`: '' }}
            {{ j + 1 !== e.serialNumbers.length ? ',': '' }}
          </span>
        </div>
        <hr/>
      </div>
    </div>
    <div class="cod-qr">
      <div class="cod">
        COD: {{ item.codAmount }} บาท
      </div>
      <div class="qr">
        <canvas :id="`${id}-qr`" height="108" width="108"></canvas>
      </div>
    </div>
    <div class="greeting-message" v-html="slip.message">
    </div>
  </div>
</template>

<script>
import QRCode from "qrcode";

export default {
  name: "FulfillmentPackingSlip",
  computed: {
    id: function () {
      return this.$route.params.key;
    },
    message: function () {
      let x = `ใบจัดส่งสินค้า (โรสโกลด์)
      บริษัท โรสโกลด์ (ไทยแลนด์) จำกัด
      คลังกระจายสินค้า 10/5 ถ.พระรามที่ 2 ซอย 30
      แขวงจอมทอง เขตจอมทอง กทม. 10150`;
      x = x.replaceAll('\n', '<br/>');
      return x;
    },
    item: function () {
      const {key} = this.$route.params;
      return this.$store.state.packingSlip.items[key];
    },
    slip: function () {
      let slip = this.$store.state.slip.slip;
      slip.message = slip.message.replaceAll('\n', '<br/>');
      return slip;
    }
  },
  mounted() {
    QRCode.toCanvas(document.getElementById(`${this.id}-qr`), this.item.trackingCode, {
      width: 100,
      margin: 0,
    });
  }
}
</script>

<style scoped>
.slip {
  display: flex;
  flex-direction: column;
  width: 8cm;
  background-color: white;
  margin: 1mm auto;
  overflow: hidden;
  position: relative;
  box-sizing: border-box;
  page-break-after: always;
  border: 1px solid #666666;
}

.message {
  vertical-align: center;
  height: 116px;
  padding-top: 10px;
  text-align: center;
  border-bottom: 1px solid black;
}

.greeting-message {
  font-size: 14px;
  padding: 5px;
  border-bottom: 1px solid black;
}

.ref-no {
  height: 30px;
  text-align: center;
  padding-top: 10px;
  border-bottom: 1px solid black;
}

.items {
  min-height: 50px;
  padding: 10px;
  border-bottom: 1px solid black;
}

.cod-qr {
  display: flex;
  flex-direction: row;
  height: 115px;
  border-bottom: 1px solid black;
}

.cod-qr .cod {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 50%;
  border-right: 1px solid black;
}

.cod-qr .qr {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 50%;
}

.barcode {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  height: 75px;
}
</style>