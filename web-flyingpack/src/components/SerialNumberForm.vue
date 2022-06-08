<template>
  <div style="margin-bottom: 10px; padding-bottom: 10px; border-bottom: 1px solid #c2c2c2">
    <div style="display: flex; justify-content: space-between; margin-bottom: 10px">
      <h6 style="width: 50%">{{ title }}</h6>
      <div style="display: flex; width: 50%; justify-content: flex-end">
        <b-input type="number" v-model="value" style="width: 25%"></b-input>
        <b-button @click="onClickPlus()">+</b-button>
      </div>
    </div>
    <div v-for="(s, i) in serialNumbers" :key="i" style="display: flex; justify-content: space-between">
      <b-form-input
          :id="serialNumberElmIDs[i]"
          :placeholder="`Start #${i + 1}`"
          v-model="s.serialNumberStart"
          type="text"
          @keyup="onKeyup($event, i)"
          :state="s.valid"
          style="width: 40%"
      ></b-form-input>
      <p>-</p>
      <b-form-input
          :placeholder="`End`"
          v-model="s.serialNumberEnd"
          type="text"
          style="width: 40%"
          :disabled="!s.serialNumberStart"
      ></b-form-input>
      <b-button :disabled="!!s.anOrderProductSerialNumberId" @click="onClickDelete(i)" variant="danger">-</b-button>
    </div>
  </div>
</template>

<script>
import {AnSerialNumber} from "@/entities/AnSerialNumber";

export default {
  name: "SerialNumberForm",
  props: {
    //Props
    index: Number,
    title: String,
    serialNumbers: Array,
    serialRegex: String,
    //Events
    onChangeSerialNumbers: Function,
  },
  data: function () {
    return {
      uniqueStr: (Math.random() + 1).toString(36).substring(7),
      idCount: 0,
      serialNumberElmIDs: [],
      value: 1,
    }
  },
  created() {
    if (this.serialNumbers.length === 0) {
      this._addSerialNumber(1);
    }
    //Sync serialNumberElmIDs with serialNumbers
    for (let i = 0; i < this.serialNumbers.length; i++) {
      this.serialNumberElmIDs.push(`${this.uniqueStr}:${this.idCount++}`);
    }
  },
  methods: {
    async _addSerialNumber(n) {
      if (n > 0) {
        let serialNumbers = [...this.serialNumbers];
        for (let i = 0; i < n; i++) {
          serialNumbers.push(new AnSerialNumber({}));
          this.serialNumberElmIDs.push(`${this.uniqueStr}:${this.idCount++}`);
        }

        //Callback to update parent component.
        await this.onChangeSerialNumbers(this.index, [...serialNumbers]);
      }
    },
    async _removeSerialNumber(i) {
      let serialNumbers = [...this.serialNumbers];
      serialNumbers.splice(i, 1);
      this.serialNumberElmIDs.splice(i, 1);

      //Callback to update parent component.
      await this.onChangeSerialNumbers(this.index, [...serialNumbers]);
    },
    async onKeyup(event, i) {
      //On hit ENTER/RETURN
      if (event.keyCode === 13) {
        //If run to final serialNumber
        if (i + 1 === this.serialNumbers.length) {
          await this._addSerialNumber(1);
        }

        //Focus next serialNumber
        document.getElementById(this.serialNumberElmIDs[i + 1]).focus();
      } else {
        if (this.serialRegex) {
          let re = new RegExp(this.serialRegex);
          this.serialNumbers[i].valid = re.test(this.serialNumbers[i].serialNumberStart);
          if (this.serialNumbers[i].serialNumberEnd) {
            this.serialNumbers[i].valid &= re.test(this.serialNumbers[i].serialNumberEnd);
          }
        }
      }
    },
    async onClickPlus() {
      await this._addSerialNumber(this.value);
    },
    async onClickDelete(i) {
      await this._removeSerialNumber(i);
    },
  },
}
</script>

<style scoped>

</style>