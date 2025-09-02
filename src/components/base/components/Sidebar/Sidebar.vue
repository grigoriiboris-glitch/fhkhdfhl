<template>
  <div class="sidebar-wrapper ">
    <nav class="d-flex flex-column" :class="{ sidebar: true, sidebarStatic, sidebarOpened }"
      @mouseenter="sidebarMouseEnter" @mouseleave="sidebarMouseLeave">
      <header class="logo">
        <router-link :to="{ name: 'chat-bots' }"><span class="primary-word">constructor</span> <span
            class="secondary-word">
            Bot </span></router-link>
      </header>
      <ul class="nav" v-if="currentItems.length">
        <NavLink v-for="link in currentItems" :key="link.title" :activeItem="activeItem" :header="link.title"
          :link="link.to" :iconName="`typcn ${link.icon}`" :index="link.to"
          :childrenLinks="link && 'items' in link ? link.items : []" :routes="allRoutes" />

      </ul>
      <template v-if="User.isSuperAdmin()">
        <ul class="nav" v-if="adminItems.length">
          <h4 class="navTitle">
            {{ $t('Management') }}
          </h4>
          <NavLink v-for="link in adminItems" :key="link.title" :activeItem="activeItem" :header="link.title"
            :link="link.to" :iconName="`typcn ${link.icon}`" :index="link.to"
            :childrenLinks="link && 'items' in link ? link.items : []" :routes="allRoutes" />
        </ul>
      </template>


      <!-- <h5 class="navTitle">
        LABELS
      </h5>
      <ul class="sidebarLabels">
        <li>
          <a href="#">
            <i class="fa fa-circle text-danger"/>
            <span class="labelName">Core</span>
          </a>
        </li>
      </ul> -->

      <div class="sidebarAlerts">
        <div class="p-3">
          <!-- <div v-if="User.isAdmin()" class="flex flex-column mb-1">
                <h4>Тариф 2</h4>
                <router-link :to="{ name: 'tariffs' }"><el-button type="primary">Улучшить</el-button></router-link>
              </div> -->

        </div>
      </div>

      <div class="flex mt-auto mb-5 justify-content-start ml-4">
        <router-link :to="{ name: 'cabinets' }"><el-button type="warning">{{$t('Offices')}}</el-button></router-link>
      </div>
    </nav>

  </div>
</template>

<script setup>
import { ref, computed, defineProps, onBeforeMount, onMounted, getCurrentInstance } from 'vue';


import useAuth from '@/composition/useAuth';
import { useRouter, useRoute } from 'vue-router';

import { useStore } from 'vuex';
import NavLink from './NavLink/NavLink.vue';
import isScreen from '../../core/screenHelper';
//import useEvents from '@/composition/events/useEvents';
//import t from '@/plugins/i18n'
const { proxy } = getCurrentInstance();
const t = proxy.$t;

const router = useRouter();
const store = useStore();

const allRoutes = router.getRoutes();

const sidebarStatic = computed(() => store.state.layout.sidebarStatic);
const sidebarOpened = computed(() => !store.state.layout.sidebarClose);
const activeItem = computed(() => store.state.layout.sidebarActiveElement);

const route = useRoute();
const { User } = useAuth();

// const props = defineProps({
//   modelValue: {
//     type: Array,
//     default: []
//   }
// });


