<template>
  <div class="big-box">
    <div class="box">
      <b-row>
        <b-col cols="6">
          <b-form inline style="display: flex" @submit.prevent="onSubmitSearch">
            <b-form-input
                id="reference-no-input"
                v-model="referenceNo"
                placeholder="เลขอ้างอิง"
                required
                type="text"
            ></b-form-input>
            <b-overlay
                :show="loading.get"
                opacity="0.4"
                spinner-small
                spinner-variant="primary"
            >
              <b-button style="margin-left: 5px; color: white" type="submit" variant="primary">ค้นหา</b-button>
            </b-overlay>
          </b-form>
        </b-col>
        <b-col cols="6" style="display: flex; justify-content: space-between;">
          <b-button @click="onClickPrintSlip">พิมพ์สลิป</b-button>
          <div>
            robot:
            <b-overlay
                :show="loadingRobot.sendOrder"
                class="d-inline-block"
                opacity="0.4"
                spinner-small
                spinner-variant="primary"
            >
              <b-button @click="onClickRobotSendOrder">สั่ง</b-button>
            </b-overlay>
            <b-overlay
                :show="loadingRobot.takePhoto"
                class="d-inline-block"
                opacity="0.4"
                spinner-small
                spinner-variant="primary"
            >
              <b-button @click="onClickRobotTakePhoto">รูป</b-button>
            </b-overlay>
            <!-- <b-overlay
                :show="loadingRobot.serialNumber"
                class="d-inline-block"
                opacity="0.4"
                spinner-small
                spinner-variant="primary"
            >
              <b-button @click="onClickRobotSerialNumber">Serial Number</b-button>
            </b-overlay> -->
          </div>
        </b-col>
      </b-row>
    </div>

    <!-- Alert Box-->
    <b-row>
      <b-col>
        <b-alert
            :show="alertStatus"
            :variant="alertVariant"
        >
          {{ alertMessage }}
        </b-alert>
      </b-col>
    </b-row>

    <!-- Order -->
    <b-row v-if="hasOrder">
      <!-- Order Card -->
      <b-col cols="4">
        <div class="box">
          <h3>{{ order.referenceNo }}</h3>
          <b-overlay :show="!hasOrder" rounded="lg" opacity="0.6">
            <template #overlay>
              <div class="d-flex align-items-center">
              </div>
            </template>
            <p><b>ชื่อ</b>: {{ order.destination.name }}</p>
            <p><b>สถานะ</b>: {{ order.fulfillmentStatusString }}</p>
            <p><b>เบอร์โทร</b>: {{ order.destination.phoneNumber }}</p>
            <p><b>รหัสไปรษณีย์</b>: {{ order.destination.postcode }}</p>
            <p><b>ขนส่ง</b>: {{ order.courierName }}</p>
            <p><b>COD</b>: {{ order.codAmount }}</p>
            <p><b>เลข Tracking</b>: {{ order.trackingCode }}</p>
            <div style="background-color: rgb(235, 235, 224); padding: 5px">
              <div style="display: flex; justify-content: space-between">
                <b>Serial Number</b>
                <b-button @click="onClickEditSerialNumber" :disabled="!hasOrder">แก้ไข</b-button>
              </div>
              <p v-for="(o, i) in order.items" :key="i">
                {{ o.productCode }} ({{ o.quantity }}):
                <span
                    v-if="o.serialNumbers.filter(e => e.serialNumberStart === '').length === o.serialNumbers.length">ไม่มี</span>
                <span v-else>{{
                    o.serialNumbers.map(e => e.serialNumberStart + (e.serialNumberEnd ? `-${e.serialNumberEnd}` : '')).join(', ')
                  }}</span>
              </p>
            </div>
          </b-overlay>
        </div>
      </b-col>

      <!-- Robotic -->
      <b-col cols="5">
        <b-modal
            id="modal-image"
            title="ผลการถ่ายรูป"
            size="lg"
        >
          <b-img-lazy :src="imageBase64" width="500"></b-img-lazy>
        </b-modal>
      </b-col>
      <b-col cols="3">

      </b-col>
    </b-row>

    <!-- Serial Number Editor-->
    <b-modal
        id="modal-serial-number"
        title="Serial number"
        cancel-title="ยกเลิก"
        @ok.prevent="onSubmitSerialForm"
    >
      <div v-for="(o, i) in order.items" :key="i">
        <SerialNumberForm
            :title="`${o.productCode}: ${o.quantity}`"
            :index="i"
            :serial-numbers="o.serialNumbers"
            :on-change-serial-numbers="onChangeSerialNumbers"
        />
      </div>
      <template #modal-footer="{ok}">
        <b-overlay
            :show="loading.update"
            class="d-inline-block"
            opacity="0.4"
            spinner-small
            spinner-variant="primary"
        >
          <b-button variant="primary" @click="ok">ตกลง</b-button>
        </b-overlay>
      </template>
    </b-modal>
  </div>
