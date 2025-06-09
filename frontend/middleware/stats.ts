export default defineNuxtRouteMiddleware(async (to, from) => {
    const authStore = useAuthStore();
    const statsStore = useStatsStore();
    
    if (authStore.isAuthenticated) {
      console.log('Fetching stats');
      await statsStore.fetchStats();
    }
  });