<template>
  <b-container style="margin-top: 80px">
    <b-row>
      <!--Couriers & Products-->
      <b-col cols="4">
        <b-card no-body title="Card Title">
          <b-card-header header-tag="nav">
            <b-nav card-header tabs>
              <b-nav-item :active="showCourier" @click="onSelectCourierTab">
                <b-card-title>เลือกบริการขนส่ง</b-card-title>
              </b-nav-item>
              <b-nav-item :active="!showCourier" @click="onSelectProductTab">
                <b-card-title>เลือกผลิตภัณฑ์</b-card-title>
              </b-nav-item>
            </b-nav>
          </b-card-header>

          <!-- Courier -->
          <PosItemList v-if="showCourier">
            <PosItemCard
              v-for="courier in couriers"
              :key="courier.name"
              :disabled="disabledStates.courierCard"
              v-on:click="onSelectCourier(courier)"
            >
              <template #header>
                <h6 class="mb-0">{{ courier.name }}</h6>
              </template>
              <b-img :alt="courier.name" :src="courier.logo" thumbnail />
            </PosItemCard>
          </PosItemList>
          <!-- Product -->
          <PosItemList v-if="!showCourier">
            <PosItemCard
              v-for="product in products"
              :key="product.id"
              :disabled="disabledStates.productCard"
              v-on:click="onSelectProduct(product)"
            >
              <template #header>
                <h6 class="mb-0">{{ product.name }}</h6>
              </template>
              <template #footer>
                <h6 class="mb-0">{{ product.price }} บาท</h6>
              </template>
            </PosItemCard>
          </PosItemList>
        </b-card>
      </b-col>
      <!--Bill-->
      <b-col cols="8">
        <b-card no-body title="Card Title">
          <b-card-header header-tag="nav">
            <b-card-title>รายการสินค้า</b-card-title>
          </b-card-header>

          <b-card-body>
            <!--index, name, quantity, price, totalPrice, codAmount, handle-->
            <b-table
              :fields="orderFields"
              :items="orderItems"
              responsive="sm"
              striped
            >
              <template #cell(index)="data">
                <p>{{ data.index + 1 }}</p>
              </template>
              <template #cell(quantity)="data">
                <b-button-group v-if="data.item.type === 'product'">
                  <b-button
                    :disabled="
                      disabledStates.decProductBtn ||
                        order[data.item.index].quantity <= 1
                    "
                    size="sm"
                    @click="onClickDecProductQuantity(data.item.index)"
                    >-
                  </b-button>
                  <b-form-input
                    v-model="order[data.item.index].quantity"
                    :disabled="disabledStates.productAmountInp"
                    size="sm"
                    style="width: 80px;"
                    type="number"
                    @blur="onBlurProductQuantity(data.item.index)"
                  ></b-form-input>
                  <b-button
                    :disabled="disabledStates.incProductBtn"
                    size="sm"
                    @click="onClickIncProductQuantity(data.item.index)"
                    >+
                  </b-button>
                </b-button-group>
                <p v-else>{{ data.item.quantity }}</p>
              </template>
              <template #cell(handle)="data">
                <b-button-group>
                  <b-button
                    v-show="data.item.type === 'parcel'"
                    :disabled="disabledStates.editParcelBtn"
                    variant="info"
                    @click="onClickEditParcel(data.item.index)"
                    >แก้ไข
                  </b-button>
                  <b-button
                    :disabled="disabledStates.deleteOrderBtn"
                    variant="danger"
                    @click="onClickDelOrder(data.item.index)"
                    >ลบ
                  </b-button>
                </b-button-group>
              </template>
            </b-table>
            <hr />
            <b-row>
              <b-col>
                <h4>ราคารวม</h4>
              </b-col>
              <b-col class="text-right">
                <h4>{{ orderTotalPrice }} บาท</h4>
              </b-col>
            </b-row>
          </b-card-body>

          <b-card-footer>
            <b-row class="m-2">
              <b-col>
                <b-button
                  :disabled="disabledStates.receiveMoneyBtn"
                  class="w-100 p-3"
                  variant="success"
                  @click="onClickReceiveMoney"
                >
                  รับเงิน
                </b-button>
              </b-col>
            </b-row>
            <b-row class="m-2">
              <b-col>
                <b-button
                  :disabled="disabledStates.receiptBtn"
                  class="w-100 p-3"
                  variant="primary"
                  @click="onClickReceipt"
                >
                  ใบเสร็จ
                </b-button>
              </b-col>
              <b-col>
                <b-dropdown
                  :disabled="disabledStates.labelBtn"
                  class="w-100 p-3"
                  text="ใบปะหน้า"
                  variant="info"
                >
                <!-- <b-dropdown-item
                  @click="
                    onClickLabel('sticker-4x6')
                  "
                  >Sticker4x6</b-dropdown-item
                > -->
                <b-dropdown-item @click="onClickLabel('sticker-8x8')"
                  >Sticker8x8</b-dropdown-item
                >
                </b-dropdown>
              </b-col>
              <b-col>
                <b-button
                  :disabled="disabledStates.newBillBtn"
                  class="w-100 p-3"
                  variant="danger"
                  @click="onClickNewBill"
                >
                  เปิดบิลใหม่
                </b-button>
              </b-col>
            </b-row>
          </b-card-footer>
        </b-card>
      </b-col>
    </b-row>
    <!--sender form-->
    <b-modal
      id="modal-sender-form"
      cancel-title="ยกเลิก"
      ok-title="ยืนยัน"
      title="ฟอร์มยืนยันผู้ใช้บริการส่งพัสดุ ( บัตรประชาชน )"
      @ok="onSubmitSenderForm"
    >
      <b-form @submit.prevent="onSubmitSenderForm">
        <!--ประเภทลูกค้า-->
        <b-form-group id="sender-type" label="ประเภทลูกค้า">
          <b-form-radio-group
            v-model="senderForm.senderType"
            name="sender-type"
          >
            <b-form-radio
              v-model="senderForm.senderType"
              name="natural-person"
              value="1"
              >บุคคลธรรมดา</b-form-radio
            >
            <b-form-radio
              v-model="senderForm.senderType"
              name="juristic-person"
              value="2"
              >นิติบุคคล</b-form-radio
            >
          </b-form-radio-group>
        </b-form-group>
        <!--ชื่อ-นามสกุล / ชื่อบริษัท | วันเกิด (เดิอน/วัน/ค.ศ.)-->
        <b-row>
          <b-col>
            <b-form-group
              id="sender-name"
              :invalid-feedback="senderForm.error.name.message"
              :state="senderForm.error.name.valid"
              label="* ชื่อ-นามสกุล / ชื่อบริษัท"
              label-for="sender-name"
            >
              <b-form-input
                id="sender-name"
                v-model="senderForm.name"
                :state="senderForm.error.name.valid"
                required
                type="email"
              ></b-form-input>
            </b-form-group>
          </b-col>
          <b-col>
            <b-form-group
              id="sender-birthdate"
              label="* วันเกิด (เดือน/วัน/ค.ศ.)"
              label-for="sender-birthdate"
            >
              <b-form-input
                id="sender-birthdate"
                v-model="senderForm.birthDate"
                required
                type="date"
              ></b-form-input>
            </b-form-group>
          </b-col>
        </b-row>
        <!--หมายเลขโทรศัพท์-->
        <b-form-group
          id="sender-phone-number"
          :invalid-feedback="senderForm.error.phoneNumber.message"
          :state="senderForm.error.phoneNumber.valid"
          label="* หมายเลขโทรศัพท์"
          label-for="sender-phone-number"
        >
          <b-form-input
            id="sender-phone-number"
            v-model="senderForm.phoneNumber"
            :state="senderForm.error.phoneNumber.valid"
            required
            type="tel"
          ></b-form-input>
        </b-form-group>
        <!--หมายเลขบัตรประชาชน-->
        <b-form-group
          id="sender-national-id-number"
          :invalid-feedback="senderForm.error.nationalIdNumber.message"
          :state="senderForm.error.nationalIdNumber.valid"
          label="* หมายเลขบัตรประชาชน"
          label-for="sender-national-id-number"
        >
          <b-form-input
            id="sender-national-id-number"
            v-model="senderForm.nationalIdNumber"
            :state="senderForm.error.nationalIdNumber.valid"
            required
            type="tel"
          ></b-form-input>
        </b-form-group>
        <!--* ที่อยู่-->
        <b-form-group
          id="sender-address"
          :invalid-feedback="senderForm.error.address.message"
          :state="senderForm.error.address.valid"
          label="* ที่อยู่"
          label-for="sender-address"
        >
          <b-form-input
            id="sender-address"
            v-model="senderForm.address"
            :state="senderForm.error.address.valid"
            required
            type="text"
          ></b-form-input>
        </b-form-group>
        <!--* ตำบล/แขวง | * อำเภอ/เขต-->
        <b-row>
          <b-col>
            <ThailandAddressAutoComplete
              v-model="senderForm.district"
              label="* ตำบล/แขวง"
              size="default"
              type="district"
              @select="onSelectSenderAddress"
            />
            <b-form-invalid-feedback :state="senderForm.error.district.valid">
              {{ senderForm.error.district.message }}
            </b-form-invalid-feedback>
          </b-col>
          <b-col>
            <ThailandAddressAutoComplete
              v-model="senderForm.state"
              label="* อำเภอ/เขต"
              size="default"
              type="amphoe"
              @select="onSelectSenderAddress"
            />
            <b-form-invalid-feedback :state="senderForm.error.state.valid">
              {{ senderForm.error.state.message }}
            </b-form-invalid-feedback>
          </b-col>
        </b-row>
        <!--* จังหวัด | * รหัสไปรษณีย์-->
        <b-row>
          <b-col>
            <ThailandAddressAutoComplete
              v-model="senderForm.province"
              label="* จังหวัด"
              size="default"
              type="province"
              @select="onSelectSenderAddress"
            />
            <b-form-invalid-feedback :state="senderForm.error.province.valid">
              {{ senderForm.error.province.message }}
            </b-form-invalid-feedback>
          </b-col>
          <b-col>
            <ThailandAddressAutoComplete
              v-model="senderForm.postcode"
              label="* รหัสไปรษณีย์"
              size="default"
              type="zipcode"
              @select="onSelectSenderAddress"
            />
            <b-form-invalid-feedback :state="senderForm.error.postcode.valid">
              {{ senderForm.error.postcode.message }}
            </b-form-invalid-feedback>
          </b-col>
        </b-row>
        <!--เลขประจำตัวผู้เสียภาษี-->
        <b-form-group
          id="sender-tax-id-number"
          label="เลขประจำตัวผู้เสียภาษี"
          label-for="sender-tax-id-number"
        >
          <b-form-input
            id="sender-tax-id-number"
            v-model="senderForm.taxIdNumber"
            required
            type="tel"
          ></b-form-input>
        </b-form-group>
      </b-form>
    </b-modal>
    <!--parcel form ข้อมูลการจัดส่ง create/edit-->
    <ParcelForm
      ref="parcelForm"
      :busy="loading.getPrice || loading.createOrder"
      :disabled-cancel-button="loading.getPrice || loading.createOrder"
      :loading-ok-button="loading.getPrice || loading.createOrder"
      :on-submit="onSubmitParcelForm"
      :parcel-form="parcelForm"
    >
      <template v-slot:afterOriginTitle>
        <b-row>
          <b-col cols="6">
            <b-form-checkbox
              id="parcel-form-same-address"
              v-model="sameAddressAsSender"
              :unchecked-value="false"
              :value="true"
              name="parcel-form-same-address"
              @change="onChangeSameAddress"
            >
              ที่อยู่ตามบัตรประชาชน
            </b-form-checkbox>
          </b-col>
        </b-row>
      </template>
    </ParcelForm>
    <!--ข้อมูลการจัดส่ง, overview-->
    <b-modal
      id="modal-overview-parcel"
      :busy="loading.createOrder"
      cancel-title="ยกเลิก"
      ok-title="ยืนยัน"
      size="lg"
      title="ข้อมูลการจัดส่ง"
    >
      <h6>ต้นทาง: {{ this.overviewParcel.origin }}</h6>
      <h6>ปลายทาง: {{ this.overviewParcel.destination }}</h6>
      <h6>ขนาดพัสดุ: {{ this.overviewParcel.dimension }}</h6>
      <h6>น้ำหนัก: {{ this.overviewParcel.weight }}</h6>
      <h6>ราคา: {{ this.overviewParcel.price }}</h6>
      <template #modal-footer="{cancel}">
        <b-button :disabled="loading.createOrder" @click="cancel"
          >ยกเลิก</b-button
        >
        <b-button
          :disabled="loading.createOrder"
          variant="info"
          @click="onClickEditOverviewParcel"
          >แก้ไข</b-button
        >
        <b-overlay
          :show="loading.createOrder"
          class="d-inline-block"
          opacity="0.4"
          spinner-small
          spinner-variant="primary"
        >
          <b-button
            :disabled="loading.createOrder"
            variant="primary"
            @click="onClickOkOverviewParcel"
            >ตกลง</b-button
          >
        </b-overlay>
      </template>
    </b-modal>
    <!--ข้อมูลการรับเงิน-->
    <b-modal
      id="modal-receive-money"
      :busy="loading.createOrder"
      cancel-title="ยกเลิก"
      ok-title="ยืนยัน"
      title="ข้อมูลการรับเงิน"
      @ok="onSubmitReceiveMoneyForm"
    >
      <b-form @submit.prevent="onSubmitReceiveMoneyForm">
        <b-row>
          <b-col>
            <h6>ยอดรวมทั้งหมด</h6>
          </b-col>
          <b-col class="text-right">
            <h4>{{ orderTotalPrice }} บาท</h4>
          </b-col>
        </b-row>
        <hr />
        <b-row>
          <b-col>
            <h6>รับเงิน</h6>
          </b-col>
          <b-col class="text-right">
            <b-form-input
              id="receive-money-input"
              v-model="receivedMoney.value"
              :state="receivedMoney.valid"
              size="lg"
              type="number"
            />
            <b-form-invalid-feedback :state="receivedMoney.valid">
              จำนวนเงินน้อยกว่าราคารวม
            </b-form-invalid-feedback>
          </b-col>
        </b-row>
        <hr />
        <b-row>
          <b-col>
            <h6>เงินทอน</h6>
          </b-col>
          <b-col>
            <h4 class="text-right">
              {{ receivedMoney.value - orderTotalPrice }} บาท
            </h4>
          </b-col>
        </b-row>
      </b-form>
      <template #modal-footer="{cancel, ok}">
        <b-button :disabled="loading.createOrder" @click="cancel">
          ยกเลิก
        </b-button>
        <b-overlay
          :show="loading.createOrder"
          class="d-inline-block"
          opacity="0.4"
          spinner-small
          spinner-variant="primary"
        >
          <b-button
            :disabled="!receivedMoney.valid"
            variant="primary"
            @click="ok"
            >ตกลง</b-button
          >
        </b-overlay>
      </template>
    </b-modal>
  </b-container>
