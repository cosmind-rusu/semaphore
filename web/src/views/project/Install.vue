<template xmlns:v-slot="http://www.w3.org/1999/XSL/Transform">
  <div v-if="project != null">
    <v-dialog
      :max-width="600"
      v-model="detailsDialog"
      persistent
      :transition="false"
      scrollable
    >
      <v-card>
        <v-card-title>
          {{ detailsVersion.semver }}
          <v-spacer></v-spacer>
          <v-btn icon @click="detailsDialog = false">
            <v-icon>mdi-close</v-icon>
          </v-btn>
        </v-card-title>
        <v-divider></v-divider>
        <v-card-text class="pa-0">
          <iframe
            title=""
            style="border: 0; display: block; padding: 0; margin: 0; width: 100%; height: 400px;"
            :src="`https://www.semui.co/releases/semaphore-v${(detailsVersion.semver || '').replaceAll('.', '_')}-premium/fragment`"
          ></iframe>
        </v-card-text>
      </v-card>
    </v-dialog>

    <v-dialog
      v-model="installationGuideDialog"
      max-width="800"
      persistent
      :transition="false"
    >
      <v-card>
        <v-card-title>
          Installation
          <v-spacer></v-spacer>
          <v-btn
            @click="installationGuideDialog = false"
            icon
          >
            <v-icon>mdi-close</v-icon>
          </v-btn>
        </v-card-title>

        <v-card-text class="text-xs-center">
          <pre style="white-space: pre-wrap">{{ installationCommand }}</pre>
        </v-card-text>

        <v-btn
          v-if="project.licenseKey"
          style="position: absolute; right: 20px; bottom: 20px;"
          icon
          @click="copyToClipboard(installationCommand)"
        >
          <v-icon>mdi-content-copy</v-icon>
        </v-btn>
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
      >Download
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
          class="justify-start"
        >
          <v-expansion-panel-header>
            <h2 style="flex: unset; width: 100px;">
              {{ version.semver }}
              <span
                style="display: block; font-size: 14px; font-weight: normal;"
              >{{ version.date }}</span>
            </h2>

            <v-btn
              style="max-width: 100px;"
              class="ml-4"
              color="primary"
              outlined
              @click="showDetails(version)"
            >
              Details
            </v-btn>
            <div class="pl-6">{{ version.description }}</div>

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
                      class="mr-3"
                      style="text-decoration: none !important;"
                      color="primary"
                      @click="showInstallationGuide(asset, version)"
                    >
                      Install
                    </v-btn>
                    <v-btn
                      v-if="asset.platform !== 'docker'"
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

function getInstallationFileName(assets, version) {
  return `semaphore_${version.semver}-premium_${assets.platform}_${assets.architecture[0]}.${assets.extension}`;
}

function getInstallationFileURL(assets, version) {
  return `https://www.semui.co/uploads/v${version.semver}-premium/${version.id}/${getInstallationFileName(assets, version)}`;
}

const EXTENSION_ICONS = {
  deb: {
    icon: 'debian',
    color: 'blue',
    guide(assets, version) {
      return `wget ${getInstallationFileURL(assets, version)}

sudo dpkg -i ${getInstallationFileName(assets, version)}`;
    },
  },
  rpm: {
    icon: 'centos',
    color: 'red',
    guide(assets, version) {
      return `wget ${getInstallationFileURL(assets, version)}

sudo yum install ${getInstallationFileName(assets, version)}`;
    },
  },
  zip: {
    icon: 'folder-zip',
    color: 'orange',
    guide(assets, version) {
      return `Invoke-WebRequest -Uri ("${getInstallationFileURL(assets, version)}") -OutFile semaphore.zip

Expand-Archive -Path semaphore.zip  -DestinationPath ./

./semaphore setup
./semaphore server --config=./config.json`;
    },
  },
  'tar.gz': {
    icon: 'archive',
    color: 'gray',
    guide(assets, version) {
      return `wget ${getInstallationFileURL(assets, version)}

tar xf ${getInstallationFileName(assets, version)}

./semaphore setup
./semaphore server --config=./config.json`;
    },
  },
  '': {
    icon: '',
    color: '',
    guide() {
      return '';
    },
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
      installationGuideDialog: null,
      installationGuideVersion: {},
      installationGuideAsset: {},

      detailsDialog: null,
      detailsVersion: {},

      versions: [{
        semver: '2.9.107',
        id: '17b51e42-3e81-467b-835f-105ade00e0f4',
        date: 'May 24, 2024',
        description: 'Fixed BoltDB database path in Docker.',
        dockerAssets: {
          terraform: '1.8.2',
          tofu: '1.7.0',
          pulumi: '3.116.1',
        },
      }, {
        semver: '2.9.101',
        id: 'be4b13f0-b5f4-4786-b34f-d03d5f6fff57',
        date: 'May 22, 2024',
        description: 'Added Enum and Secret survey variables.',
      }, {
        semver: '2.9.83',
        id: 'e67dec00-3bef-4ded-ba90-f72e110c07b0',
        date: 'May 8, 2024',
        description: 'Fixed Bash and Terraform in Docker image.',
      }, {
        semver: '2.9.76',
        id: '18b71336-157b-481c-9828-0433a3e8a48b',
        date: 'May 4, 2024',
        description: 'Fixed activation process.',
      }, {
        semver: '2.9.75',
        id: '2cc7896f-4d0f-4921-994c-e330d012605e',
        date: 'April 29, 2024',
        description: 'Fixed SSH using wrong username.',
      }, {
        semver: '2.9.73',
        id: '594d728f-5f07-4314-8f8e-23feb4ad4dfb',
        date: 'April 27, 2024',
        description: 'Terraform support bugfixes.',
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

  computed: {
    installationCommand() {
      if (!this.installationGuideAsset.platform) {
        return '';
      }

      if (this.installationGuideAsset.platform === 'docker') {
        return `docker run -p 3000:3000 --name semaphore \\
    -e SEMAPHORE_DB_DIALECT=bolt \\
    -e SEMAPHORE_ADMIN=admin \\
    -e SEMAPHORE_ADMIN_PASSWORD=changeme \\
    -e SEMAPHORE_ADMIN_NAME=Admin \\
    -e SEMAPHORE_ADMIN_EMAIL=admin@localhost \\
    -d semaphoreui/semaphore:v${this.installationGuideVersion.semver}-premium`;
      }

      if (!this.installationGuideAsset.extension) {
        return '';
      }

      const extInfo = EXTENSION_ICONS[this.installationGuideAsset.extension];

      return extInfo.guide(this.installationGuideAsset, this.installationGuideVersion);
    },
  },

  async created() {
    await this.refreshProject();
  },

  methods: {

    async copyToClipboard(text) {
      try {
        await window.navigator.clipboard.writeText(text);
        EventBus.$emit('i-snackbar', {
          color: 'success',
          text: 'The command has been copied to the clipboard.',
        });
      } catch (e) {
        EventBus.$emit('i-snackbar', {
          color: 'error',
          text: `Can't copy the command: ${e.message}`,
        });
      }
    },

    showDetails(version) {
      this.detailsVersion = version;
      this.detailsDialog = true;
      // eslint-disable-next-line no-restricted-globals
      event.preventDefault();
      // eslint-disable-next-line no-restricted-globals
      event.stopPropagation();
    },

    showInstallationGuide(asset, version) {
      this.installationGuideDialog = true;
      this.installationGuideVersion = version;
      this.installationGuideAsset = asset;
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
