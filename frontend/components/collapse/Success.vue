<template>
  <details class="collapse bg-base-200 mb-4">
    <summary class="collapse-title text-xl font-medium p-6 relative">
      Succès
    </summary>
    <div class="collapse-content px-6">
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <!-- Achievement Card -->
        <div
          class="card bg-base-300 p-4"
          v-for="achievement in achievementsStore.achievements"
          :key="achievement.name"
        >
          <div class="flex items-center gap-4">
            <div
              :class="achievement.unlockedAt ? 'bg-green-700' : 'bg-gray-700'"
              class="p-3 rounded-full"
            >
              <!-- Utilisation de l'image si disponible, sinon utilisation de l'icône par défaut -->
              <img
                v-if="achievement.image"
                :src="achievement.image"
                alt="Achievement icon"
                class="h-6 w-6 object-cover"
              />
              <svg
                v-else
                xmlns="http://www.w3.org/2000/svg"
                class="h-6 w-6"
                :class="
                  achievement.unlockedAt ? 'text-green-200' : 'text-gray-400'
                "
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"
                />
              </svg>
            </div>
            <div>
              <h3
                class="font-bold"
                :class="
                  achievement.unlockedAt ? 'text-green-400' : 'text-gray-400'
                "
              >
                {{ achievement.name }}
              </h3>
              <p class="text-sm text-gray-500">{{ achievement.description }}</p>
              <p
                v-if="achievement.unlockedAt"
                class="text-xs text-green-400 mt-1"
              >
                Débloqué le
                {{ new Date(achievement.unlockedAt).toLocaleDateString() }}
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </details>
</template>

<script setup lang="ts">
import { useAchievementsStore } from "@/stores/achievements";

const achievementsStore = useAchievementsStore();
achievementsStore.fetchAchievements();
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
