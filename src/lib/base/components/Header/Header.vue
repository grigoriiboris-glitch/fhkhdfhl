<template>
  <nav :class="[navbarTypeClass, 'header-' + navbarColorScheme]"
    class="navbar app-header d-print-none navbar-light navbar-expand navbar-static-type header-light">
    <div class="container-fluid">
      <button class="navbar-toggler d-block" type="button" @click="toggleSidebar">
        <img :src="Menu" alt="menu" />
      </button>

      <router-link class="navbar-brand d-md-none" :to="{ name: 'Home' }">
        <i class="fa fa-circle text-primary me-1"></i>
        <i class="fa fa-circle text-danger"></i>
        sing
        <i class="fa fa-circle text-danger me-1"></i>
        <i class="fa fa-circle text-primary"></i>
      </router-link>

      <div class="collapse navbar-collapse ml-4 justify-content-between">
        <ul class="navbar-nav me-auto">
          <li class="nav-item">
            <bread-crumbs></bread-crumbs>
          </li>
        </ul>
        <ul v-if="1 == 2" class="nav">
          <!-- <form class="navbar-form d-none d-lg-block" role="search">
            <div class="form-group">
              <div class="input-group input-group-no-border ml-4">
                <input class="form-control" id="main-search" type="text" placeholder="Search Dashboard">
                <span class="input-group-append">
                  <span class="input-group-text">
                    <i class="la la-search"></i>
                  </span>
                </span>
              </div>
            </div>
          </form> -->
        </ul>
        <ul class="navbar-nav ms-auto"   @mouseleave="dropdown = false"  @mouseenter="dropdown = true">
          <li class="nav-item b-nav-dropdown dropdown notificationsMenu d-md-down-none ms-2">
            <a class="nav-link dropdown-toggle align-items-center d-flex"
              @click="dropdownNotification = !dropdownNotification">
              <span class="avatar rounded-circle mr-2">
                <!-- <img v-if="user.avatar" :src="user.avatar" class="rounded-circle" alt="user" /> -->
                <span>{{ firstUserLetter }}</span>
              </span>
              <span>{{ User.name }}</span>
              <!-- <span class="badge bg-danger ms-2">9</span> -->

              <!-- <img class="px-2 dropdown-arrow" :class="{ active: dropdownNotification }" :src="CaretDown"
                alt="caretDown" /> -->
            </a>
            <!-- <ul :class="{ show:dropdownNotification}" class="dropdown-menu notificationsWrapper py-0 animate__animated animate__animated-fast animate__fadeIn dropdown-menu-right ">

              <Notifications />
            </ul> -->
          </li>

          <li @click.stop="dropdown = !dropdown" class="nav-item b-nav-dropdown dropdown settingsDropdown d-sm-down-none align-items-center d-flex">
            <a role="button" aria-haspopup="true" aria-expanded="false"
              class="nav-link dropdown-toggle dropdown-toggle-no-caret ">
              <img :src="Settings" alt="settings" class="px-2"></a>

            <ul @mouseleave="dropdown = false" v-if="User" :class="dropdown ? 'show' : ''" class="dropdown-menu dropdown-menu-right ">
              <router-link class="dropdown-item" :to="{
                name: 'profile',
                params: { user_id: 1/*User.id*/ }
              }"><img :src="Userimg" class="me-2" />{{ $t('My profile') }}</router-link>

              <div class="dropdown-divider"></div>
              <router-link class="dropdown-item" :to="{
                name: 'tariffs'
              }"><img :src="Document" class="me-2" />{{ $t('Pricing plans') }}</router-link>

              <!-- <router-link class="dropdown-item" :to="{
                name: 'tariffs'
              }"><img :src="Document" class="me-2" />{{ 'Кейсы интеграций' }}</router-link> -->

              <!-- <a class="dropdown-item" href="/inbox">
                <img :src="Envelope" class="me-2" /> Inbox
                <span class="badge bg-dark ms-2">9</span>
              </a> -->
              <div class="dropdown-divider"></div>
              <button class="dropdown-item" @click="logout">
                <img :src="Cancel" class="me-2" /> {{ $t('exit') }}
              </button>
            </ul>
          </li>
        </ul>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { computed, ref, onMounted } from 'vue';
import { useStore } from 'vuex';
import useAuth from '../../composition/useAuth';

import Menu from '../../assets/sidebar/Fill/Menu.svg';
import Exchange from '../../assets/sidebar/Fill/Exchange.svg';
import Cross from '../../assets/sidebar/Fill/Cross.svg';
import Search from '../../assets/sidebar/Fill/Search.svg';
import CaretDown from '../../assets/sidebar/Fill/Caret down.svg';
import Settings from '../../assets/sidebar/Outline/Settings-alt.svg';
import Userimg from '../../assets/sidebar/Outline/User.svg';
import Document from '../../assets/sidebar/Outline/Document.svg';
import Envelope from '../../assets/sidebar/Outline/Envelope.svg';
import Cancel from '../../assets/sidebar/Outline/Cancel.svg';
import Notifications from '../Notifications/Notifications';
import BreadCrumbs from './BreadCrumbs.vue';



