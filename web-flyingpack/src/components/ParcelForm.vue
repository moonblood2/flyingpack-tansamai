<template>
  <b-modal
      id="modal-parcel-form"
      :busy="busy"
      :no-close-on-backdrop="busy"
      :no-close-on-esc="busy"
      :title="`ข้อมูลการจัดส่ง [${this.parcelForm.courierName}]`"
      :visible="visible"
      cancel-title="ยกเลิก"
      ok-title="ยืนยัน"
      size="lg"
      @cancel="handleCancel"
      @close="handleCancel"
      @ok="handleOk"
  >
    <b-form @submit.prevent="handleOk">
      <div v-if="showOrigin">
        <!--ต้นทาง-->
        <h5>ต้นทาง</h5>
        <slot name="afterOriginTitle"></slot>
        <!--ชื่อ-นามสกุล | หมายเลขโทรศัพท์-->
        <b-row>
          <b-col>
            <b-form-group
                id="parcel-name"
                :invalid-feedback="parcelForm.origin.error.name.message"
                :state="parcelForm.origin.error.name.valid"
                label="* ชื่อ-นามสกุล"
                label-for="parcel-name-input"
            >
              <b-form-input
                  id="parcel-name-input"
                  v-model="parcelForm.origin.name"
                  :state="parcelForm.origin.error.name.valid"
                  required
                  type="email"
              ></b-form-input>
            </b-form-group>
          </b-col>
          <b-col>
            <b-form-group
                id="parcel-phone-number"
                :invalid-feedback="parcelForm.origin.error.phoneNumber.message"
                :state="parcelForm.origin.error.phoneNumber.valid"
                label="* หมายเลขโทรศัพท์"
                label-for="parcel-phone-number-input"
            >
              <b-form-input
                  id="parcel-phone-number-input"
                  v-model="parcelForm.origin.phoneNumber"
                  :state="parcelForm.origin.error.phoneNumber.valid"
                  required
                  type="tel"
              ></b-form-input>
            </b-form-group>
          </b-col>
        </b-row>
        <!--* ที่อยู่ | * ตำบล/แขวง-->
        <b-row>
          <b-col>
            <b-form-group
                id="parcel-address"
                :invalid-feedback="parcelForm.origin.error.address.message"
                :state="parcelForm.origin.error.address.valid"
                label="* ที่อยู่"
                label-for="parcel-address-input"
            >
              <b-form-input
                  id="parcel-address-input"
                  v-model="parcelForm.origin.address"
                  :state="parcelForm.origin.error.address.valid"
                  required
                  type="text"
              ></b-form-input>
            </b-form-group>
          </b-col>
          <b-col>
            <ThailandAddressAutoComplete v-model="parcelForm.origin.district" label="* ตำบล/แขวง"
                                         size="default"
                                         type="district" @select="onSelectParcelOriginAddress"/>
            <b-form-invalid-feedback :state="parcelForm.origin.error.district.valid">
              {{ parcelForm.origin.error.district.message }}
            </b-form-invalid-feedback>
          </b-col>
        </b-row>
        <!--* อำเภอ/เขต | * จังหวัด | * รหัสไปรษณีย์-->
        <b-row>
          <b-col>
            <ThailandAddressAutoComplete v-model="parcelForm.origin.state" label="* อำเภอ/เขต"
                                         size="default"
                                         type="amphoe" @select="onSelectParcelOriginAddress"/>
            <b-form-invalid-feedback :state="parcelForm.origin.error.state.valid">
              {{ parcelForm.origin.error.state.message }}
            </b-form-invalid-feedback>
          </b-col>
          <b-col>
            <ThailandAddressAutoComplete v-model="parcelForm.origin.province" label="* จังหวัด"
                                         size="default"
                                         type="province" @select="onSelectParcelOriginAddress"/>
            <b-form-invalid-feedback :state="parcelForm.origin.error.province.valid">
              {{ parcelForm.origin.error.province.message }}
            </b-form-invalid-feedback>
          </b-col>
          <b-col>
            <ThailandAddressAutoComplete v-model="parcelForm.origin.postcode" label="* รหัสไปรษณีย์"
                                         size="default"
                                         type="zipcode" @select="onSelectParcelOriginAddress"/>
            <b-form-invalid-feedback :state="parcelForm.origin.error.postcode.valid">
              {{ parcelForm.origin.error.postcode.message }}
            </b-form-invalid-feedback>
          </b-col>
        </b-row>
        <hr/>
      </div>

      <!--ปลายทาง-->
      <h5>ปลายทาง</h5>
      <!--ชื่อ-นามสกุล | หมายเลขโทรศัพท์-->
      <b-row>
        <b-col>
          <b-form-group
              id="parcel-name"
              :invalid-feedback="parcelForm.destination.error.name.message"
              :state="parcelForm.destination.error.name.valid"
              label="* ชื่อ-นามสกุล"
              label-for="parcel-name-input"
          >
            <b-form-input
                id="parcel-name-input"
                v-model="parcelForm.destination.name"
                :state="parcelForm.destination.error.name.valid"
                required
                type="email"
            ></b-form-input>
          </b-form-group>
        </b-col>
        <b-col>
          <b-form-group
              id="parcel-phone-number"
              :invalid-feedback="parcelForm.destination.error.phoneNumber.message"
              :state="parcelForm.destination.error.phoneNumber.valid"
              label="* หมายเลขโทรศัพท์"
              label-for="parcel-phone-number-input"
          >
            <b-form-input
                id="parcel-phone-number-input"
                v-model="parcelForm.destination.phoneNumber"
                :state="parcelForm.destination.error.phoneNumber.valid"
                required
                type="tel"
            ></b-form-input>
          </b-form-group>
        </b-col>
      </b-row>
      <!--* ที่อยู่ | * ตำบล/แขวง-->
      <b-row>
        <b-col>
          <b-form-group
              id="parcel-address"
              :invalid-feedback="parcelForm.destination.error.address.message"
              :state="parcelForm.destination.error.address.valid"
              label="* ที่อยู่"
              label-for="parcel-address-input"
          >
            <b-form-input
                id="parcel-address-input"
                v-model="parcelForm.destination.address"
                :state="parcelForm.destination.error.address.valid"
                required
                type="text"
            ></b-form-input>
          </b-form-group>
        </b-col>
        <b-col>
          <ThailandAddressAutoComplete v-model="parcelForm.destination.district" label="* ตำบล/แขวง"
                                       size="default"
                                       type="district" @select="onSelectParcelDestinationAddress"/>
          <b-form-invalid-feedback :state="parcelForm.destination.error.district.valid">
            {{ parcelForm.destination.error.district.message }}
          </b-form-invalid-feedback>
        </b-col>
      </b-row>
      <!--* อำเภอ/เขต | * จังหวัด | * รหัสไปรษณีย์-->
      <b-row>
        <b-col>
          <ThailandAddressAutoComplete v-model="parcelForm.destination.state" label="* อำเภอ/เขต"
                                       size="default"
                                       type="amphoe" @select="onSelectParcelDestinationAddress"/>
          <b-form-invalid-feedback :state="parcelForm.destination.error.state.valid">
            {{ parcelForm.destination.error.state.message }}
          </b-form-invalid-feedback>
        </b-col>
        <b-col>
          <ThailandAddressAutoComplete v-model="parcelForm.destination.province" label="* จังหวัด"
                                       size="default"
                                       type="province" @select="onSelectParcelDestinationAddress"/>
          <b-form-invalid-feedback :state="parcelForm.destination.error.province.valid">
            {{ parcelForm.destination.error.province.message }}
          </b-form-invalid-feedback>
        </b-col>
        <b-col>
          <ThailandAddressAutoComplete v-model="parcelForm.destination.postcode" label="* รหัสไปรษณีย์"
                                       size="default"
                                       type="zipcode" @select="onSelectParcelDestinationAddress"/>
          <b-form-invalid-feedback :state="parcelForm.destination.error.postcode.valid">
            {{ parcelForm.destination.error.postcode.message }}
          </b-form-invalid-feedback>
        </b-col>
      </b-row>
      <hr/>
      <!--รายละเอียดเพิ่มเติม-->
      <h5>รายละเอียดเพิ่มเติม</h5>
      <b-row>
        <b-col>
          <b-form-group
              id="parcel-weight"
              :invalid-feedback="parcelForm.error.weight.message"
              :state="parcelForm.error.weight.valid"
              label="* น้ำหนัก ( กรัม )"
              label-for="parcel-weight-input"
          >
            <b-form-input
                id="parcel-weight-input"
                v-model="parcelForm.weight"
                :state="parcelForm.error.weight.valid"
                required
                type="number"
            ></b-form-input>
          </b-form-group>
        </b-col>
        <b-col>
          <b-form-group
              id="parcel-shape"
              :invalid-feedback="`${parcelForm.error.width.message}, ${parcelForm.error.length.message} and ${parcelForm.error.height.message}`"
              :state="parcelForm.error.width.valid && parcelForm.error.length.valid && parcelForm.error.height.valid"
              label="* ขนาดพัสดุ ( เซนติเมตร )"
              label-for="parcel-shape-input"
          >
            <b-row>
              <b-col>
                <b-form-input
                    id="parcel-shape-width"
                    v-model="parcelForm.width"
                    :state="parcelForm.error.width.valid"
                    placeholder="กว้าง"
                    required
                    type="number"
                ></b-form-input>
              </b-col>
              <b-col>
                <b-form-input
                    id="parcel-shape-length"
                    v-model="parcelForm.length"
                    :state="parcelForm.error.length.valid"
                    placeholder="ยาว"
                    required
                    type="number"
                ></b-form-input>
              </b-col>
              <b-col>
                <b-form-input
                    id="parcel-shape-height"
                    v-model="parcelForm.height"
                    :state="parcelForm.error.height.valid"
                    placeholder="สูง"
                    required
                    type="number"
                ></b-form-input>
              </b-col>
            </b-row>
          </b-form-group>
        </b-col>
      </b-row>
      <b-row>
        <b-col cols="6">
          <b-form-group
              id="parcel-cod-amount"
              :invalid-feedback="parcelForm.enableCOD ? parcelForm.error.codAmount.message: ''"
              :state="parcelForm.enableCOD ? parcelForm.error.codAmount.valid: null"
              label="เก็บเงินปลายทาง ( บาท )"
              label-for="parcel-cod-amount-input"
          >
            <b-form-input
                id="parcel-cod-amount-input"
                v-model="parcelForm.codAmount"
                :disabled="!parcelForm.enableCOD"
                :state="parcelForm.enableCOD ? parcelForm.error.codAmount.valid: null"
                required
                type="number"
            ></b-form-input>
          </b-form-group>
        </b-col>
      </b-row>
      <!--รายละเอียดสินค้า-->
      <div v-if="this.parcelForm.type === parcelType.AN_PARCEL">
        <hr/>
        <h5>รายละเอียดสินค้า</h5>
        <b-row
            v-for="(item, i) in parcelForm.items"
            :key="i"
            style="margin-bottom: 10px"
        >
          <b-col cols="6">
            <p>
              {{ item.productCode }}: {{ item.quantity }} ชิ้น
            </p>
          </b-col>
          <b-col cols="6">
            <SerialNumberForm :serial-numbers="item.serialNumbers"/>
          </b-col>
        </b-row>
      </div>
    </b-form>
    <template #modal-footer="{cancel, ok}">
      <b-button :disabled="disabledCancelButton" @click="cancel">ยกเลิก</b-button>
      <b-overlay
          :show="loadingOkButton"
          class="d-inline-block"
          opacity="0.4"
          spinner-small
          spinner-variant="primary"
      >
        <b-button variant="primary" @click="ok">ตกลง</b-button>
      </b-overlay>
    </template>
  </b-modal>
