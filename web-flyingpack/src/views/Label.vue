<template>
  <div v-if="size === 'sticker-4x6'" id="monochrome">
    <Sticker4x6
      v-for="(parcel, index) in parcels"
      :id="`label_${index}`"
      :key="`sticker-4x6_${index}`"
      :parcel="parcel"
    />
  </div>
  <div v-else-if="size === 'sticker-8x8'">
    <Sticker8x8
      v-for="(parcel, index) in parcels"
      :id="`label_${index}`"
      :key="`sticker-8x8_${index}`"
      :parcel="parcel"
    />
  </div>
  <div v-else-if="size === 'sticker-100x75'">
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
    parcels: function() {
      const { key } = this.$route.params;
      return this.$store.state.label.parcel[key];
    },
  },
};
</script>
