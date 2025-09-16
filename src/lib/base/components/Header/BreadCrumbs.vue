<template>
  <el-breadcrumb separator="/">
    <el-breadcrumb-item v-for="(item, i) in breadCrumb" :key="i">
      <template v-if="item.meta?.excluded != true">
        <router-link :to="{ name: item.name, params: $route.params }">
          {{ item.meta.title }}
        </router-link>
      </template>
    </el-breadcrumb-item>
  </el-breadcrumb>
</template>
<style lang="scss">
.el-breadcrumb {
  font-size: 16px;
}
</style>
<script>
import { useRouter, useRoute } from 'vue-router';
import { ref, onMounted, computed ,getCurrentInstance} from 'vue';

export default {
  setup() {
    const router = useRouter();
    const route = useRoute();
    const { proxy } = getCurrentInstance();
    const t = proxy.$t;

    const parent = ref([]);

    const breadCrumb = computed(() => {
      let breadCrumbs = [];
      parent.value = [];

      let arr = getParent(router.currentRoute.value);

      if (arr.length) {
        breadCrumbs.unshift(...arr);
      }

      return breadCrumbs;
    });

    const getParent = curRoute => {
      return curRoute.matched;
    };

    return {
      breadCrumb,
    };
  }
};
</script>