const store = useStore();

const sidebarClose = computed(() => store.state.layout.sidebarClose);
const sidebarStatic = computed(() => store.state.layout.sidebarStatic);
const navbarColorScheme = computed(() => store.state.layout.navbarColorScheme);

const switchSidebar = () => store.dispatch('layout/switchSidebar');
const changeSidebarActive = (value) => store.dispatch('layout/changeSidebarActive', value);
const toggleSidebar = () => store.dispatch('layout/toggleSidebar');
const logoutUser = () => store.dispatch('auth/logoutUser');

const dropdown = ref(false);
const dropdownNotification = ref(false);
const firstUserLetter = computed(() => (User.value.name || User.value.email || 'P')[0].toUpperCase());
const navbarTypeClass = computed(() => `navbar-${store.state.layout.navbarType}-type`);

const handleWindowResize = () => {
  const width = window.innerWidth;
  if (width <= 768 && sidebarStatic.value) {
    toggleSidebar();
    changeSidebarActive(null);
  }
};

const switchSidebarMethod = () => {
  if (!sidebarClose.value) {
    switchSidebar();
    changeSidebarActive(null);
  } else {
    switchSidebar();
    const paths = window.location.pathname.split('/');
    paths.pop();
    changeSidebarActive(paths.join('/'));
  }
};

const toggleSidebarMethod = () => {
  if (sidebarStatic.value) {
    toggleSidebar();
    changeSidebarActive(null);
  } else {
    toggleSidebar();
    const paths = window.location.pathname.split('/');
    paths.pop();
    changeSidebarActive(paths.join('/'));
  }
};


const { User, logout } = useAuth();

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

</script>


<style lang="scss">
@import '../../styles/app';

.app-header.navbar {
  height: $navbar-height;
  z-index: 100;
  padding: 0 1.85rem 0;
  background: var(--navbar-bg);
  box-shadow: var(--navbar-shadow);
  transition: $transition-base;
  border: none;
  font-weight: 500;
  justify-content: flex-start;

  .dropdown-arrow {
    transition: transform .3s;

    &.active {
      transform: rotate(180deg);
    }
  }

  .nav {
    height: 100%;
    padding: 0;

    .nav-item {

      .nav-link,
      .nav-link>a {
        display: flex;
        align-items: center;
        height: 100%;
        position: relative;
        padding: 0.5rem;
        color: $navbar-link-color;

        @include hover {
          color: $navbar-link-hover-color;
          background: $navbar-link-hover-bg;
          text-decoration: none;
        }

        .la {
          font-size: 20px;

          @include media-breakpoint-down(sm) {
            font-size: 27px;
          }
        }
      }
    }
  }

  &.header-dark {
    .nav {
      .nav-item {

        .nav-link,
        .nav-link>a {
          color: $white;

          @include hover {
            color: rgba($white, 0.7);
          }
        }
      }
    }

    .input-group-text {
      color: $white;
    }

    input::placeholder {
      color: rgba($white, 0.7);
    }

    .navbarBrand {
      color: $white;
    }
  }

  @include media-breakpoint-down(md) {
    padding: 0 $spacer/2;
  }

  &.navbar-floating-type {
    margin: $spacer $content-padding 0;
    border-radius: $border-radius;

    @media (max-width: breakpoint-max(sm)) {
      margin-left: $content-padding-sm;
      margin-right: $content-padding-sm;
    }
  }

  .form-group {
    width: 300px;

    input {
      background-color: #F9FAFE;
      border-radius: 0;

      &::placeholder {
        color: #4A5056;
        font-weight: $font-weight-normal;
      }
    }

    .input-group-prepend {
      margin-right: 0;
    }

    .input-group-text {
      border: none;
      background-color: #F9FAFE !important;
      transition: background-color ease-in-out 0.15s;
      border-radius: 0;
    }
  }

  .avatar {
    display: flex;
    align-items: center;
    justify-content: center;
    overflow: hidden;
    height: 40px;
    width: 40px;
    background: #b8e5ff;
    font-weight: 600;
    font-size: 18px;
  }

  .avatar-badge {
    background-color: #0D2B47;
  }

  .navbarBrand {
    position: absolute;
    left: 0;
    right: 0;
    top: 0;
    bottom: 0;
    display: flex;
    justify-content: center;
    align-items: center;
    font-weight: 700;
    font-size: 1.25rem;
    pointer-events: none;

    i {
      font-size: 10px;
    }
  }

  .notificationsMenu {
    .dropdown-menu {
      left: auto !important;
      right: 0 !important; 
      top: $navbar-height !important;
    }
  }

  .dropdown-toggle::after {
    display: none;
  }

  .settingsDropdown {
.dropdown-menu {
  top:55px;
}
    .dropdown-item:focus {
      outline: none;
    }

    ul {
      width: 15rem;

      li {
        padding: 4px 0;
      }
    }
  }

  .headerSvgFlipColor {
    color: var(--navbar-icon-bg) !important;

    :global {
      .bg-primary {
        background-color: var(--navbar-icon-bg) !important;
      }
    }
  }
}
</style>
