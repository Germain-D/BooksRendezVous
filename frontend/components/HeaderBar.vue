<template>
    <div class="mb-16 mt-16 mx-4 text-center relative">
        <div class="absolute top-0 left-1/2 -translate-x-1/2 w-32 h-32 bg-primary/20 rounded-full blur-3xl"></div>
        <h1 class="text-6xl font-bold mb-6 tracking-tight">
          <span class="bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">{{ pseudo }}</span>
          
        </h1>
        <h1 class="text-6xl font-bold mb-6 tracking-tight">
        <span class="text-base-content/80">B<span class="bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">oo</span>ksRendezVous</span></h1>
        <p class="text-lg text-base-content/60 max-w-2xl mx-auto leading-relaxed">
          {{ description }}
        </p>
      </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';

/* add props */

const props = defineProps({
  pseudo: String,
});


const description = "Explorez ma bibliothèque personnelle : découvrez les livres que je lis actuellement, ceux que j'ai terminés et ma liste de lecture à venir.";
</script>