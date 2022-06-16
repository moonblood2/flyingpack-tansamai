<template>
  <div class="big-box">
    <div class="text-box">
      <h3 style="display: inline">รายงาน Fulfillment</h3>
    </div>
    <div class="box">
      <div class="tool-box-1">
        <b-row style="margin-bottom: 10px">
          <b-col cols="3">
            <b-form-input v-model="form.startDate" type="date"></b-form-input>
            <b-form-timepicker
              v-model="form.startTime"
              :hour12="false"
              no-close-button
            >
            </b-form-timepicker>
          </b-col>
          <b-col cols="3">
            <b-form-input v-model="form.endDate" type="date"></b-form-input>
            <b-form-timepicker
              v-model="form.endTime"
              :hour12="false"
              no-close-button
            >
            </b-form-timepicker>
          </b-col>
          <b-col cols="3">
            <b-form-input
              v-model="form.keyWord"
              cols="3"
              placeholder="ref..., track..., ชื่อผู้ส่ง"
              type="text"
            ></b-form-input>
            <b-button @click="onClickExportSerialNumber" variant="success"
              >serial number</b-button
            >
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
              <!-- <b-form-checkbox-group
                v-model="form.productsIds"
                :options="productOptions"
                disabled-field="notEnabled"
                style="padding-left: 10px"
                text-field="name"
                value-field="item"
              ></b-form-checkbox-group> -->

              <!-- Tansamai ADD -->
              <b-form-group v-slot="{ ariaDescribedby }">
                <b-form-checkbox-group
                  id="filter_prod_group"
                  v-model="productsIdSelected"
                  :options="productOptions"
                  :aria-describedby="ariaDescribedby"
                  name="filter_prod"
                ></b-form-checkbox-group>
              </b-form-group>
            </b-dropdown>
          </b-col>
          <b-col cols="3">
            <b-form-select
              v-model="form.fulfillmentOfRow"
              :options="fulfillmentOfRowOptions"
              @change="onChangeRowTotal"
            ></b-form-select>
          </b-col>
          <b-col
            cols="3"
            style="
              display: flex;
              align-items: center;
              justify-content: flex-end;
            "
          >
            <b>ทั้งหมด {{ anOrder.totalItem }} รายการ</b>
          </b-col>
        </b-row>
        <hr />
        <b-row>
          <b-col cols="4" style="display: flex; justify-content: flex-start">
            <b-button-group>
              <b-button variant="info" @click="onClickSelectAll">All</b-button>
              <b-button @click="onClickClearSelected">Clear</b-button>
            </b-button-group>
          </b-col>
          <b-col cols="4"></b-col>
          <b-col cols="4" style="display: flex; justify-content: flex-end">
            <b-button-group>
              <b-overlay
                :show="loading.createOrder"
                opacity="0.4"
                spinner-small
                spinner-variant="primary"
              >
                <b-button
                  :disabled="disableCreateOrderBtn"
                  variant="info"
                  @click="onClickCreateOrder"
                >
                  ทำรายการ
                </b-button>
              </b-overlay>
              <b-dropdown text="ใบปะหน้า">
                <b-dropdown-item
                  @click="onClickLabel('sticker-4x6', selected.selectedItems)"
                  >Sticker4x6</b-dropdown-item
                >
                <b-dropdown-item
                  @click="onClickLabel('sticker-8x8', selected.selectedItems)"
                  >Sticker8x8</b-dropdown-item
                >
                <b-dropdown-item
                  @click="
                    onClickLabel('sticker-100x75', selected.selectedItems)
                  "
                  >Sticker100x75</b-dropdown-item
                >
              </b-dropdown>
            </b-button-group>
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
          selectable
          style="max-height: 100%"
        >
          <template #table-colgroup="scope">
            <col
              v-for="field in scope.fields"
              :key="field.key"
              :style="{ ...field.style }"
            />
          </template>
          <template #cell(index)="data">
            {{ perPage * (anOrder.currentPage - 1) + (data.index + 1) }}
          </template>
          <template #cell(createdAt)="data">
            {{ data.item.createdAt.split(".")[0] }}
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
          <template #cell(destinationBrief)="data">
            {{ `${data.item.destination.name}` }}
          </template>
          <template #cell(trackingCode)="data">
            {{ `${data.item.trackingCode}` }}
          </template>
          <template #cell(itemsBrief)="data">
            <p
              v-for="(v, i) in data.item.items"
              :key="i"
              style="font-size: 16px"
            >
              {{ `${v["productCode"]}: ${v["quantity"]}` }}
            </p>
          </template>
          <template #cell(handle)="data">
            <b-button-group>
              <!-- disable if status is packed or cancel or  -->
              <b-button
                :disabled="
                  data.item.fulfillmentStatus === 1 ||
                  data.item.fulfillmentStatus === 3
                "
                variant="primary"
                @click="onClickDo(data.index)"
              >
                ทำ
              </b-button>
              <!-- disable if status is not packed or cancel -->
              <b-button variant="warning" @click="onClickEdit(data.index)">
                แก้
              </b-button>
              <!-- disable if status is cancel -->
              <b-button
                :disabled="data.item.fulfillmentStatus === 3"
                variant="danger"
                @click="onClickCancel(data.index)"
              >
                ยกเลิก
              </b-button>
            </b-button-group>
          </template>
        </b-table>
      </div>
      <div style="display: flex">
        <b-pagination
          v-model="anOrder.currentPage"
          :per-page="perPage"
          :total-rows="anOrder.totalItem"
          class="mx-auto"
          @change="onChangePage"
        ></b-pagination>
      </div>
    </div>
    <ParcelForm
      ref="parcel-form"
      :busy="loading.createOrder"
      :disabled-cancel-button="loading.createOrder"
      :loading-ok-button="loading.createOrder"
      :on-cancel="onCancelParcelForm"
      :on-submit="onSubmitParcelForm"
      :parcel-form="parcelForm"
      :show-origin="false"
      :disable-validation="true"
    />
    <b-modal
      id="create-order-result-modal"
      :no-close-on-backdrop="true"
      title="ผลการทำรายการ"
    >
      <h5>สำเร็จ ({{ createOrderModal.successOrders.length }})</h5>
      <div>
        <p v-for="(e, i) in createOrderModal.successOrders" :key="i">
          {{ e.order.referenceNo }}
        </p>
        <b-dropdown
          class="w-100 p-3"
          text="ใบปะหน้า"
          variant="info"
          :disabled="createOrderModal.successOrders.length === 0"
        >
          <b-dropdown-item
            @click="
              onClickLabel(
                'sticker-8x8',
                createOrderModal.successOrders.map((x) => x.order)
              )
            "
            >Sticker8x8</b-dropdown-item
          >
          <b-dropdown-item
            @click="
              onClickLabel(
                'sticker-100x75',
                createOrderModal.successOrders.map((x) => x.order)
              )
            "
            >Sticker100x75</b-dropdown-item
          >
        </b-dropdown>
      </div>
      <h5>ไม่สำเร็จ ({{ createOrderModal.failOrders.length }})</h5>
      <div>
        <p v-for="(e, i) in createOrderModal.failOrders" :key="i">
          {{ e.order.referenceNo }} ({{ e.result.message }})
        </p>
      </div>
    </b-modal>
    <b-modal
      v-if="selectedCancelIndex !== -1"
      id="cancel-confirm-modal"
      :busy="loading.cancelOrder"
      :title="
        selectedCancelIndex !== -1
          ? `ต้องการยกเลิก ${items[selectedCancelIndex].referenceNo} ?`
          : 'เลือกก่อน'
      "
    >
      <template #modal-footer="{ cancel }">
        <b-button :disabled="loading.createOrder" @click="cancel"
          >ยกเลิก</b-button
        >
        <b-overlay
          :show="loading.cancelOrder"
          class="d-inline-block"
          opacity="0.4"
          spinner-small
          spinner-variant="primary"
        >
          <b-button variant="info" @click="onClickConfirmCancelOrder"
            >ตกลง</b-button
          >
        </b-overlay>
      </template>
    </b-modal>
    <b-modal
      v-if="selectedEditIndex !== -1"
      id="edit-confirm-modal"
      :busy="loading.editOrder"
      :title="
        selectedEditIndex !== -1
          ? `ต้องการแก้ไข ${items[selectedEditIndex].referenceNo} ?`
          : 'เลือกก่อน'
      "
    >
      <div>
        <div style="display: flex; justify-content: space-between">
          <h6>Tracking Code:</h6>
          <b-button variant="success" @click="onClickEdit_AddTrackingCode()"
            >+</b-button
          >
        </div>
        <div
          v-for="(v, i) in editForm.trackingCode"
          :key="i"
          style="display: flex"
        >
          <b-input
            placeholder="trackingCode"
            v-model="editForm.trackingCode[i]"
            type="text"
          ></b-input>
          <b-button variant="danger" @click="onClickEdit_RemoveTrackingCode(i)"
            >-</b-button
          >
        </div>
      </div>
      <div>
        <p>สถานะ Fulfillment</p>
        <b-form-select
          v-model="editForm.fulfillmentStatus"
          :options="fulfillmentStatusOptions"
        ></b-form-select>
      </div>
      <template #modal-footer="{ cancel }">
        <b-button :disabled="loading.editOrder" @click="cancel"
          >ยกเลิก</b-button
        >
        <b-overlay
          :show="loading.editOrder"
          class="d-inline-block"
          opacity="0.4"
          spinner-small
          spinner-variant="primary"
        >
          <b-button variant="info" @click="onClickConfirmEdit">ตกลง</b-button>
        </b-overlay>
      </template>
    </b-modal>
  </div>
