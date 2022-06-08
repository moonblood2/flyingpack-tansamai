<template>
  <div class="big-box">
    <div class="text-box">
      <h3 style="display: inline">ภาพรวม Fulfillment</h3>
    </div>
    <div class="box">
      <div class="tool-box-1">
        <b-row>
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
      </div>
    </div>
    <b-row>
      <b-col cols="6">
        <FulfillmentPivotTable
            :busy="loading.get"
            :items="fulfillmentPivotTableItems"
        />
      </b-col>
      <b-col cols="6">
        <FulfillmentOrderProductSummaryTable
            :busy="loading.get"
            :items="fulfilmentOrderProductSummaryItems"
        />
      </b-col>
    </b-row>
  </div>
</template>

<script>
import FulfillmentPivotTable from "@/components/FulfillmentPivotTable";
import FulfillmentOrderProductSummaryTable from "@/components/FulfillmentOrderProductSummaryTable";

import {getOrderFulfillmentSummary} from "@/api/agent-network";

import "@/styles/common.css";
import {currentDate} from "@/utils/date";

export default {
  name: "DashboardFulfillment",
  components: {
    FulfillmentPivotTable,
    FulfillmentOrderProductSummaryTable,
  },
  data() {
    return {
      loading: {
        get: false,
      },
      form: {
        startDate: currentDate(),
        endDate: currentDate(),
        startTime: "00:00:00",
        endTime: "23:59:59",
      },

      fulfillmentPivotTableItems: [],
      fulfilmentOrderProductSummaryItems: [],
    }
  },
  created() {
    this.onClickGet();
  },
  methods: {
    async onClickGet() {
      this.loading.get = true;
      try {
        const res = await getOrderFulfillmentSummary(
            `${this.form.startDate} ${this.form.startTime}`,
            `${this.form.endDate} ${this.form.endTime}`,
        );
        if (res.data.data && res.data.data.orderSummary && res.data.data.orderProductSummary) {
          this.fulfillmentPivotTableItems = res.data.data.orderSummary;
          this.fulfilmentOrderProductSummaryItems = res.data.data.orderProductSummary;
        }
      } catch (error) {
        console.log(error);
      }
      this.loading.get = false;
    }
  },
}
</script>