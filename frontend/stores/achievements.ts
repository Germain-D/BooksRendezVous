import { defineStore } from 'pinia';



// Achievement interface
interface Achievement {
    name: string;
    description: string;
    type: string;
    unlockedAt: string;
    image: string;

}

// Store state interface
interface AchievementsState {
    achievements: Achievement[];
    loading: boolean;
    error: string | null;
}

export const useAchievementsStore = defineStore('achievements', {
    state: (): AchievementsState => ({
        achievements: [],
        loading: false,
        error: null
    }),

    getters: {
        getAchievements: (state) => state.achievements,

    },

    actions: {
        async fetchAchievements() {
            this.loading = true;
            try {
                const config = useRuntimeConfig();
                const response = await fetch(`${config.public.BACKEND_URL}/api/achievements`, {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': 'Bearer ' + localStorage.getItem('jwt')
                    },
                    credentials: 'include'
                });
                
                const data = await response.json();
                this.achievements = data.achievements;
                console.log("Achievements fetched");
                console.log(data);
                this.loading = false;
            } catch (error) {
                this.error = error instanceof Error ? error.message : 'Failed to fetch achievements';
                this.loading = false;
            }
        },


        clearAchievements() {
            this.achievements = [];
            this.loading = false;
            this.error = null;
        }
    }
});