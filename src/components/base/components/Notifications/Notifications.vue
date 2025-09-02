<template>
  <section class="notifications navbar-notifications">
    <header class="cardHeader card-header">
      <div class="text-center mb-sm">
        <strong>You have 13 notifications</strong>
      </div>
    </header>
    <NewNotifictionsList v-if="newNotifications" />
    <NotifictionsList v-else-if="notificationsTabSelected === 1" />

    <footer class="cardFooter text-sm card-footer">
      <span class="fs-mini">Synced at: 21 Apr 2014 18:36</span>
      <el-button
        variant="link"
        @click="loadNotifications"
        :class="{disabled: isLoad, 'btnNotificationsReload btn-xs float-end py-0': true}"
      >
        <span v-if="isLoad"><i class="la la-refresh la-spin" /> Loading...</span>
        <i v-else class="la la-refresh" />
      </el-button>
    </footer>
  </section>
</template>

<script setup>
import { ref } from 'vue';
import NotifictionsList from './NotificationsDemo/NotificationsList';
import NewNotifictionsList from './NotificationsDemo/NewNotificationsList';
import Messages from './NotificationsDemo/Messages';

const notificationsTabSelected = ref(1);
const newNotifications = ref(null);
const isLoad = ref(false);

const changeNotificationsTab = (tab) => {
  notificationsTabSelected.value = tab;
  newNotifications.value = null;
};

const loadNotifications = () => {
  isLoad.value = true;
  
  setTimeout(() => {
    newNotifications.value = 'new notifications component';
    isLoad.value = false;
  }, 1500);
};
</script>

<style src="./Notifications.scss" lang="scss" />
