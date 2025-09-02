import { createStore } from 'vuex';
import layout from './layout';
import dashboard from './dashboard';

const store = createStore({
  modules: {
    layout,
    dashboard,
  },
});

export default store;
