<template>
  <b-card
      :class="{'item-card': !disabled}"
      style="min-width: 150px; max-width: 150px; margin: 0 10px 10px 10px;"
      footer-tag="footer"
      header-tag="header"
      v-on:click="disabled ? null: $emit('click')"
  >
    <template #header>
      <slot name="header"></slot>
    </template>
    <slot></slot>
    <template #footer>
      <slot name="footer"></slot>
    </template>
  </b-card>
</template>

<script>
export default {
  props: {
    disabled: Boolean
  },
  name: "PosItemCard",
}
</script>

<style scoped>
.item-card:hover {
  cursor: pointer;
  transition: .4s;
  border: 1px solid rgba(0, 0, 0, .05);
}
</style>