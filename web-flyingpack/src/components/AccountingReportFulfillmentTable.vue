<template>
  <div class="big-box">
    <div class="text-box">
      <h3 style="display: inline">รายงานทางบัญชี Fulfillment</h3>
    </div>
    <div class="box">
      <div class="tool-box-1">
        <b-row>
          <b-col cols="6">
            <b-form-select
              v-model="form.userId"
              :options="userOptions"
            ></b-form-select>
          </b-col>
        </b-row>
        <b-row>
          <b-col cols="3">
            <b-form-input
              v-model="form.startDate"
              cols="3"
              type="date"
            ></b-form-input>
          </b-col>
          <b-col cols="3">
            <b-form-input
              v-model="form.endDate"
              cols="3"
              type="date"
            ></b-form-input>
          </b-col>
          <b-col cols="3">
            <b-form-input
              v-model="form.keyWord"
              cols="3"
              placeholder="referenceNo"
              type="text"
            ></b-form-input>
          </b-col>
          <b-col cols="2">
            <b-form-select
              v-model="form.fulfillmentStatus"
              :options="fulfillmentStatusOptions"
            ></b-form-select>
          </b-col>
          <b-col cols="1">
            <b-overlay
              :show="loading.get"
              class="d-inline-block"
              opacity="0.4"
              spinner-small
              spinner-variant="primary"
            >
              <b-button
                class="search-box"
                size="lg"
                variant="primary"
                @click="onClickGetOrder"
              >
                <b-icon aria-label="Help" class="search" icon="search"></b-icon>
              </b-button>
            </b-overlay>
          </b-col>
        </b-row>
        <b-row>
          <b-col cols="3">
            <b-form-select
              v-model="form.courierCode"
              :options="courierCodeOptions"
            ></b-form-select>
          </b-col>
          <b-col cols="3">
            <b-dropdown
              id="product-dropdown"
              dropright
              text="สินค้า"
              variant="outline-secondary"
            >
              <b-form-checkbox-group
                v-model="form.productsIds"
                :options="productOptions"
                disabled-field="notEnabled"
                style="padding-left: 10px"
                text-field="name"
                value-field="item"
              ></b-form-checkbox-group>
            </b-dropdown>
          </b-col>
          <b-col cols="3"></b-col>
          <b-col cols="3">
            <b-button variant="success" @click="onClickExport">Export</b-button>
          </b-col>
        </b-row>
        <hr />
        <b-row>
          <b-col v-if="!!anOrder.totalFulfillmentServiceCharge" cols="3">
            <h6>
              รวมค่าบริการ: {{ anOrder.totalFulfillmentServiceCharge }} บาท
            </h6>
          </b-col>
          <b-col cols="6"></b-col>
          <b-col
            cols="3"
            style="display: flex; align-items: center; justify-content: flex-end;"
          >
            <b>ทั้งหมด {{ anOrder.totalItem }} รายการ</b>
          </b-col>
        </b-row>
      </div>
      <div class="tb">
        <b-table
          id="report-fulfillment-table"
          ref="report-fulfillment-table"
          :busy="loading.get"
          :fields="fields"
          :fixed="true"
          :items="items"
          :show-empty="true"
          :sticky-header="true"
          empty-text="ไม่มีรายการให้แสดง"
          style="max-height: 100%"
        >
          <template #table-colgroup="scope">
            <col
              v-for="field in scope.fields"
              :key="field.key"
              :style="{ ...field.style }"
            />
          </template>
          <template #cell(status)="data">
            <b-badge
              :variant="
                data.item.fulfillmentStatus === 1
                  ? 'success'
                  : data.item.fulfillmentStatus === 2
                  ? ''
                  : data.item.fulfillmentStatus === 3
                  ? 'danger'
                  : 'warning'
              "
            >
              {{ data.item.fulfillmentStatusString }}
            </b-badge>
          </template>
        </b-table>
      </div>
    </div>
  </div>
</template>

<script>
import { parseAnParcel } from "@/entities";

import { currentDate } from "@/utils/date";
import { getUsers } from "@/api/shipping";
import {
  getAnProductsByUserId,
  getOrderFulfillmentPrice,
} from "@/api/agent-network";

import "@/styles/common.css";
import { UserRoles } from "@/entities/User";

import XLSX from "xlsx";
import { AnCourier } from "@/entities/AnCourier";

