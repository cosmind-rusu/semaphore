<template>
  <div style="height: 100%;">
    <v-toolbar flat>
      <v-app-bar-nav-icon style="color: white;" @click="showDrawer()"></v-app-bar-nav-icon>
      <v-toolbar-title>Support Center</v-toolbar-title>
    </v-toolbar>

    <iframe
      ref="iframe"
      @load="onLoad"
      class="crisp_chat"
      title="Support Desk"
      :src="`/helpdesk/projects/${projectId}/${pathMatch || 'dashboard'}`"/>
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
    pathMatch: String,
  },

  methods: {
    async trackNavigation() {
      const iframeWindow = this.$refs.iframe.contentWindow;
      const { pathname } = iframeWindow.location;

      const [, path] = pathname.match(/\/helpdesk\/projects\/\d+\/(.*)$/);

      await this.$router.push({
        path: `/project/${this.projectId}/helpdesk/${path}`,
      });
    },

    onLoad() {
      const iframeWindow = this.$refs.iframe.contentWindow;
      const iframeDocument = iframeWindow.document;

      iframeDocument.body.addEventListener('click', () => this.trackNavigation());
      iframeDocument.body.addEventListener('submit', () => this.trackNavigation());

      const originalPushState = iframeWindow.history.pushState;
      const self = this;
      iframeWindow.history.pushState = function (...args) {
        originalPushState.apply(this, args);
        self.trackNavigation();
      };

      const originalReplaceState = iframeWindow.history.replaceState;
      iframeWindow.history.replaceState = function (...args) {
        originalReplaceState.apply(this, args);
        self.trackNavigation();
      };

      // MutationObserver to detect DOM changes
      const observer = new MutationObserver((mutations) => {
        mutations.forEach((mutation) => {
          if (mutation.type === 'childList') {
            self.trackNavigation();
          }
        });
      });

      observer.observe(iframeDocument.body, { childList: true, subtree: true });
    },

    showDrawer() {
      EventBus.$emit('i-show-drawer');
    },
  },
};

</script>
