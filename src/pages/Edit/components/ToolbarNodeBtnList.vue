<template>
  <div class="toolbarNodeBtnList" :class="[dir, { isDark: isDark }]">
    <div @click="$router.push({ name: 'Home' })" class="toolbarBtn">
      <span class="icon iconfont iconhoutui-shi"></span>
      <span class="text" >{{ $t('mindmap.myMaps')}}</span>
    </div>
    <template v-for="item in list" :key="item">
      <div class="toolbarBtn" :class="{
        disabled: 'disabled' in buttonConfig[item] ? buttonConfig[item].disabled() : false,
        active: 'active' in buttonConfig[item] ? buttonConfig[item].active() : false
      }" @click="emit('event' in buttonConfig ? buttonConfig.event : 'execCommand', buttonConfig[item].command)">
        <span class="icon iconfont" :class="buttonConfig[item].icon"></span>
        <span class="text">{{ $t(buttonConfig[item].text) }}</span>
      </div>
    </template>
  </div>
</template>

<script setup>

import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useStore } from 'vuex'
import bus from '@/utils/bus.js'

const props = defineProps({
  dir: {
    type: String,
    default: 'h' // h（水平排列）、v（垂直排列）
  },
  list: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['execCommand', 'startPainter', 'showNodeImage', 'showNodeLink', 'showNodeNote', 'showNodeTag', 'createAssociativeLine'])

const store = useStore()

// Реактивные переменные
const activeNodes = ref([])
const backEnd = ref(false)
const forwardEnd = ref(true)
const readonly = ref(false)
const isFullDataFile = ref(false)
const isInPainter = ref(false)

// Computed свойства
const isDark = computed(() => store.state.isDark)

const hasRoot = computed(() => {
  return activeNodes.value.findIndex(node => node.isRoot) !== -1
})

const hasGeneralization = computed(() => {
  return activeNodes.value.findIndex(node => node.isGeneralization) !== -1
})

// Методы
const onModeChange = (mode) => {
  readonly.value = mode === 'readonly'
}

const onNodeActive = (args) => {
  activeNodes.value = [...args[1]]
}

const onBackForward = (index, len) => {
  backEnd.value = index <= 0
  forwardEnd.value = index >= len - 1
}

const onPainterStart = () => {
  isInPainter.value = true
}

const onPainterEnd = () => {
  isInPainter.value = false
}

const showNodeIcon = () => {
  bus.emit('close_node_icon_toolbar')
  store.commit('setActiveSidebar', 'nodeIconSidebar')
}

const showFormula = () => {
  store.commit('setActiveSidebar', 'formulaSidebar')
}

const emitBus = (...args) => bus.emit(...args)

// Хуки жизненного цикла
onMounted(() => {
  bus.on('mode_change', onModeChange)
  bus.on('node_active', onNodeActive)
  bus.on('back_forward', onBackForward)
  bus.on('painter_start', onPainterStart)
  bus.on('painter_end', onPainterEnd)
})

onUnmounted(() => {
  bus.off('mode_change', onModeChange)
  bus.off('node_active', onNodeActive)
  bus.off('back_forward', onBackForward)
  bus.off('painter_start', onPainterStart)
  bus.off('painter_end', onPainterEnd)
})

// Конфигурация кнопок
const buttonConfig = {
  back: {
    icon: 'iconhoutui-shi',
    text: 'toolbar.undo',
    command: 'BACK',
    disabled: () => readonly.value || backEnd.value
  },
  forward: {
    icon: 'iconqianjin1',
    text: 'toolbar.redo',
    command: 'FORWARD',
    disabled: () => readonly.value || forwardEnd.value
  },
  painter: {
    icon: 'iconjiedian',
    text: 'toolbar.painter',
    event: 'startPainter',
    disabled: () => activeNodes.value.length <= 0 || hasGeneralization.value,
    active: () => isInPainter.value
  },
  siblingNode: {
    icon: 'iconjiedian',
    text: 'toolbar.insertSiblingNode',
    command: 'INSERT_NODE',
    disabled: () => activeNodes.value.length <= 0 || hasRoot.value || hasGeneralization.value
  },
  childNode: {
    icon: 'icontianjiazijiedian',
    text: 'toolbar.insertChildNode',
    command: 'INSERT_CHILD_NODE',
    disabled: () => activeNodes.value.length <= 0 || hasGeneralization.value
  },
  deleteNode: {
    icon: 'iconshanchu',
    text: 'toolbar.deleteNode',
    command: 'REMOVE_NODE',
    disabled: () => activeNodes.value.length <= 0
  },
  image: {
    icon: 'iconimage',
    text: 'toolbar.image',
    event: 'showNodeImage',
    disabled: () => activeNodes.value.length <= 0
  },
  icon: {
    icon: 'iconxiaolian',
    text: 'toolbar.icon',
    action: showNodeIcon,
    disabled: () => activeNodes.value.length <= 0
  },
  link: {
    icon: 'iconchaolianjie',
    text: 'toolbar.link',
    event: 'showNodeLink',
    disabled: () => activeNodes.value.length <= 0
  },
  note: {
    icon: 'iconflow-Mark',
    text: 'toolbar.note',
    event: 'showNodeNote',
    disabled: () => activeNodes.value.length <= 0
  },
  tag: {
    icon: 'iconbiaoqian',
    text: 'toolbar.tag',
    event: 'showNodeTag',
    disabled: () => activeNodes.value.length <= 0
  },
  summary: {
    icon: 'icongaikuozonglan',
    text: 'toolbar.summary',
    command: 'ADD_GENERALIZATION',
    disabled: () => activeNodes.value.length <= 0 || hasRoot.value || hasGeneralization.value
  },
  associativeLine: {
    icon: 'iconlianjiexian',
    text: 'toolbar.associativeLine',
    event: 'createAssociativeLine',
    disabled: () => activeNodes.value.length <= 0 || hasGeneralization.value
  },
  formula: {
    icon: 'icongongshi',
    text: 'toolbar.formula',
    action: showFormula,
    disabled: () => activeNodes.value.length <= 0 || hasGeneralization.value
  }
}

</script>

<style lang="less" scoped>
.toolbarNodeBtnList {
  display: flex;

  &.isDark {
    .toolbarBtn {
      color: hsla(0, 0%, 100%, 0.9);

      .icon {
        background: transparent;
        border-color: transparent;
      }

      &:hover {
        &:not(.disabled) {
          .icon {
            background: hsla(0, 0%, 100%, 0.05);
          }
        }
      }

      &.disabled {
        color: #54595f;
      }
    }
  }

  .toolbarBtn {
    display: flex;
    justify-content: center;
    flex-direction: column;
    cursor: pointer;
    margin-right: 20px;

    &:last-of-type {
      margin-right: 0;
    }

    &:hover {
      &:not(.disabled) {
        .icon {
          background: #f5f5f5;
        }
      }
    }

    &.active {
      .icon {
        background: #f5f5f5;
      }
    }

    &.disabled {
      color: #bcbcbc;
      cursor: not-allowed;
      pointer-events: none;
    }

    .icon {
      display: flex;
      height: 26px;
      background: #fff;
      border-radius: 4px;
      border: 1px solid #e9e9e9;
      justify-content: center;
      flex-direction: column;
      text-align: center;
      padding: 0 5px;
    }

    .text {
      margin-top: 3px;
    }
  }

  &.v {
    display: block;
    width: 120px;
    flex-wrap: wrap;

    .toolbarBtn {
      flex-direction: row;
      justify-content: flex-start;
      margin-bottom: 10px;
      width: 100%;
      margin-right: 0;

      &:last-of-type {
        margin-bottom: 0;
      }

      .icon {
        margin-right: 10px;
      }

      .text {
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }
    }
  }
}
</style>
