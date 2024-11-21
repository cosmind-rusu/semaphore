<template xmlns:v-slot="http://www.w3.org/1999/XSL/Transform">
  <div>

    <v-toolbar flat>
      <v-app-bar-nav-icon @click="showDrawer()"></v-app-bar-nav-icon>
      <v-toolbar-title>
        {{ projectType === 'premium' ? 'Manage Subscription' : $t('dashboard') }}
      </v-toolbar-title>
    </v-toolbar>

    <DashboardMenu :project-id="projectId" :project-type="projectType"
      :can-update-project="can(USER_PERMISSIONS.updateProject)" />

    <div class="px-6 pt-8" style="max-width: 1200px;">
      <v-row>
        <v-col md="6" v-for="post in posts" :key="post.title">
          <v-card>
            <v-img :src="post.image" height="200px"></v-img>

            <v-card-title>
              {{ post.title }}
            </v-card-title>

            <v-card-subtitle>
              {{ post.subtitle }}
            </v-card-subtitle>

            <v-card-actions>
              <v-btn text @click="post.show = !post.show">
                Read more
              </v-btn>

              <v-spacer></v-spacer>

              <v-btn icon @click="post.show = !post.show">
                <v-icon>{{ post.show ? 'mdi-chevron-up' : 'mdi-chevron-down' }}</v-icon>
              </v-btn>
            </v-card-actions>

            <v-expand-transition>
              <div v-show="post.show">
                <v-divider></v-divider>

                <v-card-text>
                  <div v-html="post.text"></div>
                </v-card-text>
              </div>
            </v-expand-transition>
          </v-card>
        </v-col>
      </v-row>
    </div>

  </div>
</template>
<script>
import DashboardMenu from '@/components/DashboardMenu.vue';
import PermissionsCheck from '@/components/PermissionsCheck';
import { USER_PERMISSIONS } from '@/lib/constants';

export default {
  props: {
    projectId: Number,
    projectType: String,
  },
  mixins: [PermissionsCheck],
  components: {
    DashboardMenu,
  },

  data() {
    return {
      USER_PERMISSIONS,
      posts: [{
        image: 'https://semaphoreui.com/uploads/v2.10.41/poster.webp',
        title: 'Project Runners',
        subtitle: 'Allow joining a remote runner to a specified project.',
        text: '<ol>'
          + '<li>Project Runners allow dedicated resources to be allocated to specific projects, avoiding contention and improving performance.</li>'
          + '<li>You can scale runners independently for each project based on workload, making it easier to manage high-demand projects without impacting others.</li>'
          + '<li>Administrators can assign and monitor runners at the project level, providing better control over execution environments and resource utilization.</li>'
          + '</ol>',
        show: false,
      }],
    };
  },

};
</script>
