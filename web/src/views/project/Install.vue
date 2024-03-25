<template xmlns:v-slot="http://www.w3.org/1999/XSL/Transform">
  <div>

    <v-toolbar flat>
      <v-app-bar-nav-icon @click="showDrawer()"></v-app-bar-nav-icon>
      <v-toolbar-title>{{ $t('dashboard') }}</v-toolbar-title>
    </v-toolbar>

    <v-tabs show-arrows class="pl-4">
      <v-tab
        v-if="projectType === 'premium'"
        key="install"
        :to="`/project/${projectId}/install`"
      >Setup
      </v-tab>

      <v-tab
        v-if="projectType === ''"
        key="history"
        :to="`/project/${projectId}/history`"
      >{{ $t('history') }}
      </v-tab>
      <v-tab key="activity" :to="`/project/${projectId}/activity`">{{ $t('activity') }}</v-tab>
      <v-tab key="settings" :to="`/project/${projectId}/settings`">{{ $t('settings') }}</v-tab>
      <v-tab
        key="billing"
        :to="`/project/${projectId}/billing`"
      >Billing
      </v-tab>
    </v-tabs>

    <div class="pa-4" v-if="false">
      <v-alert
        dense
        text
        type="warning"
        style="display: inline-block"
      >
        You don't have an activated plan.
      </v-alert>
    </div>
    <div class="pa-4" v-else>
      <div v-for="version of versions" :key="version.id">

        <h2 class="pl-4">
          {{ version.semver }}
          <v-chip small class="ml-2 mb-1" color="error">Premium</v-chip>
        </h2>

        <v-simple-table>
          <template v-slot:default>
            <thead>
            <tr>
              <th class="text-left">
                Platform
              </th>
              <th class="text-left">
                Architecture
              </th>
              <th class="text-left">
                Type
              </th>
              <th class="text-left">
              </th>
            </tr>
            </thead>
            <tbody>
            <tr
              v-for="asset in assets"
              :key="getAssetUrl(asset, version)"
            >
              <td>
                <v-icon
                  :color="PLATFORM_ICONS[asset.platform].color"
                  class="mr-2"
                >
                  mdi-{{ PLATFORM_ICONS[asset.platform].icon }}
                </v-icon>
                {{ asset.platform }}
              </td>
              <td>
                <code
                  :style="{
                    background: asset.architecture === 'amd64' ? 'green' : 'lightblue',
                    color: asset.architecture === 'amd64' ? 'white' : 'black',
                  }"
                >{{ asset.architecture }}</code>
              </td>
              <td>
                <v-icon
                  :color="EXTENSION_ICONS[asset.extension].color"
                  class="mr-2"
                >
                  mdi-{{ EXTENSION_ICONS[asset.extension].icon }}
                </v-icon>
                {{ asset.extension }}
              </td>
              <td>
                <v-btn
                  style="text-decoration: none !important;"
                  color="primary"
                  :href="getAssetUrl(asset, version)"
                >
                  Download
                </v-btn>
              </td>
            </tr>
            </tbody>
          </template>
        </v-simple-table>
      </div>
    </div>
  </div>
</template>
<style lang="scss">
</style>
<script>
import EventBus from '@/event-bus';

const PLATFORM_ICONS = {
  windows: {
    icon: 'microsoft-windows',
    color: 'blue',
  },
  freebsd: {
    icon: 'freebsd',
    color: 'red',
  },
  darwin: {
    icon: 'apple',
    color: 'grey',
  },
  linux: {
    icon: 'linux',
    color: 'black',
  },
};

const EXTENSION_ICONS = {
  deb: {
    icon: 'debian',
    color: 'blue',
  },
  rpm: {
    icon: 'centos',
    color: 'red',
  },
  zip: {
    icon: 'folder-zip',
    color: 'orange',
  },
  'tar.gz': {
    icon: 'archive',
    color: 'gray',
  },
};

export default {
  components: {},
  props: {
    projectId: Number,
    projectType: String,
  },

  data() {
    return {
      PLATFORM_ICONS,
      EXTENSION_ICONS,
      deleteProjectDialog: null,
      versions: [{
        semver: '2.9.63',
        id: '07238e2b-6cc1-422d-b1f9-0deb45cf4d93',
      }],
      assets: [{
        platform: 'linux',
        architecture: 'amd64',
        extension: 'tar.gz',
      }, {
        platform: 'linux',
        architecture: 'arm64',
        extension: 'tar.gz',
      }, {
        platform: 'linux',
        architecture: 'amd64',
        extension: 'deb',
      }, {
        platform: 'linux',
        architecture: 'arm64',
        extension: 'deb',
      }, {
        platform: 'linux',
        architecture: 'amd64',
        extension: 'rpm',
      }, {
        platform: 'linux',
        architecture: 'arm64',
        extension: 'rpm',
      }, {
        platform: 'freebsd',
        architecture: 'amd64',
        extension: 'tar.gz',
      }, {
        platform: 'freebsd',
        architecture: 'arm64',
        extension: 'tar.gz',
      }, {
        platform: 'darwin',
        architecture: 'amd64',
        extension: 'tar.gz',
      }, {
        platform: 'windows',
        architecture: 'amd64',
        extension: 'zip',
      }],
    };
  },

  methods: {
    getAssetUrl(asset, version) {
      return `https://www.semui.co/uploads/v${version.semver}-premium/${version.id}/semaphore_${version.semver}-premium_${asset.platform}_${asset.architecture}.${asset.extension}`;
    },

    showDrawer() {
      EventBus.$emit('i-show-drawer');
    },

    onError(e) {
      EventBus.$emit('i-snackbar', {
        color: 'error',
        text: e.message,
      });
    },

  },
};
</script>