export default {
  name: "AccountingReportFulfillmentTable",
  async created() {
    //Get Users
    try {
      const res = await getUsers();
      let users = [];
      //Filter users get only AgentNetwork Member
      for (const user of res.data) {
        if (user.role === UserRoles.AGENT_NETWORK_MEMBER) {
          users.push({
            id: user["id"],
            email: user["email"],
            name: user["name"],
            role: user["role"],
            roleString: user["role_string"],
          });
        }
      }
      this.users = users;
      if (this.users[0]) {
        this.form.userId = users[0].id;
      }
    } catch (error) {
      console.log(error);
    }

    //Set field, Meta
    this.fields = [
      {
        key: "index",
        label: "ลำดับ",
        sortable: false,
        style: { width: "75px" },
      },
      {
        key: "createdAt",
        label: "ได้รับเมื่อ",
        sortable: false,
        style: { width: "125px" },
      },
      {
        key: "status",
        label: "สถานะ",
        sortable: false,
        style: { width: "75px" },
      },
      {
        key: "referenceNo",
        label: "referenceNo",
        sortable: false,
        style: { width: "150px" },
      },
      {
        key: "desName",
        label: "ผู้รับ",
        sortable: false,
        style: { width: "200px" },
      },
      {
        key: "desPhoneNumber",
        label: "เบอร์โทร",
        sortable: false,
        style: { width: "150px" },
      },
      {
        key: "desAddress",
        label: "ที่อยู่",
        sortable: false,
        style: { width: "150px" },
      },
      {
        key: "desSubdistrict",
        label: "ตำบล/แขวง",
        sortable: false,
        style: { width: "150px" },
      },
      {
        key: "desDistrict",
        label: "อำเภอ/เขต",
        sortable: false,
        style: { width: "150px" },
      },
      {
        key: "desProvince",
        label: "จังหวัด",
        sortable: false,
        style: { width: "150px" },
      },
      {
        key: "desPostcode",
        label: "รหัสไปรษณีย์",
        sortable: false,
        style: { width: "100px" },
      },
      {
        key: "courierName",
        label: "ขนส่ง",
        sortable: false,
        style: { width: "150px" },
      },
      {
        key: "trackingCode",
        label: "tracking code",
        sortable: false,
        style: { width: "150px" },
      },
      {
        key: "codAmount",
        label: "COD",
        sortable: false,
        style: { width: "100px" },
      },
      {
        key: "weight",
        label: "น้ำหนัก",
        sortable: false,
        style: { width: "100px" },
      },
      {
        key: "dimension",
        label: "ขนาด",
        sortable: false,
        style: { width: "100px" },
      },
      {
        key: "widthLengthHeight",
        label: "กว้าง x ยาว x สูง",
        sortable: false,
        style: { width: "150px" },
      },
    ];

    if (this.form.userId) {
      //Get products
      try {
        const getProductsRes = await getAnProductsByUserId(this.form.userId);
        this.products = getProductsRes.data.data;
        //Set field, Product
        for (const p of this.products) {
          this.fields.push({
            key: p["productCode"],
            label: p["productCode"],
            sortable: false,
            style: { width: "100px" },
          });
        }
      } catch (error) {
        console.log(error);
      }

      //Get orders
      await this.onClickGetOrder();
    }

    //Set field, Price
    this.fields = [
      ...this.fields,
      {
        key: "itemQuantitySum",
        label: "จำนวนสินค้ารวม",
        sortable: false,
        style: { width: "150px" },
      },
      {
        key: "parcelCost",
        label: "ราคาต้นทุน(ขนส่ง)",
        sortable: false,
        style: { width: "150px" },
      },
      {
        key: "parcelSellingPrice",
        label: "ราคาขาย(ขนส่ง)",
        sortable: false,
        style: { width: "100px" },
      },
      {
        key: "fulfillmentServiceCharge",
        label: "ค่าบริการ",
        sortable: false,
        style: { width: "100px" },
      },
    ];
  },
  data() {
    return {
      //Loading status
      loading: {
        get: false,
      },

      //Form options
      fulfillmentStatusOptions: [
        { value: -1, text: "สถานะ" },
        { value: 1, text: "จัดส่งแล้ว" },
        { value: 2, text: "ยังไม่จัดส่ง" },
        { value: 3, text: "ยกเลิก" },
      ],
      courierCodeOptions: [
        { value: 0, text: "ขนส่ง" },
        { value: AnCourier.FLASH.code, text: AnCourier.FLASH.text },
        { value: AnCourier.KERRY.code, text: AnCourier.KERRY.text },
        { value: AnCourier.SCG.code, text: AnCourier.SCG.text },
        { value: AnCourier.EMS.code, text: AnCourier.EMS.text },
        { value: AnCourier.EMS_WORLD.code, text: AnCourier.EMS_WORLD.text },
        { value: AnCourier.MESSENGER.code, text: AnCourier.MESSENGER.text },
      ],

      //Form data
      form: {
        startDate: currentDate(),
        endDate: currentDate(),
        keyWord: "",
        fulfillmentStatus: -1,
        courierCode: 0,
        productsIds: [],
        userId: "",
      },

      //Table
      fields: [],

      //Order data
      anOrder: {
        previousPage: null,
        currentPage: 1,
        nextPage: null,
        firstPage: 1,
        lastPage: 1,
        isFirstPage: true,
        isLastPage: false,
        totalPage: 1,
        totalItem: 0,
        orders: [],
      },

      //Users data
      users: [],

      //Product data
      products: [],

      //Per page use to request order.
      perPage: 999999,
    };
  },
  computed: {
    //items contain data for each order, separate in two case, data to display and and data to consume in API.
    items: {
      get() {
        if (
          !this.anOrder.orders ||
          !this.anOrder ||
          this.anOrder.totalItem === 0
        )
          return [];
        //Map to order parcel.
        let items = [];
        for (let i = 0; i < this.anOrder.orders.length; i++) {
          const e = this.anOrder.orders[i];
          let mapAnOrderItems = {}; //Key=ProductCode, Value=Quantity
          for (const p of this.products) {
            mapAnOrderItems[p.productCode] = 0;
          }
          for (const item of e.items) {
            mapAnOrderItems[item.productCode] += item.quantity;
          }

          let anParcel = parseAnParcel(e);
          let item = {
            anParcel: anParcel,
            index: i + 1,
            createdAt: anParcel.createdAt.split(".")[0],
            status: anParcel.fulfillmentStatusString,
            fulfillmentStatus: anParcel.fulfillmentStatus,
            fulfillmentStatusString: anParcel.fulfillmentStatusString,
            referenceNo: anParcel.referenceNo,
            desName: anParcel.destination.name,
            desPhoneNumber: anParcel.destination.phoneNumber,
            desAddress: anParcel.destination.address,
            desSubdistrict: anParcel.destination.district,
            desDistrict: anParcel.destination.state,
            desProvince: anParcel.destination.province,
            desPostcode: anParcel.destination.postcode,
            courierName: anParcel.courierName,
            trackingCode: anParcel.trackingCode,
            codAmount: anParcel.codAmount,
            weight: anParcel.weight,
            dimension: anParcel.width + anParcel.length + anParcel.height,
            widthLengthHeight: `${anParcel.width} x ${anParcel.length} x ${anParcel.height}`,
            ...mapAnOrderItems, //Populate mapAnOrderItems for viewing in Accounting.
            itemQuantitySum: anParcel.anOrderItemQuantitySum,
            parcelCost: `${
              anParcel.anOrderPrice ? anParcel.anOrderPrice.parcelCost : 0
            }`,
            parcelSellingPrice: `${
              anParcel.anOrderPrice
                ? anParcel.anOrderPrice.parcelSellingPrice
                : 0
            }`,
            fulfillmentServiceCharge: `${
              anParcel.anOrderPrice
                ? anParcel.anOrderPrice.fulfillmentServiceCharge
                : 0
            }`,
          };
          items.push(item);
        }
        return [...items];
      },
      set(value) {
        return value;
      },
    },
    //productOptions is a list of all products of each AgentNetwork Member.
    productOptions: function() {
      let productOptions = [];
      if (this.products) {
        for (const product of this.products) {
          productOptions.push({
            item: product.anProductId,
            name: product.productCode,
          });
        }
      }
      return [...productOptions];
    },
    //userOptions is a list of all AgentNetwork Members.
    userOptions: function() {
      let userOptions = [];
      if (this.users) {
        for (const user of this.users) {
          userOptions.push({
            value: user.id,
            text: user.name,
          });
        }
      }
      return userOptions;
    },
  },
  methods: {
    //onClickGetOrder request API order fulfillment with price.
    async onClickGetOrder() {
      if (this.form.userId) {
        this.loading.get = true;
        try {
          const res = await getOrderFulfillmentPrice(
            `${this.form.startDate} 00:00:00`,
            `${this.form.endDate} 23:59:59`,
            this.anOrder.currentPage,
            this.perPage,
            this.form.keyWord,
            this.form.fulfillmentStatus,
            this.form.courierCode,
            this.form.productsIds,
            this.form.userId
          );
          if (res.data.data) {
            this.anOrder = res.data.data;
            console.log(this.anOrder);
          }
          this.loading.get = false;
        } catch (error) {
          this.loading.get = false;
          console.log(error);
        }
      }
    },
    //onClickExport export table to .xlsx file.
    onClickExport() {
      /**Prepare worksheet header and data.*/
      let headerLabel = [];
      let headerKey = [];
      let data = [];
      for (const f of this.fields) {
        if (f.key === "index") {
          continue;
        }
        headerKey.push(f.key);
        headerLabel.push(f.label);
      }
      let items = [...this.items];
      for (const item of items) {
        let row = [];
        for (const key of headerKey) {
          row.push(item[key]);
        }
        data.push(row);
      }
      let wsData = [headerLabel, ...data];
      /**Create new workbook.*/
      const wb = XLSX.utils.book_new();
      /**Create new worksheet.*/
      const ws = XLSX.utils.aoa_to_sheet(wsData);
      const wsName = "Sheet1";
      XLSX.utils.book_append_sheet(wb, ws, wsName);
      /**Write file and save.*/
      const fileName = `fulfillment-report[${this.form.startDate}-${
        this.form.endDate
      }][${Date.now()}].xlsx`;
      XLSX.writeFile(wb, fileName);
    },
  },
};
</script>