</template>

<script>
import SerialNumberForm from "@/components/SerialNumberForm";

import ThailandAddressAutoComplete from "vue-thailand-address-autocomplete";
import {AnParcel, Parcel} from "@/entities";
import {parcelType} from "@/entities/Parcel";

export default {
  name: "ParcelForm",
  components: {
    SerialNumberForm,
    ThailandAddressAutoComplete,
  },
  props: {
    visible: Boolean,
    busy: Boolean,
    parcelForm: {
      type: [Parcel, AnParcel],
      default: function () {
        return new Parcel({});
      },
    },
    title: String,
    showOrigin: Boolean,
    loadingOkButton: Boolean,
    disabledCancelButton: Boolean,
    onSubmit: Function,
    onCancel: Function,

    disableValidation: Boolean,
  },
  data() {
    return {
      parcelType: parcelType,
    }
  },
  methods: {
    onSelectParcelOriginAddress(address) {
      this.parcelForm.origin.district = address.district
      this.parcelForm.origin.state = address.amphoe
      this.parcelForm.origin.province = address.province
      this.parcelForm.origin.postcode = address.zipcode
    },
    onSelectParcelDestinationAddress(address) {
      this.parcelForm.destination.district = address.district
      this.parcelForm.destination.state = address.amphoe
      this.parcelForm.destination.province = address.province
      this.parcelForm.destination.postcode = address.zipcode
    },
    onChangeSerialNumber(value, indexOfAnOrderItem, indexOfSerialNumber) {
      this.parcelForm.anOrderItems[indexOfAnOrderItem].serialNumbers[indexOfSerialNumber] = value;
    },
    handleOk(bvModalEvt) {
      bvModalEvt.preventDefault();
      if (this.disableValidation) {
        this.onSubmit(bvModalEvt);
      } else {
        const {valid} = this.parcelForm.validate();
        if (valid) {
          this.onSubmit(bvModalEvt);
        }
      }
    },
    handleCancel(bvModalEvt) {
      this.onCancel(bvModalEvt);
    },
    show() {
      this.$bvModal.show("modal-parcel-form");
    },
    hide() {
      this.$bvModal.hide("modal-parcel-form");
    }
  },
}
</script>

<style scoped>

</style>