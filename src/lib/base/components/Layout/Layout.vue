<template>
  <div
    id="base"
    :class="[{ root: true, sidebarClose, sidebarStatic }, 'sing-dashboard', 'sidebar-' + sidebarColorName, 'sidebar-' + sidebarType, 'navbar-' + navbarColorName]">
    <Sidebar />
    <div class="wrap">
      <Header />
      <transition name="router-animation">
       
      </transition>
      <div class="contentbar" :class="$route.name !== 'constructor' ? 'mt-3 ml-4 mr-4': ''">

        <div  v-if="$route.meta.title && $route.name !== 'constructor'" class="row">
        <div class="mt-3 col-md-12">
          <h2>{{ $t($route.meta.title) }}</h2>
        </div>
        </div>
        <router-view :class="$route.name !== 'constructor' ? 'mt-4': ''" />
      </div>
      <footer></footer>
    </div>

  </div>
</template>

<script setup>
import { computed, ref, onMounted, onBeforeUnmount } from 'vue';
import { useStore } from 'vuex';
import Sidebar from '../Sidebar/Sidebar';
import Header from '../Header/Header';

const store = useStore();

const sidebarClose = computed(() => store.state.layout.sidebarClose);
const sidebarStatic = computed(() => store.state.layout.sidebarStatic);
const sidebarColorName = computed(() => store.state.layout.sidebarColorName);
const navbarColorName = computed(() => store.state.layout.navbarColorName);
const sidebarType = computed(() => store.state.layout.sidebarType);
const helperOpened = computed(() => store.state.layout.helperOpened);

const switchSidebar = () => store.dispatch('layout/switchSidebar');
const changeSidebarActive = (value) => store.dispatch('layout/changeSidebarActive', value);
const toggleSidebar = () => store.dispatch('layout/toggleSidebar');

const handleWindowResize = () => {
  const width = window.innerWidth;
  if (width <= 768 && sidebarStatic.value) {
    toggleSidebar();
    changeSidebarActive(null);
  }
};

onMounted(() => {
  const staticSidebar = localStorage.getItem('sidebarStatic') ? JSON.parse(localStorage.getItem('sidebarStatic')) : false;
  if (staticSidebar) {
    store.state.layout.sidebarStatic = true;
  } else if (!sidebarClose.value) {
    setTimeout(() => {
      switchSidebar();
      changeSidebarActive(null);
    }, 2500);
  }
  
  handleWindowResize();
  window.addEventListener('resize', handleWindowResize);
});

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleWindowResize);
});
</script>

<style scoped lang="scss">
@import '../../styles/app.scss';

.contentbar {
  //margin-top: 75px;
  //padding: 30px;
  //margin-bottom: 30px;
}
.border-left-primary {
  border-left: 0.25rem solid #4e73df !important;
}

.typcn {
  color: rgb(143, 143, 143);
}

.root {
  height: 100%;
  position: relative;
  left: 0;
  transition: left $sidebar-transition-time ease-in-out;
}

.wrap {
  position: relative;
  min-height: 100%;
  display: flex;
  margin-left: 64px;
  flex-direction: column;
  left: $sidebar-width-open - $sidebar-width-closed;
  right: 0;
  transition: left $sidebar-transition-time ease-in-out, margin-left $sidebar-transition-time ease-in-out;

  @media (max-width: breakpoint-max(sm)) {
    margin-left: 0;
    left: $sidebar-width-open;
  }
}

.sidebarClose div.wrap {
  left: 0;
}

.sidebarStatic .wrap {
  //transition: none;
  left: 0;
  margin-left: $sidebar-width-open;
}

.content {
  position: relative;
  display: flex;
  flex-direction: column;
  flex-grow: 1;
  padding: $content-padding $content-padding (
    $content-padding + 20px
  );
background-color: $body-bg;

@media (max-width: breakpoint-max(sm)) {
  padding: 20px $content-padding-sm ($content-padding + 30px);
}

// hammers disallows text selection, allowing it for large screens
@media (min-width: breakpoint-min(sm)) {
  user-select: auto !important;
}
}

.contentFooter {
  position: absolute;
  bottom: 15px;
  color: $text-muted;
}
</style>
