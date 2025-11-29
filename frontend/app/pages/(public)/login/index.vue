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

function login() {
  if (validate()) {
    connect(username.value, password.value);
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
  Login
  <NuxtLink to="/signup" class="signup-link">Need an account? Sign up</NuxtLink>
</template>
<style scoped lang="scss">
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
</style>