const leftbarItems = ref([
  { title: t('Pricing plans'), icon: 'typcn-spanner', to: 'tariffs', items: [], roles: [1] },
  { title: t('Integration cases'), icon: 'typcn-spanner', to: '', items: [], roles: [1] },
  {
    title: t('Instructions'),
    icon: 'typcn-spanner',
    to: '',
    items: [{ title: t('Pricing plans'), icon: 'typcn-spanner', to: '', items: [] }],
    roles: [1]
  },
  {
    title: t('Knowledge base'),
    icon: 'typcn-book',
    to: 'help',
    roles: [1, 2, 3]
  },
]);
const leftbarItemsBot = ref([
  {
    title: t('Analytics'),
    icon: 'typcn-chart-bar',
    to: 'bot-statistic',
    roles: [1, 2]
  },
  {
    title: t('Constructor'),
    icon: 'typcn-device-laptop',
    to: 'constructor',
    roles: [1, 2]
  },
  {
    title: t('Dialogues'),
    icon: 'typcn-messages',
    to: 'chat',
    roles: [1, 2, 3]
  },
  {
    title: t('Users'),
    icon: 'typcn-group',
    to: 'users',
    roles: [1, 2]
  },
  {
    title: t('Mailing'),
    icon: 'typcn-mail',
    to: 'mailings',
    roles: [1, 2]
  },

  {
    title: t('Notifications'),
    icon: 'typcn-volume',
    to: 'events',
    roles: [1, 2]
  },
  {
    title: t('Integrations'),
    icon: 'typcn-flow-switch',
    to: 'integration',
    roles: [1, 2]
  },

  {
    title: t('Knowledge base'),
    icon: 'typcn-book',
    to: 'help',
    roles: [1, 2, 3]
  },
  {
    title: t('Settings'),
    icon: 'typcn-cog-outline',
    to: 'settings',
    roles: [1, 2]
  }
]);

const adminItems = ref([

  {
    title: t('Customers'),
    admin: true,
    icon: 'typcn-credit-card',
    to: 'users_admin'
  },
  {
    title: t('Payments'),
    admin: true,
    icon: 'typcn-credit-card',
    to: 'payments_admin'
  },
  {
    title: t('Pricing plans'),
    admin: true,
    icon: 'typcn-document-text',
    to: 'tariffs_admin'
  },
  {
    title: t('Payment systems'),
    admin: true,
    icon: 'typcn-input-checked',
    to: 'payment_systems_admin'
  },
  {
    title: t('Settings'),
    admin: true,
    icon: 'typcn-cog-outline',
    to: 'appSettings'
  },
  {
    title: t('Articles'),
    admin: true,
    icon: 'typcn-book',
    to: 'article_admin'
  },
  {
    title: t('Categories'),
    admin: true,
    icon: 'typcn-book',
    to: 'category_admin'
  }

]);

const currentItems = computed(() => {
  let arr = leftbarItems.value;

  if (route.params.id) {
    arr = leftbarItemsBot.value;
  }

  return arr;
});

//const { listenCommonEvents } = useEvents();


const setActiveByRoute = () => {
  const paths = router.currentRoute.value.fullPath.split('/');
  paths.pop();
  store.dispatch('layout/changeSidebarActive', paths.join('/'));
};

const sidebarMouseEnter = () => {
  if (!sidebarStatic.value && (isScreen('lg') || isScreen('xl'))) {
    store.dispatch('layout/switchSidebar', false);
    setActiveByRoute();
  }
};

const sidebarMouseLeave = () => {
  if (!sidebarStatic.value && (isScreen('lg') || isScreen('xl'))) {
    store.dispatch('layout/switchSidebar', true);
    store.dispatch('layout/changeSidebarActive', null);
  }
};

onMounted(async () => {
  setActiveByRoute();
})
</script>

<!-- Sidebar styles should be scoped -->
<style lang="scss" scoped>
@import '../../styles/app';

.sidebar {
  position: absolute;
  width: $sidebar-width-open;
  background-color: var(--sidebar-bg-color);
  color: var(--sidebar-color);
  transition: $transition-base;
  height: 100vh;
}

.sidebar-wrapper {
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;
  overflow-y: auto;
  overflow-x: hidden;
  width: $sidebar-width-open;
  transition: width $sidebar-transition-time ease-in-out, left $sidebar-transition-time ease-in-out;

  @include sidebar-scrollbar();
}

.sidebar-transparent .sidebar-wrapper:hover::-webkit-scrollbar-thumb {
  background-color: transparent;
}

