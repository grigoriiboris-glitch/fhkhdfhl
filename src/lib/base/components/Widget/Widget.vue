<template>
  <div :class="['widget', className, { collapsed: state === 'collapse', fullscreened: state === 'fullscreen', loading: fetchingData }]">
    <div v-if="title && !customHeader" class="card-header">{{ title }}</div>
    <div v-if="title && customHeader" class="card-header" v-html="title"></div>
    <div class="widget-controls d-flex justify-content-end p-2">
      <button v-if="settings" class="btn btn-light btn-sm"><i class="la la-cog"></i></button>
      <button v-if="refresh" @click="loadWidgster" class="btn btn-light btn-sm">
        <i class="la la-refresh"></i>
      </button>
      <button v-if="fullscreen && state !== 'fullscreen'" @click="changeState('fullscreen')" class="btn btn-light btn-sm">
        <i class="glyphicon glyphicon-resize-full"></i>
      </button>
      <button v-if="fullscreen && state === 'fullscreen'" @click="changeState('default')" class="btn btn-light btn-sm">
        <i class="glyphicon glyphicon-resize-small"></i>
      </button>
      <button v-if="collapse && state !== 'collapse'" @click="changeState('collapse')" class="btn btn-light btn-sm">
        <i class="la la-angle-down"></i>
      </button>
      <button v-if="collapse && state === 'collapse'" @click="changeState('default')" class="btn btn-light btn-sm">
        <i class="la la-angle-up"></i>
      </button>
      <button v-if="close" @click="closeWidget" class="btn btn-light btn-sm">
        <i class="la la-remove"></i>
      </button>
    </div>
    <div class="widget-body card-body" v-show="state !== 'collapse'">
      <div v-if="fetchingData && showLoader" class="text-center">
        <span class="spinner-border spinner-border-sm"></span>
      </div>
      <slot v-else></slot>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Widget',
  props: {
    title: String,
    className: String,
    customHeader: Boolean,
    settings: Boolean,
    refresh: Boolean,
    fullscreen: Boolean,
    collapse: Boolean,
    close: Boolean,
    fetchingData: Boolean,
    showLoader: Boolean
  },
  data() {
    return {
      state: 'default'
    };
  },
  methods: {
    changeState(newState) {
      this.state = newState;
    },
    loadWidgster() {
      this.$emit('load');
    },
    closeWidget() {
      this.$emit('close');
    }
  }
};
</script>

<style lang="scss" >
@import '../../styles/app';

.title {
  margin-top: 0;
  color: $widget-title-color;

  @include clearfix();
}

:global .h-100 {
  height: 100%;
}

.widget {
  display: block;
  position: relative;
  margin-bottom: $grid-gutter-width;
  padding: $widget-padding-vertical $widget-padding-horizontal;
  background: $widget-bg-color;
  border-radius: $border-radius-sm;
  box-shadow: var(--widget-shadow);

  &.loading {
    min-height: 150px;
  }

  [control='collapse'] {
    display: unset;
  }

  [control='expand'] {
    display: none;
  }

  [control='fullscreen'] {
    display: unset;
  }

  [control='restore'] {
    display: none;
  }

  &.collapsed {
    min-height: unset;

    [control='collapse'] {
      display: none;
    }

    [control='expand'] {
      display: unset;
    }
  }

  &.fullscreened {
    position: fixed;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
    margin: 0;
    z-index: 10000;

    [control='fullscreen'] {
      display: none;
    }

    [control='restore'] {
      display: unset;
    }
  }

  > header {
    margin: (-$widget-padding-vertical) (-$widget-padding-horizontal);
    padding: $widget-padding-vertical $widget-padding-horizontal;

    h1,
    h2,
    h3,
    h4,
    h5,
    h6 {
      margin: 0;
    }
  }

  :global {
    .loader {
      position: absolute;
      top: 0;
      bottom: 0;
      left: 0;
      right: 0;

      .spinner {
        position: absolute;
        top: 50%;
        width: 100%; //ie fix
        margin-top: -10px;
        font-size: 20px;
        text-align: center;
      }
    }

    .widget-body.p-0 {
      margin: $widget-padding-vertical (-$widget-padding-horizontal) (-$widget-padding-vertical);

      + footer {
        margin-top: $widget-padding-vertical;
      }
    }
  }

  &:global.bg-transparent {
    box-shadow: none;
  }
}

.widgetBody {
  @include clearfix();

  > footer {
    margin: $spacer/2 (-$widget-padding-horizontal) (-$widget-padding-vertical);
    padding: 10px 20px;
  }
}

.widgetControls + .widgetBody {
  margin-top: $widget-padding-vertical;
}

.widgetControls,
:global(.widget-controls) {
  position: absolute;
  z-index: 1;
  top: 0;
  right: 0;
  padding: 14px;
  font-size: $font-size-sm;

  a {
    padding: 1px 4px;
    border-radius: 4px;
    color: rgba($black, 0.4);

    @include transition(color 0.15s ease-in-out);

    &:hover {
      color: rgba($black, 0.1);
      text-decoration: none;
    }

    .la {
      position: relative;
      top: 2px;
    }

    .glyphicon {
      font-size: 0.7rem;
    }
  }
}

.inverse {
  top: 2px;
  position: relative;
  margin-left: 3px;

  :global {
    .glyphicon {
      vertical-align: baseline;
    }
  }
}

:global {
  .widget-image {
    position: relative;
    overflow: hidden;
    margin: (-$widget-padding-vertical) (-$widget-padding-horizontal);
    border-radius: $border-radius;

    > img {
      max-width: 100%;
      border-radius: $border-radius $border-radius 0 0;
      transition: transform 0.15s ease;
    }

    &:hover > img {
      transform: scale(1.1, 1.1);
    }

    .title {
      position: absolute;
      top: 0;
      left: 0;
      margin: 20px;
    }

    .info {
      position: absolute;
      top: 0;
      right: 0;
      margin: 20px;
    }
  }

  .widget-footer-bottom {
    position: absolute;
    bottom: 0;
    width: 100%;
  }

  .widget-sm {
    height: 230px;
  }

  .widget-md {
    height: 373px;
  }

  .widget-padding-md {
    padding: $widget-padding-vertical $widget-padding-horizontal;
  }

  .widget-padding-lg {
    padding: $widget-padding-vertical*2 $widget-padding-horizontal*2;
  }

  .widget-body-container {
    position: relative;
    height: 100%;
  }

  .widget-top-overflow,
  .widget-middle-overflow {
    position: relative;
    margin: 0 (-$widget-padding-horizontal);

    > img {
      max-width: 100%;
    }
  }

  .widget-top-overflow {
    margin-top: (-$widget-padding-vertical);
    border-top-left-radius: $border-radius;
    border-top-right-radius: $border-radius;
    overflow: hidden;

    > img {
      border-top-left-radius: $border-radius;
      border-top-right-radius: $border-radius;
    }

    > .btn-toolbar {
      position: absolute;
      top: 0;
      right: 0;
      z-index: 1;
      margin-right: $widget-padding-horizontal;

      @include media-breakpoint-up(md) {
        top: auto;
        bottom: 0;
      }
    }
  }

  .widget-icon {
    opacity: 0.5;
    font-size: 42px;
    height: 60px;
    line-height: 45px;
    display: inline-block;
  }
}

.widget-loader {
  position: absolute;
  top: 0;
  left: 0;
}

</style>