</template>

<script>
import SerialNumberForm from "@/components/SerialNumberForm";
import {AnParcel, parseAnParcel} from "@/entities";

import {createSerialNumbers, getOrderByReferenceNo, getSlip, updateSerialNumbers} from "@/api/agent-network";
import {addOrder, getPackingData, takePhoto} from "@/api/robot";
import {PACKING_SLIP_ADD_ITEM} from "@/store/actions/packing-slip";
import {SLIP_SET} from "@/store/actions/slip";
import {BlobToDataURL} from "@/utils/image";
import {currentDate} from "@/utils/date";

import "@/styles/common.css";

export default {
  name: "OrderFulfillmentChecker",
  components: {
    SerialNumberForm
  },
  data() {
    return {
      hasSubmit: false,
      loading: {
        get: false,
        update: false,
      },
      loadingRobot: {
        sendOrder: false,
        takePhoto: false,
        serialNumber: false,
      },
      hasOrder: false,
      referenceNo: "",
      order: new AnParcel({}),

      //Image
      hasImage: false,
      imageBase64: "",

      //Serial
      serialNumber: "",

      fields: [],

      //Alert
      alertStatus: false,
      alertMessage: "",
      alertVariant: "info",

      //Robotic
      roboticPackingData: [],
    }
  },
  async created() {
    if (!this.$store.state.slip.hasSet) {
      try {
        const res = await getSlip();
        if (res.data && res.data.data && res.data.data.slip) {
          let slip = res.data.data.slip;
          this.$store.commit(SLIP_SET, slip)
        }
      } catch (e) {
        console.log(e)
      }
    }
  },
  methods: {
    setAlert(status = false, message = "", variant = "info") {
      this.alertStatus = status;
      this.alertMessage = message;
      this.alertVariant = variant;

      setTimeout(() => {
        this.alertStatus = false;
      }, 3000);
    },

    /**Search**/

    async onSubmitSearch() {
      this.hasSubmit = true;
      if (this.referenceNo) {
        this.loading.get = true;
        try {
          const res = await getOrderByReferenceNo(this.referenceNo);
          this.loading.get = false;
          if (res.data && res.data.data.order) {
            this.order = parseAnParcel(res.data.data.order);
            this.hasOrder = true;
            for (let i = 0; i < this.order.items.length; i++) {
              for (let j = 0; j < this.order.items[i].serialNumbers; j++) {
                this.order.items[i].serialNumbers[j].valid = true;
              }
            }
          } else {
            this.order = new AnParcel({});
            this.hasOrder = false;
          }
        } catch (error) {
          console.log(error);
          this.loading.get = false;
          this.hasOrder = false;
        }
      } else {
        new AnParcel({});
        this.hasOrder = false;
      }
    },

    /**Serial number**/

    onClickEditSerialNumber() {
      this.$bvModal.show("modal-serial-number");
    },

    onChangeSerialNumbers(i, serialNumbers) {
      this.order.items[i].serialNumbers = serialNumbers;
    },

    async onSubmitSerialForm(bvModalEvt) {
      await bvModalEvt.preventDefault();
      this.loading.update = true;

      try {
        //Categorize serialNumbers, which serial should create or update or delete.
        let serialsNumberForUpdate = [];
        let itemsForCreate = [];
        for (let i = 0; i < this.order.items.length; i++) {
          let item = this.order.items[i];
          let serialNumberForCreate = [];
          for (let j = 0; j < item.serialNumbers.length; j++) {
            let s = this.order.items[i].serialNumbers[j];
            if (s.anOrderProductSerialNumberId !== "") {
              serialsNumberForUpdate.push(s);
            } else if (s.serialNumberStart !== "" || s.serialNumberEnd !== "") {
              serialNumberForCreate.push(s);
            }
          }
          if (serialNumberForCreate.length > 0) {
            itemsForCreate.push({
              anOrderProductId: item.anOrderProductId,
              serialNumbers: serialNumberForCreate,
            });
          }
        }
        if (itemsForCreate.length > 0) {
          await createSerialNumbers(this.order.anOrderId, itemsForCreate);
        }
        if (serialsNumberForUpdate.length > 0) {
          await updateSerialNumbers(this.order.anOrderId, serialsNumberForUpdate);
        }
      } catch (error) {
        console.log(error);
      }

      this.loading.update = false;
      this.$bvModal.hide("modal-serial-number");
    },

    /**Slip**/

    onClickPrintSlip() {
      const key = Math.random().toString(36).substring(2, 7);
      this.$store.commit(PACKING_SLIP_ADD_ITEM, {key: key, items: this.order})
      let routeData = this.$router.resolve({name: 'PackingSlip', params: {key: key}});
      const url = `${routeData.href}`;
      window.open(url, '_blank');
    },

    /**Robotic**/

    async getPackingData() {
      try {
        const res = await getPackingData(currentDate());
        if (res.data) {
          this.roboticPackingData = res.data;
        }
      } catch (err) {
        console.log(err);
      }
    },

    async onClickRobotSendOrder() {
      this.loadingRobot.sendOrder = true;
      this.setAlert(false);
      let referenceNo = this.order.referenceNo;
      let order = {};

      let hasRoboticSKU = true;
      for (const i of this.order.items) {
        if (!i['roboticSKU']) {
          hasRoboticSKU = false;
          break;
        }
        order[i['roboticSKU']] = i['quantity'];
      }

      if (!hasRoboticSKU) {
        this.setAlert(true, "ไม่สามารถทำรายการได้: ไม่มี robotic SKU", "warning");
      } else {
        try {
          await addOrder(referenceNo, order);
          this.setAlert(true, `${referenceNo}: สั่งทำรายการสำเร็จ`, "success");
        } catch (e) {
          let err = {msg: ""};
          if (e.response && e.response.data) {
            err = e.response.data;
          }
          this.setAlert(true, `ไม่สามารถทำรายการได้: error: ${err.msg}`, "warning");
        }
      }
      this.loadingRobot.sendOrder = false;
    },

    async onClickRobotTakePhoto() {
      this.setAlert(false);
      this.loadingRobot.takePhoto = true;
      this.hasImage = false;
      let referenceNo = this.order.referenceNo;
      try {
        const res = await takePhoto(referenceNo);
        if (res.data) {
          const r = await BlobToDataURL(res.data);
          if (r) {
            this.hasImage = true;
            this.imageBase64 = r;
          }
        }
      } catch (e) {
        console.log(e);
      }
      this.loadingRobot.takePhoto = false;

      if (this.hasImage) {
        this.$bvModal.show("modal-image");
      }
    },

    async onClickRobotSerialNumber() {
      this.setAlert(false);
      this.loadingRobot.takePhoto = true;
      this.hasImage = false;
      let referenceNo = this.order.referenceNo;
      try {
        const res = await takePhoto(referenceNo);
        if (res.data) {
          const r = await BlobToDataURL(res.data);
          if (r) {
            this.hasImage = true;
            this.imageBase64 = r;
          }
        }
      } catch (e) {
        console.log(e);
      }
      this.loadingRobot.takePhoto = false;

      if (this.hasImage) {
        this.$bvModal.show("modal-image");
      }
    },

    async onClickRoboticPackingCard(i, v) {
      this.referenceNo = v.packing_no;
      await this.onSubmitSearch();
    },
  }
}
</script>