.sidebarClose .sidebar-wrapper {
  width: $sidebar-width-closed;

  @include media-breakpoint-down(sm) {
    width: 0;
  }
}

.sidebarStatic .sidebar-wrapper {
  width: $sidebar-width-open;
}

.logo {
  margin: 15px 0;
  font-size: 18px;
  width: 100%;
  font-weight: $font-weight-normal;
  text-align: center;
  text-shadow: 4px 4px 15px rgba(92, 176, 255, 0.6);
  transition: width $sidebar-transition-time ease-in-out;
  text-transform: uppercase;

  a {
    color: var(--logo-color);
    padding: 0 9px;
    text-decoration: none;
    white-space: nowrap;
  }

  .secondary-word {
    opacity: 1;
    transition: all $sidebar-transition-time;
  }

  .primary-word {
    color: var(--logo-color);
    font-weight: $font-weight-normal;
  }
}

.generator-link {
  display: block;
  color: var(--sidebar-color-item) !important;
  text-decoration: none;
  cursor: pointer;
  font-size: 16px !important;
}

.sidebarClose .logo {
  width: $sidebar-width-closed;

  .secondary-word {
    opacity: 0;
  }
}

.sidebarStatic .logo {
  width: 100%;
  transition: none;

  .secondary-word {
    opacity: 1;
  }
}

.nav {
  width: 100%;
  padding-bottom: 10px;
  overflow-y: auto;
  overflow-x: hidden;
}

.navTitle {
  margin: 35px 0 24px 22px;
  font-size: $font-size-mini;
  font-weight: $font-weight-bold;
  transition: opacity $sidebar-transition-time ease-in-out;
  color: var(--sidebar-nav-title-color);

  &.first {
    margin-top: 46px;
  }

  @media (min-width: breakpoint-min(lg)) {
    opacity: 1;
  }
}

.sidebarClose .navTitle {
  opacity: 0;
}

.sidebarStatic .navTitle {
  opacity: 1;
  transition: none;
}

.actionLink {
  color: #aaa;
  float: right;
  margin-right: 22px;
  margin-top: -1px;

  .la {
    font-size: $font-size-sm;
    margin-top: 4px;
  }
}

.labelName {
  opacity: 1;
  font-size: $font-size-base;
  color: var(--sidebar-item-active);
  transition: opacity $sidebar-transition-time ease-in-out;
}

.sidebarClose .labelName {
  opacity: 0;
}

.sidebarStatic .labelName {
  transition: none;
  opacity: 1;
}

.glyphiconSm {
  font-size: 9px;
}

.sidebarLabels {
  list-style-type: none;
  padding: 11px 15px 11px 24px;

  >li+li {
    margin-top: 10px;
  }

  li>a {
    font-size: $font-size-mini;
    color: var(--sidebar-color);
    text-decoration: none;

    >i {
      font-size: 11px;
      vertical-align: 1px;
      margin-right: 20px;
      color: var(--sidebar-item-active);
      transition: margin-left $sidebar-transition-time ease-in-out;
    }
  }

  li {
    padding-bottom: 28px;
  }
}

.sidebarStatic {
  .sidebarLabels {
    list-style-type: none;
  }

  .sidebarLabels>li>a>i {
    transition: none;
    margin-left: 0;
  }
}

.sidebarAlerts {
  margin-bottom: $spacer * 2;
  transition: opacity $sidebar-transition-time ease-in-out;
  opacity: 1;
  color: var(--sidebar-item-active);
  font-size: 13px;
  font-weight: 400;
}

.sidebarClose .sidebarAlerts {
  opacity: 0;
}

.sidebarStatic .sidebarAlerts {
  opacity: 1;
  transition: none;
}

.sidebarAlert {
  background: transparent;
  margin-bottom: 0;
  padding-right: 22px;
  padding-left: 22px;
}

.sidebarProgress {
  background: var(--sidebar-progress-bg-color);
}

.groupTitle {
  margin-bottom: 15px;
}
</style>
