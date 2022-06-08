<template>
  <div class="robotic-packing-card" @click="$emit('click')">
    {{ data.packing_no }} <br/>
    {{ data.packing_status }}
  </div>
</template>

<script>
export default {
  name: "RoboticPackingCard",
  props: {
    data: Object,
  },
}
</script>

<style scoped>
.robotic-packing-card {
  height: 50px;
  border: 1px solid #555555;
  border-radius: 3px;
  margin-bottom: 2px;
  padding-left: 2px;
  cursor: pointer;
}
.robotic-packing-card:hover {
  box-shadow: 1px 1px 15px -2px #000000;
}
</style>