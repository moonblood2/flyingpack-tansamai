<template>
  <div class="big-box">
    <div class="text-box">
      <h3>สรุปข้อมูลการขายพัสดุ</h3>
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
          <!--          <li>-->
          <!--            <logo for="search">ค้นหาข้อมูล</logo>-->
          <!--            <input id="search" class="input1" type="text"/>-->
          <!--          </li>-->
          <!--          <li>-->
          <!--            <p>ค้นหาขนส่ง-->
          <!--            <p/>-->
          <!--            <b-form-select v-model="selected" :options="options" class="select">ค้นหาข้อมูล</b-form-select>-->
          <!--          </li>-->
          <!--          <li>-->
          <!--            <p>Export File</p>-->
          <!--            <b-button variant="success">Excel</b-button>-->
          <!--            <b-button variant="primary">Print</b-button>-->
          <!--          </li>-->
        </b-row>
      </div>
      <!--      <div class="tool-box-2">-->
      <!--        <ul>-->
      <!--          <li>-->
      <!--            <p>ค้นหาข้อมูล-->
      <!--            <p/>-->
      <!--            <b-form-select v-model="selected2" :options="options2"></b-form-select>-->
      <!--          </li>-->
      <!--          <li>-->
      <!--            <p>ค้นหาข้อมูล-->
      <!--            <p/>-->
      <!--            <input class="input1" type="text">-->
      <!--          </li>-->
      <!--          <li>-->
      <!--            <b-button class="searchbox" size="lg" variant="primary">-->
      <!--              <b-icon aria-logo="Help" class="search" icon="search"></b-icon>-->
      <!--            </b-button>-->
      <!--          </li>-->
      <!--        </ul>-->
      <!--      </div>-->
      <!--      <p class="tracking">( หมายเลขใบเสร็จ , หมายเลข TRACKING )</p>-->
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
<!--          <template #cell(handle)="data">-->
<!--            <b-button-group>-->
<!--              <b-dropdown-->
<!--                  class="w-100 p-3"-->
<!--                  text="ใบปะหน้า"-->
<!--                  variant="info"-->
<!--              >-->
<!--                <b-dropdown-item @click="onClickLabel('sticker-8x8', data.item)">Sticker8x8</b-dropdown-item>-->
<!--                <b-dropdown-item @click="onClickLabel('sticker-4x6', data.item)">Sticker4x6</b-dropdown-item>-->
<!--              </b-dropdown>-->
<!--            </b-button-group>-->
<!--          </template>-->
        </b-table>
      </div>
    </div>
  </div>
</template>

<script>
import {getOrderParcel} from "@/api/shipping";
import {currentDate} from "@/utils/date";

import "@/styles/common.css";

export default {
  name: "ReportOrderParcelTable",
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
        {key: "desPostcode", label: 'รหัสไปรษณีย์ผู้รับ', sortable: false},
        {key: "desName", label: 'ชื่อผู้รับ', sortable: false},
        {key: "desPhoneNumber", label: 'หมายเลขโทรศัพท์ผู้รับ	', sortable: false},
        {key: "trackingCode", label: 'หมายเลข TRACKING', sortable: false},
        {key: "providerName", label: 'ชื่อขนส่ง', sortable: false},
        {key: "weight", label: 'น้ำหนัก', sortable: false},
        {key: "dimension", label: 'ขนาดพัสดุ (ซม.)', sortable: false},
        {key: "codAmount", label: 'ยอดเรียกเก็บ COD', sortable: false},
        {key: "price", label: 'ราคารวม', sortable: false},
        {key: "notice", label: 'หมายเหตุ', sortable: false},
        // {key: "handle", label: "จัดการ"},
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
        //Convert data to Parcel object
        // let parcel = new Parcel({
        //   name: "",
        //   trackingCode: e["order_parcel"]["tracking_code"],
        //   codAmount: e["order_parcel"]["cod_amount"],
        //   origin: new ContactInfo({
        //     name: e["origin"]["name"],
        //     phoneNumber: e["origin"]["phone_number"],
        //     address: e["origin"]["address"],
        //     district: e["origin"]["district"],
        //     state: e["origin"]["state"],
        //     province: e["origin"]["province"],
        //     postcode: e["origin"]["postcode"],
        //   }),
        //   destination: new ContactInfo({
        //     name: e["destination"]["name"],
        //     phoneNumber: e["destination"]["phone_number"],
        //     address: e["destination"]["address"],
        //     district: e["destination"]["district"],
        //     state: e["destination"]["state"],
        //     province: e["destination"]["province"],
        //     postcode: e["destination"]["postcode"],
        //   })
        // })

        let createdAt = String(e["order_parcel"]["created_at"]).split('T')
        itemsTmp.push({
          "index": i + 1,
          "createdDate": `${createdAt[0]} \n ${createdAt[1].split('.')[0]}`,
          "desPostcode": e["destination"]["postcode"],
          "desName": e["destination"]["name"],
          "desPhoneNumber": e["destination"]["phone_number"],
          "trackingCode": e["order_parcel"]["tracking_code"],
          "providerName": "-",
          "weight": e["order_parcel"]["weight"],
          "dimension": `${e["order_parcel"]["width"]} x ${e["order_parcel"]["length"]} x ${e["order_parcel"]["height"]}`,
          "codAmount": e["order_parcel"]["cod_amount"],
          "price": e["order_parcel"]["price"],
          "notice": "-",
          // "parcel": parcel,
        })
      }
      return [...itemsTmp]
    },
    totalPrice: function () {
      if (this.data === null) return 0.0
      let totalPrice = 0.0
      for (let i = 0; i < this.data.length; i++) {
        totalPrice += this.data[i]["order_parcel"]["price"]
      }
      return totalPrice
    },
  },
  methods: {
    async onClickGet() {
      this.loading.get = true
      try {
        const res = await getOrderParcel(this.form.startDate, this.form.endDate)
        this.data = res.data === null ? [] : res.data
        this.loading.get = false
      } catch (err) {
        this.loading.get = false
        throw new Error(err)
      }
    },
    // onClickLabel(size, item) {
    //   const key = Math.random().toString(36).substring(2, 7);
    //   this.$store.commit(LABEL_ADD_PARCEL, {key: key, parcels: Object(item.parcel)});
    //   let routeData = this.$router.resolve({name: 'Label', params: {key: key}});
    //   const url = `${routeData.href}?size=${size}`;
    //   window.open(url, '_blank');
    // },
  },
}
</script>