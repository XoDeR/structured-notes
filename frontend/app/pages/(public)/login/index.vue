<script setup lang="ts">

const userStore = useUserStore();
const router = useRouter();

// form fields
const username = ref('');
const password = ref('');

// form fields errors
const errors = ref({
  username: '',
  password: '',
  general: '',
});

// password related
const showPassword = ref(false);

function togglePassword() {
  showPassword.value = !showPassword.value;
}

function validate(): boolean {
  let valid = true;

  if (!username.value) {
    errors.value.username = 'Username is required';
    valid = false;
  }
  else {
    errors.value.username = '';
  }

  if (!password.value) {
    errors.value.password = 'Password is required';
    valid = false;
  }
  else {
    errors.value.password = '';
  }

  return valid;
}

async function login() {
  if (validate()) {
    await connect(username.value, password.value);
  }
}

async function connect(username: string, password: string) {
  const result = await userStore.login(username, password);
  if (result.success) {
    router.push('/dashboard');
  } else {
    errors.value.general = result.errorMessage!;
  }
}

</script>

<template>
  <div class="container">
    <div class="body-container">
      <h1>Login</h1>
      <form @submit.prevent="login">
        <div class="form-group">
          <label for="username">Username</label>
          <input id="username" v-model="username" type="username" :class="{ 'is-invalid': errors.username }" />
          <p v-if="errors.username" class="invalid-feedback">{{ errors.username }}</p>
        </div>
        <div class="form-group">
          <label for="password">Password</label>
          <div class="password-input">
            <input id="password" :key="`password-${showPassword}`" v-model="password"
              :type="showPassword ? 'text' : 'password'" :class="{ 'is-invalid': errors.password }" />
            <button type="button" class="password-toggle" @click="togglePassword">
              <div class="eye-icon" :class="{ show: showPassword }">
                <Icon v-if="showPassword" name="eye" />
                <Icon v-else name="eye_off" />
              </div>
            </button>
          </div>
          <p v-if="errors.password" class="invalid-feedback">{{ errors.password }}</p>
        </div>
        <button type="submit" class="btn">Login</button>
        <p v-if="errors.general" class="invalid-feedback">{{ errors.general }}</p>
        <p class="forgot-password-link">Forgot your password? <NuxtLink to="/login/request-reset">Click here</NuxtLink>
        </p>
      </form>
      <NuxtLink to="/signup" class="signup-link">Need an account? Sign up</NuxtLink>
    </div>
  </div>
</template>
<style scoped lang="scss">
.container {
  display: flex;
  width: 95%;
  height: 100%;
  margin: 0 auto;
  flex-direction: column;
  justify-content: space-between;
  padding-top: 1.5rem;
}

.body-container {
  display: flex;
  width: 100%;
  max-width: 600px;
  margin: 0 auto 10%;
  align-items: center;
  flex-direction: column;
  justify-content: center;
}

form {
  width: 100%;
}

.form-group {
  display: flex;
  width: 100%;
  flex-direction: column;
  gap: 6px;
  margin-bottom: 0.8rem;

  label {
    font-size: 0.95rem;
    font-weight: 600;
  }
}

input {
  width: 100%;
  box-sizing: border-box;
  padding: 0.6rem;
  border: 2px solid var(--border-color);
  border-radius: 8px;
  font-size: 1rem;
  background: var(--bg-color);
  transition: all 0.2s ease;

  &:focus {
    border-color: var(--primary);
    box-shadow: 0 0 0 3px rgb(var(--primary-rgb), 0.1);
    outline: none;
  }
}

.password-input {
  position: relative;
}

.password-toggle {
  position: absolute;
  top: 50%;
  right: 8px;
  display: flex;
  padding: 8px;
  border: none;
  border-radius: 50%;
  background: transparent;
  transition: all 0.2s ease;
  align-items: center;
  cursor: pointer;
  justify-content: center;
  transform: translateY(-50%);

  &:hover {
    background: rgb(0 0 0 / 5%);
  }

  &:active {
    transform: translateY(-50%) scale(0.95);
  }
}

.eye-icon {
  display: flex;
  width: 20px;
  height: 20px;
  transition: all 0.3s ease;
  align-items: center;
  justify-content: center;

  &.show {
    transform: scale(1.1);
  }

  font-size: 32px;
}

.signup-link {
  display: block;
  font-weight: 500;
  color: var(--primary);
  text-align: center;
  transition: all 0.2s ease;
  margin-bottom: 1rem;
  text-decoration: none;

  &:hover {
    color: var(--primary-dark);
    text-decoration: underline;
  }
}

.btn {
  width: 100%;
  padding: 0.9rem;
  border: none;
  border-radius: 12px;
  font-size: 1.1rem;
  font-weight: 600;
  color: white;
  background-color: var(--default);
  transition: all 0.2s ease;
  cursor: pointer;
  margin-top: 1rem;

  &:hover {
    background: var(--default-dark);
    box-shadow: 0 4px 12px rgb(var(--default-rgb), 0.3);
    transform: translateY(-1px);
  }

  &:active {
    box-shadow: 0 2px 6px rgb(var(--default-rgb), 0.3);
    transform: translateY(0);
  }
}

.forgot-password-link {
  font-size: 16px;
  font-weight: 500;
  text-align: center;
  margin-top: 1rem;

  a {
    font-weight: 600;
    color: var(--primary);
    transition: all 0.2s ease;
    text-decoration: none;

    &:hover {
      color: var(--primary-dark);
      text-decoration: underline;
    }
  }
}

.is-invalid {
  border-color: var(--red) !important;
}

.invalid-feedback {
  margin: 0;
  font-size: 0.8rem;
  color: var(--red);
  text-align: center;
}
</style>