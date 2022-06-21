<template>
  <div id="sticker_8x8">
    <div v-if="size === 'sticker-8x8'">
      <div class="preview">
        <button id="btn-preview" @click="print_preview">Print Preview</button>
      </div>
      <Sticker8x8
        v-for="(parcel, index) in parcels"
        :id="`label_${index}`"
        :key="`sticker-8x8_${index}`"
        :parcel="parcel"
      />
    </div>
    <div v-else>
      <h1>Sticker Not Found</h1>
    </div>
  </div>
</template>

<script>
import Sticker8x8 from "@/components/labels/Sticker8x8";
import printJs from "print-js";

export default {
  name: "Label8x8",
  components: {
    Sticker8x8,
  },
  data() {
    return {};
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
      var style = "@page { size: 3.2in } ";
      style += ".preview {display: none !important;} ";
      style +=
        "@media print { .sticker-flash {width: 100% !important; height: 100% !important; margin: 0px !important; border: 0px !important;} .sticker-flash .logo img { width: 30% !important; height: 20px !important;}}";

      // Sticker Style
      style +=
        ".sticker-flash {width: 3.2in; height: 3.2in; padding: 1mm; margin: 0px auto 10px auto; overflow: hidden; position: relative; box-sizing: border-box; page-break-after: always; border: 1px solid #666666; background-color: #fff;}";
      style += ".sticker-flash .logo {display: inline-block; width: 100%; flex-direction: row; justify-content: space-around; text-align: center; border-bottom: 0;}";
      style += ".sticker-flash .barcode {width: 90%; margin: 0px 5%; display: inline-block; flex-direction: column; align-items: center; border-bottom: 0; text-align: center;} .sticker-flash .barcode img {width: 100%;}";
      style += ".sticker-flash .origin {height: auto; padding: 1mm; border: 1px solid #666; border-bottom: 0; font-size: 11px;} .sticker-flash .destination {height: auto; padding: 1mm; border: 1px solid #666; border-bottom: 0; font-size: 11px;}";
      style += ".sticker-flash .sorting {display: flex; height: 35px; border: 1px solid #666; border-bottom: 0;} .sticker-flash .sorting .b1 {width: 60%; font-size: 14px; text-align: center;}";
      style += ".sticker-flash .sorting .b2 {width: 40%; color: white; font-size: 25px; font-weight: bold; text-align: center; background: black;} .sticker-flash .destination .first-line {font-size: 13px;}";
      style += ".sticker-flash .other {display: flex; height: 50px; border: 1px solid #666;} .sticker-flash .other .item-list { width: 50%; font-size: 9px; border-right: 1px solid #666;} .sticker-flash .other .cod { width: 50%; margin: auto; text-align: center; font-size: 14px; font-weight: bold;}";

      printJs({
        printable: "sticker_8x8",
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
  size: 3.2in;
  margin: 0px !important;
  padding: 0px !important;
} */
</style>

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

/* @media print {
  .preview {
    display: none !important;
  }
  .sticker-flash {
    width: 100% !important;
    height: 100% !important;
    margin: 0px !important;
    border: 0px !important;
  }

  .sticker-flash .logo img {
    width: 30% !important;
    height: 20px !important;
  }
} */
</style>
