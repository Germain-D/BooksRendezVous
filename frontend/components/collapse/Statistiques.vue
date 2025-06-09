<template>
  <details class="collapse bg-base-200 mb-4">
    <summary class="collapse-title text-xl font-medium p-6">
      Statistiques
    </summary>
    <div class="collapse-content px-6">
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div class="stat">
          <div class="stat-title">Totals Livres</div>
          <div class="stat-value">{{ totalBooks }}</div>
        </div>
        <div class="stat">
          <div class="stat-title">Livres lus</div>
          <div class="stat-value">{{ completedBooks }}</div>
        </div>
        <div class="stat">
          <div class="stat-title">En cours</div>
          <div class="stat-value">{{ readingBooks }}</div>
        </div>
        <div class="stat">
          <div class="stat-title">À lire</div>
          <div class="stat-value">{{ toReadBooks }}</div>
        </div>
        <div class="stat">
          <div class="stat-title">Coup de coeur</div>
          <div class="stat-value">{{ favoriteBooks }}</div>
        </div>
        <div class="stat">
          <div class="stat-title">Nombres de Pages</div>
          <div class="stat-value">{{ totalPages }}</div>
        </div>
        <div class="stat">
          <div class="stat-title">Note moyenne</div>
          <div class="stat-value">{{ averageRating }}</div>
        </div>
      </div>
    </div>
  </details>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { useStatsStore } from "@/stores/stats";

const statsStore = useStatsStore();

const totalBooks = computed(() => statsStore.totalBooks);
const completedBooks = computed(() => statsStore.completedBooks);
const readingBooks = computed(() => statsStore.readingBooks);
const toReadBooks = computed(() => statsStore.toReadBooks);
const favoriteBooks = computed(() => statsStore.favoriteBooks);
const totalPages = computed(() => statsStore.totalPages);
const averageRating = computed(() =>
  statsStore.averageRating ? statsStore.averageRating.toFixed(2) : 0
);
</script>

<style scoped>
.collapse {
  @apply rounded-lg overflow-hidden transition-all duration-200;
}

.collapse-title {
  @apply cursor-pointer hover:bg-base-300 transition-colors;
}

.collapse-content {
  @apply pt-0;
}

.collapse[open] summary {
  @apply border-b border-base-300;
}
</style>
