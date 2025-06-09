<template>
  <div class="max-w-7xl m-16">
    <h1 class="text-3xl font-bold mb-8">Mon Profil</h1>

    <CollapsePersonalData />

    <CollapsePasswordChange />

    <CollapseProfilVisibility />

    <CollapseStatistiques />
    <CollapseSuccess :ready="false" />
  </div>
</template>

<script setup lang="ts">
// Page configuration
definePageMeta({
  layout: "connected",
  middleware: ["guard", "stats"],
});
</script>
