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