<template>
  <img
      :id="id"
      :height="height"
      :width="width"
  />
</template>

<script>
import JsBarcode from "jsbarcode";

export default {
  name: "ReferenceNoBarcode",
  props: {
    id: String,
    referenceNo: String,
    width: Number,
    height: Number,
  },
  mounted() {
    JsBarcode(`#${this.id}`, this.referenceNo, {
      format: "CODE128",
      width: 5,
      height: 100,
      margin: 0,
      displayValue: false,
    });
  }
}
</script>