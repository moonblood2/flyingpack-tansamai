<template>
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
</template>

<script>
import Sticker100x75 from "@/components/labels/Sticker100x75";

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
      window.print();
    },
  },
};
</script>

<style>
.preview {
  width: 100%;
  display: inline-block;
  text-align: center;
  padding: 10px 0px;
}
#btn-preview{
  display: inline-block;
  font-weight: 400;
  font-size: 1rem;
  padding: .5rem 1rem;
  line-height: 1.5;
  border-radius: .3rem;
  color: #fff;
  background-color: #007bff;
  text-align: center;
  vertical-align: middle;
  box-shadow: unset;
  border: 0px;
  cursor: pointer;
}
#btn-preview:hover{
  background-color: #0069d9;
}


@page {
  size: 3.9in 3in;
  margin: 0mm !important;
  padding: 0mm !important;
}

@media print {
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
  .address1, .address2{
    font-size: 12px;
  }
  .product {
    font-size: 10px;
  }
  .cod {
    height: 90% !important;
  }
}
</style>
