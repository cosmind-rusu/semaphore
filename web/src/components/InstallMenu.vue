<template>
  <v-tabs show-arrows class="pl-4">

    <v-tab
      key="artifacts"
      :to="`/project/${projectId}`"
    >
      Artifacts
    </v-tab>

  </v-tabs>
</template>
<script>

export default {

  props: {
    projectId: Number,
    projectType: String,
    canUpdateProject: Boolean,
  },

  data() {
    return {
      id: null,
    };
  },

  computed: {
    isCloud() {
      return process.env.VUE_APP_BUILD_TYPE === 'cloud';
    },
  },

};
</script>
