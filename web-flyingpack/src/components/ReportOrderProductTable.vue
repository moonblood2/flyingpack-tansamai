<template>
  <div class="big-box">
    <div class="text-box">
      <h3 style="display: inline">สรุปข้อมูลการขายสินค้า</h3><h5 style="display: inline"> (ไม่รวมการส่งพัสดุ)</h5>
      <h4>ยอดรวมทั้งหมด : {{ totalPrice }} บาท</h4>
    </div>
    <div class="box">
      <div class="tool-box-1">
        <b-row>
          <b-col><label>ค้นหาวันที่ขายพัสดุ (วันเริ่มต้น - วันสุดท้าย)</label></b-col>
        </b-row>
        <b-row>
          <b-col cols="3">
            <b-form-input v-model="form.startDate" cols="3" type="date"></b-form-input>
          </b-col>
          <b-col cols="3">
            <b-form-input v-model="form.endDate" cols="3" type="date"></b-form-input>
          </b-col>
          <b-col cols="1">
            <b-overlay
                :show="loading.get"
                class="d-inline-block"
                opacity="0.4"
                spinner-small
                spinner-variant="primary"
            >
              <b-button class="search-box" size="lg" variant="primary" @click="onClickGet">
                <b-icon aria-label="Help" class="search" icon="search"></b-icon>
              </b-button>
            </b-overlay>
          </b-col>
        </b-row>
        <!--          <li>-->
        <!--            <p>ค้นหาข้อมูล-->
        <!--            <p/>-->
        <!--            <input class="input1" type="text">-->
        <!--          </li>-->
        <!--          <li>-->
        <!--            <p>Export File</p>-->
        <!--            <b-button variant="success">Excel</b-button>-->
        <!--            <b-button variant="primary">Print</b-button>-->
        <!--          </li>-->
      </div>
      <div class="tb">
        <b-table
            :fields="fields"
            :items="items"
            :show-empty="true"
            :sticky-header="true"
            empty-text="ไม่มีรายการให้แสดง"
        >
          <template #head()="scope">
            <div class="text-nowrap">
              {{ scope.label }}
            </div>
          </template>
        </b-table>
      </div>
    </div>
  </div>
</template>

<script>
import {currentDate} from "@/utils/date";
import {getOrderProduct} from "@/api/shipping";

import "@/styles/common.css";

export default {
  name: "ReportOrderProductTable",
  created() {
    this.onClickGet()
  },
  data() {
    return {
      form: {
        startDate: currentDate(),
        endDate: currentDate(),
      },
      fields: [
        {key: "index", label: 'ลำดับ', sortable: false},
        {key: "createdDate", label: 'วันที่สร้างรายการ', sortable: false},
        {key: "details", label: 'รายละเอียด', sortable: false},
        {key: "price", label: 'ราคาขาย/ชิ้น (บาท)', sortable: false},
        {key: "quantity", label: 'จำนวนทั้งหมด', sortable: false},
        {key: "totalPrice", label: 'ราคาทั้งหมด (บาท)', sortable: false},
      ],
      data: [],
      loading: {
        get: false,
      },
    }
  },
  computed: {
    items: function () {
      if (this.data === null) return []
      let itemsTmp = []
      for (let i = 0; i < this.data.length; i++) {
        const e = this.data[i]
        itemsTmp.push({
          "index": i + 1,
          "createdDate": String(e["order_product"]["created_at"]).split('T')[0],
          "details": e["product"]["name"],
          "price": e["product"]["price"],
          "quantity": e["order_product"]["quantity"],
          "totalPrice": parseFloat(e["order_product"]["quantity"]) * parseFloat(e["product"]["price"]),
        })
      }
      return [...itemsTmp]
    },
    totalPrice: function () {
      if (this.data === null) return 0.0
      let totalPrice = 0.0
      for (let i = 0; i < this.data.length; i++) {
        totalPrice += parseFloat(this.data[i]["order_product"]["quantity"]) * parseFloat(this.data[i]["product"]["price"])
      }
      return totalPrice
    },
  },
  methods: {
    async onClickGet() {
      this.loading.get = true
      try {
        const res = await getOrderProduct(this.form.startDate, this.form.endDate)
        this.data = res.data === null ? [] : res.data
        this.loading.get = false
      } catch (err) {
        this.loading.get = false
        throw new Error(err)
      }
    },
  },
}
</script>