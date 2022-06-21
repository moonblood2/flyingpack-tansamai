<template>
  <div id="sticker_100x75">
    <div v-if="size === 'sticker-100x75'">
      <div class="preview">
        <button id="btn-preview" @click="print_preview">Print Preview</button>
      </div>
      <Sticker100x75
        v-for="(parcel, index) in parcels"
        :id="`label_${index}`"
        :key="`sticker-100x75_${index}`"
        :parcel="parcel"
      />
    </div>
    <div v-else>
      <h1>Sticker Not Found</h1>
    </div>
  </div>
</template>

<script>
import Sticker100x75 from "@/components/labels/Sticker100x75";
import printJs from "print-js";

export default {
  name: "Label100x75",
  components: {
    Sticker100x75,
  },
  data() {
    return {
      size: "sticker-100x75",
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
      //window.print();

      var style = "@page { size: 3.9in 3in } ";
      style += ".preview {display: none !important;} ";
      style +=
        "@media print { .sticker {width: 100% !important; height: 100% !important; margin: 0px !important; border: 0px !important;} .note-open {font-size: 8px !important;} .address1, .address2 {font-size: 12px;} .product {font-size: 10px;} .cod {height: 90% !important;}}";

      // Sticker Style
      style +=
        ".sticker {width: 3.93701in; height: 2.95276in; padding: 2mm; margin-bottom: 2px; box-sizing: border-box; background-color: white; margin-left: auto; margin-right: auto; border: 1px solid black; page-break-after: always;}";
      style +=
        ".header {display: flex; flex-direction: column; justify-content: center; align-items: center; width: 100%;}";
      style +=
        ".header-logo {display: flex; justify-content: center; width: 100%;} .header-logo-1 {width: 30%; display: flex; justify-content: center;} .header-logo-2 {display: flex; justify-content: center; width: 70%; font-size: 25px; letter-spacing: 2px;}";
      style +=
        ".title {border: 1px solid black; border-bottom: 0; width: 100%; text-align: left;} .title div {padding-left: 5px; font-size: 12px;}";
      style +=
        ".tracking {font-size: 12px; text-transform: uppercase; padding-bottom: 2px;}";
      style +=
        ".address1 {display: flex;flex-direction: row;border: 1px solid black;border-bottom: 0;width: 100%;height: 20%;font-size: 12px;text-align: left;}";
      style +=
        ".address2 {display: flex;flex-direction: row;border: 1px solid black;border-bottom: 0;border-top: 0;width: 100%;height: 20%;font-size: 12px;text-align: left;} .address1 div, .address2 div {padding-left: 5px;}";
      style +=
        ".code1 {background-color: black;display: flex;justify-content: center;align-items: center;color: white;width: 70%;height: 160%;font-size: 50px;} .code2 {display: flex;justify-content: flex-end;align-items: flex-end;color: red;width: 70%;height: 100%;font-size: 9px;padding-top: 2px;}";
      style +=
        ".note-open {display: flex;justify-content: flex-end;color: red;font-size: 11px;padding: 3px 0;} .foot {display: flex;flex-direction: row;border: 1px solid black;width: 100%;height: 20%;} .bg {display: flex;width: 100%;height: 100%;justify-content: center;align-items: center;border-right: 1px solid black;font-size: 10px;}";
      style += ".cod {display: flex;width: 100%;height: 100%;justify-content: center;align-items: center;border-right: 1px solid black;border-top: 1px solid black;} .foot1 {display: flex;flex-direction: column;width: 20%;height: 100%;}";
      style += ".product {display: flex;flex-direction: column;justify-content: space-between;width: 100%;height: 100%;font-size: 15px;padding: 2px;text-align: left;}";
      style += ".foot-print-date {display: flex;justify-content: space-between;font-size: 8px;}";

      printJs({
        printable: "sticker_100x75",
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
/* @page {
  size: 3.9in 3in;
  margin: 0px !important;
  padding: 0px !important;
} */
</style>

<style>
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

/* @media print {
  .preview {
    display: none !important;
  }
  .sticker {
    width: 100% !important;
    height: 100% !important;
    margin: 0px !important;
    border: 0px !important;
  }
  .note-open {
    font-size: 8px !important;
  }
  .address1,
  .address2 {
    font-size: 12px;
  }
  .product {
    font-size: 10px;
  }
  .cod {
    height: 90% !important;
  }
} */
</style>
