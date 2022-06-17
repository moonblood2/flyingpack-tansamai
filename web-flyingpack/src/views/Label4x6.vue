<template>
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
</template>

<script>
import Sticker4x6 from "@/components/labels/Sticker4x6";

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
      window.print();
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
  size: 4in 6in !important;
  margin: 0mm !important;
  padding: 0mm !important;
}

@media print {
  .preview {
    display: none !important;
  }
  .sticker {
    width: 100%;
    height: 100%;
    margin: 0px !important;
    border: 0px !important;
  }
}
</style>
