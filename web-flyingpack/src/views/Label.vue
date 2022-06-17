<template>
  <div v-if="size === 'sticker-4x6'" id="monochrome">
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
  <div v-else-if="size === 'sticker-8x8'">
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
  <div v-else-if="size === 'sticker-100x75'">
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
</template>

<script>
import Sticker8x8 from "@/components/labels/Sticker8x8";
import Sticker4x6 from "@/components/labels/Sticker4x6";
import Sticker100x75 from "@/components/labels/Sticker100x75";

export default {
  name: "Label",
  components: {
    Sticker8x8,
    Sticker4x6,
    Sticker100x75,
  },
  data() {
    return {
      size: "sticker-8x8",
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

@media print {
  .preview {
    display: none !important;
  }
}
</style>