</template>

<script>
import Vue from "vue";
import ParcelForm from "@/components/ParcelForm";

import XLSX from "xlsx";

import {
  anFulfillmentStatus,
  AnParcel,
  ContactInfo,
  parseAnParcel,
  Sender,
  SpOrderParcelShippopFlash,
  anFulfillmentStatusToString,
  AnCourier,
} from "@/entities";

import { LABEL_ADD_PARCEL } from "@/store/actions/label";
import { currentDate } from "@/utils/date";
import {
  cancelAnOrder,
  getAnProducts,
  getOrderFulfillment,
  updateOrderFulfillment,
  createOrderByWebhook,
} from "@/api/agent-network";

import env from "../constants/env";

import "@/styles/common.css";

export default {
  name: "ReportFulfillmentTable",
  components: {
    ParcelForm,
  },
  async created() {
    await this.onClickGetOrder();

    //Get products
    try {
      const getProductsRes = await getAnProducts();
      this.products = getProductsRes.data.data;
    } catch (error) {
      console.log(error);
    }
  },
  data() {
    return {
      //Loading status
      loading: {
        get: false,
        createOrder: false,
        editOrder: false,
        cancelOrder: false,
      },

      //Default address
      defaultSenderForm: new Sender({
        senderType: 2,
        nationalIdNumber: env.JNA_ID,
        birthDate: env.JNA_BIRTH_DATE,
        name: env.JNA_NAME,
        phoneNumber: env.JNA_PHONE_NUMBER,
        address: env.JNA_ADDRESS,
        district: env.JNA_DISTRICT,
        state: env.JNA_STATE,
        province: env.JNA_PROVINCE,
        postcode: env.JNA_POSTCODE,
      }),
      defaultOriginAddress: new ContactInfo({
        name: env.JNA_NAME,
        phoneNumber: env.JNA_PHONE_NUMBER,
        address: env.JNA_ADDRESS,
        district: env.JNA_DISTRICT,
        state: env.JNA_STATE,
        province: env.JNA_PROVINCE,
        postcode: env.JNA_POSTCODE,
      }),

      //Form options
      fulfillmentOfRowOptions: [
        { value: -1, text: "จำนวนข้อมูลที่แสดง" },
        { value: 50, text: "50" },
        { value: 100, text: "100" },
        { value: 500, text: "500" },
        { value: 1000, text: "1000" },
        { value: 3000, text: "3000" },
      ],
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
      productOptions: [
        { value: 1, text: "SERUM:1" },
        { value: 2, text: "SERUM:2" },
        { value: 3, text: "SERUM:3" },
        { value: 4, text: "SERUM:4" },
        { value: 5, text: "SERUM:5" },
      ],

      //Form data
      form: {
        startDate: currentDate(),
        endDate: currentDate(),
        startTime: "00:00:00",
        endTime: "23:59:59",
        keyWord: "",
        fulfillmentStatus: -1,
        courierCode: 0,
        fulfillmentOfRow: -1,
        // productsIds: [],
      },

      // Tansamai ADD
      productsIdSelected: [],
      productIds: [],

      //Table
      fields: [
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
          key: "destinationBrief",
          label: "ผู้รับ",
          sortable: false,
          style: { width: "200px" },
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
          style: { width: "100px" },
        },
        {
          key: "itemsBrief",
          label: "สินค้า",
          sortable: false,
          style: { width: "200px", ["font-size"]: "16px" },
        },
        {
          key: "codAmount",
          label: "COD",
          sortable: false,
          style: { width: "100px" },
        },
        {
          key: "handle",
          label: "จัดการ",
          sortable: false,
          style: { width: "50px" },
        },
      ],

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

      //Product data
      products: [],

      //Per page use to request order.
      perPage: 100,

      memSelectedItems: {},
      createOrderModal: {
        completed: false,
        successOrders: [],
        failOrders: [],
        message: "",
      },
      disableCreateOrderBtn: false,

      //ParcelForm for created one by one order.
      selectedIndex: 0,
      parcelForm: new AnParcel({}),

      editForm: {
        trackingCode: "",
        fulfillmentStatus: -1,
      },

      selectedEditIndex: -1,
      selectedCancelIndex: -1,
    };
  },
  computed: {
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

          // Tansamai ADD
          if (this.productsIdSelected.length > 0) {
            for (let j = 0; j < this.productsIdSelected.length; j++) {
              if (e.items[0].quantity == this.productsIdSelected[j]) {
                let anParcel = parseAnParcel(e);
                anParcel.origin = this.defaultOriginAddress;
                items.push(anParcel);
              }
            }
          } else {
            let anParcel = parseAnParcel(e);
            anParcel.origin = this.defaultOriginAddress;
            items.push(anParcel);
            //console.log("anParcel : ", anParcel);
          }
        }
        return [...items];
      },
      set(value) {
        return value;
      },
    },
    selected() {
      let anParcels = [];
      let selectedIndex = [];
      let selectedItems = [];

      let createdOrder = null;
      for (let i = 0; i < this.perPage; i++) {
        if (this.$refs["report-fulfillment-table"].isRowSelected(i)) {
          anParcels.push(this.items[i]);
          selectedIndex.push(i);
          selectedItems.push(this.items[i]);
          if (this.items[i].fulfillmentStatus === anFulfillmentStatus.PACKED) {
            createdOrder = this.items[i];
          }
        }
      }
      return { anParcels, selectedIndex, selectedItems, createdOrder };
    },
    // productOptions: function() {
    //   let productOptions = [];
    //   if (this.products) {
    //     for (const product of this.products) {
    //       productOptions.push({
    //         item: product.anProductId,
    //         name: product.productCode,
    //       });
    //     }
    //   }
    //   return [...productOptions];
    // },
  },
  methods: {
    onChangeRowTotal() {
      this.perPage = this.form.fulfillmentOfRow;
    },
    async onClickGetOrder() {
      this.loading.get = true;
      try {
        const res = await getOrderFulfillment(
          `${this.form.startDate} ${this.form.startTime}`,
          `${this.form.endDate} ${this.form.endTime}`,
          this.anOrder.currentPage,
          this.perPage,
          this.form.keyWord,
          this.form.fulfillmentStatus,
          this.form.courierCode,
          // this.form.productsIds
          this.productsIds
        );
        if (res.data.data) {
          this.anOrder = res.data.data;
        }
        this.loading.get = false;
      } catch (error) {
        this.loading.get = false;
        console.log(error);
      }
    },
    onChangePage(page) {
      this.anOrder.currentPage = page;
      this.onClickGetOrder();
    },

    onClickSelectAll() {
      this.$refs["report-fulfillment-table"].selectAllRows();
    },
    onClickClearSelected() {
      this.$refs["report-fulfillment-table"].clearSelected();
    },

    //onClickCreateOrder select order from selected item in table then request for created order API.
    async onClickCreateOrder() {
      const { selectedIndex, anParcels, selectedItems, createdOrder } =
        this.selected;
      if (selectedItems.length <= 0) {
        return;
      }
      if (createdOrder !== null) {
        await this.$bvModal.msgBoxOk(
          `${createdOrder.referenceNo} ได้ทำรายการไปแล้ว`,
          {
            title: "ผิดพลาด",
            size: "sm",
            buttonSize: "md",
            okVariant: "danger",
            centered: true,
          }
        );
        return;
      }

      this.createOrderModal.successOrders = [];
      this.createOrderModal.failOrders = [];

      this.loading.createOrder = true;

      try {
        const res = await createOrderByWebhook(
          this.defaultOriginAddress,
          anParcels
        );
        const { data } = res;
        console.log("data: ", data);
        //Check status
        if (data.code === 1) {
          this.hasCreatedOrder = true;
          //Populate tracking code, and Flash-sorting-code to order parcel
          if (data.data && data.data.orders) {
            const { orders, results } = data.data;
            for (let i = 0; i < orders.length; i++) {
              if (orders[i].fulfillmentStatus === anFulfillmentStatus.PACKED) {
                // orders[i].trackingCode = orders[i].trackingCode;
                this.createOrderModal.successOrders.push({
                  order: orders[i],
                  result: results[i],
                });

                this.items[selectedIndex[i]].trackingCode =
                  orders[i].trackingCode;
                this.items[selectedIndex[i]].trackingCode =
                  orders[i].trackingCode;
                this.items[selectedIndex[i]].spOrderParcelId =
                  orders[i].spOrderParcelId;
                this.items[selectedIndex[i]].fulfillmentStatus =
                  orders[i].fulfillmentStatus;
                this.items[selectedIndex[i]].fulfillmentStatusString =
                  orders[i].fulfillmentStatusString;
                if (orders[i].spOrderParcelShippopFlash) {
                  const { spOrderParcelShippopFlash } = orders[i];
                  this.items[selectedIndex[i]].spOrderParcelShippopFlash =
                    new SpOrderParcelShippopFlash({
                      dstCode: spOrderParcelShippopFlash.dstCode,
                      sortCode: spOrderParcelShippopFlash.sortCode,
                      sortingLineCode:
                        spOrderParcelShippopFlash.sortingLineCode,
                    });
                }
              } else {
                this.createOrderModal.failOrders.push({
                  order: orders[i],
                  result: results[i],
                });
              }
            }
          }
          Vue.set(this, "items", [...this.items]);
        }
      } catch (error) {
        console.log(error);
      }

      this.loading.createOrder = false;
      this.$bvModal.show("create-order-result-modal");
    },
    async onClickLabel(size, parcels) {
      if (parcels.length <= 0) {
        return;
      }
      console.log(parcels);
      const key = Math.random().toString(36).substring(2, 7);
      this.$store.commit(LABEL_ADD_PARCEL, { key: key, parcels: parcels });

      let routeData = this.$router.resolve({
        name: "Label",
        params: { key: key },
      });
      const url = `${routeData.href}?size=${size}`;
      console.log(url);
      window.open(url, "_blank");
    },

    //onClickDo select item in table and link selected item with ParcelForm.
    //When ParcelForm edit order details, order details in table and order details in ParcelForm
    //will chang together.
    onClickDo(index) {
      this.onClickClearSelected();
      this.$refs["report-fulfillment-table"].selectRow(index);
      this.selectedIndex = index;
      this.parcelForm = this.items[index];
      this.$refs["parcel-form"].show();
    },
    async onClickEdit(index) {
      this.$bvModal.show("edit-confirm-modal");
      this.selectedEditIndex = index;
      this.editForm.trackingCode = [...this.items[index].trackingCode];
      this.editForm.fulfillmentStatus = this.items[index].fulfillmentStatus;
    },
    onClickEdit_AddTrackingCode() {
      this.editForm.trackingCode.push("");
    },
    onClickEdit_RemoveTrackingCode(i) {
      this.editForm.trackingCode.splice(i, 1);
    },
    async onClickConfirmEdit() {
      let index = this.selectedEditIndex;
      this.loading.editOrder = true;
      let ok = false;

      try {
        await updateOrderFulfillment([
          {
            anOrderId: this.items[index]["anOrderId"],
            trackingCode: this.editForm.trackingCode,
            fulfillmentStatus: this.editForm.fulfillmentStatus,
          },
        ]);
        ok = true;
      } catch (error) {
        console.log(error);
      }

      if (ok) {
        // this.items[index]['trackingCode'] = this.editForm.trackingCode[0];
        this.items[index]["trackingCode"] = this.editForm.trackingCode;
        this.items[index]["fulfillmentStatus"] =
          this.editForm.fulfillmentStatus;
        this.items[index]["fulfillmentStatusString"] =
          anFulfillmentStatusToString[this.editForm.fulfillmentStatus];
      }
      this.$bvModal.hide("edit-confirm-modal");
      this.loading.editOrder = false;
    },

    onClickCancel(index) {
      this.$bvModal.show("cancel-confirm-modal");
      this.selectedCancelIndex = index;
    },
    async onClickConfirmCancelOrder() {
      this.loading.cancelOrder = true;
      try {
        await cancelAnOrder(this.items[this.selectedCancelIndex].anOrderId);
        this.items[this.selectedCancelIndex].fulfillmentStatus = 3;
        this.items[this.selectedCancelIndex].fulfillmentStatusString = "ยกเลิก";
      } catch (error) {
        console.log(error);
      }
      this.loading.cancelOrder = false;
      this.$bvModal.hide("cancel-confirm-modal");
    },
    //onSubmitParcelForm call onClickCreateOrder.
    async onSubmitParcelForm() {
      await this.onClickCreateOrder();
      this.$refs["parcel-form"].hide();
    },
    //onCancelParcelForm reset selected item
    onCancelParcelForm() {
      this.onClickClearSelected();
      this.selectedIndex = 0;
    },

    onClickExportSerialNumber() {
      let headerLabel = ["referenceNo", "serialNumber"];
      let data = [];
      for (const o of this.anOrder.orders) {
        for (const i of o.items) {
          for (const s of i.serialNumbers) {
            data.push([o.referenceNo, s.serialNumber]);
          }
        }
      }
      let wsData = [headerLabel, ...data];

      /**Create new workbook.*/
      const wb = XLSX.utils.book_new();
      /**Create new worksheet.*/
      const ws = XLSX.utils.aoa_to_sheet(wsData);
      const wsName = "Sheet1";
      XLSX.utils.book_append_sheet(wb, ws, wsName);
      /**Write file and save.*/
      const fileName = `serial-number[${this.form.startDate}-${
        this.form.endDate
      }][${Date.now()}].xlsx`;
      XLSX.writeFile(wb, fileName);
    },
  },
};
</script>

<style scoped>
.form-group {
  margin: 0px;
}
#filter_prod_group {
  padding: 0.5rem 1rem;
}
</style>
