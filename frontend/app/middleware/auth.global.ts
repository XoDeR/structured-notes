export default defineNuxtRouteMiddleware(async (to, _) => {
  if (!import.meta.client) return;
  const userAuth = localStorage.getItem('isLoggedIn');

  if (to.fullPath.startsWith('/dashboard') && !userAuth) {
    return navigateTo({
      path: '/login',
      query: { redirect: to.fullPath },
    });
  }

  // NOT IMPLEMENTED
  /*
  if (to.fullPath.startsWith('/dashboard/admin')) {
    const user = await useUserStore().fetch();
    if (user && user.role === 2) return;
    return navigateTo({
      path: '/dashboard',
    });
  }
    */

  if (to.fullPath.startsWith('/login') || to.fullPath.startsWith('/signup')) {
    if (userAuth)
      return navigateTo({
        path: '/dashboard',
      });
  }
});