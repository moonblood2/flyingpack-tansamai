<template>
  <div id="sticker_4x6">
    <div v-if="size === 'sticker-4x6'">
      <div class="preview">
        <button id="btn-preview" @click="print_preview">Print Preview</button>
      </div>
      <Sticker4x6
        v-for="(parcel, index) in parcels"
        :id="`label_${index}`"
        :key="`sticker-4x6_${index}`"
        :parcel="parcel"
      />
    </div>
    <div v-else>
      <h1>Sticker Not Found</h1>
    </div>
  </div>
</template>

<script>
import Sticker4x6 from "@/components/labels/Sticker4x6";
import printJs from "print-js";

export default {
  name: "Label4x6",
  components: {
    Sticker4x6,
  },
  data() {
    return {
      size: "sticker-4x6",
    };
  },
  created() {
    if (this.$route.query.size) {
      this.size = this.$route.query.size;
    }
  },
  computed: {
    parcels: function () {
      const { key } = this.$route.params;
      return this.$store.state.label.parcel[key];
    },
  },
  methods: {
    print_preview() {
      // Page Preview Style
      var style = "@page { size: 4in 6in } ";
      style += ".preview {display: none !important;} ";
      style +=
        "@media print { .sticker {width: 100%;height: 100%;margin: 0px !important;border: 0px !important;}}";

      // Sticker Style
      style +=
        ".sticker {padding: 2mm; box-sizing: border-box; background-color: white; page-break-after: always;}";
      style +=
        ".header {display: flex; flex-direction: column; justify-content: center; align-items: center; width: 100%; padding: 0.5rem 0rem;}";
      style +=
        ".header-logo {display: flex; justify-content: center; width: 100%;} .header-logo-1 { width: 30%; display: flex; justify-content: center;} .header-logo-2 { display: flex; justify-content: center; width: 70%; font-size: 25px; letter-spacing: 2px;}";
      style += ".title {border: 1px solid black; border-bottom: 0; width: 100%; text-align: left;} .title div { font-size: 12px; padding: 0.5rem 5px;}";
      style += ".tracking {font-size: 12px; text-transform: uppercase; padding-bottom: 2px;} .address1 {display: flex; flex-direction: row; border: 1px solid black; border-bottom: 0; width: 100%; height: 15%; font-size: 12px; text-align: left;} .address2 { display: flex; flex-direction: row; border: 1px solid black; border-bottom: 0; border-top: 0; width: 100%; font-size: 12px; text-align: left;} .address1 div, .address2 div {padding: 0.5rem 5px;}";
      style += ".code1 {background-color: black; display: flex; justify-content: center; align-items: center; color: white; width: 50%; height: 80%; font-size: 50px;} .code2 {height: 50%;}";
      style += ".price_body {width: 100%; display: inline-block; padding: 0.5rem 1px; text-align: center; background-color: #000; color: #fff; font-size: 20px; }";
      style += ".shipping-code {width: 100%; display: inline-block; padding: 0.5rem 1px; text-align: center;} .shipping-code span {font-size: 12px;}";
      style += ".note-open {display: inline-block; width: 96%; font-size: 8px; padding: 0px 2%; text-align: right;} .foot { display: flex; flex-direction: row; border: 1px solid black; width: 100%; height: 20%;}";
      style += ".foot {display: flex; flex-direction: row; border: 1px solid black; width: 100%; height: 20%;} .bg {display: flex; width: 100%; height: 100%; justify-content: center; align-items: center; border-right: 1px solid black; font-size: 10px;}";
      style += ".cod {display: flex; width: 100%; height: 100%; justify-content: center; align-items: center; border-right: 1px solid black; border-top: 1px solid black;} .foot1 {display: flex; flex-direction: column; width: 40%; height: 100%; min-width: 100px;}";
      style += ".product {display: inline-block; flex-direction: column; justify-content: space-between; width: 96%; height: auto; min-height: 50px; font-size: 12px; text-align: left; padding: 0.5rem 2%; border: 1px solid #000;}";
      style += ".foot-print-date {display: flex; justify-content: space-between; font-size: 8px;}";

      printJs({
        printable: "sticker_4x6",
        type: "html",
        style: style,
        honorColor: true,
        scanStyles: false,
        onPrintDialogClose: () => console.log("The print dialog was closed"),
        onError: (e) => console.log(e),
      });
    },
  },
};
</script>

<style scoped>
.preview {
  width: 100%;
  display: inline-block;
  text-align: center;
  padding: 10px 0px;
}
#btn-preview {
  display: inline-block;
  font-weight: 400;
  font-size: 1rem;
  padding: 0.5rem 1rem;
  line-height: 1.5;
  border-radius: 0.3rem;
  color: #fff;
  background-color: #007bff;
  text-align: center;
  vertical-align: middle;
  box-shadow: unset;
  border: 0px;
  cursor: pointer;
}
#btn-preview:hover {
  background-color: #0069d9;
}

/* @page {
  size: 4in 6in;
  margin: 0mm !important;
  padding: 0mm !important;
} */
</style>
