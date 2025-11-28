<template>
  Sign Up
  <NuxtLink to="/login" class="login-link">Already have an account? Log in</NuxtLink>
</template>
<style scoped lang="scss">
.login-link {
  display: block;
  margin: 1.5rem 0;
  font-weight: 500;
  color: var(--primary);
  text-align: center;
  transition: all 0.2s ease;
  text-decoration: none;

  &:hover {
    color: var(--primary-dark);
    text-decoration: underline;
  }
}
</style>