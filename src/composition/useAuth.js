import { reactive, computed, toRef } from "vue";
import router from "@/router";

import User from "@/models/User";
import errors from "@/consts/errors";
import randomNum from "../helpers/randomNum";
import error from "../helpers/error";

const ls = localStorage;

const state = reactive({
  //captcha check
  lang: {},
  authenticated: false,
  user: null,
  User: null,
  access_token: null,
  expiration: parseInt(ls.getItem("expiration")) || undefined,
  loading: false,
  status: null,
  roles: [
    {
      value: 1,
      label: 'admin'
    },
    {
      value: 2,
      label: 'moderator'
    },
    {
      value: 3,
      label: 'manager'
    }
  ],
  statuses: [
    {
      value: 1,
      label: 'wait'
    },
    {
      value: 2,
      label: 'active'
    },
  ]
});

const access_token = computed(() => {
  if (!state.access_token || !state.expiration) return;
  if (Date.now() >= state.expiration) return;
  return state.access_token;
});

const login = async params => {
  await axios.get('/sanctum/csrf-cookie')
  await axios.post('/api/login', params).then((res) => {
    if (res.data.token) {
      state.access_token = res.data.token;

      ls.setItem('token', res.data.token);
    }

  }).catch(() => {
    ElNotification({
      title: 'Ошибка',
      message: 'неверные данные',
      type: 'error',
    })
  })
  await me()

  router.push("/");
};

const me = async () => {
  return axios({
    method: 'get',
    url: '/api/user_'

  }).then((response) => {
    state.authenticated = true;
    state.user = response.data.user;
    state.User = new User(state.user);

    if (state.User.data.lang && state.User.data.lang !== localStorage.getItem('lang')) {
      localStorage.setItem('lang', state.User.data.lang);
      location.reload();
    }

  }).catch(() => {
    state.authenticated = false;
    state.user = null;
    state.User = null;
  })
};

const update = async (params) => {
  return axios({
    method: 'post',
    url: '/api/user/update',
    data: params

  }).then((res) => {
    state.user = res.data;
    state.User = new User(state.user);
    ElNotification({
      title: '',
      message: 'Success',
      duration: 3000,
    });
  }).catch(() => {
    ElNotification({
      title: 'Error',
      message: 'wrong',
      type: 'error',
    })
  })
};

const register = async params => {
  params.hash = randomNum.generateObfuscatedString(randomNum.getNumbersForCurrentMinutesList());
  return axios({
    method: 'post',
    url: '/api/register',
    data: params

  }).then(({ data }) => {

    log(router);
    if (data.success) {
      router.push("/login").then(() => {
        ElNotification({
          title: 'Регистрация прошла успешно',
          active: true,
          message: 'please, confirm your mail',
          type: 'success',
          duration: 0,
        });
      });
    }
  }).catch((e) => {
    ElNotification({
      title: 'Error fields',
      message: e.response ? error.getError(e) : 'неверные данные',
      type: 'error',
      duration: 3000,
      active: true,
    })

    state.status = "invalidData";
  })
};

const logout = async () => {

  await axios({
    method: 'post',
    url: '/api/logout',
    data: {
      'token_id': 10
    }
  })
  state.user = null;

  state.User = null;
  state.authenticated = false;
  state.access_token = null;
  router.push("/login");
};

const getLang = (lang = 'ru') => {
  return new Promise(async (resolve) => {
    try {
      const response = await fetch(`/api/langue/${lang}`);

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || `HTTP error ${response.status}`);
      }

    resolve(data);
    } catch (e) {
      console.error('Language fetch error:', e);
      ElNotification({
        title: 'Language Error',
        message: errors.something_wrong,
        type: 'error',
      });
      resolve(false);
    }
  })

};

const getRole = (role_id) => {
  let role = state.roles.filter(el => el.value === role_id);

  return role.length ? role[0].label : 'hidden';
}

const getStatus = (id) => {
  let role = state.statuses.filter(el => el.value === id);

  return role.length ? role[0].label : 'hidden';
}

export default () => {
  return {
    authenticated: toRef(state, "authenticated"),
    loading: toRef(state, "loading"),
    status: toRef(state, "status"),
    access_token,
    User: toRef(state, "User"),
    roles: toRef(state, "roles"),
    me,
    getRole,
    getStatus,
    login,
    register,
    logout,
    update,
    getLang
  };
};
