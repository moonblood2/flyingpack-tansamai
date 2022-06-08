<template>
  <div class="box">
    <h6>รายงานประจำวัน</h6>
    <div class="tb">
      <b-table
          id="fulfillment-pivot-table"
          ref="fulfillment-pivot-table"
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
        <template #cell(labels)="data">
          <span style="display: flex; justify-content: flex-start">
          {{ data.item.labels }}
          </span>
        </template>
        <template #cell(count)="data">
          <span style="display: flex; justify-content: flex-end">
          {{ data.item.count }}
          </span>
        </template>
        <template #foot(index)>
          -
        </template>
        <template #foot(labels)>
          <b>รวม</b>
        </template>
        <template #foot(count)>
          <b>{{ pivotTotal.total }}</b>
        </template>
        <template #foot(packedCount)>
          <b>{{ pivotTotal.packedTotal }}</b>
        </template>
        <template #foot(notPackedCount)>
          <b>{{ pivotTotal.notPackedTotal }}</b>
        </template>
        <template #foot(cancelCount)>
          <b>{{ pivotTotal.cancelTotal }}</b>
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
        {key: "labels", label: 'สินค้า', style: {width: '150px', 'text-align': 'start'}},
        {key: "count", label: 'จำนวน', style: {width: '50px'}},
        {key: "packedCount", label: 'ทำ', style: {width: '50px'}},
        {key: "notPackedCount", label: 'ยังไม่ทำ', style: {width: '50px'}},
        {key: "cancelCount", label: 'ยกเลิก', style: {width: '50px'}},
      ],
    }
  },
  computed: {
    pivotTotal: function () {
      let total = 0, packedTotal = 0, notPackedTotal = 0, cancelTotal = 0;
      if (this.items) {
        for (const i of this.items) {
          total += parseInt(i.count);
          notPackedTotal += parseInt(i.notPackedCount);
          packedTotal += parseInt(i.packedCount);
          cancelTotal += parseInt(i.cancelCount);
        }
      }
      return {
        total: total,
        packedTotal: packedTotal,
        notPackedTotal: notPackedTotal,
        cancelTotal: cancelTotal,
      }
    }
  }
}
</script>