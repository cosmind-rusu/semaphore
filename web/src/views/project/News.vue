<template xmlns:v-slot="http://www.w3.org/1999/XSL/Transform">
  <div>

    <v-toolbar flat>
      <v-app-bar-nav-icon @click="showDrawer()"></v-app-bar-nav-icon>
      <v-toolbar-title>
        {{ $t('dashboard') }}
        <SubscriptionLabel v-if="projectType === 'premium'" />
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
                  {{ post.text }}
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
import SubscriptionLabel from '@/components/SubscriptionLabel.vue';
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
    DashboardMenu, SubscriptionLabel,
  },

  data() {
    return {
      USER_PERMISSIONS,
      posts: [{
        image: 'https://opengraph.githubassets.com/144b5aed2d05ec90597b53a971560e5084a2bf58a3b79ad3c72e853e2fbc8bc7/semaphoreui/semaphore/releases/tag/v2.10.35',
        title: 'Project Runners',
        subtitle: 'Allow to join remote runner to specified project.',
        text: "I'm a thing. But, like most politicians, he promised more than he could deliver. "
          + "You won't have time for sleeping, soldier, not with all the bed making you'll be doing. "
          + "Then we'll go with that data file! Hey, you add a one and two zeros to that or we walk! "
          + "You're going to do his laundry? I've got to find a way to escape.",
        show: false,
      }],
    };
  },

};
</script>
