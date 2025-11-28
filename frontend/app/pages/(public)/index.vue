<script setup lang="ts">

//// Test
const { API } = useApi();
console.log("API", API);
////
</script>

<template>
  <main>
    <section>
      <div class="buttons">
        <NuxtLink :prefetch="false" to="/dashboard" class="btn primary glow">Get Started</NuxtLink>
      </div>
    </section>
  </main>
  <h1>Home Page</h1>
  <p>Structured Notes</p>
</template>

<style scoped lang="scss">
main {
  padding: 0 2rem;
  margin-top: 1em;
  overflow-x: hidden;
}

.buttons {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-top: 1.5rem;
}

.btn {
  padding: 0.9rem 1.8rem;
  border-radius: 50px;
  font-weight: 600;
  text-decoration: none;
}

.btn.primary {
  color: var(--primary);
  background: white;
}

.btn.glow {
  box-shadow: 0 10px 30px rgb(0 0 0 / 20%), 0 0 0 rgb(255 255 255 / 60%);
}
</style>