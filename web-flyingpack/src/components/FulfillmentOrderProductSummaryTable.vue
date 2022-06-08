<template>
  <div class="box">
    <h6>รายงานสินค้าประจำวัน</h6>
    <div class="tb">
      <b-table
          id="fulfillment-order-product-summary-table"
          ref="fulfillment-order-product-summary-table"
          :busy="busy"
          :fields="fields"
          :fixed="true"
          :items="items"
          :show-empty="true"
          :sticky-header="true"
          empty-text="ไม่มีรายการให้แสดง"
          foot-clone
          style="max-height: 100%"
      >
        <template #table-colgroup="scope">
          <col
              v-for="field in scope.fields"
              :key="field.key"
              :style="{ ...field.style }"
          >
        </template>
        <template #cell(index)="data">
          {{ data.index + 1 }}
        </template>
        <template #cell(productCode)="data">
          <span style="display: flex; justify-content: flex-start">
          {{ data.item.productCode }}
          </span>
        </template>
        <template #cell(sum)="data">
          <span style="display: flex; justify-content: flex-end">
          {{ data.item.sum }}
          </span>
        </template>
        <template #foot(index)>
          -
        </template>
        <template #foot(productCode)>
          <b>รวม</b>
        </template>
        <template #foot(sum)>
          <b>{{ total }}</b>
        </template>
      </b-table>
    </div>
  </div>
</template>

<script>
import "@/styles/common.css";

export default {
  name: "FulfillmentPivotTable",
  props: {
    items: Array,
    busy: Boolean,
  },
  data() {
    return {
      loading: {
        get: false,
      },
      fields: [
        {key: "index", label: 'ลำดับ', style: {width: '50px'}},
        {key: "productCode", label: 'สินค้า', style: {width: '100px', 'text-align': 'start'}},
        {key: "sum", label: 'จำนวน', style: {width: '75px'}},
      ],
    }
  },
  computed: {
    total: function () {
      if (this.items) {
        let total = 0;
        for (const i of this.items) {
          total += parseInt(i.sum);
        }
        return total;
      }
      return 0;
    }
  }
}
</script>