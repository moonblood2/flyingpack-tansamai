<template>
  <div class="big-box">
    <div class="text-box">
      <h3 style="display: inline">รายงาน COD</h3>
    </div>
    <div class="box">
      <div class="tool-box-1">
        <b-row>
          <b-col cols="6">
            <b-form-select v-model="form.userId" :options="userOptions"></b-form-select>
          </b-col>
        </b-row>
        <b-row>
          <b-col cols="3">
            <b-form-input v-model="form.startDate" cols="3" type="date"></b-form-input>
          </b-col>
          <b-col cols="3">
            <b-form-input v-model="form.endDate" cols="3" type="date"></b-form-input>
          </b-col>
          <b-col cols="2">
            <b-form-select v-model="form.dateType" :options="dateTypeOptions"></b-form-select>
          </b-col>
          <b-col cols="3">
            <b-form-input v-model="form.keyWord" placeholder="ref..., track..., ชื่อผู้ส่ง"
                          type="text"></b-form-input>
          </b-col>
          <b-col cols="1">
            <b-overlay
                :show="loading.get"
                class="d-inline-block"
                opacity="0.4"
                spinner-small
                spinner-variant="primary"
            >
              <b-button class="search-box" size="lg" variant="primary" @click="onClickGetOrder">
                <b-icon aria-label="Help" class="search" icon="search"></b-icon>
              </b-button>
            </b-overlay>
          </b-col>
        </b-row>
        <hr>
        <b-row>
          <b-col cols="6">
            <b-button variant="success" @click="onClickExport">Export</b-button>
            <b-button variant="success" @click="onClickExportForTMB">Export สำหรับ TMB</b-button>
          </b-col>
          <b-col cols="3"></b-col>
          <b-col cols="3" style="display: flex; align-items: center; justify-content: flex-end;">
            ยอดเงินรวม {{ total | numberWithCommas }} บาท
          </b-col>
        </b-row>
        <hr>
        <b-row>
          <b-col cols="4" style="display: flex; justify-content: flex-start">
            <b-button-group>
              <b-button variant="info" @click="onClickSelectAll">All</b-button>
              <b-button @click="onClickClearSelected">Clear</b-button>
            </b-button-group>
          </b-col>
          <b-col cols="4">
          </b-col>
          <b-col cols="4" style="display: flex; justify-content: flex-end">
            <b-form-input v-model="form.jnaCodTransferredDate" type="date"></b-form-input>
            <b-overlay
                :show="loading.update"
                class="d-inline-block"
                opacity="0.4"
                spinner-small
                spinner-variant="primary"
            >
              <b-button variant="info" @click="onClickUpdate">โอนคืน (JNA)</b-button>
            </b-overlay>
          </b-col>
        </b-row>
        <hr>
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
            >
          </template>
          <template #cell(fulfillmentStatus)="data">
            <b-badge :variant="
            data.item.fulfillmentStatus === 1 ? 'success':
            data.item.fulfillmentStatus === 2 ? '':
            data.item.fulfillmentStatus === 3 ? 'danger': 'warning'"
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
import {parseAnParcel} from "@/entities";

import {currentDate} from "@/utils/date";
import {getUsers} from "@/api/shipping";
import {getOrderCod, updateOrderFulfillment} from "@/api/agent-network";

import "@/styles/common.css";
import {UserRoles} from "@/entities/User";

import XLSX from 'xlsx';

