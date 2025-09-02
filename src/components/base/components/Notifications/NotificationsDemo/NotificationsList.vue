<template>
  <ul class="listGroup thin-scroll">
    <li class="list-group-item listGroupItem d-flex align-items-start" v-for="(item, index) in notifications" :key="index">
      <span class="notificationIcon thumb-sm mr-3">
        <template v-if="item.img">
          <img class="rounded-circle" :src="item.img" :alt="item.alt" />
        </template>
        <template v-else>
          <i :class="item.icon" class="fa-lg"></i>
        </template>
      </span>
      <div class="">
        <p class="m-0 overflow-hidden">
          <span v-html="item.text"></span>
          <time class="d-block text-muted small">
            {{ item.time }}
          </time>
        </p>
        <div v-if="item.actions" class="mt-2">
          <el-button size="small" type="success" class="me-1">Allow</el-button>
          <el-button size="small" type="danger">Deny</el-button>
        </div>
      </div>
    </li>
  </ul>
</template>

<script setup>
import { ref } from 'vue';

const notifications = ref([
  {
    img: '/assets/people/a3.jpg',
    alt: 'woman',
    text: '1 new user just signed up! Check out <a href="#">Monica Smith</a>\'s account.',
    time: '2 mins ago'
  },
  {
    img: '/assets/people/a5.jpg',
    alt: 'man',
    text: 'User <a href="#">Jeff</a> registered',
    time: '12:18 AM',
    actions: true
  },
  {
    icon: 'fa fa-shield',
    text: 'Instructions for changing your Envato Account password. Please check your <a href="#">security page</a>.',
    time: '12:18 AM'
  },
  {
    icon: 'fa fa-facebook text-white',
    text: 'New <strong>76</strong> Facebook likes received.',
    time: '15 Apr 2014',
    bg: 'bg-primary rounded-lg'
  },
  {
    icon: 'fa fa-circle-o text-white',
    text: 'Dark matter detected.',
    time: '15 Apr 2014',
    bg: 'bg-gray-dark circle circle-lg'
  }
]);
</script>

<style src="./ListGroup.scss" lang="scss" />