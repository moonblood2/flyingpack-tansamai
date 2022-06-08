<template>
  <div>
    <ThailandAutoComplete v-model="address.district" type="district" @select="select" label="ตำบล" size="medium"
                          placeholder="ตำบล..."/>
    <ThailandAutoComplete v-model="address.amphoe" type="amphoe" @select="select" label="อำเภอ" size="medium"
                          placeholder="อำเภอ..."/>
    <ThailandAutoComplete v-model="address.province" type="province" @select="select" label="จังหวัด" size="medium"
                          placeholder="จังหวัด..."/>
    <ThailandAutoComplete v-model="address.zipcode" type="zipcode" @select="select" label="รหัสไปรษณีย์" size="medium"
                          placeholder="รหัสไปรษณีย์..."/>
  </div>
</template>

<script>
import ThailandAutoComplete from 'vue-thailand-address-autocomplete'

export default {
  name: "ThaiAddressAutocomplete",
  components: {
    ThailandAutoComplete
  },
  data() {
    return {
      address: {
        district: '',
        amphoe: '',
        province: '',
        zipcode: ''
      }
    }
  },
  watch: {
    address: {
      handler(val) {
        this.$emit('input', val)
      },
      deep: true
    }
  },
  methods: {
    select(val) {
      this.address = val
      this.$emit('select', val)
    }
  }
}
</script>

<style scoped>

</style>