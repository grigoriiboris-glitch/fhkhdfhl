<template>
  <li v-if="childrenLinks && childrenLinks.length" :class="{ headerLink: true, className }">
    <div @click="() => togglePanelCollapse(link)">
      <a class="sidebar-link" :class="{ 'router-link-exact-active router-link-active': isActive }">
        <span class="icon">
          <icon :name="iconName" :size="24" />
        </span>
        <div class="d-flex  align-items-center">
          <div>{{ header }} </div>
        </div>
        <div v-if="childrenLinks.length"
          :class="{ caretWrapper: true, carretActive: isActive && headerLinkWasClicked }">
          <i class="ml-5 fa fa-angle-right" />
        </div>
      </a>
    </div>
    <!-- <transition name="slide-fade"> -->
    <div id="collapsemain" style="" :class="{ show: isActive && headerLinkWasClicked }" class="collapse">
      <ul class="sub-menu">
        <NavLink v-for="childrenLink in childrenLinks" :activeItem="activeItem" :header="childrenLink.title"
          :index="childrenLink.to" :link="childrenLink.to" :childrenLinks="childrenLink.childrenLinks"
          :key="childrenLink.to" 
          :routes="routes"
          />
      </ul>
    </div>
    <!-- </transition> -->
  </li>
  <li v-else class="headerLink">
    <router-link :to="{ name: link, params: getRouteParams(link) }" class="sidebar-link"
      :class="{ 'router-link-exact-active router-link-active': isActive }">
      <span class="icon mr-2">
        <icon :name="iconName" :size="24" />
      </span>
      <div class="d-flex  align-items-center">
        <div>{{ header }} </div>
      </div> <sup v-if="label" :class="'text-' + labelColor" class="headerLabel">{{ label }}</sup>
      <span v-if="badge" class="badge rounded-pill bg-danger">{{ badge }}</span>
    </router-link>
  </li>
</template>

<script setup>
import { computed, ref } from 'vue';
import { useStore } from 'vuex';
import { useRoute } from 'vue-router';

const route = useRoute();
const store = useStore();

const props = defineProps({
  badge: { type: String, default: '' },
  header: { type: String, default: '' },
  iconName: { type: String, default: '' },
  iconImg: { type: String, default: '' },
  headerLink: { type: String, default: '' },
  link: { type: String, default: '' },
  childrenLinks: { type: Array, default: null },
  className: { type: String, default: '' },
  isHeader: { type: Boolean, default: false },
  deep: { type: Number, default: 0 },
  activeItem: {},
  label: { type: String },
  labelColor: { type: String, default: 'warning' },
  index: { type: String },
  routes: {type: Array, required:true}
});

const getRouteParams = (name) => {
  let params = route.params;
  let data = {};
  props.routes.forEach(el => {
    if (el.name === name) {
      Object.keys(params).forEach(param => {
        if (el.path.includes(`:${param}`)) {
          data[param] = params[param];
        }
      });
    }
  })

  return data;
}

const headerLinkWasClicked = ref(false);

const togglePanelCollapse = (link) => {
  if (props.childrenLinks && !props.childrenLinks.length) {
    return;
  }
  if (!props.activeItem !== link) {
    store.dispatch('layout/changeSidebarActive', link);
  } else {
    store.dispatch('layout/changeSidebarActive', null);

  }
  log(link, props.activeItem, props.index)
  headerLinkWasClicked.value = !headerLinkWasClicked.value;
};

const fullIconName = computed(() => `fi ${props.iconName}`);
const isActive = computed(() => {
  if (props.link != '' && props.activeItem === props.link) {
    return true;
  }
  // if (!props.activeItem) {
  //   return false
  // }
  // props.activeItem.forEach(el => {
  //   if (el.to === props.index) {
  //     return true;
  //   }
  // });

  // return false;

  //&& headerLinkWasClicked.value
});
</script>

<style lang="scss" scoped >
@import '../../../styles/app';

.headerLink {
  width: 100%;
  overflow-x: hidden;

  @media (min-width: map_get($grid-breakpoints, lg))
  and (min-height: $screen-lg-height), (max-width: map_get($grid-breakpoints, md) - 1px) {
    font-size: 13px;
  }

  a {
    display: block;
    color: var(--sidebar-color);
    text-decoration: none;
    cursor: pointer;
  }

  &:last-child > a {
    border-bottom: 1px solid $sidebar-item-border-color;
  }

  > a,
  > div a {
    align-items: center;
    position: relative;
    padding-left: 64px;
    line-height: 58px;
    border-top: 1px solid $sidebar-item-border-color;
    font-size: 14px;
    font-weight: $font-weight-normal;

    &:hover {
      background-color: var(--sidebar-item-hover-bg-color);
    }

    > i {
      margin-right: 7px;
    }
  }

  .icon {
    font-size: $font-size-larger;
    display: flex;
    justify-content: center;
    align-items: center;
    position: absolute;
    top: 10px;
    left: 16px;
    width: 32px;
    height: 32px;
    line-height: 28px;
    text-align: center;

    @media (min-width: map_get($grid-breakpoints, lg))
    and (min-height: $screen-lg-height), (max-width: map_get($grid-breakpoints, md) - 1px) {
      top: 12px;
    }
  }

  .badge {
    float: right;
    line-height: 10px;
    margin-top: 20px;
    margin-right: 20px;
    font-size: 0.875em;
    background: var(--sidebar-badge-bg);

    @media (min-width: map_get($grid-breakpoints, lg)) and (min-height: $screen-lg-height), (max-width: map_get($grid-breakpoints, md) - 1px) {
      margin-top: 16px;
    }
  }
  .collapse > ul > li a {
    padding-left: 50px;
  }
  #collapsemenu > ul > li a {
    padding-bottom: 14px;
  }
}

.headerLabel {
  font-weight: 600;
}

.caretWrapper {
  display: flex;
  align-items: center;
  margin-left: auto;
  margin-right: 22px;
  color: var(--sidebar-color);

  i {
    @include transition(transform 0.3s ease-in-out);
  }
}

.caretWrapper i {
  font-size: $font-size-larger;
}
.carretActive i {
  transform: rotate(90deg);
}

div a.router-link-active {
  color: var(--sidebar-item-active);
  font-size: $font-size-larger;
  font-weight: $font-weight-semi-bold;

  .icon {
    border-radius: 50%;
    background-color: var(--sidebar-item-active-bg);

    i {
      color: var(--sidebar-icon-active);
    }
  }
}

.collapse,
.collapsing {
  border: none;
  box-shadow: none;
  margin: 0;
  border-radius: 0;

  a {
    line-height: 20px !important;

    &.router-link-active {
      font-weight: $font-weight-semi-bold;
      color: var(--sidebar-item-active);
    }
  }

  ul {
    background: var(--sidebar-action-bg);
    padding: $spacer;

    li {
      list-style: none;
    }

    a {
      padding: 10px 20px 10px 26px;
      font-size: $font-size-mini;

      &:hover {
        background-color: var(--sidebar-item-hover-bg-color);
      }
    }
  }
}

.sidebar-link{
  position: relative;
 .fa-angle-right{
  position: absolute;
  right: 20px;
  top: 21px;
 }
 div{
  font-size: 16px;
 }

}

</style>
