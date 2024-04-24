<template>
  <div>
    <v-toolbar flat>
      <v-app-bar-nav-icon style="color: white;" @click="showDrawer()"></v-app-bar-nav-icon>
      <v-toolbar-title>Service Desk</v-toolbar-title>
    </v-toolbar>

    <iframe
      class="crisp_chat"
      title="Support Desk"
      :src="`/helpdesk/projects/${projectId}/admin/panel/tickets/my-tickets`" />
  </div>

</template>
<style>
.crisp_chat {
  position: absolute;
  width: 100%;
  height: calc(100% - 64px);
  bottom: 0;
  display: block;
  border: 0;
}
</style>
<script>

import EventBus from '@/event-bus';

export default {
  created() {
    window.$crisp.push(['do', 'chat:hide']);
  },

  props: {
    projectId: Number,
  },

  methods: {
    newTicket() {

    },
    showDrawer() {
      EventBus.$emit('i-show-drawer');
    },
  },
};

</script>
