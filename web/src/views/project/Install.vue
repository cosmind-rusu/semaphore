<template xmlns:v-slot="http://www.w3.org/1999/XSL/Transform">
  <div v-if="project != null">

    <v-dialog
      v-model="dockerGuideDialog"
      max-width="600"
      persistent
      :transition="false"
    >
      <v-card>
        <v-card-title>
          Docker usage
          <v-spacer></v-spacer>
          <v-btn
            @click="dockerGuideDialog = false"
            icon
          >
            <v-icon>mdi-close</v-icon>
          </v-btn>
        </v-card-title>
        <v-card-text class="text-xs-center pb-0">
          <pre>
docker run -p 3000:3000 --name semaphore \
    -e SEMAPHORE_DB_DIALECT=bolt \
    -e SEMAPHORE_ADMIN=admin \
    -e SEMAPHORE_ADMIN_PASSWORD=changeme \
    -e SEMAPHORE_ADMIN_NAME=Admin \
    -e SEMAPHORE_ADMIN_EMAIL=admin@localhost \
    -d semaphoreui/semaphore:v{{dockerGuideVersion.semver}}-premium
          </pre>
        </v-card-text>
      </v-card>
    </v-dialog>

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

    <div class="pa-4" v-if="project.plan === 'free'">
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
      <v-expansion-panels accordion focusable :value="0">
        <v-expansion-panel
          v-for="version of versions"
          :key="version.id"
        >
          <v-expansion-panel-header>
            <h2 class="pl-4 pt-1">
              {{ version.semver }}
              <v-chip v-if="version.latest" small class="ml-2 mb-1" color="success">Latest</v-chip>
            </h2>
          </v-expansion-panel-header>
          <v-expansion-panel-content>
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
                      class="mr-2"
                      :style="{
                    background: architecture === 'amd64' ? 'green' : 'lightblue',
                    color: architecture === 'amd64' ? 'white' : 'black',
                  }"
                      v-for="architecture in asset.architecture"
                      :key="architecture"
                    >{{ architecture }}</code>
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
                      v-if="asset.platform === 'docker'"
                      style="text-decoration: none !important;"
                      color="primary"
                      @click="showDockerGuide(asset, version)"
                    >
                      Install
                    </v-btn>
                    <v-btn
                      v-else
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
          </v-expansion-panel-content>
        </v-expansion-panel>
      </v-expansion-panels>
    </div>
  </div>
</template>
<style lang="scss">
</style>
<script>
import EventBus from '@/event-bus';
import axios from 'axios';

const PLATFORM_ICONS = {
  windows: {
    icon: 'microsoft-windows',
    color: 'blue',
  },
  docker: {
    icon: 'docker',
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
  snap: {
    icon: 'docker',
    color: 'blue',
  },
  '': {
    icon: '',
    color: '',
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
      project: null,
      deleteProjectDialog: null,
      dockerGuideDialog: null,
      dockerGuideVersion: {},

      versions: [{
        semver: '2.9.73',
        id: '594d728f-5f07-4314-8f8e-23feb4ad4dfb',
        latest: true,
      }, {
        semver: '2.9.72',
        id: '863ba8c7-f4ea-4921-8883-8af1ca254a6c',
      }, {
        semver: '2.9.63',
        id: '07238e2b-6cc1-422d-b1f9-0deb45cf4d93',
      }],
      assets: [{
        platform: 'linux',
        architecture: ['amd64'],
        extension: 'tar.gz',
      }, {
        platform: 'linux',
        architecture: ['arm64'],
        extension: 'tar.gz',
      }, {
        platform: 'linux',
        architecture: ['amd64'],
        extension: 'deb',
      }, {
        platform: 'linux',
        architecture: ['arm64'],
        extension: 'deb',
      }, {
        platform: 'linux',
        architecture: ['amd64'],
        extension: 'rpm',
      }, {
        platform: 'linux',
        architecture: ['arm64'],
        extension: 'rpm',
      }, {
        platform: 'docker',
        architecture: ['amd64', 'arm64'],
        extension: '',
      }, {
        platform: 'freebsd',
        architecture: ['amd64'],
        extension: 'tar.gz',
      }, {
        platform: 'freebsd',
        architecture: ['arm64'],
        extension: 'tar.gz',
      }, {
        platform: 'darwin',
        architecture: ['amd64'],
        extension: 'tar.gz',
      }, {
        platform: 'windows',
        architecture: ['amd64'],
        extension: 'zip',
      }],
    };
  },

  watch: {
    projectId() {
      this.refreshProject();
    },
  },

  async created() {
    await this.refreshProject();
  },

  methods: {
    showDockerGuide(asset, version) {
      this.dockerGuideDialog = true;
      this.dockerGuideVersion = version;
    },

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
    async refreshProject() {
      this.project = (await axios({
        method: 'get',
        url: `/billing/projects/${this.projectId}`,
        responseType: 'json',
      })).data;
    },

  },
};
</script>
