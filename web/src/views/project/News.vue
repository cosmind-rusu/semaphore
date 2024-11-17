<template xmlns:v-slot="http://www.w3.org/1999/XSL/Transform">
    <div>

        <v-toolbar flat>
            <v-app-bar-nav-icon @click="showDrawer()"></v-app-bar-nav-icon>
            <v-toolbar-title>
                {{ $t('dashboard') }}
                <SubscriptionLabel v-if="projectType === 'premium'" />
            </v-toolbar-title>
        </v-toolbar>

        <DashboardMenu
            :project-id="projectId"
            :project-type="projectType"
            :can-update-project="can(USER_PERMISSIONS.updateProject)"
        />

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
    };
  },

};
</script>