</template>

<script>
import Vue from "vue";
import axios from "axios";

import ThailandAddressAutoComplete from "vue-thailand-address-autocomplete";
import PosItemList from "@/components/PosItemList";
import PosItemCard from "@/components/PosItemCard";
import ParcelForm from "@/components/ParcelForm";

import env from "@/constants/env";
import {
  ContactInfo,
  Parcel,
  parcelType,
  Product,
  Sender,
  SpOrderParcelShippopFlash,
} from "@/entities";
import { createOrder, getParcelPrice } from "@/api/shipping";

import { LABEL_ADD_PARCEL } from "@/store/actions/label";
import { POS_SET_TYPE } from "@/store/actions/pos";

export default {
  name: "PosForm",
  components: {
    ParcelForm,
    PosItemList,
    PosItemCard,
    ThailandAddressAutoComplete,
  },
  async created() {
    await axios
      .get(`${env.VUE_APP_SERVICE_SHIPPING_URL}/courier-and-product/`, {
        headers: { Authorization: this.$store.state.auth.token },
      })
      .then((res) => {
        this.couriers = res.data.courier;
        this.products = res.data.product;
      })
      .catch((error) => {
        console.log(error);
      });
  },
  beforeDestroy() {
    //Reset type to normal parcel
    this.$store.commit(POS_SET_TYPE, 1);
  },
  data() {
    return {
      //Static data, JNA information
      jnaSenderForm: new Sender({
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
      jnaAddress: new ContactInfo({
        name: env.JNA_NAME,
        phoneNumber: env.JNA_PHONE_NUMBER,
        address: env.JNA_ADDRESS,
        district: env.JNA_DISTRICT,
        state: env.JNA_STATE,
        province: env.JNA_PROVINCE,
        postcode: env.JNA_POSTCODE,
      }),

      //Couriers and products
      showCourier: true, //For select between courier list and product list.
      couriers: [], //For display courier list.
      products: [], //For display product list.

      //Type
      type: this.$store.state.pos.type, //Type of order, create by user = 1 | receive order from AgentNetwork = 2.

      //Sender
      senderForm: new Sender({}),
      submittedSenderForm: false, //To decide when or when not to display the sender form.

      //Parcel
      parcelForm: new Parcel({}),
      parcelFormMode: "create", //There are 2 modes which are "create" and "edit".
      sameAddressAsSender: false, //It's true when user use same address as sender.

      //Editing parcel
      editParcelForm: new Parcel({}),
      editingParcelIndex: null, //Index of order.

      //Order
      order: [], //Order keeps both parcel and product. Use index as key, add -> O(1), remove O(n). Consume with API.
      orderParcels: [], //Use for render only orderParcel, logo.

      //Receive money
      receivedMoney: {
        value: 0,
        valid: false,
      }, //Use in modal-receive-money.

      //Table show summary of order, both orderProduct and orderParcel.
      orderFields: [
        { key: "index", label: "#" },
        { key: "name", label: "ชื่อสินค้า" },
        { key: "quantity", label: "จำนวน" },
        { key: "price", label: "ราคา" },
        { key: "totalPrice", label: "ราคารวม" },
        { key: "codAmount", label: "ยอด COD" },
        { key: "handle", label: "จัดการ" },
      ], //For bootstrap table fields. Header are index, courierName, quantity, price, totalPrice, codAmount, handle | #	ชื่อสินค้า	จำนวน	ราคา	ราคารวม	ยอด COD	จัดการ

      //True when has created order, disable all button in POS page.
      hasCreatedOrder: false,

      //Loading when Request API
      loading: {
        getPrice: false,
        createOrder: false,
      },
    };
  },
  computed: {
    //For bootstrap table, map order to items.
    orderItems: function() {
      let items = [];
      for (let i = 0; i < this.order.length; i++) {
        const val = this.order[i];
        if (val instanceof Parcel) {
          items.push({
            index: i,
            type: "parcel",
            name: `${val.courierName} [${val.origin.name} - ${val.destination.name}]`,
            quantity: 1,
            price: val.price,
            totalPrice: val.price,
            codAmount: val.codAmount,
          });
        } else if (val instanceof Product) {
          items.push({
            index: i,
            type: "product",
            name: val.name,
            quantity: val.quantity,
            price: val.price,
            totalPrice: val.totalPrice,
            codAmount: "-",
          });
        }
      }
      return items;
    },
    //Total price.
    orderTotalPrice: function() {
      let totalPrice = 0.0;
      for (const e of this.order) {
        if (e instanceof Parcel) {
          totalPrice += e.price;
        } else if (e instanceof Product) {
          totalPrice += e.price * e.quantity;
        }
      }
      return totalPrice;
    },
    //overviewParcel is data use in modal-overview-parcel.
    overviewParcel: function() {
      const {
        origin,
        destination,
        width,
        length,
        height,
        weight,
        price,
      } = this.parcelForm;
      return {
        origin: `${origin.name}  ${origin.address}  ${origin.district}  ${origin.state}  ${origin.province}  ${origin.postcode}`,
        destination: `${destination.name}  ${destination.address}  ${destination.district}  ${destination.state}  ${destination.province}  ${destination.postcode}`,
        dimension: `${width}x${length}x${height}`,
        weight: `${weight} กรัม`,
        price: `${price} บาท`,
      };
    },
    // Computed disable status of product/courier card,
    // inc/dec product button, product amount input, delete order button,
    // receive money button, logo button, receipt button, and new bill button.
    disabledStates: function() {
      return {
        productCard: this.hasCreatedOrder,
        courierCard: this.hasCreatedOrder,
        incProductBtn: this.hasCreatedOrder,
        decProductBtn: this.hasCreatedOrder,
        productAmountInp: this.hasCreatedOrder,
        editParcelBtn: this.hasCreatedOrder,
        deleteOrderBtn: this.hasCreatedOrder,
        receiveMoneyBtn: this.order.length <= 0 || this.hasCreatedOrder,
        labelBtn: !this.hasCreatedOrder || this.orderParcels.length === 0,
        receiptBtn: true,
        newBillBtn: !this.hasCreatedOrder,
      };
    },
  },
  watch: {
    // To compute valid status of receivedMoney
    "receivedMoney.value": function(newValue) {
      this.receivedMoney.valid = parseFloat(newValue) >= this.orderTotalPrice;
    },
  },
  methods: {
    //Handlers for switch tab between courier list and product list.
    onSelectCourierTab() {
      this.showCourier = true;
    },
    onSelectProductTab() {
      this.showCourier = false;
    },

    //Handlers for courier card and product card.
    async onSelectCourier(courier) {
      await this.openSenderForm();
      //Setting selected courier's name, provider_code, courier_code and enable_cod.
      this.parcelForm.courierName = courier.name;
      this.parcelForm.providerCode = courier.provider_code;
      this.parcelForm.courierCode = courier.courier_code;
      this.parcelForm.enableCOD = courier.enable_cod;
      //When user select courier set mode to "create".
      this.parcelFormMode = "create";
      this.onChangeSameAddress(true);
      this.$refs.parcelForm.show();
    },
    async onSelectProduct(product) {
      await this.openSenderForm();
      //Search product in order by Linear search.
      let i = 0;
      for (; i < this.order.length; i++) {
        if (
          this.order[i] instanceof Product &&
          this.order[i].id === product.id
        ) {
          break;
        }
      }
      // if already existed then increase its quantity, if not add new product to order.
      if (i < this.order.length) {
        this.onClickIncProductQuantity(i);
      } else {
        this.order.push(
          new Product({
            id: product.id,
            name: product.name,
            price: product.price,
            quantity: 1,
            totalPrice: product.price,
          })
        );
      }
    },

    //Handler for modal-sender-form
    openSenderForm() {
      return new Promise((resolve) => {
        if (!this.submittedSenderForm) {
          this.$bvModal.show("modal-sender-form");
        }
        const waitForSubmitted = () => {
          if (this.submittedSenderForm) {
            resolve(1);
          }
          setTimeout(waitForSubmitted, 50);
        };
        waitForSubmitted();
      });
    },
    onSubmitSenderForm(bvModalEvent) {
      //Validate
      const { valid } = this.senderForm.validate();
      if (valid) {
        //If success, save sender data to order, close sender form and open parcel form.
        this.submittedSenderForm = true;
      } else {
        bvModalEvent.preventDefault();
      }
    },

    //Handler for ParcelForm
    onChangeSameAddress(value) {
      this.sameAddressAsSender = value;
      // Check if the user clicks the same address as the sender,
      // if true set the parcel's origin address the same as the sender.
      if (value === true) {
        this.parcelForm.origin = new ContactInfo({
          name: this.senderForm.name,
          phoneNumber: this.senderForm.phoneNumber,
          address: this.senderForm.address,
          district: this.senderForm.district,
          state: this.senderForm.state,
          province: this.senderForm.province,
          postcode: this.senderForm.postcode,
        });
      } else {
        this.parcelForm.origin = new ContactInfo({});
      }
    },
    onChangeSerialNumber(value, indexOfAnOrderItem, indexOfSerialNumber) {
      this.parcelForm.anOrderItems[indexOfAnOrderItem].serialNumbers[
        indexOfSerialNumber
      ] = value;
    },
    async onSubmitParcelForm(bvModalEvent) {
      bvModalEvent.preventDefault();
      //Validate
      const { valid } = this.parcelForm.validate();
      if (valid) {
        if (this.parcelForm.type === parcelType.PARCEL) {
          this.loading.getPrice = true;
          //POST request getParcelPrice, if success show modal to display price with only close button and Add parcel to order.parcels.
          try {
            const res = await getParcelPrice(this.parcelForm.toPayload());
            this.loading.getPrice = false;
            if (res.data[0]) {
              const data = res.data[0];
              if (data.status) {
                //Close ParcelForm
                this.$refs.parcelForm.hide();
                //Open modal-overview-parcel
                this.$bvModal.show("modal-overview-parcel");
                this.parcelForm.price = data.price;
              }
            } else {
              //TODO: Show error message.
            }
          } catch (error) {
            console.log(error);
            this.loading.getPrice = false;
          }
        } else {
          console.error("Parcel type mismatch.");
        }
      }
    },

    //Handlers for modal-overview-parcel
    async onClickOkOverviewParcel() {
      if (this.parcelFormMode === "create") {
        this.order.push(this.parcelForm.clone());
      } else if (this.parcelFormMode === "edit") {
        //Check if not set editingParcelIndex.
        if (this.editingParcelIndex === null) {
          throw new Error("editingParcelIndex is null");
        }
        //Use Vue.set to replace order in array, to trigger vue update dom.
        Vue.set(this.order, this.editingParcelIndex, this.parcelForm.clone());
      } else {
        throw new Error("parcelFormMode is neither create nor edit.");
      }

      if (this.type === 2) {
        this.loading.createOrder = true;
        await this.createOrder();
        this.loading.createOrder = false;
      }

      //Reset form.
      this.parcelForm = new Parcel({});
      this.$bvModal.hide("modal-overview-parcel");
    },
    onClickEditOverviewParcel() {
      //Open ParcelForm
      this.$refs.parcelForm.show();
      //Close modal-overview-form
      this.$bvModal.hide("modal-overview-parcel");
    },

    //Handlers for select address of thailand address auto-compete.
    onSelectSenderAddress(address) {
      this.senderForm.district = address.district;
      this.senderForm.state = address.amphoe;
      this.senderForm.province = address.province;
      this.senderForm.postcode = address.zipcode;
    },
    onSelectParcelOriginAddress(address) {
      this.parcelForm.origin.district = address.district;
      this.parcelForm.origin.state = address.amphoe;
      this.parcelForm.origin.province = address.province;
      this.parcelForm.origin.postcode = address.zipcode;
    },
    onSelectParcelDestinationAddress(address) {
      this.parcelForm.destination.district = address.district;
      this.parcelForm.destination.state = address.amphoe;
      this.parcelForm.destination.province = address.province;
      this.parcelForm.destination.postcode = address.zipcode;
    },

    //Handlers for increase/decrease quantity of product order.
    //Decrease product quantity
    onClickDecProductQuantity(orderIndex) {
      this.order[orderIndex].quantity -= 1;
    },
    //Increase product quantity
    onClickIncProductQuantity(orderIndex) {
      this.order[orderIndex].quantity += 1;
    },
    //On user out of focus, check if product quantity is in acceptable range.
    onBlurProductQuantity(orderIndex) {
      if (this.order[orderIndex].quantity < 1) {
        this.order[orderIndex].quantity = 1;
      }
    },

    //Handler for delete order both product and parcel.
    onClickDelOrder(orderIndex) {
      this.order.splice(orderIndex, 1);
    },

    //Handler for edit parcel by ParcelForm.
    onClickEditParcel(orderIndex) {
      this.parcelFormMode = "edit";
      this.editingParcelIndex = orderIndex;
      this.parcelForm = this.order[orderIndex].clone();
      this.$refs.parcelForm.show();
    },

    //Handler for check-bill button group.
    //Onclick receive money (รับเงิน) button
    onClickReceiveMoney() {
      //Open modal-receive-money
      this.$bvModal.show("modal-receive-money");
    },
    //Onclick receipt (ใบเสร็จ) button
    onClickReceipt() {},
    //Onclick label (ใบปะหน้า) button
    onClickLabel(size) {
      const key = Math.random()
        .toString(36)
        .substring(2, 7);
      this.$store.commit(LABEL_ADD_PARCEL, {
        key: key,
        parcels: this.orderParcels,
      });
      let routeData = this.$router.resolve({
        name: "Label",
        params: { key: key },
      });
      const url = `${routeData.href}?size=${size}`;
      window.open(url, "_blank");
    },
    //Onclick new bill (เปิดบิลใหม่) button
    onClickNewBill() {
      location.reload();
    },

    //Onsubmit submit modal-receive-money
    async onSubmitReceiveMoneyForm(bvModalEvent) {
      bvModalEvent.preventDefault();
      this.loading.createOrder = true;
      await this.createOrder();
      this.loading.createOrder = false;
      this.$bvModal.hide("modal-receive-money");
    },

    async createOrder() {
      let response = false;
      //Map order to products and parcels.
      let products = [],
        parcels = [];
      for (const e of this.order) {
        if (e instanceof Product) {
          products.push(e.toPayload());
        } else if (e instanceof Parcel) {
          this.orderParcels.push(e);
          parcels.push(e.toPayload());
        } else {
          throw new Error("e is neither Product nor Parcel");
        }
      }
      try {
        const res = await createOrder({
          sender: this.senderForm.toPayload(),
          products: products,
          parcels: parcels,
          paymentMethod: 1,
        });
        const { data } = res;
        //Check status
        if (data.status === true) {
          this.hasCreatedOrder = true;
          //Populate tracking code to order parcel
          if (data.parcels !== undefined) {
            for (let i = 0; i < data.parcels.length; i++) {
              this.orderParcels[i].trackingCode =
                data.parcels[i]["tracking_code"];
              if (data.parcels[i]["shippop_flash_sorting_code"]) {
                this.orderParcels[
                  i
                ].spOrderParcelShippopFlash = new SpOrderParcelShippopFlash({
                  dstCode:
                    data.parcels[i]["shippop_flash_sorting_code"]["dst_code"],
                  sortCode:
                    data.parcels[i]["shippop_flash_sorting_code"]["sort_code"],
                  sortingLineCode:
                    data.parcels[i]["shippop_flash_sorting_code"][
                      "sorting_line_code"
                    ],
                });
              }
            }
          }
          response = true;
          await this.showMsgBoxOrderCreating();
        } else {
          await this.showMsgBoxOrderCreating("ไม่สำเร็จ กรุณาทำรายการใหม่");
        }
      } catch (error) {
        console.log(error);
      }

      return response;
    },

    async showMsgBoxOrderCreating(message = "สำเร็จ") {
      await this.$bvModal.msgBoxOk(message, {
        title: "การทำรายการ",
        size: "sm",
        buttonSize: "md",
        okVariant: "primary",
        headerClass: "p-2 border-bottom-0",
        footerClass: "p-2 border-top-0",
        centered: true,
      });
    },
  },
};
</script>
