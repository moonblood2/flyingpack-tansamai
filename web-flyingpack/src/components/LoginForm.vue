<template>
  <b-container style="margin-top: 80px">
    <b-row class="d-flex justify-content-center">
      <b-col class="login-600">
        <b-card>
          <b-form @submit.prevent="onSubmit">
            <h5>เข้าสู่ระบบ</h5>
            <b-form-group
                id="input-group-1"
                :invalid-feedback="email.validation.message"
                :state="email.validation.status"
                label-for="input-1"
            >
              <b-form-input
                  id="input-1"
                  v-model="email.value"
                  :state="email.validation.status"
                  name="email"
                  placeholder="Email"
                  required
                  type="email"
                  @keyup="onFormChange($event)"
              ></b-form-input>
            </b-form-group>
            <b-form-group
                id="input-group-2"
                :invalid-feedback="password.validation.message"
                :state="password.validation.status"
                label-for="input-2"
            >
              <b-form-input
                  id="input-2"
                  v-model="password.value"
                  :state="password.validation.status"
                  name="password"
                  placeholder="Password"
                  required
                  type="password"
                  @keyup="onFormChange($event)"
              ></b-form-input>
            </b-form-group>
            <b-overlay
                :show="authLoading"
                class="d-inline-block"
                opacity="0.6"
                rounded
                spinner-small
                spinner-variant="primary"
            >
              <b-button type="submit" variant="primary">Login</b-button>
            </b-overlay>
            <b-form-invalid-feedback
                :state="authValidation.valid"
            >
              {{ authValidation.details }}
            </b-form-invalid-feedback>
          </b-form>
        </b-card>
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
import {mapGetters} from 'vuex'
import {AUTH_REQUEST} from '@/store/actions/auth'
import {getDefaultPath} from "@/router";

export default {
  name: "LoginForm",
  data() {
    return {
      email: {
        value: "",
        validator: {
          required: true,
          pattern: "email"
        },
        validation: {status: null, message: ""}
      },
      password: {
        value: "",
        validator: {
          required: true
        },
        validation: {status: null, message: ""}
      },
      formValid: null
    }
  },
  computed: {
    ...mapGetters(['authLoading', 'authValidation'])
  },
  watch: {
    authValidation(newValue) {
      if (newValue.valid === false) {
        this.email.validation.status = false
        this.password.validation.status = false
      }
    }
  },
  methods: {
    onFormChange(event) {
      const name = event.target.name;
      const value = event.target.value;

      let updatedElement = {...this[name]};
      const validatorObject = this.checkValidator(value, updatedElement.validator);
      updatedElement.value = value;
      updatedElement.validation = {
        status: validatorObject.status,
        message: validatorObject.message
      };
      this[name] = {...updatedElement};
    },
    checkValidator(value, rule) {
      let valid = true;
      let message = "";

      if (valid && rule.required && value.length === 0) {
        valid = false;
        message = "ต้องกรอก"
      }

      if (valid && rule.pattern === "email") {
        if (/^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$/.test(value) == false) {
          valid = false;
          message = "ต้องเป็น e-mail เท่านั้น"
        }
      }

      return {status: valid, message: message};
    },
    onSubmit() {
      this.$store.dispatch(AUTH_REQUEST, {email: this.email.value, password: this.password.value})
          .then(() => {
            this.$router.push(getDefaultPath());
          })
    }
  }
}
</script>

<style scoped>
.login-600 {
  max-width: 600px;
}
</style>