export default {
  name: "AccountingReportCodTable",
  async created() {
    //Get Users
    try {
      const res = await getUsers();
      let users = [];
      //Filter users get only AgentNetwork Member
      for (const user of res.data) {
        if (user.role === UserRoles.AGENT_NETWORK_MEMBER) {
          users.push({
            id: user['id'],
            email: user['email'],
            name: user['name'],
            role: user['role'],
            roleString: user['role_string']
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
      {key: "index", label: 'ลำดับ', sortable: false, style: {width: '75px'}},
      {key: "fulfillmentStatus", label: 'สถานะบรรจุ', sortable: false, style: {width: '125px'}},
      {key: "statusCompletedDate", label: 'วันที่ส่งสำเร็จ', sortable: false, style: {width: '125px'}},
      {key: "shippingStatus", label: 'สถานะจัดส่ง', sortable: false, style: {width: '125px'}},
      {key: "trackingCode", label: 'tracking code', sortable: false, style: {width: '150px'}},
      {key: "referenceNo", label: 'referenceNo', sortable: false, style: {width: '150px'}},
      {key: "desName", label: 'ผู้รับ', sortable: false, style: {width: '200px'}},
      {key: "codTransferredDate", label: 'วันที่โอนคืน', sortable: false, style: {width: '125px'}},
      {key: "jnaCodTransferredDate", label: 'วันที่โอนคืน (JNA)', sortable: false, style: {width: '125px'}},
      {key: "codStatus", label: 'สถานะ COD', sortable: false, style: {width: '125px'}},
      {key: "codAmount", label: 'ยอดเงิน COD', sortable: false, style: {width: '125px'}},
      {key: "codTransferredDate", label: 'วันที่โอนคืน', sortable: false, style: {width: '125px'}},
      {key: "bank", label: 'ธนาคาร', sortable: false, style: {width: '125px'}},
      {key: "accountName", label: 'ชื่อบัญชี', sortable: false, style: {width: '125px'}},
      {key: "accountNo", label: 'หมายเลขบัญชี', sortable: false, style: {width: '125px'}},
      {key: "email", label: 'อีเมลล์', sortable: false, style: {width: '125px'}},
    ];

    //Get orders
    if (this.form.userId) {
      await this.onClickGetOrder();
    }
  },
  data() {
    return {
      //Loading status
      loading: {
        get: false,
        update: false,
      },

      //Form data
      form: {
        startDate: currentDate(),
        endDate: currentDate(),
        dateType: 1,
        keyWord: "",
        userId: "",
        jnaCodTransferredDate: currentDate(),
      },

      //Form options
      dateTypeOptions: [
        {value: 1, text: 'วันที่ได้รับรายการ'},
        {value: 2, text: 'วันที่ส่งสำเร็จ'},
        {value: 3, text: 'วันที่โอนคืน'},
      ],

      //Table
      fields: [],

      //Order data
      orders: [],

      //Users data
      users: [],

      //Per page use to request order.
      perPage: 999999,
    }
  },
  computed: {
    //items contain data for each order, separate in two case, data to display and and data to consume in API.
    items: {
      get() {
        if (!this.orders) return [];
        //Map to order parcel.
        let items = [];
        for (let i = 0; i < this.orders.length; i++) {
          const e = this.orders[i];
          let anParcel = parseAnParcel(e);
          let item = {
            anParcel: anParcel,
            index: i + 1,
            createdAt: anParcel.createdAt.split(".")[0],
            codTransferredDate: anParcel.codTransferredDate.split(" ")[0],
            jnaCodTransferredDate: anParcel.jnaCodTransferredDate.split(" ")[0],
            statusCompletedDate: anParcel.statusCompletedDate.split(" ")[0],
            fulfillmentStatus: anParcel.fulfillmentStatus,
            fulfillmentStatusString: anParcel.fulfillmentStatusString,
            shippingStatus: anParcel.shippingStatus,
            codStatus: anParcel.codStatus,
            referenceNo: anParcel.referenceNo,
            desName: anParcel.destination.name,
            trackingCode: anParcel.trackingCode,
            codAmount: anParcel.codAmount,
            bank: anParcel.bankAccount.bank,
            accountName: anParcel.bankAccount.accountName,
            accountNo: anParcel.bankAccount.accountNo,
            email: anParcel.bankAccount.email,
          }
          items.push(item);
        }
        return [...items];
      },
      set(value) {
        return value;
      }
    },
    //userOptions is a list of all AgentNetwork Members.
    userOptions: function () {
      let userOptions = [];
      if (this.users) {
        for (const user of this.users) {
          userOptions.push({
            value: user.id,
            text: user.name,
          })
        }
      }
      return userOptions;
    },
    total: function () {
      let total = 0;
      for (const i of this.items) {
        total += i.anParcel.codAmount;
      }
      return total;
    },
  },
  methods: {
    onClickSelectAll() {
      this.$refs["report-fulfillment-table"].selectAllRows();
    },
    onClickClearSelected() {
      this.$refs["report-fulfillment-table"].clearSelected();
    },

    async onClickUpdate() {
      const ok = await this.$bvModal.msgBoxConfirm('ต้องการอัพเดตวันที่โอนคืนโดย JNA', {
        title: 'Please Confirm',
        size: 'sm',
        buttonSize: 'sm',
        okTitle: 'YES',
        cancelTitle: 'NO',
        footerClass: 'p-2',
        hideHeaderClose: false,
        centered: true
      })

      if (ok) {
        this.loading.update = true;
        let selectedIndexes = [];
        let totalItem = this.items.length;
        let updateItems = [];

        for (let i = 0; i < totalItem; i++) {
          if (this.$refs["report-fulfillment-table"].isRowSelected(i)) {
            selectedIndexes.push(i);
          }
        }
        for (const i of selectedIndexes) {
          updateItems.push({
            anOrderId: this.items[i].anParcel.anOrderId,
            jnaCodTransferredDate: this.form.jnaCodTransferredDate,
          });
        }

        try {
          const res = await updateOrderFulfillment(updateItems, this.form.userId);
          if (res.data.code === 1) {
            let j = 0;
            for (const i of selectedIndexes) {
              this.items[i].jnaCodTransferredDate = updateItems[j].jnaCodTransferredDate;
              j++;
            }
          }
        } catch (e) {
          console.log(e);
        }
        this.loading.update = false;
      }
    },

    //onClickGetOrder request API order fulfillment with price.
    async onClickGetOrder() {
      if (this.form.userId) {
        this.loading.get = true;
        try {
          const res = await getOrderCod(
              this.form.userId,
              `${this.form.startDate} 00:00:00`,
              `${this.form.endDate} 23:59:59`,
              this.form.dateType,
              this.form.keyWord,
          );
          if (res.data.data.orders) {
            this.orders = res.data.data.orders;
          }
          this.loading.get = false;
        } catch (error) {
          this.loading.get = false;
          console.log(error);
        }
      }
    },

    onClickExport() {
      /**Prepare worksheet header and data.*/
      let headerLabel = [];
      let data = [];
      for (const f of this.fields) {
        headerLabel.push(f.key)
      }
      for (const i of this.items) {
        data.push([
          i.index,
          i.fulfillmentStatusString,
          i.statusCompletedDate,
          i.shippingStatus,
          i.trackingCode,
          i.referenceNo,
          i.desName,
          i.codTransferredDate,
          i.jnaCodTransferredDate,
          i.codStatus,
          i.codAmount,
          i.codTransferredDate,
          i.bank,
          i.accountName,
          i.accountNo,
          i.email,
        ])
      }
      let wsData = [
        headerLabel,
        ...data,
      ];
      /**Create new workbook.*/
      const wb = XLSX.utils.book_new();
      /**Create new worksheet.*/
      const ws = XLSX.utils.aoa_to_sheet(wsData);
      const wsName = "Sheet1";
      XLSX.utils.book_append_sheet(wb, ws, wsName);
      /**Write file and save.*/
      const fileName = `report_cod_[${this.form.startDate}-${this.form.endDate}][${Date.now()}].xlsx`
      XLSX.writeFile(wb, fileName);
    },

    //onClickExportForTBB export table to .xlsx file.
    onClickExportForTMB() {
      /**Prepare worksheet header and data.*/
      let headerLabel = [
        "select Vendor",
        "Payee Name",
        "Amount",
        "Payment Ref.",
        "",
        "Payee A/C",
        "Payee Bank",
        "Advise Mode",
        "Fax Number",
        "Email",
        "Mobile No.",
        "Charge on",
        "tracking code",
        "ผู้รับสินค้า",
      ];
      let data = [];
      let items = [...this.items];
      for (let i = 0; i < items.length; i++) {
        const anParcel = items[i].anParcel;
        data.push([
          i + 1,
          anParcel.bankAccount.accountName,
          anParcel.codAmount,
          anParcel.referenceNo.replace("ALE", ""),
          "",
          anParcel.bankAccount.accountNo,
          `${anParcel.bankAccount.fiCode}: ${anParcel.bankAccount.fiName}`,
          "EMAIL",
          "",
          anParcel.bankAccount.email,
          "",
          "OUR",
          anParcel.trackingCode,
          anParcel.destination.name,
        ]);
      }
      let wsData = [
        headerLabel,
        ...data,
      ];
      /**Create new workbook.*/
      const wb = XLSX.utils.book_new();
      /**Create new worksheet.*/
      const ws = XLSX.utils.aoa_to_sheet(wsData);
      const wsName = "Sheet1";
      XLSX.utils.book_append_sheet(wb, ws, wsName);
      /**Write file and save.*/
      const fileName = `cod_[${this.form.startDate}-${this.form.endDate}][${Date.now()}].xlsx`
      XLSX.writeFile(wb, fileName);
    }
  },
}
</script>